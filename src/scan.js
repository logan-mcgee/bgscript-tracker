import core from '@actions/core';

import fs from 'node:fs';
import path from 'node:path';

import { fetch } from 'undici';

import crypto from 'node:crypto';

const RPF_NAME = 'bg_ng_2802_0';
const RPF_URL = `http://prod.cloud.rockstargames.com/titles/gta5/pcros/bgscripts/${RPF_NAME}.rpf`;

async function main() {
  if (!fs.existsSync('bgscripts')) fs.mkdirSync('bgscripts');
  const rpfPath = path.join('bgscripts', RPF_NAME);
  if (!fs.existsSync(rpfPath)) fs.mkdirSync(rpfPath);
  if (!fs.existsSync(path.join(rpfPath, 'last-modified.txt'))) fs.writeFileSync(path.join(rpfPath, 'last-modified.txt'), 'Thu, 01 Jan 1970 00:00:00 GMT');

  const response = await fetch(RPF_URL);

  if (response.status !== 200) {
    core.setFailed(`Failed to fetch bgscript: ${response.status} ${response.statusText}`);
    return;
  }

  if (response.headers.get('content-type') !== 'application/octet-stream') {
    core.setFailed(`Invalid content-type: ${response.headers.get('content-type')}`);
    return;
  }

  const lastModified = response.headers.get('last-modified');
  const lastModifiedDate = new Date(lastModified);

  if (lastModifiedDate.toString() === 'Invalid Date') {
    core.setFailed(`Invalid last-modified date: ${lastModified}`);
    return;
  }

  const lastModifiedFile = path.join(rpfPath, 'last-modified.txt');
  const lastModifiedFileDate = new Date(fs.readFileSync(lastModifiedFile, 'utf8'));

  if (lastModifiedFileDate.toString() === 'Invalid Date') {
    core.setFailed(`Invalid last-modified date in file: ${lastModifiedFileDate}`);
    return;
  }

  if (lastModifiedDate.getTime() === lastModifiedFileDate.getTime()) {
    core.info('bgscript is up to date');
    process.exit(0);
    return;
  }

  core.info('New bgscript available');

  // get content as buffer
  const buffer = Buffer.from(await response.arrayBuffer());

  const fileNameDate = lastModifiedDate.toISOString().replaceAll(':', '-');
  const fileName = `${fileNameDate}-${crypto.createHash('sha256').update(buffer).digest('hex')}.rpf`;

  fs.writeFileSync(path.join(rpfPath, fileName), buffer);

  fs.writeFileSync(lastModifiedFile, lastModified);
}

// eslint-disable-next-line unicorn/prefer-top-level-await
main();