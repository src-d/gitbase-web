const envVars = window.gitbasepg || { SERVER_URL: '', SELECT_LIMIT: 100 };

const serverUrl = envVars.SERVER_URL || window.location.origin;
const selectLimit = envVars.SELECT_LIMIT;

const apiUrl = url => `${serverUrl}${url}`;

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
  const url = new URL(apiUrl('/export'));
  url.searchParams.append('query', sql);
  return url.toString();
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

function parseCode(language, content, filter, customServerUrl) {
  return apiCall('/parse', {
    method: 'POST',
    body: {
      language,
      content,
      filter,
      serverUrl: customServerUrl
    }
  }).then(res => {
    if (res.data.status !== 0) {
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

export default {
  query,
  schema,
  queryExport,
  detectLang,
  parseCode,
  getLanguages,
  filterUAST
};
