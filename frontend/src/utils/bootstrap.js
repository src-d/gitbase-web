import { Button } from 'react-bootstrap';
import { bootstrapUtils } from 'react-bootstrap/lib/utils';

function initButtonStyles() {
  bootstrapUtils.addStyle(Button, 'gbpl-secondary');
  bootstrapUtils.addStyle(Button, 'gbpl-secondary-tint-2-link');
  bootstrapUtils.addStyle(Button, 'gbpl-tertiary');
  bootstrapUtils.addStyle(Button, 'gbpl-tertiary-tint-2-link');
  bootstrapUtils.addStyle(Button, 'gbpl-primary-tint-2');
}

export default initButtonStyles;
