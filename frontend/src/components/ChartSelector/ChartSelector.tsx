import React, { useState } from 'react';
import TimeOpenChart from '../Charts/TimeOpenChart';
import StatusDistributionChart from '../Charts/StatusDistributionChart';
import TimeSpentChart from '../Charts/TimeSpentChart';
import PriorityChart from '../Charts/PriorityChart';
import './ChartSelector.scss';

interface ChartSelectorProps {
  projectKey: string;
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


const ChartSelector: React.FC<ChartSelectorProps> = ({ projectKey }) => {
  const [selectedChart, setSelectedChart] = useState<ChartType>('timeOpen');

  const charts: { id: ChartType; name: string }[] = [
    { id: 'timeOpen', name: 'Время в открытом состоянии' },
    { id: 'statusDistribution', name: 'Распределение по статусам' },
    { id: 'timeSpent', name: 'Затраченное время' },
    { id: 'priority', name: 'По приоритетам' },
  ];

  const renderChart = () => {
    switch (selectedChart) {
      case 'timeOpen':
        return <TimeOpenChart projectKey={projectKey} />;
      case 'statusDistribution':
        return <StatusDistributionChart projectKey={projectKey} />;
      case 'timeSpent':
        return <TimeSpentChart projectKey={projectKey} />;
      case 'priority':
        return <PriorityChart projectKey={projectKey} />;
      default:
        return <TimeOpenChart projectKey={projectKey} />;
    }
  };

  return (
    <div className="chart-selector">
      <h3 className="chart-selector-title">Аналитика проекта {projectKey}</h3>
      
      <div className="chart-tabs">
        {charts.map((chart) => (
          <button
            key={chart.id}
            className={`chart-tab ${selectedChart === chart.id ? 'active' : ''}`}
            onClick={() => setSelectedChart(chart.id)}
          >
            {chart.name}
          </button>
        ))}
      </div>
      
      <div className="chart-container">
        {renderChart()}
      </div>
    </div>
  );
};

export default ChartSelector;