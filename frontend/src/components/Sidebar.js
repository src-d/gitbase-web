import React from 'react';
import PropTypes from 'prop-types';
import SplitPane from 'react-split-pane';
import Schema from './Schema';
import SampleQueries from './SampleQueries';
import './Sidebar.less';

function Sidebar({ schema, onTableClick, onExampleClick }) {
  return (
    <div className="sidebar">
      <h3 className="header">{'{d}'} Gitbase Playground</h3>
      <div className="main">
        <SplitPane split="horizontal" defaultSize={200} minSize={100}>
          <Schema schema={schema} onTableClick={onTableClick} />
          <SampleQueries onExampleClick={onExampleClick} />
        </SplitPane>
      </div>
      <div className="footer">
        <a
          href="https://sourced.tech"
          target="_blank"
          rel="noopener noreferrer"
        >
          source{'{d}'}
        </a>
        <a
          href="https://github.com/src-d/gitbase"
          target="_blank"
          rel="noopener noreferrer"
        >
          gitbase
        </a>
        <a
          href="https://github.com/src-d/go-git"
          target="_blank"
          rel="noopener noreferrer"
        >
          go-git
        </a>
        <a href="https://doc.bblf.sh" target="_blank" rel="noopener noreferrer">
          babelfish
        </a>
      </div>
    </div>
  );
}

Sidebar.propTypes = {
  schema: Schema.propTypes.schema,
  onTableClick: PropTypes.func,
  onExampleClick: PropTypes.func
};

export default Sidebar;
