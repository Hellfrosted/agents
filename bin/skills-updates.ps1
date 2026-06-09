Set-StrictMode -Version Latest
$ErrorActionPreference = "Continue"

function Set-Utf8ConsoleEncoding {
    $utf8NoBom = [Text.UTF8Encoding]::new($false)

    try {
        [Console]::InputEncoding = $utf8NoBom
    } catch {}

    try {
        [Console]::OutputEncoding = $utf8NoBom
    } catch {}

    $script:OutputEncoding = $utf8NoBom
}

Set-Utf8ConsoleEncoding

$mode = "help"
$target = ""
$targets = New-Object System.Collections.Generic.List[string]
$scope = ""
$globalOptionUsed = $false
$displayName = "skills-updates"

if ($args.Count -ge 2 -and $args[0] -eq "--cmd-name") {
    $displayName = $args[1]
    $args = @($args | Select-Object -Skip 2)
}

function Show-Help {
    function Write-HelpCommand {
        param(
            [string]$Command,
            [string]$Description
        )

        Write-Host ("  {0,-42}" -f $Command) -ForegroundColor Green -NoNewline
        Write-Host $Description -ForegroundColor Gray
    }

    function Write-HelpNote {
        param([string]$Message)

        Write-Host "  - $Message" -ForegroundColor DarkGray
    }

    Write-Host "Usage:" -ForegroundColor Cyan
    if ($displayName -eq "sk-up") {
        Write-HelpCommand $displayName "Show help"
        Write-HelpCommand "$displayName -h" "Show help"
        Write-HelpCommand "$displayName -l" "List installed skills without checking upstream"
        Write-Host ""
        Write-Host "Check:" -ForegroundColor Cyan
        Write-HelpCommand "$displayName -g" "List global skill status"
        Write-HelpCommand "$displayName -d <skill>" "Show terminal diff for one skill"
        Write-HelpCommand "$displayName -z [skill1 skill2 ...]" "Open Zed diff viewer"
        Write-Host ""
        Write-Host "Install:" -ForegroundColor Cyan
        Write-HelpCommand "$displayName -i" "Install changed skills"
        Write-HelpCommand "$displayName -i [skill1 skill2 ...]" "Install one or more skills"
        Write-HelpCommand "$displayName -i <source-url>" "Install a skill package URL without wiping the lockfile"
        Write-Host ""
        Write-Host "Uninstall:" -ForegroundColor Cyan
        Write-HelpCommand "$displayName -r <skill1 skill2 ...>" "Uninstall one or more global skills"
        Write-Host ""
        Write-Host "Skip:" -ForegroundColor Cyan
        Write-HelpCommand "$displayName -s <skill>" "Skip current upstream diff"
        Write-HelpCommand "$displayName -u <skill>" "Remove saved skip"
        Write-HelpCommand "$displayName -S" "List saved skips"
        Write-Host ""
        Write-Host "Notes:" -ForegroundColor Cyan
        Write-HelpNote "Checks compare installed skill folders against upstream content, not just lockfile hashes."
        Write-HelpNote "Source repos are cached locally and fetched in parallel."
        Write-HelpNote "Skips are tied to the current upstream tree hash and expire when upstream changes."
        Write-HelpNote "Named installs run: pnpm dlx skills@latest add <source> -g -y --agent universal --skill <skill-name>"
        Write-HelpNote "Source installs run: pnpm dlx skills@latest add <source> -g -y --agent universal"
        Write-HelpNote "Uninstalls run: pnpm dlx skills@latest remove -g -y --agent universal --skill <skill-name>"
        Write-HelpNote "Uninstalls also remove the global installed skill directory, saved skip, and lockfile entry."
        return
    }

    Write-HelpCommand $displayName "Show help"
    Write-HelpCommand "$displayName --help" "Show help"
    Write-HelpCommand "$displayName --list" "List installed skills without checking upstream"
    Write-Host ""
    Write-Host "Check:" -ForegroundColor Cyan
    Write-HelpCommand "$displayName --global" "List global skill status"
    Write-HelpCommand "$displayName --diff <skill>" "Show terminal diff for one skill"
    Write-HelpCommand "$displayName --zed [skill1 ...]" "Open Zed diff viewer"
    Write-Host ""
    Write-Host "Install:" -ForegroundColor Cyan
    Write-HelpCommand "$displayName --install" "Install changed skills"
    Write-HelpCommand "$displayName --install [skill1 ...]" "Install named skills"
    Write-HelpCommand "$displayName --install <source-url>" "Install a skill package URL without wiping the lockfile"
    Write-Host ""
    Write-Host "Uninstall:" -ForegroundColor Cyan
    Write-HelpCommand "$displayName --remove <skill1 ...>" "Uninstall named global skills"
    Write-Host ""
    Write-Host "Skip:" -ForegroundColor Cyan
    Write-HelpCommand "$displayName --skip <skill>" "Skip current upstream diff"
    Write-HelpCommand "$displayName --unskip <skill>" "Remove saved skip"
    Write-HelpCommand "$displayName --skips" "List saved skips"
    Write-Host ""
    Write-Host "Notes:" -ForegroundColor Cyan
    Write-HelpNote "Checks compare installed skill folders against upstream content, not just lockfile hashes."
    Write-HelpNote "Source repos are cached locally and fetched in parallel."
    Write-HelpNote "Skips are tied to the current upstream tree hash and expire when upstream changes."
    Write-HelpNote "Named installs run: pnpm dlx skills@latest add <source> -g -y --agent universal --skill <skill-name>"
    Write-HelpNote "Source installs run: pnpm dlx skills@latest add <source> -g -y --agent universal"
    Write-HelpNote "Uninstalls run: pnpm dlx skills@latest remove -g -y --agent universal --skill <skill-name>"
    Write-HelpNote "Uninstalls also remove the global installed skill directory, saved skip, and lockfile entry."
}

function Write-ColoredLine {
    param(
        [string]$Message,
        [ConsoleColor]$Color = [ConsoleColor]::Gray
    )

    Write-Host $Message -ForegroundColor $Color
}

function Write-StatusLine {
    param([string]$Message)

    switch -Regex ($Message) {
        '^\s*(\d+\.\s+)?OK\s+' { Write-ColoredLine $Message Green; break }
        '^\s*(\d+\.\s+)?UPDATE\s+' { Write-ColoredLine $Message Yellow; break }
        '^\s*(\d+\.\s+)?SKIP\s+' { Write-ColoredLine $Message DarkYellow; break }
        '^\s*(\d+\.\s+)?MISSING\s+' { Write-ColoredLine $Message Magenta; break }
        '^\s*(\d+\.\s+)?ERROR\s+' { Write-ColoredLine $Message Red; break }
        '^\s*(\d+\.\s+)?DIFF\s+' { Write-ColoredLine $Message Yellow; break }
        '^\s*(\d+\.\s+)?INSTALL\s+' { Write-ColoredLine $Message Cyan; break }
        '^\s*(\d+\.\s+)?UNINSTALL\s+' { Write-ColoredLine $Message Cyan; break }
        '^\s*(\d+\.\s+)?CHECK\s+' { Write-ColoredLine $Message Cyan; break }
        '^\s*(\d+\.\s+)?(FETCH|CLONE|COMPARE|READY)\s+' { Write-ColoredLine $Message DarkGray; break }
        default { Write-Output $Message }
    }
}

