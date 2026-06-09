# Shell Setup

Light rebuild notes for matching Windows, WSL, and Linux shell UX.

This is enough to recreate the shell feel: prompt, key model, fuzzy pickers,
smart directory jumping, and WSL Bash editing. It does not replicate private
state such as SSH profiles, Tabby vault contents, Atuin sync login, or existing
history databases.

## Goal

Use one prompt config and the same shell muscle memory everywhere:

- `Ctrl-R`: Atuin history search
- `Ctrl-T`: fzf file picker
- `Alt-C`: fzf directory picker
- `Tab`: PSReadLine menu completion in PowerShell; ble.sh-enhanced completion
  in WSL Bash
- `z` / `zi`: zoxide smart cd / interactive cd

Atuin sync and AI are intentionally not configured. Atuin can become the main
history search workflow, but Bash and PSReadLine still keep their own history
unless separately disabled.

## Prompt

Starship is the shared prompt engine.

Config location:

- Windows: `C:\Users\nguco\.config\starship.toml`
- WSL Ubuntu: `/home/crunch/.config/starship.toml`
- Linux/Arch: `~/.config/starship.toml`
- Repo template: `docs/shell/starship.toml`

The current config uses portable modules only: OS, shell, directory, git, command duration, inline time, and prompt character. It uses Powerline separators, so use a Nerd Font or another font with Powerline glyphs in the terminal. Color roles: blue/mauve for OS and shell, peach/yellow for git, teal for path, green/red only for command success or failure.

## Windows

Prerequisites on a clean Windows machine:

```powershell
winget install --id Microsoft.PowerShell --source winget
winget install --id Git.Git --source winget
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex
```

Set the repo path for the commands below:

```powershell
$repo = 'E:\dev\agents-toolkit'
```

Install tools with Scoop:

```powershell
scoop bucket add extras
scoop bucket add nerd-fonts
scoop install tabby
scoop install starship fzf zoxide atuin clink
scoop install CascadiaCode-NF-Mono
Install-PackageProvider NuGet -Scope CurrentUser -Force
Set-PSRepository PSGallery -InstallationPolicy Trusted
Install-Module PSReadLine -Scope CurrentUser -Force
Install-Module PSFzf -Scope CurrentUser -Force
```

Copy the shared prompt config:

```powershell
New-Item -ItemType Directory -Force -Path "$HOME\.config" | Out-Null
Copy-Item "$repo\docs\shell\starship.toml" "$HOME\.config\starship.toml" -Force
```

PowerShell 7 profile path:

```powershell
$PROFILE
```

Use this shell initialization block:

```powershell
$env:STARSHIP_CONFIG = Join-Path $HOME '.config\starship.toml'
if (Get-Command starship.exe -ErrorAction SilentlyContinue) {
    Invoke-Expression (&starship.exe init powershell)
}
```

Use this interactive tooling block in the same profile:

```powershell
if (Get-Module -ListAvailable -Name PSReadLine) {
    Import-Module PSReadLine
    Set-PSReadLineOption -EditMode Windows -BellStyle None
    if (-not [Console]::IsOutputRedirected) {
        Set-PSReadLineOption -PredictionSource History -PredictionViewStyle ListView
    }
    Set-PSReadLineKeyHandler -Chord Tab -Function MenuComplete
}

if ((Get-Command fzf.exe -ErrorAction SilentlyContinue) -and (Get-Module -ListAvailable -Name PSFzf)) {
    Import-Module PSFzf
    Set-PsFzfOption -PSReadlineChordProvider 'Ctrl+t' -PSReadlineChordSetLocation 'Alt+c'
}

$zoxideCommand = Get-Command zoxide.exe -ErrorAction SilentlyContinue
if ($zoxideCommand) {
    $zoxideInit = & $zoxideCommand.Source init powershell
    Invoke-Expression ($zoxideInit -join [Environment]::NewLine)
}

$atuinCommand = Get-Command atuin.exe -ErrorAction SilentlyContinue
if ($atuinCommand) {
    $atuinInit = & $atuinCommand.Source init powershell --disable-up-arrow --disable-ai
    Invoke-Expression ($atuinInit -join [Environment]::NewLine)
}
```

Command Prompt uses Clink through current-user AutoRun:

```cmd
%USERPROFILE%\scoop\apps\clink\current\clink.bat autorun install
```

Use the full `current` path for Scoop so AutoRun survives Clink updates.

## Tabby

Apply these manually in Tabby Settings. The shell feel is portable without
copying private profile data:

- Font: `CaskaydiaCove NFM`, or the installed Cascadia Code Nerd Font Mono
  family name if Tabby displays it differently
- Theme: Catppuccin Mocha
- Dark mode with tabs on the left
- Default profile: PowerShell
- Pane muscle memory:
  - `Alt-Shift-D`: split right
  - `Alt-Shift-S`: split bottom
  - `Alt-Arrow`: move between panes
  - `Alt-Shift-Arrow`: resize panes
  - `Ctrl-Shift-Z`: maximize pane
  - `Alt-W`: close pane

Do not copy Tabby's encrypted `vault` block or imported SSH profile cache to a
new machine. Recreate SSH profiles there from that machine's own SSH config.

## WSL Ubuntu

