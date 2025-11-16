import React, { useState, useEffect } from 'react';
import { crateAPI, customerAPI } from '../services/api';

function Crates() {
  const [crates, setCrates] = useState([]);
  const [customers, setCustomers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
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
      await crateAPI.create({
        ...formData,
        cratesIssued: parseInt(formData.cratesIssued),
        cratesReturned: parseInt(formData.cratesReturned)
      });
      setShowModal(false);
      setFormData({ customerId: '', cratesIssued: 0, cratesReturned: 0, notes: '' });
      fetchData();
    } catch (err) {
      alert('Failed to add crate entry: ' + (err.response?.data || err.message));
    }
  };

  if (loading) return <div className="loading">Loading crate ledger...</div>;

  return (
    <div className="container">
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '20px' }}>
        <h1>Crate Management</h1>
        <button onClick={() => setShowModal(true)} className="btn btn-primary">Add Entry</button>
      </div>

      <div className="card">
        <h2>Customer Crate Balances</h2>
        <table className="table">
          <thead>
            <tr>
              <th>Customer</th>
              <th>Balance</th>
            </tr>
          </thead>
          <tbody>
            {customers.map(customer => {
              const customerCrates = crates.filter(c => c.customerId === customer.id);
              const balance = customerCrates.length > 0 ? customerCrates[0].balance : 0;
              return (
                <tr key={customer.id}>
                  <td>{customer.name}</td>
                  <td>
                    <span className={`badge ${balance > 0 ? 'badge-warning' : 'badge-success'}`}>
                      {balance} crates
                    </span>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      <div className="card">
        <h2>Recent Crate Transactions</h2>
        <table className="table">
          <thead>
            <tr>
              <th>Date</th>
              <th>Customer</th>
              <th>Issued</th>
              <th>Returned</th>
              <th>Balance</th>
              <th>Notes</th>
            </tr>
          </thead>
          <tbody>
            {crates.length === 0 ? (
              <tr><td colSpan="6" style={{ textAlign: 'center', padding: '20px' }}>No crate entries found</td></tr>
            ) : (
              crates.map(crate => {
                const customer = customers.find(c => c.id === crate.customerId);
                return (
                  <tr key={crate.id}>
                    <td>{new Date(crate.date).toLocaleDateString()}</td>
                    <td>{customer?.name || 'Unknown'}</td>
                    <td>{crate.cratesIssued}</td>
                    <td>{crate.cratesReturned}</td>
                    <td>{crate.balance}</td>
                    <td>{crate.notes || '-'}</td>
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
              <h2>Add Crate Entry</h2>
              <button className="close-btn" onClick={() => setShowModal(false)}>×</button>
            </div>
            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label>Customer*</label>
                <select value={formData.customerId} onChange={(e) => setFormData({...formData, customerId: e.target.value})} required>
                  <option value="">Select Customer</option>
                  {customers.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
                </select>
              </div>
              <div className="form-group">
                <label>Crates Issued</label>
                <input type="number" value={formData.cratesIssued} onChange={(e) => setFormData({...formData, cratesIssued: e.target.value})} />
              </div>
              <div className="form-group">
                <label>Crates Returned</label>
                <input type="number" value={formData.cratesReturned} onChange={(e) => setFormData({...formData, cratesReturned: e.target.value})} />
              </div>
              <div className="form-group">
                <label>Notes</label>
                <textarea value={formData.notes} onChange={(e) => setFormData({...formData, notes: e.target.value})} />
              </div>
              <div style={{ display: 'flex', gap: '10px', marginTop: '20px' }}>
                <button type="submit" className="btn btn-primary">Add Entry</button>
                <button type="button" onClick={() => setShowModal(false)} className="btn btn-secondary">Cancel</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

export default Crates;
