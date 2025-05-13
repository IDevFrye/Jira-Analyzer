// Layout.tsx
import React from 'react';
import Header from '../Header/Header';
import Sidebar from '../Sidebar/Sidebar';
import './Layout.scss';

const Layout: React.FC<{children: React.ReactNode}> = ({ children }) => {
  return (
    <div className="layout">
      <Header />
      <div className="layout-content">
        {/* <Sidebar /> */}
        <main className="layout-main">
          {children}
        </main>
      </div>
    </div>
  );
};

export default Layout;