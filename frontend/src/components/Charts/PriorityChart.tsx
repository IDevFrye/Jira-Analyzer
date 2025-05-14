import React, { useEffect, useState } from 'react';
import { Doughnut } from 'react-chartjs-2';
import axios from 'axios';
import './Chart.scss';

interface PriorityChartProps {
  projectKey: string;
}

export interface PriorityData {
  priority: string;
  count: number;
}

const PriorityChart: React.FC<PriorityChartProps> = ({ projectKey }) => {
  const [data, setData] = useState<PriorityData[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    axios.get('/api/v1/analytics/priority', { params: { project: projectKey } })
      .then(res => {
        setData(res.data.data);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, [projectKey]);

  if (loading) return <div className="chart-loading">Загрузка данных...</div>;

  const backgroundColors = [
    '#ef4444', // Critical - red
    '#f97316', // High - orange
    '#f59e0b', // Medium - amber
    '#84cc16'  // Low - lime
  ];

  const chartData = {
    labels: data.map(item => item.priority),
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
        text: 'Распределение задач по приоритетам',
        font: {
          size: 16
        }
      },
      legend: {
        position: 'right' as const,
      }
    },
    cutoutPercentage: 70
  };

  return (
    <div className="chart-wrapper">
      <Doughnut data={chartData} options={options} />
    </div>
  );
};

export default PriorityChart;