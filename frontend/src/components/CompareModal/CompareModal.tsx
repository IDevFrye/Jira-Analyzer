import React, { useState } from 'react';
import { Project } from '../../types/models';
import CompareStatsTable from './CompareStatsTable';
import CompareCharts from './CompareCharts';
import './CompareModal.scss';

interface CompareModalProps {
  projects: Project[];
  onClose: () => void;
}

const CompareModal: React.FC<CompareModalProps> = ({ projects, onClose }) => {
  const [activeTab, setActiveTab] = useState<'stats' | 'charts'>('stats');

  return (
    <div className="modal-overlay">
      <div className="compare-modal">
        <button className="close-button" onClick={onClose}>×</button>
        <h2>Сравнение проектов</h2>
        
        <div className="tabs">
          <button 
            className={activeTab === 'stats' ? 'active' : ''}
            onClick={() => setActiveTab('stats')}
          >
            Сухая статистика
          </button>
          <button 
            className={activeTab === 'charts' ? 'active' : ''}
            onClick={() => setActiveTab('charts')}
          >
            Графики
          </button>
        </div>
        
        <div className="modal-content">
          {activeTab === 'stats' ? (
            <CompareStatsTable projects={projects} />
          ) : (
            <CompareCharts projects={projects} />
          )}
        </div>
      </div>
    </div>
  );
};

export default CompareModal;