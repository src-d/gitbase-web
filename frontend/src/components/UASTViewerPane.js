import React from 'react';
import PropTypes from 'prop-types';
import FlatUASTViewer from 'uast-viewer';
import { Button } from 'react-bootstrap';
import './UASTViewerPane.less';
import ParseModeSwitcher from './ParseModeSwitcher';

const ROOT_ID = 1;

function getSearchResults(flatUast) {
  if (!flatUast) {
    return null;
  }

  const rootNode = flatUast[ROOT_ID];
  if (!rootNode) {
    return null;
  }

  if (Array.isArray(rootNode.n)) {
    return rootNode.n.map(c => c.id);
  }

  return null;
}

function NotFound() {
  return <div>Nothing found</div>;
}

function UASTViewerPane({
  loading,
  uastViewerProps,
  showLocations,
  filter,
  handleShowLocationsChange,
  handleFilterChange,
  handleSearch,
  mode,
  handleModeChange
}) {
  let content = null;

  if (loading) {
    content = <div>loading...</div>;
  } else if (uastViewerProps.flatUast) {
    const searchResults = getSearchResults(uastViewerProps.flatUast);
    const rootIds = searchResults || [ROOT_ID];

    if (searchResults && !searchResults.length) {
      content = <NotFound />;
    } else {
      content = (
        <FlatUASTViewer
          {...uastViewerProps}
          rootIds={rootIds}
          showLocations={showLocations}
        />
      );
    }
  }

  let modeSwitcher = null;
  if (mode) {
    modeSwitcher = (
      <div className="uast-mode-wrapper">
        <ParseModeSwitcher mode={mode} handleModeChange={handleModeChange} />
      </div>
    );
  }

  return (
    <div className="uast-viewer-pane">
      <div className="show-locations-wrapper">
        <label>
          <input
            type="checkbox"
            checked={showLocations}
            onChange={handleShowLocationsChange}
          />
          <span>Show locations</span>
        </label>
      </div>
      <div className="uast-query-wrapper">
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
          <Button bsStyle="gbpl-secondary" type="submit">
            SEARCH
          </Button>
          <Button
            bsStyle="gbpl-primary-tint-2-link"
            href="https://doc.bblf.sh/using-babelfish/uast-querying.html"
            target="_blank"
          >
            Help
          </Button>
        </form>
      </div>
      {modeSwitcher}
      {content}
    </div>
  );
}

UASTViewerPane.propTypes = {
  loading: PropTypes.bool,
  uastViewerProps: PropTypes.object,
  showLocations: PropTypes.bool,
  filter: PropTypes.string,
  handleShowLocationsChange: PropTypes.func.isRequired,
  handleFilterChange: PropTypes.func.isRequired,
  handleSearch: PropTypes.func.isRequired,
  // If mode is empty the mode selector will be hidden
  mode: PropTypes.string,
  // Mandatory if mode is set
  handleModeChange: PropTypes.func
};

export default UASTViewerPane;
