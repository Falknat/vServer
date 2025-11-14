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
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleMode           = kernel32.NewProc("SetConsoleMode")
	procGetConsoleMode           = kernel32.NewProc("GetConsoleMode")
	procCreateMutex              = kernel32.NewProc("CreateMutexW")
	procCloseHandle              = kernel32.NewProc("CloseHandle")
	procCreateJobObject          = kernel32.NewProc("CreateJobObjectW")
	procAssignProcessToJobObject = kernel32.NewProc("AssignProcessToJobObject")
	procSetInformationJobObject  = kernel32.NewProc("SetInformationJobObject")
)

const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
const JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE = 0x2000

var mutexHandle syscall.Handle
var jobHandle syscall.Handle

func init() {
	enableVirtualTerminal()
	createJobObject()
}

func createJobObject() {
	handle, _, _ := procCreateJobObject.Call(0, 0)
	if handle == 0 {
		return
	}

	jobHandle = syscall.Handle(handle)

	// Устанавливаем флаг автоматического завершения дочерних процессов
	type JOBOBJECT_EXTENDED_LIMIT_INFORMATION struct {
		BasicLimitInformation struct {
			PerProcessUserTimeLimit uint64
			PerJobUserTimeLimit     uint64
			LimitFlags              uint32
			MinimumWorkingSetSize   uintptr
			MaximumWorkingSetSize   uintptr
			ActiveProcessLimit      uint32
			Affinity                uintptr
			PriorityClass           uint32
			SchedulingClass         uint32
		}
		IoInfo                [48]byte
		ProcessMemoryLimit    uintptr
		JobMemoryLimit        uintptr
		PeakProcessMemoryUsed uintptr
		PeakJobMemoryUsed     uintptr
	}

	var limitInfo JOBOBJECT_EXTENDED_LIMIT_INFORMATION
	limitInfo.BasicLimitInformation.LimitFlags = JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE

	procSetInformationJobObject.Call(
		uintptr(jobHandle),
		9, // JobObjectExtendedLimitInformation
		uintptr(unsafe.Pointer(&limitInfo)),
		unsafe.Sizeof(limitInfo),
	)

	// Добавляем текущий процесс в Job Object
	currentProcess, _ := syscall.GetCurrentProcess()
	procAssignProcessToJobObject.Call(uintptr(jobHandle), uintptr(currentProcess))
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

	// Скрываем окно процесса для GUI приложений
	process.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}

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

	// Закрываем Job Object - это автоматически убьёт все дочерние процессы
	if jobHandle != 0 {
		procCloseHandle.Call(uintptr(jobHandle))
		jobHandle = 0
	}
}
