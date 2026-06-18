#!/usr/bin/env node

import { spawnSync } from "node:child_process";
import { existsSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const repoRoot = dirname(dirname(fileURLToPath(import.meta.url)));
const hookPath = join(repoRoot, ".agents", "skills", "impeccable", "scripts", "hook.mjs");

if (!existsSync(hookPath)) process.exit(0);

const result = spawnSync(process.execPath, [hookPath], {
  env: process.env,
  input: process.stdin.isTTY ? undefined : await readStdin(),
  stdio: [process.stdin.isTTY ? "inherit" : "pipe", "inherit", "inherit"],
});

if (result.error) {
  console.error(result.error.message);
  process.exit(1);
}

process.exit(result.status ?? (result.signal ? 1 : 0));

async function readStdin() {
  const chunks = [];
  for await (const chunk of process.stdin) chunks.push(chunk);
  return Buffer.concat(chunks);
}
