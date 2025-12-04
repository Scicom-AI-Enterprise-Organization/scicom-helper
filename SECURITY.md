# Security Policy

## Reporting Security Vulnerabilities

We take the security of `scicom-helper` seriously. If you discover a security vulnerability, please follow these steps:

### 🔒 Private Disclosure

**DO NOT** open a public GitHub issue for security vulnerabilities.

Instead, please report security issues privately:

1. **Email:** [adha.sahar@scicom.com.my](mailto:adha.sahar@scicom.com.my)
2. **Subject:** "SECURITY: [Brief description]"
3. **Include:**
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

### Response Timeline

- **Initial Response:** Within 48 hours
- **Status Update:** Within 1 week
- **Resolution:** Varies based on severity and complexity

### What to Expect

1. We'll acknowledge your report within 48 hours
2. We'll investigate and provide regular updates
3. Once fixed, we'll coordinate disclosure timing with you
4. We'll credit you in the release notes (unless you prefer anonymity)

## Automated Security Scanning

This project uses automated security scanning with Trivy:

- **Continuous Scanning:** On every push and pull request
- **Daily Scans:** Scheduled at 6 AM UTC
- **Dependency Scanning:** Go modules and indirect dependencies
- **Severity Levels:** CRITICAL, HIGH, MEDIUM, LOW

### Viewing Scan Results

Security scan results are available in:
- [GitHub Actions](../../actions/workflows/security-scan.yml)
- [Security Tab](../../security) (SARIF reports)
- Pull Request comments (automatic summary)

## Supported Versions

We provide security updates for:

| Version | Supported          |
| ------- | ------------------ |
| Latest  | ✅ Yes             |
| < 1.0   | ❌ No              |

**Recommendation:** Always use the latest release from the [Releases page](../../releases/latest).

## Security Best Practices

When using `scicom-helper`:

### For Users

1. **Keep Updated:** Always use the latest version
2. **Verify Downloads:** Check SHA256 checksums from releases
3. **Review Permissions:** Understand what access Teleport roles grant
4. **Protect Credentials:** Never commit `.tsh` directory or SSH keys

### For Developers

1. **Update Dependencies:** Regularly run `go get -u` and check for updates
2. **Review Changes:** Check security scan results in PRs
3. **Handle Secrets:** Never commit credentials or API keys
4. **Code Review:** All changes require review before merging

## Dependencies

This project relies on:

- **Go 1.22+** - Core language runtime
- **Teleport CLI (tsh)** - Teleport access client
- Third-party Go modules (see `go.mod`)

### Dependency Scanning

Dependencies are automatically scanned for known vulnerabilities:
- On every commit
- Daily scheduled scans
- Pull request checks

## Vulnerability Disclosure

When we fix a security vulnerability:

1. **Patch Release:** We'll release a patched version immediately
2. **Security Advisory:** Published on GitHub Security Advisories
3. **Release Notes:** Security fixes highlighted in release notes
4. **Credits:** Security researchers credited (with permission)

## Known Issues and Limitations

Current security considerations:

- **Teleport Access:** Tool inherits Teleport user permissions
- **SSH Config:** Modifies `~/.ssh/config` (backups created)
- **Editor Settings:** Modifies VS Code/Cursor settings (backups created)
- **Platform Security:** Inherits OS-level security (file permissions, etc.)

## Contact

For security-related questions or concerns:

- **Security Issues:** [adha.sahar@scicom.com.my](mailto:adha.sahar@scicom.com.my)
- **General Support:** Create a [Jira ticket](https://scicom-ai-es.atlassian.net/jira/core/projects/IR/list)
- **Platform Team:** Platform Engineering Team

---

**Last Updated:** December 4, 2025
