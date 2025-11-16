import React, { useState, useEffect } from 'react';
import { inventoryAPI } from '../services/api';

function Inventory() {
  const [inventory, setInventory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [filter, setFilter] = useState('');
  const [formData, setFormData] = useState({
    name: '',
    variant: '',
    quantity: '',
    unit: 'kg',
    purchaseDate: new Date().toISOString().split('T')[0],
    expiryDate: ''
  });

  useEffect(() => {
    fetchInventory();
  }, [filter]);

  const fetchInventory = async () => {
    setLoading(true);
    try {
      const params = filter ? { status: filter } : { sort: 'expiry' };
      const res = await inventoryAPI.getAll(params);
      setInventory(res.data);
    } catch (err) {
      setError('Failed to load inventory');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await inventoryAPI.create(formData);
      setShowModal(false);
      setFormData({
        name: '',
        variant: '',
        quantity: '',
        unit: 'kg',
        purchaseDate: new Date().toISOString().split('T')[0],
        expiryDate: ''
      });
      fetchInventory();
    } catch (err) {
      alert('Failed to add inventory: ' + (err.response?.data || err.message));
    }
  };

  const handleDelete = async (id) => {
    if (window.confirm('Are you sure you want to delete this item?')) {
      try {
        await inventoryAPI.delete(id);
        fetchInventory();
      } catch (err) {
        alert('Failed to delete item');
      }
    }
  };

  const getStatus = (expiryDate) => {
    const today = new Date();
    const expiry = new Date(expiryDate);
    const daysUntil = Math.ceil((expiry - today) / (1000 * 60 * 60 * 24));
    
    if (daysUntil < 0) return { label: 'Expired', class: 'badge-danger' };
    if (daysUntil <= 3) return { label: 'Expiring Soon', class: 'badge-warning' };
    return { label: 'Fresh', class: 'badge-success' };
  };

  if (loading) return <div className="loading">Loading inventory...</div>;

  return (
    <div className="container">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <h1>Inventory Management</h1>
        <button onClick={() => setShowModal(true)} className="btn btn-primary">
          Add Inventory
        </button>
      </div>

      {error && <div className="error">{error}</div>}

      <div className="card">
        <div style={{ marginBottom: '20px' }}>
          <label style={{ marginRight: '10px' }}>Filter by status:</label>
          <select value={filter} onChange={(e) => setFilter(e.target.value)} style={{ padding: '8px', borderRadius: '4px', border: '1px solid #ddd' }}>
            <option value="">All (FEFO Sorted)</option>
            <option value="fresh">Fresh</option>
            <option value="expiring_soon">Expiring Soon</option>
            <option value="expired">Expired</option>
          </select>
        </div>

        <table className="table">
          <thead>
            <tr>
              <th>Status</th>
              <th>Name</th>
              <th>Variant</th>
              <th>Lot Number</th>
              <th>Quantity</th>
              <th>Purchase Date</th>
              <th>Expiry Date</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {inventory.length === 0 ? (
              <tr><td colSpan="8" style={{ textAlign: 'center', padding: '20px' }}>No inventory items found</td></tr>
            ) : (
              inventory.map((item) => {
                const status = getStatus(item.expiryDate);
                return (
                  <tr key={item.id}>
                    <td><span className={`badge ${status.class}`}>{status.label}</span></td>
                    <td>{item.name}</td>
                    <td>{item.variant || '-'}</td>
                    <td>{item.lotNumber}</td>
                    <td>{item.quantity} {item.unit}</td>
                    <td>{item.purchaseDate}</td>
                    <td>{item.expiryDate}</td>
                    <td>
                      <button onClick={() => handleDelete(item.id)} className="btn btn-danger" style={{ fontSize: '12px', padding: '5px 10px' }}>
                        Delete
                      </button>
                    </td>
                  </tr>
                );
              })
            )}
          </tbody>
        </table>
      </div>

      {showModal && (
        <div className="modal">
          <div className="modal-content">
            <div className="modal-header">
              <h2>Add Inventory Item</h2>
              <button className="close-btn" onClick={() => setShowModal(false)}>×</button>
            </div>
            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label>Item Name*</label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({...formData, name: e.target.value})}
                  required
                />
              </div>
              <div className="form-group">
                <label>Variant</label>
                <input
                  type="text"
                  value={formData.variant}
                  onChange={(e) => setFormData({...formData, variant: e.target.value})}
                  placeholder="e.g., Grade A, Organic"
                />
              </div>
              <div className="form-group">
                <label>Quantity*</label>
                <input
                  type="number"
                  step="0.01"
                  value={formData.quantity}
                  onChange={(e) => setFormData({...formData, quantity: e.target.value})}
                  required
                />
              </div>
              <div className="form-group">
                <label>Unit*</label>
                <select
                  value={formData.unit}
                  onChange={(e) => setFormData({...formData, unit: e.target.value})}
                  required
                >
                  <option value="kg">kg</option>
                  <option value="lot">lot</option>
                </select>
              </div>
              <div className="form-group">
                <label>Purchase Date*</label>
                <input
                  type="date"
                  value={formData.purchaseDate}
                  onChange={(e) => setFormData({...formData, purchaseDate: e.target.value})}
                  required
                />
              </div>
              <div className="form-group">
                <label>Expiry Date*</label>
                <input
                  type="date"
                  value={formData.expiryDate}
                  onChange={(e) => setFormData({...formData, expiryDate: e.target.value})}
                  required
                />
              </div>
              <div style={{ display: 'flex', gap: '10px', marginTop: '20px' }}>
                <button type="submit" className="btn btn-primary">Add Item</button>
                <button type="button" onClick={() => setShowModal(false)} className="btn btn-secondary">Cancel</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

export default Inventory;
