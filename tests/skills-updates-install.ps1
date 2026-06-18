Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$skUpWrapperPath = Join-Path $repoRoot "bin\sk-up.cmd"
$skillsUpdatesWrapperPath = Join-Path $repoRoot "bin\skills-updates.cmd"
$caseRoot = Join-Path ([IO.Path]::GetTempPath()) ("sk-up-install-test-{0}" -f [IO.Path]::GetRandomFileName())
$fakeBin = Join-Path $caseRoot "bin"
$agentsHome = Join-Path $caseRoot "agents"

function Invoke-CapturedProcess {
    param(
        [string]$FileName,
        [string]$Arguments,
        [hashtable]$Environment = @{},
        [int]$TimeoutMilliseconds = 8000
    )

    $startInfo = [Diagnostics.ProcessStartInfo]::new()
    $startInfo.FileName = $FileName
    $startInfo.Arguments = $Arguments
    $startInfo.UseShellExecute = $false
    $startInfo.RedirectStandardOutput = $true
    $startInfo.RedirectStandardError = $true
    foreach ($entry in $Environment.GetEnumerator()) {
        $startInfo.EnvironmentVariables[$entry.Key] = [string]$entry.Value
    }

    $process = [Diagnostics.Process]::Start($startInfo)
    if (-not $process.WaitForExit($TimeoutMilliseconds)) {
        $process.Kill()
        $stdout = $process.StandardOutput.ReadToEnd()
        throw "$FileName $Arguments timed out. stdout: $stdout"
    }

    return [pscustomobject]@{
        ExitCode = $process.ExitCode
        Stdout = $process.StandardOutput.ReadToEnd()
        Stderr = $process.StandardError.ReadToEnd()
    }
}

New-Item -ItemType Directory -Path $fakeBin, $agentsHome -Force | Out-Null
try {
    $skUpHelp = Invoke-CapturedProcess -FileName "cmd.exe" -Arguments "/d /c ""$skUpWrapperPath"" -h"
    if ($skUpHelp.ExitCode -ne 0 -or $skUpHelp.Stdout -notmatch "sk-up -g") {
        throw "sk-up wrapper help did not expose short aliases. stdout: $($skUpHelp.Stdout) stderr: $($skUpHelp.Stderr)"
    }

    $skillsUpdatesHelp = Invoke-CapturedProcess -FileName "cmd.exe" -Arguments "/d /c ""$skillsUpdatesWrapperPath"" --help"
    if ($skillsUpdatesHelp.ExitCode -ne 0 -or $skillsUpdatesHelp.Stdout -notmatch "skills-updates --global") {
        throw "skills-updates wrapper help did not expose long aliases. stdout: $($skillsUpdatesHelp.Stdout) stderr: $($skillsUpdatesHelp.Stderr)"
    }

    Set-Content -LiteralPath (Join-Path $agentsHome ".skill-lock.json") -Value '{"version":1,"skills":{}}' -Encoding UTF8
    Set-Content -LiteralPath (Join-Path $fakeBin "pnpm.ps1") -Encoding UTF8 -Value @'
Write-Host ("FAKE_PNPM_ARGS:" + (($args | ForEach-Object { " [$_]" }) -join ""))
exit 0
'@

    function Assert-SkUpInstall {
        param([string]$SourceUrl)

        $installArgs = '/d /s /c ""' + $skUpWrapperPath + '" -i ' + $SourceUrl + '"'
        $installResult = Invoke-CapturedProcess -FileName "cmd.exe" -Arguments $installArgs -Environment @{
            PATH = "$fakeBin;$env:PATH"
            AGENTS_HOME = $agentsHome
        }
        if ($installResult.ExitCode -ne 0) {
            throw "skills-updates install exited $($installResult.ExitCode). stdout: $($installResult.Stdout) stderr: $($installResult.Stderr)"
        }

        $escapedSourceUrl = [regex]::Escape($SourceUrl)
        if ($installResult.Stdout -notmatch "FAKE_PNPM_ARGS:\s+\[dlx\]\s+\[skills@latest\]\s+\[add\]\s+\[$escapedSourceUrl\]") {
            throw "skills-updates install did not invoke pnpm with expected args. stdout: $($installResult.Stdout)"
        }
    }

    Assert-SkUpInstall -SourceUrl "https://example.com/test.git"
    Assert-SkUpInstall -SourceUrl "https://example.com/a%2Fb.git"
} finally {
    if (Test-Path -LiteralPath $caseRoot) {
        Remove-Item -LiteralPath $caseRoot -Recurse -Force -ErrorAction SilentlyContinue
    }
}
