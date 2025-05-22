import React, { useEffect, useState } from 'react';
import { Bar } from 'react-chartjs-2';
import axios from 'axios';
import './Chart.scss';
import { config } from '../../config/config';

interface TimeOpenChartProps {
  projectKey: string;
}

export interface TimeOpenData {
  range: string;
  count: number;
}

const TimeOpenChart: React.FC<TimeOpenChartProps> = ({ projectKey }) => {
  const [data, setData] = useState<TimeOpenData[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    setLoading(true);
    setError(false);
    
    axios.get(config.api.endpoints.timeOpenAnalytics, { params: { key: projectKey } })
      .then(res => {
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
    labels: data.map(item => item.range),
    datasets: [{
      label: 'Количество задач',
      data: data.map(item => item.count),
      backgroundColor: '#3b82f6',
      borderColor: '#2563eb',
      borderWidth: 1
    }]
  };

  const options = {
    responsive: true,
    plugins: {
      title: {
        display: true,
        text: 'Время задач в открытом состоянии',
        font: {
          size: 16
        }
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        title: {
          display: true,
          text: 'Количество задач'
        },
        ticks: {
          stepSize: 1
        }
      },
      x: {
        title: {
          display: true,
          text: 'Дней в открытом состоянии'
        }
      }
    }
  };

  return (
    <div className="chart-wrapper" style={{ minHeight: '300px' }}>
      <Bar data={chartData} options={options} />
    </div>
  );
};

export default TimeOpenChart;