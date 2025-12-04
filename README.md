# Accessing EC2 via Teleport with VS Code Remote

## Overview

This guide helps developers connect to EC2 instances managed by Teleport using VS Code Remote-SSH. All EC2 instances are automatically enrolled in Teleport and accessed through short SSH aliases for convenience.

**Teleport Server:** teleport-iam.aies.scicom.dev

## Quick Start: Using Scicom-Helper CLI (Recommended)

The easiest way to get started is using the `scicom-helper` CLI tool with interactive mode.

### Installation

```bash
# Build from source
cd /path/to/teleport
make build

# Install to your PATH
make install
```

### Usage

Simply run:
```bash
scicom-helper
```

Use arrow keys to navigate the menu:
- **Teleport Setup (Login)**: Log in with GitHub SSO
- **Teleport Update Nodes**: Configure SSH for all nodes (VS Code compatible)
- **Teleport SSH**: Interactive node and login user selection, then connect
- **Exit**: Quit the tool

After running "Teleport Update Nodes", all your nodes will appear in VS Code's Remote-SSH dropdown!

For detailed CLI documentation, see the [Scicom-Helper CLI Guide](#scicom-helper-cli) section below.

---

## Manual Setup (Alternative)

## Prerequisites

- Teleport CLI installed on your machine
- VS Code with Remote-SSH extension
- EC2 access approved via Jira ticket

## Step 1: Request EC2 Access

Create a Jira ticket in the Infrastructure Request project: https://scicom-ai-es.atlassian.net/jira/core/projects/IR/list

Include the following information:

- Purpose of EC2 instance
- Specifications needed (CPU, memory, GPU)
- Expected duration of use
- Project/team information

**For Platform Engineers:** Once ticket is created, provision the EC2 instance with auto-enrollment to Teleport enabled.

## Step 2: Install Teleport CLI

### macOS

```bash
brew install teleport
```

### Linux

```bash
curl https://goteleport.com/static/install.sh | bash -s 15.0.0
```

### Windows

Download from: https://goteleport.com/download

Verify installation:

```bash
tsh version
```

## Step 3: Login to Teleport

```bash
tsh login --proxy=teleport-iam.aies.scicom.dev --auth=github-connector
```

This will open your browser for authentication. Complete the login process.

Verify access and list available nodes:

```bash
tsh ls
```

Example output:

```
Node Name      Address    Labels
-------------- ---------- -------------------------------------------
ip-10-0-101-89 ⟵ Tunnel   environment=development,managed-by=teleport
```

## Step 4: Configure SSH for Teleport

Run our automated setup script (recommended) or follow manual steps below.

### Option A: Automated Setup (Recommended)

Download and run the setup script:

```bash
curl -o ~/setup-teleport-ssh.sh https://raw.githubusercontent.com/YOUR_ORG/scripts/main/setup-teleport-ssh.sh
chmod +x ~/setup-teleport-ssh.sh
~/setup-teleport-ssh.sh
```

This script will:

- Generate Teleport SSH configuration
- Add shorthand aliases for all nodes
- Configure VS Code compatibility
- Backup your existing SSH config

### Option B: Manual Setup

#### 1. Generate Teleport Configuration

```bash
tsh config --proxy=teleport.aies.scicom.dev > ~/teleport-ssh-config.txt
```

#### 2. Append to SSH Config

```bash
cat ~/teleport-ssh-config.txt >> ~/.ssh/config
```

#### 3. Add Wildcard Alias Pattern

Add this to your `~/.ssh/config` after the Teleport configuration:

```
# Custom aliases - add ip-* to Teleport patterns
Host ip-* *.teleport.aies.scicom.dev teleport.aies.scicom.dev
    UserKnownHostsFile "/Users/YOUR_USERNAME/.tsh/known_hosts"
    IdentityFile "/Users/YOUR_USERNAME/.tsh/keys/teleport.aies.scicom.dev/YOUR_TSH_USER"
    CertificateFile "/Users/YOUR_USERNAME/.tsh/keys/teleport.aies.scicom.dev/YOUR_TSH_USER-ssh/teleport.aies.scicom.dev-cert.pub"

Host ip-* *.teleport.aies.scicom.dev !teleport.aies.scicom.dev
    Port 3022
    ProxyCommand "/usr/local/bin/tsh" proxy ssh --cluster=teleport.aies.scicom.dev --proxy=teleport.aies.scicom.dev:443 %r@%h:%p

# Wildcard alias definition
Host ip-*
    HostName %h.teleport.aies.scicom.dev
    User ubuntu
```

**Important:** Replace `YOUR_USERNAME` and `YOUR_TSH_USER` with your actual values.

## Step 5: Connect via SSH

Now you can connect using just the node name:

```bash
ssh ip-10-0-101-89
```

Instead of the full command:

```bash
ssh ubuntu@ip-10-0-101-89.teleport.aies.scicom.dev
```

## Step 6: Configure VS Code Remote-SSH

1. Open VS Code
2. Install the **Remote - SSH** extension if not already installed
3. Press `Cmd+Shift+P` (macOS) or `Ctrl+Shift+P` (Windows/Linux)
4. Type "Remote-SSH: Connect to Host"
5. Select your EC2 instance (e.g., `ip-10-0-101-89`)
6. VS Code will connect through Teleport automatically

Your EC2 instances will appear in the VS Code Remote-SSH host list with their short names.

## Automation Scripts

### Script 1: Setup Teleport SSH Config

**File:** `scripts/setup-teleport-ssh.sh`

This script automates the entire SSH configuration process. See the script file for details.

### Script 2: Update Node Aliases

**File:** `scripts/update-teleport-nodes.sh`

This script fetches all available nodes from Teleport and adds specific aliases to SSH config. See the script file for details.

### Script 3: Combined Setup

**File:** `scripts/teleport-full-setup.sh`

This script combines both setup and node update. See the script file for details.

## Usage Instructions

### Initial Setup

```bash
# Download all scripts
curl -o ~/setup-teleport-ssh.sh https://raw.githubusercontent.com/YOUR_ORG/scripts/main/setup-teleport-ssh.sh
curl -o ~/update-teleport-nodes.sh https://raw.githubusercontent.com/YOUR_ORG/scripts/main/update-teleport-nodes.sh
curl -o ~/teleport-full-setup.sh https://raw.githubusercontent.com/YOUR_ORG/scripts/main/teleport-full-setup.sh

# Make executable
chmod +x ~/setup-teleport-ssh.sh ~/update-teleport-nodes.sh ~/teleport-full-setup.sh

# Run complete setup
~/teleport-full-setup.sh
```

### Updating Node List

When new EC2 instances are added or removed, refresh your node aliases:

```bash
~/update-teleport-nodes.sh
```

This ensures VS Code Remote-SSH shows all available nodes.

## Troubleshooting

### Connection Hangs or Timeout

**Issue:** SSH connection hangs or times out

**Solution:**

- Ensure you're logged in to Teleport: `tsh login --proxy=teleport.aies.scicom.dev:443`
- Check Teleport status: `tsh status`
- Verify node is accessible: `tsh ls`
- Try connecting with verbose output: `ssh -vvv ip-10-0-101-89`

### VS Code Can't Find Host

**Issue:** Node doesn't appear in VS Code Remote-SSH

**Solution:**

- Run the node update script: `~/update-teleport-nodes.sh`
- Reload VS Code
- Check `~/.ssh/config` contains your node

### Certificate Expired

**Issue:** "Certificate has expired" error

**Solution:**

```bash
tsh login --proxy=teleport.aies.scicom.dev:443
```

Teleport certificates expire after 12 hours by default.

### Permission Denied

**Issue:** "Permission denied (publickey)"

**Solution:**

- Verify you have access: `tsh ls`
- Check if node name is correct
- Ensure SSH config has correct IdentityFile path
- Contact platform team if you need access to specific nodes

### Wrong User

**Issue:** Connecting as wrong user

**Solution:**

The default user is `ubuntu`. If you need a different user, specify it:

```bash
ssh different_user@ip-10-0-101-89
```

Or update the alias in `~/.ssh/config`:

```
Host ip-10-0-101-89
    HostName ip-10-0-101-89.teleport.aies.scicom.dev
    User different_user
```

## Best Practices

- **Regular Updates:** Run `update-teleport-nodes.sh` weekly to keep node list current
- **Session Management:** Teleport sessions expire after 12 hours - re-login when needed
- **Security:** Never share your Teleport credentials or SSH keys
- **Naming Convention:** EC2 instances follow the pattern `ip-10-0-XXX-XX`
- **Backup:** The scripts automatically backup your SSH config before making changes

## Support

For issues or questions:

- Create a ticket in Infrastructure Request: https://scicom-ai-es.atlassian.net/jira/core/projects/IR/list
- Contact: Platform Engineering Team
- Email: adha.sahar@scicom.com.my

## Additional Resources

- [Teleport Documentation](https://goteleport.com/docs/)
- [VS Code Remote-SSH](https://code.visualstudio.com/docs/remote/ssh)
- [Official Teleport VS Code Guide](https://goteleport.com/docs/enroll-resources/server-access/guides/vscode/)

---

---

## Scicom-Helper CLI

The `scicom-helper` is a Go-based CLI tool that provides an interactive interface for managing Teleport access.

### Features

- **Interactive Mode**: Arrow-key navigation for all operations
- **Teleport Setup**: GitHub SSO authentication
- **Automatic SSH Config**: Updates `~/.ssh/config` with all nodes
- **VS Code Integration**: Nodes appear in VS Code Remote-SSH dropdown
- **Node & Login Selection**: Interactive selection of nodes and login users (ubuntu, root, ec2-user, etc.)
- **Safe Updates**: Automatically backs up SSH config before changes
- **Future-Proof**: Designed to be extended with more features

### Installation

#### Prerequisites

1. Teleport CLI (tsh) installed: https://goteleport.com/download
2. Go 1.22+ (for building)

#### Build and Install

```bash
# Clone/navigate to repository
cd /path/to/teleport

# Build binary
make build

# Install to /usr/local/bin
make install

# Verify installation
scicom-helper
```

#### Build for Multiple Platforms

```bash
make build-all
```

Outputs:
- `build/scicom-helper-darwin-amd64` (macOS Intel)
- `build/scicom-helper-darwin-arm64` (macOS Apple Silicon)
- `build/scicom-helper-linux-amd64` (Linux x86_64)
- `build/scicom-helper-linux-arm64` (Linux ARM64)

### Usage Guide

#### 1. First Time Setup

```bash
scicom-helper
```

Select: `Teleport Setup (Login)`
- Logs you into Teleport using GitHub SSO
- Opens browser for authentication
- Verifies login status

#### 2. Configure SSH

Select: `Teleport Update Nodes (Update SSH config)`
- Backs up existing SSH config
- Fetches all Teleport-managed nodes
- Generates optimized SSH configuration
- Enables VS Code Remote-SSH integration

The tool creates a managed section in `~/.ssh/config`:
```
# BEGIN SCICOM-HELPER TELEPORT CONFIG
# ... auto-generated configuration ...
# END SCICOM-HELPER TELEPORT CONFIG
```

Running this command again will safely update only this section.

#### 3. Connect to Nodes

Select: `Teleport SSH (Connect to a node)`
- Lists all available nodes
- Use arrow keys to select a node
- Lists available login users for the selected node
- Use arrow keys to select a login user (ubuntu, root, ec2-user, etc.)
- Automatically connects via `tsh ssh` with the selected login
- Press Ctrl+D or type `exit` to disconnect

#### 4. VS Code Remote-SSH

After running "Teleport Update Nodes":
1. Open VS Code
2. Press `Cmd+Shift+P` (macOS) or `Ctrl+Shift+P` (Windows/Linux)
3. Select "Remote-SSH: Connect to Host"
4. See all your nodes in the dropdown
5. Click to connect!

### Configuration

The CLI uses these settings:
- **Proxy**: `teleport-iam.aies.scicom.dev`
- **Auth Method**: `github-connector` (GitHub SSO)
- **Default User**: `ubuntu`
- **Port**: `3022` (for proxy connections)

### Project Structure

```
.
├── main.go              # Entry point
├── cmd/
│   ├── root.go          # CLI framework & interactive menu
│   ├── setup.go         # Teleport login
│   ├── update_nodes.go  # SSH config management
│   ├── ssh.go           # Interactive SSH connection
│   └── utils.go         # Helper functions
├── Makefile             # Build automation
├── go.mod               # Go dependencies
└── README.md            # This file
```

### Development

```bash
# Install dependencies
go mod download

# Build
go build -o scicom-helper .

# Run directly (without installing)
go run .

# Run tests
make test

# Clean build artifacts
make clean
```

### Troubleshooting

#### "tsh: command not found"
Install Teleport CLI: https://goteleport.com/download

#### "Not logged in to Teleport"
Run the "Teleport Setup" option first.

#### "No nodes found"
Verify access: `tsh ls`
Contact platform team if you need node access.

#### SSH connection fails
1. Check login: `tsh status`
2. Update nodes: Run "Teleport Update Nodes"
3. Verify node list: `tsh ls`

### Extending the CLI

The Platform Engineering team can add new features:

1. Create new file in `cmd/` directory
2. Implement feature function
3. Add option to menu in `cmd/root.go`
4. Update this README

Example for adding a new feature:
```go
// cmd/new_feature.go
package cmd

func newFeature() error {
    // Implementation
    return nil
}

// Add to cmd/root.go menu:
Options: []string{
    "Teleport Setup (Login)",
    "Teleport Update Nodes",
    "Teleport SSH",
    "New Feature",  // Add here
    "Exit",
}
```

---

**Last Updated:** December 3, 2025
**Maintained By:** Platform Engineering Team
