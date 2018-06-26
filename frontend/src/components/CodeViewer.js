import React, { Component } from 'react';
import { Modal, Button } from 'react-bootstrap';
import PropTypes from 'prop-types';
import SplitPane from 'react-split-pane';
import UASTViewer, { Editor, withUASTEditor } from 'uast-viewer';
import Switch from 'react-switch';
import api from '../api';
import './CodeViewer.less';

function EditorPane({ languages, language, handleLangChange, editorProps }) {
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
  handleLangChange: PropTypes.func.isRequired,
  editorProps: PropTypes.object
};

const ROOT_ID = 1;
const SEARCH_RESULTS_TYPE = 'Search results';

function getSearchResults(uast) {
  if (!uast) {
    return null;
  }

  const rootNode = uast[ROOT_ID];
  if (!rootNode) {
    return null;
  }

  if (rootNode.InternalType === SEARCH_RESULTS_TYPE) {
    return rootNode.Children;
  }

  return null;
}

function NotFound() {
  return <div>Nothing found</div>;
}

function UASTViewerPane({
  uastViewerProps,
  showLocations,
  useCustomServer,
  customServer,
  filter,
  handleShowLocationsChange,
  handleUseCustomServerChange,
  handleCustomServerChange,
  handleFilterChange,
  handleSearch
}) {
  const searchResults = getSearchResults(uastViewerProps.uast);
  const rootIds = searchResults || [ROOT_ID];

  let content = null;
  if (uastViewerProps.uast) {
    if (searchResults && !searchResults.length) {
      content = <NotFound />;
    } else {
      content = (
        <UASTViewer
          {...uastViewerProps}
          rootIds={rootIds}
          showLocations={showLocations}
        />
      );
    }
  }

  return (
    <div>
      <div>
        <label>
          <input
            type="checkbox"
            checked={showLocations}
            onChange={handleShowLocationsChange}
          />Show locations
        </label>
        <label>
          <input
            type="checkbox"
            value={useCustomServer}
            onChange={handleUseCustomServerChange}
          />Custom bblfsh server
        </label>
        {useCustomServer ? (
          <input
            type="text"
            value={customServer}
            onChange={handleCustomServerChange}
          />
        ) : null}
      </div>
      <div>
        <form
          onSubmit={e => {
            e.preventDefault();
            handleSearch();
          }}
        >
          <input
            type="text"
            placeholder="UAST Query"
            value={filter}
            onChange={handleFilterChange}
          />{' '}
          <Button type="submit">Search</Button>{' '}
          <Button
            href="https://doc.bblf.sh/using-babelfish/uast-querying.html"
            target="_blank"
          >
            Help
          </Button>
        </form>
      </div>
      {content}
    </div>
  );
}

UASTViewerPane.propTypes = {
  uastViewerProps: PropTypes.object,
  showLocations: PropTypes.bool,
  useCustomServer: PropTypes.bool,
  customServer: PropTypes.string,
  filter: PropTypes.string,
  handleShowLocationsChange: PropTypes.func.isRequired,
  handleUseCustomServerChange: PropTypes.func.isRequired,
  handleCustomServerChange: PropTypes.func.isRequired,
  handleFilterChange: PropTypes.func.isRequired,
  handleSearch: PropTypes.func.isRequired
};

function EditorUASTSpitPane({
  languages,
  editorProps,
  uastViewerProps,
  showLocations,
  useCustomServer,
  customServer,
  filter,
  handleLangChange,
  handleShowLocationsChange,
  handleUseCustomServerChange,
  handleCustomServerChange,
  handleFilterChange,
  handleSearch
}) {
  return (
    <SplitPane split="vertical" defaultSize={250} minSize={175}>
      <EditorPane
        languages={languages}
        language={editorProps.languageMode}
        handleLangChange={handleLangChange}
        editorProps={editorProps}
      />
      <UASTViewerPane
        uastViewerProps={uastViewerProps}
        showLocations={showLocations}
        useCustomServer={useCustomServer}
        customServer={customServer}
        filter={filter}
        handleShowLocationsChange={handleShowLocationsChange}
        handleUseCustomServerChange={handleUseCustomServerChange}
        handleCustomServerChange={handleCustomServerChange}
        handleFilterChange={handleFilterChange}
        handleSearch={handleSearch}
      />
    </SplitPane>
  );
}

EditorUASTSpitPane.propTypes = {
  languages: EditorPane.propTypes.languages,
  editorProps: PropTypes.object,
  uastViewerProps: PropTypes.object,
  showLocations: PropTypes.bool,
  useCustomServer: PropTypes.bool,
  customServer: PropTypes.string,
  filter: PropTypes.string,
  handleLangChange: PropTypes.func.isRequired,
  handleShowLocationsChange: PropTypes.func.isRequired,
  handleUseCustomServerChange: PropTypes.func.isRequired,
  handleCustomServerChange: PropTypes.func.isRequired,
  handleFilterChange: PropTypes.func.isRequired,
  handleSearch: PropTypes.func.isRequired
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
      error: null,
      showLocations: false,
      useCustomServer: false,
      customServer: '0.0.0.0:9432',
      filter: ''
    };

    this.handleLangChange = this.handleLangChange.bind(this);
    this.handleShowUastChange = this.handleShowUastChange.bind(this);
    this.parseCode = this.parseCode.bind(this);
    this.removeError = this.removeError.bind(this);
    this.handleShowLocationsChange = this.handleShowLocationsChange.bind(this);
    this.handleUseCustomServerChange = this.handleUseCustomServerChange.bind(
      this
    );
    this.handleCustomServerChange = this.handleCustomServerChange.bind(this);
    this.handleFilterChange = this.handleFilterChange.bind(this);
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
      .parseCode(
        this.state.language,
        this.props.code,
        this.state.filter,
        this.state.useCustomServer ? this.state.customServer : undefined
      )
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

  handleShowLocationsChange() {
    this.setState({ showLocations: !this.state.showLocations });
  }

  handleUseCustomServerChange() {
    this.setState({ useCustomServer: !this.state.useCustomServer });
  }

  handleCustomServerChange(e) {
    this.setState({ customServer: e.target.value });
  }

  handleFilterChange(e) {
    this.setState({ filter: e.target.value });
  }

  render() {
    const { showModal, onHide, code, languages } = this.props;
    const {
      loading,
      language,
      showUast,
      uast,
      error,
      showLocations,
      useCustomServer,
      customServer,
      filter
    } = this.state;

    if (loading) {
      return 'loading';
    }

    return (
      <Modal show={showModal} onHide={onHide} bsSize="large">
        <Modal.Header closeButton>
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
                uast={uast}
                showLocations={showLocations}
                useCustomServer={useCustomServer}
                customServer={customServer}
                filter={filter}
                handleLangChange={this.handleLangChange}
                handleShowLocationsChange={this.handleShowLocationsChange}
                handleUseCustomServerChange={this.handleUseCustomServerChange}
                handleCustomServerChange={this.handleCustomServerChange}
                handleFilterChange={this.handleFilterChange}
                handleSearch={this.parseCode}
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
              handleLangChange={this.handleLangChange}
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
