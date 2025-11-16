import React, { useState } from 'react';
import { forecastAPI } from '../services/api';
import { Card, Button, Badge, Input } from '../components/ui';
import { TrendingUp, Sparkles, Calendar } from 'lucide-react';

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
      <div style={{ marginBottom: '24px' }}>
        <h1 style={{ marginBottom: '8px', display: 'flex', alignItems: 'center', gap: '12px' }}>
          <Sparkles size={28} style={{ color: 'var(--color-primary)' }} />
          AI-Powered Demand Forecasting
        </h1>
        <p style={{ color: 'var(--text-secondary)', fontSize: '14px' }}>
          Generate demand forecasts using Google Gemini AI based on historical sales data
        </p>
      </div>

      <Card>
        <h2 style={{ marginBottom: '20px', fontSize: '18px', fontWeight: '600' }}>Generate Forecast</h2>
        <form onSubmit={handleGenerate}>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
            {/* Forecast Parameters */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '16px' }}>
                Forecast Parameters
              </h3>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                <Input
                  label="Item Name"
                  type="text"
                  value={itemName}
                  onChangeText={setItemName}
                  placeholder="e.g., Tomato, Potato, Onion"
                  required
                  fullWidth
                />
                <Input
                  label="Forecast Period (days)"
                  type="number"
                  value={days}
                  onChangeText={setDays}
                  placeholder="7"
                  min="1"
                  max="30"
                  leftIcon={<Calendar size={18} />}
                  helperText="Enter number of days to forecast (1-30)"
                  required
                  fullWidth
                />
              </div>
            </div>

            {/* Generate Button */}
            <div style={{ paddingTop: '4px' }}>
              <Button 
                type="submit" 
                variant="primary" 
                loading={loading}
                leftIcon={<TrendingUp size={20} />}
                fullWidth
                style={{ padding: '14px' }}
              >
                {loading ? 'Generating Forecast...' : 'Generate AI Forecast'}
              </Button>
            </div>
          </div>
        </form>
      </Card>

      {error && <div className="error" style={{ marginTop: '20px', padding: '16px', backgroundColor: 'var(--color-danger-light)', color: 'var(--color-danger)', borderRadius: '8px' }}>{error}</div>}

      {forecast && (
        <Card style={{ marginTop: '24px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px', flexWrap: 'wrap', gap: '12px' }}>
            <h2 style={{ fontSize: '18px', fontWeight: '600' }}>Forecast Results for {forecast.itemName}</h2>
            <Badge variant="info" size="md">
              Confidence: {forecast.confidence}
            </Badge>
          </div>
          
          <p style={{ color: 'var(--text-secondary)', fontSize: '14px', marginBottom: '20px' }}>
            Generated at: {new Date(forecast.generatedAt).toLocaleString()}
          </p>

          <div style={{ overflowX: 'auto' }}>
            <table className="table">
              <thead>
                <tr>
                  <th>Date</th>
                  <th style={{ textAlign: 'right' }}>Predicted Demand</th>
                  <th>Unit</th>
                </tr>
              </thead>
              <tbody>
                {forecast.predictions.map((pred, idx) => (
                  <tr key={idx}>
                    <td style={{ color: 'var(--text-secondary)' }}>{pred.date}</td>
                    <td style={{ textAlign: 'right', fontWeight: '600', color: 'var(--text-primary)' }}>{pred.quantity.toFixed(2)}</td>
                    <td style={{ color: 'var(--text-secondary)' }}>{pred.unit}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <div style={{ marginTop: '20px', padding: '16px', backgroundColor: 'var(--bg-secondary)', borderRadius: '8px' }}>
            <h3 style={{ marginTop: 0, fontSize: '16px', fontWeight: '600', color: 'var(--text-primary)' }}>Summary</h3>
            <p style={{ color: 'var(--text-primary)' }}>
              <strong>Average Daily Demand:</strong>{' '}
              {(forecast.predictions.reduce((sum, p) => sum + p.quantity, 0) / forecast.predictions.length).toFixed(2)}{' '}
              {forecast.predictions[0]?.unit}
            </p>
            <p style={{ color: 'var(--text-primary)' }}>
              <strong>Total Forecast:</strong>{' '}
              {forecast.predictions.reduce((sum, p) => sum + p.quantity, 0).toFixed(2)}{' '}
              {forecast.predictions[0]?.unit}
            </p>
          </div>
        </Card>
      )}

      <Card style={{ marginTop: '24px' }}>
        <h3 style={{ fontSize: '16px', fontWeight: '600', marginBottom: '12px', color: 'var(--text-primary)' }}>About AI Forecasting</h3>
        <p style={{ color: 'var(--text-primary)' }}>
          This feature uses Google Gemini AI to analyze your historical sales data and predict future demand.
          The AI considers:
        </p>
        <ul style={{ color: 'var(--text-primary)' }}>
          <li>Historical sales trends and patterns</li>
          <li>Seasonal variations</li>
          <li>Day of week effects</li>
          <li>Typical demand patterns for perishable goods</li>
        </ul>
        <p style={{ fontSize: '14px', color: 'var(--text-secondary)', marginTop: '15px' }}>
          Note: To enable AI forecasting, set the GEMINI_API_KEY environment variable. 
          Without an API key, the system will generate simple stub forecasts.
        </p>
      </Card>
    </div>
  );
}

export default Forecasting;
