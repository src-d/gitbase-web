const CracoLessPlugin = require('craco-less');
const { ESLINT_MODES } = require('@craco/craco');

module.exports = {
  plugins: [{ plugin: CracoLessPlugin }],
  eslint: {
    mode: ESLINT_MODES.file
  }
};
