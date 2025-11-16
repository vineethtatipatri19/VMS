import React, { useState, useEffect } from 'react';
import { transactionAPI, customerAPI, inventoryAPI } from '../services/api';
import { Card, Button, Badge, Input, Select, Modal, useToast } from '../components/ui';
import { Receipt, Plus, Search, DollarSign, ShoppingCart, Trash2, Edit2 } from 'lucide-react';
import DeleteConfirmationModal from '../components/DeleteConfirmationModal';

function Transactions() {
  const toast = useToast();
  const [transactions, setTransactions] = useState([]);
  const [customers, setCustomers] = useState([]);
  const [inventory, setInventory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [editingTransaction, setEditingTransaction] = useState(null);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deletingTransaction, setDeletingTransaction] = useState(null);
  const [filter, setFilter] = useState('');
  const [searchTerm, setSearchTerm] = useState('');
  const [formData, setFormData] = useState({
    customerId: '',
    type: 'sale',
    paymentAmount: '',
    paymentMethod: 'cash',
    paymentReference: '',
    dueDate: '',
    discountAmount: '',
    discountPercentage: '',
    taxAmount: '',
    taxPercentage: '',
    saleType: 'cash',
    deliveryStatus: 'pending',
    deliveryDate: '',
    deliveryAddress: '',
    notes: '',
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
    
    // Validate customer selection
    if (!formData.customerId) {
      alert('Please select a customer');
      return;
    }
    
    try {
      const payload = {
        customerId: formData.customerId,
        type: formData.type
      };
      
      if (formData.type === 'payment') {
        const amount = parseFloat(formData.paymentAmount);
        if (isNaN(amount) || amount <= 0) {
          toast.error('Please enter a valid payment amount greater than 0');
          return;
        }
        payload.paymentAmount = amount;
      } else {
        // Validate sale items
        const items = formData.items.map(item => {
          const qty = parseFloat(item.quantity);
          const price = parseFloat(item.pricePerUnit);
          
          if (!item.inventoryLotId) {
            throw new Error('Please select an inventory item');
          }
          if (isNaN(qty) || qty <= 0) {
            throw new Error('Please enter valid quantities greater than 0');
          }
          if (isNaN(price) || price <= 0) {
            throw new Error('Please enter valid prices greater than 0');
          }
          
          return {
            ...item,
            quantity: qty,
            pricePerUnit: price
          };
        });
        payload.items = items;
      }
      
      await transactionAPI.create(payload);
      toast.success('Transaction created successfully');
      handleCloseModal();
      fetchData();
    } catch (err) {
      toast.error('Failed to create transaction: ' + (err.response?.data || err.message));
    }
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setFormData({
      customerId: '',
      type: 'sale',
      paymentAmount: '',
      paymentMethod: 'cash',
      paymentReference: '',
      dueDate: '',
      discountAmount: '',
      discountPercentage: '',
      taxAmount: '',
      taxPercentage: '',
      saleType: 'cash',
      deliveryStatus: 'pending',
      deliveryDate: '',
      deliveryAddress: '',
      notes: '',
      items: [{ inventoryLotId: '', itemName: '', quantity: '', pricePerUnit: '', unit: 'kg' }]
    });
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

  const removeItem = (index) => {
    const newItems = formData.items.filter((_, i) => i !== index);
    if (newItems.length === 0) {
      newItems.push({ inventoryLotId: '', itemName: '', quantity: '', pricePerUnit: '', unit: 'kg' });
    }
    setFormData({ ...formData, items: newItems });
  };

  const filteredTransactions = transactions.filter(tx => {
    const customer = customers.find(c => c.id === tx.customerId);
    const searchLower = searchTerm.toLowerCase();
    return searchTerm === '' ||
      (customer?.name && customer.name.toLowerCase().includes(searchLower)) ||
      tx.type.toLowerCase().includes(searchLower);
  });

  const handleEditClick = (transaction) => {
    setEditingTransaction({
      ...transaction,
      paymentAmount: transaction.paymentAmount || '',
      discountAmount: transaction.discountAmount || '',
      taxAmount: transaction.taxAmount || '',
      paymentMethod: transaction.paymentMethod || 'cash',
      paymentReference: transaction.paymentReference || '',
      dueDate: transaction.dueDate || '',
      notes: transaction.notes || '',
      invoiceNumber: transaction.invoiceNumber || '',
      saleType: transaction.saleType || 'cash',
      deliveryStatus: transaction.deliveryStatus || 'pending',
      deliveryDate: transaction.deliveryDate || '',
      deliveryAddress: transaction.deliveryAddress || ''
    });
    setShowEditModal(true);
  };

  const handleEditSubmit = async (e) => {
    e.preventDefault();
    try {
      const payload = {
        paymentAmount: editingTransaction.paymentAmount ? parseFloat(editingTransaction.paymentAmount) : undefined,
        totalAmount: editingTransaction.totalAmount ? parseFloat(editingTransaction.totalAmount) : undefined,
        discountAmount: editingTransaction.discountAmount ? parseFloat(editingTransaction.discountAmount) : undefined,
        taxAmount: editingTransaction.taxAmount ? parseFloat(editingTransaction.taxAmount) : undefined,
        paymentMethod: editingTransaction.paymentMethod,
        paymentReference: editingTransaction.paymentReference,
        dueDate: editingTransaction.dueDate,
        notes: editingTransaction.notes,
        invoiceNumber: editingTransaction.invoiceNumber,
        saleType: editingTransaction.saleType,
        deliveryStatus: editingTransaction.deliveryStatus,
        deliveryDate: editingTransaction.deliveryDate,
        deliveryAddress: editingTransaction.deliveryAddress
      };

      // Remove undefined values
      Object.keys(payload).forEach(key => {
        if (payload[key] === undefined || payload[key] === '') {
          delete payload[key];
        }
      });

      await transactionAPI.update(editingTransaction.id, payload);
      toast.success('Transaction updated successfully');
      setShowEditModal(false);
      setEditingTransaction(null);
      fetchData();
    } catch (err) {
      console.error('Error updating transaction:', err);
      toast.error('Failed to update transaction');
    }
  };

  const handleDelete = (transaction) => {
    setDeletingTransaction(transaction);
    setShowDeleteModal(true);
  };

  const confirmDelete = async (reason, attestation) => {
    try {
      await transactionAPI.delete(deletingTransaction.id, { reason, attestation });
      toast.success('Transaction deleted successfully');
      setShowDeleteModal(false);
      setDeletingTransaction(null);
      fetchData();
    } catch (err) {
      toast.error(err.response?.data || 'Failed to delete transaction');
    }
  };

  if (loading) return <div className="loading">Loading transactions...</div>;

  return (
    <div className="container">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px', flexWrap: 'wrap', gap: '16px' }}>
        <div>
          <h1 style={{ marginBottom: '8px' }}>Transaction Ledger</h1>
          <p style={{ color: 'var(--text-secondary)', fontSize: '14px' }}>Track sales and payments with customers</p>
        </div>
        <Button variant="success" leftIcon={<Plus size={20} />} onClick={() => setShowModal(true)}>
          New Transaction
        </Button>
      </div>

      <Card>
        <div style={{ display: 'flex', gap: '12px', marginBottom: '20px', flexWrap: 'wrap' }}>
          <div style={{ flex: '1', minWidth: '200px' }}>
            <Input
              placeholder="Search by customer or type..."
              value={searchTerm}
              onChangeText={setSearchTerm}
              leftIcon={<Search size={18} />}
              fullWidth
            />
          </div>
          <div style={{ minWidth: '150px' }}>
            <Select
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              options={[
                { value: '', label: 'All Types' },
                { value: 'sale', label: 'Sales' },
                { value: 'payment', label: 'Payments' }
              ]}
              fullWidth
            />
          </div>
        </div>

        <div style={{ overflowX: 'auto' }}>
          <table className="table">
            <thead>
              <tr>
                <th>Date</th>
                <th>Invoice #</th>
                <th>Type</th>
                <th>Customer</th>
                <th>Payment</th>
                <th style={{ textAlign: 'right' }}>Discount</th>
                <th style={{ textAlign: 'right' }}>Tax</th>
                <th style={{ textAlign: 'right' }}>Amount</th>
                <th>Status</th>
                <th style={{ width: '100px' }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {filteredTransactions.length === 0 ? (
                <tr>
                  <td colSpan="10" style={{ textAlign: 'center', padding: '40px' }}>
                    <Receipt size={48} style={{ opacity: 0.3, marginBottom: '16px' }} />
                    <p style={{ color: 'var(--text-secondary)' }}>
                      {searchTerm ? 'No transactions match your search' : 'No transactions found'}
                    </p>
                  </td>
                </tr>
              ) : (
                filteredTransactions.map((tx) => {
                  const customer = customers.find(c => c.id === tx.customerId);
                  const isOverdue = tx.dueDate && new Date(tx.dueDate) < new Date() && tx.type === 'sale';
                  
                  return (
                    <tr key={tx.id}>
                      <td style={{ color: 'var(--text-secondary)' }}>
                        {new Date(tx.date).toLocaleDateString()}
                      </td>
                      <td style={{ fontFamily: 'monospace', fontSize: '13px', color: 'var(--text-secondary)' }}>
                        {tx.invoiceNumber || '-'}
                      </td>
                      <td>
                        <Badge 
                          variant={tx.type === 'sale' ? 'success' : 'info'} 
                          size="sm"
                          dot
                        >
                          <span style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                            {tx.type === 'sale' ? <ShoppingCart size={12} /> : <DollarSign size={12} />}
                            {tx.type.charAt(0).toUpperCase() + tx.type.slice(1)}
                          </span>
                        </Badge>
                      </td>
                      <td>
                        <div>
                          <div style={{ fontWeight: '500' }}>{customer?.name || 'Unknown'}</div>
                          {tx.saleType && tx.saleType !== 'cash' && (
                            <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>
                              {tx.saleType === 'credit' ? 'Credit Sale' : 'On Credit'}
                            </div>
                          )}
                        </div>
                      </td>
                      <td>
                        {tx.paymentMethod ? (
                          <Badge variant="neutral" size="sm">
                            {tx.paymentMethod.toUpperCase()}
                          </Badge>
                        ) : '-'}
                      </td>
                      <td style={{ textAlign: 'right', color: 'var(--color-warning)' }}>
                        {tx.discountAmount && tx.discountAmount > 0 ? (
                          <span>-₹{parseFloat(tx.discountAmount).toFixed(2)}</span>
                        ) : '-'}
                      </td>
                      <td style={{ textAlign: 'right', color: 'var(--text-secondary)' }}>
                        {tx.taxAmount && tx.taxAmount > 0 ? (
                          <span>+₹{parseFloat(tx.taxAmount).toFixed(2)}</span>
                        ) : '-'}
                      </td>
                      <td style={{ textAlign: 'right', fontWeight: '600', color: 'var(--color-success)' }}>
                        ₹{tx.totalAmount?.toFixed(2) || '0.00'}
                      </td>
                      <td>
                        {tx.type === 'sale' && tx.deliveryStatus ? (
                          <Badge 
                            variant={
                              tx.deliveryStatus === 'delivered' ? 'success' : 
                              tx.deliveryStatus === 'in_transit' ? 'info' : 
                              'neutral'
                            }
                            size="sm"
                          >
                            {tx.deliveryStatus === 'in_transit' ? 'In Transit' : 
                             tx.deliveryStatus.charAt(0).toUpperCase() + tx.deliveryStatus.slice(1)}
                          </Badge>
                        ) : isOverdue ? (
                          <Badge variant="danger" size="sm">Overdue</Badge>
                        ) : '-'}
                      </td>
                      <td>
                        <div style={{ display: 'flex', gap: '8px' }}>
                          <Button 
                            variant="secondary" 
                            size="small"
                            onClick={() => handleEditClick(tx)}
                          >
                            <Edit2 size={14} />
                          </Button>
                          <Button 
                            variant="ghost" 
                            size="small"
                            onClick={() => handleDelete(tx)}
                            style={{ color: 'var(--color-danger)' }}
                          >
                            <Trash2 size={14} />
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
        title="New Transaction"
        size="lg"
      >
        <form onSubmit={handleSubmit}>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '20px', maxHeight: '70vh', overflowY: 'auto', padding: '4px' }}>
            {/* Transaction Details Section */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                Transaction Details
              </h3>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                <Select
                  label="Customer"
                  value={formData.customerId}
                  onChange={(e) => setFormData({...formData, customerId: e.target.value})}
                  options={[
                    { value: '', label: 'Select Customer' },
                    ...customers.map(c => ({ value: c.id, label: c.name }))
                  ]}
                  required
                  fullWidth
                />
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                  <Select
                    label="Transaction Type"
                    value={formData.type}
                    onChange={(e) => setFormData({...formData, type: e.target.value})}
                    options={[
                      { value: 'sale', label: 'Sale' },
                      { value: 'payment', label: 'Payment' }
                    ]}
                    required
                    fullWidth
                  />
                  {formData.type === 'sale' && (
                    <Select
                      label="Sale Type"
                      value={formData.saleType}
                      onChange={(e) => setFormData({...formData, saleType: e.target.value})}
                      options={[
                        { value: 'cash', label: 'Cash Sale' },
                        { value: 'credit', label: 'Credit Sale' }
                      ]}
                      required
                      fullWidth
                    />
                  )}
                </div>
              </div>
            </div>

            {/* Payment Information */}
            {formData.type === 'payment' ? (
              <div>
                <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                  Payment Information
                </h3>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                    <Input
                      label="Payment Amount"
                      type="number"
                      step="0.01"
                      value={formData.paymentAmount}
                      onChangeText={(value) => setFormData({...formData, paymentAmount: value})}
                      placeholder="0.00"
                      leftIcon={<DollarSign size={18} />}
                      required
                      fullWidth
                    />
                    <Select
                      label="Payment Method"
                      value={formData.paymentMethod}
                      onChange={(e) => setFormData({...formData, paymentMethod: e.target.value})}
                      options={[
                        { value: 'cash', label: 'Cash' },
                        { value: 'upi', label: 'UPI' },
                        { value: 'card', label: 'Card' },
                        { value: 'bank_transfer', label: 'Bank Transfer' },
                        { value: 'cheque', label: 'Cheque' }
                      ]}
                      required
                      fullWidth
                    />
                  </div>
                  {formData.paymentMethod !== 'cash' && (
                    <Input
                      label="Payment Reference / Transaction ID"
                      value={formData.paymentReference}
                      onChangeText={(value) => setFormData({...formData, paymentReference: value})}
                      placeholder="Enter reference number"
                      fullWidth
                    />
                  )}
                </div>
              </div>
            ) : (
              <>
                {/* Sale Items Section */}
                <div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
                    <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px' }}>
                      Sale Items ({formData.items.length})
                    </h3>
                  </div>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                    {formData.items.map((item, idx) => (
                      <Card key={idx} style={{ marginBottom: '4px', padding: '16px', backgroundColor: 'var(--bg-secondary)', border: '1px solid var(--border-light)' }}>
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
                          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                            <div style={{ 
                              width: '28px', 
                              height: '28px', 
                              borderRadius: '50%', 
                              backgroundColor: 'var(--color-primary-light)', 
                              color: 'var(--color-primary)',
                              display: 'flex',
                              alignItems: 'center',
                              justifyContent: 'center',
                              fontSize: '13px',
                              fontWeight: '600'
                            }}>
                              {idx + 1}
                            </div>
                            <span style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-primary)' }}>Item #{idx + 1}</span>
                          </div>
                          {formData.items.length > 1 && (
                            <Button
                              type="button"
                              variant="ghost"
                              size="sm"
                              leftIcon={<Trash2 size={14} />}
                              onClick={() => removeItem(idx)}
                              style={{ color: 'var(--color-danger)' }}
                            >
                              Remove
                            </Button>
                          )}
                        </div>
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '14px' }}>
                          <Select
                            label="Inventory Item"
                            value={item.inventoryLotId}
                            onChange={(e) => updateItem(idx, 'inventoryLotId', e.target.value)}
                            options={[
                              { value: '', label: 'Select Item' },
                              ...inventory.map(inv => ({ 
                                value: inv.id, 
                                label: `${inv.name} (${inv.quantity} ${inv.unit} available)` 
                              }))
                            ]}
                            required
                            fullWidth
                          />
                          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                            <Input
                              label="Quantity"
                              type="number"
                              step="0.01"
                              value={item.quantity}
                              onChangeText={(value) => updateItem(idx, 'quantity', value)}
                              placeholder="0.00"
                              required
                              fullWidth
                            />
                            <Input
                              label="Price per Unit (₹)"
                              type="number"
                              step="0.01"
                              value={item.pricePerUnit}
                              onChangeText={(value) => updateItem(idx, 'pricePerUnit', value)}
                              placeholder="0.00"
                              leftIcon={<DollarSign size={18} />}
                              required
                              fullWidth
                            />
                          </div>
                        </div>
                      </Card>
                    ))}
                  </div>
                  <Button
                    type="button"
                    variant="outline"
                    leftIcon={<Plus size={18} />}
                    onClick={addItem}
                    fullWidth
                    style={{ marginTop: '8px' }}
                  >
                    Add Another Item
                  </Button>
                </div>

                {/* Financial Details (Discounts & Tax) */}
                <div>
                  <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                    Discounts & Tax
                  </h3>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                    <Input
                      label="Discount Amount (₹)"
                      type="number"
                      step="0.01"
                      value={formData.discountAmount}
                      onChangeText={(value) => setFormData({...formData, discountAmount: value})}
                      placeholder="0.00"
                      fullWidth
                    />
                    <Input
                      label="Discount Percentage (%)"
                      type="number"
                      step="0.01"
                      value={formData.discountPercentage}
                      onChangeText={(value) => setFormData({...formData, discountPercentage: value})}
                      placeholder="0.00"
                      fullWidth
                    />
                    <Input
                      label="Tax Amount (₹)"
                      type="number"
                      step="0.01"
                      value={formData.taxAmount}
                      onChangeText={(value) => setFormData({...formData, taxAmount: value})}
                      placeholder="0.00"
                      fullWidth
                    />
                    <Input
                      label="Tax Percentage (%)"
                      type="number"
                      step="0.01"
                      value={formData.taxPercentage}
                      onChangeText={(value) => setFormData({...formData, taxPercentage: value})}
                      placeholder="0.00"
                      fullWidth
                    />
                  </div>
                </div>

                {/* Payment & Due Date */}
                <div>
                  <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                    Payment Details
                  </h3>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                    <Select
                      label="Payment Method"
                      value={formData.paymentMethod}
                      onChange={(e) => setFormData({...formData, paymentMethod: e.target.value})}
                      options={[
                        { value: 'cash', label: 'Cash' },
                        { value: 'upi', label: 'UPI' },
                        { value: 'card', label: 'Card' },
                        { value: 'bank_transfer', label: 'Bank Transfer' },
                        { value: 'cheque', label: 'Cheque' }
                      ]}
                      required
                      fullWidth
                    />
                    {formData.saleType === 'credit' && (
                      <Input
                        label="Due Date"
                        type="date"
                        value={formData.dueDate}
                        onChangeText={(value) => setFormData({...formData, dueDate: value})}
                        fullWidth
                      />
                    )}
                  </div>
                  {formData.paymentMethod !== 'cash' && (
                    <Input
                      label="Payment Reference / Transaction ID"
                      value={formData.paymentReference}
                      onChangeText={(value) => setFormData({...formData, paymentReference: value})}
                      placeholder="Enter reference number"
                      fullWidth
                      style={{ marginTop: '12px' }}
                    />
                  )}
                </div>

                {/* Delivery Tracking */}
                <div>
                  <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                    Delivery Information
                  </h3>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                      <Select
                        label="Delivery Status"
                        value={formData.deliveryStatus}
                        onChange={(e) => setFormData({...formData, deliveryStatus: e.target.value})}
                        options={[
                          { value: 'pending', label: 'Pending' },
                          { value: 'in_transit', label: 'In Transit' },
                          { value: 'delivered', label: 'Delivered' }
                        ]}
                        fullWidth
                      />
                      <Input
                        label="Delivery Date"
                        type="date"
                        value={formData.deliveryDate}
                        onChangeText={(value) => setFormData({...formData, deliveryDate: value})}
                        fullWidth
                      />
                    </div>
                    <div>
                      <label style={{ display: 'block', fontSize: '14px', fontWeight: '500', marginBottom: '6px', color: 'var(--text-primary)' }}>
                        Delivery Address
                      </label>
                      <textarea 
                        value={formData.deliveryAddress} 
                        onChange={(e) => setFormData({...formData, deliveryAddress: e.target.value})}
                        placeholder="Enter delivery address"
                        rows="2"
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
                      />
                    </div>
                  </div>
                </div>
              </>
            )}

            {/* Notes Section */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                Additional Notes
              </h3>
              <div>
                <label style={{ display: 'block', fontSize: '14px', fontWeight: '500', marginBottom: '6px', color: 'var(--text-primary)' }}>
                  Notes
                </label>
                <textarea 
                  value={formData.notes} 
                  onChange={(e) => setFormData({...formData, notes: e.target.value})}
                  placeholder="Additional notes or remarks"
                  rows="2"
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
                />
              </div>
            </div>

            {/* Action Buttons */}
            <div style={{ display: 'flex', gap: '12px', paddingTop: '20px', borderTop: '1px solid var(--border-light)' }}>
              <Button type="submit" variant="primary" fullWidth>
                Create Transaction
              </Button>
              <Button type="button" variant="secondary" onClick={handleCloseModal} fullWidth>
                Cancel
              </Button>
            </div>
          </div>
        </form>
      </Modal>

      {/* Edit Transaction Modal */}
      <Modal
        isOpen={showEditModal}
        onClose={() => {
          setShowEditModal(false);
          setEditingTransaction(null);
        }}
        title="Edit Transaction"
        size="lg"
      >
        {editingTransaction && (
          <form onSubmit={handleEditSubmit}>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '20px', maxHeight: '70vh', overflowY: 'auto', padding: '4px' }}>
              
              {/* Transaction Info */}
              <div style={{ padding: '12px', background: 'var(--background-hover)', borderRadius: '8px' }}>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '8px', fontSize: '13px' }}>
                  <div>
                    <span style={{ color: 'var(--text-secondary)' }}>Type:</span>{' '}
                    <Badge variant={editingTransaction.type === 'sale' ? 'success' : 'info'} size="sm">
                      {editingTransaction.type}
                    </Badge>
                  </div>
                  <div>
                    <span style={{ color: 'var(--text-secondary)' }}>Date:</span>{' '}
                    <strong>{new Date(editingTransaction.date).toLocaleDateString()}</strong>
                  </div>
                  {editingTransaction.updatedAt && (
                    <>
                      <div style={{ gridColumn: '1 / -1', marginTop: '8px', paddingTop: '8px', borderTop: '1px solid var(--border-light)' }}>
                        <span style={{ color: 'var(--text-secondary)', fontSize: '12px' }}>
                          Last updated: {new Date(editingTransaction.updatedAt).toLocaleString()} by {editingTransaction.updatedBy || 'system'}
                        </span>
                      </div>
                    </>
                  )}
                </div>
              </div>

              {/* Payment Details */}
              {editingTransaction.type === 'payment' && (
                <div>
                  <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', marginBottom: '16px' }}>
                    Payment Details
                  </h3>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                    <Input
                      label="Payment Amount"
                      type="number"
                      step="0.01"
                      value={editingTransaction.paymentAmount}
                      onChangeText={(value) => setEditingTransaction({...editingTransaction, paymentAmount: value})}
                      fullWidth
                    />
                    <Select
                      label="Payment Method"
                      value={editingTransaction.paymentMethod}
                      onChange={(e) => setEditingTransaction({...editingTransaction, paymentMethod: e.target.value})}
                      options={[
                        { value: 'cash', label: 'Cash' },
                        { value: 'upi', label: 'UPI' },
                        { value: 'card', label: 'Card' },
                        { value: 'bank_transfer', label: 'Bank Transfer' },
                        { value: 'cheque', label: 'Cheque' }
                      ]}
                      fullWidth
                    />
                  </div>
                  <div style={{ marginTop: '12px' }}>
                    <Input
                      label="Payment Reference"
                      value={editingTransaction.paymentReference}
                      onChangeText={(value) => setEditingTransaction({...editingTransaction, paymentReference: value})}
                      placeholder="Transaction ID, Cheque #, etc."
                      fullWidth
                    />
                  </div>
                </div>
              )}

              {/* Sale Details */}
              {editingTransaction.type === 'sale' && (
                <div>
                  <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', marginBottom: '16px' }}>
                    Sale Details
                  </h3>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                    <Input
                      label="Total Amount"
                      type="number"
                      step="0.01"
                      value={editingTransaction.totalAmount}
                      onChangeText={(value) => setEditingTransaction({...editingTransaction, totalAmount: value})}
                      fullWidth
                    />
                    <Input
                      label="Discount Amount"
                      type="number"
                      step="0.01"
                      value={editingTransaction.discountAmount}
                      onChangeText={(value) => setEditingTransaction({...editingTransaction, discountAmount: value})}
                      fullWidth
                    />
                    <Input
                      label="Tax Amount"
                      type="number"
                      step="0.01"
                      value={editingTransaction.taxAmount}
                      onChangeText={(value) => setEditingTransaction({...editingTransaction, taxAmount: value})}
                      fullWidth
                    />
                    <Input
                      label="Invoice Number"
                      value={editingTransaction.invoiceNumber}
                      onChangeText={(value) => setEditingTransaction({...editingTransaction, invoiceNumber: value})}
                      fullWidth
                    />
                    <Select
                      label="Sale Type"
                      value={editingTransaction.saleType}
                      onChange={(e) => setEditingTransaction({...editingTransaction, saleType: e.target.value})}
                      options={[
                        { value: 'cash', label: 'Cash Sale' },
                        { value: 'credit', label: 'Credit Sale' }
                      ]}
                      fullWidth
                    />
                    <Input
                      label="Due Date"
                      type="date"
                      value={editingTransaction.dueDate}
                      onChangeText={(value) => setEditingTransaction({...editingTransaction, dueDate: value})}
                      fullWidth
                    />
                  </div>
                </div>
              )}

              {/* Delivery Information */}
              {editingTransaction.type === 'sale' && (
                <div>
                  <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', marginBottom: '16px' }}>
                    Delivery Information
                  </h3>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                    <Select
                      label="Delivery Status"
                      value={editingTransaction.deliveryStatus}
                      onChange={(e) => setEditingTransaction({...editingTransaction, deliveryStatus: e.target.value})}
                      options={[
                        { value: 'pending', label: 'Pending' },
                        { value: 'in_transit', label: 'In Transit' },
                        { value: 'delivered', label: 'Delivered' },
                        { value: 'cancelled', label: 'Cancelled' }
                      ]}
                      fullWidth
                    />
                    <Input
                      label="Delivery Date"
                      type="datetime-local"
                      value={editingTransaction.deliveryDate ? editingTransaction.deliveryDate.slice(0, 16) : ''}
                      onChangeText={(value) => setEditingTransaction({...editingTransaction, deliveryDate: value})}
                      fullWidth
                    />
                  </div>
                  <div style={{ marginTop: '12px' }}>
                    <Input
                      label="Delivery Address"
                      value={editingTransaction.deliveryAddress}
                      onChangeText={(value) => setEditingTransaction({...editingTransaction, deliveryAddress: value})}
                      placeholder="Enter delivery address"
                      fullWidth
                    />
                  </div>
                </div>
              )}

              {/* Notes */}
              <div>
                <label style={{ display: 'block', marginBottom: '8px', fontSize: '14px', fontWeight: '500', color: 'var(--text-primary)' }}>
                  Notes
                </label>
                <textarea
                  value={editingTransaction.notes}
                  onChange={(e) => setEditingTransaction({...editingTransaction, notes: e.target.value})}
                  placeholder="Add notes about this transaction..."
                  rows={3}
                  style={{
                    width: '100%',
                    padding: '12px',
                    border: '1px solid var(--border)',
                    borderRadius: '8px',
                    fontSize: '14px',
                    fontFamily: 'inherit',
                    resize: 'vertical'
                  }}
                />
              </div>

              {/* Action Buttons */}
              <div style={{ display: 'flex', gap: '12px', paddingTop: '20px', borderTop: '1px solid var(--border-light)' }}>
                <Button type="submit" variant="primary" fullWidth>
                  Update Transaction
                </Button>
                <Button 
                  type="button" 
                  variant="secondary" 
                  onClick={() => {
                    setShowEditModal(false);
                    setEditingTransaction(null);
                  }} 
                  fullWidth
                >
                  Cancel
                </Button>
              </div>
            </div>
          </form>
        )}
      </Modal>

      <DeleteConfirmationModal
        isOpen={showDeleteModal}
        onClose={() => {
          setShowDeleteModal(false);
          setDeletingTransaction(null);
        }}
        onConfirm={confirmDelete}
        title="Delete Transaction"
        entityType="Transaction"
        entityName={deletingTransaction ? `${deletingTransaction.type.toUpperCase()} - ${deletingTransaction.customerName}` : ''}
        warningMessage="This will mark the transaction as deleted. All transaction records will be preserved for audit purposes. Customer balances may need to be adjusted."
      />
    </div>
  );
}

export default Transactions;
