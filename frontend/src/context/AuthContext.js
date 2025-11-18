import React, { createContext, useContext, useState, useEffect } from 'react';
import { authAPI } from '../services/api';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(localStorage.getItem('token'));
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (token) {
      localStorage.setItem('token', token);
    } else {
      localStorage.removeItem('token');
    }
  }, [token]);

  const login = async (email, password) => {
    setLoading(true);
    try {
      const response = await authAPI.login({ email, password });
      const data = response.data.data || response.data;
      const newToken = data.token;
      setToken(newToken);
      setUser({ email });
      return { success: true };
    } catch (error) {
      const errorMessage = error.response?.data?.error?.message || 
                          error.response?.data?.message || 
                          error.message ||
                          'Login failed. Please try again.';
      return { 
        success: false, 
        error: errorMessage
      };
    } finally {
      setLoading(false);
    }
  };

  const register = async (name, email, password) => {
    setLoading(true);
    try {
      await authAPI.register({ name, email, password });
      return { success: true };
    } catch (error) {
      const errorMessage = error.response?.data?.error?.message || 
                          error.response?.data?.message || 
                          error.message ||
                          'Registration failed. Please try again.';
      return { 
        success: false, 
        error: errorMessage
      };
    } finally {
      setLoading(false);
    }
  };

  const logout = () => {
    setToken(null);
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, token, login, register, logout, loading }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};
