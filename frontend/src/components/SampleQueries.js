import React from 'react';
import PropTypes from 'prop-types';
import './SampleQueries.less';
import ExampleIcon from '../icons/example-query.svg';

function SampleQueries({ onExampleClick, exampleQueries }) {
  return (
    <div className="sample-queries">
      <div className="title">Sample Queries</div>
      <div className="list">
        {exampleQueries.map((q, i) => (
          <div
            key={i}
            className="query"
            title={q.name}
            onClick={() => onExampleClick(q.sql)}
          >
            <ExampleIcon className="small-icon" />
            {q.name}
          </div>
        ))}
      </div>
    </div>
  );
}

SampleQueries.propTypes = {
  onExampleClick: PropTypes.func,
  exampleQueries: PropTypes.arrayOf(
    PropTypes.shape({
      name: PropTypes.string.isRequired,
      sql: PropTypes.string.isRequired
    }).isRequired
  )
};

export default SampleQueries;
