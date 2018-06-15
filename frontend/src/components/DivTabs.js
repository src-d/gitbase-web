/* eslint-disable */

/*
Copy-past from react-bootstrap
We can't just extend original component and redefine renderTab method because it exports uncontrollable HOC.

We need to use <div> instead of <a> as tab titles due to Firefox behaviour of <input> inside <a>.
It raise onBlur event on click. Example: http://jsfiddle.net/bc0u9hrq/5/
*/

import React from 'react';
import PropTypes from 'prop-types';
import requiredForA11y from 'prop-types-extra/lib/isRequiredForA11y';
import uncontrollable from 'uncontrollable';

import {
  Nav,
  NavItem,
  TabContainer as UncontrolledTabContainer,
  TabContent
} from 'react-bootstrap';
import { bsClass as setBsClass } from 'react-bootstrap/lib/utils/bootstrapUtils';
import ValidComponentChildren from 'react-bootstrap/lib/utils/ValidComponentChildren';

const TabContainer = UncontrolledTabContainer.ControlledComponent;

const propTypes = {
  /**
   * Mark the Tab with a matching `eventKey` as active.
   *
   * @controllable onSelect
   */
  activeKey: PropTypes.any,

  /**
   * Navigation style
   */
  bsStyle: PropTypes.oneOf(['tabs', 'pills']),

  animation: PropTypes.bool,

  id: requiredForA11y(
    PropTypes.oneOfType([PropTypes.string, PropTypes.number])
  ),

  /**
   * Callback fired when a Tab is selected.
   *
   * ```js
   * function (
   *   Any eventKey,
   *   SyntheticEvent event?
   * )
   * ```
   *
   * @controllable activeKey
   */
  onSelect: PropTypes.func,

  /**
   * Wait until the first "enter" transition to mount tabs (add them to the DOM)
   */
  mountOnEnter: PropTypes.bool,

  /**
   * Unmount tabs (remove it from the DOM) when it is no longer visible
   */
  unmountOnExit: PropTypes.bool
};

const defaultProps = {
  bsStyle: 'tabs',
  animation: true,
  mountOnEnter: false,
  unmountOnExit: false
};

function getDefaultActiveKey(children) {
  let defaultActiveKey;
  ValidComponentChildren.forEach(children, child => {
    if (defaultActiveKey == null) {
      defaultActiveKey = child.props.eventKey;
    }
  });

  return defaultActiveKey;
}

class DivTabs extends React.Component {
  renderTab(child) {
    const { title, eventKey, disabled, tabClassName } = child.props;
    if (title == null) {
      return null;
    }

    return (
      <NavItem
        eventKey={eventKey}
        disabled={disabled}
        className={tabClassName}
        componentClass="div"
      >
        {title}
      </NavItem>
    );
  }

  render() {
    const {
      id,
      onSelect,
      animation,
      mountOnEnter,
      unmountOnExit,
      bsClass,
      className,
      style,
      children,
      activeKey = getDefaultActiveKey(children),
      ...props
    } = this.props;

    return (
      <TabContainer
        id={id}
        activeKey={activeKey}
        onSelect={onSelect}
        className={className}
        style={style}
      >
        <div>
          <Nav {...props} role="tablist">
            {ValidComponentChildren.map(children, this.renderTab)}
          </Nav>

          <TabContent
            bsClass={bsClass}
            animation={animation}
            mountOnEnter={mountOnEnter}
            unmountOnExit={unmountOnExit}
          >
            {children}
          </TabContent>
        </div>
      </TabContainer>
    );
  }
}

DivTabs.propTypes = propTypes;
DivTabs.defaultProps = defaultProps;

setBsClass('tab', DivTabs);

export default uncontrollable(DivTabs, { activeKey: 'onSelect' });
