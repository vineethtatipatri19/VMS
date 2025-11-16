import React from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './Layout.css';

function Layout({ children }) {
  const { logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const isActive = (path) => location.pathname === path;

  return (
    <div className="layout">
      <nav className="sidebar">
        <div className="sidebar-header">
          <h2>PGVMS</h2>
          <p>Vendor Management</p>
        </div>
        <ul className="nav-menu">
          <li className={isActive('/dashboard') || isActive('/') ? 'active' : ''}>
            <Link to="/dashboard">📊 Dashboard</Link>
          </li>
          <li className={isActive('/inventory') ? 'active' : ''}>
            <Link to="/inventory">📦 Inventory</Link>
          </li>
          <li className={isActive('/customers') ? 'active' : ''}>
            <Link to="/customers">👥 Customers</Link>
          </li>
          <li className={isActive('/transactions') ? 'active' : ''}>
            <Link to="/transactions">💰 Transactions</Link>
          </li>
          <li className={isActive('/crates') ? 'active' : ''}>
            <Link to="/crates">📦 Crates</Link>
          </li>
          <li className={isActive('/forecasting') ? 'active' : ''}>
            <Link to="/forecasting">🔮 Forecasting</Link>
          </li>
          <li className={isActive('/reports') ? 'active' : ''}>
            <Link to="/reports">📄 Reports</Link>
          </li>
        </ul>
        <div className="sidebar-footer">
          <button onClick={handleLogout} className="btn btn-secondary logout-btn">
            Logout
          </button>
        </div>
      </nav>
      <main className="main-content">
        {children}
      </main>
    </div>
  );
}

export default Layout;
