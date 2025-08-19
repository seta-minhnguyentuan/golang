import axios from 'axios';

// Base URLs for your services (updated to match your actual services)
export const USER_SERVICE_URL = 'http://localhost:8080';  // User Service port
export const ASSET_SERVICE_URL = 'http://localhost:7070/api/v1';  // Asset Service (assuming same port for now)

// Create axios instances for each service
export const userApi = axios.create({
  baseURL: USER_SERVICE_URL,
});

export const assetApi = axios.create({
  baseURL: ASSET_SERVICE_URL,
});

// Add auth token to requests for authenticated endpoints
userApi.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token && config.url?.includes('/teams')) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

assetApi.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    // All asset API endpoints require authentication
    config.headers.Authorization = `Bearer ${token}`;
    console.log('Asset API Request:', config.method?.toUpperCase(), config.url, 'with token:', token.substring(0, 20) + '...');
  } else {
    console.warn('Asset API Request without token:', config.method?.toUpperCase(), config.url);
  }
  
  // Log the full request URL for CORS debugging
  const fullUrl = (config.baseURL || '') + (config.url || '');
  console.log('Full request URL:', fullUrl);
  console.log('Current origin:', window.location.origin);
  
  return config;
});

// Error handling interceptors
userApi.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('User API Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

assetApi.interceptors.response.use(
  (response) => {
    console.log('Asset API Success:', response.config.method?.toUpperCase(), response.config.url, response.status);
    return response;
  },
  (error) => {
    console.error('Asset API Error:', {
      method: error.config?.method?.toUpperCase(),
      url: error.config?.url,
      status: error.response?.status,
      statusText: error.response?.statusText,
      data: error.response?.data,
      message: error.message
    });
    
    if (error.response?.status === 401) {
      console.error('Authentication failed. Please check if you are logged in and have a valid token.');
      // Optionally redirect to login page or clear invalid token
      // localStorage.removeItem('token');
      // localStorage.removeItem('user');
    }
    
    return Promise.reject(error);
  }
);
