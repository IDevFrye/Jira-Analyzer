import React, { useEffect, useState } from 'react';
import './ProjectCard.scss';
import { config } from '../../config/config';
import axios from 'axios';

interface Project {
  Id: string;
  Key: string;
  Name: string;
  self: string;
}

interface ProjectCardProps {
  Id: number;
  Key: string;
  Name: string;
  self: string;
  isAdded: boolean;
  onAction: (key: string, action: 'add' | 'remove') => Promise<void>;
}

const ProjectCard: React.FC<ProjectCardProps> = ({
  Id,
  Key,
  Name,
  self,
  isAdded,
  onAction
}) => {
  const [loading, setLoading] = useState(false);


  const handleButtonClick = async () => {
    setLoading(true);
    try {
      await onAction(Key, isAdded ? 'remove' : 'add');
    } catch (error) {
      console.error('Action failed:', error);
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
            href={self || '#'}
            target="_blank"
            rel="noopener noreferrer"
            className="project-link"
          >
            Перейти
          </a>
          <button
            className={`project-action-button ${isAdded ? 'remove' : 'add'}`}
            onClick={handleButtonClick}
            disabled={loading}
          >
            {loading ? (
              '...'
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