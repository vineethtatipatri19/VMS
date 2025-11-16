import React, { useState } from 'react';
import { forecastAPI } from '../services/api';

function Forecasting() {
  const [itemName, setItemName] = useState('');
  const [days, setDays] = useState(7);
  const [loading, setLoading] = useState(false);
  const [forecast, setForecast] = useState(null);
  const [error, setError] = useState('');

  const handleGenerate = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setForecast(null);
    
    try {
      const res = await forecastAPI.generate({
        itemName,
        days: parseInt(days),
        historical: true
      });
      setForecast(res.data);
    } catch (err) {
      setError('Failed to generate forecast: ' + (err.response?.data || err.message));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container">
      <h1>AI-Powered Demand Forecasting</h1>
      <p style={{ color: '#666', marginBottom: '20px' }}>
        Generate demand forecasts using Google Gemini AI based on historical sales data
      </p>

      <div className="card">
        <h2>Generate Forecast</h2>
        <form onSubmit={handleGenerate}>
          <div className="form-group">
            <label>Item Name*</label>
            <input
              type="text"
              value={itemName}
              onChange={(e) => setItemName(e.target.value)}
              placeholder="e.g., Tomato, Potato"
              required
            />
          </div>
          <div className="form-group">
            <label>Forecast Period (days)*</label>
            <input
              type="number"
              value={days}
              onChange={(e) => setDays(e.target.value)}
              min="1"
              max="30"
              required
            />
          </div>
          <button type="submit" className="btn btn-primary" disabled={loading}>
            {loading ? 'Generating Forecast...' : 'Generate Forecast'}
          </button>
        </form>
      </div>

      {error && <div className="error">{error}</div>}

      {forecast && (
        <div className="card">
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
            <h2>Forecast Results for {forecast.itemName}</h2>
            <span className="badge badge-info">
              Confidence: {forecast.confidence}
            </span>
          </div>
          
          <p style={{ color: '#666', fontSize: '14px', marginBottom: '20px' }}>
            Generated at: {new Date(forecast.generatedAt).toLocaleString()}
          </p>

          <table className="table">
            <thead>
              <tr>
                <th>Date</th>
                <th>Predicted Demand</th>
                <th>Unit</th>
              </tr>
            </thead>
            <tbody>
              {forecast.predictions.map((pred, idx) => (
                <tr key={idx}>
                  <td>{pred.date}</td>
                  <td>{pred.quantity.toFixed(2)}</td>
                  <td>{pred.unit}</td>
                </tr>
              ))}
            </tbody>
          </table>

          <div style={{ marginTop: '20px', padding: '15px', backgroundColor: '#f8f9fa', borderRadius: '4px' }}>
            <h3 style={{ marginTop: 0 }}>Summary</h3>
            <p>
              <strong>Average Daily Demand:</strong>{' '}
              {(forecast.predictions.reduce((sum, p) => sum + p.quantity, 0) / forecast.predictions.length).toFixed(2)}{' '}
              {forecast.predictions[0]?.unit}
            </p>
            <p>
              <strong>Total Forecast:</strong>{' '}
              {forecast.predictions.reduce((sum, p) => sum + p.quantity, 0).toFixed(2)}{' '}
              {forecast.predictions[0]?.unit}
            </p>
          </div>
        </div>
      )}

      <div className="card">
        <h3>About AI Forecasting</h3>
        <p>
          This feature uses Google Gemini AI to analyze your historical sales data and predict future demand.
          The AI considers:
        </p>
        <ul>
          <li>Historical sales trends and patterns</li>
          <li>Seasonal variations</li>
          <li>Day of week effects</li>
          <li>Typical demand patterns for perishable goods</li>
        </ul>
        <p style={{ fontSize: '14px', color: '#666', marginTop: '15px' }}>
          Note: To enable AI forecasting, set the GEMINI_API_KEY environment variable. 
          Without an API key, the system will generate simple stub forecasts.
        </p>
      </div>
    </div>
  );
}

export default Forecasting;
