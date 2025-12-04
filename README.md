# Scicom Helper - Teleport EC2 Access Tool

## Overview

`scicom-helper` is a CLI tool that simplifies connecting to EC2 instances managed by Teleport. It provides an interactive interface for managing Teleport access and automatically configures SSH for VS Code Remote-SSH integration.

**Teleport Server:** `teleport-iam.aies.scicom.dev`

## Prerequisites

Install Teleport CLI (tsh):

**macOS:**
```bash
# Download the latest version from https://goteleport.com/download
curl -O https://cdn.teleport.dev/teleport-18.4.2.pkg
sudo installer -pkg teleport-18.4.2.pkg -target /
```

**Note:** Check [Teleport Downloads](https://goteleport.com/download) for the latest version number.

**Linux:**
```bash
curl https://goteleport.com/static/install.sh | bash -s 15.0.0
```

**Windows:**
Use WSL2 and follow Linux instructions above.

Install WSL2: https://learn.microsoft.com/en-us/windows/wsl/install

Verify installation:
```bash
tsh version
```

## Quick Start

### 1. Download the Binary

Download the latest release for your platform:

**macOS (Apple Silicon):**
```bash
curl -L -o scicom-helper https://github.com/YOUR_ORG/teleport/releases/latest/download/scicom-helper-darwin-arm64
chmod +x scicom-helper
sudo mv scicom-helper /usr/local/bin/
```

**Linux (AMD64):**
```bash
curl -L -o scicom-helper https://github.com/YOUR_ORG/teleport/releases/latest/download/scicom-helper-linux-amd64
chmod +x scicom-helper
sudo mv scicom-helper /usr/local/bin/
```

**Linux (ARM64):**
```bash
curl -L -o scicom-helper https://github.com/YOUR_ORG/teleport/releases/latest/download/scicom-helper-linux-arm64
chmod +x scicom-helper
sudo mv scicom-helper /usr/local/bin/
```

**Verify Installation:**
```bash
scicom-helper
```

### 2. Login to Teleport

Run the tool and select **"Teleport Setup (Login)"**:

```bash
scicom-helper
```

This will:
- Open your browser for GitHub SSO authentication
- Log you into Teleport
- Verify your access

### 3. Configure VS Code (Recommended)

Select **"Configure VS Code for Teleport"**

This will automatically set the required VS Code setting:
- Sets `remote.SSH.useLocalServer = false`
- Required for Teleport SSH connections to work properly
- Backs up your existing VS Code settings

**Note:** Restart VS Code after running this step.

### 4. Update SSH Configuration

Select **"Teleport Update Nodes (Update SSH config)"**

This will:
- Backup your existing `~/.ssh/config`
- Fetch all accessible nodes from Teleport
- Generate optimized SSH configuration
- Enable VS Code Remote-SSH integration

**IMPORTANT:** Re-run this step whenever you gain access to new EC2 instances!

### 5. Connect to Nodes

You can now connect in three ways:

**Option A: Via the CLI**

Select **"Teleport SSH (Connect to a node)"**
- Choose a node from the list
- Select a login user (ubuntu, root, etc.)
- Connect automatically

**Option B: Via VS Code Remote-SSH**

1. Open VS Code
2. Press `Cmd+Shift+P` (macOS) or `Ctrl+Shift+P` (Windows/Linux)
3. Type "Remote-SSH: Connect to Host"
4. Select your node from the dropdown
5. Connect!

**Option C: Via SSH Command**

```bash
ssh <node-name>
```

Example:
```bash
ssh ip-172-31-16-103
```

## Features

- **Interactive Mode**: Arrow-key navigation for all operations
- **GitHub SSO**: Seamless authentication via GitHub
- **Auto SSH Config**: Automatically updates `~/.ssh/config` with all accessible nodes
- **VS Code Integration**: Nodes appear in VS Code Remote-SSH dropdown
- **Smart Login Detection**: Automatically detects and prioritizes available logins (ubuntu > root > others)
- **Safe Updates**: Backs up SSH config before making changes

## Important Notes

### When to Re-run "Update Nodes"

Run **"Teleport Update Nodes"** again whenever:
- You've been granted access to new EC2 instances
- New instances have been added to Teleport
- You want to refresh your SSH configuration
- Your nodes don't appear in VS Code Remote-SSH

### Session Expiry

Teleport sessions expire after **12 hours**. If you see authentication errors:

1. Run `scicom-helper`
2. Select **"Teleport Setup (Login)"** again
3. Re-authenticate via browser

## Troubleshooting

### "tsh: command not found"
Install Teleport CLI from prerequisites section above.

### "Not logged in to Teleport"
Run `scicom-helper` and select **"Teleport Setup (Login)"**.

### "No nodes found"
- Verify access: `tsh ls`
- Request EC2 access via Jira: https://scicom-ai-es.atlassian.net/jira/core/projects/IR/list
- Contact Platform Engineering team

### "SSH connection failed"
1. Check Teleport login: `tsh status`
2. Re-run: **"Teleport Update Nodes"**
3. Verify node list: `tsh ls`

### Nodes not appearing in VS Code
1. Re-run: **"Teleport Update Nodes"**
2. Restart VS Code
3. Check `~/.ssh/config` contains your nodes

## Build from Source

If you prefer to build the tool yourself:

### Prerequisites
- Go 1.22 or higher

### Build and Install

```bash
# Clone the repository
cd /path/to/teleport

# Build binary
make build

# Install to /usr/local/bin
make install

# Verify installation
scicom-helper
```

### Build for All Platforms

```bash
make build-all
```

Outputs:
- `build/scicom-helper-darwin-amd64` (macOS Intel)
- `build/scicom-helper-darwin-arm64` (macOS Apple Silicon)
- `build/scicom-helper-linux-amd64` (Linux x86_64)
- `build/scicom-helper-linux-arm64` (Linux ARM64)

## Project Structure

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

## Development

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

## Creating a Release

The project uses GitHub Actions for automated releases.

### Option 1: Automated Release (Recommended)

1. Go to [Actions → Create Release](../../actions/workflows/create-release.yml)
2. Click **"Run workflow"**
3. Select version bump type:
   - **patch**: Bug fixes (1.0.0 → 1.0.1)
   - **minor**: New features (1.0.0 → 1.1.0)
   - **major**: Breaking changes (1.0.0 → 2.0.0)
4. Click **"Run workflow"**

GitHub Actions will automatically:
- Calculate the new version
- Create and push the tag
- Build binaries
- Create a GitHub Release

### Option 2: Manual Release

```bash
make release-check  # Review release checklist
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

For detailed instructions, see [RELEASING.md](RELEASING.md).

## Support

For issues or questions:
- Create a Jira ticket: https://scicom-ai-es.atlassian.net/jira/core/projects/IR/list
- Contact: Platform Engineering Team
- Email: adha.sahar@scicom.com.my

## Additional Resources

- [Teleport Documentation](https://goteleport.com/docs/)
- [VS Code Remote-SSH](https://code.visualstudio.com/docs/remote/ssh)
- [Official Teleport VS Code Guide](https://goteleport.com/docs/enroll-resources/server-access/guides/vscode/)

---

**Last Updated:** December 4, 2025
**Maintained By:** Platform Engineering Team
