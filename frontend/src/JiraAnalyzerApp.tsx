import { BrowserRouter } from 'react-router-dom';
import Router from './router';
import Layout from './components/Layout/Layout';
import './config/chartConfig';

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
