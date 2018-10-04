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

  link(text, url, version) {
    return (
      <div>
        <a href={url} target="_blank" rel="noopener noreferrer">
          <LinkIcon className="small-icon" />
          {text}
        </a>
        {version && <span className="version">{version}</span>}
      </div>
    );
  }

  render() {
    const {
      schema,
      version,
      onTableClick,
      onExampleClick,
      exampleQueries
    } = this.props;
    const { collapsed } = this.state;

    return (
      <div className={`sidebar ${collapsed ? 'collapsed' : ''}`}>
        <div className="header">
          <h3>gitbase web</h3>
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
          {this.link(
            'babelfish',
            'https://doc.bblf.sh',
            version ? version.bblfsh : ''
          )}
          {this.link(
            'gitbase',
            'https://docs.sourced.tech/gitbase',
            version ? version.gitbase : ''
          )}
          {this.link(
            'gitbase-web',
            'https://github.com/src-d/gitbase-web',
            version ? version.version : ''
          )}
          {this.link('source{d} © 2018', 'https://sourced.tech')}
        </div>
      </div>
    );
  }
}

Sidebar.propTypes = {
  schema: Schema.propTypes.schema,
  version: PropTypes.shape({
    version: PropTypes.string.isRequired,
    bblfsh: PropTypes.string.isRequired,
    gitbase: PropTypes.string.isRequired
  }),
  onTableClick: PropTypes.func,
  onExampleClick: PropTypes.func,
  exampleQueries: SampleQueries.propTypes.exampleQueries
};

export default Sidebar;
