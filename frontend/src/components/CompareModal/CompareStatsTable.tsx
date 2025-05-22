import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Project, ProjectStats } from '../../types/models';
import './CompareStatsTable.scss';
import { config } from '../../config/config';

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
            axios.get(config.api.endpoints.projectStats(project.Id))
          )
        );
        
        const formattedData = responses.map((res, index) => ({
          Id: projects[index].Id,
          Key: projects[index].Key,
          Name: projects[index].Name,
          self: projects[index].self,
          allIssuesCount: res.data.total_issues,
          openIssuesCount: res.data.open_issues,
          closeIssuesCount: res.data.closed_issues,
          reopenedIssuesCount: res.data.reopened_issues,
          resolvedIssuesCount: res.data.resolved_issues,
          progressIssuesCount: res.data.in_progress_issues,
          averageTime: Number(res.data.avg_resolution_time_h.toFixed(2)),
          averageIssuesCount: Number(res.data.avg_created_per_day_7d.toFixed(2)) 
        }));
        
        setStatsData(formattedData);
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
          {statsData.length > 0 ? (
            stats.map(stat => (
              <tr key={stat.key}>
                <td>{stat.label}</td>
                {statsData.map((projectData, index) => (
                  <td key={`${projects[index].Id}-${stat.key}`}>
                    {JSON.stringify(projectData[stat.key as keyof ProjectStats])}
                  </td>
                ))}
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan={projects.length + 1}>Нет данных для отображения</td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
};

export default CompareStatsTable;