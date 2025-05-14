// Layout.tsx
import React from 'react';
import Header from '../Header/Header';
import './Layout.scss';

const Layout: React.FC<{children: React.ReactNode}> = ({ children }) => {
  return (
    <div className="layout">
      <Header />
      <div className="layout-content">
        <main className="layout-main">
          {children}
        </main>
      </div>
    </div>
  );
};

export default Layout;