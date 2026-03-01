import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import ProjectsPage from './ProjectsPage';
import axios from 'axios';
import { config } from '../../config/config';

jest.mock('axios');
jest.mock('../../config/config', () => ({
  config: {
    api: {
      endpoints: {
        connectorProjects: '/api/connector/projects',
        projects: '/api/projects',
        deleteProject: (id: number) => `/api/projects/${id}`,
        updateProject: '/api/connector/updateProject',
      },
    },
  },
}));

const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('ProjectsPage', () => {
  beforeEach(() => {
    mockedAxios.get.mockReset();
  });

  test('загружает и отображает проекты', async () => {
    mockedAxios.get
      .mockResolvedValueOnce({
        data: {
          projects: [
            { id: '1', key: 'PRJ1', name: 'Project One', self: 'url1' },
          ],
          pageInfo: { pageCount: 1 },
        },
      })
      .mockResolvedValueOnce({
        data: [
          { id: '10', key: 'PRJ1', name: 'Project One', self: 'url1-db' },
        ],
      });

    render(<ProjectsPage />);

    await waitFor(() => {
      expect(screen.getByText('Project One')).toBeInTheDocument();
    });

    expect(mockedAxios.get).toHaveBeenCalledWith(
      config.api.endpoints.connectorProjects,
      expect.any(Object),
    );
    expect(mockedAxios.get).toHaveBeenCalledWith(
      config.api.endpoints.projects,
    );
  });

  test('показывает пустое состояние при отсутствии проектов', async () => {
    mockedAxios.get
      .mockResolvedValueOnce({
        data: {
          projects: [],
          pageInfo: { pageCount: 1 },
        },
      })
      .mockResolvedValueOnce({ data: [] });

    render(<ProjectsPage />);

    await waitFor(() => {
      expect(screen.getByText('Нет доступных проектов')).toBeInTheDocument();
    });
  });

  test('фильтрация по поиску и очистка поля', async () => {
    mockedAxios.get
      .mockResolvedValueOnce({
        data: {
          projects: [
            { id: '1', key: 'PRJ1', name: 'Alpha', self: 'url1' },
            { id: '2', key: 'PRJ2', name: 'Beta', self: 'url2' },
          ],
          pageInfo: { pageCount: 1 },
        },
      })
      .mockResolvedValueOnce({ data: [] })
      // повторный запрос после изменения строки поиска
      .mockResolvedValueOnce({
        data: {
          projects: [
            { id: '2', key: 'PRJ2', name: 'Beta', self: 'url2' },
          ],
          pageInfo: { pageCount: 1 },
        },
      })
      .mockResolvedValueOnce({ data: [] });

    render(<ProjectsPage />);

    await waitFor(() => {
      expect(screen.getByText('Alpha')).toBeInTheDocument();
    });

    const input = screen.getByPlaceholderText('Поиск проектов...');
    fireEvent.change(input, { target: { value: 'Beta' } });

    await waitFor(() => {
      expect(screen.getByText('Beta')).toBeInTheDocument();
    });

    const clearButton = screen.getAllByRole('button').find((btn) =>
      btn.className.includes('clear-search-btn'),
    );
    if (clearButton) {
      fireEvent.click(clearButton);
    }

    await waitFor(() => {
      expect((input as HTMLInputElement).value).toBe('');
    });
  });
});

