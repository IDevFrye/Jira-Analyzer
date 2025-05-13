import React, { useEffect, useState } from 'react';
import axios from 'axios';
import AnalyticsChart from '../../components/AnalyticsChart/AnalyticsChart';
import { Project } from '../../types/models';

const TASKS = [
  { id: 1, name: 'Время в открытом состоянии' },
  { id: 2, name: 'Распределение по статусам' },
  { id: 3, name: 'Залогированное время' },
  { id: 4, name: 'Распределение по приоритетам' },
];

const TasksPage: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedProject, setSelectedProject] = useState('');
  const [selectedTasks, setSelectedTasks] = useState<number[]>([]);

  useEffect(() => {
    axios.get('/api/v1/projects').then((res) => {
      setProjects(res.data);
    });
  }, []);

  const toggleTask = (id: number) => {
    setSelectedTasks((prev) =>
      prev.includes(id) ? prev.filter((t) => t !== id) : [...prev, id]
    );
  };

  const handleProcess = () => {
    selectedTasks.forEach((taskId) => {
      axios.post(`/api/v1/graph/make/${taskId}`, {
        project: selectedProject,
      });
    });
  };

  return (
    <div>
      <h1>Аналитика задач</h1>
      <select onChange={(e) => setSelectedProject(e.target.value)}>
        <option value=''>Выберите проект</option>
        {projects.map((p) => (
          <option key={p.Id} value={p.Key}>{p.Name}</option>
        ))}
      </select>

      <div>
        {TASKS.map((t) => (
          <label key={t.id}>
            <input
              type="checkbox"
              checked={selectedTasks.includes(t.id)}
              onChange={() => toggleTask(t.id)}
            />
            {t.name}
          </label>
        ))}
      </div>

      <button
        disabled={selectedTasks.length === 0 || !selectedProject}
        onClick={handleProcess}
      >
        Обработать
      </button>

      <div>
        {selectedTasks.map((taskId) => (
          <AnalyticsChart key={taskId} taskId={taskId} projectKey={selectedProject} />
        ))}
      </div>
    </div>
  );
};

export default TasksPage;
