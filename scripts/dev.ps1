<#
.SYNOPSIS
    PayBridge AI — Windows PowerShell Developer CLI
.DESCRIPTION
    Native PowerShell alternative to Makefile for local Windows development.
.EXAMPLE
    .\scripts\dev.ps1 setup
    .\scripts\dev.ps1 up
    .\scripts\dev.ps1 status
    .\scripts\dev.ps1 test
#>

param (
    [Parameter(Position=0)]
    [ValidateSet("setup", "up", "down", "status", "test", "lint", "codegen", "migrate-up", "migrate-status", "help")]
    [string]$Command = "help"
)

$ErrorActionPreference = "Continue"
$ComposeFile = "infra/docker/docker-compose.yml"
$DbUrl = if ($env:DATABASE_URL) { $env:DATABASE_URL } else { "postgres://paybridge_admin:paybridge_dev_secret@localhost:5432/paybridge_dev?sslmode=disable" }

function Show-Help {
    Write-Host "PayBridge AI Windows Developer Commands:" -ForegroundColor Cyan
    Write-Host "  .\scripts\dev.ps1 setup          - Verify environment prerequisites"
    Write-Host "  .\scripts\dev.ps1 up             - Start local Docker infrastructure"
    Write-Host "  .\scripts\dev.ps1 down           - Stop local Docker infrastructure"
    Write-Host "  .\scripts\dev.ps1 status         - Check status of local containers"
    Write-Host "  .\scripts\dev.ps1 test           - Run backend tests"
    Write-Host "  .\scripts\dev.ps1 lint           - Run Go vet and static checks"
    Write-Host "  .\scripts\dev.ps1 codegen        - Regenerate Go and TS types from OpenAPI"
    Write-Host "  .\scripts\dev.ps1 migrate-up     - Run database migrations"
    Write-Host "  .\scripts\dev.ps1 migrate-status - Inspect database migration status"
}

switch ($Command) {
    "setup" {
        Write-Host "==> Checking prerequisites..." -ForegroundColor Yellow
        if (Get-Command go -ErrorAction SilentlyContinue) {
            Write-Host "[OK] Go: $(go version)" -ForegroundColor Green
        } else {
            Write-Host "[FAIL] Go is not installed or not in PATH." -ForegroundColor Red
        }
        if (Get-Command node -ErrorAction SilentlyContinue) {
            Write-Host "[OK] Node.js: $(node -v)" -ForegroundColor Green
        } else {
            Write-Host "[FAIL] Node.js is not installed or not in PATH." -ForegroundColor Red
        }
        if (Get-Command docker -ErrorAction SilentlyContinue) {
            Write-Host "[OK] Docker: $(docker --version)" -ForegroundColor Green
        } else {
            Write-Host "[WARN] Docker CLI is not installed or not in PATH." -ForegroundColor Yellow
        }
        if (-not (Test-Path ".env")) {
            if (Test-Path ".env.example") {
                Copy-Item ".env.example" ".env"
                Write-Host "[OK] Created .env from .env.example" -ForegroundColor Green
            }
        } else {
            Write-Host "[OK] .env file exists." -ForegroundColor Green
        }
    }

    "up" {
        Write-Host "==> Starting PayBridge Docker containers..." -ForegroundColor Yellow
        if (Get-Command docker -ErrorAction SilentlyContinue) {
            docker compose -f $ComposeFile up -d
            docker compose -f $ComposeFile ps
        } else {
            Write-Host "[FAIL] Docker command not available. Please install Docker Desktop for Windows." -ForegroundColor Red
        }
    }

    "down" {
        Write-Host "==> Stopping PayBridge Docker containers..." -ForegroundColor Yellow
        if (Get-Command docker -ErrorAction SilentlyContinue) {
            docker compose -f $ComposeFile down
        } else {
            Write-Host "[FAIL] Docker command not available." -ForegroundColor Red
        }
    }

    "status" {
        Write-Host "==> PayBridge Infrastructure Status:" -ForegroundColor Cyan
        if (Get-Command docker -ErrorAction SilentlyContinue) {
            docker compose -f $ComposeFile ps
        } else {
            Write-Host "[FAIL] Docker command not available." -ForegroundColor Red
        }
    }

    "test" {
        Write-Host "==> Running Go backend tests..." -ForegroundColor Yellow
        if (Test-Path "services/api/go.mod") {
            Push-Location services/api
            go test -v ./...
            Pop-Location
        } else {
            Write-Host "services/api/go.mod not yet created." -ForegroundColor Yellow
        }
    }

    "lint" {
        Write-Host "==> Running Go vet..." -ForegroundColor Yellow
        if (Test-Path "services/api/go.mod") {
            Push-Location services/api
            go vet ./...
            Pop-Location
        }
    }

    "migrate-up" {
        Write-Host "==> Running goose migrations..." -ForegroundColor Yellow
        if (Get-Command goose -ErrorAction SilentlyContinue) {
            goose -dir services/api/db/migrations postgres $DbUrl up
        } else {
            Write-Host "[FAIL] goose not found. Run: go install github.com/pressly/goose/v3/cmd/goose@latest" -ForegroundColor Red
        }
    }

    "migrate-status" {
        Write-Host "==> Checking goose migration status..." -ForegroundColor Yellow
        if (Get-Command goose -ErrorAction SilentlyContinue) {
            goose -dir services/api/db/migrations postgres $DbUrl status
        } else {
            Write-Host "[FAIL] goose not found." -ForegroundColor Red
        }
    }

    default {
        Show-Help
    }
}