function Write-CliError {
    param([string]$Message)

    $shortMessage = $Message -replace "^skills-updates:\s*", ""
    $line = "ERROR   $shortMessage"
    $previousColor = [Console]::ForegroundColor
    try {
        [Console]::ForegroundColor = [ConsoleColor]::Red
        [Console]::Error.WriteLine($line)
    } finally {
        [Console]::ForegroundColor = $previousColor
    }
}

function Test-TargetMatch {
    param([string]$Name)

    return $targets.Count -eq 0 -or $targets.Contains($Name)
}

function Test-InstallSourceArgument {
    param([string]$Value)

    return (
        $Value -match "^(https?|ssh)://" -or
        $Value -match "^git@" -or
        $Value -match "\.git($|[#?])" -or
        $Value -match "^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+($|/)"
    )
}

foreach ($arg in $args) {
    switch ($arg) {
        { $_ -in @("--global", "-g") } {
            $globalOptionUsed = $true
            $scope = "global"
            if ($mode -eq "help") {
                $mode = "summary"
            }
            continue
        }
        { $_ -in @("--list", "-l") } {
            $mode = "list"
            continue
        }
        { $_ -in @("--diff", "-d") } {
            $scope = "global"
            $mode = "diff"
            continue
        }
        { $_ -in @("--zed", "--gui", "-z") } {
            $scope = "global"
            $mode = "zed"
            continue
        }
        { $_ -in @("--install-all", "install-all") } {
            $mode = "install-all"
            continue
        }
        { $_ -in @("-i", "--install") } {
            $mode = "install"
            continue
        }
        { $_ -in @("-r", "--remove", "--uninstall") } {
            $mode = "uninstall"
            continue
        }
        { @("-s", "--skip") -ccontains $_ } {
            $mode = "skip"
            continue
        }
        { @("-u", "--unskip") -ccontains $_ } {
            $mode = "unskip"
            continue
        }
        { @("-S", "--skips") -ccontains $_ } {
            $mode = "skips"
            continue
        }
        { $_ -in @("--help", "-h") } {
            Show-Help
            exit 0
        }
        default {
            if ($arg.StartsWith("-")) {
                Write-CliError "skills-updates: unknown option: $arg"
                exit 1
            }
            $targets.Add($arg)
        }
    }
}

if ($targets.Count -eq 1) {
    $target = $targets[0]
}

if ($mode -eq "help") {
    Show-Help
    exit 0
}

if ($scope -and $scope -ne "global") {
    Write-CliError "skills-updates: only global scope is supported"
    exit 1
}

if ($scope -eq "global" -and $mode -in @("install", "install-all")) {
    Write-CliError "skills-updates: -g cannot be used with -i; use $displayName -i instead"
    exit 1
}

if ($globalOptionUsed -and $mode -eq "zed") {
    Write-CliError "skills-updates: -g cannot be used with -z; use $displayName -z [skill] instead"
    exit 1
}

if ($globalOptionUsed -and $mode -eq "diff") {
    Write-CliError "skills-updates: -g cannot be used with -d; use $displayName -d <skill> instead"
    exit 1
}

if ($globalOptionUsed -and $mode -eq "summary" -and $targets.Count -gt 0) {
    Write-CliError "skills-updates: -g does not accept a skill name; use $displayName -d <skill> instead"
    exit 1
}

if (-not $scope -and $mode -notin @("list", "install", "install-all", "uninstall", "skip", "unskip", "skips")) {
    Write-CliError "skills-updates: use -g to check global skills"
    Write-Output ""
    Show-Help
    exit 1
}

if ($mode -eq "diff" -and $targets.Count -ne 1) {
    Write-CliError "skills-updates: diff mode requires exactly one skill name"
    exit 1
}

if ($mode -eq "install" -and $targets.Count -eq 0) {
    $mode = "install-all"
}

if ($mode -eq "uninstall" -and $targets.Count -eq 0) {
    Write-CliError "skills-updates: uninstall requires at least one skill name"
    exit 1
}

if ($mode -in @("skip", "unskip") -and $targets.Count -ne 1) {
    Write-CliError "skills-updates: $mode requires exactly one skill name"
    exit 1
}

$agentsHome = if ($env:AGENTS_HOME) {
    $env:AGENTS_HOME
} elseif ($env:USERPROFILE) {
    Join-Path $env:USERPROFILE ".agents"
} else {
    Join-Path $HOME ".agents"
}

$lock = Join-Path $agentsHome ".skill-lock.json"
$skillsDir = Join-Path $agentsHome "skills"
$stateRoot = if ($env:LOCALAPPDATA) {
    Join-Path $env:LOCALAPPDATA "skills-updates"
} else {
    Join-Path ([IO.Path]::GetTempPath()) "skills-updates-state"
}
$skipFile = Join-Path $stateRoot "skips.json"
$repoRoot = Join-Path $stateRoot "repos"
$maxRepoJobs = 24
$compareTempPaths = New-Object System.Collections.Generic.List[string]
$compareTempFiles = New-Object System.Collections.Generic.List[string]
$repoRunspacePool = $null
$preserveCompareTempPaths = $false

function Show-InstalledSkills {
    if (-not (Test-Path -LiteralPath $skillsDir -PathType Container)) {
        Write-StatusLine "OK      no installed skills found"
        return
    }

    $skillDirs = @(Get-ChildItem -LiteralPath $skillsDir -Directory -ErrorAction SilentlyContinue | Sort-Object -Property Name)
    if ($skillDirs.Count -eq 0) {
        Write-StatusLine "OK      no installed skills found"
        return
    }

    foreach ($skillDir in $skillDirs) {
        Write-Output $skillDir.Name
    }
}

function Get-CacheKey {
    param([string]$Value)

    $sha = [Security.Cryptography.SHA256]::Create()
    try {
        $bytes = [Text.Encoding]::UTF8.GetBytes($Value)
        $hash = $sha.ComputeHash($bytes)
        return -join ($hash | ForEach-Object { $_.ToString("x2") })
    } finally {
        $sha.Dispose()
    }
}

