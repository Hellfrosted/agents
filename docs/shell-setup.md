# Shell Setup

Light rebuild notes for matching Windows, WSL, and Linux shell UX.

## Goal

Use one prompt config and the same shell muscle memory everywhere:

- `Ctrl-R`: Atuin history search
- `Ctrl-T`: fzf file picker
- `Alt-C`: fzf directory picker
- `Tab`: normal shell completion
- `z` / `zi`: zoxide smart cd / interactive cd

Atuin sync and AI are intentionally not configured.

## Prompt

Starship is the shared prompt engine.

Config location:

- Windows: `C:\Users\nguco\.config\starship.toml`
- WSL Ubuntu: `/home/crunch/.config/starship.toml`
- Linux/Arch: `~/.config/starship.toml`

The current config uses portable modules only: OS, shell, directory, git, Node, Python, Rust, Go, command duration, time, and prompt character.

## Windows

Install tools with Scoop:

```powershell
scoop install starship fzf zoxide atuin clink
Install-Module PSFzf -Scope CurrentUser
```

PowerShell 7 profile:

```powershell
$env:STARSHIP_CONFIG = Join-Path $HOME '.config\starship.toml'
if (Get-Command starship.exe -ErrorAction SilentlyContinue) {
    Invoke-Expression (&starship.exe init powershell)
}
```

Interactive tooling in the same profile:

```powershell
Import-Module PSReadLine
Set-PSReadLineOption -EditMode Windows -BellStyle None
Set-PSReadLineKeyHandler -Chord Tab -Function MenuComplete

Import-Module PSFzf
Set-PsFzfOption -PSReadlineChordProvider 'Ctrl+t' -PSReadlineChordSetLocation 'Alt+c'

Invoke-Expression ((zoxide init powershell) -join [Environment]::NewLine)
Invoke-Expression ((atuin init powershell --disable-up-arrow --disable-ai) -join [Environment]::NewLine)
```

Command Prompt uses Clink through current-user AutoRun:

```cmd
%USERPROFILE%\scoop\apps\clink\current\clink.bat autorun install
```

Use the full `current` path for Scoop so AutoRun survives Clink updates.

## WSL Ubuntu

Install managed Ubuntu packages:

```bash
sudo apt update
sudo apt install -y fzf zoxide atuin starship fd-find bat eza direnv
mkdir -p ~/.local/bin ~/.config
ln -sf /usr/bin/fdfind ~/.local/bin/fd
ln -sf /usr/bin/batcat ~/.local/bin/bat
```

Add to `~/.bashrc` after existing PATH/toolchain setup:

```bash
if [ -n "${BASH_VERSION:-}" ] && [ -n "${PS1:-}" ]; then
    eval "$(starship init bash)"

    export FZF_DEFAULT_COMMAND='fd --type f --hidden --follow --exclude .git --exclude node_modules --exclude .cache --exclude .cargo/registry --exclude .local/share/pnpm/store'
    export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"
    export FZF_ALT_C_COMMAND='fd --type d --hidden --follow --exclude .git --exclude node_modules --exclude .cache --exclude .cargo/registry --exclude .local/share/pnpm/store'

    if [ -t 0 ] && [ -t 1 ]; then
        . /usr/share/doc/fzf/examples/completion.bash
        . /usr/share/doc/fzf/examples/key-bindings.bash
        bind -r "\C-r" 2>/dev/null || true
    fi

    eval "$(zoxide init bash)"
    [ -t 0 ] && [ -t 1 ] && eval "$(atuin init bash --disable-up-arrow)"
    eval "$(direnv hook bash)"
fi
```

Copy the Starship config:

```bash
cp /mnt/c/Users/nguco/.config/starship.toml ~/.config/starship.toml
```

## Arch Linux

Install from official repos:

```bash
sudo pacman -S --needed starship fzf zoxide atuin fd bat eza direnv ripgrep tmux
```

For Bash, use the same `~/.bashrc` block as WSL, but Arch uses normal binary names, so no `fd` or `bat` symlinks are needed.

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
- `UpArrow` remains normal shell history.
