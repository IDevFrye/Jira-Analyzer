import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Project } from '../../types/models';
import CompareSelector from '../../components/CompareSelector/CompareSelector';
import CompareModal from '../../components/CompareModal/CompareModal';
import './ComparePage.scss';

const ComparePage: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedProjects, setSelectedProjects] = useState<Project[]>([]);
  const [showModal, setShowModal] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchProjects = async () => {
      try {
        const response = await axios.get('/api/v1/projects');
        setProjects(response.data);
      } catch (error) {
        console.error('Error fetching projects:', error);
      } finally {
        setLoading(false);
      }
    };
    
    fetchProjects();
  }, []);

  const toggleProject = (project: Project) => {
    if (selectedProjects.some(p => p.Id === project.Id)) {
      setSelectedProjects(selectedProjects.filter(p => p.Id !== project.Id));
    } else if (selectedProjects.length < 3) {
      setSelectedProjects([...selectedProjects, project]);
    }
  };

  return (
    <div className="compare-page">
      <h1 className="page-title">Сравнение проектов</h1>
      
      {loading ? (
        <div className="loading-state">
          <div className="loading-spinner"></div>
          <p>Загрузка проектов...</p>
        </div>
      ) : (
        <>
          <CompareSelector
            projects={projects}
            selectedProjects={selectedProjects}
            onSelect={toggleProject}
          />
          
          {selectedProjects.length >= 2 && (
            <button 
              className="compare-button"
              onClick={() => setShowModal(true)}
            >
              Сравнить выбранные проекты ({selectedProjects.length})
            </button>
          )}
          
          {showModal && (
            <CompareModal 
              projects={selectedProjects}
              onClose={() => setShowModal(false)}
            />
          )}
        </>
      )}
    </div>
  );
};

export default ComparePage;