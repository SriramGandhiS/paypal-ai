Write-Host "================================================================="
Write-Host "PayBridge AI -- Golden Demo B Automated Scenario Runner"
Write-Host "Executing Go Runtime Lifecycle: cmd/demob/main.go"
Write-Host "================================================================="

$prevDir = Get-Location
try {
    Set-Location "$PSScriptRoot\..\services\api"
    go run cmd/demob/main.go
} finally {
    Set-Location $prevDir
}
