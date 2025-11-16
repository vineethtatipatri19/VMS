import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { dashboardAPI } from '../services/api';
import { 
  Users, 
  Package, 
  TrendingUp, 
  TrendingDown,
  DollarSign, 
  AlertTriangle,
  Box,
  Activity,
  Plus
} from 'lucide-react';
import { Line, Bar, Doughnut } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js';
import { Card, Button, Badge } from '../components/ui';
import { format } from 'date-fns';
import './Dashboard.css';

// Register ChartJS components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

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
      setActivity(activityRes.data || []);
    } catch (err) {
      setError('Failed to load dashboard data');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // Chart configurations
  const lineChartData = {
    labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
    datasets: [
      {
        label: 'Sales',
        data: [1200, 1900, 3000, 2500, 2800, 3200, 3500],
        fill: true,
        backgroundColor: 'rgba(37, 99, 235, 0.1)',
        borderColor: 'rgb(37, 99, 235)',
        tension: 0.4,
      },
    ],
  };

  const barChartData = {
    labels: ['Tomatoes', 'Onions', 'Potatoes', 'Carrots', 'Cabbage'],
    datasets: [
      {
        label: 'Quantity Sold',
        data: [65, 59, 80, 81, 56],
        backgroundColor: 'rgba(16, 185, 129, 0.8)',
      },
    ],
  };

  const doughnutChartData = {
    labels: ['Expired', 'Expiring Soon', 'Fresh'],
    datasets: [
      {
        data: [stats?.expiredItems || 0, stats?.expiringSoonItems || 0, 50],
        backgroundColor: [
          'rgba(239, 68, 68, 0.8)',
          'rgba(245, 158, 11, 0.8)',
          'rgba(16, 185, 129, 0.8)',
        ],
      },
    ],
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: true,
        position: 'bottom',
      },
    },
  };

  if (loading) {
    return (
      <div className="dashboard-loading">
        <div className="spinner-large"></div>
        <p>Loading dashboard...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="dashboard-error">
        <AlertTriangle size={48} />
        <p>{error}</p>
        <Button onClick={fetchData}>Retry</Button>
      </div>
    );
  }

  const kpiCards = [
    {
      title: 'Total Customers',
      value: stats?.totalCustomers || 0,
      icon: Users,
      color: 'primary',
      trend: '+12%',
      trendUp: true,
    },
    {
      title: 'Unreturned Crates',
      value: stats?.unreturnedCrates || 0,
      icon: Box,
      color: 'warning',
      trend: '-5%',
      trendUp: false,
    },
    {
      title: 'Outstanding Balance',
      value: `₹${stats?.outstandingBalance?.toFixed(0) || '0'}`,
      icon: DollarSign,
      color: 'danger',
      trend: '+8%',
      trendUp: true,
    },
    {
      title: 'Items Expiring Soon',
      value: stats?.expiringSoonItems || 0,
      icon: AlertTriangle,
      color: 'warning',
      alert: true,
    },
  ];

  const salesCards = [
    {
      title: "Today's Sales",
      value: `₹${stats?.todaysSales?.toFixed(2) || '0.00'}`,
      icon: Activity,
      color: 'success',
    },
    {
      title: 'Month Sales',
      value: `₹${stats?.monthSales?.toFixed(2) || '0.00'}`,
      icon: TrendingUp,
      color: 'info',
    },
    {
      title: 'Expired Items',
      value: stats?.expiredItems || 0,
      icon: Package,
      color: 'danger',
    },
  ];

  return (
    <div className="dashboard">
      {/* Header */}
      <div className="dashboard-header">
        <div>
          <h1>Dashboard</h1>
          <p className="dashboard-subtitle">Welcome back! Here's your business overview.</p>
        </div>
        <div className="dashboard-actions">
          <Link to="/transactions">
            <Button variant="primary" leftIcon={<Plus size={20} />}>
              New Sale
            </Button>
          </Link>
        </div>
      </div>

      {/* KPI Grid */}
      <div className="kpi-grid">
        {kpiCards.map((kpi, index) => {
          const Icon = kpi.icon;
          return (
            <Card key={index} className={`kpi-card kpi-${kpi.color}`} hoverable>
              <div className="kpi-content">
                <div className="kpi-icon-wrapper">
                  <Icon size={24} />
                </div>
                <div className="kpi-details">
                  <p className="kpi-title">{kpi.title}</p>
                  <h2 className="kpi-value">{kpi.value}</h2>
                  {kpi.trend && (
                    <div className="kpi-trend">
                      {kpi.trendUp ? (
                        <TrendingUp size={16} />
                      ) : (
                        <TrendingDown size={16} />
                      )}
                      <span>{kpi.trend}</span>
                    </div>
                  )}
                  {kpi.alert && <Badge variant="warning">Alert</Badge>}
                </div>
              </div>
            </Card>
          );
        })}
      </div>

      {/* Secondary Stats */}
      <div className="secondary-stats">
        {salesCards.map((stat, index) => {
          const Icon = stat.icon;
          return (
            <Card key={index} className="stat-card" padding="md">
              <div className="stat-header">
                <div className={`stat-icon stat-icon-${stat.color}`}>
                  <Icon size={20} />
                </div>
                <h3 className="stat-value">{stat.value}</h3>
              </div>
              <p className="stat-title">{stat.title}</p>
            </Card>
          );
        })}
      </div>

      {/* Charts Row */}
      <div className="charts-row">
        <Card title="Sales Trend" className="chart-card">
          <div className="chart-container">
            <Line data={lineChartData} options={chartOptions} />
          </div>
        </Card>

        <Card title="Top Products" className="chart-card">
          <div className="chart-container">
            <Bar data={barChartData} options={chartOptions} />
          </div>
        </Card>
      </div>

      {/* Bottom Section */}
      <div className="bottom-section">
        {/* Recent Activity */}
        <Card 
          title="Recent Activity" 
          className="activity-card"
          headerAction={
            <Link to="/transactions">
              <Button variant="ghost" size="sm">View All</Button>
            </Link>
          }
        >
          {activity && activity.length > 0 ? (
            <div className="activity-list">
              {activity.slice(0, 5).map((item, idx) => (
                <div key={idx} className="activity-item">
                  <div className="activity-indicator">
                    <div className={`activity-dot activity-dot-${item.type === 'sale' ? 'success' : 'info'}`} />
                  </div>
                  <div className="activity-content">
                    <p className="activity-description">{item.description}</p>
                    <span className="activity-time">
                      {item.timestamp && format(new Date(item.timestamp), 'MMM dd, HH:mm')}
                    </span>
                  </div>
                  <Badge variant={item.type === 'sale' ? 'success' : 'info'}>
                    {item.type}
                  </Badge>
                </div>
              ))}
            </div>
          ) : (
            <div className="empty-state">
              <Activity size={48} />
              <p>No recent activity</p>
            </div>
          )}
        </Card>

        {/* Inventory Status */}
        <Card title="Inventory Status" className="inventory-card">
          <div className="chart-container-small">
            <Doughnut data={doughnutChartData} options={chartOptions} />
          </div>
          <div className="inventory-legend">
            <div className="legend-item">
              <span className="legend-color legend-danger"></span>
              <span>Expired: {stats?.expiredItems || 0}</span>
            </div>
            <div className="legend-item">
              <span className="legend-color legend-warning"></span>
              <span>Expiring Soon: {stats?.expiringSoonItems || 0}</span>
            </div>
            <div className="legend-item">
              <span className="legend-color legend-success"></span>
              <span>Fresh: 50</span>
            </div>
          </div>
        </Card>
      </div>

      {/* Quick Actions */}
      <Card title="Quick Actions" className="quick-actions-card">
        <div className="quick-actions-grid">
          <Link to="/inventory">
            <Button variant="outline" fullWidth leftIcon={<Package size={20} />}>
              Add Inventory
            </Button>
          </Link>
          <Link to="/customers">
            <Button variant="outline" fullWidth leftIcon={<Users size={20} />}>
              Add Customer
            </Button>
          </Link>
          <Link to="/transactions">
            <Button variant="outline" fullWidth leftIcon={<DollarSign size={20} />}>
              New Transaction
            </Button>
          </Link>
          <Link to="/reports">
            <Button variant="outline" fullWidth leftIcon={<TrendingUp size={20} />}>
              Generate Report
            </Button>
          </Link>
        </div>
      </Card>
    </div>
  );
}

export default Dashboard;
