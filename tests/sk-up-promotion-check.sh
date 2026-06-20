#!/usr/bin/env bash
set -euo pipefail

ROOT="$(unset CDPATH; cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
TMPDIR="$(mktemp -d)"

cleanup() {
  find "$TMPDIR" -mindepth 1 -delete
  rmdir "$TMPDIR"
}
trap cleanup EXIT

if [ -e "$ROOT/bin/skills-updates.ps1" ]; then
  printf 'retired PowerShell implementation still exists\n' >&2
  exit 1
fi

if grep -Eiq 'powershell|skills-updates\.ps1' "$ROOT/bin/sk-up.cmd" "$ROOT/bin/skills-updates.cmd"; then
  printf 'Windows wrappers still reference PowerShell fallback\n' >&2
  exit 1
fi

grep -Fq 'sk-up.exe' "$ROOT/bin/sk-up.cmd"
grep -Fq 'sk-up.exe' "$ROOT/bin/skills-updates.cmd"
grep -Fq 'SK_UP_ENTRYPOINT=skills-updates' "$ROOT/bin/skills-updates.cmd"

SK_UP_DIST_DIR="$TMPDIR/dist" SK_UP_VERSION=promotion "$ROOT/bin/build-sk-up-release.sh" >/tmp/sk-up-promotion-release.out

python3 - "$TMPDIR/dist" <<'PY'
import pathlib
import sys
import zipfile

dist = pathlib.Path(sys.argv[1])
required = {
    "sk-up-promotion-windows-amd64/sk-up.exe",
    "sk-up-promotion-windows-amd64/sk-up.cmd",
    "sk-up-promotion-windows-amd64/skills-updates.cmd",
    "sk-up-promotion-windows-arm64/sk-up.exe",
    "sk-up-promotion-windows-arm64/sk-up.cmd",
    "sk-up-promotion-windows-arm64/skills-updates.cmd",
}
seen = set()
for archive in dist.glob("sk-up-promotion-windows-*.zip"):
    with zipfile.ZipFile(archive) as zf:
        names = set(zf.namelist())
        forbidden = sorted(name for name in names if name.endswith("/skills-updates.exe"))
        if forbidden:
            print("forbidden Windows archive entries:", *forbidden, sep="\n", file=sys.stderr)
            raise SystemExit(1)
        seen.update(names)
missing = sorted(required - seen)
if missing:
    print("missing Windows archive entries:", *missing, sep="\n", file=sys.stderr)
    raise SystemExit(1)
PY

go build -o "$TMPDIR/sk-up" "$ROOT/cmd/sk-up"
mkdir -p "$TMPDIR/agents/skills/alpha"
"$TMPDIR/sk-up" -h >/tmp/sk-up-promotion-help.out
"$TMPDIR/sk-up" -l --agents-home "$TMPDIR/agents" --cache-dir "$TMPDIR/cache" --state-dir "$TMPDIR/state" >/tmp/sk-up-promotion-list.out
"$TMPDIR/sk-up" -I owner/repo --dry-run --json --agents-home "$TMPDIR/agents" --cache-dir "$TMPDIR/cache" --state-dir "$TMPDIR/state" >/tmp/sk-up-promotion-dryrun.json

grep -Fq 'sk-up -g' /tmp/sk-up-promotion-help.out
grep -Fxq 'alpha' /tmp/sk-up-promotion-list.out
python3 - <<'PY'
import json
import pathlib

doc = json.loads(pathlib.Path("/tmp/sk-up-promotion-dryrun.json").read_text(encoding="utf-8"))
if not doc.get("ok") or not doc.get("dryRun"):
    raise SystemExit("dry-run JSON missing ok/dryRun")
actions = doc.get("actions") or []
if not actions or actions[0].get("action") != "install-source":
    raise SystemExit("dry-run JSON missing install-source action")
PY

rm /tmp/sk-up-promotion-release.out /tmp/sk-up-promotion-help.out /tmp/sk-up-promotion-list.out /tmp/sk-up-promotion-dryrun.json
