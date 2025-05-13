import React from 'react';
import './Sidebar.scss';

const Sidebar: React.FC = () => {
  return (
    <aside className="sidebar">
      <div className="sidebar-section">
        <h3 className="sidebar-title">Фильтры</h3>
        <div className="sidebar-filters">
          <label className="sidebar-filter">
            <input type="checkbox" className="sidebar-checkbox" />
            <span>Активные проекты</span>
          </label>
          <label className="sidebar-filter">
            <input type="checkbox" className="sidebar-checkbox" />
            <span>Архивные</span>
          </label>
        </div>
      </div>
      
      <div className="sidebar-section">
        <h3 className="sidebar-title">Быстрый доступ</h3>
        <nav className="sidebar-nav">
          <a href="#" className="sidebar-link">Избранное</a>
          <a href="#" className="sidebar-link">Недавние</a>
          <a href="#" className="sidebar-link">Шаблоны</a>
        </nav>
      </div>
      
      <div className="sidebar-section">
        <h3 className="sidebar-title">Статистика</h3>
        <div className="sidebar-stats">
          <div className="sidebar-stat">
            <span>Проектов</span>
            <span className="sidebar-stat-value">24</span>
          </div>
          <div className="sidebar-stat">
            <span>Задач</span>
            <span className="sidebar-stat-value">156</span>
          </div>
        </div>
      </div>
    </aside>
  );
};

export default Sidebar;