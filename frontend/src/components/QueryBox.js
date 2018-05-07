import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Grid, Row, Col, Button } from 'react-bootstrap';
import './QueryBox.less';

class QueryBox extends Component {
  render() {
    return (
      <Grid className="QueryBox">
        <Row>
          <Col xs={12}>
            <textarea
              rows="7"
              placeholder="Enter an SQL query"
              value={this.props.sql}
              onChange={e => this.props.handleTextChange(e.target.value)}
            />
          </Col>
        </Row>
        <Row>
          <Col xs={12}>
            <Button onClick={this.props.handleSubmit}>RUN</Button>
          </Col>
        </Row>
      </Grid>
    );
  }
}

QueryBox.propTypes = {
  sql: PropTypes.string.isRequired,
  handleTextChange: PropTypes.func.isRequired,
  handleSubmit: PropTypes.func.isRequired
};

export default QueryBox;
