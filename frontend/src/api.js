const envVars = window.gitbasepg || { SERVER_URL: '', SELECT_LIMIT: 100 };

const serverUrl = envVars.SERVER_URL;
const selectLimit = envVars.SELECT_LIMIT;

// if serverUrl is unset, this replaces the leading / in order to work behind proxies with a path other than /.
const apiUrl = url => `${serverUrl}${url}`.replace(/^\/+/g, '');

function statusError(resp) {
  return new Error(resp.statusText || `${resp.status} Error`);
}

function checkStatus(resp) {
  if (resp.status < 200 || resp.status >= 300) {
    return resp
      .json()
      .catch(() => {
        throw statusError(resp);
      })
      .then(json => {
        if (json.errors) {
          throw json.errors;
        }
        throw statusError(resp);
      });
  }
  return resp;
}

function normalizeError(err) {
  if (typeof err === 'object') {
    // error from server
    if (err.title) {
      return err.title;
    }
    // fetch abort
    if (err.name === 'AbortError') {
      return 'The user aborted the request.';
    }
    // javascript error
    if (err.message) {
      return err.message;
    }
    // weird object as error, shouldn't really happen
    return JSON.stringify(err);
  }
  if (typeof err === 'string') {
    return err;
  }
  return 'Internal error';
}

function normalizeErrors(err) {
  if (Array.isArray(err)) {
    return err.map(e => normalizeError(e));
  }
  return [normalizeError(err)];
}

function apiCall(url, options = {}) {
  const fetchOptions = {
    ...options,
    headers: {
      ...options.headers
    }
  };

  if (options.body) {
    if (options.formData) {
      fetchOptions.body = options.body;
    } else {
      fetchOptions.body = JSON.stringify(options.body);
      fetchOptions.headers['Content-Type'] = 'application/json';
    }
  }

  return fetch(apiUrl(url), fetchOptions)
    .then(checkStatus)
    .then(resp => resp.json())
    .then(json => {
      if (json.errors) {
        throw json.errors;
      }
      return json;
    })
    .catch(err => Promise.reject(normalizeErrors(err)));
}

function query(sql, signal) {
  const startTime = new Date();
  return apiCall(`/query`, {
    method: 'POST',
    body: {
      query: sql,
      limit: selectLimit
    },
    signal
  }).then(res => {
    res.meta.elapsedTime = new Date() - startTime;
    return res;
  });
}

function queryExport(sql) {
  const rawUrl = apiUrl('/export');
  const params = new URLSearchParams();
  params.append('query', sql);
  const url = `${rawUrl}?${params.toString()}`;
  return url;
}

/* Returns an array in the form:
[
  {
    "table": "refs",
    "columns": [
      {
        "name": "repository_id",
        "type": "TEXT"
      },
      ...
    ]
  },
  ...
]
*/
function schema() {
  return apiCall(`/schema`).then(res => res.data);
}

function detectLang(content, filename) {
  return apiCall('/detect-lang', {
    method: 'POST',
    body: {
      content,
      filename
    }
  }).then(res => res.data);
}

const defaultUastMode = 'semantic';
const uastModes = ['semantic', 'annotated', 'native'];

function parseCode(language, content, mode, filter, customServerUrl) {
  return apiCall('/parse', {
    method: 'POST',
    body: {
      language,
      content,
      mode,
      filter,
      serverUrl: customServerUrl
    }
  }).then(res => {
    if (res.status !== 200) {
      throw normalizeErrors(res.data.errors);
    }
    return res.data.uast;
  });
}

function getLanguages() {
  return apiCall(`/get-languages`).then(res => res.data);
}

function filterUAST(protobufs, filter) {
  return apiCall('/filter', {
    method: 'POST',
    body: {
      protobufs,
      filter
    }
  }).then(res => res.data);
}

function version() {
  return apiCall(`/version`).then(res => res.data);
}

export default {
  query,
  schema,
  queryExport,
  detectLang,
  parseCode,
  getLanguages,
  filterUAST,
  version,
  uastModes,
  defaultUastMode
};
