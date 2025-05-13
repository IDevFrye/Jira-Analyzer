import React from 'react';
import { BrowserRouter } from 'react-router-dom';
import Router from './router';
import Header from './components/Header/Header';
import Sidebar from './components/Sidebar/Sidebar';
import Layout from './components/Layout/Layout';

const JiraAnalyzerApp = () => (
  <BrowserRouter>
    <div className="layout">
      <main>
        <Layout>
          <Router />
        </Layout>
      </main>
    </div>
  </BrowserRouter>
);

export default JiraAnalyzerApp;
