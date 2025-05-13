import React, { useEffect, useState } from 'react';
import axios from 'axios';
import StatsCard from '../../components/StatsCard/StatsCard';
import { Project } from '../../types/models';
import { FiFolder, FiSearch } from 'react-icons/fi';
import './MyProjectsPage.scss';

const MyProjectsPage: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [filteredProjects, setFilteredProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');

  // Загрузка проектов
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

  // Фильтрация проектов
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

  return (
    <div className="my-projects">
      <div className="my-projects-header">
        <h1 className="my-projects-title">Мои проекты</h1>
        <input
          type="text"
          className="my-projects-search"
          placeholder="Поиск проектов..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </div>

      {loading ? (
        <div className="my-projects-loading">
          <div className="loading-spinner"></div>
          Загрузка проектов...
        </div>
      ) : filteredProjects.length === 0 ? (
          <div className="my-projects-empty">
            {search ? (
              <>
                <FiSearch className="empty-icon search-icon" />
                <p>Ничего не найдено по запросу "{search}"</p>
              </>
            ) : (
              <>
                <FiFolder className="empty-icon folder-icon" />
                <p>У вас пока нет проектов</p>
              </>
            )}
            {search && (
              <button 
                className="clear-search"
                onClick={() => setSearch('')}
              >
                Очистить поиск
              </button>
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