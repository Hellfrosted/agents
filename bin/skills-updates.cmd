@echo off
setlocal EnableExtensions

set "SK_UP_EXE=%~dp0sk-up.exe"
if not exist "%SK_UP_EXE%" (
    >&2 echo ERROR   missing promoted Go updater: %SK_UP_EXE%
    exit /b 2
)

for /f "tokens=2 delims=:" %%A in ('chcp') do set "SK_UP_PREVIOUS_CODEPAGE=%%A"
for /f "tokens=* delims= " %%A in ("%SK_UP_PREVIOUS_CODEPAGE%") do set "SK_UP_PREVIOUS_CODEPAGE=%%A"
chcp 65001 >nul
set "SK_UP_PREVIOUS_ENTRYPOINT=%SK_UP_ENTRYPOINT%"
set "SK_UP_ENTRYPOINT=skills-updates"
"%SK_UP_EXE%" %*
set "SK_UP_EXIT_CODE=%ERRORLEVEL%"
if defined SK_UP_PREVIOUS_ENTRYPOINT (
    set "SK_UP_ENTRYPOINT=%SK_UP_PREVIOUS_ENTRYPOINT%"
) else (
    set "SK_UP_ENTRYPOINT="
)
if defined SK_UP_PREVIOUS_CODEPAGE chcp %SK_UP_PREVIOUS_CODEPAGE% >nul
exit /b %SK_UP_EXIT_CODE%
