@echo off
setlocal EnableExtensions

set "WSL_EXE=%CODEX_WSL_EXE%"
if not defined WSL_EXE set "WSL_EXE=%SystemRoot%\System32\wsl.exe"

set "WSL_DISTRO_ARG="
if defined CODEX_WSL_DISTRO set "WSL_DISTRO_ARG=-d %CODEX_WSL_DISTRO%"

set "WSL_USER_ARG="
if defined CODEX_WSL_USER set "WSL_USER_ARG=-u %CODEX_WSL_USER%"

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

if not defined CODEX_WSL_PROXY (
    for /f "usebackq delims=" %%A in (`"%WSL_EXE%" %WSL_DISTRO_ARG% %WSL_USER_ARG% --cd / --exec printenv HOME`) do set "CODEX_WSL_HOME=%%A"
    if not defined CODEX_WSL_HOME (
        echo {"error": "Failed to resolve HOME inside WSL"}
        exit /b 1
    )
    set "CODEX_WSL_PROXY=%CODEX_WSL_HOME%/.local/bin/codex-wsl-proxy-runner.sh"
)

:: Invoke bash explicitly and let the WSL-side runner resolve the configured
:: node binary before launching the JS proxy.
"%WSL_EXE%" %WSL_DISTRO_ARG% %WSL_USER_ARG% --cd / --exec env "CODEX_HOME=%CODEX_HOME%" "CODEX_WSL_PROXY=%CODEX_WSL_PROXY%" "CODEX_WSL_PROXY_DISTRO=%CODEX_WSL_DISTRO%" "T3CODE_WINDOWS_CWD=%T3CODE_SHIM_CWD%" bash "%CODEX_WSL_PROXY%" %*

exit /b %ERRORLEVEL%
