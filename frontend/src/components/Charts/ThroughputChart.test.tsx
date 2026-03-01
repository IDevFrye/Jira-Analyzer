import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import ThroughputChart from './ThroughputChart';
import axios from 'axios';

jest.mock('axios');
jest.mock('react-chartjs-2', () => ({
  Line: () => <div data-testid="line-chart" />,
}));

const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('ThroughputChart', () => {
  test('отображает данные графика при успешной загрузке', async () => {
    mockedAxios.get.mockResolvedValueOnce({
      data: [
        { date: '2025-01-01', count: 5 },
        { date: '2025-01-02', count: 3 },
      ],
    });

    render(<ThroughputChart projectKey="PRJ1" />);

    expect(screen.getByText('Загрузка данных...')).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByTestId('line-chart')).toBeInTheDocument();
    });
  });

  test('отображает сообщение при отсутствии данных', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: [] });

    render(<ThroughputChart projectKey="PRJ1" />);

    await waitFor(() => {
      expect(screen.getByText('Нет данных для отображения')).toBeInTheDocument();
    });
  });

  test('отображает ошибку при неуспешной загрузке', async () => {
    mockedAxios.get.mockRejectedValueOnce(new Error('network error'));

    render(<ThroughputChart projectKey="PRJ1" />);

    await waitFor(() => {
      expect(screen.getByText('Ошибка загрузки данных')).toBeInTheDocument();
    });
  });
});

