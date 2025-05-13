import React from 'react';
import ReactDOM from 'react-dom/client';
import JiraAnalyzerApp from './JiraAnalyzerApp';
// import './styles/main.css';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
  <React.StrictMode>
    <JiraAnalyzerApp />
  </React.StrictMode>
);