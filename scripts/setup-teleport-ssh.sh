#!/bin/bash

# Setup Teleport SSH Configuration
# This script configures SSH to work with Teleport-managed EC2 instances

set -e

TELEPORT_PROXY="teleport.aies.scicom.dev"
SSH_CONFIG="$HOME/.ssh/config"
BACKUP_CONFIG="$SSH_CONFIG.backup.$(date +%Y%m%d_%H%M%S)"

echo "=== Teleport SSH Setup Script ==="
echo ""

# Check if tsh is installed
if ! command -v tsh &> /dev/null; then
    echo "Error: tsh (Teleport CLI) is not installed"
    echo "Please install from: https://goteleport.com/download"
    exit 1
fi

# Check if logged in to Teleport
if ! tsh status &> /dev/null; then
    echo "You are not logged in to Teleport. Logging in now..."
    tsh login --proxy=$TELEPORT_PROXY:443
fi

# Backup existing SSH config
if [ -f "$SSH_CONFIG" ]; then
    echo "Backing up existing SSH config to: $BACKUP_CONFIG"
    cp "$SSH_CONFIG" "$BACKUP_CONFIG"
fi

# Get Teleport configuration
echo "Generating Teleport SSH configuration..."
TSH_CONFIG=$(tsh config --proxy=$TELEPORT_PROXY)

# Check if Teleport config already exists
if grep -q "Begin generated Teleport configuration" "$SSH_CONFIG" 2>/dev/null; then
    echo "Removing existing Teleport configuration..."
    sed -i.tmp '/# Begin generated Teleport configuration/,/# End generated Teleport configuration/d' "$SSH_CONFIG"
    rm -f "$SSH_CONFIG.tmp"
fi

# Get user details for wildcard pattern
TSH_USER=$(tsh status | grep 'User:' | awk '{print $2}')
TSH_KEYS_DIR="$HOME/.tsh/keys/$TELEPORT_PROXY"

# Write Teleport config to SSH config
echo "$TSH_CONFIG" >> "$SSH_CONFIG"

# Add wildcard alias patterns
cat >> "$SSH_CONFIG" << EOF

# Custom aliases - add ip-* to Teleport patterns
Host ip-* *.teleport.aies.scicom.dev teleport.aies.scicom.dev
    UserKnownHostsFile "$HOME/.tsh/known_hosts"
    IdentityFile "$TSH_KEYS_DIR/$TSH_USER"
    CertificateFile "$TSH_KEYS_DIR/$TSH_USER-ssh/teleport.aies.scicom.dev-cert.pub"

Host ip-* *.teleport.aies.scicom.dev !teleport.aies.scicom.dev
    Port 3022
    ProxyCommand "/usr/local/bin/tsh" proxy ssh --cluster=teleport.aies.scicom.dev --proxy=teleport.aies.scicom.dev:443 %r@%h:%p

# Wildcard alias definition
Host ip-*
    HostName %h.teleport.aies.scicom.dev
    User ubuntu
EOF

echo ""
echo "=== Setup Complete! ==="
echo ""
echo "Your SSH config has been updated with:"
echo "  1. Teleport SSH configuration"
echo "  2. Wildcard alias pattern for ip-* hosts"
echo ""
echo "You can now connect using: ssh ip-10-0-101-89"
echo ""
echo "Run 'tsh ls' to see available nodes"
