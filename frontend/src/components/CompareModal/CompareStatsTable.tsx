import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Project } from '../../types/models';
import './CompareStatsTable.scss';

interface ProjectStats extends Project {
  allIssuesCount?: number;
  openIssuesCount?: number;
  closeIssuesCount?: number;
  resolvedIssuesCount?: number;
  reopenedIssuesCount?: number;
  progressIssuesCount?: number;
  averageTime?: number;
  averageIssuesCount?: number;
}

interface CompareStatsTableProps {
  projects: Project[];
}

const CompareStatsTable: React.FC<CompareStatsTableProps> = ({ projects }) => {
  const [statsData, setStatsData] = useState<ProjectStats[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const responses = await Promise.all(
          projects.map(project => 
            axios.get(`/api/v1/projects/${project.Id}`)
          )
        );
        setStatsData(responses.map(res => res.data));
      } catch (error) {
        console.error('Error fetching stats:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, [projects]);

  const stats = [
    { label: 'Общее количество задач', key: 'allIssuesCount' },
    { label: 'Количество открытых задач', key: 'openIssuesCount' },
    { label: 'Количество закрытых задач', key: 'closeIssuesCount' },
    { label: 'Количество переоткрытых задач', key: 'reopenedIssuesCount' },
    { label: 'Количество разрешенных задач', key: 'resolvedIssuesCount' },
    { label: 'Количество задач "In progress"', key: 'progressIssuesCount' },
    { 
      label: 'Среднее время выполнения (часы)', 
      key: 'averageTime', 
      format: (val?: number) => val ? val.toFixed(1) : 'N/A' 
    },
    { 
      label: 'Среднее количество задач в день', 
      key: 'averageIssuesCount',
      format: (val?: number) => val ? val.toString() : 'N/A'
    }
  ];

  if (loading) return <div className="loading">Загрузка данных...</div>;

  return (
    <div className="stats-table-container">
      <table className="stats-table">
        <thead>
          <tr>
            <th>Сухая статистика</th>
            {projects.map(project => (
              <th key={project.Id}>{project.Name}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {stats.map(stat => (
            <tr key={stat.key}>
              <td>{stat.label}</td>
              {statsData.map((projectData, index) => (
                <td key={`${projects[index].Id}-${stat.key}`}>
                  {stat.format 
                    ? stat.format(projectData[stat.key as keyof ProjectStats] as number)
                    : projectData[stat.key as keyof ProjectStats] ?? 'N/A'}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default CompareStatsTable;