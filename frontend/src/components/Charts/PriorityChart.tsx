import React, { useEffect, useState } from 'react';
import { Doughnut } from 'react-chartjs-2';
import axios from 'axios';
import './Chart.scss';
import { config } from '../../config/config';

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
  const [error, setError] = useState(false);

  useEffect(() => {
    setLoading(true);
    setError(false);
    
    axios.get(config.api.endpoints.priorityAnalytics, { params: { key: projectKey } })
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
  if (data.length === 0) return <div className="chart-no-data">Нет данных о приоритетах</div>;

  const backgroundColors = [
    '#ef4444', '#f97316', '#f59e0b', '#84cc16'
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
        font: { size: 16 }
      },
      legend: { position: 'right' as const }
    },
    cutout: '70%'
  };

  return (
    <div className="chart-wrapper-priority" style={{ minHeight: '300px' }}>
      <Doughnut data={chartData} options={options} />
    </div>
  );
};

export default PriorityChart;