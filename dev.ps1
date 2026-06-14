# Levanta backend (Go) y frontend (Vite) en un solo terminal.
# Uso: .\dev.ps1

$ErrorActionPreference = "Stop"
$Root = $PSScriptRoot

Set-Location $Root

if (-not (Test-Path ".\node_modules\concurrently")) {
  Write-Host "Instalando dependencias del monorepo..." -ForegroundColor Cyan
  npm install
}

if (-not (Test-Path ".\frontend\node_modules")) {
  Write-Host "Instalando dependencias del frontend..." -ForegroundColor Cyan
  Set-Location ".\frontend"
  npm install
  Set-Location $Root
}

npm run dev
