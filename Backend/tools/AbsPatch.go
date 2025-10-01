package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// getAbsPath получает абсолютный путь с автоматической проверкой существования для файлов
func AbsPath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("ошибка получения абсолютного пути: %v", err)
	}

	// Проверяем существование только для файлов (с расширением)
	if strings.Contains(filepath.Base(absPath), ".") {
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			return "", fmt.Errorf("файл не найден: %s", absPath)
		}
	}

	return absPath, nil
}
