# Building hyperctl with GoReleaser

This project uses [GoReleaser](https://goreleaser.com/) to build cross-platform binaries for hyperctl.

## Prerequisites

1. Install GoReleaser:
   ```bash
   # macOS
   brew install goreleaser
   
   # Linux
   curl -sfL https://goreleaser.com/static/run | bash -s -- --version
   
   # Or download from https://github.com/goreleaser/goreleaser/releases
   ```

2. Ensure you have Go installed and dependencies are available.

## Building

### Development/Snapshot Build

To build for the current platform only (for development):
```bash
goreleaser build --snapshot --clean --single-target
```

To build for all platforms (snapshot, no release):
```bash
goreleaser build --snapshot --clean
```

### Release Build

To create a full release (requires a git tag):
```bash
# Tag the release
git tag v1.0.0

# Build and release
goreleaser release --clean
```

### Building without Git Tags

For builds without creating a release:
```bash
goreleaser build --snapshot --clean
```

## Output

GoReleaser will create binaries in the `dist/` directory with the following structure:
- Linux: `dist/hyperctl_linux_amd64_v1/hyperctl`, `dist/hyperctl_linux_arm64/hyperctl`
- Windows: `dist/hyperctl_windows_amd64_v1/hyperctl.exe`, `dist/hyperctl_windows_arm64/hyperctl.exe`
- macOS: `dist/hyperctl_darwin_amd64_v1/hyperctl`, `dist/hyperctl_darwin_arm64/hyperctl`, and universal binary

Archives (tar.gz/zip) will be created in `dist/` as well.

## Configuration

The build configuration is in `.goreleaser.yml` and handles:
- Cross-platform compilation for Linux, Windows, and macOS
- Both amd64 and arm64 architectures
- Universal macOS binaries (combines amd64 and arm64)
- CGO disabled and trimpath flags
- Proper archive naming and packaging

## Migration from Shell Script

This GoReleaser configuration replaces the previous shell script and provides:
- Better cross-platform build management
- Automatic archive creation
- Checksum generation
- Release automation
- Universal binary support for macOS