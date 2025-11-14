# vServer Admin Panel Builder
$ErrorActionPreference = 'SilentlyContinue'
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

# Очищаем консоль
Clear-Host
Start-Sleep -Milliseconds 100

function Write-Step {
    param($Step, $Total, $Message)
    Write-Host "[$Step/$Total] " -ForegroundColor Cyan -NoNewline
    Write-Host $Message -ForegroundColor White
}

function Write-Success {
    param($Message)
    Write-Host "     + OK: " -ForegroundColor Green -NoNewline
    Write-Host $Message -ForegroundColor Green
}

function Write-Info {
    param($Message)
    Write-Host "     > " -ForegroundColor Yellow -NoNewline
    Write-Host $Message -ForegroundColor Yellow
}

function Write-Err {
    param($Message)
    Write-Host "     X ERROR: " -ForegroundColor Red -NoNewline
    Write-Host $Message -ForegroundColor Red
}

function Write-ProgressBar {
    param($Percent)
    $filled = [math]::Floor($Percent / 4)
    $empty = 25 - $filled
    $bar = "#" * $filled + "-" * $empty
    Write-Host "     [$bar] $Percent%" -ForegroundColor Cyan
}

Write-Host ""
Write-Host "=================================================" -ForegroundColor Magenta
Write-Host "    vServer Admin Panel Builder" -ForegroundColor Cyan
Write-Host "=================================================" -ForegroundColor Magenta
Write-Host ""

Write-Step 1 4 "Проверка go.mod..."
if (-not (Test-Path "go.mod")) {
    Write-Info "Создание go.mod..."
    go mod init vServer 2>&1 | Out-Null
    Write-Success "Создан"
} else {
    Write-Success "Найден"
}
Write-ProgressBar 25
Write-Host ""

Write-Step 2 4 "Установка зависимостей..."
go mod tidy 2>&1 | Out-Null
Write-Success "Зависимости установлены"
Write-ProgressBar 50
Write-Host ""

Write-Step 3 4 "Проверка Wails CLI..."
$null = wails version 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Info "Установка Wails CLI..."
    go install github.com/wailsapp/wails/v2/cmd/wails@latest 2>&1 | Out-Null
    Write-Success "Установлен"
} else {
    Write-Success "Найден"
}
Write-ProgressBar 75
Write-Host ""

Write-Step 4 4 "Сборка приложения..."
Write-Info "Компиляция (может занять ~10 сек)..."

wails build -f admin.go 2>&1 | Out-Null

if (Test-Path "bin\vServer-Admin.exe") {
    Write-Success "Скомпилировано"
    Write-ProgressBar 100
    Write-Host ""
    
    Write-Host "Финализация..." -ForegroundColor Cyan
    Move-Item -Path "bin\vServer-Admin.exe" -Destination "vSerf.exe" -Force 2>$null
    Write-Success "Файл перемещён: vSerf.exe"
    
    if (Test-Path "bin") { Remove-Item -Path "bin" -Recurse -Force 2>$null }
    if (Test-Path "windows") { Remove-Item -Path "windows" -Recurse -Force 2>$null }
    Write-Success "Временные файлы удалены"
    Write-Host ""
    
    Write-Host "=================================================" -ForegroundColor Green
    Write-Host "  УСПЕШНО СОБРАНО!" -ForegroundColor Green
    Write-Host "  Файл: " -ForegroundColor Green -NoNewline
    Write-Host "vSerf.exe" -ForegroundColor Cyan
    Write-Host "=================================================" -ForegroundColor Green
} else {
    Write-Err "Ошибка компиляции"
    Write-Host ""
    
    Write-Host "=================================================" -ForegroundColor Red
    Write-Host "  ОШИБКА СБОРКИ!" -ForegroundColor Red
    Write-Host "=================================================" -ForegroundColor Red
}

Write-Host ""
