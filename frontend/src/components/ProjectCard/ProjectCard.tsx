import React, { useState } from 'react';
import axios from 'axios';
import { Project } from '../../types/models';
import './ProjectCard.scss';

interface ProjectCardProps extends Project {
  onUpdate: () => void;
}

const ProjectCard: React.FC<ProjectCardProps> = ({ Id, Name, Key, Url, onUpdate }) => {
  const [isAdded, setIsAdded] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleAction = async () => {
    setLoading(true);
    try {
      if (isAdded) {
        await axios.delete(`/api/v1/projects/${Id}`);
      } else {
        await axios.post('/api/v1/connector/updateProject', { project: Key });
      }
      setIsAdded(!isAdded);
      onUpdate();
    } catch (error) {
      console.error('Error:', error);
    } finally {
      setLoading(false);
    }
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
            onClick={handleAction}
            className={`project-action-button ${isAdded ? 'remove' : 'add'}`}
            disabled={loading}
          >
            {loading ? (
              'Загрузка...'
            ) : isAdded ? (
              'Удалить'
            ) : (
              'Добавить'
            )}
          </button>
        </div>
      </div>
    </div>
  );
};

export default ProjectCard;