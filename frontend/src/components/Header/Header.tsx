import React from 'react';
import { Link } from 'react-router-dom';
import logo from '../../assets/logo.png';
import './Header.scss';

const Header: React.FC = () => {
  return (
    <header className="header">
      <nav className="header-nav">
        <ul className="header-list">

            <img src={logo} alt="Логотип" className="header-logo"></img>
          <li className="header-item">
            <Link to="/projects" className="header-link">Проекты</Link>
          </li>
          <li className="header-item">
            <Link to="/tasks" className="header-link">Задачи</Link>
          </li>
          <li className="header-item">
            <Link to="/compare" className="header-link">Сравнение</Link>
          </li>
          <li className="header-item">
            <Link to="/my-projects" className="header-link">Мои проекты</Link>
          </li>
        </ul>
      </nav>
    </header>
  );
};

export default Header;