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

**Linux (Ubuntu/Debian):**
```bash
sudo curl https://apt.releases.teleport.dev/gpg \
  -o /usr/share/keyrings/teleport-archive-keyring.asc
source /etc/os-release
echo "deb [signed-by=/usr/share/keyrings/teleport-archive-keyring.asc] \
https://apt.releases.teleport.dev/${ID?} ${VERSION_CODENAME?} stable/v18" \
| sudo tee /etc/apt/sources.list.d/teleport.list > /dev/null

sudo apt-get update
sudo apt-get install teleport
```

**Linux (CentOS/RHEL):**
```bash
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo \
  https://yum.releases.teleport.dev/teleport.repo
sudo yum install teleport
```

**Other Linux distributions:** See [Teleport Linux Installation](https://goteleport.com/docs/installation/linux/)

**Windows:**

Download the Windows installer from [Teleport Downloads](https://goteleport.com/download) or use the [community installer](https://goteleport.com/docs/installation/windows/).

**Verify installation:**
```bash
tsh version
```

## Quick Start

### 1. Download and Install

Go to the [latest release page](../../releases/latest) and follow the installation instructions for your platform.

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
- **IMPORTANT:** When prompted, grant access to the **AIES-Infra** organization
- Log you into Teleport
- Verify your access

### 3. Configure Editors (Recommended)

Select **"Configure VS Code for Teleport"**

This will automatically configure both VS Code and Cursor (if installed):
- Sets `remote.SSH.useLocalServer = false`
- Required for Teleport SSH connections to work properly
- Backs up your existing settings
- **Windows Users:** Fixes "posix_spawnp: No such file or directory" error in Remote-SSH

**Note:** Restart your editor after running this step.

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

**Option B: Via VS Code/Cursor Remote-SSH**

1. Open VS Code or Cursor
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
- **Editor Integration**: Automatically configures VS Code and Cursor for Remote-SSH
- **Windows Support**: Native Windows support without WSL2, fixes "posix_spawnp" error
- **Smart Login Detection**: Automatically detects and prioritizes available logins (ubuntu > root > others)
- **Safe Updates**: Backs up SSH config and editor settings before making changes

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

### Nodes not appearing in VS Code/Cursor
1. Re-run: **"Teleport Update Nodes"**
2. Restart your editor
3. Check `~/.ssh/config` contains your nodes

### "posix_spawnp: No such file or directory" error (Windows)
1. Run: **"Configure VS Code for Teleport"** to set `remote.SSH.useLocalServer = false`
2. Restart your editor
3. Try connecting again

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
- `build/scicom-helper-windows-amd64.exe` (Windows x86_64)

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

## Security

This project includes comprehensive automated security scanning with multiple tools.

### Security Scanning Tools

**SAST (Static Application Security Testing):**
- **[gosec](https://github.com/securego/gosec)** - Go-specific security scanner
- **[Semgrep](https://semgrep.dev/)** - Multi-language pattern-based SAST

**Dependency Scanning:**
- **[Trivy](https://github.com/aquasecurity/trivy)** - Comprehensive vulnerability scanner

### Automated Security Scans

Security scans run automatically:
- **On every push** to main branch
- **On every pull request**
- **Daily** at 6 AM UTC
- **Manually** via GitHub Actions workflow dispatch

### What Gets Scanned

1. **Source Code (SAST)**
   - Security vulnerabilities in Go code
   - Common coding mistakes and anti-patterns
   - OWASP Top 10 security issues
   - CWE Top 25 weaknesses

2. **Dependencies (SCA)**
   - Go module vulnerabilities
   - Transitive dependencies
   - Known CVEs in packages
   - License compliance issues

### Scan Results

View security scan results:
1. Go to the [Actions tab](../../actions/workflows/security-scan.yml)
2. Click on the latest workflow run
3. View the summary or download detailed reports from artifacts

Results are also available in:
- **GitHub Security tab** - SARIF format for native GitHub integration
- **Pull Request comments** - Automatic vulnerability summary (Trivy)
- **Workflow artifacts** - HTML, JSON, SARIF, and text reports

### Vulnerability Severity

The scanner checks for:
- 🔴 **Critical** - Immediate action required
- 🟠 **High** - Should be addressed soon
- 🟡 **Medium** - Review and plan remediation
- 🟢 **Low** - Monitor and address when possible

**Note:** PRs with critical or high severity vulnerabilities will fail the security check.

### Ignoring False Positives

To ignore false positives or accepted risks, add CVE IDs to [`.trivyignore`](.trivyignore):

```
# Example
CVE-2023-12345  # False positive - not applicable to our use case
```

### Reporting Security Issues

For security concerns, please see [SECURITY.md](SECURITY.md).

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
- [Cursor Editor](https://cursor.sh/)
- [Official Teleport VS Code Guide](https://goteleport.com/docs/enroll-resources/server-access/guides/vscode/)

---

**Last Updated:** December 4, 2025
**Maintained By:** Platform Engineering Team
