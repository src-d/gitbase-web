import React from 'react';
import { Glyphicon } from 'react-bootstrap';
import './SampleQueries.less';

const queries = [
  {
    name: 'load all java files',
    sql: 'select * from commits'
  },
  {
    name: 'get uast from code',
    sql: 'select * from commits'
  },
  {
    name: 'top 50 repositories by something very long string',
    sql: 'select * from commits'
  }
];

function SampleQueries() {
  return (
    <div className="sample-queries">
      <div className="title">Sample Queries</div>
      <div className="list">
        {queries.map((q, i) => (
          <div key={i} className="query" title={q.name}>
            <Glyphicon glyph="list" />
            {q.name}
          </div>
        ))}
      </div>
    </div>
  );
}

export default SampleQueries;
