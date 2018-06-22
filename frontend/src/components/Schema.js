import React, { Component } from 'react';
import PropTypes from 'prop-types';
import './Schema.less';
import OpenIcon from '../icons/open-tree.svg';
import CloseIcon from '../icons/close-tree.svg';
import ColumnIcon from '../icons/tree-column.svg';

class SchemaTable extends Component {
  constructor(props) {
    super(props);

    this.state = { expanded: false };

    this.toggle = this.toggle.bind(this);
  }

  toggle() {
    this.setState({ expanded: !this.state.expanded });
  }

  render() {
    const { table, columns, onTableClick } = this.props;
    const icon = this.state.expanded ? (
      <CloseIcon className="small-icon" onClick={this.toggle} />
    ) : (
      <OpenIcon className="small-icon" onClick={this.toggle} />
    );

    return (
      <div className="schema-table list">
        <div className="name">
          {icon}
          <span onClick={() => onTableClick && onTableClick(table)}>
            {table}
          </span>
        </div>
        {this.state.expanded && (
          <div className="columns">
            {columns.map((c, i) => (
              <div key={i} className="column">
                <ColumnIcon className="small-icon" />
                {c.name}
              </div>
            ))}
          </div>
        )}
      </div>
    );
  }
}

SchemaTable.propTypes = {
  table: PropTypes.string.isRequired,
  columns: PropTypes.arrayOf(
    PropTypes.shape({
      name: PropTypes.string.isRequired,
      type: PropTypes.string.isRequired
    })
  ).isRequired,
  onTableClick: PropTypes.func
};

function Schema({ schema, onTableClick }) {
  if (!schema) {
    return null;
  }

  return (
    <div className="schema">
      {schema.map(item => (
        <SchemaTable key={item.table} {...item} onTableClick={onTableClick} />
      ))}
    </div>
  );
}

Schema.propTypes = {
  schema: PropTypes.arrayOf(
    PropTypes.shape({
      table: SchemaTable.propTypes.table,
      columns: SchemaTable.propTypes.columns
    })
  ),
  onTableClick: PropTypes.func
};

export default Schema;
