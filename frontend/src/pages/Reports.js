import React, { useState, useEffect } from 'react';
import { reportAPI, customerAPI } from '../services/api';

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
      setCustomers(res.data);
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
        <h1>Reports & Analytics</h1>
        
        <div className="card">
          <h2>Generate Report</h2>
          <form onSubmit={handleGenerate}>
            <div className="form-group">
              <label>Report Type*</label>
              <select value={reportType} onChange={(e) => setReportType(e.target.value)} required>
                <option value="sales">Sales Report</option>
                <option value="inventory">Inventory Report</option>
                <option value="customer">Customer Report</option>
              </select>
            </div>
            
            {reportType !== 'inventory' && (
              <>
                <div className="form-group">
                  <label>Start Date</label>
                  <input
                    type="date"
                    value={startDate}
                    onChange={(e) => setStartDate(e.target.value)}
                  />
                </div>
                <div className="form-group">
                  <label>End Date</label>
                  <input
                    type="date"
                    value={endDate}
                    onChange={(e) => setEndDate(e.target.value)}
                  />
                </div>
              </>
            )}
            
            {reportType === 'customer' && (
              <div className="form-group">
                <label>Customer*</label>
                <select value={customerId} onChange={(e) => setCustomerId(e.target.value)} required>
                  <option value="">Select Customer</option>
                  {customers.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
                </select>
              </div>
            )}
            
            <button type="submit" className="btn btn-primary" disabled={loading}>
              {loading ? 'Generating...' : 'Generate Report'}
            </button>
          </form>
        </div>

        {error && <div className="error">{error}</div>}
      </div>

      {report && (
        <div className="card">
          <div className="no-print" style={{ marginBottom: '20px' }}>
            <button onClick={handlePrint} className="btn btn-secondary">🖨️ Print Report</button>
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
                    {report.transactions.map(tx => (
                      <tr key={tx.id}>
                        <td>{new Date(tx.date).toLocaleDateString()}</td>
                        <td>₹{tx.totalAmount?.toFixed(2)}</td>
                      </tr>
                    ))}
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
                <h2>{report.customerName}</h2>
                <p>Generated: {new Date(report.generatedAt).toLocaleString()}</p>
                {startDate && <p>Period: {startDate} to {endDate || 'Present'}</p>}
              </div>
              
              <div style={{ marginBottom: '30px' }}>
                <h2>Summary</h2>
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: '20px' }}>
                  <div><strong>Total Sales:</strong> ₹{report.totalSales?.toFixed(2)}</div>
                  <div><strong>Total Payments:</strong> ₹{report.totalPayments?.toFixed(2)}</div>
                  <div><strong>Outstanding Balance:</strong> ₹{report.outstandingBalance?.toFixed(2)}</div>
                  <div><strong>Crate Balance:</strong> {report.crateBalance} crates</div>
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
                  {report.transactions.map(tx => (
                    <tr key={tx.id}>
                      <td>{new Date(tx.date).toLocaleDateString()}</td>
                      <td>{tx.type}</td>
                      <td>₹{tx.totalAmount?.toFixed(2)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </>
          )}
        </div>
      )}
    </div>
  );
}

export default Reports;
