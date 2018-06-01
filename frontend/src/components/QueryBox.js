import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Row, Col, Button } from 'react-bootstrap';
import { Controlled as CodeMirror } from 'react-codemirror2';

import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/sql/sql';
import 'codemirror/addon/display/placeholder';
import 'codemirror/addon/edit/matchbrackets';
import 'codemirror/addon/hint/show-hint.css';
import 'codemirror/addon/hint/show-hint';
import 'codemirror/addon/hint/sql-hint';

import './QueryBox.less';

class QueryBox extends Component {
  constructor(props) {
    super(props);

    this.state = {
      schema: undefined,
      codeMirrorTables: {}
    };
  }

  static getDerivedStateFromProps(nextProps, prevState) {
    if (nextProps.schema === prevState.schema) {
      return null;
    }

    return {
      schema: nextProps.schema,
      codeMirrorTables: QueryBox.schemaToCodeMirror(nextProps.schema)
    };
  }

  static schemaToCodeMirror(schema) {
    if (!schema) {
      return {};
    }

    return schema.reduce(
      (prevVal, currVal) => ({
        ...prevVal,
        [currVal.table]: currVal.columns.map(col => col.name)
      }),
      {}
    );
  }

  render() {
    const { codeMirrorTables } = this.state;

    const options = {
      mode: 'text/x-mariadb',
      smartIndent: true,
      lineNumbers: true,
      matchBrackets: true,
      autofocus: true,
      placeholder: 'Enter an SQL query',
      extraKeys: {
        'Ctrl-Space': 'autocomplete',
        'Ctrl-Enter': () => this.props.handleSubmit()
      },
      hintOptions: {
        tables: codeMirrorTables
      }
    };

    return (
      <div className="query-box">
        <Row className="codemirror-row">
          <Col xs={12} className="codemirror-col">
            <CodeMirror
              value={this.props.sql}
              options={options}
              onBeforeChange={(editor, data, value) => {
                this.props.handleTextChange(value);
              }}
            />
          </Col>
        </Row>
        <Row className="button-row">
          <Col xsOffset={6} xs={3}>
            <Button
              className="pull-right"
              disabled={!this.props.exportUrl}
              href={this.props.exportUrl}
              target="_blank"
            >
              EXPORT
            </Button>
          </Col>
          <Col xs={3}>
            <Button
              className="pull-right"
              disabled={this.props.enabled === false}
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
  schema: PropTypes.arrayOf(
    PropTypes.shape({
      table: PropTypes.string.isRequired,
      columns: PropTypes.arrayOf(
        PropTypes.shape({
          name: PropTypes.string.isRequired,
          type: PropTypes.string.isRequired
        })
      ).isRequired
    })
  ),
  enabled: PropTypes.bool,
  handleTextChange: PropTypes.func.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  exportUrl: PropTypes.string
};

export default QueryBox;