function Get-SkillLockMutexName {
    $lockKey = try {
        [IO.Path]::GetFullPath($lock).ToLowerInvariant()
    } catch {
        $lock.ToLowerInvariant()
    }

    return "Local\skills-updates-lock-$((Get-CacheKey $lockKey).Substring(0, 32))"
}

function Invoke-WithSkillLockMutex {
    param(
        [scriptblock]$ScriptBlock,
        [ref]$Completed
    )

    $mutex = [Threading.Mutex]::new($false, (Get-SkillLockMutexName))
    $hasLock = $false
    if ($Completed) {
        $Completed.Value = $false
    }

    try {
        $hasLock = $mutex.WaitOne([TimeSpan]::FromMinutes(10))
        if (-not $hasLock) {
            Write-CliError "skills-updates: timed out waiting for lockfile guard: $lock"
            return
        }

        if ($Completed) {
            $Completed.Value = $true
        }
        & $ScriptBlock
    } finally {
        if ($hasLock) {
            $mutex.ReleaseMutex()
        }
        $mutex.Dispose()
    }
}

function Invoke-WithSkillStateTransaction {
    param(
        [scriptblock]$ScriptBlock,
        [ref]$Completed
    )

    $transactionScriptBlock = $ScriptBlock
    Invoke-WithSkillLockMutex -Completed $Completed -ScriptBlock {
        $snapshot = Read-SkillLockSnapshot
        Save-SkillLockBackup -Snapshot $snapshot
        & $transactionScriptBlock $snapshot
    }
}

function Complete-SkillStateTransaction {
    Remove-SkillLockBackup
}

function Restore-SkillStateAfterSkillsCommand {
    param([object]$BeforeSnapshot)

    Restore-SkillLockAfterPnpx -BeforeSnapshot $BeforeSnapshot
}

function Restore-SkillStateSnapshot {
    param([object]$Snapshot)

    return Restore-SkillLockSnapshotExact -Snapshot $Snapshot
}

function Invoke-SkillsAdd {
    param(
        [string]$SourceUrl,
        [string]$Name = ""
    )

    if (-not (Get-Command pnpm -ErrorAction SilentlyContinue)) {
        Write-CliError "skills-updates: pnpm is required to install skills"
        return $false
    }

    if ($Name) {
        Write-StatusLine "INSTALL $Name from $SourceUrl"
    } else {
        Write-StatusLine "INSTALL $SourceUrl"
    }

    $skillsArgs = New-Object System.Collections.Generic.List[string]
    $skillsArgs.Add("skills@latest")
    $skillsArgs.Add("add")
    $skillsArgs.Add($SourceUrl)
    $skillsArgs.Add("-g")
    $skillsArgs.Add("-y")
    $skillsArgs.Add("--agent")
    $skillsArgs.Add("universal")
    if ($Name) {
        $skillsArgs.Add("--skill")
        $skillsArgs.Add($Name)
    }

    $status = [pscustomobject]@{ Ok = $false }
    $completed = $false
    Invoke-WithSkillStateTransaction -Completed ([ref]$completed) -ScriptBlock {
        param([object]$lockBeforeInstall)

        & pnpm dlx @skillsArgs | ForEach-Object { Write-Host $_ }
        $status.Ok = $LASTEXITCODE -eq 0
        Restore-SkillStateAfterSkillsCommand -BeforeSnapshot $lockBeforeInstall
        Complete-SkillStateTransaction
    }
    return $completed -and $status.Ok
}

function Install-Skill {
    param(
        [string]$Name,
        [string]$SourceUrl
    )

    return Invoke-SkillsAdd -SourceUrl $SourceUrl -Name $Name
}

function Install-SkillSource {
    param([string]$SourceUrl)

    return Invoke-SkillsAdd -SourceUrl $SourceUrl
}

function Get-InstalledSkillDirectory {
    param([string]$Name)

    if ([string]::IsNullOrWhiteSpace($Name) -or $Name -in @(".", "..") -or (Split-Path -Leaf $Name) -ne $Name) {
        Write-CliError "skills-updates: invalid skill name for uninstall: $Name"
        return $null
    }

    $baseDir = [IO.Path]::GetFullPath($skillsDir)
    $candidateDir = [IO.Path]::GetFullPath((Join-Path $skillsDir $Name))
    $basePrefix = if (
        $baseDir.EndsWith([IO.Path]::DirectorySeparatorChar) -or
        $baseDir.EndsWith([IO.Path]::AltDirectorySeparatorChar)
    ) {
        $baseDir
    } else {
        "$baseDir$([IO.Path]::DirectorySeparatorChar)"
    }

    if (-not $candidateDir.StartsWith($basePrefix, [StringComparison]::OrdinalIgnoreCase)) {
        Write-CliError "skills-updates: invalid skill name for uninstall: $Name"
        return $null
    }

    return $candidateDir
}

function Remove-InstalledSkillDirectory {
    param([string]$Name)

    $installedDir = Get-InstalledSkillDirectory -Name $Name
    if (-not $installedDir) {
        return $false
    }

    if (-not (Test-Path -LiteralPath $installedDir -PathType Container)) {
        Write-StatusLine "OK      no installed directory for $Name"
        return $true
    }

    try {
        Remove-Item -LiteralPath $installedDir -Recurse -Force -ErrorAction Stop
    } catch {
        Write-CliError "skills-updates: could not remove installed directory for $Name`: $($_.Exception.Message)"
        return $false
    }

    Write-StatusLine "OK      removed installed directory for $Name"
    return $true
}

function Uninstall-Skill {
    param([string]$Name)

    if (-not (Get-Command pnpm -ErrorAction SilentlyContinue)) {
        Write-CliError "skills-updates: pnpm is required to uninstall skills"
        return $false
    }

    Write-StatusLine "UNINSTALL $Name"
    $status = [pscustomobject]@{
        Ok = $false
        LockEntryRemoved = $false
    }
    $completed = $false
    Invoke-WithSkillStateTransaction -Completed ([ref]$completed) -ScriptBlock {
        param([object]$lockBeforeUninstall)

        & pnpm dlx skills@latest remove -g -y --agent universal --skill $Name | ForEach-Object { Write-Host $_ }
        $exitCode = Get-Variable -Name LASTEXITCODE -Scope Global -ErrorAction SilentlyContinue
        if (-not $exitCode) {
            $status.Ok = $true
        } else {
            $status.Ok = $exitCode.Value -eq 0
        }

        if (-not $status.Ok) {
            Restore-SkillStateAfterSkillsCommand -BeforeSnapshot $lockBeforeUninstall
            Complete-SkillStateTransaction
            return
        }

        if (-not (Remove-InstalledSkillDirectory -Name $Name)) {
            $status.Ok = $false
            if (Restore-SkillStateSnapshot -Snapshot $lockBeforeUninstall) {
                Complete-SkillStateTransaction
            }
            return
        }

        $state = Read-SkillLock
        if (Remove-SkillLockEntry -State $state -Name $Name) {
            try {
                Write-SkillLock -State $state
                $writtenState = Read-SkillLock
                if (Test-SkillLockEntry -State $writtenState -Name $Name) {
                    throw "lock entry remained after write"
                }
                $status.LockEntryRemoved = $true
            } catch {
                Write-CliError "skills-updates: could not remove $Name from lockfile: $($_.Exception.Message)"
                $status.Ok = $false
                if (Restore-SkillStateSnapshot -Snapshot $lockBeforeUninstall) {
                    Complete-SkillStateTransaction
                }
                return
            }
        }

        Remove-SavedSkillSkip -Name $Name | Out-Null
        Complete-SkillStateTransaction
    }

    if (-not ($completed -and $status.Ok)) {
        return $false
    }

    if ($status.LockEntryRemoved) {
        Write-StatusLine "OK      removed $Name from lockfile"
    } else {
        Write-StatusLine "OK      no lock entry for $Name"
    }

    return $true
}

