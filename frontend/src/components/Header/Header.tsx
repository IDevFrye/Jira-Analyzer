import React from 'react';
import { Link } from 'react-router-dom';
import logo from '../../assets/logo.png';
import './Header.scss';

const Header: React.FC = () => {
  return (
    <header className="header">
      <nav className="header-nav">
        <ul className="header-list">
          <div className="header-flex">
            <img src={logo} alt="Логотип" className="header-logo"></img>
          </div>
          <div className="header-flex2">
            <li className="header-item">
              <Link to="/projects" className="header-link">Проекты</Link>
            </li>
            <li className="header-item">
              <Link to="/my-projects" className="header-link">Мои проекты</Link>
            </li>
            <li className="header-item">
              <Link to="/compare" className="header-link">Сравнение</Link>
            </li>
          </div>
          <div className="header-flex">
            <span className="header-title">
              Jira Analyzer
              <div className="aurora">
                <div className="aurora__item"></div>
                <div className="aurora__item"></div>
                <div className="aurora__item"></div>
                <div className="aurora__item"></div>
              </div>
            </span>
          </div>
        </ul>
      </nav>
    </header>
  );
};

export default Header;