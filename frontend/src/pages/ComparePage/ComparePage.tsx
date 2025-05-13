import React, { useState } from 'react';
import CompareTable from '../../components/CompareTable/CompareTable';

const ComparePage: React.FC = () => {
  const [selectedProjects, setSelectedProjects] = useState<string[]>([]);

  const toggleProject = (key: string) => {
    if (selectedProjects.includes(key)) {
      setSelectedProjects(selectedProjects.filter((k) => k !== key));
    } else if (selectedProjects.length < 3) {
      setSelectedProjects([...selectedProjects, key]);
    } else {
      alert('Можно выбрать максимум 3 проекта');
    }
  };

  return (
    <div>
      <h1>Сравнение проектов</h1>
      <CompareTable selectedProjectKeys={selectedProjects} onSelect={toggleProject} />
    </div>
  );
};

export default ComparePage;