function Read-SkillLock {
    Repair-SkillLockFromBackup | Out-Null

    if (-not (Test-Path -LiteralPath $lock -PathType Leaf)) {
        return $null
    }

    try {
        return Read-RawSkillLock | ConvertFrom-Json
    } catch {
        Write-CliError "skills-updates: could not read lockfile: $lock"
        return $null
    }
}

function Read-SkillLockSnapshot {
    Repair-SkillLockFromBackup | Out-Null

    $snapshot = [pscustomobject]@{
        Exists = $false
        Raw = $null
        State = $null
    }

    if (-not (Test-Path -LiteralPath $lock -PathType Leaf)) {
        return $snapshot
    }

    $snapshot.Exists = $true
    try {
        $snapshot.Raw = Read-RawSkillLock
        if ($snapshot.Raw) {
            $snapshot.State = $snapshot.Raw | ConvertFrom-Json
        }
    } catch {
        Write-CliError "skills-updates: could not snapshot lockfile before install: $lock"
    }

    return $snapshot
}

function Get-SkillLockWritePath {
    try {
        $item = Get-Item -LiteralPath $lock -ErrorAction SilentlyContinue
        if ($item -and $item.LinkType -eq "SymbolicLink" -and $item.Target -and $item.Target.Count -gt 0) {
            $target = [string]$item.Target[0]
            if ([IO.Path]::IsPathRooted($target)) {
                return $target
            }

            return Join-Path (Split-Path -Parent $lock) $target
        }
    } catch {
        return $lock
    }

    return $lock
}

function Get-SkillLockBackupPath {
    return "$lock.sk-up-backup"
}

function Test-RawJsonObject {
    param([AllowNull()][string]$Raw)

    if ([string]::IsNullOrWhiteSpace($Raw)) {
        return $false
    }

    try {
        $state = $Raw | ConvertFrom-Json
        return $state -is [pscustomobject]
    } catch {
        return $false
    }
}

function Read-RawSkillLock {
    if (-not (Test-Path -LiteralPath $lock -PathType Leaf)) {
        return $null
    }

    return Get-Content -LiteralPath $lock -Raw
}

function Save-SkillLockBackup {
    param([object]$Snapshot)

    if (-not $Snapshot -or -not $Snapshot.Exists -or -not (Test-RawJsonObject -Raw $Snapshot.Raw)) {
        return
    }

    $backupPath = Get-SkillLockBackupPath
    $backupDir = Split-Path -Parent $backupPath
    New-Item -ItemType Directory -Path $backupDir -Force | Out-Null
    $tempBackup = Join-Path $backupDir (".skill-lock.json.sk-up-backup.tmp-{0}" -f [IO.Path]::GetRandomFileName())
    try {
        Write-Utf8NoBomFile -Path $tempBackup -Raw $Snapshot.Raw
        Move-Item -LiteralPath $tempBackup -Destination $backupPath -Force
    } finally {
        if (Test-Path -LiteralPath $tempBackup -PathType Leaf) {
            Remove-Item -LiteralPath $tempBackup -Force -ErrorAction SilentlyContinue
        }
    }
}

function Remove-SkillLockBackup {
    $backupPath = Get-SkillLockBackupPath
    if (Test-Path -LiteralPath $backupPath -PathType Leaf) {
        Remove-Item -LiteralPath $backupPath -Force -ErrorAction SilentlyContinue
    }
}

function Repair-SkillLockFromBackup {
    $backupPath = Get-SkillLockBackupPath
    if (-not (Test-Path -LiteralPath $backupPath -PathType Leaf)) {
        return $false
    }

    $currentRaw = Read-RawSkillLock
    if (Test-RawJsonObject -Raw $currentRaw) {
        return $false
    }

    $backupRaw = Get-Content -LiteralPath $backupPath -Raw
    if (-not (Test-RawJsonObject -Raw $backupRaw)) {
        return $false
    }

    Write-RawSkillLock -Raw $backupRaw
    Write-StatusLine "OK      restored lockfile from interrupted sk-up backup"
    return $true
}

function Write-Utf8NoBomFile {
    param(
        [string]$Path,
        [AllowNull()][string]$Raw
    )

    if ($null -eq $Raw) {
        $Raw = ""
    }

    $utf8NoBom = [Text.UTF8Encoding]::new($false)
    [IO.File]::WriteAllText($Path, $Raw, $utf8NoBom)
}

function Write-RawSkillLock {
    param([AllowNull()][string]$Raw)

    if ($null -eq $Raw) {
        $Raw = ""
    }

    $writePath = Get-SkillLockWritePath
    $lockDir = Split-Path -Parent $writePath
    New-Item -ItemType Directory -Path $lockDir -Force -ErrorAction Stop | Out-Null
    $tempLock = Join-Path $lockDir (".skill-lock.json.tmp-{0}" -f [IO.Path]::GetRandomFileName())
    try {
        Write-Utf8NoBomFile -Path $tempLock -Raw $Raw
        Move-Item -LiteralPath $tempLock -Destination $writePath -Force -ErrorAction Stop
    } finally {
        if (Test-Path -LiteralPath $tempLock -PathType Leaf) {
            Remove-Item -LiteralPath $tempLock -Force -ErrorAction SilentlyContinue
        }
    }
}

function Write-SkillLock {
    param([object]$State)

    $json = $State | ConvertTo-Json -Depth 20
    Write-RawSkillLock -Raw ($json + [Environment]::NewLine)
}

