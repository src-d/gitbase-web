import React from 'react';
import PropTypes from 'prop-types';
import UASTViewer from 'uast-viewer';
import { Button } from 'react-bootstrap';
import './UASTViewerPane.less';

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
  loading,
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
  let content = null;

  if (loading) {
    content = <div>loading...</div>;
  } else if (uastViewerProps.uast) {
    const searchResults = getSearchResults(uastViewerProps.uast);
    const rootIds = searchResults || [ROOT_ID];

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
        {typeof useCustomServer !== 'undefined' ? (
          <label>
            <input
              type="checkbox"
              value={useCustomServer}
              onChange={handleUseCustomServerChange}
            />
            <span>Custom bblfsh server</span>
          </label>
        ) : null}
        {useCustomServer ? (
          <input
            type="text"
            value={customServer}
            onChange={handleCustomServerChange}
          />
        ) : null}
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
            className="edit-query"
            bsStyle="gbpl-primary-tint-2-link"
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
  loading: PropTypes.bool,
  uastViewerProps: PropTypes.object,
  showLocations: PropTypes.bool,
  useCustomServer: PropTypes.bool,
  customServer: PropTypes.string,
  filter: PropTypes.string,
  handleShowLocationsChange: PropTypes.func.isRequired,
  handleUseCustomServerChange: PropTypes.func,
  handleCustomServerChange: PropTypes.func,
  handleFilterChange: PropTypes.func.isRequired,
  handleSearch: PropTypes.func.isRequired
};

export default UASTViewerPane;
