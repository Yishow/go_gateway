@echo off
setlocal
REM Go Gateway launcher - ASCII only to avoid codepage issues

REM Ensure Go and nodejs on PATH (resolve tools directly to avoid
REM substring false-positives such as "Go\bin" matching "...\Cargo\bin")
where go >nul 2>&1 || set "PATH=%PATH%;C:\Program Files\Go\bin"
echo %PATH% | "%SystemRoot%\System32\findstr.exe" /i /l /c:"%USERPROFILE%\go\bin" >nul || set "PATH=%PATH%;%USERPROFILE%\go\bin"
where node >nul 2>&1 || set "PATH=%PATH%;C:\Program Files\nodejs"
where pwsh >nul 2>&1 || set "PATH=%PATH%;C:\Program Files\PowerShell\7"

REM Enable pnpm via corepack if missing
where pnpm >nul 2>&1
if %errorlevel% neq 0 (
    where corepack >nul 2>&1 && call corepack enable >nul 2>&1
)

cd /d "%~dp0"
where pwsh >nul 2>&1
if %errorlevel%==0 (
    pwsh -NoProfile -ExecutionPolicy Bypass -File "%~dp0start.ps1" %*
) else (
    powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0start.ps1" %*
)
set "EC=%errorlevel%"
if "%~1"=="" if not "%EC%"=="0" (
    echo.
    echo [start.cmd] ExitCode=%EC% - press any key to close...
    pause >nul
)
exit /b %EC%
