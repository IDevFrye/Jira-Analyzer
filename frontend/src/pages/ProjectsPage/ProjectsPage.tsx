import React, { useEffect, useState } from 'react';
import axios from 'axios';
import ProjectCard from '../../components/ProjectCard/ProjectCard';
import { Project } from '../../types/models';
import { FiSearch, FiFolder, FiX } from 'react-icons/fi';
import './ProjectsPage.scss';

const ProjectsPage: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const [pageCount, setPageCount] = useState(1);
  const [loading, setLoading] = useState(true);
  const [isSearchFocused, setIsSearchFocused] = useState(false);

  useEffect(() => {
    setLoading(true);
    axios
      .get(`/api/v1/connector/projects`, {
        params: { page, limit: 9, search },
      })
      .then((res) => {
        setProjects(res.data.Projects);
        setPageCount(res.data.PageInfo.pageCount);
        setLoading(false);
      })
      .catch(() => {
        setLoading(false);
      });
  }, [page, search]);

  const handleClearSearch = () => {
    setSearch('');
  };

  return (
    <div className="projects-container">
      <div className="projects-header">
        <h1 className="projects-title">Все проекты</h1>
        <div className={`search-container ${isSearchFocused ? 'focused' : ''}`}>
          <FiSearch className="search-icon" />
          <input
            type="text"
            className="projects-search"
            placeholder="Поиск проектов..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            onFocus={() => setIsSearchFocused(true)}
            onBlur={() => setIsSearchFocused(false)}
          />
          {search && (
            <button className="clear-search-btn" onClick={handleClearSearch}>
              <FiX />
            </button>
          )}
        </div>
      </div>
      
      {loading ? (
        <div className="loading-state">
          <div className="loading-spinner"></div>
          <p>Загрузка проектов...</p>
        </div>
      ) : projects.length === 0 ? (
        <div className="empty-state">
          {search ? (
            <>
              <div className="search-empty-content">
                <FiSearch className="empty-icon" />
                <p>Не найдено проектов по запросу <strong>"{search}"</strong></p>
                <button 
                  className="clear-search-btn large"
                  onClick={handleClearSearch}
                >
                  Очистить поиск
                </button>
              </div>
            </>
          ) : (
            <>
              <FiFolder className="empty-icon" />
              <p>Нет доступных проектов</p>
            </>
          )}
        </div>
      ) : (
        <>
          <div className="projects-grid">
            {projects.map((p) => (
              <ProjectCard
                key={p.Id}
                Id={p.Id}
                Name={p.Name}
                Key={p.Key}
                Url={p.Url}
                onUpdate={() => {}}
              />
            ))}
          </div>
          
          <div className="projects-pagination">
            <button 
              className="pagination-button" 
              disabled={page === 1} 
              onClick={() => setPage(page - 1)}
            >
              &larr; Назад
            </button>
            <span className="pagination-info">{page} / {pageCount}</span>
            <button 
              className="pagination-button" 
              disabled={page === pageCount} 
              onClick={() => setPage(page + 1)}
            >
              Вперед &rarr;
            </button>
          </div>
        </>
      )}
    </div>
  );
};

export default ProjectsPage;