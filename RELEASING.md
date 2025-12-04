# Release Process

This document describes how to create a new release of `scicom-helper`.

## Prerequisites

- Commit and push all changes
- Ensure all tests pass: `make test`
- Understand [Semantic Versioning](https://semver.org/)

## Quick Release (Automated) - Recommended

### 1. Navigate to GitHub Actions

Go to: [Actions → Create Release](../../actions/workflows/create-release.yml)

### 2. Run the Workflow

1. Click **"Run workflow"** button
2. Select the **version bump type**:
   - **patch** (v1.0.0 → v1.0.1): Bug fixes, backward compatible
   - **minor** (v1.0.0 → v1.1.0): New features, backward compatible
   - **major** (v1.0.0 → v2.0.0): Breaking changes
3. *(Optional)* Enter a **custom version** to override automatic calculation
4. Click **"Run workflow"** to start

### 3. Workflow Execution

The workflow will automatically:
1. Fetch the latest tag (or start from v0.0.0 if none exists)
2. Calculate the new version based on your selection
3. Create an annotated git tag
4. Push the tag to GitHub
5. Trigger the Release workflow

### 4. Monitor Progress

- Watch the [Create Release workflow](../../actions/workflows/create-release.yml) complete
- The [Release workflow](../../actions/workflows/release.yml) will automatically start
- Check the [Releases page](../../releases) for the published release

---

## Manual Release (Alternative)

If you prefer manual control:

### 1. Run Pre-release Check

```bash
make release-check
```

### 2. Create Tag

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
```

Version format: `vMAJOR.MINOR.PATCH`
- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality, backward compatible
- **PATCH**: Bug fixes, backward compatible

### 3. Push Tag

```bash
git push origin v1.0.0
```

### 4. GitHub Actions Workflow

Once the tag is pushed, GitHub Actions will automatically:

1. Build binaries for all platforms:
   - macOS ARM64 (`darwin/arm64`)
   - Linux AMD64 (`linux/amd64`)
   - Linux ARM64 (`linux/arm64`)
   - Windows AMD64 (`windows/amd64`)

2. Generate SHA256 checksums

3. Create a GitHub Release with:
   - All binary artifacts
   - Checksums file
   - Auto-generated release notes
   - Installation instructions

### 5. Monitor the Build

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
- `scicom-helper-windows-amd64.exe` - Windows x86_64 binary
- `checksums.txt` - SHA256 checksums for all binaries

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

### Windows

```powershell
# Download using PowerShell
Invoke-WebRequest -Uri "https://github.com/YOUR_ORG/teleport/releases/download/v1.0.0/scicom-helper-windows-amd64.exe" -OutFile "scicom-helper.exe"

# Run the binary
.\scicom-helper.exe
```

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
