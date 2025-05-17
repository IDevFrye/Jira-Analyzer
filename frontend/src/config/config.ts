// Для TypeScript определим тип для process.env
declare global {
  namespace NodeJS {
    interface ProcessEnv {
      REACT_APP_API_BASE_URL?: string;
    }
  }
}

const DEFAULT_API_BASE_URL = 'http://localhost:8000';
const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || DEFAULT_API_BASE_URL;

export const config = {
  api: {
    baseUrl: API_BASE_URL,
    endpoints: {
      // Projects
      projects: `${API_BASE_URL}/api/v1/projects`,
      projectStats: (id: number) => `${API_BASE_URL}/api/v1/projects/${id}`,
      deleteProject: (id: number) => `${API_BASE_URL}/api/v1/projects/${id}`,
      
      // Connector
      connectorProjects: `${API_BASE_URL}/api/v1/connector/projects`,
      updateProject: `${API_BASE_URL}/api/v1/connector/updateProject`,
      
      // Single project analytics
      timeOpenAnalytics: `${API_BASE_URL}/api/v1/analytics/time-open`,
      statusDistribution: `${API_BASE_URL}/api/v1/analytics/status-distribution`,
      timeSpentAnalytics: `${API_BASE_URL}/api/v1/analytics/time-spent`,
      priorityAnalytics: `${API_BASE_URL}/api/v1/analytics/priority`,
      
      // Comparison analytics
      compareTimeOpen: `${API_BASE_URL}/api/v1/compare/time-open`,
      compareStatusDistribution: `${API_BASE_URL}/api/v1/compare/status-distribution`,
      compareTimeSpent: `${API_BASE_URL}/api/v1/compare/time-spent`,
      comparePriority: `${API_BASE_URL}/api/v1/compare/priority`
    }
  }
};