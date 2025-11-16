import React, { useState, useEffect } from 'react';
import { wastageAPI, inventoryAPI } from '../services/api';
import { Card, Button, Badge, Input, Select, Modal, useToast } from '../components/ui';
import { AlertTriangle, Plus, Search, Calendar, Edit2, Trash2 } from 'lucide-react';
import DeleteConfirmationModal from '../components/DeleteConfirmationModal';

function Wastage() {
  const toast = useToast();
  const [wastageLog, setWastageLog] = useState([]);
  const [inventory, setInventory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingWastage, setEditingWastage] = useState(null);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deletingWastage, setDeletingWastage] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [formData, setFormData] = useState({
    inventoryItemId: '',
    itemName: '',
    quantity: '',
    unit: 'kg',
    reason: 'expired',
    reasonDetails: '',
    costValue: '',
    loggedBy: '',
    photoUrl: ''
  });

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [wastageRes, invRes] = await Promise.all([
        wastageAPI.getAll(),
        inventoryAPI.getAll()
      ]);
      setWastageLog(wastageRes.data || []);
      setInventory(invRes.data || []);
    } catch (err) {
      console.error('Error fetching data:', err);
      toast.error('Failed to load wastage data');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!formData.itemName || !formData.quantity || !formData.reason) {
      toast.error('Please fill in all required fields');
      return;
    }

    try {
      const payload = {
        ...formData,
        quantity: parseFloat(formData.quantity),
        costValue: formData.costValue ? parseFloat(formData.costValue) : null
      };

      if (editingWastage) {
        await wastageAPI.update(editingWastage.id, payload);
        toast.success('Wastage entry updated successfully');
      } else {
        await wastageAPI.create(payload);
        toast.success('Wastage entry added successfully');
      }
      
      handleCloseModal();
      fetchData();
    } catch (err) {
      toast.error('Failed to save wastage entry: ' + (err.response?.data || err.message));
    }
  };

  const handleEdit = (wastage) => {
    setEditingWastage(wastage);
    setFormData({
      inventoryItemId: wastage.inventoryItemId || '',
      itemName: wastage.itemName,
      quantity: wastage.quantity.toString(),
      unit: wastage.unit,
      reason: wastage.reason,
      reasonDetails: wastage.reasonDetails || '',
      costValue: wastage.costValue?.toString() || '',
      loggedBy: wastage.loggedBy || '',
      photoUrl: wastage.photoUrl || ''
    });
    setShowModal(true);
  };

  const handleDelete = (wastage) => {
    setDeletingWastage(wastage);
    setShowDeleteModal(true);
  };

  const confirmDelete = async (reason, attestation) => {
    try {
      await wastageAPI.delete(deletingWastage.id, { reason, attestation });
      toast.success('Wastage entry deleted successfully');
      setShowDeleteModal(false);
      setDeletingWastage(null);
      fetchData();
    } catch (err) {
      toast.error(err.response?.data || 'Failed to delete wastage entry');
    }
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setEditingWastage(null);
    setFormData({
      inventoryItemId: '',
      itemName: '',
      quantity: '',
      unit: 'kg',
      reason: 'expired',
      reasonDetails: '',
      costValue: '',
      loggedBy: '',
      photoUrl: ''
    });
  };

  const handleInventorySelect = (itemId) => {
    const item = inventory.find(i => i.id === itemId);
    if (item) {
      setFormData(prev => ({
        ...prev,
        inventoryItemId: itemId,
        itemName: item.name,
        unit: item.unit,
        costValue: (item.costPrice || 0) * (prev.quantity || 0)
      }));
    }
  };

  const calculateTotalWastage = () => {
    return wastageLog.reduce((sum, log) => sum + (log.costValue || 0), 0);
  };

  const getReasonBadge = (reason) => {
    const variants = {
      expired: 'danger',
      damaged: 'warning',
      spoiled: 'warning',
      pest: 'danger',
      other: 'neutral'
    };
    return <Badge variant={variants[reason] || 'neutral'}>{reason.toUpperCase()}</Badge>;
  };

  const filteredWastage = wastageLog.filter(log => {
    const searchLower = searchTerm.toLowerCase();
    return searchTerm === '' ||
      log.itemName?.toLowerCase().includes(searchLower) ||
      log.reason?.toLowerCase().includes(searchLower);
  });

  if (loading) return <div className="loading">Loading wastage log...</div>;

  return (
    <div className="container">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px', flexWrap: 'wrap', gap: '16px' }}>
        <div>
          <h1 style={{ marginBottom: '8px' }}>Wastage Log</h1>
          <p style={{ color: 'var(--text-secondary)', fontSize: '14px' }}>Track damaged, expired, and spoiled inventory</p>
        </div>
        <Button variant="primary" leftIcon={<Plus size={20} />} onClick={() => setShowModal(true)}>
          Log Wastage
        </Button>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: '16px', marginBottom: '24px' }}>
        <Card>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <div style={{ padding: '12px', background: 'var(--danger-light)', borderRadius: '8px' }}>
              <AlertTriangle size={24} color="var(--danger)" />
            </div>
            <div>
              <p style={{ fontSize: '13px', color: 'var(--text-secondary)', marginBottom: '4px' }}>Total Wastage Value</p>
              <p style={{ fontSize: '24px', fontWeight: 'bold', color: 'var(--danger)' }}>₹{calculateTotalWastage().toFixed(2)}</p>
            </div>
          </div>
        </Card>
        <Card>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <div style={{ padding: '12px', background: 'var(--warning-light)', borderRadius: '8px' }}>
              <Calendar size={24} color="var(--warning)" />
            </div>
            <div>
              <p style={{ fontSize: '13px', color: 'var(--text-secondary)', marginBottom: '4px' }}>Total Entries</p>
              <p style={{ fontSize: '24px', fontWeight: 'bold' }}>{wastageLog.length}</p>
            </div>
          </div>
        </Card>
      </div>

      <div style={{ marginBottom: '20px' }}>
        <Input
          placeholder="Search by item name or reason..."
          value={searchTerm}
          onChangeText={setSearchTerm}
          leftIcon={<Search size={18} />}
          fullWidth
        />
      </div>

      <Card>
        <h2 style={{ marginBottom: '20px', fontSize: '18px', fontWeight: '600' }}>Wastage Entries</h2>
        <div style={{ overflowX: 'auto' }}>
          <table className="table">
            <thead>
              <tr>
                <th>Date</th>
                <th>Item Name</th>
                <th>Quantity</th>
                <th>Reason</th>
                <th>Cost Value</th>
                <th>Logged By</th>
                <th>Notes</th>
                <th style={{ textAlign: 'right' }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {filteredWastage.length === 0 ? (
                <tr>
                  <td colSpan="8" style={{ textAlign: 'center', padding: '40px', color: 'var(--text-secondary)' }}>
                    <AlertTriangle size={48} style={{ opacity: 0.3, marginBottom: '16px' }} />
                    <p>No wastage entries found</p>
                    <p style={{ fontSize: '14px', marginTop: '8px' }}>Click "Log Wastage" to record damaged or expired items</p>
                  </td>
                </tr>
              ) : (
                filteredWastage.map((log) => (
                  <tr key={log.id}>
                    <td>{new Date(log.loggedAt).toLocaleDateString()}</td>
                    <td style={{ fontWeight: '500' }}>{log.itemName}</td>
                    <td>{log.quantity} {log.unit}</td>
                    <td>{getReasonBadge(log.reason)}</td>
                    <td style={{ color: 'var(--danger)', fontWeight: '500' }}>₹{(log.costValue || 0).toFixed(2)}</td>
                    <td>{log.loggedBy || '—'}</td>
                    <td style={{ maxWidth: '200px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                      {log.reasonDetails || '—'}
                    </td>
                    <td>
                      <div style={{ display: 'flex', gap: '8px', justifyContent: 'flex-end' }}>
                        <Button 
                          variant="ghost" 
                          size="sm" 
                          leftIcon={<Edit2 size={16} />}
                          onClick={() => handleEdit(log)}
                        >
                          Edit
                        </Button>
                        <Button 
                          variant="ghost" 
                          size="sm" 
                          leftIcon={<Trash2 size={16} />}
                          onClick={() => handleDelete(log)}
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

      {showModal && (
        <Modal onClose={handleCloseModal} title={editingWastage ? "Edit Wastage Entry" : "Log Wastage Entry"}>
          <form onSubmit={handleSubmit}>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
              <Select
                label="Select Inventory Item (Optional)"
                value={formData.inventoryItemId}
                onChange={(e) => handleInventorySelect(e.target.value)}
              >
                <option value="">Select an item...</option>
                {inventory.map(item => (
                  <option key={item.id} value={item.id}>
                    {item.name} - {item.quantity} {item.unit}
                  </option>
                ))}
              </Select>

              <Input
                label="Item Name *"
                value={formData.itemName}
                onChangeText={(val) => setFormData({ ...formData, itemName: val })}
                placeholder="Enter item name"
                required
              />

              <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '16px' }}>
                <Input
                  label="Quantity *"
                  type="number"
                  step="0.01"
                  value={formData.quantity}
                  onChangeText={(val) => {
                    setFormData({ 
                      ...formData, 
                      quantity: val,
                      costValue: inventory.find(i => i.id === formData.inventoryItemId)?.costPrice * parseFloat(val || 0) || formData.costValue
                    });
                  }}
                  placeholder="0.00"
                  required
                />

                <Select
                  label="Unit *"
                  value={formData.unit}
                  onChange={(e) => setFormData({ ...formData, unit: e.target.value })}
                >
                  <option value="kg">kg</option>
                  <option value="lot">lot</option>
                </Select>
              </div>

              <Select
                label="Reason *"
                value={formData.reason}
                onChange={(e) => setFormData({ ...formData, reason: e.target.value })}
              >
                <option value="expired">Expired</option>
                <option value="damaged">Damaged</option>
                <option value="spoiled">Spoiled</option>
                <option value="pest">Pest Damage</option>
                <option value="other">Other</option>
              </Select>

              <Input
                label="Reason Details"
                value={formData.reasonDetails}
                onChangeText={(val) => setFormData({ ...formData, reasonDetails: val })}
                placeholder="Additional details..."
              />

              <Input
                label="Cost Value"
                type="number"
                step="0.01"
                value={formData.costValue}
                onChangeText={(val) => setFormData({ ...formData, costValue: val })}
                placeholder="₹0.00"
              />

              <Input
                label="Logged By"
                value={formData.loggedBy}
                onChangeText={(val) => setFormData({ ...formData, loggedBy: val })}
                placeholder="Name of person logging"
              />

              <div style={{ display: 'flex', gap: '12px', marginTop: '8px' }}>
                <Button type="button" variant="secondary" onClick={handleCloseModal} style={{ flex: 1 }}>
                  Cancel
                </Button>
                <Button type="submit" variant="primary" style={{ flex: 1 }}>
                  {editingWastage ? 'Update Entry' : 'Log Wastage'}
                </Button>
              </div>
            </div>
          </form>
        </Modal>
      )}

      <DeleteConfirmationModal
        isOpen={showDeleteModal}
        onClose={() => {
          setShowDeleteModal(false);
          setDeletingWastage(null);
        }}
        onConfirm={confirmDelete}
        title="Delete Wastage Entry"
        entityType="Wastage Entry"
        entityName={deletingWastage ? `${deletingWastage.itemName} - ${deletingWastage.quantity} ${deletingWastage.unit}` : ''}
        warningMessage="This will mark the wastage entry as deleted. The entry record will be preserved for audit purposes."
      />
    </div>
  );
}

export default Wastage;
