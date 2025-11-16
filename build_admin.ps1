# vServer Admin Panel Builder
$ErrorActionPreference = 'Stop'
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

Clear-Host
Start-Sleep -Milliseconds 100

$script:LogBuffer = @()
$script:MaxLogLines = 8

function Write-Step {
    param($Step, $Total, $Message)
    Write-Host "[$Step/$Total] " -ForegroundColor Cyan -NoNewline
    Write-Host $Message -ForegroundColor White
}

function Write-Success {
    param($Message)
    Write-Host "     ✓ " -ForegroundColor Green -NoNewline
    Write-Host $Message -ForegroundColor Gray
}

function Write-Info {
    param($Message)
    Write-Host "     → " -ForegroundColor Yellow -NoNewline
    Write-Host $Message -ForegroundColor Gray
}

function Write-Err {
    param($Message, $Details = "")
    Write-Host ""
    Write-Host "     ✗ ОШИБКА: " -ForegroundColor Red -NoNewline
    Write-Host $Message -ForegroundColor Red
    if ($Details) {
        Write-Host ""
        Write-Host "     Детали:" -ForegroundColor Yellow
        $Details -split "`n" | Select-Object -First 15 | ForEach-Object {
            $line = $_.Trim()
            if ($line) {
                Write-Host "       $line" -ForegroundColor DarkYellow
            }
        }
    }
}

function Write-ProgressBar {
    param($Percent, $Message = "")
    $filled = [math]::Floor($Percent / 4)
    $empty = 25 - $filled
    $bar = "█" * $filled + "░" * $empty
    Write-Host "     [$bar] $Percent%" -ForegroundColor Cyan -NoNewline
    if ($Message) {
        Write-Host " - $Message" -ForegroundColor DarkGray
    } else {
        Write-Host ""
    }
}

function Update-LogWindow {
    param($NewLine)
    
    $script:LogBuffer += $NewLine
    if ($script:LogBuffer.Count -gt $script:MaxLogLines) {
        $script:LogBuffer = $script:LogBuffer[(-$script:MaxLogLines)..-1]
    }
    
    $cursorTop = [Console]::CursorTop
    Write-Host ("     ┌" + "─" * 70 + "┐") -ForegroundColor DarkGray
    
    for ($i = 0; $i -lt $script:MaxLogLines; $i++) {
        if ($i -lt $script:LogBuffer.Count) {
            $line = $script:LogBuffer[$i]
            if ($line.Length -gt 68) {
                $line = $line.Substring(0, 65) + "..."
            }
            Write-Host "     │ " -ForegroundColor DarkGray -NoNewline
            Write-Host $line.PadRight(68) -ForegroundColor DarkCyan -NoNewline
            Write-Host "│" -ForegroundColor DarkGray
        } else {
            Write-Host ("     │" + " " * 70 + "│") -ForegroundColor DarkGray
        }
    }
    
    Write-Host ("     └" + "─" * 70 + "┘") -ForegroundColor DarkGray
}

function Invoke-WithLiveLog {
    param(
        [scriptblock]$ScriptBlock,
        [string]$Activity
    )
    
    $script:LogBuffer = @()
    Update-LogWindow ""
    
    $job = Start-Job -ScriptBlock $ScriptBlock -ArgumentList $PWD
    
    $spinnerChars = @('⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏')
    $spinnerIndex = 0
    
    while ($job.State -eq 'Running') {
        $output = Receive-Job -Job $job 2>&1
        if ($output) {
            foreach ($line in $output) {
                $lineStr = $line.ToString().Trim()
                if ($lineStr -and $lineStr -notmatch '^\s*$') {
                    $script:LogBuffer += $lineStr
                    if ($script:LogBuffer.Count -gt $script:MaxLogLines) {
                        $script:LogBuffer = $script:LogBuffer[(-$script:MaxLogLines)..-1]
                    }
                }
            }
        }
        
        [Console]::SetCursorPosition(0, [Console]::CursorTop - $script:MaxLogLines - 2)
        Write-Host ("     ┌" + "─" * 70 + "┐") -ForegroundColor DarkGray
        
        for ($i = 0; $i -lt $script:MaxLogLines; $i++) {
            if ($i -lt $script:LogBuffer.Count) {
                $line = $script:LogBuffer[$i]
                if ($line.Length -gt 68) {
                    $line = $line.Substring(0, 65) + "..."
                }
                Write-Host "     │ " -ForegroundColor DarkGray -NoNewline
                Write-Host $line.PadRight(68) -ForegroundColor DarkCyan -NoNewline
                Write-Host "│" -ForegroundColor DarkGray
            } else {
                Write-Host ("     │" + " " * 70 + "│") -ForegroundColor DarkGray
            }
        }
        
        Write-Host "     └─ " -ForegroundColor DarkGray -NoNewline
        Write-Host $spinnerChars[$spinnerIndex] -ForegroundColor Cyan -NoNewline
        Write-Host (" $Activity " + "─" * (64 - $Activity.Length)) -ForegroundColor DarkGray -NoNewline
        Write-Host "┘" -ForegroundColor DarkGray
        
        $spinnerIndex = ($spinnerIndex + 1) % $spinnerChars.Count
        Start-Sleep -Milliseconds 100
    }
    
    $output = Receive-Job -Job $job 2>&1
    if ($output) {
        foreach ($line in $output) {
            $lineStr = $line.ToString().Trim()
            if ($lineStr -and $lineStr -notmatch '^\s*$') {
                $script:LogBuffer += $lineStr
            }
        }
    }
    
    $result = @{
        Success = $job.State -eq 'Completed'
        Output = $script:LogBuffer -join "`n"
        Error = $null
    }
    
    if ($job.State -ne 'Completed') {
        $result.Error = "Job failed with state: $($job.State)"
    }
    
    Remove-Job -Job $job -Force
    return $result
}

