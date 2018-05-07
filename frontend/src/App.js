import React, { Component } from 'react';
import { Helmet } from 'react-helmet';
import { Grid, Row, Col } from 'react-bootstrap';
import QueryBox from './components/QueryBox';
import ResultsTable from './components/ResultsTable';
import api from './api';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      sql: `SELECT COUNT(*) as num_commits, month, repo_id, committer_email
	FROM (
		SELECT
			MONTH(committer_when) as month,
			r.id as repo_id,
			committer_email
		FROM repositories r
		INNER JOIN refs ON refs.repository_id = r.id AND refs.name = 'HEAD'
		INNER JOIN commits c ON YEAR(committer_when) = 2018 AND history_idx(refs.hash, c.hash) >= 0
	) as t
GROUP BY committer_email, month, repo_id`,
      response: undefined,
      debug: 'Debug info will appear here.'
    };

    this.handleTextChange = this.handleTextChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleTextChange(text) {
    this.setState({ sql: text });
  }

  handleSubmit() {
    this.setState({ response: undefined, debug: 'Loading...' });

    api
      .query(this.state.sql)
      .then(json => {
        this.setState({ response: json, debug: JSON.stringify(json, null, 2) });
      })
      .catch(e => this.setState({ debug: e }));
  }

  render() {
    return (
      <div className="App">
        <Helmet>
          <title>gitbase-playground</title>
        </Helmet>
        <Grid>
          <Row>
            <Col xs={10} xsOffset={1}>
              <QueryBox
                sql={this.state.sql}
                handleTextChange={this.handleTextChange}
                handleSubmit={this.handleSubmit}
              />
            </Col>
          </Row>
          <Row>
            <Col xs={10} xsOffset={1}>
              <ResultsTable response={this.state.response} />
            </Col>
          </Row>
          <Row>
            <Col xs={10} xsOffset={1}>
              <div>
                <pre>{this.state.debug}</pre>
              </div>
            </Col>
          </Row>
        </Grid>
      </div>
    );
  }
}

export default App;
