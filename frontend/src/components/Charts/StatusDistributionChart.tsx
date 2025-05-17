import React, { useEffect, useState } from 'react';
import { Pie } from 'react-chartjs-2';
import axios from 'axios';
import './Chart.scss';
import { config } from '../../config/config';

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
  const [error, setError] = useState(false);

  useEffect(() => {
    setLoading(true);
    setError(false);
    
    axios.get(config.api.endpoints.statusDistribution, { params: { key: projectKey } })
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
  if (data.length === 0) return <div className="chart-no-data">Нет данных о статусах</div>;

  const backgroundColors = [
    '#ef4444', '#f59e0b', '#10b981', '#3b82f6', '#8b5cf6'
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
        font: { size: 16 }
      },
      legend: { position: 'right' as const }
    }
  };

  return (
    <div className="chart-wrapper-status" style={{ minHeight: '300px' }}>
      <Pie data={chartData} options={options} />
    </div>
  );
};

export default StatusDistributionChart;