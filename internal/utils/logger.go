package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

// Уровни логирования
// Используются для фильтрации сообщений по важности
type LogLevel int

const (
	LevelError LogLevel = iota // Ошибки
	LevelWarn                  // Предупреждения
	LevelInfo                  // Информационные сообщения
)

// Цвета для вывода логов в консоль
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGreen  = "\033[32m"
)

// Глобальные переменные для логирования
var (
	logLevel    = LevelInfo // Текущий уровень логирования
	once        sync.Once   // Для инициализации логгера
	enableColor = true      // Включение/отключение цвета
)

// Устанавливает уровень логирования
func SetLogLevel(level LogLevel) {
	logLevel = level
}

// Включает или отключает цветной вывод логов
func SetLogColor(enable bool) {
	enableColor = enable
}

// Возвращает префикс для сообщения лога в зависимости от уровня
func getPrefix(level LogLevel) string {
	if !enableColor {
		switch level {
		case LevelError:
			return "[ERROR] "
		case LevelWarn:
			return "[WARN] "
		case LevelInfo:
			return "[INFO] "
		default:
			return "[LOG] "
		}
	}
	switch level {
	case LevelError:
		return colorRed + "[ERROR]" + colorReset + " "
	case LevelWarn:
		return colorYellow + "[WARN]" + colorReset + " "
	case LevelInfo:
		return colorBlue + "[INFO]" + colorReset + " "
	default:
		return colorGreen + "[LOG]" + colorReset + " "
	}
}

// Форматирует и выводит лог-сообщение
func logf(level LogLevel, format string, v ...interface{}) {
	if level > logLevel {
		return
	}
	once.Do(func() {
		log.SetOutput(os.Stdout)
		log.SetFlags(0)
	})
	// Определяем имя файла, функцию и строку вызова
	pc, file, line, ok := runtime.Caller(2)
	funcName := "?"
	fileName := file
	if ok {
		if lastSlash := lastIndex(file, "/"); lastSlash != -1 {
			fileName = file[lastSlash+1:]
		} else if lastSlash := lastIndex(file, "\\"); lastSlash != -1 {
			fileName = file[lastSlash+1:]
		}
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			funcName = fn.Name()
			if lastSlash := lastIndex(funcName, "/"); lastSlash != -1 {
				funcName = funcName[lastSlash+1:]
			}
			if dot := lastIndex(funcName, "."); dot != -1 {
				funcName = funcName[dot+1:]
			}
		}
	}
	// Формируем и выводим сообщение
	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006/01/02 - 15:04:05")
	log.Printf("%s%s [%s:%d %s] %s", getPrefix(level), timestamp, fileName, line, funcName, msg)
}

// Вспомогательная функция для поиска последнего вхождения подстроки
func lastIndex(s, sep string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if len(s)-i >= len(sep) && s[i:i+len(sep)] == sep {
			return i
		}
	}
	return -1
}

// Логирует ошибку
func Error(format string, v ...interface{}) {
	logf(LevelError, format, v...)
}

// Логирует предупреждение
func Warn(format string, v ...interface{}) {
	logf(LevelWarn, format, v...)
}

// Логирует информационное сообщение
func Info(format string, v ...interface{}) {
	logf(LevelInfo, format, v...)
}
