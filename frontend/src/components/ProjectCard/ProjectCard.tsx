import React, { useState } from 'react';
import axios from 'axios';
import { Project } from '../../types/models';
import './ProjectCard.scss';
import { config } from '../../config/config';

interface ProjectCardProps extends Project {
  onUpdate: () => void;
  isAdded: boolean;
  Id: number;
  Name: string; 
  Key: string; 
  self: string;
}

const ProjectCard: React.FC<ProjectCardProps> = ({ 
  Id, Name, Key, self, onUpdate, isAdded 
}) => {
  const [loading, setLoading] = useState(false);

  const handleAction = async () => {
    console.log(self)
    setLoading(true);
    try {
      if (isAdded) {
        const response = await axios.get(config.api.endpoints.projects);
        const projectsFromServer = response.data; // Предполагаем, что это массив проектов вида { id: string, key: string, ... }
  
        const projectToDelete = projectsFromServer.find((p: any) => p.key === Key);
        
        if (!projectToDelete) {
          throw new Error("Проект не найден на сервере");
        }
  
        await axios.delete(config.api.endpoints.deleteProject(projectToDelete.id));
      } else {
        await axios.post(
          config.api.endpoints.updateProject,
          null,
          { params: { project: Key } }
        );
      }
      onUpdate();
      console.log(self)
    } catch (error) {
      console.error('Error:', error);
      alert(`Не удалось ${isAdded ? 'удалить' : 'добавить'} проект`);
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
            href={self || "#"} 
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