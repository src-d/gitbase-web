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
import SuccessIcon from '../icons/success-query.svg';
import ErrorIcon from '../icons/error-query.svg';

function ResultInfo({ result }) {
  if (!result) {
    return null;
  }

  if (result.response && result.response.meta) {
    return (
      <span className="meta meta-success">
        <SuccessIcon className="big-icon" />Showing rows (query took{' '}
        {result.response.meta.elapsedTime / 1000} seconds)
      </span>
    );
  }

  if (result.errorMsg) {
    return (
      <span className="meta meta-error">
        <ErrorIcon className="big-icon" />Query Failed - {result.errorMsg}
      </span>
    );
  }

  return null;
}

ResultInfo.propTypes = {
  result: PropTypes.object
};

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
    const { result } = this.props;
    const { codeMirrorTables } = this.state;

    const options = {
      mode: 'text/x-mariadb',
      smartIndent: true,
      lineNumbers: false,
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
      <div className="query-box-padding full-height full-width">
        <div className="query-box full-height full-width">
          <Row className="codemirror-row no-spacing">
            <Col xs={12} className="codemirror-col no-spacing">
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
            <Col xs={7} className="meta-wrapper no-spacing">
              <ResultInfo result={result} />
            </Col>
            <Col xs={5} className="buttons-wrapper no-spacing">
              <Button
                bsStyle="gbpl-secondary-tint-2-link"
                disabled={!this.props.exportUrl}
                href={this.props.exportUrl}
                target="_blank"
              >
                EXPORT
              </Button>
              <Button
                className="run-query"
                bsStyle="gbpl-secondary"
                disabled={this.props.enabled === false}
                onClick={this.props.handleSubmit}
              >
                RUN QUERY
              </Button>
            </Col>
          </Row>
        </div>
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
  result: PropTypes.object,
  enabled: PropTypes.bool,
  handleTextChange: PropTypes.func.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  exportUrl: PropTypes.string
};

export default QueryBox;
