import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ReactTable from 'react-table';
import 'react-table/react-table.css';
import './ResultsTable.less';

class ResultsTable extends Component {
  render() {
    const columns = this.props.response.meta.headers.map(col => ({
      Header: col,
      id: col,
      accessor: row => {
        const v = row[col];
        switch (typeof v) {
          case 'boolean':
            return v.toString();
          case 'object':
            return JSON.stringify(v, null, 2);
          default:
            return v;
        }
      }
    }));

    return (
      <ReactTable
        className="ResultsTable"
        data={this.props.response.data}
        columns={columns}
      />
    );
  }
}

ResultsTable.propTypes = {
  // Must be a success response
  response: PropTypes.object.isRequired
};

export default ResultsTable;
