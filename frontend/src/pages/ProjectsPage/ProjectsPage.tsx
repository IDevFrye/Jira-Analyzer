import React, { useEffect, useState } from 'react';
import axios from 'axios';
import ProjectCard from '../../components/ProjectCard/ProjectCard';
import { Project } from '../../types/models';
import './ProjectsPage.scss';

const ProjectsPage: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const [pageCount, setPageCount] = useState(1);

  useEffect(() => {
    axios
      .get(`/api/v1/connector/projects`, {
        params: { page, limit: 9, search },
      })
      .then((res) => {
        setProjects(res.data.Projects);
        setPageCount(res.data.PageInfo.pageCount);
      });
  }, [page, search]);

  return (
    <div className="projects-container">
      <div className="projects-header">
        <h1 className="projects-title">All Projects</h1>
        <input
          type="text"
          className="projects-search"
          placeholder="Search projects..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </div>
      
      <div className="projects-grid">
        {projects.map((p) => (
          <ProjectCard
            key={p.Id}
            Id={p.Id}
            Name={p.Name}
            Key={p.Key}
            Url={p.Url}
            onAdd={() => {}}
          />
        ))}
      </div>
      
      <div className="projects-pagination">
        <button 
          className="pagination-button" 
          disabled={page === 1} 
          onClick={() => setPage(page - 1)}
        >
          &larr; Previous
        </button>
        <span className="pagination-info">{page} / {pageCount}</span>
        <button 
          className="pagination-button" 
          disabled={page === pageCount} 
          onClick={() => setPage(page + 1)}
        >
          Next &rarr;
        </button>
      </div>
    </div>
  );
};

export default ProjectsPage;