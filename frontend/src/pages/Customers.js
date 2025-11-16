import React, { useState, useEffect } from 'react';
import { customerAPI } from '../services/api';
import { Card, Button, Badge, Input, Select, Modal, useToast } from '../components/ui';
import { Users, Plus, Edit2, Trash2, Search, Phone, MapPin, CheckCircle, XCircle } from 'lucide-react';
import DeleteConfirmationModal from '../components/DeleteConfirmationModal';

function Customers() {
  const toast = useToast();
  const [customers, setCustomers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingCustomer, setEditingCustomer] = useState(null);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deletingCustomer, setDeletingCustomer] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [formData, setFormData] = useState({
    name: '',
    address: '',
    contactNumber: '',
    email: '',
    whatsappNumber: '',
    alternateContact: '',
    customerType: 'retail',
    businessName: '',
    gstin: '',
    creditLimit: '',
    paymentTermsDays: '',
    status: 'active',
    kycDocumentType: '',
    kycDocumentNumber: '',
    aadhaarVerified: false,
    notes: ''
  });

  useEffect(() => {
    fetchCustomers();
  }, []);

  const fetchCustomers = async () => {
    setLoading(true);
    try {
      const res = await customerAPI.getAll();
      setCustomers(res.data);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      if (editingCustomer) {
        await customerAPI.update(editingCustomer.id, formData);
        toast.success('Customer updated successfully');
      } else {
        await customerAPI.create(formData);
        toast.success('Customer added successfully');
      }
      handleCloseModal();
      fetchCustomers();
    } catch (err) {
      toast.error(editingCustomer ? 'Failed to update customer' : 'Failed to add customer');
    }
  };

  const handleEdit = (customer) => {
    setEditingCustomer(customer);
    setFormData({
      name: customer.name,
      address: customer.address || '',
      contactNumber: customer.contactNumber || '',
      email: customer.email || '',
      whatsappNumber: customer.whatsappNumber || '',
      alternateContact: customer.alternateContact || '',
      customerType: customer.customerType || 'retail',
      businessName: customer.businessName || '',
      gstin: customer.gstin || '',
      creditLimit: customer.creditLimit || '',
      paymentTermsDays: customer.paymentTermsDays || '',
      status: customer.status || 'active',
      kycDocumentType: customer.kycDocumentType || '',
      kycDocumentNumber: customer.kycDocumentNumber || '',
      aadhaarVerified: customer.aadhaarVerified || false,
      notes: customer.notes || ''
    });
    setShowModal(true);
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setEditingCustomer(null);
    setFormData({
      name: '',
      address: '',
      contactNumber: '',
      email: '',
      whatsappNumber: '',
      alternateContact: '',
      customerType: 'retail',
      businessName: '',
      gstin: '',
      creditLimit: '',
      paymentTermsDays: '',
      status: 'active',
      kycDocumentType: '',
      kycDocumentNumber: '',
      aadhaarVerified: false,
      notes: ''
    });
  };

  const handleDelete = (customer) => {
    setDeletingCustomer(customer);
    setShowDeleteModal(true);
  };

  const confirmDelete = async (reason, attestation) => {
    try {
      await customerAPI.delete(deletingCustomer.id, { reason, attestation });
      toast.success('Customer deleted successfully');
      setShowDeleteModal(false);
      setDeletingCustomer(null);
      fetchCustomers();
    } catch (err) {
      toast.error(err.response?.data || 'Failed to delete customer');
    }
  };

  const filteredCustomers = customers.filter(customer => {
    const searchLower = searchTerm.toLowerCase();
    return searchTerm === '' ||
      customer.name.toLowerCase().includes(searchLower) ||
      (customer.contactNumber && customer.contactNumber.includes(searchTerm)) ||
      (customer.address && customer.address.toLowerCase().includes(searchLower));
  });

  if (loading) return <div className="loading">Loading customers...</div>;

  return (
    <div className="container">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px', flexWrap: 'wrap', gap: '16px' }}>
        <div>
          <h1 style={{ marginBottom: '8px' }}>Customer Management</h1>
          <p style={{ color: 'var(--text-secondary)', fontSize: '14px' }}>Manage your customers and their KYC details</p>
        </div>
        <Button variant="primary" leftIcon={<Plus size={20} />} onClick={() => setShowModal(true)}>
          Add Customer
        </Button>
      </div>

      <Card>
        <div style={{ marginBottom: '20px' }}>
          <Input
            placeholder="Search by name, phone, or address..."
            value={searchTerm}
            onChangeText={setSearchTerm}
            leftIcon={<Search size={18} />}
            fullWidth
          />
        </div>
        <div style={{ overflowX: 'auto' }}>
          <table className="table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Type</th>
                <th>Contact</th>
                <th>Credit Limit</th>
                <th>Balance</th>
                <th>Status</th>
                <th style={{ textAlign: 'right' }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {filteredCustomers.length === 0 ? (
                <tr>
                  <td colSpan="7" style={{ textAlign: 'center', padding: '40px' }}>
                    <Users size={48} style={{ opacity: 0.3, marginBottom: '16px' }} />
                    <p style={{ color: 'var(--text-secondary)' }}>
                      {searchTerm ? 'No customers match your search' : 'No customers found'}
                    </p>
                  </td>
                </tr>
              ) : (
                filteredCustomers.map((customer) => {
                  const creditUsed = customer.currentBalance || 0;
                  const creditLimit = customer.creditLimit || 0;
                  const creditPercent = creditLimit > 0 ? (creditUsed / creditLimit) * 100 : 0;
                  
                  return (
                    <tr key={customer.id}>
                      <td>
                        <div>
                          <div style={{ fontWeight: '500' }}>{customer.name}</div>
                          {customer.businessName && <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>{customer.businessName}</div>}
                        </div>
                      </td>
                      <td>
                        <Badge 
                          variant={customer.customerType === 'b2b' ? 'primary' : customer.customerType === 'wholesale' ? 'info' : 'neutral'} 
                          size="sm"
                        >
                          {customer.customerType ? customer.customerType.toUpperCase() : 'RETAIL'}
                        </Badge>
                      </td>
                      <td>
                        {customer.contactNumber ? (
                          <div style={{ display: 'flex', alignItems: 'center', gap: '6px', color: 'var(--text-secondary)' }}>
                            <Phone size={14} />
                            {customer.contactNumber}
                          </div>
                        ) : '-'}
                      </td>
                      <td style={{ fontSize: '14px' }}>
                        {creditLimit > 0 ? `₹${creditLimit.toFixed(2)}` : '-'}
                      </td>
                      <td>
                        {creditLimit > 0 ? (
                          <div>
                            <div style={{ fontWeight: '600', color: creditUsed > 0 ? 'var(--color-danger)' : 'var(--color-success)' }}>
                              ₹{creditUsed.toFixed(2)}
                            </div>
                            <div style={{ fontSize: '11px', color: 'var(--text-secondary)' }}>
                              {creditPercent.toFixed(1)}% used
                            </div>
                          </div>
                        ) : '-'}
                      </td>
                      <td>
                        <Badge 
                          variant={customer.status === 'active' ? 'success' : customer.status === 'blocked' ? 'danger' : 'warning'} 
                          size="sm"
                          dot
                        >
                          {customer.status ? customer.status.charAt(0).toUpperCase() + customer.status.slice(1) : 'Active'}
                        </Badge>
                      </td>
                      <td>
                        <div style={{ display: 'flex', gap: '8px', justifyContent: 'flex-end' }}>
                          <Button 
                            variant="ghost" 
                            size="sm" 
                            leftIcon={<Edit2 size={16} />}
                            onClick={() => handleEdit(customer)}
                          >
                            Edit
                          </Button>
                          <Button 
                            variant="ghost" 
                            size="sm" 
                            leftIcon={<Trash2 size={16} />}
                            onClick={() => handleDelete(customer)}
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
        title={editingCustomer ? 'Edit Customer' : 'Add Customer'}
        size="lg"
      >
        <form onSubmit={handleSubmit}>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '20px', maxHeight: '70vh', overflowY: 'auto', padding: '4px' }}>
            {/* Basic Information Section */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                Basic Information
              </h3>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                <Input
                  label="Full Name"
                  value={formData.name}
                  onChangeText={(value) => setFormData({...formData, name: value})}
                  placeholder="Enter customer name"
                  required
                  fullWidth
                />
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                  <Select
                    label="Customer Type"
                    value={formData.customerType}
                    onChange={(e) => setFormData({...formData, customerType: e.target.value})}
                    options={[
                      { value: 'retail', label: 'Retail' },
                      { value: 'b2b', label: 'B2B' },
                      { value: 'b2c', label: 'B2C' },
                      { value: 'wholesale', label: 'Wholesale' }
                    ]}
                    required
                    fullWidth
                  />
                  <Select
                    label="Status"
                    value={formData.status}
                    onChange={(e) => setFormData({...formData, status: e.target.value})}
                    options={[
                      { value: 'active', label: 'Active' },
                      { value: 'inactive', label: 'Inactive' },
                      { value: 'blocked', label: 'Blocked' }
                    ]}
                    required
                    fullWidth
                  />
                </div>
              </div>
            </div>

            {/* Contact Information */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                Contact Information
              </h3>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                  <Input
                    label="Contact Number"
                    type="tel"
                    value={formData.contactNumber}
                    onChangeText={(value) => setFormData({...formData, contactNumber: value})}
                    placeholder="Enter phone number"
                    leftIcon={<Phone size={18} />}
                    fullWidth
                  />
                  <Input
                    label="WhatsApp Number"
                    type="tel"
                    value={formData.whatsappNumber}
                    onChangeText={(value) => setFormData({...formData, whatsappNumber: value})}
                    placeholder="WhatsApp number"
                    fullWidth
                  />
                </div>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                  <Input
                    label="Email"
                    type="email"
                    value={formData.email}
                    onChangeText={(value) => setFormData({...formData, email: value})}
                    placeholder="customer@example.com"
                    fullWidth
                  />
                  <Input
                    label="Alternate Contact"
                    type="tel"
                    value={formData.alternateContact}
                    onChangeText={(value) => setFormData({...formData, alternateContact: value})}
                    placeholder="Alternate number"
                    fullWidth
                  />
                </div>
              </div>
            </div>

            {/* Business Details */}
            {(formData.customerType === 'b2b' || formData.customerType === 'wholesale') && (
              <div>
                <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                  Business Details
                </h3>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                  <Input
                    label="Business Name"
                    value={formData.businessName}
                    onChangeText={(value) => setFormData({...formData, businessName: value})}
                    placeholder="Enter business name"
                    fullWidth
                  />
                  <Input
                    label="GSTIN"
                    value={formData.gstin}
                    onChangeText={(value) => setFormData({...formData, gstin: value})}
                    placeholder="GST Identification Number (15 digits)"
                    maxLength="15"
                    fullWidth
                  />
                </div>
              </div>
            )}

            {/* Credit Management */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                Credit Management
              </h3>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                <Input
                  label="Credit Limit (₹)"
                  type="number"
                  step="0.01"
                  value={formData.creditLimit}
                  onChangeText={(value) => setFormData({...formData, creditLimit: value})}
                  placeholder="0.00"
                  fullWidth
                />
                <Input
                  label="Payment Terms (Days)"
                  type="number"
                  value={formData.paymentTermsDays}
                  onChangeText={(value) => setFormData({...formData, paymentTermsDays: value})}
                  placeholder="e.g., 30"
                  fullWidth
                />
              </div>
              {editingCustomer && editingCustomer.currentBalance != null && (
                <div style={{ marginTop: '12px', padding: '12px', backgroundColor: 'var(--bg-secondary)', borderRadius: '6px', border: '1px solid var(--border-light)' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <span style={{ fontSize: '13px', color: 'var(--text-secondary)' }}>Current Balance:</span>
                    <span style={{ fontSize: '16px', fontWeight: '600', color: editingCustomer.currentBalance > 0 ? 'var(--color-danger)' : 'var(--color-success)' }}>
                      ₹{parseFloat(editingCustomer.currentBalance).toFixed(2)}
                    </span>
                  </div>
                </div>
              )}
            </div>

            {/* Address Section */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                Address Details
              </h3>
              <div style={{ marginBottom: '0' }}>
                <label style={{ display: 'block', fontSize: '14px', fontWeight: '500', marginBottom: '6px', color: 'var(--text-primary)' }}>
                  <MapPin size={14} style={{ display: 'inline', marginRight: '6px', verticalAlign: 'middle' }} />
                  Full Address
                </label>
                <textarea 
                  value={formData.address} 
                  onChange={(e) => setFormData({...formData, address: e.target.value})}
                  placeholder="Enter complete address"
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
                />
              </div>
            </div>

            {/* KYC Verification Section */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                KYC Documents
              </h3>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                  <Select
                    label="Document Type"
                    value={formData.kycDocumentType}
                    onChange={(e) => setFormData({...formData, kycDocumentType: e.target.value})}
                    options={[
                      { value: '', label: 'Select Document' },
                      { value: 'aadhaar', label: 'Aadhaar Card' },
                      { value: 'pan', label: 'PAN Card' },
                      { value: 'passport', label: 'Passport' },
                      { value: 'driving_license', label: 'Driving License' },
                      { value: 'voter_id', label: 'Voter ID' }
                    ]}
                    fullWidth
                  />
                  <Input
                    label="Document Number"
                    value={formData.kycDocumentNumber}
                    onChangeText={(value) => setFormData({...formData, kycDocumentNumber: value})}
                    placeholder="Enter document number"
                    fullWidth
                  />
                </div>
                <div style={{ 
                  padding: '16px', 
                  backgroundColor: 'var(--bg-secondary)', 
                  borderRadius: '8px',
                  border: '1px solid var(--border-light)'
                }}>
                  <label style={{ display: 'flex', alignItems: 'center', gap: '12px', cursor: 'pointer' }}>
                    <input 
                      type="checkbox" 
                      checked={formData.aadhaarVerified} 
                      onChange={(e) => setFormData({...formData, aadhaarVerified: e.target.checked})}
                      style={{ 
                        width: '20px', 
                        height: '20px', 
                        cursor: 'pointer',
                        accentColor: 'var(--color-primary)'
                      }}
                    />
                    <div>
                      <span style={{ fontSize: '15px', fontWeight: '500', color: 'var(--text-primary)', display: 'block' }}>
                        KYC Verified
                      </span>
                      <span style={{ fontSize: '13px', color: 'var(--text-secondary)' }}>
                        Customer identity has been verified
                      </span>
                    </div>
                  </label>
                </div>
              </div>
            </div>

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
                  placeholder="Additional notes or remarks about the customer"
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
                />
              </div>
            </div>

            {/* Action Buttons */}
            <div style={{ display: 'flex', gap: '12px', marginTop: '8px', paddingTop: '20px', borderTop: '1px solid var(--border-light)' }}>
              <Button type="submit" variant="primary" fullWidth>
                {editingCustomer ? 'Update Customer' : 'Add Customer'}
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
          setDeletingCustomer(null);
        }}
        onConfirm={confirmDelete}
        title="Delete Customer"
        entityType="Customer"
        entityName={deletingCustomer ? deletingCustomer.name : ''}
        warningMessage="This will mark the customer as deleted. All customer records and transaction history will be preserved for audit purposes."
      />
    </div>
  );
}

export default Customers;
