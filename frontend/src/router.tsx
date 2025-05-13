import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import ProjectsPage from './pages/ProjectsPage/ProjectsPage';
import MyProjectsPage from './pages/MyProjectsPage/MyProjectsPage';
import TasksPage from './pages/TasksPage/TasksPage';
import ComparePage from './pages/ComparePage/ComparePage';

const Router = () => (
  <Routes>
    <Route path="/" element={<Navigate to="/projects" />} />
    <Route path="/projects" element={<ProjectsPage />} />
    <Route path="/my-projects" element={<MyProjectsPage />} />
    <Route path="/issues" element={<TasksPage />} />
    <Route path="/compare" element={<ComparePage />} />
  </Routes>
);

export default Router;
