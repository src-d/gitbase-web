import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Row, Col, Button } from 'react-bootstrap';
import './QueryBox.less';

class QueryBox extends Component {
  render() {
    return (
      <div className="query-box">
        <Row>
          <Col xs={12}>
            <textarea
              autoFocus="true"
              rows="7"
              placeholder="Enter an SQL query"
              value={this.props.sql}
              onChange={e => this.props.handleTextChange(e.target.value)}
            />
          </Col>
        </Row>
        <Row>
          <Col xsOffset={9} xs={3}>
            <Button
              className="pull-right"
              disabled={!this.props.enabled}
              onClick={this.props.handleSubmit}
            >
              RUN
            </Button>
          </Col>
        </Row>
      </div>
    );
  }
}

QueryBox.propTypes = {
  sql: PropTypes.string.isRequired,
  enabled: PropTypes.bool.isRequired,
  handleTextChange: PropTypes.func.isRequired,
  handleSubmit: PropTypes.func.isRequired
};

export default QueryBox;
