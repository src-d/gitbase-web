import React from 'react';
import { shallow } from 'enzyme';
import renderer from 'react-test-renderer';
import ParseModeSwitcher from './ParseModeSwitcher';

test('it renders correctly', () => {
  const tree = renderer
    .create(
      <ParseModeSwitcher mode={'semantic'} handleModeChange={() => null} />
    )
    .toJSON();

  expect(tree).toMatchSnapshot();
});

test('it calls handleModeChange when the input is clicked', () => {
  const spy = jest.fn();
  const wrapper = shallow(
    <ParseModeSwitcher mode={'semantic'} handleModeChange={spy} />
  );
  wrapper
    .find('[value="annotated"]')
    .simulate('change', { target: { value: 'annotated' } });
  expect(spy.mock.calls.length).toBe(1);
});
