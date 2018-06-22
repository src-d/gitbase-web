import React, { Component } from 'react';
import PropTypes from 'prop-types';
import SplitPane from 'react-split-pane';
import Schema from './Schema';
import SampleQueries from './SampleQueries';
import './Sidebar.less';
import CollapseIcon from '../icons/collapse-left-column.svg';
import LinkIcon from '../icons/links.svg';

class Sidebar extends Component {
  constructor(props) {
    super(props);

    this.state = { collapsed: false };

    this.handleToggle = this.handleToggle.bind(this);
  }

  handleToggle() {
    this.setState({ collapsed: !this.state.collapsed });
  }

  link(text, url) {
    return (
      <div>
        <a href={url} target="_blank" rel="noopener noreferrer">
          <LinkIcon className="small-icon" />
          {text}
        </a>
      </div>
    );
  }

  render() {
    const { schema, onTableClick, onExampleClick, exampleQueries } = this.props;
    const { collapsed } = this.state;

    return (
      <div className={`sidebar ${collapsed ? 'collapsed' : ''}`}>
        <div className="header">
          <h3>gitbase playgroun{'{d}'}</h3>
          <CollapseIcon className="big-icon" onClick={this.handleToggle} />
        </div>
        <div className="main">
          <SplitPane
            split="horizontal"
            defaultSize={300}
            minSize={0}
            maxSize={-15}
          >
            <Schema schema={schema} onTableClick={onTableClick} />
            <SampleQueries
              onExampleClick={onExampleClick}
              exampleQueries={exampleQueries}
            />
          </SplitPane>
        </div>
        <div className="footer list">
          {this.link('go-git', 'https://github.com/src-d/go-git')}
          {this.link('babelfish', 'https://doc.bblf.sh')}
          {this.link('source{d} © 2018', 'https://sourced.tech')}
        </div>
      </div>
    );
  }
}

Sidebar.propTypes = {
  schema: Schema.propTypes.schema,
  onTableClick: PropTypes.func,
  onExampleClick: PropTypes.func,
  exampleQueries: SampleQueries.propTypes.exampleQueries
};

export default Sidebar;
