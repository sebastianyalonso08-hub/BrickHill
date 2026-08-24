$ErrorActionPreference = "Stop"
$base = "https://brickhill.onrender.com"
$installDir = Join-Path $env:LOCALAPPDATA "BrickHill"
$zip = Join-Path $env:TEMP "BrickHillClient.zip"

Write-Host "Brick Hill Client Installer"
Write-Host "Downloading client files..."
Invoke-WebRequest -Uri "$base/BrickHillClient.zip" -OutFile $zip -UseBasicParsing

if (!(Test-Path $installDir)) { New-Item -ItemType Directory -Path $installDir -Force | Out-Null }
Expand-Archive -Path $zip -DestinationPath $installDir -Force
Remove-Item $zip -Force -ErrorAction SilentlyContinue

$launcher = Join-Path $installDir "BrickHillLauncher.exe"
$game = Join-Path $installDir "Brick_Hill_Multiplayer.exe"
$dll = Join-Path $installDir "BRICK.dll"
$bridge = Join-Path $installDir "BrickHillNetworkBridge.exe"
$shim = Join-Path $installDir "wsock32.dll"

foreach ($file in @($launcher,$game,$dll,$bridge,$shim)) {
  if (!(Test-Path $file)) { throw "Installation failed: missing $file" }
}

$command = '"' + $launcher + '" "%1"'
$protocol = "HKCU:\Software\Classes\brickhill"
New-Item -Path "$protocol\shell\open\command" -Force | Out-Null
Set-ItemProperty -Path $protocol -Name "(Default)" -Value "URL:Brick Hill Protocol"
Set-ItemProperty -Path $protocol -Name "URL Protocol" -Value ""
Set-ItemProperty -Path "$protocol\shell\open\command" -Name "(Default)" -Value $command

Write-Host ""
Write-Host "Brick Hill Client installed successfully."
Write-Host "Installed to: $installDir"
Write-Host "You can now return to brickhill.onrender.com and press Play."
