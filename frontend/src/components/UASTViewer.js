import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { transformer } from 'uast-viewer';
import UASTViewerPane from './UASTViewerPane';
import api from '../api';

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

  render() {
    const { uast, error, loading } = this.state;
    const { showLocations, filter } = this.state;
    const uastViewerProps = { uast };

    if (error) {
      return <div>{error}</div>;
    }

    return (
      <UASTViewerPane
        loading={loading}
        uastViewerProps={uastViewerProps}
        showLocations={showLocations}
        filter={filter}
        handleShowLocationsChange={this.handleShowLocationsChange}
        handleFilterChange={this.handleFilterChange}
        handleSearch={this.handleSearch}
      />
    );
  }
}

UASTViewer.propTypes = {
  uast: PropTypes.array,
  protobufs: PropTypes.string
};

export default UASTViewer;
