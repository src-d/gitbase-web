import React, { Component } from 'react';
import PropTypes from 'prop-types';
import SplitPane from 'react-split-pane';
import { Glyphicon } from 'react-bootstrap';
import Schema from './Schema';
import SampleQueries from './SampleQueries';
import './Sidebar.less';

class Sidebar extends Component {
  constructor(props) {
    super(props);

    this.state = { collapsed: false };

    this.handleToggle = this.handleToggle.bind(this);
  }

  handleToggle() {
    this.setState({ collapsed: !this.state.collapsed });
  }

  render() {
    const { schema, onTableClick, onExampleClick } = this.props;
    const { collapsed } = this.state;
    const togglerIcon = collapsed ? 'chevron-right' : 'chevron-left';

    return (
      <div className={`sidebar ${collapsed ? 'collapsed' : ''}`}>
        <div className="header">
          <h3>{'{d}'} Gitbase Playground </h3>
          <Glyphicon onClick={this.handleToggle} glyph={togglerIcon} />
        </div>
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
          <a
            href="https://doc.bblf.sh"
            target="_blank"
            rel="noopener noreferrer"
          >
            babelfish
          </a>
        </div>
      </div>
    );
  }
}

Sidebar.propTypes = {
  schema: Schema.propTypes.schema,
  onTableClick: PropTypes.func,
  onExampleClick: PropTypes.func
};

export default Sidebar;
