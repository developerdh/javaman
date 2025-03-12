# Windows platform installation script

# Get the latest release information
$repo = "developerdh/javaman"
$latest = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest"
$version = $latest.tag_name

# Create installation directory
$installDir = "$env:USERPROFILE\.javaman"
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

# Download the corresponding binary file
$assetName = "javaman_windows_amd64.zip"
$asset = $latest.assets | Where-Object { $_.name -eq $assetName }
if ($asset) {
    $downloadUrl = $asset.browser_download_url
    $zipPath = Join-Path $installDir "javaman.zip"
    
    Write-Host "Downloading javaman $version..."
    Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath
    
    # Extract files
    Expand-Archive -Path $zipPath -DestinationPath $installDir -Force
    Remove-Item $zipPath
    
    # Add to PATH environment variable
    $currentPath = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User)
    if ($currentPath -notlike "*$installDir*") {
        [Environment]::SetEnvironmentVariable(
            "Path",
            "$currentPath;$installDir",
            [EnvironmentVariableTarget]::User
        )
    }
    
    Write-Host "Javaman has been successfully installed!"
    Write-Host "Please reopen the terminal to apply the environment variables."
} else {
    Write-Host "Error: The corresponding download file was not found."
    exit 1
}

Write-Host "Press any key to continue..."
Read-Host | Out-Null