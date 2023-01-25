module.exports = {
  'env': {
    'node': true,
    'es2021': true,
  },
  'extends': [
    'eslint:recommended', 
    'plugin:unicorn/all'
  ],
  'overrides': [
  ],
  'parserOptions': {
    'ecmaVersion': 'latest',
    'sourceType': 'module',
    'allowImportExportEverywhere': true,
  },
  'ignorePatterns': [
    'node_modules',
    '*.lock',
  ],
  'rules': {
    'indent': [
      'error',
      2
    ],
    'linebreak-style': [
      'error',
      'windows'
    ],
    'quotes': [
      'error',
      'single',
      {
        'avoidEscape': true,
        'allowTemplateLiterals': true
      }
    ],
    'semi': [
      'error',
      'always'
    ],
    'unicorn/no-keyword-prefix': ['off'],
    'unicorn/no-process-exit': ['off']
  }
};
