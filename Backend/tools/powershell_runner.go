package tools

import (
	"bufio"
	"os/exec"
	"strings"
	"sync"
)

// Глобальная persistent PowerShell сессия
var (
	psCmd    *exec.Cmd
	psStdin  *bufio.Writer
	psStdout *bufio.Scanner
	psMutex  sync.Mutex
)

// RunPersistentScript выполняет команды через постоянную PowerShell сессию
func RunPersistentScript(commands []string) string {
	psMutex.Lock()
	defer psMutex.Unlock()

	// Инициализируем если еще не запущен
	if psCmd == nil {
		psCmd = exec.Command("powershell", "-NoExit", "-Command", "-")
		stdin, _ := psCmd.StdinPipe()
		stdout, _ := psCmd.StdoutPipe()
		psStdin = bufio.NewWriter(stdin)
		psStdout = bufio.NewScanner(stdout)
		psCmd.Start()
	}

	// Выполняем команды
	fullCommand := strings.Join(commands, "; ")
	psStdin.WriteString(fullCommand + "; Write-Output '---END---'\n")
	psStdin.Flush()

	// Читаем результат - только последняя строка с данными
	var lastLine string
	for psStdout.Scan() {
		line := psStdout.Text()
		if line == "---END---" {
			break
		}
		if strings.TrimSpace(line) != "" {
			lastLine = line
		}
	}

	return lastLine
}

// RunPowerShellCommand выполняет PowerShell команду и возвращает результат
// Если ошибка - возвращает текст ошибки в строке
func RunPScode(command string) string {
	cmd := exec.Command("powershell", "-Command", command)

	output, err := cmd.Output()
	if err != nil {
		return "ERROR: " + err.Error()
	}

	return strings.TrimSpace(string(output))
}

// RunPowerShellScript выполняет несколько команд PowerShell
func RunPowerShellScript(commands []string) string {
	// Объединяем команды через точку с запятой
	fullCommand := strings.Join(commands, "; ")

	return RunPScode(fullCommand)
}
