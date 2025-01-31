package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// LoggersInterface .
type LoggersInterface interface {
	Info(message string, args ...interface{})
	Success(message string, args ...interface{})
	Trace(message string, args ...interface{})
	Error(message string, err error)
	Warn(message string, args ...interface{})
	Fatal(message string, err error)
	Debug(message string, args ...interface{})
	Critical(message string, err error)
	Panic(message string, err error)
	Security(message string, args ...interface{})
	Close() error
}

// MyLogger .
type MyLogger struct {
	consoleLogger *log.Logger
	fileLogger    *log.Logger
	file          *os.File
}

// NewLogger создает новый логгер.
// Если logFilePath указан, логи будут записываться в файл и в консоль.
// Если logFilePath пуст, логи будут записываться только в консоль.
func NewLogger(logFilePath string) (LoggersInterface, error) {
	var consoleWriter io.Writer = os.Stdout
	var fileWriter io.Writer
	var logFile *os.File

	// Если указан путь к файлу, открываем его
	if logFilePath != "" {
		dir := filepath.Dir(logFilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("не удалось создать директорию для логов: %w", err)
		}

		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("не удалось открыть файл логов: %w", err)
		}
		logFile = file
		fileWriter = logFile
	} else {
		fileWriter = io.Discard
	}

	// Два раздельных логгера
	consoleLogger := log.New(consoleWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	fileLogger := log.New(fileWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	return &MyLogger{
		consoleLogger: consoleLogger,
		fileLogger:    fileLogger,
		file:          logFile,
	}, nil
}

// Close закрывает файл логов (если он был открыт).
func (l *MyLogger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// formatMessage форматирует сообщение для логирования.
func formatMessage(file string, line int, level string, color string, message string, args ...interface{}) string {
	caller := fmt.Sprintf("%s:%d", filepath.Base(file), line)

	var str strings.Builder
	str.WriteString("[")
	if color != "" {
		str.WriteString(color)
	}
	str.WriteString(level)
	if color != "" {
		str.WriteString(colorReset)
	}
	str.WriteString("]")
	str.WriteString(" ")
	str.WriteString(caller)
	str.WriteString(" ")
	str.WriteString(message)

	if len(args) > 0 {
		str.WriteString(" ")
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				str.WriteString(fmt.Sprintf("%v: %v ", args[i], args[i+1]))
			} else {
				str.WriteString(fmt.Sprintf("%v", args[i]))
			}
		}
	}

	return strings.TrimSpace(str.String())
}

// logWithCallerInfo выводит отформатированное сообщение в лог.
func (l *MyLogger) logWithCallerInfo(file string, line int, level string, color string, message string, args ...interface{}) {
	formattedMessageConsole := formatMessage(file, line, level, color, message, args...)
	formattedMessageFile := formatMessage(file, line, level, "", message, args...)

	if l.consoleLogger != nil {
		l.consoleLogger.Println(formattedMessageConsole)
	}

	if l.fileLogger != nil {
		l.fileLogger.Println(formattedMessageFile)
	}
}

// Info записывает информационное сообщение в лог вместе с контекстом вызова функции.
// Параметры args содержат дополнительные данные для сообщения.
func (l *MyLogger) Info(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if l.consoleLogger != nil {
		l.logWithCallerInfo(file, line, "INFO", colorGreen, message, args...)
	} else {
		log.Print("No logger available.")
	}
}

// Success записывает логирования успешного завершения операций
func (l *MyLogger) Success(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if l.consoleLogger != nil {
		l.logWithCallerInfo(file, line, "SUCCESS", colorBlue, message, args...)
	} else {
		log.Print("No logger available.")
	}
}

// Trace записывает максимально детального логирования,
// например, для отладки сложных операций или отслеживания выполнения каждого шага в коде.
func (l *MyLogger) Trace(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if l.consoleLogger != nil {
		l.logWithCallerInfo(file, line, "TRACE", colorGray, message, args...)
	} else {
		log.Print("No logger available.")
	}
}

// Error записывает сообщение об ошибке в лог вместе с контекстом вызова функции.
// Параметр err содержит ошибку, связанную с данным сообщением.
func (l *MyLogger) Error(message string, err error) {
	_, file, line, _ := runtime.Caller(1)
	if l.consoleLogger != nil {
		l.logWithCallerInfo(file, line, "ERROR", colorRed, message, err)
	} else {
		log.Print("No logger available.")
	}
}

// Warn записывает предупреждение в лог вместе с контекстом вызова функции.
// Параметры args содержат дополнительные данные для сообщения.
func (l *MyLogger) Warn(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if l.consoleLogger != nil {
		l.logWithCallerInfo(file, line, "WARN", colorYellow, message, args...)
	} else {
		log.Print("No logger available.")
	}
}

// Debug записывает информационное сообщение в лог вместе с контекстом вызова функции.
// Параметры args содержат дополнительные данные для сообщения.
func (l *MyLogger) Debug(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if l.consoleLogger != nil {
		l.logWithCallerInfo(file, line, "DEBUG", colorCyan, message, args...)
	} else {
		log.Print("No logger available.")
	}
}

// Fatal записывает фатальное сообщение в лог вместе с контекстом вызова функции
// и завершает приложение с кодом ошибки 1.
// Параметр err содержит ошибку, связанную с данным сообщением.
func (l *MyLogger) Fatal(message string, err error) {
	_, file, line, _ := runtime.Caller(1)
	if l.consoleLogger != nil {
		l.logWithCallerInfo(file, line, "FATAL", colorOrange, "%s: %v", message, err)
		os.Exit(1) // Завершаем приложение с кодом ошибки
	} else {
		log.Print("No logger available.")
	}
}

// Critical записывает логирования критических ошибок,
// которые требуют немедленного внимания, но не приводят к завершению программы (в отличие от Fatal).
func (l *MyLogger) Critical(message string, err error) {
	_, file, line, _ := runtime.Caller(1)
	if l.consoleLogger != nil {
		l.logWithCallerInfo(file, line, "CRITICAL", colorMagenta, message, err)
	} else {
		log.Print("No logger available.")
	}
}

// Panic записывает логирования ситуаций, которые приводят к панике (panic) в программе.
func (l *MyLogger) Panic(message string, err error) {
	_, file, line, _ := runtime.Caller(1)
	if l.consoleLogger != nil {
		l.logWithCallerInfo(file, line, "PANIC", colorWhite, message, err)
		panic(err)
	} else {
		log.Print("No logger available.")
	}
}

// Security записывает логирования событий, связанных с безопасностью (например, попытки несанкционированного доступа).
func (l *MyLogger) Security(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if l.consoleLogger != nil {
		l.logWithCallerInfo(file, line, "SECURITY", colorPurple, message, args...)
	} else {
		log.Print("No logger available.")
	}
}
