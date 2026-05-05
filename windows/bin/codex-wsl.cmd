@echo off
setlocal

:: Configuration
set "WSL_EXE=%SystemRoot%\System32\wsl.exe"
set "WSL_DISTRO=Ubuntu"
set "WSL_USER=crunch"
set "CODEX_HOME=/home/crunch/.codex"
set "CODEX_WSL_PROXY=/home/crunch/.local/bin/codex-wsl-proxy-runner.sh"
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

:: Execution
:: Invoke bash explicitly and let the WSL-side runner resolve the configured
:: default NVM node binary before launching the JS proxy.
"%WSL_EXE%" -d %WSL_DISTRO% -u %WSL_USER% --cd / --exec env "CODEX_HOME=%CODEX_HOME%" "CODEX_WSL_PROXY_DISTRO=%WSL_DISTRO%" "T3CODE_WINDOWS_CWD=%T3CODE_SHIM_CWD%" bash %CODEX_WSL_PROXY% %*

exit /b %ERRORLEVEL%
