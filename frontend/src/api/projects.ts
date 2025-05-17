import axios from 'axios';

const API = '/api/v1';

export const getAllProjects = () => axios.get(`${API}/projects`);
export const getProjectStats = (id: number) => axios.get(`${API}/projects/${id}`);
export const deleteProject = (id: number) => axios.delete(`${API}/projects/${id}`);
export const searchAvailableProjects = (params: { limit: number; page: number; search?: string }) =>
  axios.get(`${API}/connector/projects`, { params });

export const downloadProject = (key: string) =>
  axios.post(`${API}/connector/updateProject`, null, { params: { project: key } });
