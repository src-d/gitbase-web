import React, { Component } from 'react';
import { Helmet } from 'react-helmet';
import { Grid, Row, Col } from 'react-bootstrap';
import SplitPane from 'react-split-pane';
import Sidebar from './components/Sidebar';
import QueryBox from './components/QueryBox';
import TabbedResults from './components/TabbedResults';
import api from './api';
import './App.less';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      sql: `SELECT COUNT(*) as num_commits, month, repo_id, committer_email FROM (
  SELECT MONTH(committer_when) as month, r.repository_id as repo_id, committer_email
  FROM repositories r
  INNER JOIN refs ON refs.repository_id = r.repository_id AND refs.ref_name = 'HEAD'
  INNER JOIN commits c ON YEAR(committer_when) = 2018 AND history_idx(refs.commit_hash, c.commit_hash) >= 0
) as t
GROUP BY committer_email, month, repo_id`,
      results: new Map(),
      schema: undefined
    };

    this.handleTextChange = this.handleTextChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleRemoveResult = this.handleRemoveResult.bind(this);

    this.uniqueKey = 0;
  }

  handleTextChange(text) {
    this.setState({ sql: text });
  }

  setResult(key, result) {
    if (!this.state.results.has(key)) {
      // Tab was removed, ignore results
      return;
    }

    const newResults = new Map(this.state.results);
    newResults.set(key, result);

    this.setState({ results: newResults });
  }

  handleSubmit() {
    const { sql } = this.state;
    const key = ++this.uniqueKey;

    const loadingResults = new Map(this.state.results);
    loadingResults.set(key, { sql, loading: true });

    this.setState({ results: loadingResults });

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

  componentDidMount() {
    this.loadSchema();
  }

  handleRemoveResult(key) {
    const newResults = new Map(this.state.results);
    newResults.delete(key);

    this.setState({ results: newResults });
  }

  render() {
    const { results } = this.state;

    let resultsElem = '';
    if (results.size > 0) {
      resultsElem = (
        <Col xs={12} className="full-height">
          <TabbedResults
            results={results}
            handleRemoveResult={this.handleRemoveResult}
            handleEditQuery={this.handleTextChange}
          />
        </Col>
      );
    }
    return (
      <div className="app">
        <Helmet>
          <title>Gitbase Playground</title>
        </Helmet>
        <Grid className="full-height" fluid={true}>
          <Row className="full-height">
            <Col xs={3} className="full-height">
              <Sidebar schema={this.state.schema} />
            </Col>
            <Col xs={9} className="full-height">
              <SplitPane split="horizontal" defaultSize={250} minSize={100}>
                <Grid className="full-height full-width">
                  <Row className="query-box-row">
                    <Col xs={12} className="full-height">
                      <QueryBox
                        sql={this.state.sql}
                        schema={this.state.schema}
                        handleTextChange={this.handleTextChange}
                        handleSubmit={this.handleSubmit}
                        exportUrl={api.queryExport(this.state.sql)}
                      />
                    </Col>
                  </Row>
                </Grid>
                <Grid className="full-height full-width">
                  <Row className="results-row">{resultsElem}</Row>
                </Grid>
              </SplitPane>
            </Col>
          </Row>
        </Grid>
      </div>
    );
  }
}

export default App;
