#!/bin/bash

# Complete Teleport SSH Setup
# Runs both initial setup and node alias update

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo "=== Complete Teleport SSH Setup ==="
echo ""

# Run initial setup
if [ -f "$SCRIPT_DIR/setup-teleport-ssh.sh" ]; then
    bash "$SCRIPT_DIR/setup-teleport-ssh.sh"
else
    echo "Error: setup-teleport-ssh.sh not found"
    exit 1
fi

echo ""
echo "Waiting 2 seconds before updating node list..."
sleep 2
echo ""

# Update node aliases
if [ -f "$SCRIPT_DIR/update-teleport-nodes.sh" ]; then
    bash "$SCRIPT_DIR/update-teleport-nodes.sh"
else
    echo "Error: update-teleport-nodes.sh not found"
    exit 1
fi

echo ""
echo "=== All Done! ==="
echo ""
echo "Your Teleport SSH setup is complete and node aliases are updated."
echo "You can now connect using 'ssh <node-name>' or through VS Code Remote-SSH"
