package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

const STARTING_GAMEVERSION = 3258
const GAMEVERSION_RANGE = 150

const SUBID_RANGE = 30

var GEN_GAMEVERSIONS = map[int]int{
	8: 3179,
	9: 554,
}

var PLATFORMS = map[string]int{
	"pcros":   8,
	"ps4":     8,
	"xboxone": 8,

	"ps5":    9,
	"xboxsx": 9,
}

const RPF_URL = "http://prod.cloud.rockstargames.com/titles/gta5/%s/bgscripts/bg_ng_%d_%d.rpf"

var ErrNotFound = errors.New("rpf not found")

func downloadRPF(platform string, gameVersion int, subId int) ([]byte, string, error) {
	//fmt.Printf("Attempting to download %s %d %d\n", platform, gameVersion, subId)

	resp, err := http.Get(fmt.Sprintf(RPF_URL, platform, gameVersion, subId))
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, "", ErrNotFound
	} else if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("expected 200, got %d", resp.StatusCode)
	}

	lastModified := resp.Header.Get("last-modified")
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return out, lastModified, nil
}

func hashBytes(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

func dirExists(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func main() {
	bwg := NewBoundedWaitGroup(100)

	for platform, gen := range PLATFORMS {
		fmt.Printf("Doing platform %s, gen %d\n", platform, gen)
		for gameVersion := GEN_GAMEVERSIONS[gen]; gameVersion < GEN_GAMEVERSIONS[gen]+GAMEVERSION_RANGE; gameVersion++ {
			for subId := 0; subId < SUBID_RANGE; subId++ {

				bwg.Add(1)

				go func(p string, g, gV, s int) {
					defer bwg.Done()

					file, lastModified, err := downloadRPF(p, gV, s)
					if err != nil {
						if !errors.Is(err, ErrNotFound) {
							fmt.Printf("Failed to download rpf for %s %d %d: %v\n", p, gV, s, err)
						}

						return
					}

					filePath := path.Join("bgscripts", fmt.Sprintf("gen%d", g), fmt.Sprintf("bg_ng_%d_%d", gV, s), p)
					if !dirExists(filePath) {
						os.MkdirAll(filePath, 0755)
					}

					fileHash := hashBytes(file)

					lastModifiedTime, err := time.Parse(time.RFC1123, lastModified)
					if err != nil {
						fmt.Printf("Failed to parse Last-Modified header: %v\n", err)
						return
					}
					newTime := lastModifiedTime.Unix()

					filePath = path.Join(filePath, fmt.Sprintf("%d_%s.rpf", newTime, fileHash))

					os.WriteFile(filePath, file, 0755)
				}(platform, gen, gameVersion, subId)

			}
		}

		bwg.Wait()
	}
}