function Restore-SkillLockSnapshotExact {
    param(
        [object]$Snapshot
    )

    if (-not $Snapshot) {
        return $false
    }

    try {
        if (-not $Snapshot.Exists) {
            if (Test-Path -LiteralPath $lock -PathType Leaf) {
                Remove-Item -LiteralPath $lock -Force -ErrorAction Stop
                Write-StatusLine "OK      restored lockfile snapshot"
            }
            return $true
        }

        Write-RawSkillLock -Raw $Snapshot.Raw
        Write-StatusLine "OK      restored lockfile snapshot"
        return $true
    } catch {
        Write-CliError "skills-updates: could not restore lockfile snapshot: $($_.Exception.Message)"
        return $false
    }
}

function Set-JsonProperty {
    param(
        [object]$Object,
        [string]$Name,
        [object]$Value
    )

    if ((Get-JsonPropertyNames -Object $Object) -ccontains $Name) {
        $Object.$Name = $Value
    } else {
        $Object | Add-Member -NotePropertyName $Name -NotePropertyValue $Value
    }
}

function Get-JsonPropertyNames {
    param([object]$Object)

    if (-not $Object) {
        return @()
    }

    return @($Object.PSObject.Properties | ForEach-Object { $_.Name })
}

function Merge-JsonObjectProperties {
    param(
        [object]$Target,
        [object]$Source
    )

    if (-not $Target -or -not $Source) {
        return 0
    }

    $addedCount = 0
    foreach ($property in $Source.PSObject.Properties) {
        if (-not ((Get-JsonPropertyNames -Object $Target) -ccontains $property.Name)) {
            $Target | Add-Member -NotePropertyName $property.Name -NotePropertyValue $property.Value
            $addedCount += 1
        } elseif (
            $property.Value -is [pscustomobject] -and
            $Target.$($property.Name) -is [pscustomobject]
        ) {
            $addedCount += Merge-JsonObjectProperties -Target $Target.$($property.Name) -Source $property.Value
        }
    }

    return $addedCount
}

function Restore-SkillLockAfterPnpx {
    param(
        [object]$BeforeSnapshot
    )

    if (-not $BeforeSnapshot -or -not $BeforeSnapshot.Exists) {
        return
    }

    $BeforeState = $BeforeSnapshot.State
    $afterState = Read-SkillLock
    if (-not $afterState) {
        Write-RawSkillLock -Raw $BeforeSnapshot.Raw
        Write-StatusLine "OK      restored lockfile after skills command"
        return
    }

    if (-not $BeforeState -or -not ((Get-JsonPropertyNames -Object $BeforeState) -ccontains "skills")) {
        return
    }

    if (-not ((Get-JsonPropertyNames -Object $afterState) -ccontains "skills") -or -not $afterState.skills) {
        Set-JsonProperty -Object $afterState -Name "skills" -Value ([pscustomobject]@{})
    } elseif (-not ($afterState.skills -is [pscustomobject])) {
        Set-JsonProperty -Object $afterState -Name "skills" -Value ([pscustomobject]@{})
    }

    $restoredCount = 0
    $preservedCount = 0
    foreach ($property in $BeforeState.skills.PSObject.Properties) {
        if (-not ((Get-JsonPropertyNames -Object $afterState.skills) -ccontains $property.Name)) {
            $afterState.skills | Add-Member -NotePropertyName $property.Name -NotePropertyValue $property.Value
            $restoredCount += 1
        } elseif ($afterState.skills.$($property.Name) -is [pscustomobject]) {
            $preservedCount += Merge-JsonObjectProperties -Target $afterState.skills.$($property.Name) -Source $property.Value
        }
    }

    foreach ($property in $BeforeState.PSObject.Properties) {
        if ($property.Name -in @("version", "skills")) {
            continue
        }

        if (-not ((Get-JsonPropertyNames -Object $afterState) -ccontains $property.Name)) {
            $afterState | Add-Member -NotePropertyName $property.Name -NotePropertyValue $property.Value
            $preservedCount += 1
        } elseif (
            $property.Value -is [pscustomobject] -and
            $afterState.$($property.Name) -is [pscustomobject]
        ) {
            $preservedCount += Merge-JsonObjectProperties -Target $afterState.$($property.Name) -Source $property.Value
        }
    }

    if (($restoredCount + $preservedCount) -gt 0) {
        Write-SkillLock -State $afterState
        $messageParts = New-Object System.Collections.Generic.List[string]
        if ($restoredCount -gt 0) {
            $messageParts.Add("$restoredCount existing lockfile entr$(if ($restoredCount -eq 1) { 'y' } else { 'ies' })")
        }
        if ($preservedCount -gt 0) {
            $messageParts.Add("$preservedCount lockfile field$(if ($preservedCount -eq 1) { '' } else { 's' })")
        }
        Write-StatusLine "OK      preserved $($messageParts -join ', ')"
    }
}

function Remove-SkillLockEntry {
    param(
        [object]$State,
        [string]$Name
    )

    if (-not (Test-SkillLockEntry -State $State -Name $Name)) {
        return $false
    }

    $State.skills.PSObject.Properties.Remove($Name)
    return $true
}

function Test-SkillLockEntry {
    param(
        [object]$State,
        [string]$Name
    )

    if (-not $State -or -not ((Get-JsonPropertyNames -Object $State) -contains "skills")) {
        return $false
    }

    $names = @($State.skills.PSObject.Properties | ForEach-Object { $_.Name })
    return $names -ccontains $Name
}

function Remove-SkillLockEntryFromFile {
    param([string]$Name)

    $status = [pscustomobject]@{ Removed = $false }
    $completed = $false
    Invoke-WithSkillLockMutex -Completed ([ref]$completed) -ScriptBlock {
        $state = Read-SkillLock
        if (Remove-SkillLockEntry -State $state -Name $Name) {
            Write-SkillLock -State $state
            $status.Removed = $true
            return
        }
    }

    return $completed -and $status.Removed
}

function Read-SkipState {
    if (-not (Test-Path -LiteralPath $skipFile -PathType Leaf)) {
        return [pscustomobject]@{ skips = [pscustomobject]@{} }
    }

    try {
        $state = Get-Content -LiteralPath $skipFile -Raw | ConvertFrom-Json
        if (-not ((Get-JsonPropertyNames -Object $state) -contains "skips")) {
            return [pscustomobject]@{ skips = [pscustomobject]@{} }
        }
        return $state
    } catch {
        return [pscustomobject]@{ skips = [pscustomobject]@{} }
    }
}

function Write-SkipState {
    param([object]$State)

    New-Item -ItemType Directory -Path $stateRoot -Force | Out-Null
    $State | ConvertTo-Json -Depth 5 | Set-Content -LiteralPath $skipFile -Encoding UTF8
}

