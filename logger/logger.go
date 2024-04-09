package logger

import (
	"fmt"
	"log"
	"os"
)

var logger *Logger

// Logger is a slight improvement over golang's built-in log package
type Logger struct {
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
	FatalLogger *log.Logger
}

func init() {
	logger = &Logger{
		InfoLogger:  log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
		WarnLogger:  log.New(os.Stdout, "[WARN] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
		ErrorLogger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
		FatalLogger: log.New(os.Stderr, "[FATAL] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}
}

func Info(v ...any) {
	logger.InfoLogger.Output(2, fmt.Sprint(v...))
}

func Infof(format string, v ...any) {
	logger.InfoLogger.Output(2, fmt.Sprintf(format, v...))
}

func Warn(v ...any) {
	logger.WarnLogger.Output(2, fmt.Sprint(v...))
}

func Warnf(format string, v ...any) {
	logger.WarnLogger.Output(2, fmt.Sprintf(format, v...))
}

func Error(v ...any) {
	logger.ErrorLogger.Output(2, fmt.Sprint(v...))
}

func Errorf(format string, v ...any) {
	logger.ErrorLogger.Output(2, fmt.Sprintf(format, v...))
}

func Fatal(v ...any) {
	logger.FatalLogger.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...any) {
	logger.FatalLogger.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}
