@echo off
setlocal

set "BINARY_NAME=devtool.exe"
set "INSTALL_DIR=%USERPROFILE%\bin"

echo Building %BINARY_NAME%...
go build -o %BINARY_NAME% main.go

if %ERRORLEVEL% NEQ 0 (
    echo Build failed! Please check your Go code.
    exit /b %ERRORLEVEL%
)

echo Build successful.

if not exist "%INSTALL_DIR%" (
    echo Creating directory %INSTALL_DIR%...
    mkdir "%INSTALL_DIR%"
)

echo Installing %BINARY_NAME% to %INSTALL_DIR%...
move /Y "%BINARY_NAME%" "%INSTALL_DIR%\%BINARY_NAME%"

if %ERRORLEVEL% EQU 0 (
    echo %BINARY_NAME% installed successfully to %INSTALL_DIR%
) else (
    echo Installation failed!
    exit /b 1
)

echo Checking if %INSTALL_DIR% is in PATH...
echo %PATH% | find /i "%INSTALL_DIR%" > nul
if %ERRORLEVEL% NEQ 0 (
    echo Adding %INSTALL_DIR% to User PATH...
    setx PATH "%PATH%;%INSTALL_DIR%"
    echo PATH updated. Please restart your terminal/command prompt to use '%BINARY_NAME%' globally.
) else (
    echo %INSTALL_DIR% is already in PATH.
)

endlocal
