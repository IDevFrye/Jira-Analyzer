import React, { useEffect, useState } from 'react';
import { 
  TimeOpenComparisonChart,
  StatusDistributionComparisonChart,
  TimeSpentComparisonChart,
  PriorityComparisonChart
} from '../Charts/ComparisonCharts';
import './CompareCharts.scss';

interface CompareChartsProps {
  projects: { Key: string }[];
}

export type ChartType = 
  | 'timeOpen' 
  | 'statusDistribution' 
  | 'timeSpent' 
  | 'priority';

export interface ChartData {
  labels: string[];
  datasets: {
    label: string;
    data: number[];
    backgroundColor: string[];
    borderColor?: string[];
    borderWidth?: number;
  }[];
}

const CompareCharts: React.FC<CompareChartsProps> = ({ projects }) => {
  const [chartType, setChartType] = useState<'timeOpen' | 'statusDistribution' | 'timeSpent' | 'priority'>('timeOpen');
  const [loading, setLoading] = useState(false);

  const projectKeys = projects.map(p => p.Key);

  if (loading) return <div className="comparison-loading">Загрузка данных...</div>;

  return (
    <div className="compare-charts">
      <div className="chart-tabs">
        <button
          className={`chart-tab ${chartType === 'timeOpen' ? 'active' : ''}`}
          onClick={() => setChartType('timeOpen')}
        >
          Время в открытом состоянии
        </button>
        <button
          className={`chart-tab ${chartType === 'statusDistribution' ? 'active' : ''}`}
          onClick={() => setChartType('statusDistribution')}
        >
          Распределение по статусам
        </button>
        <button
          className={`chart-tab ${chartType === 'timeSpent' ? 'active' : ''}`}
          onClick={() => setChartType('timeSpent')}
        >
          Затраченное время
        </button>
        <button
          className={`chart-tab ${chartType === 'priority' ? 'active' : ''}`}
          onClick={() => setChartType('priority')}
        >
          По приоритетам
        </button>
      </div>
      
      <div className="chart-container">
        {chartType === 'timeOpen' && <TimeOpenComparisonChart projectKeys={projectKeys} />}
        {chartType === 'statusDistribution' && <StatusDistributionComparisonChart projectKeys={projectKeys} />}
        {chartType === 'timeSpent' && <TimeSpentComparisonChart projectKeys={projectKeys} />}
        {chartType === 'priority' && <PriorityComparisonChart projectKeys={projectKeys} />}
      </div>
    </div>
  );
};

export default CompareCharts;