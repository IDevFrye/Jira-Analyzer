import React from 'react';
import { Project } from '../../types/models';
import './CompareSelector.scss';

interface CompareSelectorProps {
  projects: Project[];
  selectedProjects: Project[];
  onSelect: (project: Project) => void;
}

const CompareSelector: React.FC<CompareSelectorProps> = ({ 
  projects, 
  selectedProjects, 
  onSelect 
}) => {
  return (
    <div className="compare-selector">
      <div className="projects-list">
        <h3>Доступные проекты</h3>
        <ul>
          {projects.map(project => (
            <li 
              key={project.Id} 
              className={selectedProjects.some(p => p.Id === project.Id) ? 'selected' : ''}
              onClick={() => onSelect(project)}
            >
              {project.Name} ({project.Key})
            </li>
          ))}
        </ul>
      </div>
      
      <div className="selected-projects">
        <h3>Выбранные для сравнения</h3>
        {selectedProjects.length === 0 ? (
          <p className="empty-message">Выберите для сравнения от 2 до 3 проектов</p>
        ) : (
          <ul>
            {selectedProjects.map(project => (
              <li 
                key={project.Id}
                onClick={() => onSelect(project)}
              >
                {project.Name} ({project.Key})
                <span className="remove-btn">×</span>
              </li>
            ))}
            
          </ul>
        )}
      </div>
    </div>
  );
};

export default CompareSelector;