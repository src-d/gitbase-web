import React, { Component } from 'react';
import { Helmet } from 'react-helmet';
import { Grid, Row, Col } from 'react-bootstrap';
import QueryBox from './components/QueryBox';
import TabbedResults from './components/TabbedResults';
import api from './api';
import './App.less';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      sql: `SELECT COUNT(*) as num_commits, month, repo_id, committer_email FROM (
  SELECT MONTH(committer_when) as month, r.id as repo_id, committer_email
  FROM repositories r
  INNER JOIN refs ON refs.repository_id = r.id AND refs.name = 'HEAD'
  INNER JOIN commits c ON YEAR(committer_when) = 2018 AND history_idx(refs.hash, c.hash) >= 0
) as t
GROUP BY committer_email, month, repo_id`,
      results: new Map()
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
      .then(response => this.setResult(key, { sql, response }))
      .catch(msgArr =>
        this.setResult(key, { sql, errorMsg: msgArr.join('; ') })
      );
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
        <Col xs={12}>
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
        <Grid>
          <Row className="query-row">
            <Col xs={12}>
              <QueryBox
                sql={this.state.sql}
                handleTextChange={this.handleTextChange}
                handleSubmit={this.handleSubmit}
              />
            </Col>
          </Row>
          <Row className="results-row">{resultsElem}</Row>
        </Grid>
      </div>
    );
  }
}

export default App;
