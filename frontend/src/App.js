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
      results: [],
      loading: false
    };

    this.handleTextChange = this.handleTextChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleRemoveResult = this.handleRemoveResult.bind(this);
  }

  handleTextChange(text) {
    this.setState({ sql: text });
  }

  handleSubmit() {
    if (this.state.loading) {
      // Should not happen, but just in case
      return;
    }

    const { sql } = this.state;

    this.setState({
      loading: true,
      results: [...this.state.results, { sql, loading: true }]
    });

    // After a success or failure we pop the last entry {loading:true} and
    // push a new result. This works because the 'run' button is disabled
    // until the request finishes
    api
      .query(this.state.sql)
      .then(response => {
        this.setState({
          results: [...this.state.results.slice(0, -1), { sql, response }],
          loading: false
        });
      })
      .catch(msgArr => {
        this.setState({
          results: [
            ...this.state.results.slice(0, -1),
            { sql, errorMsg: msgArr.join('; ') }
          ],
          loading: false
        });
      });
  }

  handleRemoveResult(index) {
    this.setState({
      results: [
        ...this.state.results.slice(0, index),
        ...this.state.results.slice(index + 1)
      ]
    });
  }

  render() {
    const { results } = this.state;

    let resultsElem = '';
    if (results.length > 0) {
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
                enabled={!this.state.loading}
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
