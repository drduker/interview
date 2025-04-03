#!/bin/bash
set -e

# Define colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Setting up Kind cluster for Secure Artifact Manager...${NC}"

# Check if kind is installed
if ! command -v kind &> /dev/null; then
    echo -e "${RED}Error: kind is not installed. Please install it first:${NC}"
    echo "https://kind.sigs.k8s.io/docs/user/quick-start/#installation"
    exit 1
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo -e "${RED}Error: kubectl is not installed. Please install it first:${NC}"
    echo "https://kubernetes.io/docs/tasks/tools/"
    exit 1
fi

# Check if docker is running
if ! docker info &> /dev/null; then
    echo -e "${RED}Error: Docker is not running. Please start Docker and try again.${NC}"
    exit 1
fi

# Create kind cluster with configuration for ingress
echo -e "${GREEN}Creating Kind cluster...${NC}"
cat <<EOF | kind create cluster --name artifact-manager --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF

# Install NGINX Ingress Controller
echo -e "${GREEN}Installing NGINX Ingress Controller...${NC}"
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

# Wait for ingress controller to be ready
echo -e "${GREEN}Waiting for NGINX Ingress Controller to be ready...${NC}"
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s

# Build and load Docker images into Kind cluster
echo -e "${GREEN}Building and loading Docker images...${NC}"
docker build -t artifact-manager-backend:latest ./backend
docker build -t artifact-manager-frontend:latest ./frontend

kind load docker-image artifact-manager-backend:latest --name artifact-manager
kind load docker-image artifact-manager-frontend:latest --name artifact-manager

# Apply Kubernetes configurations
echo -e "${GREEN}Applying Kubernetes configurations...${NC}"
kubectl apply -f kubernetes/configmap-secrets.yaml
kubectl apply -f kubernetes/storage.yaml

echo -e "${GREEN}Deploying PostgreSQL database...${NC}"
kubectl apply -f kubernetes/db-deployment.yaml
kubectl apply -f kubernetes/services.yaml

echo -e "${GREEN}Waiting for database to be ready...${NC}"
kubectl wait --for=condition=Ready pods -l app=artifact-manager,component=database --timeout=120s

echo -e "${GREEN}Deploying backend and frontend...${NC}"
kubectl apply -f kubernetes/backend-deployment.yaml
kubectl apply -f kubernetes/frontend-deployment.yaml

echo -e "${GREEN}Configuring Ingress...${NC}"
kubectl apply -f kubernetes/ingress.yaml

echo -e "${YELLOW}================================================${NC}"
echo -e "${GREEN}Setup complete!${NC}"
echo -e "${YELLOW}================================================${NC}"
echo -e "Add the following to your /etc/hosts file:"
echo -e "${GREEN}127.0.0.1 artifact-manager.local${NC}"
echo -e "Then access the application at: ${GREEN}http://artifact-manager.local${NC}"
echo -e "${YELLOW}================================================${NC}"