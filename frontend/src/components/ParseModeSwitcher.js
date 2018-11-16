import React from 'react';
import PropTypes from 'prop-types';
import api from '../api';

function ParseModeSwitcher({ mode, handleModeChange }) {
  return (
    <div style={{ padding: '20px' }}>
      {api.uastModes.map(m => (
        <label
          key={m}
          style={{
            marginRight: '20px',
            textTransform: 'capitalize'
          }}
        >
          <input
            type="radio"
            value={m}
            checked={mode === m}
            onChange={e => handleModeChange(e.target.value)}
          />{' '}
          {m}
        </label>
      ))}
    </div>
  );
}

ParseModeSwitcher.propTypes = {
  mode: PropTypes.string.isRequired,
  handleModeChange: PropTypes.func.isRequired
};

export default ParseModeSwitcher;
