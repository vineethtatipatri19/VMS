import React, { useState, useEffect } from 'react';
import { expiryAlertsAPI, inventoryAPI } from '../services/api';
import { Card, Button, Badge, Input, useToast } from '../components/ui';
import { AlertCircle, Bell, CheckCircle, Search, Calendar, Trash2 } from 'lucide-react';
import DeleteConfirmationModal from '../components/DeleteConfirmationModal';
import { normalizeExpiryAlert, normalizeInventory, formatDate } from '../utils/dataHelpers';

function ExpiryAlerts() {
  const toast = useToast();
  const [alerts, setAlerts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deletingAlert, setDeletingAlert] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [filter, setFilter] = useState('all'); // all, unacknowledged, acknowledged

  useEffect(() => {
    fetchAlerts();
  }, []);

  const fetchAlerts = async () => {
    setLoading(true);
    try {
      const [alertsRes, inventoryRes] = await Promise.all([
        expiryAlertsAPI.getAll(),
        inventoryAPI.getAll()
      ]);
      const alertsData = alertsRes.data.data || alertsRes.data || [];
      const inventoryData = inventoryRes.data.data || inventoryRes.data || [];
      
      // Map inventory items by ID for quick lookup
      const inventoryMap = {};
      inventoryData.forEach(item => {
        const normalizedItem = normalizeInventory(item);
        inventoryMap[normalizedItem.id] = normalizedItem;
      });
      
      // Add item names to alerts
      const alertsWithItems = alertsData.map(alert => {
        const normalized = normalizeExpiryAlert(alert);
        const inventoryItem = inventoryMap[normalized.inventoryItemId];
        return {
          ...normalized,
          itemName: inventoryItem ? inventoryItem.name : 'Unknown Item'
        };
      });
      
      setAlerts(alertsWithItems);
    } catch (err) {
      console.error('Error fetching alerts:', err);
      toast.error('Failed to load expiry alerts');
    } finally {
      setLoading(false);
    }
  };

  const handleAcknowledge = async (alertId) => {
    try {
      await expiryAlertsAPI.acknowledge(alertId);
      toast.success('Alert acknowledged');
      fetchAlerts();
    } catch (err) {
      toast.error('Failed to acknowledge alert: ' + (err.response?.data || err.message));
    }
  };

  const handleDelete = (alert) => {
    setDeletingAlert(alert);
    setShowDeleteModal(true);
  };

  const confirmDelete = async (reason, attestation) => {
    try {
      await expiryAlertsAPI.delete(deletingAlert.id, { reason, attestation });
      toast.success('Expiry alert deleted successfully');
      setShowDeleteModal(false);
      setDeletingAlert(null);
      fetchAlerts();
    } catch (err) {
      toast.error(err.response?.data || 'Failed to delete alert');
    }
  };

  const getDaysUntilExpiryColor = (days) => {
    if (days < 0) return 'var(--danger)';
    if (days <= 2) return 'var(--danger)';
    if (days <= 5) return 'var(--warning)';
    return 'var(--text-primary)';
  };

  const getDaysUntilExpiryText = (days) => {
    if (days < 0) return 'Expired';
    if (days === 0) return 'Expires Today';
    if (days === 1) return '1 Day';
    return `${days} Days`;
  };

  const getUrgencyBadge = (days) => {
    if (days < 0) return <Badge variant="danger">EXPIRED</Badge>;
    if (days <= 2) return <Badge variant="danger">CRITICAL</Badge>;
    if (days <= 5) return <Badge variant="warning">URGENT</Badge>;
    return <Badge variant="info">UPCOMING</Badge>;
  };

  const filteredAlerts = alerts.filter(alert => {
    const searchLower = searchTerm.toLowerCase();
    const matchesSearch = searchTerm === '' ||
      alert.itemName?.toLowerCase().includes(searchLower);
    
    const matchesFilter = filter === 'all' ||
      (filter === 'unacknowledged' && !alert.acknowledged) ||
      (filter === 'acknowledged' && alert.acknowledged);
    
    return matchesSearch && matchesFilter;
  });

  const unacknowledgedCount = alerts.filter(a => !a.acknowledged).length;
  const criticalCount = alerts.filter(a => a.daysUntilExpiry <= 2 && !a.acknowledged).length;

  if (loading) return <div className="loading">Loading expiry alerts...</div>;

  return (
    <div className="container">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px', flexWrap: 'wrap', gap: '16px' }}>
        <div>
          <h1 style={{ marginBottom: '8px' }}>Expiry Alerts</h1>
          <p style={{ color: 'var(--text-secondary)', fontSize: '14px' }}>Monitor items approaching expiration</p>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: '16px', marginBottom: '24px' }}>
        <Card>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <div style={{ padding: '12px', background: 'var(--danger-light)', borderRadius: '8px' }}>
              <AlertCircle size={24} color="var(--danger)" />
            </div>
            <div>
              <p style={{ fontSize: '13px', color: 'var(--text-secondary)', marginBottom: '4px' }}>Critical Alerts</p>
              <p style={{ fontSize: '24px', fontWeight: 'bold', color: 'var(--danger)' }}>{criticalCount}</p>
            </div>
          </div>
        </Card>
        <Card>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <div style={{ padding: '12px', background: 'var(--warning-light)', borderRadius: '8px' }}>
              <Bell size={24} color="var(--warning)" />
            </div>
            <div>
              <p style={{ fontSize: '13px', color: 'var(--text-secondary)', marginBottom: '4px' }}>Pending Alerts</p>
              <p style={{ fontSize: '24px', fontWeight: 'bold' }}>{unacknowledgedCount}</p>
            </div>
          </div>
        </Card>
        <Card>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <div style={{ padding: '12px', background: 'var(--primary-light)', borderRadius: '8px' }}>
              <Calendar size={24} color="var(--primary)" />
            </div>
            <div>
              <p style={{ fontSize: '13px', color: 'var(--text-secondary)', marginBottom: '4px' }}>Total Alerts</p>
              <p style={{ fontSize: '24px', fontWeight: 'bold' }}>{alerts.length}</p>
            </div>
          </div>
        </Card>
      </div>

      <div style={{ display: 'flex', gap: '12px', marginBottom: '20px', flexWrap: 'wrap' }}>
        <Input
          placeholder="Search by item name..."
          value={searchTerm}
          onChangeText={setSearchTerm}
          leftIcon={<Search size={18} />}
          style={{ flex: '1', minWidth: '250px' }}
        />
        <div style={{ display: 'flex', gap: '8px' }}>
          <Button 
            variant={filter === 'all' ? 'primary' : 'secondary'}
            onClick={() => setFilter('all')}
          >
            All
          </Button>
          <Button 
            variant={filter === 'unacknowledged' ? 'primary' : 'secondary'}
            onClick={() => setFilter('unacknowledged')}
          >
            Pending
          </Button>
          <Button 
            variant={filter === 'acknowledged' ? 'primary' : 'secondary'}
            onClick={() => setFilter('acknowledged')}
          >
            Acknowledged
          </Button>
        </div>
      </div>

      <Card>
        <h2 style={{ marginBottom: '20px', fontSize: '18px', fontWeight: '600' }}>Expiry Alerts</h2>
        <div style={{ overflowX: 'auto' }}>
          <table className="table">
            <thead>
              <tr>
                <th>Urgency</th>
                <th>Item Name</th>
                <th>Expiry Date</th>
                <th>Days Until Expiry</th>
                <th>Alert Date</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {filteredAlerts.length === 0 ? (
                <tr>
                  <td colSpan="7" style={{ textAlign: 'center', padding: '40px', color: 'var(--text-secondary)' }}>
                    <CheckCircle size={48} style={{ opacity: 0.3, marginBottom: '16px' }} />
                    <p>No expiry alerts found</p>
                    <p style={{ fontSize: '14px', marginTop: '8px' }}>All items are monitored for upcoming expiration</p>
                  </td>
                </tr>
              ) : (
                filteredAlerts.map((alert) => (
                  <tr key={alert.id} style={{ opacity: alert.acknowledged ? 0.6 : 1 }}>
                    <td>{getUrgencyBadge(alert.daysUntilExpiry)}</td>
                    <td style={{ fontWeight: '500' }}>{alert.itemName || '—'}</td>
                    <td>{new Date(alert.expiryDate).toLocaleDateString()}</td>
                    <td style={{ color: getDaysUntilExpiryColor(alert.daysUntilExpiry), fontWeight: '500' }}>
                      {getDaysUntilExpiryText(alert.daysUntilExpiry)}
                    </td>
                    <td>{new Date(alert.alertDate).toLocaleDateString()}</td>
                    <td>
                      {alert.acknowledged ? (
                        <Badge variant="success">Acknowledged</Badge>
                      ) : (
                        <Badge variant="warning">Pending</Badge>
                      )}
                    </td>
                    <td>
                      <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
                        {!alert.acknowledged && (
                          <Button 
                            variant="primary" 
                            size="small"
                            onClick={() => handleAcknowledge(alert.id)}
                          >
                            Acknowledge
                          </Button>
                        )}
                        {alert.acknowledged && alert.acknowledgedBy && (
                          <span style={{ fontSize: '13px', color: 'var(--text-secondary)' }}>
                            by {alert.acknowledgedBy}
                          </span>
                        )}
                        <Button 
                          variant="ghost" 
                          size="small"
                          leftIcon={<Trash2 size={14} />}
                          onClick={() => handleDelete(alert)}
                          style={{ color: 'var(--color-danger)' }}
                        >
                          Delete
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </Card>

      <DeleteConfirmationModal
        isOpen={showDeleteModal}
        onClose={() => {
          setShowDeleteModal(false);
          setDeletingAlert(null);
        }}
        onConfirm={confirmDelete}
        title="Delete Expiry Alert"
        entityType="Expiry Alert"
        entityName={deletingAlert ? `${deletingAlert.itemName} - Expires ${new Date(deletingAlert.expiryDate).toLocaleDateString()}` : ''}
        warningMessage="This will mark the expiry alert as deleted. The alert record will be preserved for audit purposes."
      />
    </div>
  );
}

export default ExpiryAlerts;
