@echo off
setlocal enabledelayedexpansion

:: Set installation directory
set "INSTALL_DIR=%USERPROFILE%\.javaman"

:: Create installation directory
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

:: Download temporary files
set "TEMP_JSON=%TEMP%\javaman_release.json"
set "TEMP_ZIP=%TEMP%\javaman.zip"

echo Fetching latest version information...
powershell -Command "& {Invoke-RestMethod -Uri 'https://api.github.com/repos/developerdh/javaman/releases/latest' | ConvertTo-Json | Set-Content -Path '%TEMP_JSON%'}"

:: Parse version information and download URL
for /f "tokens=* usebackq" %%a in (`powershell -Command "& {$json = Get-Content '%TEMP_JSON%' | ConvertFrom-Json; $asset = $json.assets | Where-Object { $_.name -eq 'javaman_windows_amd64.zip' }; $asset.browser_download_url}"`) do (
    set "DOWNLOAD_URL=%%a"
)

if not defined DOWNLOAD_URL (
    echo Error: Download file not found
    goto :cleanup
)

:: Download file
echo Downloading javaman...
powershell -Command "& {Invoke-WebRequest -Uri '!DOWNLOAD_URL!' -OutFile '%TEMP_ZIP%'}"

:: Extract file
echo Extracting files...
powershell -Command "& {Expand-Archive -Path '%TEMP_ZIP%' -DestinationPath '%INSTALL_DIR%' -Force}"

:: Set user environment variable
echo Configuring environment variables...
for /f "tokens=* usebackq" %%a in (`powershell -Command "& {[Environment]::GetEnvironmentVariable('Path', 'User')}"`) do (
    set "USER_PATH=%%a"
)

echo !USER_PATH! | findstr /C:"%INSTALL_DIR%" >nul
if errorlevel 1 (
    powershell -Command "& {[Environment]::SetEnvironmentVariable('Path', [Environment]::GetEnvironmentVariable('Path', 'User') + ';%INSTALL_DIR%', 'User')}"
)

echo javaman installation completed!
echo Please reopen the command prompt to apply the environment variables.

:cleanup
:: Clean up temporary files
if exist "%TEMP_JSON%" del "%TEMP_JSON%"
if exist "%TEMP_ZIP%" del "%TEMP_ZIP%"

echo.
echo Press any key to exit...
pause >nul
endlocal