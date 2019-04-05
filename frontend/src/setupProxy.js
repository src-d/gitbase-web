const proxy = require('http-proxy-middleware');

module.exports = function(app) {
  [
    '/query',
    '/export',
    '/schema',
    '/detect-lang',
    '/parse',
    '/get-languages',
    '/filter',
    '/version'
  ].forEach(path => app.use(proxy(path, { target: 'http://localhost:8080' })));
};
