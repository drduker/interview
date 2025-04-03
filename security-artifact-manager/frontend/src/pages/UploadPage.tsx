// src/pages/UploadPage.tsx
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { artifactApi } from '../services/api';
import './UploadPage.css';

const UploadPage: React.FC = () => {
  const navigate = useNavigate();
  const [file, setFile] = useState<File | null>(null);
  const [name, setName] = useState<string>('');
  const [description, setDescription] = useState<string>('');
  const [isUploading, setIsUploading] = useState<boolean>(false);
  const [uploadProgress, setUploadProgress] = useState<number>(0);
  const [error, setError] = useState<string | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const uploadFile = e.target.files[0];
      setFile(uploadFile);
      
      // If no name set, use the file name without extension
      if (!name) {
        const fileName = uploadFile.name.split('.');
        fileName.pop(); // Remove extension
        setName(fileName.join('.'));
      }
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!file) {
      setError('Please select a file to upload');
      return;
    }
    
    if (!name.trim()) {
      setError('Please provide a name for the artifact');
      return;
    }
    
    setIsUploading(true);
    setError(null);
    
    const formData = new FormData();
    formData.append('file', file);
    formData.append('name', name);
    formData.append('description', description);
    
    try {
      // Create a custom axios request with upload progress tracking
      await artifactApi.upload(formData);
      
      // Redirect to dashboard on success
      navigate('/');
    } catch (err) {
      console.error(err);
      setError('Failed to upload artifact. Please try again.');
      setIsUploading(false);
    }
  };

  return (
    <div className="upload-page">
      <h1>Upload Artifact</h1>
      
      <form onSubmit={handleSubmit} className="upload-form">
        {error && <div className="error-message">{error}</div>}
        
        <div className="form-group">
          <label htmlFor="file">Artifact File*</label>
          <input
            type="file"
            id="file"
            onChange={handleFileChange}
            required
            disabled={isUploading}
          />
          {file && (
            <div className="file-info">
              <strong>Selected file:</strong> {file.name} ({(file.size / (1024 * 1024)).toFixed(2)} MB)
            </div>
          )}
        </div>
        
        <div className="form-group">
          <label htmlFor="name">Name*</label>
          <input
            type="text"
            id="name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
            disabled={isUploading}
            placeholder="Enter artifact name"
          />
        </div>
        
        <div className="form-group">
          <label htmlFor="description">Description</label>
          <textarea
            id="description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            disabled={isUploading}
            placeholder="Enter artifact description (optional)"
            rows={4}
          />
        </div>
        
        {isUploading && (
          <div className="upload-progress">
            <div 
              className="progress-bar" 
              style={{ width: `${uploadProgress}%` }}
            ></div>
            <div className="progress-label">{uploadProgress}%</div>
          </div>
        )}
        
        <div className="form-actions">
          <button
            type="button"
            className="btn btn-secondary"
            onClick={() => navigate('/')}
            disabled={isUploading}
          >
            Cancel
          </button>
          <button
            type="submit"
            className="btn btn-primary"
            disabled={isUploading || !file}
          >
            {isUploading ? 'Uploading...' : 'Upload Artifact'}
          </button>
        </div>
      </form>
    </div>
  );
};

export default UploadPage;