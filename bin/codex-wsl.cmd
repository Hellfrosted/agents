@echo off
setlocal EnableExtensions EnableDelayedExpansion

set "WSL_EXE=%CODEX_WSL_EXE%"
if not defined WSL_EXE set "WSL_EXE=%SystemRoot%\System32\wsl.exe"

if not defined CODEX_WSL_DISTRO set "CODEX_WSL_DISTRO=Ubuntu"

set "WSL_DISTRO_ARG="
if defined CODEX_WSL_DISTRO set "WSL_DISTRO_ARG=-d %CODEX_WSL_DISTRO%"

set "WSL_USER_ARG="
if defined CODEX_WSL_USER set "WSL_USER_ARG=-u %CODEX_WSL_USER%"

if defined CODEX_WSL_SHIM_DEBUG (
    echo codex-wsl.cmd: WSL_EXE="%WSL_EXE%"
    echo codex-wsl.cmd: WSL_DISTRO_ARG="%WSL_DISTRO_ARG%"
    echo codex-wsl.cmd: WSL_USER_ARG="%WSL_USER_ARG%"
)

if not defined CODEX_WSL_PROXY_IDLE_TIMEOUT_MS set "CODEX_WSL_PROXY_IDLE_TIMEOUT_MS=1800000"
if not defined CODEX_WSL_PROXY_DISTRO set "CODEX_WSL_PROXY_DISTRO=%CODEX_WSL_DISTRO%"

set "CODEX_WSL_ENV_NAMES=CODEX_WSL_PROXY_JS CODEX_WSL_PROXY_TARGET CODEX_WSL_PROXY_DEBUG_LOG CODEX_WSL_PROXY_SKILLS_TIMEOUT_MS CODEX_SKILLS_DIRS CODEX_SKILL_ROOTS LAZYCODEX_AUTO_UPDATE_DISABLED OMO_CODEX_AUTO_UPDATE_DISABLED LAZYCODEX_CONFIG_MIGRATION_DISABLED OMO_CODEX_CONFIG_MIGRATION_DISABLED OMO_CODEX_DISABLE_POSTHOG OMO_CODEX_SEND_ANONYMOUS_TELEMETRY OMO_DISABLE_POSTHOG OMO_SEND_ANONYMOUS_TELEMETRY OPENAI_API_KEY OPENAI_BASE_URL HTTP_PROXY HTTPS_PROXY NO_PROXY SSL_CERT_FILE SSL_CERT_DIR NODE_EXTRA_CA_CERTS"
set "CODEX_WSLENV_APPEND="
for %%V in (%CODEX_WSL_ENV_NAMES%) do if defined %%V call set "CODEX_WSLENV_APPEND=%%CODEX_WSLENV_APPEND%%:%%V"
if defined CODEX_WSLENV_APPEND (
    if defined WSLENV (
        set "WSLENV=%WSLENV%%CODEX_WSLENV_APPEND%"
    ) else (
        set "WSLENV=%CODEX_WSLENV_APPEND:~1%"
    )
)

set "T3CODE_SHIM_CWD=%__CD__%"
if not defined T3CODE_SHIM_CWD set "T3CODE_SHIM_CWD=%CD%"
set "T3CODE_SHIM_DRIVE=%T3CODE_SHIM_CWD:~0,2%"
set "T3CODE_SHIM_REMOTE="
for /f "tokens=1,2,*" %%A in ('net use %T3CODE_SHIM_DRIVE% 2^>nul ^| findstr /C:"Remote name"') do set "T3CODE_SHIM_REMOTE=%%C"
if defined T3CODE_SHIM_REMOTE set "T3CODE_SHIM_CWD=%T3CODE_SHIM_REMOTE%%T3CODE_SHIM_CWD:~2%"
if not "%T3CODE_SHIM_CWD:~3%"=="" if "%T3CODE_SHIM_CWD:~-1%"=="\" set "T3CODE_SHIM_CWD=%T3CODE_SHIM_CWD:~0,-1%"

:: Validation
if not exist "%WSL_EXE%" (
    echo {"error": "Failed to find wsl.exe at %WSL_EXE%"}
    exit /b 1
)

if not defined CODEX_WSL_PROXY set "CODEX_WSL_PROXY="

:: Invoke bash explicitly and let the WSL-side runner resolve the configured
:: node binary before launching the proxy. T3code keeps Codex app-server
:: sessions open, so the WSL proxy also enforces an idle reaper.
if defined CODEX_WSL_SHIM_DEBUG echo codex-wsl.cmd: launching "%WSL_EXE%" %WSL_DISTRO_ARG% %WSL_USER_ARG% --cd / --exec env CODEX_WSL_PROXY="%CODEX_WSL_PROXY%" bash -lc "exec \"${CODEX_WSL_PROXY:-$HOME/.local/bin/codex-wsl-proxy-runner.sh}\" \"$@\"" codex-wsl-proxy %*
"%WSL_EXE%" %WSL_DISTRO_ARG% %WSL_USER_ARG% --cd / --exec env "CODEX_HOME=%CODEX_HOME%" "CODEX_WSL_PROXY=%CODEX_WSL_PROXY%" "CODEX_WSL_PROXY_DISTRO=%CODEX_WSL_PROXY_DISTRO%" "CODEX_WSL_PROXY_IDLE_TIMEOUT_MS=%CODEX_WSL_PROXY_IDLE_TIMEOUT_MS%" "T3CODE_WINDOWS_CWD=%T3CODE_SHIM_CWD%" /bin/bash -lc "exec \"${CODEX_WSL_PROXY:-$HOME/.local/bin/codex-wsl-proxy-runner.sh}\" \"$@\"" codex-wsl-proxy %*

exit /b %ERRORLEVEL%
