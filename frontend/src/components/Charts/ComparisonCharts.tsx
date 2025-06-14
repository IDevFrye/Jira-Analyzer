import React, { useEffect, useState } from 'react';
import { Bar, Pie, Doughnut } from 'react-chartjs-2';
import './ComparisonCharts.scss';
import { config } from '../../config/config';
import axios from 'axios';
import { darken, lighten } from 'polished';

const BASE_COLORS = ['#ef4444', '#f97316', '#f59e0b', '#84cc16', '#06b6d4', '#8b5cf6', '#ec4899'];

function getBaseColorMap(keys: string[]): Record<string, string> {
  const map: Record<string, string> = {};
  keys.forEach((key, index) => {
    map[key] = BASE_COLORS[index % BASE_COLORS.length];
  });
  return map;
}

function generateShadedColors(baseMap: Record<string, string>, projectIndex: number, totalProjects: number): string[] {
  return Object.keys(baseMap).map(label => {
    const base = baseMap[label];
    if (projectIndex === 0) return base;

    const factor =  0.35 * (projectIndex / Math.max(1, totalProjects - 1));
    try {
      const shaded = lighten(factor, base);
      return shaded || base;
    } catch {
      return base;
    }
  });
}


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


interface DataItem {
  range?: string;
  status?: string;
  user?: string;
  time?: number;
  priority?: string;
  count: number;
}

interface ProjectData {
  project: string;
  data: DataItem[];
}

interface ApiResponse {
  projects: ProjectData[];
}

interface ChartDataset {
  label: string;
  data: number[];
  backgroundColor: string | string[];
  borderColor?: string | string[];
  borderWidth?: number;
}

interface ChartData {
  labels: string[];
  datasets: ChartDataset[];
}

interface ApiResponse {
  projects: ProjectData[];
}

interface TimeOpenData {
  [projectKey: string]: Array<{
    range: string;
    count: number;
  }>;
}

interface StatusDistributionData {
  [projectKey: string]: {
    [status: string]: number;
  };
}

interface TimeSpentData {
  [projectKey: string]: {
    authors: Array<{
      author: string;
      total_time_spent: number;
    }>;
  };
}

interface PriorityData {
  [projectKey: string]: {
    [priority: string]: number;
  };
}

// constants/chart.ts
export const CHART_COLORS = [
  '#3b82f6', // blue
  '#ef4444', // red
  '#10b981', // green
  '#f59e0b', // yellow
  '#8b5cf6', // purple
  '#f97316', // orange
  '#06b6d4', // cyan
  '#ec4899'  // pink
];

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
    (type === 'doughnut' || type === 'pie')  ?
    (<div className="comparison-chart-pie">
      {renderChart()}
    </div>) : (
      <div className="comparison-chart">
      {renderChart()}
    </div>
    )
  );
};

export const TimeOpenComparisonChart: React.FC<{ projectKeys: string[] }> = ({ projectKeys }) => {
  const [data, setData] = useState<ChartData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get<TimeOpenData>(config.api.endpoints.compareTimeOpen, {
          params: { key: projectKeys.join(',') }
        });
        
        const projectData = response.data;
        if (!projectData || Object.keys(projectData).length === 0) {
          throw new Error('Нет данных о времени открытия задач');
        }

        // Собираем все уникальные диапазоны из всех проектов
        const allRanges = new Set<string>();
        projectKeys.forEach(key => {
          if (projectData[key]) {
            projectData[key]?.forEach(item => {
              allRanges.add(item.range);
            });
          }
        });

        const sortedRanges = Array.from(allRanges).sort((a, b) => {
          // Сортируем диапазоны по логическому порядку
          // Это примерная реализация, вам может потребоваться более точная сортировка
          const getRangeValue = (range: string) => {
            if (range === '30+') return Infinity;
            const numbers = range.split('-').map(Number);
            return numbers[0];
          };
          return getRangeValue(a) - getRangeValue(b);
        });

        const chartData: ChartData = {
          labels: sortedRanges,
          datasets: projectKeys.map((key, i) => {
            const project = projectData[key];
            const countsMap = new Map<string, number>();
            
            // Инициализируем все диапазоны нулями
            sortedRanges.forEach(range => countsMap.set(range, 0));
            
            // Заполняем реальными значениями
            if (project) {
              project.forEach(item => {
                countsMap.set(item.range, item.count);
              });
            }
            
            return {
              label: key,
              data: sortedRanges.map(range => countsMap.get(range) || 0),
              backgroundColor: CHART_COLORS[i % CHART_COLORS.length],
              borderColor: CHART_COLORS[i % CHART_COLORS.length],
              borderWidth: 1
            };
          })
        };
        
        setData(chartData);
      } catch (err: unknown) {
        setError(err instanceof Error ? err.message : 'Неизвестная ошибка');
        console.error('Ошибка загрузки времени открытия:', err);
      } finally {
        setLoading(false);
      }
    };

    if (projectKeys.length > 0) {
      fetchData();
    }
  }, [projectKeys]);

  if (loading) return <div className="chart-loading">Загрузка данных...</div>;
  if (error) return <div className="chart-error">{error}</div>;
  if (!data) return <div className="chart-no-data">Нет данных для отображения</div>;

  return (
    <ComparisonChart 
      data={data}
      title="Сравнение времени открытия задач"
      type="bar"
    />
  );
};

