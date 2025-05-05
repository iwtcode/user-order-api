package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

type LogLevel int

const (
	LevelError LogLevel = iota
	LevelWarn
	LevelInfo
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGreen  = "\033[32m"
)

var (
	logLevel    = LevelInfo
	once        sync.Once
	enableColor = true // set to false to disable color output
)

func SetLogLevel(level LogLevel) {
	logLevel = level
}

func SetLogColor(enable bool) {
	enableColor = enable
}

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
	// With color
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

func logf(level LogLevel, format string, v ...interface{}) {
	if level > logLevel {
		return
	}
	once.Do(func() {
		log.SetOutput(os.Stdout)
		log.SetFlags(0)
	})
	// Получаем информацию о месте вызова
	pc, file, line, ok := runtime.Caller(2)
	funcName := "?"
	fileName := file
	if ok {
		// Оставляем только имя файла
		if lastSlash := lastIndex(file, "/"); lastSlash != -1 {
			fileName = file[lastSlash+1:]
		} else if lastSlash := lastIndex(file, "\\"); lastSlash != -1 {
			fileName = file[lastSlash+1:]
		}
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			// Оставляем только имя функции
			funcName = fn.Name()
			if lastSlash := lastIndex(funcName, "/"); lastSlash != -1 {
				funcName = funcName[lastSlash+1:]
			}
			if dot := lastIndex(funcName, "."); dot != -1 {
				funcName = funcName[dot+1:]
			}
		}
	}
	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006/01/02 - 15:04:05")
	log.Printf("%s%s [%s:%d %s] %s", getPrefix(level), timestamp, fileName, line, funcName, msg)
}

// lastIndex returns the index of the last instance of sep in s, or -1 if sep is not present.
func lastIndex(s, sep string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if len(s)-i >= len(sep) && s[i:i+len(sep)] == sep {
			return i
		}
	}
	return -1
}

func Error(format string, v ...interface{}) {
	logf(LevelError, format, v...)
}

func Warn(format string, v ...interface{}) {
	logf(LevelWarn, format, v...)
}

func Info(format string, v ...interface{}) {
	logf(LevelInfo, format, v...)
}
