import React, { Component } from 'react';
import { Modal } from 'react-bootstrap';
import PropTypes from 'prop-types';
import SplitPane from 'react-split-pane';
import UASTViewer, { Editor, withUASTEditor } from 'uast-viewer';
import api from '../api';
import './CodeViewer.less';

function EditorPane({
  languages,
  language,
  showUast,
  handleLangChange,
  handleShowUastChange,
  editorProps
}) {
  return (
    <div className="editor-pane">
      Language:{' '}
      <select value={language} onChange={handleLangChange}>
        <option value="">Select language</option>
        {languages.map(lang => (
          <option key={lang.id} value={lang.id}>
            {lang.name}
          </option>
        ))}
      </select>
      <label>
        <input
          type="checkbox"
          checked={showUast}
          onChange={handleShowUastChange}
          disabled={!language}
        />UAST
      </label>
      <Editor {...editorProps} theme="default" />
    </div>
  );
}

EditorPane.propTypes = {
  languages: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      name: PropTypes.string.isRequired
    })
  ).isRequired,
  language: PropTypes.string,
  showUast: PropTypes.bool,
  handleLangChange: PropTypes.func.isRequired,
  handleShowUastChange: PropTypes.func.isRequired,
  editorProps: PropTypes.object
};

function EditorUASTSpitPane({
  languages,
  editorProps,
  uastViewerProps,
  showUast,
  handleLangChange,
  handleShowUastChange
}) {
  return (
    <SplitPane split="vertical" defaultSize={250} minSize={175}>
      <EditorPane
        languages={languages}
        language={editorProps.languageMode}
        showUast={showUast}
        handleLangChange={handleLangChange}
        handleShowUastChange={handleShowUastChange}
        editorProps={editorProps}
      />
      {uastViewerProps.uast ? <UASTViewer {...uastViewerProps} /> : <div />}
    </SplitPane>
  );
}

EditorUASTSpitPane.propTypes = {
  languages: EditorPane.propTypes.languages,
  editorProps: PropTypes.object,
  uastViewerProps: PropTypes.object,
  showUast: PropTypes.bool,
  handleLangChange: PropTypes.func.isRequired,
  handleShowUastChange: PropTypes.func.isRequired
};

const EditorWithUAST = withUASTEditor(EditorUASTSpitPane);

class CodeViewer extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loading: true,
      language: null,
      showUast: false,
      uast: null,
      error: null
    };

    this.handleLangChange = this.handleLangChange.bind(this);
    this.handleShowUastChange = this.handleShowUastChange.bind(this);
    this.removeError = this.removeError.bind(this);
  }

  componentDidMount() {
    api
      .detectLang(this.props.code)
      .then(res => {
        this.setState({ language: res.language });
      })
      .catch(err => {
        // we don't have UI for this error and actually it's not very important
        // user can select language manualy
        console.error(`can't detect language: ${err}`);
      })
      .then(() => this.setState({ loading: false }));
  }

  handleLangChange(e) {
    this.setState({ language: e.target.value }, () => {
      if (!this.state.language) {
        this.setState({ showUast: false });
        return;
      }

      if (this.state.showUast) {
        this.parseCode();
      }
    });
  }

  handleShowUastChange() {
    const showUast = !this.state.showUast;
    if (showUast) {
      this.parseCode();
    }
    this.setState({ showUast });
  }

  parseCode() {
    this.setState({ error: null, uast: null });

    api
      .parseCode(this.state.language, this.props.code)
      .then(res => {
        this.setState({ uast: res });
      })
      .catch(error => {
        this.setState({ error });
      });
  }

  removeError() {
    this.setState({ error: null });
  }

  render() {
    const { loading, language, showUast, uast, error } = this.state;
    const { showModal, onHide, code, languages } = this.props;

    if (loading) {
      return 'loading';
    }

    return (
      <Modal show={showModal} onHide={onHide} bsSize="large">
        <Modal.Header closeButton>
          <Modal.Title>CODE</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {showUast ? (
            <div className="code-viewer">
              <EditorWithUAST
                languages={languages}
                code={code}
                languageMode={language}
                showUast={showUast}
                handleLangChange={this.handleLangChange}
                handleShowUastChange={this.handleShowUastChange}
                uast={uast}
              />
              {error ? (
                <div className="error">
                  <button onClick={this.removeError} className="close">
                    close
                  </button>
                  {error}
                </div>
              ) : null}
            </div>
          ) : (
            <EditorPane
              languages={languages}
              language={language}
              showUast={showUast}
              handleLangChange={this.handleLangChange}
              handleShowUastChange={this.handleShowUastChange}
              editorProps={{ code, languageMode: language }}
            />
          )}
        </Modal.Body>
      </Modal>
    );
  }
}

CodeViewer.propTypes = {
  code: PropTypes.string,
  languages: EditorPane.propTypes.languages,
  showModal: PropTypes.bool.isRequired,
  onHide: PropTypes.func.isRequired
};

export default CodeViewer;
