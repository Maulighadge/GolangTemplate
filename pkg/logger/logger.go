package logger

import (
	"log"
	"os"
)

type Logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

var Log *Logger

func InitLogger() {
	Log = &Logger{
		infoLogger:  log.New(os.Stdout, "\033[34m[INFO]\033[0m ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "\033[33m[WARN]\033[0m ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "\033[31m[ERROR]\033[0m ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(os.Stdout, "\033[36m[DEBUG]\033[0m ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func Info(format string, v ...interface{}) {
	if Log != nil && Log.infoLogger != nil {
		Log.infoLogger.Printf(format, v...)
	} else {
		log.Printf("[INFO] "+format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if Log != nil && Log.warnLogger != nil {
		Log.warnLogger.Printf(format, v...)
	} else {
		log.Printf("[WARN] "+format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if Log != nil && Log.errorLogger != nil {
		Log.errorLogger.Printf(format, v...)
	} else {
		log.Printf("[ERROR] "+format, v...)
	}
}

func Debug(format string, v ...interface{}) {
	if Log != nil && Log.debugLogger != nil {
		Log.debugLogger.Printf(format, v...)
	} else {
		log.Printf("[DEBUG] "+format, v...)
	}
}
