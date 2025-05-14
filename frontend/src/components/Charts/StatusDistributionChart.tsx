import React, { useEffect, useState } from 'react';
import { Pie } from 'react-chartjs-2';
import axios from 'axios';
import './Chart.scss';

interface StatusDistributionChartProps {
  projectKey: string;
}

export interface StatusDistributionData {
  status: string;
  count: number;
}

const StatusDistributionChart: React.FC<StatusDistributionChartProps> = ({ projectKey }) => {
  const [data, setData] = useState<StatusDistributionData[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    axios.get('/api/v1/analytics/status-distribution', { params: { project: projectKey } })
      .then(res => {
        setData(res.data.data);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, [projectKey]);

  if (loading) return <div className="chart-loading">Загрузка данных...</div>;

  const backgroundColors = [
    '#ef4444', // Open - red
    '#f59e0b', // In Progress - amber
    '#10b981', // Resolved - emerald
    '#3b82f6', // Closed - blue
    '#8b5cf6'  // Reopened - violet
  ];

  const chartData = {
    labels: data.map(item => item.status),
    datasets: [{
      data: data.map(item => item.count),
      backgroundColor: backgroundColors,
      borderColor: '#fff',
      borderWidth: 2
    }]
  };

  const options = {
    responsive: true,
    plugins: {
      title: {
        display: true,
        text: 'Распределение задач по статусам',
        font: {
          size: 16
        }
      },
      legend: {
        position: 'right' as const,
      }
    }
  };

  return (
    <div className="chart-wrapper">
      <Pie data={chartData} options={options} />
    </div>
  );
};

export default StatusDistributionChart;