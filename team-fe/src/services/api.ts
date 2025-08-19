import axios from 'axios';

// Base URLs for your services
export const ASSET_SERVICE_URL = 'http://localhost:7070/api/v1';
export const USER_SERVICE_URL = 'http://localhost:8080';

// Create axios instances for each service
export const assetApi = axios.create({
  baseURL: ASSET_SERVICE_URL,
});

export const userApi = axios.create({
  baseURL: USER_SERVICE_URL,
});

// Add auth token to requests
assetApi.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

userApi.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
