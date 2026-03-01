import React from 'react';
import { render } from '@testing-library/react';
import JiraAnalyzerApp from './JiraAnalyzerApp';

jest.mock('./router', () => () => <div data-testid="router-mock" />);
jest.mock('./components/Layout/Layout', () => ({ children }: { children: React.ReactNode }) => (
  <div data-testid="layout-mock">{children}</div>
));

describe('JiraAnalyzerApp', () => {
  test('рендерит Layout и Router внутри BrowserRouter', () => {
    const { getByTestId } = render(<JiraAnalyzerApp />);
    expect(getByTestId('layout-mock')).toBeInTheDocument();
    expect(getByTestId('router-mock')).toBeInTheDocument();
  });
});

