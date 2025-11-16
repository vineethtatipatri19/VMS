import React, { useState, useEffect } from 'react';
import { dashboardAPI } from '../services/api';
import { Link } from 'react-router-dom';

function Dashboard() {
  const [stats, setStats] = useState(null);
  const [activity, setActivity] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [statsRes, activityRes] = await Promise.all([
        dashboardAPI.getStats(),
        dashboardAPI.getActivity()
      ]);
      setStats(statsRes.data);
      setActivity(activityRes.data);
    } catch (err) {
      setError('Failed to load dashboard data');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div className="loading">Loading dashboard...</div>;
  if (error) return <div className="error">{error}</div>;

  return (
    <div className="container">
      <h1>Dashboard</h1>
      
      <div className="stats-grid">
        <div className="stat-card">
          <h3>Total Customers</h3>
          <div className="value">{stats?.totalCustomers || 0}</div>
        </div>
        <div className="stat-card">
          <h3>Expiring Soon</h3>
          <div className="value" style={{color: '#f39c12'}}>{stats?.expiringSoonItems || 0}</div>
        </div>
        <div className="stat-card">
          <h3>Expired Items</h3>
          <div className="value" style={{color: '#e74c3c'}}>{stats?.expiredItems || 0}</div>
        </div>
        <div className="stat-card">
          <h3>Unreturned Crates</h3>
          <div className="value">{stats?.unreturnedCrates || 0}</div>
        </div>
        <div className="stat-card">
          <h3>Outstanding Balance</h3>
          <div className="value">₹{stats?.outstandingBalance?.toFixed(2) || '0.00'}</div>
        </div>
        <div className="stat-card">
          <h3>Today's Sales</h3>
          <div className="value" style={{color: '#27ae60'}}>₹{stats?.todaysSales?.toFixed(2) || '0.00'}</div>
        </div>
        <div className="stat-card">
          <h3>Month Sales</h3>
          <div className="value" style={{color: '#3498db'}}>₹{stats?.monthSales?.toFixed(2) || '0.00'}</div>
        </div>
      </div>

      <div className="card">
        <h2>Quick Actions</h2>
        <div style={{ display: 'flex', gap: '10px', flexWrap: 'wrap', marginTop: '15px' }}>
          <Link to="/inventory" className="btn btn-primary">Add Inventory</Link>
          <Link to="/customers" className="btn btn-primary">Add Customer</Link>
          <Link to="/transactions" className="btn btn-success">New Sale</Link>
          <Link to="/forecasting" className="btn btn-info">AI Forecast</Link>
          <Link to="/reports" className="btn btn-secondary">Generate Report</Link>
        </div>
      </div>

      {activity && activity.length > 0 && (
        <div className="card">
          <h2>Recent Activity</h2>
          <table className="table">
            <thead>
              <tr>
                <th>Type</th>
                <th>Description</th>
                <th>Time</th>
              </tr>
            </thead>
            <tbody>
              {activity.map((item, idx) => (
                <tr key={idx}>
                  <td>
                    <span className={`badge ${item.type === 'sale' ? 'badge-success' : 'badge-info'}`}>
                      {item.type}
                    </span>
                  </td>
                  <td>{item.description}</td>
                  <td>{new Date(item.timestamp).toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

export default Dashboard;
