import React, { useEffect, useState } from 'react';
import axios from 'axios';
import ProjectCard from '../../components/ProjectCard/ProjectCard';
import { Project } from '../../types/models';
import { FiSearch, FiFolder, FiX } from 'react-icons/fi';
import './ProjectsPage.scss';
import { config } from '../../config/config';

const ProjectsPage: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const [pageCount, setPageCount] = useState(1);
  const [loading, setLoading] = useState(false);
  const [isSearchFocused, setIsSearchFocused] = useState(false);
  const [addedProjects, setAddedProjects] = useState<Set<string>>(new Set());

  // Загрузка проектов
  useEffect(() => {
    const fetchProjects = async () => {
      setLoading(true);
      try {
        const res = await axios.get(config.api.endpoints.connectorProjects, {
          params: { page, limit: 9, search },
        });

        const formattedProjects = res.data?.projects?.map((p: any) => ({
          Id: p.id,
          Key: p.key,
          Name: p.name,
          self: p.self
        })) || [];

        setProjects(formattedProjects);
        setPageCount(res.data?.pageInfo?.pageCount || 1);
      } catch (error) {
        console.error("Error fetching projects:", error);
        setProjects([]);
      } finally {
        setLoading(false);
      }
    };

    fetchProjects();
  }, [page, search]);

  // Загрузка добавленных проектов
  useEffect(() => {
    const fetchAddedProjects = async () => {
      try {
        const res = await axios.get(config.api.endpoints.projects);
        // Явно указываем тип для данных
        const addedKeys = new Set<string>(res.data.map((p: { key: string }) => p.key));
        setAddedProjects(addedKeys);
      } catch (error) {
        console.error("Error fetching added projects:", error);
      }
    };
  
    fetchAddedProjects();
  }, []);

  // Обработчик действий с проектом
  const handleProjectAction = async (projectKey: string, action: 'add' | 'remove') => {
    try {
      if (action === 'remove') {
        const { data } = await axios.get<Array<{ id: string; key: string }>>(config.api.endpoints.projects);
        const projectToDelete = data.find((p) => p.key === projectKey);
        
        if (projectToDelete) {
          await axios.delete(config.api.endpoints.deleteProject(Number(projectToDelete.id)));
        }
  
        setAddedProjects(prev => {
          const newSet = new Set<string>(prev); // Явно указываем тип Set<string>
          newSet.delete(projectKey);
          return newSet;
        });
      } else {
        await axios.post(
          config.api.endpoints.updateProject,
          null,
          { params: { project: projectKey } }
        );
        setAddedProjects(prev => new Set<string>(prev).add(projectKey)); // Явно указываем тип
      }
    } catch (error) {
      console.error('Action failed:', error);
      alert(`Не удалось ${action === 'add' ? 'добавить' : 'удалить'} проект`);
      // Принудительно обновляем состояние с явным указанием типа
      const { data } = await axios.get<Array<{ key: string }>>(config.api.endpoints.projects);
      setAddedProjects(new Set<string>(data.map(p => p.key)));
    }
  };

  const handleClearSearch = () => setSearch('');

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
            onChange={(e) => {
              setSearch(e.target.value);
              setPage(1)
            }}
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
            {projects.map((project) => (
              <ProjectCard
                key={project.Id}
                Id={project.Id}
                Key={project.Key}
                Name={project.Name}
                self={project.self}
                isAdded={addedProjects.has(project.Key)}
                onAction={handleProjectAction}
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