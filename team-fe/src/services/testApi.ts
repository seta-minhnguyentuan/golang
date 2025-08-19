// Test utility to check API authentication
import { assetApi } from './api';
import { AxiosError } from 'axios';

export const testAssetAPI = async () => {
  console.log('=== Asset API Authentication Test ===');
  
  // Check if token exists
  const token = localStorage.getItem('token');
  const user = localStorage.getItem('user');
  
  console.log('Token exists:', !!token);
  console.log('User exists:', !!user);
  
  if (token) {
    console.log('Token preview:', token.substring(0, 20) + '...');
  }
  
  if (user) {
    console.log('User:', JSON.parse(user));
  }
  
  // Test health check (if available)
  try {
    const healthResponse = await assetApi.get('/health');
    console.log('Health check:', healthResponse.data);
  } catch (error) {
    const axiosError = error as AxiosError;
    console.log('Health check failed:', axiosError.response?.status, axiosError.response?.data);
  }
  
  // Test folders endpoint
  try {
    const foldersResponse = await assetApi.get('/folders');
    console.log('Folders request successful:', foldersResponse.data);
    return foldersResponse.data;
  } catch (error) {
    const axiosError = error as AxiosError;
    console.error('Folders request failed:', {
      status: axiosError.response?.status,
      statusText: axiosError.response?.statusText,
      data: axiosError.response?.data,
      headers: axiosError.response?.headers
    });
    return null;
  }
};

// Test authentication flow
export const testAuthFlow = async (email: string, password: string) => {
  console.log('=== Testing Complete Auth Flow ===');
  
  try {
    // Import userService dynamically to avoid circular dependencies
    const { userService } = await import('./userService');
    
    // Login first
    console.log('Attempting login...');
    const authPayload = await userService.login({ email, password });
    console.log('Login successful:', authPayload.user);
    
    // Now test asset API
    const folders = await testAssetAPI();
    return { success: true, folders };
    
  } catch (error) {
    console.error('Auth flow failed:', error);
    return { success: false, error };
  }
};
