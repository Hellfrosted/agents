@echo off
setlocal

powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%~dp0skills-updates.ps1" --cmd-name "%~n0" %*
exit /b %ERRORLEVEL%
