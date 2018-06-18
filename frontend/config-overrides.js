const { compose } = require('react-app-rewired');
const rewireLess = require('react-app-rewire-less');
const rewireSvgReactLoader = require('react-app-rewire-svg-react-loader');

module.exports = compose(rewireLess, rewireSvgReactLoader);
