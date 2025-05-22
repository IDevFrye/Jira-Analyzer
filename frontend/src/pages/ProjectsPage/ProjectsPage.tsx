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

  useEffect(() => {
    const fetchProjectsAndAdded = async () => {
      setLoading(true);
      try {
        const connectorRes = await axios.get(config.api.endpoints.connectorProjects, {
          params: { page, limit: 9, search },
        });
  
        const rawProjects = connectorRes.data?.projects || [];
  
        const formattedProjects = rawProjects.map((p: any) => ({
          Id: p.id,
          Key: p.key,
          Name: p.name,
          self: p.self,
        }));
  
        const addedRes = await axios.get(config.api.endpoints.projects);
        const added = addedRes.data;
  
        const addedKeys = new Set<string>();
        const selfMap = new Map<string, string>();
  
        added?.forEach((p: { key: string; self: string }) => {
          addedKeys.add(p.key);
          selfMap.set(p.key, p.self);
        });
  
        const mergedProjects = formattedProjects.map((p: Project) => ({
          ...p,
          self: selfMap.get(p.Key) || p.self
        }));
        
        setProjects(mergedProjects);
        setAddedProjects(addedKeys);
        setPageCount(connectorRes.data?.pageInfo?.pageCount || 1);
      } catch (error) {
        console.error("Error loading projects:", error);
        setProjects([]);
      } finally {
        setLoading(false);
      }
    };
  
    fetchProjectsAndAdded();
  }, [page, search]);
  

  const handleProjectAction = async (projectName: string, action: 'add' | 'remove') => {
    try {
      if (action === 'remove') {
        const { data } = await axios.get<Array<{ id: string; name: string, key: string,  self: string }>>(config.api.endpoints.projects);
        const projectToDelete = data.find((p) => {
          return p.key === projectName;
          
        });

        if (projectToDelete) {
          await axios.delete(config.api.endpoints.deleteProject(Number(projectToDelete.id)));
        }
        setAddedProjects(prev => {
          const newSet = new Set<string>(prev);
          newSet.delete(projectName);
          return newSet;
        });
      } else {
        await axios.post(
          config.api.endpoints.updateProject,
          null,
          { params: { project: projectName } }
        );
        setAddedProjects(prev => {
          const newSet = new Set<string>(prev);
          newSet.add(projectName);
          return newSet;
        });
      }
    } catch (error) {
      console.error('Action failed:', error);
      alert(`Не удалось ${action === 'add' ? 'добавить' : 'удалить'} проект`);
      const { data } = await axios.get<Array<{ name: string, key: string }>>(config.api.endpoints.projects);
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