import React, { useEffect, useState } from 'react';
import { Line } from 'react-chartjs-2';
import axios from 'axios';
import './Chart.scss';
import { config } from '../../config/config';

interface ThroughputChartProps {
  projectKey: string;
}

export interface ThroughputPoint {
  date: string;
  count: number;
}

const ThroughputChart: React.FC<ThroughputChartProps> = ({ projectKey }) => {
  const [data, setData] = useState<ThroughputPoint[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    setLoading(true);
    setError(false);

    axios
      .get(config.api.endpoints.throughputAnalytics, { params: { key: projectKey } })
      .then((res) => {
        const responseData = Array.isArray(res.data) ? res.data : res.data?.data || [];
        setData(responseData);
        setLoading(false);
      })
      .catch(() => {
        setLoading(false);
        setError(true);
      });
  }, [projectKey]);

  if (loading) return <div className="chart-loading">Загрузка данных...</div>;
  if (error) return <div className="chart-error">Ошибка загрузки данных</div>;
  if (data.length === 0) return <div className="chart-no-data">Нет данных для отображения</div>;

  const chartData = {
    labels: data.map((point) => point.date),
    datasets: [
      {
        label: 'Создано задач за день',
        data: data.map((point) => point.count),
        borderColor: '#10b981',
        backgroundColor: 'rgba(16, 185, 129, 0.2)',
        tension: 0.2,
      },
    ],
  };

  const options = {
    responsive: true,
    plugins: {
      title: {
        display: true,
        text: 'Пропускная способность (создано задач по дням)',
        font: { size: 16 },
      },
      legend: {
        position: 'top' as const,
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        title: {
          display: true,
          text: 'Количество задач',
        },
      },
      x: {
        title: {
          display: true,
          text: 'Дата создания',
        },
      },
    },
  };

  return (
    <div className="chart-wrapper" style={{ minHeight: '300px' }}>
      <Line data={chartData} options={options} />
    </div>
  );
};

export default ThroughputChart;

