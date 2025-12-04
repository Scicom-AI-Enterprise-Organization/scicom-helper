# Release Process

This document describes how to create a new release of `scicom-helper`.

## Prerequisites

- Commit and push all changes
- Ensure all tests pass: `make test`
- Decide on version number (follow [Semantic Versioning](https://semver.org/))

## Quick Release

Run the pre-release check:

```bash
make release-check
```

## Creating a Release

### 1. Tag the Release

Create an annotated tag with your version:

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
```

Version format: `vMAJOR.MINOR.PATCH`
- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality, backward compatible
- **PATCH**: Bug fixes, backward compatible

### 2. Push the Tag

```bash
git push origin v1.0.0
```

### 3. GitHub Actions Workflow

Once the tag is pushed, GitHub Actions will automatically:

1. Build binaries for all platforms:
   - macOS ARM64 (`darwin/arm64`)
   - Linux AMD64 (`linux/amd64`)
   - Linux ARM64 (`linux/arm64`)
   - **Note:** Windows users should use WSL2 with the Linux binary

2. Generate SHA256 checksums

3. Create a GitHub Release with:
   - All binary artifacts
   - Checksums file
   - Auto-generated release notes
   - Installation instructions

### 4. Monitor the Build

1. Go to the [Actions tab](../../actions) in your repository
2. Watch the "Release" workflow
3. Once complete, check the [Releases page](../../releases)

## Workflow Details

The release workflow is defined in [`.github/workflows/release.yml`](.github/workflows/release.yml).

### Build Process

For each platform, the workflow runs:

```bash
GOOS=<os> GOARCH=<arch> go build -o build/scicom-helper-<os>-<arch> .
```

### Artifacts Produced

- `scicom-helper-darwin-arm64` - macOS Apple Silicon binary
- `scicom-helper-linux-amd64` - Linux x86_64 binary
- `scicom-helper-linux-arm64` - Linux ARM64 binary
- `checksums.txt` - SHA256 checksums for all binaries

**Windows Support:** Windows users should use WSL2 (Windows Subsystem for Linux) and run the Linux binary.

## Testing the Release

After creating a release, test the binaries:

### macOS

```bash
curl -L -o scicom-helper https://github.com/YOUR_ORG/teleport/releases/download/v1.0.0/scicom-helper-darwin-arm64
chmod +x scicom-helper
./scicom-helper
```

### Linux

```bash
curl -L -o scicom-helper https://github.com/YOUR_ORG/teleport/releases/download/v1.0.0/scicom-helper-linux-amd64
chmod +x scicom-helper
./scicom-helper
```

### Windows (via WSL2)

1. Install WSL2: https://learn.microsoft.com/en-us/windows/wsl/install
2. Open WSL2 terminal
3. Follow Linux installation instructions above

## Verifying Checksums

Users can verify download integrity:

```bash
# Download the checksums file
curl -L -o checksums.txt https://github.com/YOUR_ORG/teleport/releases/download/v1.0.0/checksums.txt

# Verify a specific binary
sha256sum -c checksums.txt --ignore-missing
```

## Troubleshooting

### Build Fails

1. Check the [Actions tab](../../actions) for error logs
2. Ensure `go.mod` is up to date: `go mod tidy`
3. Test local build: `make build-all`

### Release Not Created

1. Verify tag was pushed: `git ls-remote --tags origin`
2. Check workflow permissions in repository settings
3. Ensure tag follows pattern `v*` (e.g., `v1.0.0`, not `1.0.0`)

### Binary Not Working

1. Check Go version compatibility (requires Go 1.22+)
2. Verify target platform matches binary
3. Ensure execute permissions on Unix-like systems: `chmod +x scicom-helper`

## Deleting a Release

If you need to delete a release:

```bash
# Delete the tag locally
git tag -d v1.0.0

# Delete the tag remotely
git push origin :refs/tags/v1.0.0

# Manually delete the GitHub Release from the Releases page
```

## Example Release Commands

```bash
# Full release process
git add .
git commit -m "Prepare release v1.0.0"
git push origin main
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# Watch the build
open https://github.com/YOUR_ORG/teleport/actions

# View the release
open https://github.com/YOUR_ORG/teleport/releases
```

## Release Checklist

- [ ] All changes committed and pushed
- [ ] Tests passing: `make test`
- [ ] Version number decided (semantic versioning)
- [ ] Tag created: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`
- [ ] Tag pushed: `git push origin vX.Y.Z`
- [ ] GitHub Actions workflow completed successfully
- [ ] Release appears on GitHub Releases page
- [ ] Binaries download and run correctly
- [ ] Checksums verified
- [ ] Release notes updated (if needed)

---

**Last Updated:** December 4, 2025
