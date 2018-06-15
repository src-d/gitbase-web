import React, { Component } from 'react';
import { Helmet } from 'react-helmet';
import { Grid, Row, Modal } from 'react-bootstrap';
import SplitPane from 'react-split-pane';
import UASTViewer, { Editor, transformer } from 'uast-viewer';
import 'uast-viewer/dist/default-theme.css';
import initButtonStyles from './utils/bootstrap';
import Sidebar from './components/Sidebar';
import QueryBox from './components/QueryBox';
import TabbedResults from './components/TabbedResults';
import api from './api';
import { STATUS_LOADING, STATUS_ERROR, STATUS_SUCCESS } from './state/query';
import './App.less';

const INACTIVE_TIMEOUT = 3600000;

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      sql: `/* Contributor's number of commits, each month of 2018, for each repository */

SELECT COUNT(*) as num_commits, month, repo_id, committer_name
FROM ( SELECT MONTH(committer_when) as month,
              r.repository_id as repo_id,
              committer_name
    FROM ref_commits r
    INNER JOIN commits c
        ON YEAR(c.committer_when) = 2018 AND r.commit_hash = c.commit_hash
    WHERE r.ref_name = 'HEAD'
) as t GROUP BY committer_name, month, repo_id`,
      results: new Map(),
      schema: undefined,
      history: [],
      lastResultMeta: null,

      // modal
      showModal: false,
      modalTitle: null,
      modalContent: null
    };
    this.timers = {};

    this.handleTextChange = this.handleTextChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleRemoveResult = this.handleRemoveResult.bind(this);
    this.handleTableClick = this.handleTableClick.bind(this);
    this.handleExampleClick = this.handleExampleClick.bind(this);
    this.handleModalClose = this.handleModalClose.bind(this);
    this.handleResetHistory = this.handleResetHistory.bind(this);
    this.handleSetActiveResult = this.handleSetActiveResult.bind(this);
    this.handleReload = this.handleReload.bind(this);

    this.showCode = this.showCode.bind(this);
    this.showUAST = this.showUAST.bind(this);

    this.uniqueKey = 0;
  }

  handleTextChange(text) {
    this.setState({ sql: text });
  }

  setResult(key, result) {
    const { results, history } = this.state;
    const historyIdx = history.findIndex(i => i.key === key);
    const { response } = result;

    if (historyIdx >= 0) {
      const status =
        typeof response !== 'undefined' ? STATUS_SUCCESS : STATUS_ERROR;
      const newHistory = [
        ...history.slice(0, historyIdx),
        {
          ...history[historyIdx],
          status,
          elapsedTime:
            response && response.meta ? response.meta.elapsedTime : null,
          errorMsg: result.errorMsg
        },
        ...history.slice(historyIdx + 1)
      ];
      this.setState({ history: newHistory });
    }

    if (!results.has(key)) {
      // Tab was removed, ignore results
      return;
    }

    const newResults = new Map(this.state.results);
    newResults.set(key, result);

    this.setState({
      results: newResults,
      lastResultMeta: response ? response.meta : null
    });
  }

  handleSubmit() {
    const { sql, history } = this.state;
    const key = ++this.uniqueKey;

    const loadingResults = new Map(this.state.results);
    loadingResults.set(key, { sql, loading: true });

    this.setState(
      {
        results: loadingResults,
        history: [
          {
            key,
            sql,
            datetime: new Date(),
            status: STATUS_LOADING
          },
          ...history
        ]
      },
      () => this.handleSetActiveResult(key)
    );

    api
      .query(sql)
      .then(response => {
        this.setResult(key, { sql, response });

        if (!this.state.schema) {
          // The schema was not loaded for some reason, and we know we just
          // did a successful call to the backend. Let's retry.
          this.loadSchema();
        }
      })
      .catch(msgArr =>
        this.setResult(key, { sql, errorMsg: msgArr.join('; ') })
      );
  }

  handleTableClick(table) {
    this.setState({ sql: `DESCRIBE TABLE ${table}` }, this.handleSubmit);
  }

  handleExampleClick(sql) {
    this.setState({ sql }, this.handleSubmit);
  }

  loadSchema() {
    api
      .schema()
      .then(schema => {
        if (JSON.stringify(schema) !== JSON.stringify(this.state.schema)) {
          this.setState({ schema });
        }
      })
      .catch(msgArr => {
        // TODO (@carlosms): left as console message for now, we may decide to
        // show it in the interface somehow when we have to populate the sidebar
        // eslint-disable-next-line no-console
        console.error(`Error while loading schema: ${msgArr}`);
      });
  }

  handleModalClose() {
    this.setState({ showModal: false, modalTitle: null, modalContent: null });
  }

  showCode(code) {
    this.setState({
      showModal: true,
      modalTitle: 'Source code',
      modalContent: <Editor code={code} />
    });
  }

  showUAST(uast) {
    this.setState({
      showModal: true,
      modalTitle: 'UAST',
      modalContent:
        // currently gitbase returns only 1 item, UAST of the file
        // but just in case if there is more or less we show it without viewer
        uast.length === 1 ? (
          <UASTViewer uast={transformer(uast[0])} />
        ) : (
          <pre>{uast}</pre>
        )
    });
  }

  componentDidMount() {
    this.loadSchema();
  }

  handleRemoveResult(key) {
    const newResults = new Map(this.state.results);
    newResults.delete(key);

    this.setState({ results: newResults });

    this.stopInactiveTimer(key);
  }

  handleResetHistory() {
    this.setState({ history: [] });
  }

  handleSetActiveResult(key) {
    const { results } = this.state;
    // just ignore any unknown key
    if (!results.has(key)) {
      return;
    }

    this.stopInactiveTimer(key);

    const hiddenKeys = Array.from(results.keys()).filter(k => k !== key);
    hiddenKeys.filter(k => !this.timers[k]).forEach(k => {
      this.timers[k] = window.setTimeout(
        () => this.removeResultContent(k),
        INACTIVE_TIMEOUT
      );
    });
  }

  stopInactiveTimer(key) {
    const timer = this.timers[key];
    if (timer) {
      window.clearTimeout(timer);
      delete this.timers[key];
    }
  }

  removeResultContent(key) {
    const { results } = this.state;

    const result = results.get(key);
    result.response = null;

    const newResults = new Map(results);
    newResults.set(key, result);

    this.setState({ results: newResults });
  }

  handleReload(key) {
    const { results } = this.state;

    const result = results.get(key);
    result.loading = true;

    const newResults = new Map(results);
    newResults.set(key, result);

    this.setState({ results: newResults });

    api
      .query(result.sql)
      .then(response => {
        result.response = response;
      })
      .catch(msgArr => {
        result.errorMsg = msgArr.join('; ');
      })
      .then(() => {
        const newNewResults = new Map(results);
        result.loading = false;
        newResults.set(key, result);
        this.setState({ results: newNewResults });
      });
  }

  render() {
    const { results, history } = this.state;

    return (
      <div className="app">
        <Helmet>
          <title>Gitbase Playground</title>
        </Helmet>
        <Grid className="full-height app-grid" fluid={true}>
          <Row className="main-row full-height">
            <Sidebar
              schema={this.state.schema}
              onTableClick={this.handleTableClick}
              onExampleClick={this.handleExampleClick}
            />
            <SplitPane
              className="main-split"
              split="horizontal"
              defaultSize={250}
              minSize={175}
            >
              <QueryBox
                sql={this.state.sql}
                schema={this.state.schema}
                resultMeta={this.state.lastResultMeta}
                handleTextChange={this.handleTextChange}
                handleSubmit={this.handleSubmit}
                exportUrl={api.queryExport(this.state.sql)}
              />
              <TabbedResults
                results={results}
                history={history}
                handleRemoveResult={this.handleRemoveResult}
                handleEditQuery={this.handleTextChange}
                handleResetHistory={this.handleResetHistory}
                handleSetActiveResult={this.handleSetActiveResult}
                handleReload={this.handleReload}
                showCode={this.showCode}
                showUAST={this.showUAST}
              />
            </SplitPane>
          </Row>
        </Grid>
        <Modal
          show={this.state.showModal}
          onHide={this.handleModalClose}
          bsSize="large"
        >
          <Modal.Header closeButton>
            <Modal.Title>{this.state.modalTitle}</Modal.Title>
          </Modal.Header>
          <Modal.Body>{this.state.modalContent}</Modal.Body>
        </Modal>
      </div>
    );
  }
}

initButtonStyles();

export default App;
