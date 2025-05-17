import React, { useEffect, useState } from 'react';
import axios from 'axios';
import ProjectCard from '../../components/ProjectCard/ProjectCard';
import { Project } from '../../types/models';
import { FiSearch, FiFolder, FiX } from 'react-icons/fi';
import './ProjectsPage.scss';
import { config } from '../../config/config';

interface ApiResponse {
  Projects?: Project[];
  PageInfo?: {
    pageCount: number;
    currentPage: number;
    projectsCount: number;
  };
}

const ProjectsPage: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const [pageCount, setPageCount] = useState(1);
  const [loading, setLoading] = useState(false);
  const [isSearchFocused, setIsSearchFocused] = useState(false);
const [addedProjects, setAddedProjects] = useState<Set<string>>(() => {
  const saved = localStorage.getItem('addedProjects');
  return saved ? new Set(JSON.parse(saved)) : new Set();
});

useEffect(() => {
  localStorage.setItem('addedProjects', JSON.stringify(Array.from(addedProjects)));
}, [addedProjects]);
  
  const handleToggleAdd = (key: string, isAdded: boolean) => {
    setAddedProjects(prev => {
      const newSet = new Set(prev);
      if (isAdded) {
        newSet.add(key);
      } else {
        newSet.delete(key);
      }
      return newSet;
    });
  };


  useEffect(() => {
    setLoading(true);
    
    axios
      .get(config.api.endpoints.connectorProjects, {
        params: { page, limit: 9, search },
      })
      .then((res) => {
        const { projects, pageInfo } = res.data;
        
        if (!projects) {
          throw new Error("No projects data received");
        }

        const formattedProjects = projects.map((p: any) => ({
          Id: p.id,
          Key: p.key,
          Name: p.name,
          self: p.self
        }));

        setProjects(formattedProjects);
        setPageCount(pageInfo?.pageCount || 1);
        setLoading(false);
      })
      .catch((error) => {
        console.error("Error fetching projects:", error);
        setLoading(false);
        setProjects([]);
        setPageCount(1);
      });
  }, [page, search]);

  const handleClearSearch = () => {
    setSearch('');
  };

  const handleUpdate = async () => {
    try {
      const res = await axios.get(config.api.endpoints.connectorProjects, {
        params: { page, limit: 9, search },
      });
  
      if (!res.data?.projects) {
        throw new Error("Invalid projects data received");
      }
  
      const formattedProjects = res.data.projects.map((p: any) => ({
        Id: p.id,
        Key: p.key,
        Name: p.name,
        self: p.self,
      }));
  
      const addedRes = await axios.get(config.api.endpoints.projects);
      const addedKeys = new Set<string>(addedRes.data.map((p: any) => p.key as string)); 
  
      setProjects(formattedProjects);
      setPageCount(res.data.pageInfo?.pageCount || 1);
      setAddedProjects(addedKeys);
    } catch (error) {
      console.error("Update error:", error);
      alert("Не удалось обновить данные. Пожалуйста, попробуйте позже.");
    }
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
              self={p.self}
              onUpdate={handleUpdate}
              isAdded={addedProjects.has(p.Key)}
            />
            ))}
          </div>
          
          <div className="projects-pagination">
            <button 
              className="pagination-button" 
              disabled={page === 1} 
              onClick={() => {setPage(page - 1);}}
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