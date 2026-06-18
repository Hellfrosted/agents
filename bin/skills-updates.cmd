@echo off
setlocal

if not defined SK_UP_CMD_NAME set "SK_UP_CMD_NAME=%~n0"
for /f "tokens=2 delims=:" %%A in ('chcp') do set "SK_UP_PREVIOUS_CODEPAGE=%%A"
for /f "tokens=* delims= " %%A in ("%SK_UP_PREVIOUS_CODEPAGE%") do set "SK_UP_PREVIOUS_CODEPAGE=%%A"
chcp 65001 >nul
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%~dp0skills-updates.ps1" --cmd-name "%SK_UP_CMD_NAME%" %*
set "SK_UP_EXIT_CODE=%ERRORLEVEL%"
if defined SK_UP_PREVIOUS_CODEPAGE chcp %SK_UP_PREVIOUS_CODEPAGE% >nul
exit /b %SK_UP_EXIT_CODE%
