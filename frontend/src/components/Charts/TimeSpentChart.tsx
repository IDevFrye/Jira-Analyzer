import React, { useEffect, useState } from 'react';
import { Bar } from 'react-chartjs-2';
import axios from 'axios';
import './Chart.scss';
import { config } from '../../config/config';

interface TimeSpentChartProps {
  projectKey: string;
}

export interface TimeSpentData {
  user: string;
  time: number;
}

const TimeSpentChart: React.FC<TimeSpentChartProps> = ({ projectKey }) => {
  const [data, setData] = useState<TimeSpentData[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    setLoading(true);
    setError(false);
    
    axios.get(config.api.endpoints.timeSpentAnalytics, { params: { key: projectKey } })
      .then(res => {
        const responseData = Array.isArray(res.data) ? res.data : 
                          res.data?.data ? res.data.data : [];
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
  if (data.length === 0) return <div className="chart-no-data">Нет данных о времени</div>;

  const chartData = {
    labels: data.map(item => item.user),
    datasets: [{
      label: 'Затраченное время (часы)',
      data: data.map(item => item.time),
      backgroundColor: '#10b981',
      borderColor: '#059669',
      borderWidth: 1
    }]
  };

  const options = {
    responsive: true,
    indexAxis: 'y' as const,
    plugins: {
      title: {
        display: true,
        text: 'Затраченное время по пользователям',
        font: { size: 16 }
      },
    },
    scales: {
      x: {
        beginAtZero: true,
        title: { display: true, text: 'Часы' }
      },
      y: {
        title: { display: true, text: 'Пользователи' }
      }
    }
  };

  return (
    <div className="chart-wrapper" style={{ minHeight: '300px' }}>
      <Bar data={chartData} options={options} />
    </div>
  );
};

export default TimeSpentChart;