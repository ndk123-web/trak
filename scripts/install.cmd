@echo off
setlocal EnableDelayedExpansion

echo.
echo ========================================================
echo                 TRAK Windows Installer (CMD)
echo ========================================================
echo.

set "VERSION=v1.1.0"
set "INSTALL_DIR=%USERPROFILE%\trak\bin"
set "EXE_PATH=%INSTALL_DIR%\trak.exe"

:: 1. Detect architecture
set "ARCH=amd64"
if "%PROCESSOR_ARCHITECTURE%"=="ARM64" set "ARCH=arm64"
if "%PROCESSOR_ARCHITEW6432%"=="ARM64" set "ARCH=arm64"

echo [1/4] Detected Windows (%ARCH%)

:: 2. Create install directory
if not exist "%INSTALL_DIR%" (
    mkdir "%INSTALL_DIR%"
)
echo [2/4] Install directory: %INSTALL_DIR%

:: 3. Download binary
set "DOWNLOAD_URL=https://github.com/ndk123-web/trak/releases/download/%VERSION%/trak-windows-%ARCH%.exe"
echo.
echo [3/4] Downloading TRAK %VERSION%...

curl -fsSL "%DOWNLOAD_URL%" -o "%EXE_PATH%"
if errorlevel 1 (
    echo.
    echo [ERROR] Failed to download TRAK executable. Please check your internet connection.
    exit /b 1
)

echo       Downloaded successfully.

:: 4. Add to User PATH
echo.
echo [4/4] Configuring PATH...

for /f "tokens=2*" %%A in ('reg query "HKCU\Environment" /v Path 2^>nul') do set "USER_PATH=%%B"

echo %USER_PATH% | find /i "%INSTALL_DIR%" >nul
if errorlevel 1 (
    if defined USER_PATH (
        setx PATH "%USER_PATH%;%INSTALL_DIR%" >nul
    ) else (
        setx PATH "%INSTALL_DIR%" >nul
    )
    echo       Added to USER PATH.
) else (
    echo       Already present in USER PATH.
)

:: Update current cmd session PATH
set "PATH=%INSTALL_DIR%;%PATH%"

echo.
echo ========================================================
echo             TRAK installed successfully! 🚀
echo ========================================================
echo.
echo Installed at:
echo   %EXE_PATH%
echo.
echo [TIP FOR DEVELOPERS]
echo If 'trak' is not recognized in other open windows, restart CMD
echo or open a new terminal tab to reload the updated PATH.
echo.
echo Quick Test Commands:
echo   trak list
echo   trak init lang/go
echo.

endlocal
