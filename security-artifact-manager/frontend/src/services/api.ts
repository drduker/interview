// src/services/api.ts
import axios from 'axios';

const api = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Handle 401 Unauthorized responses
    if (error.response && error.response.status === 401) {
      // Remove token from localStorage
      localStorage.removeItem('token');
      
      // Redirect to login page if not already there
      if (!window.location.pathname.includes('/login')) {
        window.location.href = '/login';
      }
    }
    
    return Promise.reject(error);
  }
);

// Artifact-related API calls
export const artifactApi = {
  // Get all artifacts
  getAll: () => api.get('/api/artifacts'),
  
  // Get artifact by ID
  getById: (id: string) => api.get(`/api/artifacts/${id}`),
  
  // Upload artifact
  upload: (formData: FormData) => api.post('/api/artifacts', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  }),
  
  // Delete artifact
  delete: (id: string) => api.delete(`/api/artifacts/${id}`),
  
  // Download artifact
  getDownloadUrl: (id: string) => `${api.defaults.baseURL}/api/artifacts/${id}/download`,
};

export default api;