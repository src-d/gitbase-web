import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Glyphicon } from 'react-bootstrap';
import './Schema.less';

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
    const glyph = this.state.expanded ? 'minus' : 'plus';

    return (
      <div className="schema-table list">
        <div className="name">
          <Glyphicon glyph={glyph} onClick={this.toggle} />
          <span onClick={() => onTableClick && onTableClick(table)}>
            {table}
          </span>
        </div>
        {this.state.expanded && (
          <div className="columns">
            {columns.map((c, i) => (
              <div key={i} className="column">
                <Glyphicon glyph="align-justify" />
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
