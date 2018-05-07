import React, { Component } from 'react';
import PropTypes from 'prop-types';
// import './ResultsTable.css';

class ResultsTable extends Component {
  render() {
    return <pre>{JSON.stringify(this.props.response, null, 2)}</pre>;
  }
}

ResultsTable.propTypes = {
  response: PropTypes.object.isRequired
};

export default ResultsTable;
