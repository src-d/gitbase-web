import React from 'react';
import renderer from 'react-test-renderer';
import ResultsTable from './ResultsTable';

describe('ResultsTable', () => {
  const noop = () => null;

  it('text with new lines should be shown as code', () => {
    const response = {
      meta: {
        headers: ['notCode', 'code'],
        types: ['TEXT', 'TEXT']
      },
      data: [{ notCode: 'not a code', code: 'this\nis\ncode' }]
    };
    const tree = renderer
      .create(
        <ResultsTable response={response} showCode={noop} showUAST={noop} />
      )
      .toJSON();

    expect(tree).toMatchSnapshot();
  });
});
