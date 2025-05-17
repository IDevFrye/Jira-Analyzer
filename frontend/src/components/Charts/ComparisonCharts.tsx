import React from 'react';
import { Bar, Pie, Doughnut } from 'react-chartjs-2';
import './ComparisonCharts.scss';

interface ComparisonChartProps {
  data: {
    labels: string[];
    datasets: {
      label: string;
      data: number[];
      backgroundColor: string | string[];
      borderColor?: string | string[];
      borderWidth?: number;
    }[];
  };
  title: string;
  type: 'bar' | 'pie' | 'doughnut';
  indexAxis?: 'x' | 'y';
}

const ComparisonChart: React.FC<ComparisonChartProps> = ({ 
  data, 
  title, 
  type,
  indexAxis
}) => {
  const options = {
    responsive: true,
    plugins: {
      title: {
        display: true,
        text: title,
        font: {
          size: 16
        }
      },
      legend: {
        position: 'top' as const,
      }
    },
    ...(indexAxis && { indexAxis })
  };

  const renderChart = () => {
    switch (type) {
      case 'bar':
        return <Bar data={data} options={options} />;
      case 'pie':
        return <Pie data={data} options={options} />;
      case 'doughnut':
        return <Doughnut data={data} options={options} />;
      default:
        return <Bar data={data} options={options} />;
    }
  };

  return (
    <div className="comparison-chart">
      {renderChart()}
    </div>
  );
};

export const TimeOpenComparisonChart: React.FC<{ projectKeys: string[] }> = ({ projectKeys }) => {
  const ranges = ['0-1', '1-2', '2-3', '3-5', '5-7', '7-10', '10-14', '14-21', '21-30', '30+'];
  const colors = ['#3b82f6', '#ef4444', '#10b981', '#f59e0b', '#8b5cf6'];
  
  const data = {
    labels: ranges,
    datasets: projectKeys.map((key, i) => ({
      label: key,
      data: ranges.map(() => Math.floor(Math.random() * 50) + 5),
      backgroundColor: colors[i % colors.length],
      borderColor: colors[i % colors.length],
      borderWidth: 1
    }))
  };

  return (
    <ComparisonChart 
      data={data}
      title="Сравнение времени задач в открытом состоянии"
      type="bar"
    />
  );
};

export const StatusDistributionComparisonChart: React.FC<{ projectKeys: string[] }> = ({ projectKeys }) => {
  const statuses = ['Open', 'In Progress', 'Resolved', 'Closed', 'Reopened'];
  const colors = ['#ef4444', '#f59e0b', '#10b981', '#3b82f6', '#8b5cf6'];
  
  const data = {
    labels: statuses,
    datasets: projectKeys.map((key, i) => ({
      label: key,
      data: statuses.map(() => Math.floor(Math.random() * 100) + 10),
      backgroundColor: colors,
      borderColor: '#fff',
      borderWidth: 2
    }))
  };

  return (
    <ComparisonChart 
      data={data}
      title="Сравнение распределения задач по статусам"
      type="doughnut"
    />
  );
};

export const TimeSpentComparisonChart: React.FC<{ projectKeys: string[] }> = ({ projectKeys }) => {
  const users = ['John Doe', 'Jane Smith', 'Mike Johnson', 'Sarah Williams', 'David Brown'];
  const colors = ['#3b82f6', '#ef4444', '#10b981', '#f59e0b', '#8b5cf6'];
  
  const data = {
    labels: users,
    datasets: projectKeys.map((key, i) => ({
      label: key,
      data: users.map(() => Math.floor(Math.random() * 80) + 5),
      backgroundColor: colors[i % colors.length],
      borderColor: colors[i % colors.length],
      borderWidth: 1
    }))
  };

  return (
    <ComparisonChart 
      data={data}
      title="Сравнение затраченного времени по пользователям"
      type="bar"
      indexAxis="y"
    />
  );
};

export const PriorityComparisonChart: React.FC<{ projectKeys: string[] }> = ({ projectKeys }) => {
  const priorities = ['Critical', 'High', 'Medium', 'Low'];
  const colors = ['#ef4444', '#f97316', '#f59e0b', '#84cc16'];
  
  const data = {
    labels: priorities,
    datasets: projectKeys.map((key, i) => ({
      label: key,
      data: priorities.map(() => Math.floor(Math.random() * 50) + 5),
      backgroundColor: colors,
      borderColor: '#fff',
      borderWidth: 2
    }))
  };

  return (
    <ComparisonChart 
      data={data}
      title="Сравнение распределения задач по приоритетам"
      type="pie"
    />
  );
};