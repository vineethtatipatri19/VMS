import React, { useState, useEffect } from 'react';
import { transactionAPI, customerAPI, inventoryAPI } from '../services/api';

function Transactions() {
  const [transactions, setTransactions] = useState([]);
  const [customers, setCustomers] = useState([]);
  const [inventory, setInventory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [filter, setFilter] = useState('');
  const [formData, setFormData] = useState({
    customerId: '',
    type: 'sale',
    paymentAmount: '',
    items: [{ inventoryLotId: '', itemName: '', quantity: '', pricePerUnit: '', unit: 'kg' }]
  });

  useEffect(() => {
    fetchData();
  }, [filter]);

  const fetchData = async () => {
    setLoading(true);
    try {
      const params = filter ? { type: filter } : {};
      const [txRes, custRes, invRes] = await Promise.all([
        transactionAPI.getAll(params),
        customerAPI.getAll(),
        inventoryAPI.getAll({ sort: 'expiry' })
      ]);
      setTransactions(txRes.data);
      setCustomers(custRes.data);
      setInventory(invRes.data);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const payload = { ...formData };
      if (formData.type === 'payment') {
        payload.paymentAmount = parseFloat(formData.paymentAmount);
        delete payload.items;
      } else {
        payload.items = formData.items.map(item => ({
          ...item,
          quantity: parseFloat(item.quantity),
          pricePerUnit: parseFloat(item.pricePerUnit)
        }));
      }
      await transactionAPI.create(payload);
      setShowModal(false);
      setFormData({
        customerId: '',
        type: 'sale',
        paymentAmount: '',
        items: [{ inventoryLotId: '', itemName: '', quantity: '', pricePerUnit: '', unit: 'kg' }]
      });
      fetchData();
    } catch (err) {
      alert('Failed to create transaction: ' + (err.response?.data || err.message));
    }
  };

  const addItem = () => {
    setFormData({
      ...formData,
      items: [...formData.items, { inventoryLotId: '', itemName: '', quantity: '', pricePerUnit: '', unit: 'kg' }]
    });
  };

  const updateItem = (index, field, value) => {
    const newItems = [...formData.items];
    newItems[index][field] = value;
    if (field === 'inventoryLotId') {
      const invItem = inventory.find(i => i.id === value);
      if (invItem) {
        newItems[index].itemName = invItem.name;
        newItems[index].unit = invItem.unit;
      }
    }
    setFormData({ ...formData, items: newItems });
  };

  if (loading) return <div className="loading">Loading transactions...</div>;

  return (
    <div className="container">
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '20px' }}>
        <h1>Transaction Ledger</h1>
        <button onClick={() => setShowModal(true)} className="btn btn-success">New Transaction</button>
      </div>

      <div className="card">
        <div style={{ marginBottom: '20px' }}>
          <label style={{ marginRight: '10px' }}>Filter:</label>
          <select value={filter} onChange={(e) => setFilter(e.target.value)} style={{ padding: '8px', borderRadius: '4px', border: '1px solid #ddd' }}>
            <option value="">All</option>
            <option value="sale">Sales</option>
            <option value="payment">Payments</option>
          </select>
        </div>

        <table className="table">
          <thead>
            <tr>
              <th>Date</th>
              <th>Type</th>
              <th>Customer</th>
              <th>Amount</th>
            </tr>
          </thead>
          <tbody>
            {transactions.length === 0 ? (
              <tr><td colSpan="4" style={{ textAlign: 'center', padding: '20px' }}>No transactions found</td></tr>
            ) : (
              transactions.map((tx) => {
                const customer = customers.find(c => c.id === tx.customerId);
                return (
                  <tr key={tx.id}>
                    <td>{new Date(tx.date).toLocaleDateString()}</td>
                    <td>
                      <span className={`badge ${tx.type === 'sale' ? 'badge-success' : 'badge-info'}`}>
                        {tx.type}
                      </span>
                    </td>
                    <td>{customer?.name || 'Unknown'}</td>
                    <td>₹{tx.totalAmount?.toFixed(2) || '0.00'}</td>
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
              <h2>New Transaction</h2>
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
                <label>Type*</label>
                <select value={formData.type} onChange={(e) => setFormData({...formData, type: e.target.value})} required>
                  <option value="sale">Sale</option>
                  <option value="payment">Payment</option>
                </select>
              </div>
              {formData.type === 'payment' ? (
                <div className="form-group">
                  <label>Payment Amount*</label>
                  <input type="number" step="0.01" value={formData.paymentAmount} onChange={(e) => setFormData({...formData, paymentAmount: e.target.value})} required />
                </div>
              ) : (
                <>
                  <h3>Sale Items</h3>
                  {formData.items.map((item, idx) => (
                    <div key={idx} style={{ border: '1px solid #ddd', padding: '10px', marginBottom: '10px', borderRadius: '4px' }}>
                      <div className="form-group">
                        <label>Inventory Item*</label>
                        <select value={item.inventoryLotId} onChange={(e) => updateItem(idx, 'inventoryLotId', e.target.value)} required>
                          <option value="">Select Item</option>
                          {inventory.map(inv => <option key={inv.id} value={inv.id}>{inv.name} ({inv.quantity} {inv.unit})</option>)}
                        </select>
                      </div>
                      <div className="form-group">
                        <label>Quantity*</label>
                        <input type="number" step="0.01" value={item.quantity} onChange={(e) => updateItem(idx, 'quantity', e.target.value)} required />
                      </div>
                      <div className="form-group">
                        <label>Price per Unit*</label>
                        <input type="number" step="0.01" value={item.pricePerUnit} onChange={(e) => updateItem(idx, 'pricePerUnit', e.target.value)} required />
                      </div>
                    </div>
                  ))}
                  <button type="button" onClick={addItem} className="btn btn-secondary" style={{ marginBottom: '15px' }}>+ Add Item</button>
                </>
              )}
              <div style={{ display: 'flex', gap: '10px', marginTop: '20px' }}>
                <button type="submit" className="btn btn-primary">Create Transaction</button>
                <button type="button" onClick={() => setShowModal(false)} className="btn btn-secondary">Cancel</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

export default Transactions;
