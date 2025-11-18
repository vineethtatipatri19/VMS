import React, { useState, useEffect } from 'react';
import { inventoryAPI } from '../services/api';
import { Card, Button, Badge, Input, Select, Modal, useToast } from '../components/ui';
import { Package, Plus, Edit2, Trash2, Search } from 'lucide-react';
import DeleteConfirmationModal from '../components/DeleteConfirmationModal';
import { normalizeInventory, formatCurrency, formatDate, getStatusVariant } from '../utils/dataHelpers';

function Inventory() {
  const toast = useToast();
  const [inventory, setInventory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [editingItem, setEditingItem] = useState(null);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deletingItem, setDeletingItem] = useState(null);
  const [filter, setFilter] = useState('');
  const [searchTerm, setSearchTerm] = useState('');
  const [formData, setFormData] = useState({
    name: '',
    variant: '',
    quantity: '',
    unit: 'kg',
    purchaseDate: new Date().toISOString().split('T')[0],
    expiryDate: '',
    category: '',
    subCategory: '',
    costPrice: '',
    sellingPrice: '',
    supplierId: '',
    supplierName: '',
    purchaseInvoice: '',
    minStockLevel: '',
    reorderPoint: '',
    shelfLifeDays: '',
    storageLocation: '',
    barcode: '',
    sku: '',
    hsnCode: '',
    gstRate: '',
    packagingType: '',
    notes: ''
  });

  useEffect(() => {
    fetchInventory();
  }, [filter]);

  const fetchInventory = async () => {
    setLoading(true);
    try {
      const params = filter ? { status: filter } : { sort: 'expiry' };
      const res = await inventoryAPI.getAll(params);
      console.log('Inventory API Response:', res);
      console.log('Response data:', res.data);
      const dataArray = res.data.data || res.data || [];
      console.log('Data array:', dataArray);
      const normalized = dataArray.map(normalizeInventory);
      console.log('Normalized inventory:', normalized);
      setInventory(normalized);
    } catch (err) {
      setError('Failed to load inventory');
      console.error('Error fetching inventory:', err);
      console.error('Error response:', err.response);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    const qty = parseFloat(formData.quantity);
    if (isNaN(qty) || qty <= 0) {
      toast.error('Please enter a valid quantity greater than 0');
      return;
    }
    
    try {
      const dataToSend = { ...formData, quantity: qty };
      
      if (editingItem) {
        await inventoryAPI.update(editingItem.id, dataToSend);
        toast.success('Inventory updated successfully');
      } else {
        await inventoryAPI.create(dataToSend);
        toast.success('Inventory added successfully');
      }
      
      handleCloseModal();
      fetchInventory();
    } catch (err) {
      toast.error(editingItem ? 'Failed to update inventory' : 'Failed to add inventory');
    }
  };

  const handleEdit = (item) => {
    setEditingItem(item);
    setFormData({
      name: item.name,
      variant: item.variant || '',
      quantity: item.quantity.toString(),
      unit: item.unit,
      purchaseDate: item.purchaseDate,
      expiryDate: item.expiryDate,
      category: item.category || '',
      subCategory: item.subCategory || '',
      costPrice: item.costPrice || '',
      sellingPrice: item.sellingPrice || '',
      supplierId: item.supplierId || '',
      supplierName: item.supplierName || '',
      purchaseInvoice: item.purchaseInvoice || '',
      minStockLevel: item.minStockLevel || '',
      reorderPoint: item.reorderPoint || '',
      shelfLifeDays: item.shelfLifeDays || '',
      storageLocation: item.storageLocation || '',
      barcode: item.barcode || '',
      sku: item.sku || '',
      hsnCode: item.hsnCode || '',
      gstRate: item.gstRate || '',
      packagingType: item.packagingType || '',
      notes: item.notes || ''
    });
    setShowModal(true);
  };

  const handleDelete = (item) => {
    setDeletingItem(item);
    setShowDeleteModal(true);
  };

  const confirmDelete = async (reason, attestation) => {
    try {
      await inventoryAPI.delete(deletingItem.id, { reason, attestation });
      toast.success('Inventory item deleted successfully');
      setShowDeleteModal(false);
      setDeletingItem(null);
      fetchInventory();
    } catch (err) {
      toast.error(err.response?.data || 'Failed to delete item');
    }
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setEditingItem(null);
    setFormData({
      name: '',
      variant: '',
      quantity: '',
      unit: 'kg',
      purchaseDate: new Date().toISOString().split('T')[0],
      expiryDate: '',
      category: '',
      subCategory: '',
      costPrice: '',
      sellingPrice: '',
      supplierId: '',
      supplierName: '',
      purchaseInvoice: '',
      minStockLevel: '',
      reorderPoint: '',
      shelfLifeDays: '',
      storageLocation: '',
      barcode: '',
      sku: '',
      hsnCode: '',
      gstRate: '',
      packagingType: '',
      notes: ''
    });
  };

  const filteredInventory = inventory.filter(item => {
    const matchesSearch = searchTerm === '' || 
      item.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      (item.variant && item.variant.toLowerCase().includes(searchTerm.toLowerCase()));
    return matchesSearch;
  });

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
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px', flexWrap: 'wrap', gap: '16px' }}>
        <div>
          <h1 style={{ marginBottom: '8px' }}>Inventory Management</h1>
          <p style={{ color: 'var(--text-secondary)', fontSize: '14px' }}>Manage your perishable goods inventory with FEFO sorting</p>
        </div>
        <Button variant="primary" leftIcon={<Plus size={20} />} onClick={() => setShowModal(true)}>
          Add Inventory
        </Button>
      </div>

      {error && <div className="error">{error}</div>}

      <Card>
        <div style={{ display: 'flex', gap: '16px', marginBottom: '20px', flexWrap: 'wrap' }}>
          <div style={{ flex: '1', minWidth: '250px' }}>
            <Input
              placeholder="Search by name or variant..."
              value={searchTerm}
              onChangeText={setSearchTerm}
              leftIcon={<Search size={18} />}
            />
          </div>
          <div style={{ minWidth: '200px' }}>
            <Select
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              options={[
                { value: '', label: 'All (FEFO Sorted)' },
                { value: 'fresh', label: 'Fresh' },
                { value: 'expiring_soon', label: 'Expiring Soon' },
                { value: 'expired', label: 'Expired' }
              ]}
            />
          </div>
        </div>

        <div style={{ overflowX: 'auto' }}>
          <table className="table">
            <thead>
              <tr>
                <th>Status</th>
                <th>Name</th>
                <th>Category</th>
                <th>Supplier</th>
                <th>Quantity</th>
                <th>Cost/Selling</th>
                <th>Margin %</th>
                <th>Expiry Date</th>
                <th style={{ textAlign: 'right' }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {filteredInventory.length === 0 ? (
                <tr>
                  <td colSpan="9" style={{ textAlign: 'center', padding: '40px' }}>
                    <Package size={48} style={{ opacity: 0.3, marginBottom: '16px' }} />
                    <p style={{ color: 'var(--text-secondary)' }}>
                      {searchTerm ? 'No items match your search' : 'No inventory items found'}
                    </p>
                  </td>
                </tr>
              ) : (
                filteredInventory.map((item) => {
                  const status = getStatus(item.expiryDate);
                  const marginPercent = item.marginPercentage || 0;
                  const marginColor = marginPercent >= 30 ? 'var(--color-success)' : marginPercent >= 10 ? 'var(--color-warning)' : 'var(--color-danger)';
                  
                  return (
                    <tr key={item.id}>
                      <td>
                        <Badge variant={status.class.replace('badge-', '')} size="sm">
                          {status.label}
                        </Badge>
                      </td>
                      <td>
                        <div>
                          <div style={{ fontWeight: '500' }}>{item.name}</div>
                          {item.variant && <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>{item.variant}</div>}
                        </div>
                      </td>
                      <td>
                        {item.category ? (
                          <Badge variant="neutral" size="sm">{item.category}</Badge>
                        ) : '-'}
                      </td>
                      <td style={{ color: 'var(--text-secondary)', fontSize: '14px' }}>
                        {item.supplierName || '-'}
                      </td>
                      <td><strong>{item.quantity}</strong> {item.unit}</td>
                      <td style={{ fontSize: '13px' }}>
                        {item.costPrice && item.sellingPrice ? (
                          <div>
                            <div style={{ color: 'var(--text-secondary)' }}>₹{parseFloat(item.costPrice).toFixed(2)}</div>
                            <div style={{ color: 'var(--color-success)', fontWeight: '500' }}>₹{parseFloat(item.sellingPrice).toFixed(2)}</div>
                          </div>
                        ) : '-'}
                      </td>
                      <td>
                        {item.marginPercentage != null ? (
                          <span style={{ fontWeight: '600', color: marginColor }}>
                            {parseFloat(item.marginPercentage).toFixed(1)}%
                          </span>
                        ) : '-'}
                      </td>
                      <td>{item.expiryDate}</td>
                      <td>
                        <div style={{ display: 'flex', gap: '8px', justifyContent: 'flex-end' }}>
                          <Button 
                            variant="ghost" 
                            size="sm" 
                            leftIcon={<Edit2 size={16} />}
                            onClick={() => handleEdit(item)}
                          >
                            Edit
                          </Button>
                          <Button 
                            variant="ghost" 
                            size="sm" 
                            leftIcon={<Trash2 size={16} />}
                            onClick={() => handleDelete(item)}
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
        title={editingItem ? 'Edit Inventory Item' : 'Add Inventory Item'}
        size="lg"
      >
        <form onSubmit={handleSubmit}>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '20px', maxHeight: '70vh', overflowY: 'auto', padding: '4px' }}>
            {/* Product Information Section */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                Product Information
              </h3>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                <Input
                  label="Item Name"
                  value={formData.name}
                  onChangeText={(value) => setFormData({...formData, name: value})}
                  placeholder="e.g., Tomato, Potato"
                  leftIcon={<Package size={18} />}
                  required
                  fullWidth
                />
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                  <Input
                    label="Variant"
                    value={formData.variant}
                    onChangeText={(value) => setFormData({...formData, variant: value})}
                    placeholder="e.g., Red, Large"
                    fullWidth
                  />
                  <Select
                    label="Category"
                    value={formData.category}
                    onChange={(e) => setFormData({...formData, category: e.target.value})}
                    options={[
                      { value: '', label: 'Select Category' },
                      { value: 'Vegetables', label: 'Vegetables' },
                      { value: 'Fruits', label: 'Fruits' },
                      { value: 'Dairy', label: 'Dairy' },
                      { value: 'Grains', label: 'Grains' },
                      { value: 'Meat', label: 'Meat & Poultry' },
                      { value: 'Seafood', label: 'Seafood' },
                      { value: 'Other', label: 'Other' }
                    ]}
                    fullWidth
                  />
                </div>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                  <Input
                    label="Sub-category"
                    value={formData.subCategory}
                    onChangeText={(value) => setFormData({...formData, subCategory: value})}
                    placeholder="e.g., Leafy Greens"
                    fullWidth
                  />
                  <Input
                    label="Barcode"
                    value={formData.barcode}
                    onChangeText={(value) => setFormData({...formData, barcode: value})}
                    placeholder="Scan or enter barcode"
                    fullWidth
                  />
                </div>
              </div>
            </div>

            {/* Pricing Section */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                Pricing & Cost
              </h3>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                <Input
                  label="Cost Price (₹)"
                  type="number"
                  step="0.01"
                  value={formData.costPrice}
                  onChangeText={(value) => setFormData({...formData, costPrice: value})}
                  placeholder="0.00"
                  fullWidth
                />
                <Input
                  label="Selling Price (₹)"
                  type="number"
                  step="0.01"
                  value={formData.sellingPrice}
                  onChangeText={(value) => setFormData({...formData, sellingPrice: value})}
                  placeholder="0.00"
                  fullWidth
                />
              </div>
              {formData.costPrice && formData.sellingPrice && parseFloat(formData.costPrice) > 0 && (
                <div style={{ marginTop: '12px', padding: '12px', backgroundColor: 'var(--bg-secondary)', borderRadius: '6px', border: '1px solid var(--border-light)' }}>
                  <span style={{ fontSize: '13px', color: 'var(--text-secondary)' }}>Profit Margin: </span>
                  <span style={{ fontSize: '15px', fontWeight: '600', color: 'var(--color-success)' }}>
                    {(((parseFloat(formData.sellingPrice) - parseFloat(formData.costPrice)) / parseFloat(formData.costPrice)) * 100).toFixed(2)}%
                  </span>
                </div>
              )}
            </div>

            {/* Supplier Information */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                Supplier Information
              </h3>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                <Input
                  label="Supplier Name"
                  value={formData.supplierName}
                  onChangeText={(value) => setFormData({...formData, supplierName: value})}
                  placeholder="Enter supplier name"
                  fullWidth
                />
                <Input
                  label="Purchase Invoice"
                  value={formData.purchaseInvoice}
                  onChangeText={(value) => setFormData({...formData, purchaseInvoice: value})}
                  placeholder="Invoice number"
                  fullWidth
                />
              </div>
            </div>

            {/* Quantity & Stock Management */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                Quantity & Stock Management
              </h3>
              <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '12px', marginBottom: '12px' }}>
                <Input
                  label="Quantity"
                  type="number"
                  step="0.01"
                  value={formData.quantity}
                  onChangeText={(value) => setFormData({...formData, quantity: value})}
                  placeholder="0.00"
                  required
                  fullWidth
                />
                <Select
                  label="Unit"
                  value={formData.unit}
                  onChange={(e) => setFormData({...formData, unit: e.target.value})}
                  options={[
                    { value: 'kg', label: 'Kilograms (kg)' },
                    { value: 'g', label: 'Grams (g)' },
                    { value: 'l', label: 'Liters (l)' },
                    { value: 'ml', label: 'Milliliters (ml)' },
                    { value: 'pcs', label: 'Pieces' },
                    { value: 'box', label: 'Box' }
                  ]}
                  required
                  fullWidth
                />
              </div>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '12px' }}>
                <Input
                  label="Min Stock Level"
                  type="number"
                  step="0.01"
                  value={formData.minStockLevel}
                  onChangeText={(value) => setFormData({...formData, minStockLevel: value})}
                  placeholder="0.00"
                  fullWidth
                />
                <Input
                  label="Reorder Point"
                  type="number"
                  step="0.01"
                  value={formData.reorderPoint}
                  onChangeText={(value) => setFormData({...formData, reorderPoint: value})}
                  placeholder="0.00"
                  fullWidth
                />
                <Input
                  label="Shelf Life (days)"
                  type="number"
                  value={formData.shelfLifeDays}
                  onChangeText={(value) => setFormData({...formData, shelfLifeDays: value})}
                  placeholder="e.g., 7"
                  fullWidth
                />
              </div>
            </div>

            {/* Dates Section */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                Important Dates
              </h3>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                <Input
                  label="Purchase Date"
                  type="date"
                  value={formData.purchaseDate}
                  onChangeText={(value) => setFormData({...formData, purchaseDate: value})}
                  required
                  fullWidth
                />
                <Input
                  label="Expiry Date"
                  type="date"
                  value={formData.expiryDate}
                  onChangeText={(value) => setFormData({...formData, expiryDate: value})}
                  required
                  fullWidth
                />
              </div>
            </div>

            {/* GST & Tax Information */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                GST & Tax Information
              </h3>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                <Input
                  label="HSN Code"
                  value={formData.hsnCode}
                  onChangeText={(value) => setFormData({...formData, hsnCode: value})}
                  placeholder="e.g., 0701"
                  fullWidth
                />
                <Select
                  label="GST Rate (%)"
                  value={formData.gstRate}
                  onChange={(e) => setFormData({...formData, gstRate: e.target.value})}
                  options={[
                    { value: '', label: 'Select GST Rate' },
                    { value: '0', label: '0% (Exempt)' },
                    { value: '5', label: '5%' },
                    { value: '12', label: '12%' },
                    { value: '18', label: '18%' },
                    { value: '28', label: '28%' }
                  ]}
                  fullWidth
                />
              </div>
            </div>

            {/* Additional Details */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                Additional Details
              </h3>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px', marginBottom: '12px' }}>
                <Input
                  label="Storage Location"
                  value={formData.storageLocation}
                  onChangeText={(value) => setFormData({...formData, storageLocation: value})}
                  placeholder="e.g., Cold Storage A"
                  fullWidth
                />
                <Select
                  label="Packaging Type"
                  value={formData.packagingType}
                  onChange={(e) => setFormData({...formData, packagingType: e.target.value})}
                  options={[
                    { value: '', label: 'Select Packaging' },
                    { value: 'Box', label: 'Box' },
                    { value: 'Bag', label: 'Bag' },
                    { value: 'Crate', label: 'Crate' },
                    { value: 'Loose', label: 'Loose' },
                    { value: 'Bottle', label: 'Bottle' },
                    { value: 'Can', label: 'Can' }
                  ]}
                  fullWidth
                />
              </div>
              <div>
                <label style={{ display: 'block', fontSize: '14px', fontWeight: '500', marginBottom: '6px', color: 'var(--text-primary)' }}>
                  Notes
                </label>
                <textarea 
                  value={formData.notes} 
                  onChange={(e) => setFormData({...formData, notes: e.target.value})}
                  placeholder="Additional notes or remarks"
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
                {editingItem ? 'Update Item' : 'Add Item'}
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
          setDeletingItem(null);
        }}
        onConfirm={confirmDelete}
        title="Delete Inventory Item"
        entityType="Inventory Item"
        entityName={deletingItem ? `${deletingItem.name}${deletingItem.variant ? ' - ' + deletingItem.variant : ''}` : ''}
        warningMessage="This will mark the inventory item as deleted. The item record and its history will be preserved for audit purposes."
      />
    </div>
  );
}

export default Inventory;
