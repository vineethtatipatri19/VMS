import React, { useState, useEffect } from 'react';
import { crateAPI, customerAPI } from '../services/api';
import { Card, Button, Badge, Input, Select, Modal, useToast } from '../components/ui';
import { Box, Plus, Search, ArrowUp, ArrowDown, AlertCircle, Edit2, Trash2 } from 'lucide-react';
import DeleteConfirmationModal from '../components/DeleteConfirmationModal';

function Crates() {
  const toast = useToast();
  const [crates, setCrates] = useState([]);
  const [customers, setCustomers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingCrate, setEditingCrate] = useState(null);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deletingCrate, setDeletingCrate] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [formData, setFormData] = useState({
    customerId: '',
    cratesIssued: 0,
    cratesReturned: 0,
    notes: ''
  });

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [cratesRes, custRes] = await Promise.all([
        crateAPI.getAll(),
        customerAPI.getAll()
      ]);
      setCrates(cratesRes.data);
      setCustomers(custRes.data);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const payload = {
        ...formData,
        cratesIssued: parseInt(formData.cratesIssued),
        cratesReturned: parseInt(formData.cratesReturned)
      };

      if (editingCrate) {
        await crateAPI.update(editingCrate.id, payload);
        toast.success('Crate entry updated successfully');
      } else {
        await crateAPI.create(payload);
        toast.success('Crate entry added successfully');
      }
      
      handleCloseModal();
      fetchData();
    } catch (err) {
      toast.error('Failed to save crate entry: ' + (err.response?.data || err.message));
    }
  };

  const handleEdit = (crate) => {
    setEditingCrate(crate);
    setFormData({
      customerId: crate.customerId,
      cratesIssued: crate.cratesIssued || 0,
      cratesReturned: crate.cratesReturned || 0,
      notes: crate.notes || ''
    });
    setShowModal(true);
  };

  const handleDelete = (crate) => {
    setDeletingCrate(crate);
    setShowDeleteModal(true);
  };

  const confirmDelete = async (reason, attestation) => {
    try {
      await crateAPI.delete(deletingCrate.id, { reason, attestation });
      toast.success('Crate entry deleted successfully');
      setShowDeleteModal(false);
      setDeletingCrate(null);
      fetchData();
    } catch (err) {
      toast.error(err.response?.data || 'Failed to delete crate entry');
    }
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setEditingCrate(null);
    setFormData({ customerId: '', cratesIssued: 0, cratesReturned: 0, notes: '' });
  };

  const filteredCrates = crates.filter(crate => {
    const customer = customers.find(c => c.id === crate.customerId);
    const searchLower = searchTerm.toLowerCase();
    return searchTerm === '' ||
      (customer?.name && customer.name.toLowerCase().includes(searchLower));
  });

  const filteredCustomers = customers.filter(customer => {
    const searchLower = searchTerm.toLowerCase();
    return searchTerm === '' ||
      customer.name.toLowerCase().includes(searchLower);
  });

  if (loading) return <div className="loading">Loading crate ledger...</div>;

  return (
    <div className="container">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px', flexWrap: 'wrap', gap: '16px' }}>
        <div>
          <h1 style={{ marginBottom: '8px' }}>Crate Management</h1>
          <p style={{ color: 'var(--text-secondary)', fontSize: '14px' }}>Track crate issuance and returns by customer</p>
        </div>
        <Button variant="primary" leftIcon={<Plus size={20} />} onClick={() => setShowModal(true)}>
          Add Entry
        </Button>
      </div>

      <div style={{ marginBottom: '20px' }}>
        <Input
          placeholder="Search by customer name..."
          value={searchTerm}
          onChangeText={setSearchTerm}
          leftIcon={<Search size={18} />}
          fullWidth
        />
      </div>

      <Card style={{ marginBottom: '24px' }}>
        <h2 style={{ marginBottom: '20px', fontSize: '18px', fontWeight: '600' }}>Customer Crate Balances</h2>
        <div style={{ overflowX: 'auto' }}>
          <table className="table">
            <thead>
              <tr>
                <th>Customer</th>
                <th style={{ textAlign: 'right' }}>Balance</th>
              </tr>
            </thead>
            <tbody>
              {filteredCustomers.length === 0 ? (
                <tr>
                  <td colSpan="2" style={{ textAlign: 'center', padding: '40px' }}>
                    <Box size={48} style={{ opacity: 0.3, marginBottom: '16px' }} />
                    <p style={{ color: 'var(--text-secondary)' }}>
                      {searchTerm ? 'No customers match your search' : 'No customers found'}
                    </p>
                  </td>
                </tr>
              ) : (
                filteredCustomers.map(customer => {
                  const customerCrates = crates.filter(c => c.customerId === customer.id);
                  const balance = customerCrates.length > 0 ? customerCrates[0].balance : 0;
                  return (
                    <tr key={customer.id}>
                      <td style={{ fontWeight: '500' }}>{customer.name}</td>
                      <td style={{ textAlign: 'right' }}>
                        <Badge 
                          variant={balance > 0 ? 'warning' : 'success'} 
                          size="sm"
                          dot
                        >
                          <span style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                            {balance > 0 && <AlertCircle size={12} />}
                            {balance} crates
                          </span>
                        </Badge>
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>
      </Card>

      <Card>
        <h2 style={{ marginBottom: '20px', fontSize: '18px', fontWeight: '600' }}>Recent Crate Transactions</h2>
        <div style={{ overflowX: 'auto' }}>
          <table className="table">
            <thead>
              <tr>
                <th>Date</th>
                <th>Customer</th>
                <th style={{ textAlign: 'center' }}>Issued</th>
                <th style={{ textAlign: 'center' }}>Returned</th>
                <th style={{ textAlign: 'center' }}>Balance</th>
                <th>Notes</th>
                <th style={{ textAlign: 'right' }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {filteredCrates.length === 0 ? (
                <tr>
                  <td colSpan="7" style={{ textAlign: 'center', padding: '40px' }}>
                    <Box size={48} style={{ opacity: 0.3, marginBottom: '16px' }} />
                    <p style={{ color: 'var(--text-secondary)' }}>
                      {searchTerm ? 'No crate entries match your search' : 'No crate entries found'}
                    </p>
                  </td>
                </tr>
              ) : (
                filteredCrates.map(crate => {
                  const customer = customers.find(c => c.id === crate.customerId);
                  return (
                    <tr key={crate.id}>
                      <td style={{ color: 'var(--text-secondary)' }}>
                        {new Date(crate.date).toLocaleDateString()}
                      </td>
                      <td style={{ fontWeight: '500' }}>{customer?.name || 'Unknown'}</td>
                      <td style={{ textAlign: 'center' }}>
                        {crate.cratesIssued > 0 ? (
                          <span style={{ display: 'inline-flex', alignItems: 'center', gap: '4px', color: 'var(--color-warning)' }}>
                            <ArrowUp size={14} />
                            {crate.cratesIssued}
                          </span>
                        ) : '-'}
                      </td>
                      <td style={{ textAlign: 'center' }}>
                        {crate.cratesReturned > 0 ? (
                          <span style={{ display: 'inline-flex', alignItems: 'center', gap: '4px', color: 'var(--color-success)' }}>
                            <ArrowDown size={14} />
                            {crate.cratesReturned}
                          </span>
                        ) : '-'}
                      </td>
                      <td style={{ textAlign: 'center', fontWeight: '600' }}>
                        {crate.balance}
                      </td>
                      <td style={{ color: 'var(--text-secondary)', fontSize: '14px' }}>
                        {crate.notes || '-'}
                      </td>
                      <td>
                        <div style={{ display: 'flex', gap: '8px', justifyContent: 'flex-end' }}>
                          <Button 
                            variant="ghost" 
                            size="sm" 
                            leftIcon={<Edit2 size={16} />}
                            onClick={() => handleEdit(crate)}
                          >
                            Edit
                          </Button>
                          <Button 
                            variant="ghost" 
                            size="sm" 
                            leftIcon={<Trash2 size={16} />}
                            onClick={() => handleDelete(crate)}
                            style={{ color: 'var(--color-danger)' }}
                          >
                            Delete
                          </Button>
                        </div>
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>
      </Card>

      <Modal
        isOpen={showModal}
        onClose={handleCloseModal}
        title={editingCrate ? "Edit Crate Entry" : "Add Crate Entry"}
        size="md"
      >
        <form onSubmit={handleSubmit}>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
            {/* Customer Selection */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                Customer Information
              </h3>
              <Select
                label="Select Customer"
                value={formData.customerId}
                onChange={(e) => setFormData({...formData, customerId: e.target.value})}
                options={[
                  { value: '', label: 'Select Customer' },
                  ...customers.map(c => ({ value: c.id, label: c.name }))
                ]}
                required
                fullWidth
              />
            </div>

            {/* Crate Movement */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                Crate Movement
              </h3>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                <div style={{ padding: '16px', backgroundColor: 'var(--bg-secondary)', borderRadius: '8px', border: '1px solid var(--border-light)' }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '12px' }}>
                    <div style={{ padding: '8px', backgroundColor: 'var(--color-warning-light)', borderRadius: '6px' }}>
                      <ArrowUp size={18} style={{ color: 'var(--color-warning)' }} />
                    </div>
                    <span style={{ fontSize: '13px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase' }}>Issued</span>
                  </div>
                  <Input
                    type="number"
                    value={formData.cratesIssued}
                    onChangeText={(value) => setFormData({...formData, cratesIssued: value})}
                    placeholder="0"
                    fullWidth
                  />
                </div>
                <div style={{ padding: '16px', backgroundColor: 'var(--bg-secondary)', borderRadius: '8px', border: '1px solid var(--border-light)' }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '12px' }}>
                    <div style={{ padding: '8px', backgroundColor: 'var(--color-success-light)', borderRadius: '6px' }}>
                      <ArrowDown size={18} style={{ color: 'var(--color-success)' }} />
                    </div>
                    <span style={{ fontSize: '13px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase' }}>Returned</span>
                  </div>
                  <Input
                    type="number"
                    value={formData.cratesReturned}
                    onChangeText={(value) => setFormData({...formData, cratesReturned: value})}
                    placeholder="0"
                    fullWidth
                  />
                </div>
              </div>
            </div>

            {/* Additional Notes */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                Additional Notes
              </h3>
              <textarea 
                value={formData.notes} 
                onChange={(e) => setFormData({...formData, notes: e.target.value})}
                placeholder="Add any notes about this transaction..."
                rows="3"
                style={{
                  width: '100%',
                  padding: '12px 14px',
                  border: '1px solid var(--border-medium)',
                  borderRadius: '8px',
                  fontSize: '16px',
                  fontFamily: 'inherit',
                  resize: 'vertical',
                  color: 'var(--text-primary)',
                  backgroundColor: 'var(--bg-primary)',
                  transition: 'border-color 0.2s'
                }}
                onFocus={(e) => e.target.style.borderColor = 'var(--color-primary)'}
                onBlur={(e) => e.target.style.borderColor = 'var(--border-medium)'}
              />
            </div>

            {/* Action Buttons */}
            <div style={{ display: 'flex', gap: '12px', paddingTop: '20px', borderTop: '1px solid var(--border-light)' }}>
              <Button type="submit" variant="primary" fullWidth>
                {editingCrate ? 'Update Entry' : 'Add Entry'}
              </Button>
              <Button type="button" variant="secondary" onClick={handleCloseModal} fullWidth>
                Cancel
              </Button>
            </div>
          </div>
        </form>
      </Modal>

      <DeleteConfirmationModal
        isOpen={showDeleteModal}
        onClose={() => {
          setShowDeleteModal(false);
          setDeletingCrate(null);
        }}
        onConfirm={confirmDelete}
        title="Delete Crate Entry"
        entityType="Crate Entry"
        entityName={deletingCrate ? `Entry for ${customers.find(c => c.id === deletingCrate.customerId)?.name || 'Unknown'}` : ''}
        warningMessage="This will mark the crate entry as deleted. The entry record will be preserved for audit purposes."
      />
    </div>
  );
}

export default Crates;