function Refresh-EnvironmentPath {
    $goBinPath = Join-Path $env:USERPROFILE "go\bin"
    if (Test-Path $goBinPath) {
        $env:Path = "$goBinPath;$env:Path"
    }
}

Write-Host ""
Write-Host "═══════════════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
Write-Host "                        vServer Admin Panel Builder" -ForegroundColor Cyan
Write-Host "═══════════════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
Write-Host ""

try {
    $projectRoot = $PWD.Path
    
    Write-Step 1 6 "Проверка go.mod..."
    Write-Host ""
    
    if (-not (Test-Path "go.mod")) {
        Write-Info "Создание go.mod..."
        $result = go mod init vServer 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Success "go.mod создан"
        } else {
            throw "Не удалось создать go.mod: $result"
        }
    } else {
        Write-Success "go.mod найден"
    }
    
    Write-ProgressBar 16
    Write-Host ""
    Start-Sleep -Milliseconds 300

    Write-Step 2 6 "Установка зависимостей..."
    Write-Host ""
    Write-Info "Загрузка и установка Go модулей (это может занять время)..."
    Write-Host ""
    
    $result = Invoke-WithLiveLog -Activity "Установка зависимостей" -ScriptBlock {
        param($dir)
        Set-Location $dir
        go mod tidy 2>&1
    }
    
    if (-not $result.Success) {
        throw "Ошибка установки зависимостей:`n$($result.Output)"
    }
    
    Write-Host ""
    Write-Success "Все зависимости установлены"
    Write-ProgressBar 33
    Write-Host ""
    Start-Sleep -Milliseconds 300

    Write-Step 3 6 "Проверка Wails CLI..."
    Write-Host ""
    
    Refresh-EnvironmentPath
    
    $wailsPath = Join-Path $env:USERPROFILE "go\bin\wails.exe"
    $wailsExists = Test-Path $wailsPath
    
    if ($wailsExists) {
        try {
            $wailsCheck = & $wailsPath version 2>&1
            $version = ($wailsCheck | Select-String -Pattern 'v\d+\.\d+\.\d+').Matches.Value
            Write-Success "Wails CLI найден ($version)"
        } catch {
            $wailsExists = $false
        }
    }
    
    if (-not $wailsExists) {
        Write-Info "Установка Wails CLI..."
        Write-Host ""
        
        $result = Invoke-WithLiveLog -Activity "Установка Wails" -ScriptBlock {
            param($dir)
            Set-Location $dir
            go install github.com/wailsapp/wails/v2/cmd/wails@latest 2>&1
        }
        
        if (-not $result.Success) {
            throw "Ошибка установки Wails CLI:`n$($result.Output)"
        }
        
        Refresh-EnvironmentPath
        
        if (-not (Test-Path $wailsPath)) {
            throw "Wails CLI установлен, но исполняемый файл не найден по пути: $wailsPath"
        }
        
        Write-Host ""
        Write-Success "Wails CLI установлен"
    }
    
    Write-ProgressBar 50
    Write-Host ""
    Start-Sleep -Milliseconds 300

    Write-Step 4 6 "Проверка конфигурации проекта..."
    Write-Host ""
    
    if (-not (Test-Path "wails.json")) {
        Write-Info "Создание wails.json..."
        
        $wailsConfig = @{
            '$schema' = "https://wails.io/schemas/config.v2.json"
            name = "vServer-Admin"
            outputfilename = "vServer-Admin"
            frontend = @{
                install = ""
                build = ""
                dev = ""
            }
            author = @{
                name = "vServer"
            }
        } | ConvertTo-Json -Depth 10
        
        $wailsConfig | Out-File -FilePath "wails.json" -Encoding UTF8
        Write-Success "wails.json создан"
    } else {
        Write-Success "wails.json найден"
    }
    
    Write-ProgressBar 66
    Write-Host ""
    Start-Sleep -Milliseconds 300

    Write-Step 5 6 "Генерация биндингов..."
    Write-Host ""
    Write-Info "Создание TypeScript/JS интерфейсов для Go методов..."
    Write-Host ""
    
    $result = Invoke-WithLiveLog -Activity "Генерация биндингов" -ScriptBlock {
        param($dir)
        Set-Location $dir
        $wailsPath = Join-Path $env:USERPROFILE "go\bin\wails.exe"
        & $wailsPath generate module 2>&1
    }
    
    if (-not $result.Success) {
        throw "Ошибка генерации биндингов:`n$($result.Output)"
    }
    
    Write-Host ""
    Write-Success "Биндинги успешно сгенерированы"
    Write-ProgressBar 83
    Write-Host ""
    Start-Sleep -Milliseconds 300

    Write-Step 6 6 "Сборка приложения..."
    Write-Host ""
    Write-Info "Компиляция приложения (обычно занимает 10-60 секунд)..."
    Write-Host ""
    
    $result = Invoke-WithLiveLog -Activity "Компиляция" -ScriptBlock {
        param($dir)
        Set-Location $dir
        $wailsPath = Join-Path $env:USERPROFILE "go\bin\wails.exe"
        & $wailsPath build -f admin.go 2>&1
    }
    
    Write-Host ""
    
    if (Test-Path "bin\vServer-Admin.exe") {
        Write-Success "Приложение успешно скомпилировано"
        Write-ProgressBar 100
        Write-Host ""
        
        Write-Info "Финализация сборки..."
        
        try {
            Move-Item -Path "bin\vServer-Admin.exe" -Destination "vSerf.exe" -Force -ErrorAction Stop
            Write-Success "Исполняемый файл перемещён: vSerf.exe"
        } catch {
            throw "Не удалось переместить файл: $_"
        }
        
        if (Test-Path "bin") { 
            Remove-Item -Path "bin" -Recurse -Force -ErrorAction SilentlyContinue
            Write-Success "Папка bin удалена"
        }
        if (Test-Path "windows") { 
            Remove-Item -Path "windows" -Recurse -Force -ErrorAction SilentlyContinue
            Write-Success "Папка windows удалена"
        }
        
        Write-Host ""
        Write-Host "═══════════════════════════════════════════════════════════════════════════" -ForegroundColor Green
        Write-Host "                            ✓ СБОРКА ЗАВЕРШЕНА!" -ForegroundColor Green
        Write-Host "═══════════════════════════════════════════════════════════════════════════" -ForegroundColor Green
        Write-Host ""
        Write-Host "     Исполняемый файл: " -ForegroundColor Gray -NoNewline
        Write-Host "vSerf.exe" -ForegroundColor Cyan
        Write-Host "     Размер: " -ForegroundColor Gray -NoNewline
        $fileSize = [math]::Round((Get-Item "vSerf.exe").Length / 1MB, 2)
        Write-Host "$fileSize MB" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "═══════════════════════════════════════════════════════════════════════════" -ForegroundColor Green
        
    } else {
        $errorOutput = $result.Output
        if ($errorOutput -match "ERROR.*") {
            $errorLines = ($errorOutput -split "`n" | Where-Object { $_ -match "ERROR|error|failed" }) -join "`n"
            throw "Ошибка компиляции:`n$errorLines`n`nПолный вывод:`n$errorOutput"
        } else {
            throw "Исполняемый файл не найден после компиляции.`nВывод команды:`n$errorOutput"
        }
    }

} catch {
    Write-Host ""
    Write-Host "═══════════════════════════════════════════════════════════════════════════" -ForegroundColor Red
    Write-Host "                            ✗ ОШИБКА СБОРКИ!" -ForegroundColor Red
    Write-Host "═══════════════════════════════════════════════════════════════════════════" -ForegroundColor Red
    Write-Err $_.Exception.Message $_.ScriptStackTrace
    Write-Host ""
    Write-Host "═══════════════════════════════════════════════════════════════════════════" -ForegroundColor Red
    Write-Host ""
    exit 1
}

Write-Host ""