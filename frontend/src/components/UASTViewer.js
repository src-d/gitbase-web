import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { expandRootIds, uastV2 } from 'uast-viewer';
import UASTViewerPane from './UASTViewerPane';
import api from '../api';
import { ReactComponent as CloseIcon } from '../icons/close-query-tab.svg';

// Same values as the ones applied by withUASTEditor in CodeViewer.js
// https://github.com/bblfsh/uast-viewer/blob/v0.2.0/src/withUASTEditor.js#L208
const ROOT_IDS = [1];
const LEVELS_EXPAND = 2;

class UASTViewer extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loading: false,
      initialFlatUast: this.transform(props.uast),
      showLocations: false,
      filter: '',
      error: null
    };

    this.handleShowLocationsChange = this.handleShowLocationsChange.bind(this);
    this.handleFilterChange = this.handleFilterChange.bind(this);
    this.handleSearch = this.handleSearch.bind(this);
    this.removeError = this.removeError.bind(this);
  }

  handleShowLocationsChange() {
    this.setState({ showLocations: !this.state.showLocations });
  }

  handleFilterChange(e) {
    this.setState({ filter: e.target.value });
  }

  handleSearch() {
    this.setState({ flatUast: null, error: null, loading: true });

    api
      .filterUAST(this.props.protobufs, this.state.filter)
      .then(uast => {
        this.setState({ flatUast: this.transform(uast) });
      })
      .catch(err => this.setState({ flatUast: null, error: err }))
      .then(() => this.setState({ loading: false }));
  }

  // Applies the uast-viewer object shape transformer, and expands the first
  // 2 levels
  transform(uast) {
    const flatUAST = uastV2.transformer(uast);
    return expandRootIds(
      flatUAST,
      ROOT_IDS,
      LEVELS_EXPAND,
      uastV2.getChildrenIds
    );
  }

  removeError() {
    this.setState({ error: null });
  }

  render() {
    const { initialFlatUast, error, loading } = this.state;
    const { showLocations, filter } = this.state;
    const uastViewerProps = { initialFlatUast };

    return (
      <div className="pg-uast-viewer">
        <UASTViewerPane
          loading={loading}
          uastViewerProps={uastViewerProps}
          showLocations={showLocations}
          filter={filter}
          handleShowLocationsChange={this.handleShowLocationsChange}
          handleFilterChange={this.handleFilterChange}
          handleSearch={this.handleSearch}
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
    );
  }
}

UASTViewer.propTypes = {
  uast: PropTypes.array,
  protobufs: PropTypes.string
};

export default UASTViewer;