function Get-SkipEntry {
    param(
        [object]$State,
        [string]$Name
    )

    $names = @($State.skips.PSObject.Properties | ForEach-Object { $_.Name })
    if ($names -contains $Name) {
        return $State.skips.$Name
    }
    return $null
}

function Set-SkipEntry {
    param(
        [object]$State,
        [string]$Name,
        [string]$RemoteHash,
        [string]$SourceUrl
    )

    $entry = [pscustomobject]@{
        remoteHash = $RemoteHash
        sourceUrl = $SourceUrl
        skippedAt = (Get-Date).ToUniversalTime().ToString("o")
    }

    $names = @($State.skips.PSObject.Properties | ForEach-Object { $_.Name })
    if ($names -contains $Name) {
        $State.skips.$Name = $entry
    } else {
        $State.skips | Add-Member -NotePropertyName $Name -NotePropertyValue $entry
    }
}

function New-RemoteComparePath {
    param(
        [string]$Repo,
        [string]$RemoteDir
    )

    $tempRoot = Join-Path ([IO.Path]::GetTempPath()) ("skills-updates-compare-{0}" -f [IO.Path]::GetRandomFileName())
    New-Item -ItemType Directory -Path $tempRoot -Force | Out-Null
    $compareTempPaths.Add($tempRoot)

    $archiveArgs = @("-C", $Repo, "archive", "--format=tar")
    $archiveFile = Join-Path ([IO.Path]::GetTempPath()) ("skills-updates-compare-{0}.tar" -f [IO.Path]::GetRandomFileName())
    $compareTempFiles.Add($archiveFile)
    $archiveArgs += "--output=$archiveFile"
    $archiveArgs += "HEAD"
    if ($RemoteDir -and $RemoteDir -ne ".") {
        $archiveArgs += $RemoteDir
    }

    & git @archiveArgs
    if ($LASTEXITCODE -ne 0) {
        Write-CliError "skills-updates: could not export clean compare tree for $RemoteDir"
        return $null
    }

    tar -xf $archiveFile -C $tempRoot
    if ($LASTEXITCODE -ne 0) {
        Write-CliError "skills-updates: could not extract clean compare tree for $RemoteDir"
        return $null
    }
    Remove-Item -LiteralPath $archiveFile -Force -ErrorAction SilentlyContinue

    if ($RemoteDir -and $RemoteDir -ne ".") {
        return Join-Path $tempRoot ($RemoteDir -replace "/", [IO.Path]::DirectorySeparatorChar)
    }

    return $tempRoot
}

function Remove-SkipEntry {
    param(
        [object]$State,
        [string]$Name
    )

    $names = @($State.skips.PSObject.Properties | ForEach-Object { $_.Name })
    if ($names -contains $Name) {
        $State.skips.PSObject.Properties.Remove($Name)
        return $true
    }
    return $false
}

function Get-SavedSkillSkips {
    $state = Read-SkipState
    return [pscustomobject]@{
        State = $state
        Names = @($state.skips.PSObject.Properties | ForEach-Object { $_.Name } | Sort-Object)
    }
}

function Get-SavedSkillSkip {
    param(
        [object]$State,
        [string]$Name
    )

    return Get-SkipEntry -State $State -Name $Name
}

function Save-SkillSkip {
    param(
        [object]$State,
        [string]$Name,
        [string]$RemoteHash,
        [string]$SourceUrl
    )

    Set-SkipEntry -State $State -Name $Name -RemoteHash $RemoteHash -SourceUrl $SourceUrl
    Write-SkipState -State $State
}

function Remove-SavedSkillSkip {
    param([string]$Name)

    $state = Read-SkipState
    if (-not (Remove-SkipEntry -State $state -Name $Name)) {
        return $false
    }

    Write-SkipState -State $state
    return $true
}

if ($mode -eq "skips") {
    $savedSkips = Get-SavedSkillSkips
    $skipState = $savedSkips.State
    $skipNames = $savedSkips.Names
    if ($skipNames.Count -eq 0) {
        Write-StatusLine "OK      no saved skips"
        exit 0
    }

    foreach ($name in $skipNames) {
        $entry = Get-SavedSkillSkip -State $skipState -Name $name
        Write-StatusLine "SKIP    $name $($entry.remoteHash)"
    }
    exit 0
}

if ($mode -eq "unskip") {
    if (Remove-SavedSkillSkip -Name $target) {
        Write-StatusLine "OK      removed skip for $target"
    } else {
        Write-StatusLine "OK      no saved skip for $target"
    }
    exit 0
}

if ($mode -eq "uninstall") {
    $failedCount = 0

    foreach ($name in $targets) {
        if (-not (Uninstall-Skill -Name $name)) {
            $failedCount += 1
            continue
        }
    }

    if ($failedCount -gt 0) {
        exit 1
    }

    exit 0
}

if ($mode -eq "install" -and $targets.Count -gt 0) {
    $installSources = New-Object System.Collections.Generic.List[string]
    $installNames = New-Object System.Collections.Generic.List[string]

    foreach ($requested in $targets) {
        if (Test-InstallSourceArgument -Value $requested) {
            $installSources.Add($requested)
        } else {
            $installNames.Add($requested)
        }
    }

    if ($installSources.Count -gt 0) {
        if ($installNames.Count -gt 0) {
            Write-CliError "skills-updates: install cannot mix source URLs and lockfile skill names"
            exit 1
        }

        $failedCount = 0
        foreach ($sourceUrl in $installSources) {
            if (-not (Install-SkillSource -SourceUrl $sourceUrl)) {
                $failedCount += 1
            }
        }

        if ($failedCount -gt 0) {
            exit 1
        }

        exit 0
    }
}

if ($mode -eq "list") {
    if ($targets.Count -gt 0) {
        Write-CliError "skills-updates: list mode does not accept skill names"
        exit 1
    }

    Show-InstalledSkills
    exit 0
}

if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
    Write-CliError "skills-updates: git is required"
    exit 1
}
if (-not (Test-Path -LiteralPath $lock -PathType Leaf)) {
    Write-CliError "skills-updates: lockfile not found: $lock"
    exit 1
}
New-Item -ItemType Directory -Path $repoRoot -Force | Out-Null

