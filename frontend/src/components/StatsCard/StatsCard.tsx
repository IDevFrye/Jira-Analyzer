import React, { useEffect, useState } from 'react';
import axios from 'axios';
import {
  FaTasks,
  FaFolderOpen,
  FaCheckCircle,
  FaLockOpen,
  FaRedo,
  FaClock,
  FaChartBar
} from 'react-icons/fa';
import { FaCalendarAlt } from 'react-icons/fa';
import { GiSandsOfTime } from 'react-icons/gi';
import Modal from '../Modal/Modal';
import ChartSelector from '../ChartSelector/ChartSelector';
import { config } from '../../config/config';
import './StatsCard.scss';

interface ProjectStats {
  id: string;
  key: string;
  name: string;
  total_issues: number;
  open_issues: number;
  closed_issues: number;
  reopened_issues: number;
  resolved_issues: number;
  in_progress_issues: number;
  avg_resolution_time_h: number;
  avg_created_per_day_7d: number;
}

interface StatsCardProps {
  projectId: string;
  projectName: string;
  projectKey: string;
}

const StatsCard: React.FC<StatsCardProps> = ({ projectId, projectName, projectKey }) => {
  const [stats, setStats] = useState<ProjectStats | null>(null);
  const [showAnalytics, setShowAnalytics] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        setLoading(true);
        setError(null);
        const response = await axios.get(config.api.endpoints.projectStats(Number(projectId)));
        setStats(response.data);
      } catch (err) {
        console.error('Error fetching project stats:', err);
        setError('Не удалось загрузить данные проекта');
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, [projectId]);

  const formatTime = (hours: number | undefined) => {
    if (hours === undefined) return 'N/A';
    const absHours = Math.abs(hours);
    return absHours.toFixed(1);
  };

  if (loading) {
    return (
      <div className="stats-card">
        <div className="stats-loading">
          <div className="loading-spinner"></div>
          <p>Загрузка данных...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="stats-card">
        <div className="stats-error">
          <p>{error}</p>
          <button 
            className="retry-button"
            onClick={() => setLoading(true)}
          >
            Повторить попытку
          </button>
        </div>
      </div>
    );
  }

  if (!stats) {
    return (
      <div className="stats-card">
        <div className="stats-error">
          <p>Данные проекта отсутствуют</p>
        </div>
      </div>
    );
  }

  return (
    <>
      <div className="stats-card">
        <div className="stats-header">
          <h3 className="stats-title">{projectName || 'Без названия'}</h3>
          <span className="stats-key">{projectKey || 'N/A'}</span>
        </div>
      
        <div className="stats-grid">
          <div className="stat-item">
            <FaTasks className="stat-icon total" />
            <div className="stat-info">
              <span className="stat-label">Всего задач</span>
              <span className="stat-value">{stats.total_issues ?? 0}</span>
            </div>
          </div>
          
          <div className="stat-item">
            <FaFolderOpen className="stat-icon open" />
            <div className="stat-info">
              <span className="stat-label">Открытых</span>
              <span className="stat-value">{stats.open_issues ?? 0}</span>
            </div>
          </div>
          
          <div className="stat-item">
            <FaCheckCircle className="stat-icon closed" />
            <div className="stat-info">
              <span className="stat-label">Закрытых</span>
              <span className="stat-value">{stats.closed_issues ?? 0}</span>
            </div>
          </div>
          
          <div className="stat-item">
            <FaLockOpen className="stat-icon resolved" />
            <div className="stat-info">
              <span className="stat-label">Разрешенных</span>
              <span className="stat-value">{stats.resolved_issues ?? 0}</span>
            </div>
          </div>
          
          <div className="stat-item">
            <FaRedo className="stat-icon reopened" />
            <div className="stat-info">
              <span className="stat-label">Переоткрытых</span>
              <span className="stat-value">{stats.reopened_issues ?? 0}</span>
            </div>
          </div>
          
          <div className="stat-item">
            <FaClock className="stat-icon in-progress" />
            <div className="stat-info">
              <span className="stat-label">В процессе</span>
              <span className="stat-value">{stats.in_progress_issues ?? 0}</span>
            </div>
          </div>
          
          <div className="stat-item highlight time-card">
            <div className="time-visualization">
              <GiSandsOfTime className="time-icon" />
              <div className="time-value-container">
                <span className="time-value">
                  {formatTime(stats.avg_resolution_time_h)}
                </span>
                <span className="time-label">часов</span>
              </div>
            </div>
            <div className="stat-info">
              <span className="stat-label">Ср. время выполнения</span>
              <span className="stat-details">на задачу</span>
            </div>
          </div>
          
          <div className="stat-item highlight calendar-card">
            <div className="calendar-visualization">
              <FaCalendarAlt className="calendar-icon" />
              <span className="calendar-value">
                {stats.avg_created_per_day_7d ?? 0}
              </span>
            </div>
            <div className="stat-info">
              <span className="stat-label">Задач в день</span>
              <span className="stat-unit">в среднем</span>
            </div>
          </div>
        </div>
        
        <button 
          className="analytics-button"
          onClick={() => setShowAnalytics(true)}
          disabled={!projectKey}
        >
          <FaChartBar className="analytics-icon" />
          Показать аналитику
        </button>
      </div>
      
      <Modal isOpen={showAnalytics} onClose={() => setShowAnalytics(false)}>
        {projectKey && <ChartSelector projectKey={projectKey} />}
      </Modal>
    </>
  );
};

export default StatsCard;