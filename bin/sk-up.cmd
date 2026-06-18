@echo off
setlocal

set "SK_UP_CMD_NAME=%~n0"
"%~dp0skills-updates.cmd" %*
