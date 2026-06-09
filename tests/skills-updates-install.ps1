Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$scriptPath = Join-Path $repoRoot "bin\skills-updates.ps1"
$caseRoot = Join-Path ([IO.Path]::GetTempPath()) ("sk-up-install-test-{0}" -f [IO.Path]::GetRandomFileName())
$fakeBin = Join-Path $caseRoot "bin"
$agentsHome = Join-Path $caseRoot "agents"

New-Item -ItemType Directory -Path $fakeBin, $agentsHome -Force | Out-Null
try {
    Set-Content -LiteralPath (Join-Path $agentsHome ".skill-lock.json") -Value '{"version":1,"skills":{}}' -Encoding UTF8
    Set-Content -LiteralPath (Join-Path $fakeBin "pnpm.ps1") -Encoding UTF8 -Value @'
Write-Host ("FAKE_PNPM_ARGS:" + (($args | ForEach-Object { " [$_]" }) -join ""))
exit 0
'@

    $command = @"
`$env:PATH = '$fakeBin;' + `$env:PATH
`$env:AGENTS_HOME = '$agentsHome'
& '$scriptPath' --cmd-name sk-up -i https://example.com/test.git
exit `$LASTEXITCODE
"@
    $encodedCommand = [Convert]::ToBase64String([Text.Encoding]::Unicode.GetBytes($command))
    $startInfo = [Diagnostics.ProcessStartInfo]::new()
    $startInfo.FileName = "powershell.exe"
    $startInfo.Arguments = "-NoProfile -ExecutionPolicy Bypass -EncodedCommand $encodedCommand"
    $startInfo.UseShellExecute = $false
    $startInfo.RedirectStandardOutput = $true
    $startInfo.RedirectStandardError = $true

    $process = [Diagnostics.Process]::Start($startInfo)
    if (-not $process.WaitForExit(8000)) {
        $process.Kill()
        $stdout = $process.StandardOutput.ReadToEnd()
        throw "skills-updates install timed out before pnpm completed. stdout: $stdout"
    }

    $stdout = $process.StandardOutput.ReadToEnd()
    $stderr = $process.StandardError.ReadToEnd()
    if ($process.ExitCode -ne 0) {
        throw "skills-updates install exited $($process.ExitCode). stdout: $stdout stderr: $stderr"
    }

    if ($stdout -notmatch "FAKE_PNPM_ARGS:\s+\[dlx\]\s+\[skills@latest\]\s+\[add\]\s+\[https://example\.com/test\.git\]") {
        throw "skills-updates install did not invoke pnpm with expected args. stdout: $stdout"
    }
} finally {
    if (Test-Path -LiteralPath $caseRoot) {
        Remove-Item -LiteralPath $caseRoot -Recurse -Force -ErrorAction SilentlyContinue
    }
}
