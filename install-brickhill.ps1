$ErrorActionPreference = "Stop"
$base = "https://brickhill.onrender.com"
$installDir = Join-Path $env:LOCALAPPDATA "BrickHill"
$zip = Join-Path $env:TEMP "BrickHillClient.zip"

Write-Host ""
Write-Host "=== Brick Hill Client Installer ===" -ForegroundColor Cyan
Write-Host "Installing to: $installDir"
Write-Host ""

New-Item -ItemType Directory -Path $installDir -Force | Out-Null
Write-Host "Downloading client..."
Invoke-WebRequest -Uri "$base/api/client/package" -OutFile $zip -UseBasicParsing
Expand-Archive -Path $zip -DestinationPath $installDir -Force
Remove-Item $zip -Force -ErrorAction SilentlyContinue

$launcher = Join-Path $installDir "BrickHillLauncher.exe"
$game = Join-Path $installDir "Brick_Hill_Multiplayer.exe"
$dll = Join-Path $installDir "BRICK.dll"
$bridge = Join-Path $installDir "BrickHillNetworkBridge.exe"
$shim = Join-Path $installDir "wsock32.dll"
foreach ($file in @($launcher,$game,$dll,$bridge,$shim)) {
  if (!(Test-Path -LiteralPath $file)) { throw "Installation failed: missing $file" }
}

# Register the custom brickhill:// URL protocol for the current Windows user.
$protocol = "HKCU:\Software\Classes\brickhill"
$commandKey = "$protocol\shell\open\command"
New-Item -Path $commandKey -Force | Out-Null
New-ItemProperty -Path $protocol -Name "URL Protocol" -PropertyType String -Value "" -Force | Out-Null
Set-ItemProperty -Path $protocol -Name "(Default)" -Value "URL:Brick Hill Protocol" -Force
$command = '"' + $launcher + '" "%1"'
Set-ItemProperty -Path $commandKey -Name "(Default)" -Value $command -Force

# Clear the per-user URL association cache when possible.
try { Start-Process -FilePath "$env:SystemRoot\System32\rundll32.exe" -ArgumentList 'url.dll,FileProtocolHandler','brickhill://test' -Wait -WindowStyle Hidden } catch {}

Write-Host ""
Write-Host "Brick Hill Client installed successfully." -ForegroundColor Green
Write-Host "Installed to: $installDir"
Write-Host "Protocol registered: brickhill://"
Write-Host ""
Write-Host "Testing the launcher..."
try { Start-Process 'brickhill://test' } catch { Write-Warning "Windows could not open the brickhill:// protocol automatically." }
Write-Host ""
Write-Host "If a Brick Hill Launcher window appeared, the installation is working."
Write-Host "Return to brickhill.onrender.com and click Play."
Write-Host ""
Read-Host "Press Enter to close"
