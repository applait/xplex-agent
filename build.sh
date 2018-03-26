#!/bin/sh

set -e

cleanRebuild () {
    rm -rf bin/
    mkdir -p bin/
}

buildServerRelease () {
    echo "Building xplex-agent server release"
    CGO_ENABLED=0 GOOS=linux go build -o bin/xplex-agent -a -ldflags '-extldflags "-static"' .
}

buildServerDev () {
    echo "Building xplex-agent server dev"
    go build -o bin/xplex-agent .
}

buildDev () {
    cleanRebuild
    buildServerDev
}

buildRelease () {
    cleanRebuild
    buildServerRelease
}

case "$1" in
    "release")
    buildRelease
    ;;
    "dev")
    buildDev
    ;;
    *)
    echo "Usage: ./build.sh [dev|release|help]"
    echo "Binaries are put in './bin/'"
    ;;
esac
