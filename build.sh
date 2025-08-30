#!/bin/bash
# Build script for hyperctl using GoReleaser
# This replaces the previous shell script approach

set -e

COMMAND=${1:-"build"}
VERSION=${2:-"snapshot"}

echo "Building hyperctl using GoReleaser..."

case "$COMMAND" in
    "build")
        echo "Building snapshot for all platforms..."
        goreleaser build --snapshot --clean
        ;;
    "build-single")
        echo "Building snapshot for current platform only..."
        goreleaser build --snapshot --clean --single-target
        ;;
    "release")
        if [[ "$VERSION" == "snapshot" ]]; then
            echo "Error: Please provide a version for release (e.g., v1.0.0)"
            echo "Usage: $0 release v1.0.0"
            exit 1
        fi
        echo "Creating release $VERSION..."
        git tag "$VERSION" 2>/dev/null || echo "Tag $VERSION already exists"
        goreleaser release --clean
        ;;
    "check")
        echo "Checking GoReleaser configuration..."
        goreleaser check
        ;;
    *)
        echo "Usage: $0 {build|build-single|release|check} [version]"
        echo ""
        echo "Commands:"
        echo "  build        - Build snapshot for all platforms"
        echo "  build-single - Build snapshot for current platform only"
        echo "  release      - Create a tagged release (requires version)"
        echo "  check        - Validate GoReleaser configuration"
        echo ""
        echo "Examples:"
        echo "  $0 build"
        echo "  $0 build-single"
        echo "  $0 release v1.0.0"
        echo "  $0 check"
        exit 1
        ;;
esac

echo "Build completed successfully!"