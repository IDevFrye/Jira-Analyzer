import React, { useEffect, useState } from 'react';
import { Bar } from 'react-chartjs-2';
import axios from 'axios';
import './Chart.scss';

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

  useEffect(() => {
    axios.get('/api/v1/analytics/time-open', { params: { project: projectKey } })
      .then(res => {
        setData(res.data.data);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, [projectKey]);

  if (loading) return <div className="chart-loading">Загрузка данных...</div>;

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
    <div className="chart-wrapper">
      <Bar data={chartData} options={options} />
    </div>
  );
};

export default TimeOpenChart;