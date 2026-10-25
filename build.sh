#!/usr/bin/env bash
# Builds Snugfit.exe and the Windows installer.
# Needs: Go 1.24+, makensis (NSIS 3) on PATH. Works on Windows (Git Bash), Linux or macOS.
set -euo pipefail
cd "$(dirname "$0")"
VERSION=$(grep -oE 'Version = "[0-9.]+"' version.go | grep -oE '[0-9.]+')
echo "Building Snugfit $VERSION"
rm -rf dist && mkdir -p dist/Snugfit
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w -H windowsgui" -o dist/Snugfit/Snugfit.exe .
cp -r assets dist/Snugfit/assets
cp LICENSE.txt CHANGELOG.md THIRD-PARTY-NOTICES.txt dist/Snugfit/
makensis -V2 -DVERSION="$VERSION" -DSRCDIR="../dist/Snugfit" -DOUTDIR="../dist" installer/snugfit.nsi
# release notes = newest "## " section of the changelog
awk '/^## /{n++} n==1' CHANGELOG.md | tail -n +2 > dist/RELEASE_NOTES.md
echo "Done: dist/Snugfit-Setup-$VERSION.exe"
