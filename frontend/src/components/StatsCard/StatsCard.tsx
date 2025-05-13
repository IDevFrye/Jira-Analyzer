import React, { useEffect, useState } from 'react';
import axios from 'axios';
import './StatsCard.scss';
import {
  FaTasks,
  FaFolderOpen,
  FaCheckCircle,
  FaLockOpen,
  FaRedo,
  FaClock,
  FaCalendarDay
} from 'react-icons/fa';
import { FaCalendarAlt } from 'react-icons/fa';
import { GiSandsOfTime } from 'react-icons/gi';


interface StatsCardProps {
  projectId: number;
}

const StatsCard: React.FC<StatsCardProps> = ({ projectId }) => {
  const [stats, setStats] = useState<any>(null);

  useEffect(() => {
    axios.get(`/api/v1/projects/${projectId}`).then((res) => {
      setStats(res.data);
    });
  }, [projectId]);

  if (!stats) return (
    <div className="stats-loading">
      <div className="loading-spinner"></div>
      Загрузка данных...
    </div>
  );

  return (
    <div className="stats-card">
      <div className="stats-header">
        <h3 className="stats-title">{stats.Name}</h3>
        <span className="stats-key">{stats.Key}</span>
      </div>
      
      <div className="stats-grid">
        <div className="stat-item">
          <FaTasks className="stat-icon total" />
          <div className="stat-info">
            <span className="stat-label">Всего задач</span>
            <span className="stat-value">{stats.allIssuesCount}</span>
          </div>
        </div>
        
        <div className="stat-item">
          <FaFolderOpen className="stat-icon open" />
          <div className="stat-info">
            <span className="stat-label">Открытых</span>
            <span className="stat-value">{stats.openIssuesCount}</span>
          </div>
        </div>
        
        <div className="stat-item">
          <FaCheckCircle className="stat-icon closed" />
          <div className="stat-info">
            <span className="stat-label">Закрытых</span>
            <span className="stat-value">{stats.closeIssuesCount}</span>
          </div>
        </div>
        
        <div className="stat-item">
          <FaLockOpen className="stat-icon resolved" />
          <div className="stat-info">
            <span className="stat-label">Разрешенных</span>
            <span className="stat-value">{stats.resolvedIssuesCount}</span>
          </div>
        </div>
        
        <div className="stat-item">
          <FaRedo className="stat-icon reopened" />
          <div className="stat-info">
            <span className="stat-label">Переоткрытых</span>
            <span className="stat-value">{stats.reopenedIssuesCount}</span>
          </div>
        </div>
        
        <div className="stat-item">
          <FaClock className="stat-icon in-progress" />
          <div className="stat-info">
            <span className="stat-label">В процессе</span>
            <span className="stat-value">{stats.progressIssuesCount}</span>
          </div>
        </div>
        
        <div className="stat-item highlight time-card">
          <div className="time-visualization">
            <GiSandsOfTime className="time-icon" />
            <div className="time-value-container">
              <span className="time-value">{stats.averageTime.toFixed(1)}</span>
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
            <span className="calendar-value">{stats.averageIssuesCount}</span>
          </div>
          <div className="stat-info">
            <span className="stat-label">Задач в день</span>
            <span className="stat-unit">в среднем</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default StatsCard;