export const StatusDistributionComparisonChart: React.FC<{ projectKeys: string[] }> = ({ projectKeys }) => {
  const [data, setData] = useState<ChartData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get<StatusDistributionData>(config.api.endpoints.compareStatusDistribution, {
          params: { key: projectKeys.join(',') }
        });

        const projectData = response.data;
        if (!projectData || Object.keys(projectData).length === 0) {
          throw new Error('Нет данных о статусах задач');
        }

        const allStatuses = new Set<string>();
        projectKeys.forEach(key => {
          const statuses = projectData[key] ? Object.keys(projectData[key]) : [];
          statuses.forEach(status => allStatuses.add(status));
        });
        const statuses = Array.from(allStatuses);
        const baseMap = getBaseColorMap(statuses);

        const chartData: ChartData = {
          labels: statuses,
          datasets: projectKeys.map((key, i) => ({
            label: key,
            data: statuses.map(status => projectData[key]?.[status] || 0),
            backgroundColor: generateShadedColors(baseMap, i, projectKeys.length),
            borderColor: '#fff',
            borderWidth: 2
          }))
        };

        setData(chartData);
      } catch (err: unknown) {
        setError(err instanceof Error ? err.message : 'Неизвестная ошибка');
        console.error('Ошибка загрузки статусов:', err);
      } finally {
        setLoading(false);
      }
    };

    if (projectKeys.length > 0) {
      fetchData();
    }
  }, [projectKeys]);

  if (loading) return <div className="chart-loading">Загрузка данных...</div>;
  if (error) return <div className="chart-error">{error}</div>;
  if (!data) return <div className="chart-no-data">Нет данных для отображения</div>;

  return (
    <ComparisonChart 
      data={data}
      title="Сравнение распределения задач по статусам"
      type="doughnut"
    />
  );
};

export const TimeSpentComparisonChart: React.FC<{ projectKeys: string[] }> = ({ projectKeys }) => {
  const [data, setData] = useState<ChartData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get<TimeSpentData>(config.api.endpoints.compareTimeSpent, {
          params: { key: projectKeys.join(',') }
        });
        
        const projectData = response.data;
        if (!projectData || Object.keys(projectData).length === 0) {
          throw new Error('Нет данных о затраченном времени');
        }

        const allAuthors = new Map<string, number>();
        projectKeys.forEach(key => {
          const authors = projectData[key]?.authors || [];
          authors.forEach(author => {
            const current = allAuthors.get(author.author) || 0;
            allAuthors.set(author.author, current + author.total_time_spent);
          });
        });
        
        const topAuthors = Array.from(allAuthors.entries())
          .sort((a, b) => b[1] - a[1])
          .slice(0, 10)
          .map(item => item[0]);
        
        const chartData: ChartData = {
          labels: topAuthors,
          datasets: projectKeys.map((key, i) => ({
            label: key,
            data: topAuthors.map(author => {
              const found = projectData[key]?.authors.find(a => a.author === author);
              return found ? found.total_time_spent : 0;
            }),
            backgroundColor: CHART_COLORS[i % CHART_COLORS.length],
            borderColor: CHART_COLORS[i % CHART_COLORS.length],
            borderWidth: 1
          }))
        };
        
        setData(chartData);
      } catch (err: unknown) {
        setError(err instanceof Error ? err.message : 'Неизвестная ошибка');
        console.error('Ошибка загрузки времени:', err);
      } finally {
        setLoading(false);
      }
    };

    if (projectKeys.length > 0) {
      fetchData();
    }
  }, [projectKeys]);

  if (loading) return <div className="chart-loading">Загрузка данных...</div>;
  if (error) return <div className="chart-error">{error}</div>;
  if (!data) return <div className="chart-no-data">Нет данных для отображения</div>;

  return (
    <ComparisonChart 
      data={data}
      title="Сравнение затраченного времени по авторам"
      type="bar"
      indexAxis="y"
    />
  );
};

export const PriorityComparisonChart: React.FC<{ projectKeys: string[] }> = ({ projectKeys }) => {
  const [data, setData] = useState<ChartData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get<PriorityData>(config.api.endpoints.comparePriority, {
          params: { key: projectKeys.join(',') }
        });

        const projectData = response.data;
        if (!projectData || Object.keys(projectData).length === 0) {
          throw new Error('Нет данных о приоритетах');
        }

        const allPriorities = new Set<string>();
        projectKeys.forEach(key => {
          const priorities = projectData[key] ? Object.keys(projectData[key]) : [];
          priorities.forEach(priority => allPriorities.add(priority));
        });
        const priorities = Array.from(allPriorities);
        const baseMap = getBaseColorMap(priorities);

        const chartData: ChartData = {
          labels: priorities,
          datasets: projectKeys.map((key, i) => ({
            label: key,
            data: priorities.map(priority => projectData[key]?.[priority] || 0),
            backgroundColor: generateShadedColors(baseMap, i, projectKeys.length),
            borderColor: '#fff',
            borderWidth: 2
          }))
        };

        setData(chartData);
      } catch (err: unknown) {
        setError(err instanceof Error ? err.message : 'Неизвестная ошибка');
        console.error('Ошибка загрузки приоритетов:', err);
      } finally {
        setLoading(false);
      }
    };

    if (projectKeys.length > 0) {
      fetchData();
    }
  }, [projectKeys]);

  if (loading) return <div className="chart-loading">Загрузка данных...</div>;
  if (error) return <div className="chart-error">{error}</div>;
  if (!data) return <div className="chart-no-data">Нет данных для отображения</div>;

  return (
    <ComparisonChart 
      data={data}
      title="Сравнение распределения задач по приоритетам"
      type="pie"
    />
  );
};
