#!/bin/bash
# Script to set up the sample kubeconfig for kontext demo
#
# This script will create a backup of your existing kubeconfig file (if any)
# and temporarily replace it with our sample version for demonstration purposes
#
# Usage: ./setup-demo-kubeconfig.sh [restore]
#   Without arguments: Set up the demo kubeconfig
#   With 'restore': Restore your original kubeconfig

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Set up directories
KUBE_DIR="$HOME/.kube"
BACKUP_DIR="$KUBE_DIR/backup_$(date +%Y%m%d%H%M%S)"
LAST_BACKUP_FILE="$KUBE_DIR/.kontext_last_backup"

# Check if we need to restore from backup
if [ "$1" == "restore" ]; then
  if [ -f "$LAST_BACKUP_FILE" ]; then
    RESTORE_DIR=$(cat "$LAST_BACKUP_FILE")
    if [ -f "$RESTORE_DIR/config" ]; then
      echo -e "${BLUE}ℹ Restoring original kubeconfig${NC}"
      cp "$RESTORE_DIR/config" "$KUBE_DIR/config"
      echo -e "${GREEN}✓ Original kubeconfig restored from $RESTORE_DIR/config${NC}"
      exit 0
    else
      echo -e "${RED}✗ Backup file not found at $RESTORE_DIR/config${NC}"
      exit 1
    fi
  else
    echo -e "${RED}✗ No backup information found. Cannot restore.${NC}"
    exit 1
  fi
fi

# Create backup of current config if it exists
if [ -f "$KUBE_DIR/config" ]; then
  echo -e "${BLUE}ℹ Backing up existing kubeconfig${NC}"
  mkdir -p "$BACKUP_DIR"
  cp "$KUBE_DIR/config" "$BACKUP_DIR/config"
  echo "$BACKUP_DIR" > "$LAST_BACKUP_FILE"
  echo -e "${GREEN}✓ Backup created at $BACKUP_DIR/config${NC}"
fi

# Copy the sample kubeconfig
echo -e "${BLUE}ℹ Setting up demo kubeconfig${NC}"
mkdir -p "$KUBE_DIR"
cp "$(pwd)/sample-kubeconfig.yaml" "$KUBE_DIR/config"

# Set permissions
chmod 600 "$KUBE_DIR/config"
echo -e "${GREEN}✓ Sample kubeconfig installed to $KUBE_DIR/config${NC}"

echo
echo -e "${YELLOW}Demo kubeconfig contains:${NC}"
echo -e "- 9 contexts with different combinations of clusters, users and namespaces"
echo -e "- 6 clusters (dev, staging, production, gke, aws-east, aws-west)"
echo -e "- 5 users (dev-user, admin-user, prod-user, gke-user, aws-user)"
echo -e "- Multiple namespaces per context"
echo
echo -e "${YELLOW}Commands to try:${NC}"
echo -e "- ./kontext              ${GREEN}# Interactive context selector${NC}"
echo -e "- ./kontext list         ${GREEN}# List all contexts${NC}"
echo -e "- ./kontext current      ${GREEN}# Show current context and namespace${NC}"
echo -e "- ./kontext switch -n    ${GREEN}# Switch context with namespace selector${NC}"
echo -e "- ./kontext ns           ${GREEN}# Switch namespace in current context${NC}"
echo
echo -e "${BLUE}ℹ When done testing, restore your original config:${NC}"
echo -e "./setup-demo-kubeconfig.sh restore"
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
