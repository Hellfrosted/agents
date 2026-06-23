Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$caseRoot = Join-Path ([IO.Path]::GetTempPath()) ("sk-up-install-test-{0}" -f [IO.Path]::GetRandomFileName())
$caseBin = Join-Path $caseRoot "bin"
$agentsHome = Join-Path $caseRoot "agents"
$cacheDir = Join-Path $caseRoot "cache"
$stateDir = Join-Path $caseRoot "state"

function Invoke-CapturedProcess {
    param(
        [string]$FileName,
        [string]$Arguments,
        [int]$TimeoutMilliseconds = 15000
    )

    $startInfo = [Diagnostics.ProcessStartInfo]::new()
    $startInfo.FileName = $FileName
    $startInfo.Arguments = $Arguments
    $startInfo.WorkingDirectory = $repoRoot
    $startInfo.UseShellExecute = $false
    $startInfo.RedirectStandardOutput = $true
    $startInfo.RedirectStandardError = $true

    $process = [Diagnostics.Process]::Start($startInfo)
    if (-not $process.WaitForExit($TimeoutMilliseconds)) {
        $process.Kill()
        $stdout = $process.StandardOutput.ReadToEnd()
        $stderr = $process.StandardError.ReadToEnd()
        throw "$FileName $Arguments timed out. stdout: $stdout stderr: $stderr"
    }

    return [pscustomobject]@{
        ExitCode = $process.ExitCode
        Stdout = $process.StandardOutput.ReadToEnd()
        Stderr = $process.StandardError.ReadToEnd()
    }
}

function Quote-CmdArgument {
    param(
        [string]$Value
    )

    if ($Value -notmatch '[\s"]') {
        return $Value
    }

    return '"' + ($Value -replace '"', '""') + '"'
}

function Invoke-CmdScript {
    param(
        [string]$ScriptPath,
        [string[]]$Arguments,
        [int]$TimeoutMilliseconds = 15000
    )

    $quoted = @($ScriptPath) + $Arguments |
        ForEach-Object { Quote-CmdArgument -Value $_ }
    $command = [string]::Join(" ", $quoted)

    Invoke-CapturedProcess `
        -FileName "cmd.exe" `
        -Arguments "/d /s /c ""$command""" `
        -TimeoutMilliseconds $TimeoutMilliseconds
}

function Assert-Ok {
    param(
        [object]$Result,
        [string]$Label
    )

    if ($Result.ExitCode -ne 0) {
        throw "$Label exited $($Result.ExitCode). stdout: $($Result.Stdout) stderr: $($Result.Stderr)"
    }
}

New-Item -ItemType Directory -Path $caseBin, $agentsHome, $cacheDir, $stateDir -Force | Out-Null
try {
    Copy-Item -LiteralPath (Join-Path $repoRoot "bin\sk-up.cmd") -Destination (Join-Path $caseBin "sk-up.cmd")
    Copy-Item -LiteralPath (Join-Path $repoRoot "bin\skills-updates.cmd") -Destination (Join-Path $caseBin "skills-updates.cmd")

    $buildSkUp = Invoke-CapturedProcess -FileName "go" -Arguments "build -o ""$caseBin\sk-up.exe"" .\cmd\sk-up"
    Assert-Ok -Result $buildSkUp -Label "go build sk-up"
    if (Test-Path -LiteralPath (Join-Path $caseBin "skills-updates.exe")) {
        throw "Windows install test must not create skills-updates.exe"
    }

    $skUpHelp = Invoke-CmdScript -ScriptPath (Join-Path $caseBin "sk-up.cmd") -Arguments @("-h")
    Assert-Ok -Result $skUpHelp -Label "sk-up help"
    if ($skUpHelp.Stdout -notmatch "sk-up -g") {
        throw "sk-up wrapper help did not expose short aliases. stdout: $($skUpHelp.Stdout)"
    }

    $skillsUpdatesHelp = Invoke-CmdScript -ScriptPath (Join-Path $caseBin "skills-updates.cmd") -Arguments @("--help")
    Assert-Ok -Result $skillsUpdatesHelp -Label "skills-updates help"
    if ($skillsUpdatesHelp.Stdout -notmatch "skills-updates --global") {
        throw "skills-updates wrapper help did not expose long aliases. stdout: $($skillsUpdatesHelp.Stdout)"
    }

    $dryRun = Invoke-CmdScript -ScriptPath (Join-Path $caseBin "sk-up.cmd") -Arguments @(
        "-I",
        "owner/repo",
        "--dry-run",
        "--json",
        "--agents-home",
        $agentsHome,
        "--cache-dir",
        $cacheDir,
        "--state-dir",
        $stateDir
    )
    Assert-Ok -Result $dryRun -Label "sk-up install-source dry-run"
    $json = $dryRun.Stdout | ConvertFrom-Json
    if (-not $json.ok -or -not $json.dryRun -or $json.actions[0].action -ne "install-source") {
        throw "install-source dry-run returned unexpected JSON: $($dryRun.Stdout)"
    }
} finally {
    if (Test-Path -LiteralPath $caseRoot) {
        Remove-Item -LiteralPath $caseRoot -Recurse -Force -ErrorAction SilentlyContinue
    }
}