try {
    $lockData = Get-Content -LiteralPath $lock -Raw | ConvertFrom-Json
    $skipState = (Get-SavedSkillSkips).State
    $zedArgs = New-Object System.Collections.Generic.List[string]
    $changedCount = 0
    $groups = @{}
    $results = @{}
    $sourcesByName = @{}
    $showProgress = $mode -in @("summary", "zed") -and $targets.Count -eq 0

    foreach ($entry in $lockData.skills.PSObject.Properties) {
        $name = $entry.Name
        if (-not (Test-TargetMatch -Name $name)) {
            continue
        }

        $value = $entry.Value
        $sourceUrl = if ((Get-JsonPropertyNames -Object $value) -contains "sourceUrl" -and $value.sourceUrl) {
            $value.sourceUrl
        } else {
            "https://github.com/$($value.source).git"
        }
        $sourcesByName[$name] = $sourceUrl
        $skillPath = $value.skillPath
        $remoteDir = (Split-Path -Parent $skillPath) -replace "\\", "/"
        if (-not $remoteDir) {
            $remoteDir = "."
        }

        $installedDir = Join-Path $skillsDir $name

        if (-not (Test-Path -LiteralPath $installedDir -PathType Container)) {
            $results[$name] = @{
                Status = "MISSING"
                Message = "MISSING $name`: installed directory not found at $installedDir"
            }
            continue
        }

        $item = [pscustomobject]@{
            Name = $name
            SourceUrl = $sourceUrl
            RemoteDir = $remoteDir
            InstalledDir = $installedDir
        }

        if (-not $groups.ContainsKey($sourceUrl)) {
            $groups[$sourceUrl] = New-Object System.Collections.Generic.List[object]
        }
        $groups[$sourceUrl].Add($item)
    }

    $sourceUrls = @($groups.Keys | Sort-Object)
    $sourceIndex = 0
    $repoJobs = New-Object System.Collections.Generic.List[object]
    $repoUpdateScript = {
        param(
            [string]$SourceUrl,
            [string]$Repo,
            [string[]]$RemoteDirs
        )

        try {
            if (Test-Path -LiteralPath (Join-Path $Repo ".git") -PathType Container) {
                git -c core.autocrlf=false -C $Repo remote set-url origin $SourceUrl 2>$null
                git -c core.autocrlf=false -C $Repo fetch --quiet --depth 1 --filter=blob:none origin HEAD 2>$null
                if ($LASTEXITCODE -ne 0) {
                    throw "git fetch failed"
                }

                git -c core.autocrlf=false -C $Repo reset --quiet --hard FETCH_HEAD 2>$null
                if ($LASTEXITCODE -ne 0) {
                    throw "git reset failed"
                }
            } else {
                if (Test-Path -LiteralPath $Repo) {
                    Remove-Item -LiteralPath $Repo -Recurse -Force -ErrorAction SilentlyContinue
                }

                git -c core.autocrlf=false clone --quiet --depth 1 --filter=blob:none --sparse $SourceUrl $Repo 2>$null
                if ($LASTEXITCODE -ne 0) {
                    throw "git clone failed"
                }
            }

            git -c core.autocrlf=false -C $Repo sparse-checkout set @RemoteDirs 2>$null
            if ($LASTEXITCODE -ne 0) {
                throw "git sparse-checkout failed"
            }

            [pscustomobject]@{
                Ok = $true
                Message = ""
            }
        } catch {
            [pscustomobject]@{
                Ok = $false
                Message = $_.Exception.Message
            }
        }
    }

    $repoRunspacePool = [runspacefactory]::CreateRunspacePool($maxRepoJobs, $maxRepoJobs)
    $repoRunspacePool.Open()
    foreach ($sourceUrl in $sourceUrls) {
        $sourceIndex += 1
        $groupItems = $groups[$sourceUrl]
        $repo = Join-Path $repoRoot (Get-CacheKey $sourceUrl)
        $sourceLabel = $sourceUrl -replace "^https://github\.com/", "" -replace "\.git$", ""
        $remoteDirs = @($groupItems | ForEach-Object { $_.RemoteDir } | Sort-Object -Unique)

        if ($showProgress) {
            Write-StatusLine "CHECK   repo $sourceIndex/$($sourceUrls.Count): $sourceLabel ($($groupItems.Count) skill(s))"
        }

        $hasRepo = Test-Path -LiteralPath (Join-Path $repo ".git") -PathType Container
        if ($showProgress) {
            if ($hasRepo) {
                Write-StatusLine "FETCH   $sourceLabel"
            } else {
                Write-StatusLine "CLONE   $sourceLabel"
            }
        }

        $job = [powershell]::Create()
        $job.RunspacePool = $repoRunspacePool
        [void]$job.AddScript($repoUpdateScript.ToString()).AddArgument($sourceUrl).AddArgument($repo).AddArgument($remoteDirs)
        $asyncResult = $job.BeginInvoke()

        $repoJobs.Add([pscustomobject]@{
            SourceUrl = $sourceUrl
            SourceLabel = $sourceLabel
            GroupItems = $groupItems
            Repo = $repo
            Job = $job
            AsyncResult = $asyncResult
        })
    }

    while ($repoJobs.Count -gt 0) {
        $waitHandles = @($repoJobs | ForEach-Object { $_.AsyncResult.AsyncWaitHandle })
        $completedIndex = [Threading.WaitHandle]::WaitAny($waitHandles)
        $repoJob = $repoJobs[$completedIndex]
        $repoJobs.RemoveAt($completedIndex)
        $sourceUrl = $repoJob.SourceUrl
        $sourceLabel = $repoJob.SourceLabel
        $groupItems = $repoJob.GroupItems
        $repo = $repoJob.Repo
        try {
            $jobResult = @($repoJob.Job.EndInvoke($repoJob.AsyncResult))[0]
        } catch {
            $jobResult = [pscustomobject]@{
                Ok = $false
                Message = $_.Exception.Message
            }
        } finally {
            $repoJob.Job.Dispose()
        }

        if (-not $jobResult -or -not $jobResult.Ok) {
            $errorDetail = if ($jobResult -and $jobResult.Message) { ": $($jobResult.Message)" } else { "" }
            foreach ($item in $groupItems) {
                $results[$item.Name] = @{
                    Status = "ERROR"
                    Message = "ERROR   $($item.Name)`: could not update local clone for $sourceUrl$errorDetail"
                }
            }
            continue
        }

        if ($showProgress) {
            Write-StatusLine "COMPARE $sourceLabel"
        }

        foreach ($item in $groupItems) {
            $remotePath = New-RemoteComparePath -Repo $repo -RemoteDir $item.RemoteDir
            if (-not $remotePath) {
                $results[$item.Name] = @{
                    Status = "ERROR"
                    Message = "ERROR   $($item.Name)`: could not export upstream compare tree"
                }
                continue
            }
            git -c core.autocrlf=false diff --quiet --ignore-cr-at-eol --no-index -- $item.InstalledDir $remotePath 2>$null
            $diffExit = $LASTEXITCODE

            if ($diffExit -eq 0) {
                $results[$item.Name] = @{
                    Status = "OK"
                    Message = "OK      $($item.Name)"
                    InstalledDir = $item.InstalledDir
                    RemotePath = $remotePath
                    SourceUrl = $item.SourceUrl
                    RemoteHash = ""
                }
            } else {
                $remoteHash = ""
                git -c core.autocrlf=false -C $repo rev-parse "HEAD:$($item.RemoteDir)" 2>$null | ForEach-Object {
                    $remoteHash = $_.Trim()
                }
                $skipEntry = Get-SavedSkillSkip -State $skipState -Name $item.Name
                if ($skipEntry -and $skipEntry.remoteHash -eq $remoteHash) {
                    $results[$item.Name] = @{
                        Status = "SKIP"
                        Message = "SKIP    $($item.Name)"
                        InstalledDir = $item.InstalledDir
                        RemotePath = $remotePath
                        SourceUrl = $item.SourceUrl
                        RemoteHash = $remoteHash
                    }
                } else {
                    $results[$item.Name] = @{
                        Status = "UPDATE"
                        Message = "UPDATE  $($item.Name)"
                        InstalledDir = $item.InstalledDir
                        RemotePath = $remotePath
                        SourceUrl = $item.SourceUrl
                        RemoteHash = $remoteHash
                    }
                }
            }
        }
    }

    if ($mode -eq "skip") {
        if (-not $results.ContainsKey($target)) {
            Write-CliError "skills-updates: skill not found in global lockfile: $target"
            exit 1
        }

        $result = $results[$target]
        if ($result.Status -eq "OK") {
            Write-StatusLine "OK      $target has no current update to skip"
            exit 0
        }

        if (-not $result.RemoteHash) {
            Write-StatusLine "ERROR   $target`: no upstream hash available to skip"
            exit 1
        }

        Save-SkillSkip -State $skipState -Name $target -RemoteHash $result.RemoteHash -SourceUrl $result.SourceUrl
        Write-StatusLine "SKIP    saved current update for $target"
        exit 0
    }

    if ($mode -in @("install", "install-all")) {
        $installNames = New-Object System.Collections.Generic.List[string]

        if ($mode -eq "install") {
            foreach ($requestedName in $targets) {
                if (-not $sourcesByName.ContainsKey($requestedName)) {
                    Write-CliError "skills-updates: skill not found in global lockfile: $requestedName"
                    exit 1
                }
            }
        }

        foreach ($entry in $lockData.skills.PSObject.Properties) {
            $name = $entry.Name
            if (-not (Test-TargetMatch -Name $name)) {
                continue
            }

            if (-not $results.ContainsKey($name)) {
                continue
            }

            $result = $results[$name]
            if ($mode -eq "install") {
                $installNames.Add($name)
            } elseif ($result.Status -in @("UPDATE", "MISSING")) {
                $installNames.Add($name)
            }
        }

        if ($installNames.Count -eq 0) {
            if ($targets.Count -gt 0) {
                Write-StatusLine "OK      requested skills are up to date"
            } else {
                Write-StatusLine "OK      all skills are up to date"
            }
            exit 0
        }

        $failedCount = 0
        foreach ($name in $installNames) {
            if (-not $sourcesByName.ContainsKey($name)) {
                Write-StatusLine "ERROR   $name`: source not found in lockfile"
                $failedCount += 1
                continue
            }

            if (-not (Install-Skill -Name $name -SourceUrl $sourcesByName[$name])) {
                $failedCount += 1
            }
        }

        if ($failedCount -gt 0) {
            exit 1
        }

        exit 0
    }

    $statusIndex = 0
    $statusEntries = @($lockData.skills.PSObject.Properties)
    if ($mode -eq "summary" -and -not $target) {
        $statusEntries = @($statusEntries | Sort-Object -Property @{ Expression = { $_.Name } })
    }

    foreach ($entry in $statusEntries) {
        $name = $entry.Name
        if (-not (Test-TargetMatch -Name $name)) {
            continue
        }

        if (-not $results.ContainsKey($name)) {
            continue
        }

        $result = $results[$name]
        if ($result.Status -eq "UPDATE") {
            $changedCount += 1
            if ($mode -eq "summary") {
                $statusIndex += 1
                Write-StatusLine ("{0,2}. {1}" -f $statusIndex, $result.Message)
            } else {
                if ($mode -eq "zed") {
                    Write-StatusLine "DIFF    $name"
                    $zedArgs.Add("--diff")
                    $zedArgs.Add($result.InstalledDir)
                    $zedArgs.Add($result.RemotePath)
                }
                if ($mode -eq "diff") {
                    Write-Output ""
                    Write-Output "===== $name ====="
                    git -c core.autocrlf=false diff --ignore-cr-at-eol --no-index --color=auto -- $result.InstalledDir $result.RemotePath
                    if ($LASTEXITCODE -notin @(0, 1)) {
                        Write-Output "ERROR   $name`: diff failed"
                    }
                }
            }
        } elseif ($result.Status -eq "OK") {
            if ($mode -eq "summary" -or $targets.Count -gt 0) {
                if ($mode -eq "summary" -and -not $target) {
                    $statusIndex += 1
                    Write-StatusLine ("{0,2}. {1}" -f $statusIndex, $result.Message)
                } else {
                    Write-StatusLine $result.Message
                }
            }
        } elseif ($result.Status -eq "SKIP") {
            if ($targets.Count -gt 0) {
                Write-StatusLine $result.Message
            }
        } else {
            Write-StatusLine $result.Message
        }
    }

    if ($mode -eq "zed") {
        if ($changedCount -eq 0) {
            if ($target) {
                Write-StatusLine "OK      $target"
            } elseif ($targets.Count -gt 0) {
                Write-StatusLine "OK      requested skills are up to date"
            } else {
                Write-StatusLine "OK      all skills are up to date"
            }
        } elseif (Get-Command zed -ErrorAction SilentlyContinue) {
            Write-ColoredLine "Opening $changedCount diff(s) in Zed." Cyan
            & zed @zedArgs
            if ($LASTEXITCODE -ne 0) {
                Write-CliError "skills-updates: zed failed to open diff viewer"
                exit 1
            }
            $preserveCompareTempPaths = $true
        } else {
            Write-CliError "skills-updates: zed command not found"
            exit 1
        }
    }
} finally {
    if ($repoRunspacePool) {
        $repoRunspacePool.Close()
        $repoRunspacePool.Dispose()
    }
    foreach ($tempFile in $compareTempFiles) {
        if (Test-Path -LiteralPath $tempFile) {
            Remove-Item -LiteralPath $tempFile -Force -ErrorAction SilentlyContinue
        }
    }
    if (-not $preserveCompareTempPaths) {
        foreach ($tempPath in $compareTempPaths) {
            if (Test-Path -LiteralPath $tempPath) {
                Remove-Item -LiteralPath $tempPath -Recurse -Force -ErrorAction SilentlyContinue
            }
        }
    }
}
