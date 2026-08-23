$ErrorActionPreference = "Stop"
$dir = Split-Path -Parent $MyInvocation.MyCommand.Path
$exe = Join-Path $dir "BrickHillLauncher.exe"
if (!(Test-Path $exe)) { throw "BrickHillLauncher.exe not found in $dir" }
$command = '"' + $exe + '" "%1"'
New-Item -Path "HKCU:\Software\Classes\brickhill" -Force | Out-Null
Set-ItemProperty -Path "HKCU:\Software\Classes\brickhill" -Name "(Default)" -Value "URL:Brick Hill Protocol"
Set-ItemProperty -Path "HKCU:\Software\Classes\brickhill" -Name "URL Protocol" -Value ""
New-Item -Path "HKCU:\Software\Classes\brickhill\shell\open\command" -Force | Out-Null
Set-ItemProperty -Path "HKCU:\Software\Classes\brickhill\shell\open\command" -Name "(Default)" -Value $command
Write-Host "Brick Hill protocol installed. The website Play button can now launch the unchanged client."
if (!(Test-Path (Join-Path $dir "wsock32.dll"))) { throw "Network shim wsock32.dll is missing." }
