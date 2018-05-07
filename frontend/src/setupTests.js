import path from 'path';
import os from 'os';

const { LocalStorage } = require('node-localstorage');

global.fetch = require('jest-fetch-mock');

global.localStorage = new LocalStorage(
  path.join(os.tmpdir(), 'node-localstorage')
);

global.window = document.defaultView;
global.window.localStorage = global.localStorage;
