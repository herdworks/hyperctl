#!/bin/sh

function build() {
    VERSION=${1:-$(version)}
    set -e
    reset
    build_cli ${VERSION}
    release ${VERSION}
}

function reset() {
    echo "\nCleaning up build artifacts"
    if [ -d "builds" ]; then rm -r builds/* 2> /dev/null; fi
    mkdir -p builds/{bin,release}
    go mod tidy
}

function version() {
    if ! $(git diff --quiet HEAD); then echo "dev"; return; fi
    if [ -n "$(git tag --points-at HEAD)" ]; then echo "$(git tag --points-at HEAD)"; return; fi
    echo "$(git rev-parse --short HEAD)"
}

function build_cli() {
    VERSION=${1:-"v0.0.0"}
    echo "\nBuilding hyperctl_${VERSION}_linux_arm64"
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -o "builds/bin/hyperctl_${VERSION}_linux_arm64" ./cmd/hyperctl
    echo "\nBuilding hyperctl_${VERSION}_linux_amd64"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o "builds/bin/hyperctl_${VERSION}_linux_amd64" ./cmd/hyperctl
    echo "\nBuilding hyperctl_${VERSION}_windows_arm64"
    CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -trimpath -o "builds/bin/hyperctl_${VERSION}_windows_arm64" ./cmd/hyperctl
    echo "\nBuilding hyperctl_${VERSION}_windows_amd64"
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -o "builds/bin/hyperctl_${VERSION}_windows_amd64" ./cmd/hyperctl
    echo "\nBuilding hyperctl_${VERSION}_darwin_arm64"
    CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -o "builds/bin/hyperctl_${VERSION}_darwin_arm64" ./cmd/hyperctl
    echo "\nBuilding hyperctl_${VERSION}_darwin_amd64"
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -o "builds/bin/hyperctl_${VERSION}_darwin_amd64" ./cmd/hyperctl
    echo "\nBuilding hyperctl_${VERSION}_darwin_universal"
    lipo -create builds/bin/hyperctl_${VERSION}_darwin_amd64 builds/bin/hyperctl_${VERSION}_darwin_arm64 -output builds/bin/hyperctl_${VERSION}_darwin_universal
    echo "\nCopying builds/bin/hyperctl_${VERSION}_darwin_universal to builds/bin/hyperctl"
    cp builds/bin/hyperctl_${VERSION}_darwin_universal builds/bin/hyperctl
}

function release() {
    VERSION=${1:-"v0.0.0"}
    echo "\nCreating ${VERSION} release archives"
    mkdir -p builds/release/${VERSION}
    for file in builds/bin/*; do
        PLATFORM=$(basename "${file}" | cut -d'_' -f3)
        ARCH=$(basename "${file}" | cut -d'_' -f4)
        if [ "${file}" == "builds/bin/hyperctl" ]; then continue; fi
        if [ -f "$file" ]; then
            echo "\nPackaging hyperctl for ${PLATFORM}/${ARCH} ("builds/release/$(basename ${file}).tar.gz")"
            mkdir -p "builds/release/${VERSION}/${PLATFORM}/${ARCH}"
            cp "${file}" "builds/release/${VERSION}/${PLATFORM}/${ARCH}/hyperctl"
            tar -czf "builds/release/$(basename ${file}).tar.gz" --directory "builds/release/${VERSION}/${PLATFORM}/${ARCH}" .
        fi
    done
}

build ${@}
