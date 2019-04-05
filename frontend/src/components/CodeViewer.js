import React, { Component } from 'react';
import { Modal } from 'react-bootstrap';
import PropTypes from 'prop-types';
import SplitPane from 'react-split-pane';
import { Editor, withUASTEditor, languageToMode } from 'uast-viewer';
import Switch from 'react-switch';
import UASTViewerPane from './UASTViewerPane';
import api from '../api';
import './CodeViewer.less';
import { ReactComponent as CloseIcon } from '../icons/close-query-tab.svg';

function EditorPane({ languages, language, handleLangChange, editorProps }) {
  return (
    <div className="editor-pane">
      <div className="language-selection">
        <span className="lang-label">LANGUAGE</span>
        <select value={language} onChange={handleLangChange}>
          <option value="">Select language</option>
          {languages.map(lang => (
            <option key={lang.id} value={lang.id}>
              {lang.name}
            </option>
          ))}
        </select>
      </div>
      <Editor
        {...editorProps}
        languageMode={languageToMode(language)}
        theme="default"
      />
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
  handleLangChange: PropTypes.func.isRequired,
  editorProps: PropTypes.object
};

function EditorUASTSpitPane({
  uastLoading,
  languages,
  editorProps,
  uastViewerProps,
  showLocations,
  filter,
  handleLangChange,
  handleShowLocationsChange,
  handleFilterChange,
  handleSearch,
  mode,
  handleModeChange
}) {
  return (
    <SplitPane split="vertical" defaultSize={500} minSize={1} maxSize={-15}>
      <EditorPane
        languages={languages}
        language={editorProps.languageMode}
        handleLangChange={handleLangChange}
        editorProps={editorProps}
      />
      <UASTViewerPane
        loading={uastLoading}
        uastViewerProps={uastViewerProps}
        showLocations={showLocations}
        filter={filter}
        handleShowLocationsChange={handleShowLocationsChange}
        handleFilterChange={handleFilterChange}
        handleSearch={handleSearch}
        mode={mode}
        handleModeChange={handleModeChange}
      />
    </SplitPane>
  );
}

EditorUASTSpitPane.propTypes = {
  uastLoading: PropTypes.bool,
  languages: EditorPane.propTypes.languages,
  editorProps: PropTypes.object,
  uastViewerProps: PropTypes.object,
  showLocations: PropTypes.bool,
  filter: PropTypes.string,
  handleLangChange: PropTypes.func.isRequired,
  handleShowLocationsChange: PropTypes.func.isRequired,
  handleFilterChange: PropTypes.func.isRequired,
  handleSearch: PropTypes.func.isRequired,
  mode: PropTypes.string.isRequired,
  handleModeChange: PropTypes.func.isRequired
};

const EditorWithUAST = withUASTEditor(EditorUASTSpitPane);

class CodeViewer extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loading: true,
      language: null,
      uastLoading: false,
      showUast: false,
      uast: null,
      error: null,
      showLocations: false,
      filter: '',
      mode: api.defaultUastMode
    };

    this.handleLangChange = this.handleLangChange.bind(this);
    this.handleShowUastChange = this.handleShowUastChange.bind(this);
    this.parseCode = this.parseCode.bind(this);
    this.removeError = this.removeError.bind(this);
    this.handleShowLocationsChange = this.handleShowLocationsChange.bind(this);
    this.handleFilterChange = this.handleFilterChange.bind(this);
    this.handleModeChange = this.handleModeChange.bind(this);
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
        // eslint-disable-next-line no-console
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
    this.setState({ error: null, uast: null, uastLoading: true });

    api
      .parseCode(
        this.state.language,
        this.props.code,
        this.state.mode,
        this.state.filter
      )
      .then(res => {
        this.setState({ uast: res });
      })
      .catch(error => {
        this.setState({ error });
      })
      .then(() => this.setState({ uastLoading: false }));
  }

  removeError() {
    this.setState({ error: null });
  }

  handleShowLocationsChange() {
    this.setState({ showLocations: !this.state.showLocations });
  }

  handleFilterChange(e) {
    this.setState({ filter: e.target.value });
  }

  handleModeChange(mode) {
    this.setState({ mode }, () => {
      if (this.state.showUast) {
        this.parseCode();
      }
    });
  }

  render() {
    const { showModal, onHide, code, languages } = this.props;
    const {
      loading,
      language,
      showUast,
      uastLoading,
      uast,
      error,
      showLocations,
      filter,
      mode
    } = this.state;

    if (loading) {
      return 'loading';
    }

    return (
      <Modal show={showModal} onHide={onHide} bsSize="large">
        <Modal.Header>
          <Modal.Title>
            CODE
            <Switch
              checked={showUast}
              onChange={this.handleShowUastChange}
              disabled={!language}
              checkedIcon={<span className="switch-text checked">UAST</span>}
              uncheckedIcon={
                <span className="switch-text unchecked">UAST</span>
              }
              width={100}
              handleDiameter={20}
              className={`code-toggler ${showUast ? 'checked' : 'unchecked'}`}
              aria-label="Toggle UAST view"
            />
            <CloseIcon className="btn-modal-close" onClick={onHide} />
          </Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {showUast ? (
            <div className="code-viewer">
              <EditorWithUAST
                languages={languages}
                code={code}
                languageMode={language}
                showUast={showUast}
                uastLoading={uastLoading}
                uast={uast}
                showLocations={showLocations}
                filter={filter}
                handleLangChange={this.handleLangChange}
                handleShowLocationsChange={this.handleShowLocationsChange}
                handleFilterChange={this.handleFilterChange}
                handleSearch={this.parseCode}
                mode={mode}
                handleModeChange={this.handleModeChange}
              />
              {error ? (
                <div className="error">
                  <div className="error-header">
                    <span>ERROR</span>
                    <CloseIcon
                      className="btn-error-close"
                      onClick={this.removeError}
                    />
                  </div>
                  <div className="error-msg">{error}</div>
                </div>
              ) : null}
            </div>
          ) : (
            <EditorPane
              languages={languages}
              language={language}
              handleLangChange={this.handleLangChange}
              editorProps={{ code }}
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
