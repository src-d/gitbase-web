module.exports = {
  extends: [
    'airbnb-base',
    'plugin:react/recommended',
    'prettier',
    'prettier/react',
  ],
  env: {
    browser: true,
    es6: true,
    node: true,
    'jest/globals': true,
  },
  plugins: ['import', 'react', 'jest', 'prettier'],
  rules: {
    'prettier/prettier': 'error',
    'import/no-extraneous-dependencies': ['error', { 'devDependencies': ['**/*.test.js', 'src/setupTests.js'] }],
    'import/no-unresolved': 0,
    'import/extensions': 0,
    'func-names': 0,
    'no-plusplus': 0,
    'class-methods-use-this': 0, // strange rule. It doesn't allow to create method render() without this
    'no-case-declarations': 0, // otherwise code is very ugly
  },
};
