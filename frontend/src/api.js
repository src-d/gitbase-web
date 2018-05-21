const envVars = window.gitbasepg || { SERVER_URL: '', SELECT_LIMIT: 100 };

const serverUrl = envVars.SERVER_URL;
const selectLimit = envVars.SELECT_LIMIT;

const apiUrl = url => `${serverUrl}${url}`;

function checkStatus(resp) {
  if (resp.status < 200 || resp.status >= 300) {
    return resp
      .json()
      .catch(() => {
        throw new Error(resp.statusText);
      })
      .then(json => {
        if (json.errors) {
          throw json.errors;
        }
        throw new Error(resp.statusText);
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

function query(sql) {
  return apiCall(`/query`, {
    method: 'POST',
    body: {
      query: sql,
      limit: selectLimit
    }
  });
}

function tables() {
  return apiCall(`/tables`);
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
  return tables()
    .then(res =>
      Promise.all(
        res.data.map(e =>
          query(`DESCRIBE TABLE ${e.table}`).then(tableRes => ({
            table: e.table,
            columns: tableRes.data
          }))
        )
      )
    )
    .catch(err => Promise.reject(normalizeErrors(err)));
}

export default {
  query,
  tables,
  schema
};
