/**
 * Utility functions to handle API response field name variations
 * Backend returns PascalCase, but we want to normalize to camelCase
 */

/**
 * Normalize customer object to handle both PascalCase and camelCase
 */
export const normalizeCustomer = (customer) => {
  if (!customer) return null;
  
  return {
    id: customer.ID || customer.id,
    name: customer.Name || customer.name,
    email: customer.Email || customer.email,
    address: customer.Address || customer.address,
    contactNumber: customer.ContactNumber || customer.contactNumber,
    alternateContact: customer.AlternateContact || customer.alternateContact,
    whatsappNumber: customer.WhatsappNumber || customer.whatsappNumber,
    photoUrl: customer.PhotoURL || customer.photoUrl,
    businessName: customer.BusinessName || customer.businessName,
    gstin: customer.GSTIN || customer.gstin,
    customerType: customer.CustomerType || customer.customerType,
    aadhaarVerified: customer.AadhaarVerified || customer.aadhaarVerified,
    kycDocumentType: customer.KYCDocumentType || customer.kycDocumentType,
    kycDocumentNumber: customer.KYCDocumentNumber || customer.kycDocumentNumber,
    creditLimit: customer.CreditLimit || customer.creditLimit || 0,
    currentBalance: customer.CurrentBalance || customer.currentBalance || 0,
    paymentTermsDays: customer.PaymentTermsDays || customer.paymentTermsDays || 0,
    interestRate: customer.InterestRate || customer.interestRate || 0,
    status: customer.Status || customer.status,
    notes: customer.Notes || customer.notes,
    tags: customer.Tags || customer.tags || [],
    lastTransactionDate: customer.LastTransactionDate || customer.lastTransactionDate,
    totalPurchases: customer.TotalPurchases || customer.totalPurchases || 0,
    loyaltyPoints: customer.LoyaltyPoints || customer.loyaltyPoints || 0,
    createdAt: customer.CreatedAt || customer.createdAt,
    updatedAt: customer.UpdatedAt || customer.updatedAt,
  };
};

/**
 * Normalize inventory item object
 */
export const normalizeInventory = (item) => {
  if (!item) return null;
  
  return {
    id: item.ID || item.id,
    name: item.Name || item.name,
    variant: item.Variant || item.variant,
    lotNumber: item.LotNumber || item.lotNumber,
    quantity: item.Quantity !== undefined ? item.Quantity : item.quantity,
    unit: item.Unit || item.unit,
    purchaseDate: item.PurchaseDate || item.purchaseDate,
    expiryDate: item.ExpiryDate || item.expiryDate,
    category: item.Category || item.category,
    subCategory: item.SubCategory || item.subCategory,
    costPrice: item.CostPrice || item.costPrice || 0,
    sellingPrice: item.SellingPrice || item.sellingPrice || 0,
    marginPercentage: item.MarginPercentage || item.marginPercentage || 0,
    supplierID: item.SupplierID || item.supplierID,
    supplierName: item.SupplierName || item.supplierName,
    purchaseInvoice: item.PurchaseInvoice || item.purchaseInvoice,
    minStockLevel: item.MinStockLevel || item.minStockLevel || 0,
    reorderPoint: item.ReorderPoint || item.reorderPoint || 0,
    shelfLifeDays: item.ShelfLifeDays || item.shelfLifeDays,
    storageLocation: item.StorageLocation || item.storageLocation,
    barcode: item.Barcode || item.barcode,
    sku: item.SKU || item.sku,
    hsnCode: item.HSNCode || item.hsnCode,
    gstRate: item.GSTRate || item.gstRate || 0,
    status: item.Status || item.status,
    weightPerUnit: item.WeightPerUnit || item.weightPerUnit,
    packagingType: item.PackagingType || item.packagingType,
    imageURL: item.ImageURL || item.imageURL,
    notes: item.Notes || item.notes,
    totalSold: item.TotalSold || item.totalSold || 0,
    totalWasted: item.TotalWasted || item.totalWasted || 0,
    lastRestockDate: item.LastRestockDate || item.lastRestockDate,
    createdAt: item.CreatedAt || item.createdAt,
    updatedAt: item.UpdatedAt || item.updatedAt,
  };
};

/**
 * Normalize transaction object
 */
export const normalizeTransaction = (txn) => {
  if (!txn) return null;
  
  return {
    id: txn.ID || txn.id,
    customerId: txn.CustomerID || txn.customerID || txn.customerId,
    date: txn.Date || txn.date,
    type: txn.Type || txn.type,
    invoiceNumber: txn.InvoiceNumber || txn.invoiceNumber,
    paymentMethod: txn.PaymentMethod || txn.paymentMethod,
    paymentAmount: txn.PaymentAmount || txn.paymentAmount || 0,
    paymentRef: txn.PaymentRef || txn.paymentRef,
    dueDate: txn.DueDate || txn.dueDate,
    isOverdue: txn.IsOverdue || txn.isOverdue || false,
    status: txn.Status || txn.status,
    totalAmount: txn.TotalAmount || txn.totalAmount || 0,
    discountAmount: txn.DiscountAmount || txn.discountAmount || 0,
    taxAmount: txn.TaxAmount || txn.taxAmount || 0,
    balanceAfter: txn.BalanceAfter || txn.balanceAfter || 0,
    saleType: txn.SaleType || txn.saleType,
    receiptSent: txn.ReceiptSent || txn.receiptSent || false,
    deliveryStatus: txn.DeliveryStatus || txn.deliveryStatus,
    deliveryDate: txn.DeliveryDate || txn.deliveryDate,
    deliveryAddress: txn.DeliveryAddress || txn.deliveryAddress,
    notes: txn.Notes || txn.notes,
    details: txn.Details || txn.details,
    createdAt: txn.CreatedAt || txn.createdAt,
  };
};

