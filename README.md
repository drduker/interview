# Secure Artifact Manager

A secure system for storing, validating, and distributing binary artifacts with proper security controls.

## Features

- Secure artifact storage with metadata tracking
- User authentication and authorization
- Artifact upload and download functionality
- Basic vulnerability scanning
- SHA256 hash verification
- Containerized deployment with Kubernetes support

## Architecture

The application consists of:

1. **Frontend**: React/TypeScript web application
2. **Backend API**: Go REST API service
3. **Database**: PostgreSQL for storing metadata and user information
4. **Storage**: Persistent storage for the actual binary artifacts

## Prerequisites

- Docker and Docker Compose for local development
- Kubernetes cluster (Kind, Minikube, or similar) for deployment
- kubectl CLI tool

## Local Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/secure-artifact-manager.git
   cd secure-artifact-manager
   ```

2. Start the development environment:
   ```bash
   docker-compose up -d
   ```

3. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

## Kubernetes Deployment

### Setting up Kind cluster (Local Kubernetes)

1. Install Kind if not already installed:
   ```bash
   # On macOS with Homebrew
   brew install kind
   
   # On Linux
   curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.14.0/kind-linux-amd64
   chmod +x ./kind
   mv ./kind /usr/local/bin/kind
   ```

2. Create a Kind cluster:
   ```bash
   kind create cluster --name artifact-manager
   ```

3. Deploy the application:
   ```bash
   # Apply ConfigMaps and Secrets
   kubectl apply -f kubernetes/configmap-secrets.yaml
   
   # Apply Storage configurations
   kubectl apply -f kubernetes/storage.yaml
   
   # Deploy Database
   kubectl apply -f kubernetes/db-deployment.yaml
   kubectl apply -f kubernetes/services.yaml
   
   # Wait for Database to be ready
   kubectl wait --for=condition=Ready pods -l app=artifact-manager,component=database
   
   # Deploy Backend
   kubectl apply -f kubernetes/backend-deployment.yaml
   
   # Deploy Frontend
   kubectl apply -f kubernetes/frontend-deployment.yaml
   
   # Configure Ingress
   kubectl apply -f kubernetes/ingress.yaml
   ```

4. Add the following to your `/etc/hosts` file:
   ```
   127.0.0.1 artifact-manager.local
   ```

5. Access the application at https://artifact-manager.local

## Security Features

- HTTPS/TLS encryption for all traffic
- JWT-based authentication with secure token handling
- Password hashing with bcrypt
- Artifact validation with SHA256 hashing
- Basic vulnerability scanning
- Container security best practices
- RBAC for authorization control
- Kubernetes security contexts

## Usage

### User Registration and Login

1. Create an admin user (first time setup):
   ```bash
   # Using the provided script
   ./scripts/create-admin.sh
   ```

2. Login with the admin credentials at the web interface.

### Artifact Management

1. **Upload Artifacts**: Use the web interface or API to upload binary artifacts.
2. **View Artifacts**: Browse and search for artifacts in the web interface.
3. **Download Artifacts**: Download artifacts with verified integrity.
4. **Scan Artifacts**: Artifacts are automatically scanned for vulnerabilities.

## API Documentation

### Authentication

```
POST /api/auth/login
POST /api/auth/logout
```

### Artifact Operations

```
GET /api/artifacts - List all artifacts
POST /api/artifacts - Upload a new artifact
GET /api/artifacts/{id} - Get artifact details
DELETE /api/artifacts/{id} - Delete an artifact
GET /api/artifacts/{id}/download - Download artifact
```

## License

[MIT License](LICENSE)