These commands assume an Ubuntu release with the relevant packages in `universe`
(verified on this workstation's current Ubuntu). On older WSL images, enable
`universe` first or use each tool's upstream install docs.

Set the repo path for the commands below:

```bash
REPO=/mnt/e/dev/agents-toolkit
```

Install managed Ubuntu packages:

```bash
sudo apt update
sudo apt install -y software-properties-common
sudo add-apt-repository -y universe
sudo apt update
sudo apt install -y git make gawk fzf zoxide atuin starship fd-find bat eza direnv
mkdir -p ~/.local/bin ~/.config
ln -sf /usr/bin/fdfind ~/.local/bin/fd
ln -sf /usr/bin/batcat ~/.local/bin/bat
```

Install ble.sh from source:

```bash
mkdir -p ~/.local/src
git clone --recursive --depth 1 --shallow-submodules https://github.com/akinomyoga/ble.sh.git ~/.local/src/ble.sh
make -C ~/.local/src/ble.sh install PREFIX=~/.local
```

Load fzf readline bindings before ble.sh, then add ble.sh before the later Starship/zoxide/Atuin block:

```bash
# fzf's Bash bindings use readline, so install them before ble.sh takes over line editing.
if command -v fzf >/dev/null 2>&1 && [ -t 0 ] && [ -t 1 ]; then
    [ -r /usr/share/doc/fzf/examples/completion.bash ] && . /usr/share/doc/fzf/examples/completion.bash
    [ -r /usr/share/doc/fzf/examples/key-bindings.bash ] && . /usr/share/doc/fzf/examples/key-bindings.bash
    bind -r "\C-r" 2>/dev/null || true
    declare -F fzf-file-widget >/dev/null 2>&1 && bind -m emacs-standard -x "\"\C-t\": fzf-file-widget"
fi

# ble.sh needs Bash line editing; skip bash -c/script-style interactive shells.
if [[ $- == *i* && -z ${BASH_EXECUTION_STRING+x} && -t 0 && -t 1 && -r "$HOME/.local/share/blesh/ble.sh" ]]; then
    source "$HOME/.local/share/blesh/ble.sh"
fi
```

Configure ble.sh to delegate prompt cursor shape to the terminal profile:

```bash
cat >> ~/.blerc <<'EOF'
# Ask ble.sh for terminal-default cursor shape at the Bash prompt.
# Cursor 0 is DECSCUSR default, not a hard-coded bar/block shape.
ble-bind -m emacs --cursor 0
EOF
```

This asks ble.sh to use the terminal default instead of choosing a hard-coded
bar/block shape. In Tabby, that keeps the normal Bash prompt aligned with the
Appearance cursor setting; full-screen TUI programs can still set their own
cursor shape while active.

Add this after existing PATH/toolchain setup:

```bash
if [ -n "${BASH_VERSION:-}" ] && [ -n "${PS1:-}" ]; then
    if command -v starship >/dev/null 2>&1; then
        eval "$(starship init bash)"
    fi

    if command -v fd >/dev/null 2>&1; then
        export FZF_DEFAULT_COMMAND='fd --type f --hidden --follow --exclude .git --exclude node_modules --exclude .cache --exclude .cargo/registry --exclude .local/share/pnpm/store'
        export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"
        export FZF_ALT_C_COMMAND='fd --type d --hidden --follow --exclude .git --exclude node_modules --exclude .cache --exclude .cargo/registry --exclude .local/share/pnpm/store'
    fi

    if command -v zoxide >/dev/null 2>&1; then
        eval "$(zoxide init bash)"
    fi

    if command -v atuin >/dev/null 2>&1 && [ -t 0 ] && [ -t 1 ]; then
        eval "$(atuin init bash --disable-up-arrow --disable-ai)"
    fi

    if command -v direnv >/dev/null 2>&1; then
        eval "$(direnv hook bash)"
    fi
fi
```

Copy the Starship config:

```bash
cp "$REPO/docs/shell/starship.toml" ~/.config/starship.toml
```

## Arch Linux

Install from official repos:

```bash
sudo pacman -S --needed starship fzf zoxide atuin fd bat eza direnv ripgrep tmux
```

For Bash, use the same `~/.bashrc` block as WSL, but Arch uses normal binary names, so no `fd` or `bat` symlinks are needed.

Install ble.sh from AUR (`blesh-git` or `blesh`) or use the source install
commands from the WSL section.

For Zsh:

```zsh
eval "$(starship init zsh)"
eval "$(zoxide init zsh)"
eval "$(atuin init zsh --disable-up-arrow)"
eval "$(direnv hook zsh)"
```

Add fzf keybindings from the Arch fzf package if not already enabled by `/etc/profile.d/fzf.*`.

## Verify

Windows PowerShell:

```powershell
starship --version
. $PROFILE
Get-PSReadLineKeyHandler -Bound | Where-Object Key -in 'Ctrl+r','Ctrl+t','Alt+c','Tab','UpArrow'
```

WSL/Linux:

```bash
starship explain >/dev/null
bash -n ~/.bashrc
bash ~/.local/share/blesh/ble.sh --version
fzf --version
zoxide --version
atuin --version
fd --version
bat --version
eza --version
```

Expected key model:

- `Ctrl-R` opens Atuin history.
- `Ctrl-T` opens fzf file picker.
- `Alt-C` opens fzf directory picker.
- `Tab` uses PSReadLine MenuComplete in PowerShell and ble.sh-enhanced
  completion in WSL Bash.
- `UpArrow` remains normal shell history.
