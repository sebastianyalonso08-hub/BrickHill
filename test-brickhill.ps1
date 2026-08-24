Write-Host "Testing Brick Hill launcher..." -ForegroundColor Cyan
try { Start-Process 'brickhill://test'; Write-Host "Windows accepted the brickhill:// URL." -ForegroundColor Green } catch { Write-Host "Windows did not accept brickhill://. Run install-brickhill.ps1 again." -ForegroundColor Red }
Read-Host "Press Enter to close"
