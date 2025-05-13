import React from 'react';
import axios from 'axios';
import { Project } from '../../types/models';
import './ProjectCard.scss';

interface ProjectCardProps extends Project {
  onAdd: (id: number) => void;
}

const ProjectCard: React.FC<ProjectCardProps> = ({ Id, Name, Key, Url, onAdd }) => {
  const handleAdd = () => {
    axios.post('/api/v1/connector/updateProject', { project: Key });
  };

  return (
    <div className="project-card">
      <div className="project-content">
        <h3 className="project-name">{Name}</h3>
        <span className="project-key">{Key}</span>
        <div className="project-footer">
          <a 
            href={Url} 
            target="_blank" 
            rel="noopener noreferrer"
            className="project-link"
          >
            Перейти
          </a>
          <button 
            onClick={handleAdd}
            className="project-add-button"
          >
            Добавить
          </button>
        </div>
      </div>
    </div>
  );
};

export default ProjectCard;