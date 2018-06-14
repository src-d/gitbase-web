import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ReactTable from 'react-table';
import 'react-table/react-table.css';
import { STATUS_LOADING, STATUS_ERROR, STATUS_SUCCESS } from '../state/query';
import './HistoryTable.less';

class HistoryTable extends Component {
  render() {
    const { items, onOpenQuery, handleReset } = this.props;
    const columns = [
      { Header: 'query', accessor: 'sql' },
      {
        Header: 'log',
        id: 'log',
        accessor: row => {
          switch (row.status) {
            case STATUS_SUCCESS:
              return `Query ran in ${row.elapsedTime / 1000}sec`;
            case STATUS_ERROR:
              return `Query failed - ${row.errorMsg}`;
            case STATUS_LOADING:
              return 'Query is running';
            default:
              return 'unknown status';
          }
        }
      },
      {
        Header: 'date',
        accessor: 'datetime',
        Cell: ({ value }) => value.toString()
      },
      {
        Header: 'action',
        accessor: 'sql',
        Cell: ({ value }) => (
          <a onClick={() => onOpenQuery(value)}>open query</a>
        )
      }
    ];

    return (
      <div className="history">
        <div className="toolbar">
          <a onClick={handleReset}>reset history</a>
        </div>
        <ReactTable
          className="results-table"
          data={items}
          columns={columns}
          defaultPageSize={10}
        />
      </div>
    );
  }
}

HistoryTable.propTypes = {
  items: PropTypes.arrayOf(
    PropTypes.shape({
      sql: PropTypes.string,
      status: PropTypes.string,
      datetime: PropTypes.object,
      elapsedTime: PropTypes.number
    })
  ),
  onOpenQuery: PropTypes.func.isRequired,
  handleReset: PropTypes.func.isRequired
};

export default HistoryTable;
