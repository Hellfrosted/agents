#!/usr/bin/env bash
set -euo pipefail

ROOT="$(unset CDPATH; cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
DIST="${SK_UP_DIST_DIR:-$ROOT/dist/sk-up}"
VERSION="${SK_UP_VERSION:-dev}"

TARGETS=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
  "windows/arm64"
)

rm -rf "$DIST"
mkdir -p "$DIST"

(
  cd "$ROOT"
  go test ./...
)

for target in "${TARGETS[@]}"; do
  goos="${target%/*}"
  goarch="${target#*/}"
  name="sk-up-$VERSION-$goos-$goarch"
  work="$DIST/$name"
  mkdir -p "$work"

  binary="sk-up"
  if [ "$goos" = "windows" ]; then
    binary="sk-up.exe"
  fi

  (
    cd "$ROOT"
    CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" go build -trimpath -ldflags="-s -w" -o "$work/$binary" ./cmd/sk-up
  )

  if [ "$goos" = "windows" ]; then
    cp "$ROOT/bin/sk-up.cmd" "$work/sk-up.cmd"
    cp "$ROOT/bin/skills-updates.cmd" "$work/skills-updates.cmd"
    (
      cd "$DIST"
      zip -qr "$name.zip" "$name"
    )
  else
    ln "$work/$binary" "$work/skills-updates"
    (
      cd "$DIST"
      tar -czf "$name.tar.gz" "$name"
    )
  fi
done

(
  cd "$DIST"
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum ./*.tar.gz ./*.zip > SHA256SUMS
  else
    shasum -a 256 ./*.tar.gz ./*.zip > SHA256SUMS
  fi
)

printf 'wrote %s\n' "$DIST"
