package utils

import (
	"context"
	"time"

	"gorm.io/gorm/logger"
)

// Структура GormLogger реализует интерфейс логгера для GORM
type GormLogger struct{}

// LogMode реализует метод интерфейса logger.Interface для установки уровня логирования
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	InfoSrc(GormSource, msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	WarnSrc(GormSource, msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	ErrorSrc(GormSource, msg, data...)
}

// Trace логирует выполнение SQL-запроса, время выполнения, количество строк и ошибку (если есть)
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	msg, rows := fc()
	if err != nil {
		ErrorSrc(GormSource, "%s | %v | %d rows | %v", msg, elapsed, rows, err)
	} else {
		InfoSrc(GormSource, "%s | %v | %d rows", msg, elapsed, rows)
	}
}
