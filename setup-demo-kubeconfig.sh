#!/bin/bash

# Set up directories
echo "Setting up directories for testing..."
KUBE_DIR="$HOME/.kube"
BACKUP_DIR="$KUBE_DIR/backup_$(date +%Y%m%d%H%M%S)"

# Create backup of current config if it exists
if [ -f "$KUBE_DIR/config" ]; then
  echo "Backing up existing kubeconfig to $BACKUP_DIR/config"
  mkdir -p "$BACKUP_DIR"
  cp "$KUBE_DIR/config" "$BACKUP_DIR/config"
fi

# Copy the sample kubeconfig
echo "Installing sample kubeconfig..."
cp "$(pwd)/sample-kubeconfig.yaml" "$KUBE_DIR/config"

# Set permissions
chmod 600 "$KUBE_DIR/config"
echo "Sample kubeconfig installed to $KUBE_DIR/config"
echo ""
echo "Available contexts in the sample kubeconfig:"
echo "--------------------------------------------"
kubectl config get-contexts
echo ""
echo "To restore your original kubeconfig, run:"
echo "cp \"$BACKUP_DIR/config\" \"$KUBE_DIR/config\""
echo ""
echo "Ready to take screenshots! Try running 'kontext' commands like:"
echo "kontext"
echo "kontext switch"
echo "kontext namespace"
echo "kontext switch production-cluster -n"
echo "kontext list"
echo "kontext current"
