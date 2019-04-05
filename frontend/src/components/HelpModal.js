import React, { Component } from 'react';
import { Modal } from 'react-bootstrap';
import PropTypes from 'prop-types';
import { ReactComponent as CloseIcon } from '../icons/close-query-tab.svg';
import './HelpModal.less';

class HelpModal extends Component {
  render() {
    const { showModal, onHide } = this.props;

    return (
      <Modal show={showModal} onHide={onHide}>
        <Modal.Header>
          <Modal.Title>
            Help
            <CloseIcon className="btn-modal-close" onClick={onHide} />
          </Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <div>
            For a reference of gitbase SQL, please read{' '}
            <strong>
              <a
                href="https://docs.sourced.tech/gitbase"
                target="_blank"
                rel="noopener noreferrer"
              >
                docs.sourced.tech/gitbase
              </a>
            </strong>
          </div>
          <h4 className="kb-shortcuts-title">Keyboard shortcuts</h4>
          <table className="keyboard-shortcuts">
            <tbody>
              <tr>
                <td className="keys">
                  <div className="key">Ctrl</div>
                  &nbsp;+&nbsp;
                  <div className="key">Enter</div>
                </td>
                <td>Run the query</td>
              </tr>
              <tr>
                <td className="keys">
                  <div className="key">Ctrl</div>
                  &nbsp;+&nbsp;
                  <div className="key">Space</div>
                </td>
                <td>Autocomplete</td>
              </tr>
            </tbody>
          </table>
        </Modal.Body>
      </Modal>
    );
  }
}

HelpModal.propTypes = {
  showModal: PropTypes.bool.isRequired,
  onHide: PropTypes.func.isRequired
};

export default HelpModal;
