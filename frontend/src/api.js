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

  return fetch(url, fetchOptions)
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
    body: { query: sql }
  });
}

export default {
  query
};
