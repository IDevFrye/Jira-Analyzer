import React, { useEffect, useState } from 'react';
import axios from 'axios';
import StatsCard from '../../components/StatsCard/StatsCard';
import { Project } from '../../types/models';
import { FiFolder, FiSearch, FiX } from 'react-icons/fi';
import './MyProjectsPage.scss';

const MyProjectsPage: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [filteredProjects, setFilteredProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [isSearchFocused, setIsSearchFocused] = useState(false);

  useEffect(() => {
    setLoading(true);
    axios.get('/api/v1/projects')
      .then((res) => {
        setProjects(res.data);
        setFilteredProjects(res.data);
        setLoading(false);
      })
      .catch(() => {
        setLoading(false);
      });
  }, []);

  useEffect(() => {
    if (search.trim() === '') {
      setFilteredProjects(projects);
    } else {
      const filtered = projects.filter(project => 
        project.Name.toLowerCase().includes(search.toLowerCase()) ||
        project.Key.toLowerCase().includes(search.toLowerCase())
      );
      setFilteredProjects(filtered);
    }
  }, [search, projects]);

  const handleClearSearch = () => {
    setSearch('');
  };

  return (
    <div className="my-projects">
      <div className="my-projects-header">
        <h1 className="my-projects-title">Мои проекты</h1>
        <div className={`search-container ${isSearchFocused ? 'focused' : ''}`}>
          <FiSearch className="search-icon" />
          <input
            type="text"
            className="my-projects-search"
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
      ) : filteredProjects.length === 0 ? (
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
              <p>У вас пока нет проектов</p>
            </>
          )}
        </div>
      ) : (
        <div className="stats-cards-container">
          {filteredProjects.map((project) => (
            <StatsCard key={project.Id} projectId={project.Id} />
          ))}
        </div>
      )}
    </div>
  );
};

export default MyProjectsPage;