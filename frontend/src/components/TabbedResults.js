import React, { Component, Fragment } from 'react';
import { Row, Col, Alert, Tab, Button } from 'react-bootstrap';
import PropTypes from 'prop-types';
import DivTabs from './DivTabs';
import ResultsTable from './ResultsTable';
import HistoryTable from './HistoryTable';
import './TabbedResults.less';
import PencilIcon from '../icons/edit-query-tab-name.svg';
import CloseIcon from '../icons/close-query-tab.svg';
import TimerIcon from '../icons/history-tab.svg';
import LoadingImg from '../icons/alex-loading-results.gif';
import SuspendedImg from '../icons/alex-suspended-tab.gif';

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
    this.handleChange = this.handleChange.bind(this);
    this.handleKeyPress = this.handleKeyPress.bind(this);
    this.handleBlur = this.handleBlur.bind(this);
  }

  handleStartEdit() {
    this.setState({ inEdit: true });
  }

  handleEndEdit() {
    this.setState({ inEdit: false });
  }

  handleChange(e) {
    this.setState({ title: e.target.value });
  }

  handleKeyPress(e) {
    if (e.key === 'Enter') {
      this.handleEndEdit();
    }
  }

  handleBlur() {
    this.handleEndEdit();
  }

  render() {
    const { tabKey } = this.props;
    const { title, inEdit } = this.state;

    if (inEdit) {
      return (
        <div ref={this.ref}>
          <input
            autoFocus
            type="text"
            className="tab-title"
            value={title}
            onChange={this.handleChange}
            onKeyPress={this.handleKeyPress}
            /* disable handlers of tab components */
            onKeyDown={e => e.stopPropagation()}
            onBlur={this.handleBlur}
          />
        </div>
      );
    }

    return (
      <div>
        <span className="tab-title">{title}</span>
        <PencilIcon
          className="btn-title"
          onClick={() => {
            this.handleStartEdit(tabKey);
          }}
        />
        <CloseIcon
          className="btn-title"
          onClick={() => {
            this.props.handleRemoveResult(tabKey);
          }}
        />
      </div>
    );
  }
}

TabTitle.propTypes = {
  tabKey: PropTypes.any.isRequired,
  active: PropTypes.bool.isRequired,
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
    this.props.handleSetActiveResult(activeKey);
  }

  render() {
    const { showCode, showUAST, history } = this.props;

    return (
      <div className="results-padding full-height full-width">
        <DivTabs
          id="tabbed-results"
          className="full-height"
          activeKey={this.state.activeKey}
          onSelect={this.handleSelect}
        >
          {Array.from(this.props.results.entries()).map(([key, query]) => {
            let content = '';
            if (key === this.state.activeKey) {
              if (query.loading) {
                content = (
                  <Fragment>
                    <Row>
                      <Col className="text-center animation-col" xs={12}>
                        <img src={LoadingImg} alt="loading animation" />
                      </Col>
                    </Row>
                    <Row>
                      <Col
                        className="text-center message-col last-message-col"
                        xs={12}
                      >
                        RUNNING QUERY
                      </Col>
                    </Row>
                  </Fragment>
                );
              } else if (query.errorMsg) {
                content = (
                  <Row className="errors-row">
                    <Col xs={12}>
                      <Alert bsStyle="danger">{query.errorMsg}</Alert>
                    </Col>
                  </Row>
                );
              } else if (query.response) {
                content = (
                  <ResultsTable
                    response={query.response}
                    showCode={showCode}
                    showUAST={showUAST}
                  />
                );
              } else {
                content = (
                  <Fragment>
                    <Row>
                      <Col className="text-center animation-col" xs={12}>
                        <img src={SuspendedImg} alt="suspended animation" />
                      </Col>
                    </Row>
                    <Row>
                      <Col className="text-center message-col" xs={12}>
                        SUSPENDED TAB
                      </Col>
                    </Row>
                    <Row>
                      <Col
                        className="text-center message-col last-message-col"
                        xs={12}
                      >
                        <Button
                          className="reload"
                          bsStyle="gbpl-tertiary"
                          onClick={() => this.props.handleReload(key)}
                        >
                          RELOAD
                        </Button>
                      </Col>
                    </Row>
                  </Fragment>
                );
              }
            }

            return (
              <Tab
                key={key}
                eventKey={key}
                title={
                  <TabTitle
                    title={query.title || query.sql}
                    tabKey={key}
                    active={key === this.state.activeKey}
                    handleRemoveResult={this.props.handleRemoveResult}
                  />
                }
              >
                <div className="query-row">
                  <div className="query-text">{query.sql}</div>
                  <div className="query-button">
                    <Button
                      className="edit-query"
                      bsStyle="gbpl-tertiary-tint-2-link"
                      onClick={() => this.props.handleEditQuery(query.sql)}
                    >
                      EDIT
                    </Button>
                  </div>
                </div>
                {content}
              </Tab>
            );
          })}
          {history.length > 0 && (
            <Tab
              key="history"
              eventKey="history"
              tabClassName="history-tab-title"
              title={
                <div className="history-tab">
                  <TimerIcon className="icon-title" />
                  <span className="tab-title">history</span>
                </div>
              }
            >
              <HistoryTable
                items={history}
                onOpenQuery={this.props.handleEditQuery}
                handleReset={this.props.handleResetHistory}
              />
            </Tab>
          )}
        </DivTabs>
      </div>
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
  handleSetActiveResult: PropTypes.func.isRequired,
  handleReload: PropTypes.func.isRequired,
  showCode: PropTypes.func.isRequired,
  showUAST: PropTypes.func.isRequired
};

export default TabbedResults;
