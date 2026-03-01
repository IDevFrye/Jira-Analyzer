import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import MyProjectsPage from './MyProjectsPage';
import axios from 'axios';
import { config } from '../../config/config';

jest.mock('axios');

const mockedAxios = axios as jest.Mocked<typeof axios>;

const mockProjectsResponse = [
  { id: '1', key: 'PRJ1', name: 'Alpha', self: 'url1' },
  { id: '2', key: 'PRJ2', name: 'Beta', self: 'url2' },
];

const mockStatsResponse = {
  id: '1',
  key: 'PRJ1',
  name: 'Alpha',
  total_issues: 10,
  open_issues: 5,
  closed_issues: 5,
  reopened_issues: 1,
  resolved_issues: 4,
  in_progress_issues: 2,
  avg_resolution_time_h: 12,
  avg_created_per_day_7d: 1.5,
};

describe('MyProjectsPage', () => {
  beforeEach(() => {
    mockedAxios.get.mockReset();
  });

  const setupAxiosMock = () => {
    mockedAxios.get.mockImplementation((url: string) => {
      if (url === config.api.endpoints.projects) {
        return Promise.resolve({ data: mockProjectsResponse });
      }
      // project stats запросы для StatsCard
      return Promise.resolve({ data: mockStatsResponse });
    });
  };

  test('загружает и отображает сохранённые проекты', async () => {
    setupAxiosMock();

    render(<MyProjectsPage />);

    await waitFor(() => {
      expect(screen.getByText('Alpha')).toBeInTheDocument();
      expect(screen.getByText('Beta')).toBeInTheDocument();
    });
  });

  test('фильтрация по поиску', async () => {
    setupAxiosMock();

    render(<MyProjectsPage />);

    await waitFor(() => {
      expect(screen.getByText('Alpha')).toBeInTheDocument();
    });

    const input = screen.getByPlaceholderText('Поиск проектов...');
    fireEvent.change(input, { target: { value: 'Beta' } });

    await waitFor(() => {
      expect(screen.queryByText('Alpha')).not.toBeInTheDocument();
      expect(screen.getByText('Beta')).toBeInTheDocument();
    });
  });
});

