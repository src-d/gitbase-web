import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { transformer } from 'uast-viewer';
import UASTViewerPane from './UASTViewerPane';
import api from '../api';
import CloseIcon from '../icons/close-query-tab.svg';

class UASTViewer extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loading: false,
      uast: transformer({
        InternalType: 'Search results',
        Children: props.uast
      }),
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
    this.setState({ uast: null, error: null, loading: true });

    api
      .filterUAST(this.props.protobufs, this.state.filter)
      .then(uast => {
        this.setState({ uast: transformer(uast) });
      })
      .catch(err => this.setState({ uast: null, error: err }))
      .then(() => this.setState({ loading: false }));
  }

  removeError() {
    this.setState({ error: null });
  }

  render() {
    const { uast, error, loading } = this.state;
    const { showLocations, filter } = this.state;
    const uastViewerProps = { uast };

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
