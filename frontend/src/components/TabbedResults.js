import React, { Component } from 'react';
import { Row, Col, Alert, Tabs, Tab, Button, Glyphicon } from 'react-bootstrap';
import PropTypes from 'prop-types';
import ResultsTable from './ResultsTable';
import HistoryTable from './HistoryTable';
import Loader from './Loader';
import './TabbedResults.less';

class TabTitle extends Component {
  constructor(props) {
    super(props);

    this.state = {
      inEdit: false,
      title: props.title
    };
    this.ref = React.createRef();

    this.handleStartEdit = this.handleStartEdit.bind(this);
    this.handleEndEdit = this.handleEndEdit.bind(this);
  }

  handleStartEdit() {
    this.setState({ inEdit: true });
  }

  handleEndEdit() {
    this.setState({ inEdit: false });
  }

  render() {
    const { tabKey } = this.props;
    const { title, inEdit } = this.state;

    if (inEdit) {
      return (
        <div ref={this.ref}>
          <input
            type="text"
            value={title}
            onChange={e => {
              this.setState({ title: e.target.value });
            }}
            onKeyPress={e => {
              if (e.key === 'Enter') {
                this.handleEndEdit();
              }
            }}
            onBlur={this.handleEndEdit}
          />
        </div>
      );
    }

    return (
      <div>
        <span className="tab-title">{title}</span>
        <Button
          className="close"
          onClick={() => {
            this.props.handleRemoveResult(tabKey);
          }}
        >
          <span aria-hidden="true">&times;</span>
        </Button>
        <Button
          className="close edit"
          onClick={() => {
            this.handleStartEdit(tabKey);
          }}
        >
          <Glyphicon glyph="pencil" />
        </Button>
      </div>
    );
  }
}

TabTitle.propTypes = {
  tabKey: PropTypes.any.isRequired,
  title: PropTypes.string.isRequired,
  handleRemoveResult: PropTypes.func.isRequired
};

class TabbedResults extends Component {
  constructor(props) {
    super(props);
    this.state = {
      activeKey: 0,
      nTabs: 0
    };

    this.handleSelect = this.handleSelect.bind(this);
  }

  static getDerivedStateFromProps(nextProps, prevState) {
    const nextNTabs = nextProps.results.size;

    if (prevState.nTabs === nextNTabs) {
      return null;
    }

    // Make the last tab active when a new one is added,
    // or when the current active tab is deleted
    const newTab = prevState.nTabs < nextNTabs;
    const lostTab = !nextProps.results.has(prevState.activeKey);

    if (newTab || lostTab) {
      return {
        activeKey: Array.from(nextProps.results.keys())[nextNTabs - 1],
        nTabs: nextNTabs
      };
    }

    return { nTabs: nextNTabs };
  }

  handleSelect(activeKey) {
    this.setState({ activeKey });
  }

  render() {
    const { showCode, showUAST, history } = this.props;

    return (
      <Tabs
        id="tabbed-results"
        className="full-height"
        activeKey={this.state.activeKey}
        onSelect={this.handleSelect}
      >
        {Array.from(this.props.results.entries()).map(([key, query]) => {
          let content = '';

          if (query.loading) {
            content = (
              <Row>
                <Col className="text-center loader-col" xs={12}>
                  <Loader />
                </Col>
              </Row>
            );
          } else if (query.errorMsg) {
            content = (
              <Row className="errors-row">
                <Col xs={12}>
                  <Alert bsStyle="danger">{query.errorMsg}</Alert>
                </Col>
              </Row>
            );
          } else {
            content = (
              <ResultsTable
                response={query.response}
                showCode={showCode}
                showUAST={showUAST}
              />
            );
          }

          return (
            <Tab
              key={key}
              eventKey={key}
              title={
                <TabTitle
                  title={query.title || query.sql}
                  tabKey={key}
                  handleRemoveResult={this.props.handleRemoveResult}
                />
              }
            >
              <Row className="query-row">
                <Col xs={12}>
                  <div className="query-text">
                    <p>{query.sql}</p>
                  </div>
                  <Button
                    className="edit-query"
                    bsStyle="link"
                    onClick={() => this.props.handleEditQuery(query.sql)}
                  >
                    edit query
                  </Button>
                </Col>
              </Row>
              {content}
            </Tab>
          );
        })}
        <Tab key="history" eventKey="history" title="History">
          <HistoryTable
            items={history}
            onOpenQuery={this.props.handleEditQuery}
            handleReset={this.props.handleResetHistory}
          />
        </Tab>
      </Tabs>
    );
  }
}

TabbedResults.propTypes = {
  // results is a Map of objects, each object may contain:
  // sql: 'string'      Required
  // loading: true      Optional, tab will show a loading animation
  // errorMsg: 'string' Optional
  // response: object   Required if loading and errorMsg are not present
  results: PropTypes.instanceOf(Map).isRequired,
  history: HistoryTable.propTypes.items,
  handleRemoveResult: PropTypes.func.isRequired,
  handleEditQuery: PropTypes.func.isRequired,
  handleResetHistory: PropTypes.func.isRequired,
  showCode: PropTypes.func.isRequired,
  showUAST: PropTypes.func.isRequired
};

export default TabbedResults;
