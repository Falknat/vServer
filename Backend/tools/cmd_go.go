//go:build windows

package tools

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procCreateMutex    = kernel32.NewProc("CreateMutexW")
	procCloseHandle    = kernel32.NewProc("CloseHandle")
)

const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004

var mutexHandle syscall.Handle

func init() {
	enableVirtualTerminal()

}

func enableVirtualTerminal() {
	handle := os.Stdout.Fd()
	var mode uint32

	// Получаем текущий режим консоли
	_, _, _ = procGetConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))

	// Добавляем флаг поддержки ANSI
	mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING

	// Устанавливаем новый режим
	_, _, _ = procSetConsoleMode.Call(uintptr(handle), uintptr(mode))
}

func RunBatScript(script string) (string, error) {
	// Создание временного файла
	tmpFile, err := os.CreateTemp("", "script-*.bat")
	if err != nil {
		return "", fmt.Errorf("ошибка создания temp-файла: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Запись скрипта в файл
	if _, err := tmpFile.WriteString(script); err != nil {
		return "", fmt.Errorf("ошибка записи в temp-файл: %w", err)
	}
	tmpFile.Close()

	// Выполняем файл через cmd
	cmd := exec.Command("cmd", "/C", tmpFile.Name())
	output, err := cmd.CombinedOutput()

	return string(output), err
}

// Функция для логирования вывода процесса в консоль
func Logs_console(process *exec.Cmd, check bool) error {

	if check {
		// Настраиваем pipes для захвата вывода
		stdout, err := process.StdoutPipe()
		CheckError(err)
		stderr, err := process.StderrPipe()
		CheckError(err)

		// Запускаем процесс
		process.Start()

		// Захватываем stdout и stderr для вывода логов
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
		}()
	} else {
		// Просто запускаем процесс без логирования
		return process.Start()
	}

	return nil
}

// CheckSingleInstance проверяет, не запущена ли программа уже через мьютекс
func CheckSingleInstance() bool {
	mutexName, _ := syscall.UTF16PtrFromString("Global\\vServer_SingleInstance")
	
	handle, _, err := procCreateMutex.Call(
		0,
		0,
		uintptr(unsafe.Pointer(mutexName)),
	)
	
	if handle == 0 {
		return false // не удалось создать мьютекс
	}
	
	mutexHandle = syscall.Handle(handle)
	
	// Если GetLastError возвращает ERROR_ALREADY_EXISTS (183), значит мьютекс уже существует
	if err.(syscall.Errno) == 183 {
		return false // программа уже запущена
	}
	
	return true // успешно создали мьютекс, программа не запущена
}

// ReleaseMutex освобождает мьютекс при завершении программы
func ReleaseMutex() {
	if mutexHandle != 0 {
		procCloseHandle.Call(uintptr(mutexHandle))
		mutexHandle = 0
	}
}