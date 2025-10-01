package json

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"vServer/Backend/tools"
)

var CPU_NAME string
var CPU_GHz string
var CPU_Cores string
var CPU_Using string
var Disk_Size string
var Disk_Free string
var Disk_Used string
var Disk_Name string
var RAM_Using string
var RAM_Total string

// Инициализация при запуске пакета
func init() {
	// Загружаем статичные данные один раз при старте
	UpdateMetrics()
	// Загружаем динамические данные в фоне
	go updateMetricsBackground()
}

func UpdateMetrics() {
	commands := []string{
		`$name = (Get-WmiObject Win32_DiskDrive | Where-Object { $_.Index -eq (Get-WmiObject Win32_DiskPartition | Where-Object { $_.DeviceID -eq ((Get-WmiObject Win32_LogicalDiskToPartition | Where-Object { $_.Dependent -match "C:" }).Antecedent -split '"')[1] }).DiskIndex }).Model`,
		`$size = "{0} GB" -f [math]::Round(((Get-PSDrive -Name C).Used + (Get-PSDrive -Name C).Free) / 1GB)`,
		"$cpuInfo = Get-CimInstance Win32_Processor",
		"$cpuCores = $cpuInfo.NumberOfCores",
		"$cpuName = $cpuInfo.Name",
		`$ram_total = [math]::Round((Get-CimInstance Win32_ComputerSystem).TotalPhysicalMemory / 1GB, 2)`,
		"Write-Output \"$cpuName|$cpuCores|$name|$size|$ram_total\"",
	}

	// Выполняем команды и получаем результат
	result := tools.RunPersistentScript(commands)

	// Парсим результат для статичных данных
	parts := strings.Split(result, "|")
	if len(parts) >= 4 {
		cpuName := strings.TrimSpace(parts[0])
		cpuCores := strings.TrimSpace(parts[1])
		diskName := strings.TrimSpace(parts[2])
		diskSize := strings.TrimSpace(parts[3])
		ramTotal := strings.TrimSpace(parts[4])

		// Обновляем глобальные переменные
		CPU_NAME = cpuName
		CPU_Cores = cpuCores
		Disk_Name = diskName
		Disk_Size = diskSize
		RAM_Total = ramTotal
	}
}

// Фоновое обновление динамических метрик (внутреннее)
func updateMetricsBackground() {

	updateDynamicMetrics := func() {
		commands := []string{
			"$cpuInfo = Get-CimInstance Win32_Processor",
			"$cpuGHz = $cpuInfo.MaxClockSpeed",
			"$cpuUsage = $cpuInfo.LoadPercentage",
			`$used = "{0} GB" -f [math]::Round(((Get-PSDrive -Name C).Used / 1GB))`,
			`$free = "{0} GB" -f [math]::Round(((Get-PSDrive -Name C).Free / 1GB))`,
			`$ram_using = [math]::Round((Get-CimInstance Win32_OperatingSystem | % {($_.TotalVisibleMemorySize - $_.FreePhysicalMemory) / 1MB}), 2)`,
			"Write-Output \"$cpuGHz|$cpuUsage|$used|$free|$ram_using\"",
		}

		// Один запуск PowerShell для динамических команд!
		result := tools.RunPersistentScript(commands)

		// Парсим результат для динамических данных
		parts := strings.Split(result, "|")
		if len(parts) >= 4 {
			cpuGHz := strings.TrimSpace(parts[0])
			cpuUsage := strings.TrimSpace(parts[1])
			diskUsed := strings.TrimSpace(parts[2])
			diskFree := strings.TrimSpace(parts[3])
			ramUsing := strings.TrimSpace(parts[4])
			// Обновляем глобальные переменные
			CPU_GHz = cpuGHz
			CPU_Using = cpuUsage
			Disk_Used = diskUsed
			Disk_Free = diskFree
			RAM_Using = ramUsing
		}
	}

	// Выполняем сразу при запуске
	updateDynamicMetrics()

	// Затем каждые 5 секунд
	for range time.NewTicker(5 * time.Second).C {
		updateDynamicMetrics()
	}

}

// Получить JSON системных метрик
func GetAllMetrics(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(GetMetricsJSON())
}

// Получить JSON метрик
func GetMetricsJSON() []byte {
	metrics := map[string]interface{}{
		"cpu_name":      CPU_NAME,
		"cpu_ghz":       CPU_GHz,
		"cpu_cores":     CPU_Cores,
		"cpu_usage":     CPU_Using,
		"disk_name":     Disk_Name,
		"disk_size":     Disk_Size,
		"disk_used":     Disk_Used,
		"disk_free":     Disk_Free,
		"ram_using":     RAM_Using,
		"ram_total":     RAM_Total,
		"server_uptime": tools.ServerUptime("get"),
	}

	data, _ := json.Marshal(metrics)
	return data
}
