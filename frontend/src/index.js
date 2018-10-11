import React from 'react';
import ReactDOM from 'react-dom';
import 'bootstrap/dist/css/bootstrap.css';
import './variables.less';
import App from './App';

ReactDOM.render(<App />, document.getElementById('root'));

if (!window.AbortController) {
  // eslint-disable-next-line no-console
  console.warn(`Cancelling gitbase queries is disabled; your browser does not support the AbortController interface.
Please check for compatibility here: https://developer.mozilla.org/en-US/docs/Web/API/AbortController#Browser_compatibility`);
}
