// src/pages/DashboardPage.tsx
import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { artifactApi } from '../services/api';
import { Artifact } from '../types/Artifact';
import './DashboardPage.css';

const DashboardPage: React.FC = () => {
  const [artifacts, setArtifacts] = useState<Artifact[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [sortField, setSortField] = useState<keyof Artifact>('uploadTime');
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc');
  const [searchTerm, setSearchTerm] = useState<string>('');

  useEffect(() => {
    fetchArtifacts();
  }, []);

  const fetchArtifacts = async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await artifactApi.getAll();
      setArtifacts(response.data);
    } catch (err) {
      setError('Failed to load artifacts');
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSort = (field: keyof Artifact) => {
    if (field === sortField) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
    } else {
      setSortField(field);
      setSortDirection('asc');
    }
  };

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
  };

  const sortedAndFilteredArtifacts = artifacts
    .filter(artifact =>
      artifact.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      artifact.description.toLowerCase().includes(searchTerm.toLowerCase())
    )
    .sort((a, b) => {
      if (a[sortField] < b[sortField]) return sortDirection === 'asc' ? -1 : 1;
      if (a[sortField] > b[sortField]) return sortDirection === 'asc' ? 1 : -1;
      return 0;
    });

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  return (
    <div className="dashboard">
      <div className="dashboard-header">
        <h1>Artifacts</h1>
        <div className="actions">
          <input
            type="text"
            placeholder="Search artifacts..."
            value={searchTerm}
            onChange={handleSearch}
            className="search-input"
          />
          <Link to="/upload" className="btn btn-primary">
            Upload New Artifact
          </Link>
        </div>
      </div>

      {isLoading ? (
        <div className="loading">Loading artifacts...</div>
      ) : error ? (
        <div className="error">{error}</div>
      ) : sortedAndFilteredArtifacts.length === 0 ? (
        <div className="empty-state">
          <p>No artifacts found. Upload your first artifact to get started.</p>
          <Link to="/upload" className="btn btn-primary">
            Upload Artifact
          </Link>
        </div>
      ) : (
        <table className="artifacts-table">
          <thead>
            <tr>
              <th onClick={() => handleSort('name')}>
                Name {sortField === 'name' && (sortDirection === 'asc' ? '↑' : '↓')}
              </th>
              <th onClick={() => handleSort('fileSize')}>
                Size {sortField === 'fileSize' && (sortDirection === 'asc' ? '↑' : '↓')}
              </th>
              <th onClick={() => handleSort('uploadTime')}>
                Upload Date {sortField === 'uploadTime' && (sortDirection === 'asc' ? '↑' : '↓')}
              </th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {sortedAndFilteredArtifacts.map((artifact) => (
              <tr key={artifact.id}>
                <td>
                  <Link to={`/artifacts/${artifact.id}`}>
                    {artifact.name}
                  </Link>
                </td>
                <td>{formatFileSize(artifact.fileSize)}</td>
                <td>{new Date(artifact.uploadTime).toLocaleString()}</td>
                <td className="actions-cell">
                  <a 
                    href={artifactApi.getDownloadUrl(artifact.id)}
                    className="btn btn-sm btn-secondary"
                    download
                  >
                    Download
                  </a>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
};

export default DashboardPage;