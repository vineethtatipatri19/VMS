import axios from 'axios';

// Detect if running in development or production
const isDevelopment = process.env.NODE_ENV === 'development';
const isCodespaces = window.location.hostname.includes('github.dev') || 
                     window.location.hostname.includes('githubpreview.dev') ||
                     window.location.hostname.includes('app.github.dev');

// In Codespaces or production, use the same host with port 8080
// In local development, use localhost:8080
let API_URL;
if (isCodespaces) {
  // Replace port 3000 with 8080 in the current URL
  const url = new URL(window.location.href);
  const hostname = url.hostname;
  // GitHub Codespaces URL pattern
  const baseUrl = hostname.replace(/-3000\./, '-8080.');
  API_URL = `${url.protocol}//${baseUrl}/api/v1`;
} else {
  API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';
}

console.log('API URL:', API_URL);

// Create axios instance with default config
const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Auth APIs
export const authAPI = {
  login: (credentials) => api.post('/login', credentials),
  register: (userData) => api.post('/register', userData),
};

// Dashboard APIs
export const dashboardAPI = {
  getStats: () => api.get('/dashboard'),
  getActivity: () => api.get('/dashboard/activity'),
};

// Inventory APIs
export const inventoryAPI = {
  getAll: (params) => api.get('/inventory', { params }),
  getById: (id) => api.get(`/inventory/${id}`),
  create: (data) => api.post('/inventory', data),
  update: (id, data) => api.put(`/inventory/${id}`, data),
  delete: (id, data) => api.delete(`/inventory/${id}`, { data }),
};

// Customer APIs
export const customerAPI = {
  getAll: () => api.get('/customers'),
  getById: (id) => api.get(`/customers/${id}`),
  create: (data) => api.post('/customers', data),
  update: (id, data) => api.put(`/customers/${id}`, data),
  delete: (id, data) => api.delete(`/customers/${id}`, { data }),
};

// Transaction APIs
export const transactionAPI = {
  getAll: (params) => api.get('/transactions', { params }),
  getById: (id) => api.get(`/transactions/${id}`),
  create: (data) => api.post('/transactions', data),
  update: (id, data) => api.put(`/transactions/${id}`, data),
  delete: (id, data) => api.delete(`/transactions/${id}`, { data }),
};

// Crate APIs
export const crateAPI = {
  getAll: (params) => api.get('/crates', { params }),
  create: (data) => api.post('/crates', data),
  update: (id, data) => api.put(`/crates/${id}`, data),
  delete: (id, data) => api.delete(`/crates/${id}`, { data }),
  getBalance: (customerId) => api.get(`/crates/balance/${customerId}`),
};

// Wastage APIs
export const wastageAPI = {
  getAll: () => api.get('/wastage'),
  create: (data) => api.post('/wastage', data),
  update: (id, data) => api.put(`/wastage/${id}`, data),
  delete: (id, data) => api.delete(`/wastage/${id}`, { data }),
};

// Expiry Alerts APIs
export const expiryAlertsAPI = {
  getAll: (params) => api.get('/expiry-alerts', { params }),
  update: (id, data) => api.put(`/expiry-alerts/${id}`, data),
  delete: (id, data) => api.delete(`/expiry-alerts/${id}`, { data }),
};

// Forecast APIs
export const forecastAPI = {
  generate: (data) => api.post('/forecast', data),
};

// Report APIs
export const reportAPI = {
  generate: (data) => api.post('/reports/generate', data),
};

export default api;
