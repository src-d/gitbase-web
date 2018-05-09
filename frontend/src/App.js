import React, { Component } from 'react';
import { Helmet } from 'react-helmet';
import { Grid, Row, Col, Alert } from 'react-bootstrap';
import QueryBox from './components/QueryBox';
import ResultsTable from './components/ResultsTable';
import Loader from './components/Loader';
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
      response: undefined,
      loading: false,
      error: undefined
    };

    this.handleTextChange = this.handleTextChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleTextChange(text) {
    this.setState({ sql: text });
  }

  handleSubmit() {
    this.setState({
      response: undefined,
      loading: true,
      error: undefined
    });

    api
      .query(this.state.sql)
      .then(json => {
        this.setState({
          response: json,
          loading: false
        });
      })
      .catch(msgArr => {
        this.setState({ loading: false, error: msgArr.join('; ') });
      });
  }

  render() {
    const { response } = this.state;

    let results = '';

    if (this.state.loading) {
      results = (
        <Col className="text-center loader-col" xs={12}>
          <Loader />
        </Col>
      );
    } else if (response && response.status === 200) {
      results = (
        <Col xs={12}>
          <ResultsTable response={response} />
        </Col>
      );
    } else if (this.state.error) {
      results = (
        <Col xs={10} xsOffset={1}>
          <Alert bsStyle="danger">{this.state.error}</Alert>
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
          <Row className="results-row">{results}</Row>
        </Grid>
      </div>
    );
  }
}

export default App;
