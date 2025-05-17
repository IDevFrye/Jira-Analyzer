import React, { useEffect, useState } from 'react';
import { Bar } from 'react-chartjs-2';
import axios from 'axios';
import './Chart.scss';

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

  useEffect(() => {
    axios.get('/api/v1/analytics/time-spent', { params: { project: projectKey } })
      .then(res => {
        setData(res.data.data);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, [projectKey]);

  if (loading) return <div className="chart-loading">Загрузка данных...</div>;

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
        font: {
          size: 16
        }
      },
    },
    scales: {
      x: {
        beginAtZero: true,
        title: {
          display: true,
          text: 'Часы'
        }
      },
      y: {
        title: {
          display: true,
          text: 'Пользователи'
        }
      }
    }
  };

  return (
    <div className="chart-wrapper">
      <Bar 
        data={chartData} 
        options={{
          ...options,
          indexAxis: 'y' as const,
        }} 
      />
    </div>
  );
};

export default TimeSpentChart;