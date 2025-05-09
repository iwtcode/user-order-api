package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// Структура customFormatter определяет форматтер для логирования с возможностью указания источника и использования цветов
// Source - источник лога, Colors - использовать ли цветной вывод
type customFormatter struct {
	Source string
	Colors bool
}

// Метод Format реализует интерфейс logrus.Formatter и форматирует запись лога по заданному шаблону
func (f *customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("02-01-2006 15:04:05")
	level := formatLevel(entry.Level, f.Colors)
	source := f.Source
	if src, ok := entry.Data["src"].(string); ok {
		source = src
	}
	msg := fmt.Sprintf("[%s] [%s] [%s] %s\n", level, timestamp, source, entry.Message)
	return []byte(msg), nil
}

// Функция formatLevel возвращает строковое представление уровня логирования с цветом или без
func formatLevel(level logrus.Level, colors bool) string {
	switch level {
	case logrus.InfoLevel:
		if colors {
			return "\033[36mINFO\033[0m"
		}
		return "INFO"
	case logrus.WarnLevel:
		if colors {
			return "\033[33mWARN\033[0m"
		}
		return "WARN"
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		if colors {
			return "\033[31mERROR\033[0m"
		}
		return "ERROR"
	case logrus.DebugLevel:
		if colors {
			return "\033[35mDEBUG\033[0m"
		}
		return "DEBUG"
	default:
		return strings.ToUpper(level.String())
	}
}

// LogMessage описывает структуру сообщения для асинхронного логирования
// level - уровень, src - источник, msg - текст
// fileLog - писать ли в файл
type LogMessage struct {
	Level   logrus.Level
	Src     string
	Msg     string
	FileLog bool
}

// Глобальные переменные для логгеров консоли и файла, а также для однократной инициализации
var (
	consoleLogger *logrus.Logger
	fileLogger    *logrus.Logger
	initOnce      sync.Once
	logChan       chan LogMessage
	workerOnce    sync.Once
)

// startLogWorker запускает воркер для асинхронного логирования
func startLogWorker() {
	workerOnce.Do(func() {
		logChan = make(chan LogMessage, 1000)
		go func() {
			for logMsg := range logChan {
				entry := consoleLogger.WithField("src", logMsg.Src)
				switch logMsg.Level {
				case logrus.ErrorLevel:
					entry.Error(logMsg.Msg)
					if fileLogger != nil && logMsg.FileLog {
						fileLogger.WithField("src", logMsg.Src).Error(logMsg.Msg)
					}
				case logrus.WarnLevel:
					entry.Warn(logMsg.Msg)
					if fileLogger != nil && logMsg.FileLog {
						fileLogger.WithField("src", logMsg.Src).Warn(logMsg.Msg)
					}
				case logrus.InfoLevel:
					entry.Info(logMsg.Msg)
					if fileLogger != nil && logMsg.FileLog {
						fileLogger.WithField("src", logMsg.Src).Info(logMsg.Msg)
					}
				case logrus.DebugLevel:
					entry.Debug(logMsg.Msg)
					if fileLogger != nil && logMsg.FileLog {
						fileLogger.WithField("src", logMsg.Src).Debug(logMsg.Msg)
					}
				}
			}
		}()
	})
}

// InitLogger инициализирует логгеры для консоли и файла (если указан путь к файлу)
func InitLogger(logFile string) {
	initOnce.Do(func() {
		consoleLogger = logrus.New()
		consoleLogger.SetFormatter(&customFormatter{Source: "LOG", Colors: true})
		consoleLogger.SetReportCaller(false)
		consoleLogger.SetOutput(os.Stdout)

		if logFile != "" {
			dir := filepath.Dir(logFile)
			if err := os.MkdirAll(dir, 0755); err != nil {
				consoleLogger.Warnf("Failed to create log directory '%s': %v. Logging only to console.", dir, err)
			} else {
				file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
				if err == nil {
					fileLogger = logrus.New()
					fileLogger.SetFormatter(&customFormatter{Source: "LOG", Colors: false})
					fileLogger.SetReportCaller(false)
					fileLogger.SetOutput(file)
				} else {
					consoleLogger.Warnf("Failed to open log file '%s': %v. Logging only to console.", logFile, err)
				}
			}
		}
		startLogWorker()
	})
}

// logWithLevel логирует сообщение на заданном уровне с указанием источника (асинхронно)
func logWithLevel(level logrus.Level, src, format string, v ...interface{}) {
	if consoleLogger == nil {
		InitLogger("")
	}
	startLogWorker()
	msg := fmt.Sprintf(format, v...)
	if src == "LOG" || src == "" {
		if pc, file, line, ok := runtime.Caller(2); ok {
			fileBase := filepath.Base(file)
			fileName := strings.TrimSuffix(fileBase, filepath.Ext(fileBase))
			funcNameFull := runtime.FuncForPC(pc).Name()
			funcParts := strings.Split(funcNameFull, ".")
			funcName := funcParts[len(funcParts)-1]
			if idx := strings.LastIndex(funcName, ")"); idx != -1 && idx+1 < len(funcName) {
				funcName = funcName[idx+1:]
			}
			src = fmt.Sprintf("%s/%s:%d", fileName, funcName, line)
		}
	}
	// Отправляем в канал для асинхронного логирования
	logChan <- LogMessage{
		Level:   level,
		Src:     src,
		Msg:     msg,
		FileLog: true,
	}
}

func Error(format string, v ...interface{})        { logWithLevel(logrus.ErrorLevel, "LOG", format, v...) }
func Warn(format string, v ...interface{})         { logWithLevel(logrus.WarnLevel, "LOG", format, v...) }
func Info(format string, v ...interface{})         { logWithLevel(logrus.InfoLevel, "LOG", format, v...) }
func Debug(format string, v ...interface{})        { logWithLevel(logrus.DebugLevel, "LOG", format, v...) }
func InfoSrc(src, format string, v ...interface{}) { logWithLevel(logrus.InfoLevel, src, format, v...) }
func WarnSrc(src, format string, v ...interface{}) { logWithLevel(logrus.WarnLevel, src, format, v...) }
func ErrorSrc(src, format string, v ...interface{}) {
	logWithLevel(logrus.ErrorLevel, src, format, v...)
}
func DebugSrc(src, format string, v ...interface{}) {
	logWithLevel(logrus.DebugLevel, src, format, v...)
}
