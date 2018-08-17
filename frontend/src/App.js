import React, { Component } from 'react';
import { Helmet } from 'react-helmet';
import { Grid, Row, Modal } from 'react-bootstrap';
import nanoid from 'nanoid';
import SplitPane from 'react-split-pane';
import 'uast-viewer/dist/default-theme.css';
import initButtonStyles from './utils/bootstrap';
import Sidebar from './components/Sidebar';
import QueryBox from './components/QueryBox';
import TabbedResults from './components/TabbedResults';
import UASTViewer from './components/UASTViewer';
import api from './api';
import { STATUS_LOADING, STATUS_ERROR, STATUS_SUCCESS } from './state/query';
import './App.less';
import CloseIcon from './icons/close-query-tab.svg';

const INACTIVE_TIMEOUT = 3600000;

function persist(key, value) {
  window.localStorage.setItem(key, JSON.stringify(value));
}

const dateFormat = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/;

function jsonReviver(key, value) {
  if (typeof value === 'string' && dateFormat.test(value)) {
    return new Date(value);
  }

  return value;
}

function loadStateFromStorage() {
  const historyJSON = window.localStorage.getItem('history');

  return {
    history: historyJSON ? JSON.parse(historyJSON, jsonReviver) : []
  };
}

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      sql: '',
      results: new Map(),
      schema: undefined,
      languages: [],
      history: [],
      lastResult: null,

      // modal
      showModal: false,
      modalTitle: null,
      modalContent: null
    };
    this.timers = {};

    this.handleTextChange = this.handleTextChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleRemoveResult = this.handleRemoveResult.bind(this);
    this.handleAbortQuery = this.handleAbortQuery.bind(this);
    this.handleTableClick = this.handleTableClick.bind(this);
    this.handleExampleClick = this.handleExampleClick.bind(this);
    this.handleModalClose = this.handleModalClose.bind(this);
    this.handleResetHistory = this.handleResetHistory.bind(this);
    this.handleSetActiveResult = this.handleSetActiveResult.bind(this);
    this.handleReload = this.handleReload.bind(this);

    this.showUAST = this.showUAST.bind(this);

    this.uniqueKey = 0;

    this.exampleQueries = [
      {
        name: 'Files named main.go',
        sql: `/* Files named main.go in HEAD */
SELECT f.repository_id, f.file_path,
       LANGUAGE(f.file_path, f.blob_content) AS lang, f.blob_content,
       UAST(f.blob_content, LANGUAGE(f.file_path, f.blob_content)) AS uast
FROM   files AS f
       JOIN commit_files cf ON
            f.repository_id=cf.repository_id AND
            f.file_path=cf.file_path AND
            f.blob_hash=cf.blob_hash AND
            f.tree_hash=cf.tree_hash
       JOIN refs ON
            cf.repository_id = refs.repository_id AND
            cf.commit_hash = refs.commit_hash
WHERE  ref_name = 'HEAD'
       AND f.file_path REGEXP('.*main.go')`
      },

      {
        name: 'Last commit for each repository',
        sql: `/* Last commit for each repository */
SELECT r.repository_id, commit_author_name, commit_author_when, commit_message
FROM   refs r
       NATURAL JOIN commits
WHERE  r.ref_name = 'HEAD'`
      },
      {
        name: 'Top repositories by commits',
        sql: `/* Top repositories by number of commits in HEAD */
SELECT repository_id, commit_count
FROM   (SELECT r.repository_id, COUNT(*) AS commit_count
        FROM   refs r
               JOIN ref_commits AS c ON r.ref_name = c.ref_name
        WHERE  r.ref_name = 'HEAD'
        GROUP  BY r.repository_id) AS q
ORDER  BY commit_count DESC
LIMIT  10 `
      },
      {
        name: 'Top languages by repository count',
        sql: `/* Top languages by repository count */
SELECT *
FROM (SELECT language, COUNT(repository_id) AS repository_count
      FROM   (SELECT DISTINCT
                r.repository_id,
                LANGUAGE(t.tree_entry_name, b.blob_content) AS language
              FROM   refs r
                      JOIN commits c ON r.commit_hash = c.commit_hash
                      JOIN commit_trees ct ON c.commit_hash = ct.commit_hash
                      JOIN tree_entries t ON ct.tree_hash = t.tree_hash
                      JOIN blobs b ON t.blob_hash = b.blob_hash
              WHERE  r.ref_name = 'HEAD') AS q1
      GROUP  BY language) AS q2
ORDER  BY repository_count DESC `
      },
      {
        name: 'Number of commits per month',
        sql: `/* Commits per committer, each month of 2018, for each repository */
SELECT COUNT(*) as num_commits, month, repository_id, committer_name, committer_email
FROM ( SELECT MONTH(committer_when) as month,
              r.repository_id,
              committer_name,
              committer_email
    FROM ref_commits r
    INNER JOIN commits c
        ON YEAR(c.committer_when) = 2018 AND r.commit_hash = c.commit_hash
    WHERE r.ref_name = 'HEAD'
) as t GROUP BY committer_email, committer_name, month, repository_id`
      },
      {
        name: 'Number of functions per file for Go',
        sql: ` /* Number of funtions per file for Golang at HEAD */
SELECT files.repository_id, files.file_path,
 ARRAY_LENGTH(UAST(
   files.blob_content,
   LANGUAGE(files.file_path, files.blob_content),
   '//*[@roleFunction and @roleDeclaration]'
 )) as functions
FROM files
NATURAL JOIN refs
WHERE
LANGUAGE(files.file_path,files.blob_content) = 'Go'
AND refs.ref_name = 'HEAD'`
      }
    ];
  }

  setState(partialState, callback) {
    super.setState(partialState, callback);

    if (typeof partialState.history !== 'undefined') {
      persist('history', partialState.history);
    }
  }

  handleTextChange(text) {
    this.setState({ sql: text });
  }

  setResult(key, result) {
    const { results, history } = this.state;
    const historyIdx = history.findIndex(i => i.key === key);
    const { response } = result;

    const status =
      typeof response !== 'undefined' ? STATUS_SUCCESS : STATUS_ERROR;

    if (historyIdx >= 0) {
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
      lastResult: result
    });
  }

  loadQuery(key, sql) {
    const { history } = this.state;
    const abortController = new window.AbortController();

    const loadingResults = new Map(this.state.results);
    loadingResults.set(key, { sql, loading: true, abortController });

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
        ],
        lastResult: null
      },
      () => this.handleSetActiveResult(key)
    );

    api
      .query(sql, abortController.signal)
      .then(response => {
        this.setResult(key, { sql, response });

        // If the schema or the list of languages were not loaded for some
        // reason, we retry now

        if (!this.state.schema) {
          this.loadSchema();
        }

        if (this.state.languages.length === 0) {
          this.loadLanguages();
        }
      })
      .catch(msgArr => {
        this.setResult(key, { sql, errorMsg: msgArr.join('; ') });
      });
  }

  handleSubmit() {
    const { sql } = this.state;
    const key = nanoid();
    this.loadQuery(key, sql);
  }

  handleReload(key) {
    const result = this.state.results.get(key);
    if (result) {
      this.loadQuery(key, result.sql);
    }
  }

  handleTableClick(table) {
    this.setState({ sql: `DESCRIBE TABLE ${table}` }, this.handleSubmit);
  }

  handleExampleClick(sql) {
    this.setState({ sql });
  }

  loadSchema() {
    api
      .schema()
      .then(schema =>
        schema.sort((a, b) => {
          if (a.table < b.table) {
            return -1;
          }
          if (a.table > b.table) {
            return 1;
          }

          return 0;
        })
      )
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

  showUAST(uast, protobufs) {
    this.setState({
      showModal: true,
      modalTitle: (
        <div>
          UAST
          <CloseIcon
            className="btn-modal-close"
            onClick={this.handleModalClose}
          />
        </div>
      ),
      modalContent: <UASTViewer uast={uast} protobufs={protobufs} />
    });
  }

  loadLanguages() {
    api
      .getLanguages()
      .then(languages => this.setState({ languages }))
      .catch(err =>
        // we don't have UI for this error
        console.error(`Can't get list of languages from bblfsh: ${err}`)
      );
  }

  componentDidMount() {
    this.setState(loadStateFromStorage());
    this.loadSchema();
    this.loadLanguages();
    this.handleExampleClick(this.exampleQueries[0].sql);
  }

  handleRemoveResult(key) {
    this.handleAbortQuery(key);

    const newResults = new Map(this.state.results);
    newResults.delete(key);

    this.setState({ results: newResults });

    this.stopInactiveTimer(key);
  }

  handleAbortQuery(key) {
    const result = this.state.results.get(key);

    if (result && result.abortController) {
      result.abortController.abort();
    }
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

  render() {
    const { results, history } = this.state;

    return (
      <div className="app">
        <Helmet>
          <title>Gitbase Web</title>
        </Helmet>
        <Grid className="full-height app-grid" fluid={true}>
          <Row className="main-row full-height">
            <Sidebar
              schema={this.state.schema}
              onTableClick={this.handleTableClick}
              onExampleClick={this.handleExampleClick}
              exampleQueries={this.exampleQueries}
            />
            <SplitPane
              className="main-split"
              split="horizontal"
              defaultSize={250}
              minSize={1}
              maxSize={-15}
            >
              <QueryBox
                sql={this.state.sql}
                schema={this.state.schema}
                result={this.state.lastResult}
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
                handleAbortQuery={this.handleAbortQuery}
                showUAST={this.showUAST}
                languages={this.state.languages}
              />
            </SplitPane>
          </Row>
        </Grid>
        <Modal
          show={this.state.showModal}
          onHide={this.handleModalClose}
          bsSize="large"
        >
          <Modal.Header>
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
