import React, { useState, useEffect } from 'react';
import { reportAPI, customerAPI } from '../services/api';
import { Card, Button, Badge, Input, Select } from '../components/ui';
import { FileText, Printer, Calendar } from 'lucide-react';

function Reports() {
  const [reportType, setReportType] = useState('sales');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [customerId, setCustomerId] = useState('');
  const [customers, setCustomers] = useState([]);
  const [report, setReport] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchCustomers();
  }, []);

  const fetchCustomers = async () => {
    try {
      const res = await customerAPI.getAll();
      const custData = res.data.data || res.data || [];
      setCustomers(custData);
    } catch (err) {
      console.error(err);
    }
  };

  const handleGenerate = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setReport(null);

    try {
      const payload = {
        type: reportType,
        startDate: startDate || undefined,
        endDate: endDate || undefined,
        customerId: reportType === 'customer' ? customerId : undefined
      };
      const res = await reportAPI.generate(payload);
      setReport(res.data);
    } catch (err) {
      setError('Failed to generate report: ' + (err.response?.data || err.message));
    } finally {
      setLoading(false);
    }
  };

  const handlePrint = () => {
    window.print();
  };

  return (
    <div className="container">
      <div className="no-print">
        <div style={{ marginBottom: '24px' }}>
          <h1 style={{ marginBottom: '8px', display: 'flex', alignItems: 'center', gap: '12px' }}>
            <FileText size={28} style={{ color: 'var(--color-primary)' }} />
            Reports & Analytics
          </h1>
        </div>
        
        <Card>
          <h2 style={{ marginBottom: '20px', fontSize: '18px', fontWeight: '600' }}>Generate Report</h2>
          <form onSubmit={handleGenerate}>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
              {/* Report Configuration */}
              <div>
                <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                  Report Configuration
                </h3>
                <Select
                  label="Report Type"
                  value={reportType}
                  onChange={(e) => setReportType(e.target.value)}
                  options={[
                    { value: 'sales', label: '📊 Sales Report' },
                    { value: 'inventory', label: '📦 Inventory Report' },
                    { value: 'customer', label: '👤 Customer Report' }
                  ]}
                  required
                  fullWidth
                />
              </div>
              
              {/* Date Range Section */}
              {reportType !== 'inventory' && (
                <div>
                  <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                    Date Range
                  </h3>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                    <Input
                      label="Start Date"
                      type="date"
                      value={startDate}
                      onChangeText={setStartDate}
                      leftIcon={<Calendar size={18} />}
                      fullWidth
                    />
                    <Input
                      label="End Date"
                      type="date"
                      value={endDate}
                      onChangeText={setEndDate}
                      leftIcon={<Calendar size={18} />}
                      fullWidth
                    />
                  </div>
                  <p style={{ fontSize: '13px', color: 'var(--text-secondary)', marginTop: '8px' }}>
                    Leave blank for all-time data
                  </p>
                </div>
              )}
              
              {/* Customer Selection */}
              {reportType === 'customer' && (
                <div>
                  <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                    Customer Selection
                  </h3>
                  <Select
                    label="Select Customer"
                    value={customerId}
                    onChange={(e) => setCustomerId(e.target.value)}
                    options={[
                      { value: '', label: 'Select Customer' },
                      ...customers.map(c => ({ value: c.id, label: c.name }))
                    ]}
                    required
                    fullWidth
                  />
                </div>
              )}
              
              {/* Generate Button */}
              <div style={{ paddingTop: '4px' }}>
                <Button 
                  type="submit" 
                  variant="primary" 
                  loading={loading}
                  leftIcon={<FileText size={20} />}
                  fullWidth
                  style={{ padding: '14px' }}
                >
                  {loading ? 'Generating Report...' : 'Generate Report'}
                </Button>
              </div>
            </div>
          </form>
        </Card>

        {error && <div className="error" style={{ marginTop: '20px', padding: '16px', backgroundColor: 'var(--color-danger-light)', color: 'var(--color-danger)', borderRadius: '8px' }}>{error}</div>}
      </div>

      {report && (
        <Card style={{ marginTop: '24px' }}>
          <div className="no-print" style={{ marginBottom: '20px' }}>
            <Button 
              onClick={handlePrint} 
              variant="secondary"
              leftIcon={<Printer size={20} />}
            >
              Print Report
            </Button>
          </div>

          {/* Sales Report */}
          {reportType === 'sales' && (
            <>
              <div style={{ textAlign: 'center', marginBottom: '30px' }}>
                <h1>Sales Report</h1>
                <p>Generated: {new Date(report.generatedAt).toLocaleString()}</p>
                {startDate && <p>Period: {startDate} to {endDate || 'Present'}</p>}
              </div>
              
              <div style={{ marginBottom: '30px' }}>
                <h2>Summary</h2>
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: '20px' }}>
                  <div>
                    <strong>Total Sales:</strong> ₹{report.totalSales?.toFixed(2)}
                  </div>
                  <div>
                    <strong>Total Transactions:</strong> {report.totalItems}
                  </div>
                </div>
              </div>

              {report.topItems && report.topItems.length > 0 && (
                <div style={{ marginBottom: '30px' }}>
                  <h2>Top Selling Items</h2>
                  <table className="table">
                    <thead>
                      <tr>
                        <th>Item</th>
                        <th>Quantity Sold</th>
                        <th>Revenue</th>
                      </tr>
                    </thead>
                    <tbody>
                      {report.topItems.map((item, idx) => (
                        <tr key={idx}>
                          <td>{item.itemName}</td>
                          <td>{item.totalQuantity}</td>
                          <td>₹{item.totalRevenue?.toFixed(2)}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}

              <div>
                <h2>All Transactions</h2>
                <table className="table">
                  <thead>
                    <tr>
                      <th>Date</th>
                      <th>Amount</th>
                    </tr>
                  </thead>
                  <tbody>
                    {report.transactions && report.transactions.length > 0 ? (
                      report.transactions.map(tx => (
                        <tr key={tx.id}>
                          <td>{new Date(tx.date).toLocaleDateString()}</td>
                          <td>₹{(tx.totalAmount || 0).toFixed(2)}</td>
                        </tr>
                      ))
                    ) : (
                      <tr>
                        <td colSpan="2" style={{ textAlign: 'center', padding: '20px' }}>No transactions found</td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </div>
            </>
          )}

          {/* Inventory Report */}
          {reportType === 'inventory' && (
            <>
              <div style={{ textAlign: 'center', marginBottom: '30px' }}>
                <h1>Inventory Report</h1>
                <p>Generated: {new Date(report.generatedAt).toLocaleString()}</p>
              </div>
              
              <div style={{ marginBottom: '30px' }}>
                <h2>Summary</h2>
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '20px' }}>
                  <div><strong>Total Items:</strong> {report.totalItems}</div>
                  <div><strong>Expiring Soon:</strong> {report.expiringSoon}</div>
                  <div><strong>Expired:</strong> {report.expired}</div>
                </div>
              </div>

              <table className="table">
                <thead>
                  <tr>
                    <th>Item</th>
                    <th>Lot Number</th>
                    <th>Quantity</th>
                    <th>Expiry Date</th>
                  </tr>
                </thead>
                <tbody>
                  {report.items.map(item => (
                    <tr key={item.id}>
                      <td>{item.name} {item.variant && `(${item.variant})`}</td>
                      <td>{item.lotNumber}</td>
                      <td>{item.quantity} {item.unit}</td>
                      <td>{item.expiryDate}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </>
          )}

          {/* Customer Report */}
          {reportType === 'customer' && (
            <>
              <div style={{ textAlign: 'center', marginBottom: '30px' }}>
                <h1>Customer Financial Report</h1>
                <h2>{report.customerName || 'Unknown Customer'}</h2>
                <p>Generated: {new Date(report.generatedAt).toLocaleString()}</p>
                {startDate && <p>Period: {startDate} to {endDate || 'Present'}</p>}
              </div>
              
              <div style={{ marginBottom: '30px' }}>
                <h2>Summary</h2>
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: '20px' }}>
                  <div><strong>Total Sales:</strong> ₹{(report.totalSales || 0).toFixed(2)}</div>
                  <div><strong>Total Payments:</strong> ₹{(report.totalPayments || 0).toFixed(2)}</div>
                  <div><strong>Outstanding Balance:</strong> ₹{(report.outstandingBalance || 0).toFixed(2)}</div>
                  <div><strong>Crate Balance:</strong> {report.crateBalance || 0} crates</div>
                </div>
              </div>

              <table className="table">
                <thead>
                  <tr>
                    <th>Date</th>
                    <th>Type</th>
                    <th>Amount</th>
                  </tr>
                </thead>
                <tbody>
                  {report.transactions && report.transactions.length > 0 ? (
                    report.transactions.map(tx => (
                      <tr key={tx.id}>
                        <td style={{ color: 'var(--text-secondary)' }}>{new Date(tx.date).toLocaleDateString()}</td>
                        <td>
                          <Badge variant={tx.type === 'sale' ? 'success' : 'info'} size="sm">
                            {tx.type}
                          </Badge>
                        </td>
                        <td style={{ fontWeight: '600', color: 'var(--text-primary)' }}>₹{(tx.totalAmount || 0).toFixed(2)}</td>
                      </tr>
                    ))
                  ) : (
                    <tr>
                      <td colSpan="3" style={{ textAlign: 'center', padding: '20px', color: 'var(--text-secondary)' }}>No transactions found for this period</td>
                    </tr>
                  )}
                </tbody>
              </table>
            </>
          )}
        </Card>
      )}
    </div>
  );
}

export default Reports;
