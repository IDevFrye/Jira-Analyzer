import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import ComparePage from './ComparePage';
import axios from 'axios';

jest.mock('axios');
jest.mock('../../components/CompareSelector/CompareSelector', () => ({
  __esModule: true,
  default: ({ projects, selectedProjects, onSelect }: any) => (
    <div>
      {projects.map((p: any) => (
        <button
          key={p.Id}
          onClick={() => onSelect(p)}
          data-testid={`project-${p.Id}`}
        >
          {p.Name}
        </button>
      ))}
      <div data-testid="selected-count">{selectedProjects.length}</div>
    </div>
  ),
}));

jest.mock('../../components/CompareModal/CompareModal', () => ({
  __esModule: true,
  default: ({ onClose }: any) => (
    <div data-testid="compare-modal">
      Modal
      <button onClick={onClose}>close</button>
    </div>
  ),
}));

const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('ComparePage', () => {
  test('загружает проекты и позволяет выбрать до трёх для сравнения', async () => {
    mockedAxios.get.mockResolvedValueOnce({
      data: [
        { id: '1', key: 'PRJ1', name: 'Alpha', self: 'url1' },
        { id: '2', key: 'PRJ2', name: 'Beta', self: 'url2' },
        { id: '3', key: 'PRJ3', name: 'Gamma', self: 'url3' },
        { id: '4', key: 'PRJ4', name: 'Delta', self: 'url4' },
      ],
    });

    render(<ComparePage />);

    await waitFor(() => {
      expect(screen.getByText('Alpha')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('project-1'));
    fireEvent.click(screen.getByTestId('project-2'));

    await waitFor(() => {
      expect(screen.getByText(/Сравнить выбранные проекты/)).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText(/Сравнить выбранные проекты/));
    expect(screen.getByTestId('compare-modal')).toBeInTheDocument();
  });
});

