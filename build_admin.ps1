# vServer Admin Panel Builder
$ErrorActionPreference = 'SilentlyContinue'
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

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

function Check-And-Install {
    param($Name, $Command, $WingetId)

    $found = Get-Command $Command -ErrorAction SilentlyContinue

    if (-not $found) {
        Write-Err "$Name not found!"
        $answer = Read-Host "     ? Install $Name via winget? (y/n)"
        if ($answer -eq 'y' -or $answer -eq 'Y') {
            Write-Info "Installing $Name..."
            winget install --id $WingetId --accept-source-agreements --accept-package-agreements 2>&1 | Out-Null
            if ($LASTEXITCODE -eq 0) {
                Write-Success "$Name installed! Restart terminal and run script again."
            } else {
                Write-Err "Failed to install. Install manually: winget install $WingetId"
            }
            Write-Host ""
            exit 1
        } else {
            Write-Err "Cannot continue without $Name"
            Write-Host ""
            exit 1
        }
    } else {
        Write-Success "$Name found"
    }
}

Write-Host ""
Write-Host "=================================================" -ForegroundColor Magenta
Write-Host "    vServer Admin Panel Builder" -ForegroundColor Cyan
Write-Host "=================================================" -ForegroundColor Magenta
Write-Host ""

# Step 0: Check dependencies
Write-Step 1 6 "Checking dependencies..."
Check-And-Install "Go" "go" "GoLang.Go"
Check-And-Install "Node.js" "node" "OpenJS.NodeJS.LTS"
Check-And-Install "npm" "npm" "OpenJS.NodeJS.LTS"
Write-ProgressBar 10
Write-Host ""

# Step 1: go.mod
Write-Step 2 6 "Checking go.mod..."
if (-not (Test-Path "go.mod")) {
    Write-Info "Creating go.mod..."
    go mod init vServer 2>&1 | Out-Null
    Write-Success "Created"
} else {
    Write-Success "Found"
}
Write-ProgressBar 20
Write-Host ""

# Step 2: Go dependencies
Write-Step 3 6 "Installing Go dependencies..."
go mod tidy 2>&1 | Out-Null
Write-Success "Dependencies installed"
Write-ProgressBar 40
Write-Host ""

# Step 3: Wails CLI
Write-Step 4 6 "Checking Wails CLI..."
$null = wails version 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Info "Installing Wails CLI..."
    go install github.com/wailsapp/wails/v2/cmd/wails@latest 2>&1 | Out-Null
    Write-Success "Installed"
} else {
    Write-Success "Found"
}
Write-ProgressBar 55
Write-Host ""

# Step 4: Vue frontend
Write-Step 5 6 "Building Vue frontend..."
Push-Location "front_vue"
Write-Info "npm install..."
npm install 2>&1 | Out-Null
Write-Info "npm run build..."
npm run build 2>&1 | Out-Null
Pop-Location
Write-Success "Frontend built"
Write-ProgressBar 75
Write-Host ""

# Step 5: Wails build
Write-Step 6 6 "Building application..."
Write-Info "Compiling (may take ~30 sec)..."
wails build 2>&1 | Out-Null

$exePath = $null
if (Test-Path "build\bin\vServer-Admin.exe") {
    $exePath = "build\bin\vServer-Admin.exe"
}
if ((-not $exePath) -and (Test-Path "bin\vServer-Admin.exe")) {
    $exePath = "bin\vServer-Admin.exe"
}

if ($exePath) {
    Write-Success "Compiled"
    Write-ProgressBar 100
    Write-Host ""
    Write-Host "Finalizing..." -ForegroundColor Cyan
    Move-Item -Path $exePath -Destination "vSerf.exe" -Force -ErrorAction SilentlyContinue
    Write-Success "File moved: vSerf.exe"
    if (Test-Path "build") { Remove-Item -Path "build" -Recurse -Force -ErrorAction SilentlyContinue }
    if (Test-Path "bin") { Remove-Item -Path "bin" -Recurse -Force -ErrorAction SilentlyContinue }
    if (Test-Path "windows") { Remove-Item -Path "windows" -Recurse -Force -ErrorAction SilentlyContinue }
    Write-Success "Temp files removed"
    Write-Host ""
    Write-Host "=================================================" -ForegroundColor Green
    Write-Host "  BUILD SUCCESS!" -ForegroundColor Green
    Write-Host "  File: " -ForegroundColor Green -NoNewline
    Write-Host "vSerf.exe" -ForegroundColor Cyan
    Write-Host "=================================================" -ForegroundColor Green
} else {
    Write-Err "Compilation error"
    Write-Host ""
    Write-Host "=================================================" -ForegroundColor Red
    Write-Host "  BUILD FAILED!" -ForegroundColor Red
    Write-Host "=================================================" -ForegroundColor Red
}

Write-Host ""
