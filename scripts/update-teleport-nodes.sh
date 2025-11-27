#!/bin/bash

# Update Teleport Node Aliases in SSH Config
# Fetches current nodes and updates SSH config with specific host entries

set -e

TELEPORT_PROXY="teleport.aies.scicom.dev"
SSH_CONFIG="$HOME/.ssh/config"
MARKER_START="# BEGIN AUTO-GENERATED TELEPORT NODES"
MARKER_END="# END AUTO-GENERATED TELEPORT NODES"

echo "=== Update Teleport Node Aliases ==="
echo ""

# Check if tsh is installed
if ! command -v tsh &> /dev/null; then
    echo "Error: tsh (Teleport CLI) is not installed"
    exit 1
fi

# Check if logged in
if ! tsh status &> /dev/null; then
    echo "You are not logged in to Teleport. Logging in now..."
    tsh login --proxy=$TELEPORT_PROXY:443
fi

# Get list of nodes
echo "Fetching node list from Teleport..."
NODES=$(tsh ls --format=json | jq -r '.[].spec.hostname' 2>/dev/null || tsh ls | tail -n +2 | awk '{print $1}')

if [ -z "$NODES" ]; then
    echo "Warning: No nodes found"
    exit 0
fi

# Create temporary file for new node configurations
TEMP_FILE=$(mktemp)

cat > "$TEMP_FILE" << EOF
$MARKER_START
# Auto-generated Teleport node aliases
# Last updated: $(date)
# Run update-teleport-nodes.sh to refresh this section

EOF

# Add each node as a specific host entry
while IFS= read -r node; do
    if [ -n "$node" ]; then
        cat >> "$TEMP_FILE" << EOF
Host $node
    HostName $node.teleport.aies.scicom.dev
    User ubuntu

EOF
    fi
done <<< "$NODES"

echo "$MARKER_END" >> "$TEMP_FILE"

# Remove old auto-generated section if it exists
if grep -q "$MARKER_START" "$SSH_CONFIG" 2>/dev/null; then
    echo "Removing old node aliases..."
    sed -i.tmp "/$MARKER_START/,/$MARKER_END/d" "$SSH_CONFIG"
    rm -f "$SSH_CONFIG.tmp"
fi

# Append new configuration
cat "$TEMP_FILE" >> "$SSH_CONFIG"
rm "$TEMP_FILE"

echo ""
echo "=== Update Complete! ==="
echo ""
echo "Added $(echo "$NODES" | wc -l) node aliases to SSH config"
echo ""
echo "Available nodes:"
echo "$NODES"
echo ""
echo "You can now use these in VS Code Remote-SSH"
