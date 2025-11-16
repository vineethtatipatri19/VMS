import React, { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { 
  LayoutDashboard, 
  Package, 
  Users, 
  DollarSign, 
  Box, 
  TrendingUp, 
  FileText,
  Menu,
  X,
  LogOut,
  ChevronLeft,
  User,
  AlertTriangle,
  Bell
} from 'lucide-react';
import './Layout.css';

function Layout({ children }) {
  const { logout, user } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const isActive = (path) => {
    if (path === '/dashboard') {
      return location.pathname === '/' || location.pathname === '/dashboard';
    }
    return location.pathname === path;
  };

  const closeMobileSidebar = () => {
    setSidebarOpen(false);
  };

  const navItems = [
    { path: '/dashboard', label: 'Dashboard', icon: LayoutDashboard },
    { path: '/inventory', label: 'Inventory', icon: Package },
    { path: '/customers', label: 'Customers', icon: Users },
    { path: '/transactions', label: 'Transactions', icon: DollarSign },
    { path: '/crates', label: 'Crates', icon: Box },
    { path: '/wastage', label: 'Wastage', icon: AlertTriangle },
    { path: '/expiry-alerts', label: 'Expiry Alerts', icon: Bell },
    { path: '/forecasting', label: 'Forecasting', icon: TrendingUp },
    { path: '/reports', label: 'Reports', icon: FileText },
  ];

  return (
    <div className="layout">
      {/* Mobile Header */}
      <header className="mobile-header">
        <button 
          className="mobile-menu-button" 
          onClick={() => setSidebarOpen(!sidebarOpen)}
          aria-label="Toggle menu"
        >
          {sidebarOpen ? <X size={24} /> : <Menu size={24} />}
        </button>
        <div className="mobile-header-title">
          <h1>PGVMS</h1>
        </div>
        <button className="mobile-user-button" aria-label="User menu">
          <User size={20} />
        </button>
      </header>

      {/* Sidebar Overlay for Mobile */}
      {sidebarOpen && (
        <div 
          className="sidebar-overlay" 
          onClick={closeMobileSidebar}
        />
      )}

      {/* Sidebar */}
      <nav className={`sidebar ${sidebarCollapsed ? 'sidebar-collapsed' : ''} ${sidebarOpen ? 'sidebar-open' : ''}`}>
        <div className="sidebar-header">
          <div className="sidebar-logo">
            <Box size={32} className="logo-icon" />
            {!sidebarCollapsed && (
              <div className="logo-text">
                <h2>PGVMS</h2>
                <p>Vendor Management</p>
              </div>
            )}
          </div>
          <button 
            className="sidebar-collapse-button desktop-only"
            onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
            aria-label="Toggle sidebar"
          >
            <ChevronLeft size={20} className={sidebarCollapsed ? 'rotated' : ''} />
          </button>
        </div>

        <ul className="nav-menu">
          {navItems.map((item) => {
            const Icon = item.icon;
            const active = isActive(item.path);
            
            return (
              <li key={item.path} className={active ? 'active' : ''}>
                <Link 
                  to={item.path} 
                  onClick={closeMobileSidebar}
                  title={item.label}
                >
                  <Icon size={20} className="nav-icon" />
                  {!sidebarCollapsed && <span className="nav-label">{item.label}</span>}
                </Link>
              </li>
            );
          })}
        </ul>

        <div className="sidebar-footer">
          {!sidebarCollapsed && (
            <div className="user-info">
              <div className="user-avatar">
                <User size={20} />
              </div>
              <div className="user-details">
                <p className="user-name">{user?.email || 'User'}</p>
                <p className="user-role">Admin</p>
              </div>
            </div>
          )}
          <button 
            onClick={handleLogout} 
            className="logout-btn"
            title="Logout"
          >
            <LogOut size={20} />
            {!sidebarCollapsed && <span>Logout</span>}
          </button>
        </div>
      </nav>

      {/* Main Content */}
      <main className="main-content">
        <div className="content-wrapper">
          {children}
        </div>
      </main>
    </div>
  );
}

export default Layout;