/**
 * Format currency
 */
export const formatCurrency = (amount) => {
  if (amount === undefined || amount === null || amount === '') return '-';
  const num = typeof amount === 'string' ? parseFloat(amount) : amount;
  if (isNaN(num)) return '-';
  return `₹${num.toLocaleString('en-IN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
};

/**
 * Format date
 */
export const formatDate = (date) => {
  if (!date) return '-';
  try {
    return new Date(date).toLocaleDateString('en-IN', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  } catch (e) {
    return '-';
  }
};

/**
 * Format date with time
 */
export const formatDateTime = (date) => {
  if (!date) return '-';
  try {
    return new Date(date).toLocaleString('en-IN', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  } catch (e) {
    return '-';
  }
};

/**
 * Get status badge variant
 */
export const getStatusVariant = (status) => {
  const statusLower = (status || '').toLowerCase();
  switch (statusLower) {
    case 'active':
    case 'available':
    case 'completed':
    case 'delivered':
    case 'paid':
      return 'success';
    case 'inactive':
    case 'pending':
    case 'partial':
      return 'warning';
    case 'blocked':
    case 'expired':
    case 'damaged':
    case 'overdue':
    case 'cancelled':
      return 'danger';
    case 'low_stock':
    case 'low stock':
      return 'warning';
    case 'out_of_stock':
    case 'out of stock':
      return 'danger';
    default:
      return 'neutral';
  }
};

/**
 * Normalize crate object
 */
export const normalizeCrate = (crate) => {
  if (!crate) return null;
  
  return {
    id: crate.ID || crate.id,
    customerId: crate.CustomerID || crate.customer_id || crate.customerId,
    customerName: crate.CustomerName || crate.customer_name || crate.customerName,
    transactionId: crate.TransactionID || crate.transaction_id || crate.transactionId,
    date: crate.Date || crate.date,
    cratesIssued: crate.CratesIssued !== undefined ? crate.CratesIssued : (crate.crates_issued !== undefined ? crate.crates_issued : crate.cratesIssued),
    cratesReturned: crate.CratesReturned !== undefined ? crate.CratesReturned : (crate.crates_returned !== undefined ? crate.crates_returned : crate.cratesReturned),
    balance: crate.Balance !== undefined ? crate.Balance : crate.balance,
    notes: crate.Notes || crate.notes,
    crateType: crate.CrateType || crate.crate_type || crate.crateType,
    crateValue: crate.CrateValue || crate.crate_value || crate.crateValue || 0,
    updatedAt: crate.UpdatedAt || crate.updated_at || crate.updatedAt,
    updatedBy: crate.UpdatedBy || crate.updated_by || crate.updatedBy,
  };
};

/**
 * Normalize wastage log object
 */
export const normalizeWastage = (wastage) => {
  if (!wastage) return null;
  
  return {
    id: wastage.ID || wastage.id,
    inventoryId: wastage.InventoryID || wastage.inventory_id || wastage.inventoryId,
    itemName: wastage.ItemName || wastage.item_name || wastage.itemName,
    lotNumber: wastage.LotNumber || wastage.lot_number || wastage.lotNumber,
    quantity: wastage.Quantity !== undefined ? wastage.Quantity : wastage.quantity,
    unit: wastage.Unit || wastage.unit,
    reason: wastage.Reason || wastage.reason,
    reasonDetails: wastage.ReasonDetails || wastage.reason_details || wastage.reasonDetails,
    costValue: wastage.CostValue || wastage.cost_value || wastage.costValue || 0,
    expiryDate: wastage.ExpiryDate || wastage.expiry_date || wastage.expiryDate,
    notes: wastage.Notes || wastage.notes || wastage.ReasonDetails || wastage.reason_details,
    recordedBy: wastage.RecordedBy || wastage.recorded_by || wastage.recordedBy || wastage.LoggedBy || wastage.logged_by,
    recordedAt: wastage.RecordedAt || wastage.recorded_at || wastage.recordedAt || wastage.LoggedAt || wastage.logged_at,
    createdAt: wastage.CreatedAt || wastage.created_at || wastage.createdAt,
    updatedAt: wastage.UpdatedAt || wastage.updated_at || wastage.updatedAt,
  };
};

/**
 * Normalize expiry alert object
 */
export const normalizeExpiryAlert = (alert) => {
  if (!alert) return null;
  
  return {
    id: alert.ID || alert.id,
    inventoryItemId: alert.InventoryItemID || alert.inventory_item_id || alert.inventoryItemId,
    inventoryItem: alert.InventoryItem ? normalizeInventory(alert.InventoryItem) : null,
    alertDate: alert.AlertDate || alert.alert_date || alert.alertDate,
    expiryDate: alert.ExpiryDate || alert.expiry_date || alert.expiryDate,
    daysUntilExpiry: alert.DaysUntilExpiry !== undefined ? alert.DaysUntilExpiry : (alert.days_until_expiry !== undefined ? alert.days_until_expiry : alert.daysUntilExpiry),
    acknowledged: alert.Acknowledged !== undefined ? alert.Acknowledged : alert.acknowledged,
    acknowledgedAt: alert.AcknowledgedAt || alert.acknowledged_at || alert.acknowledgedAt,
    acknowledgedBy: alert.AcknowledgedBy || alert.acknowledged_by || alert.acknowledgedBy,
    createdAt: alert.CreatedAt || alert.created_at || alert.createdAt,
    updatedAt: alert.UpdatedAt || alert.updated_at || alert.updatedAt,
  };
};
