package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"qnhd/pkg/file"
	"runtime"
)

type Level int

var (
	F                  *os.File
	DefaultPrefix      = ""
	DefaultCallerDepth = 2
	logger             *log.Logger
	logPrefix          = ""
	levelFlags         = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func Setup() {
	var err error
	filePath := getLogFilePath()
	fileName := getLogFileName()
	F, err = file.MustOpen(fileName, filePath)
	if err != nil {
		log.Fatalf("logging.Setup err: %v", err)
	}

	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}
func Debug(f string, v ...interface{}) {
	setPrefix(DEBUG)
	logger.Printf(f, v...)
}
func Info(f string, v ...interface{}) {
	setPrefix(INFO)
	logger.Printf(f, v...)
}
func Warn(f string, v ...interface{}) {
	setPrefix(WARNING)
	logger.Printf(f, v...)
}
func Error(f string, v ...interface{}) {
	setPrefix(ERROR)
	logger.Printf(f, v...)
}
func Fatal(f string, v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalf(f, v...)
}
func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}
	logger.SetPrefix(logPrefix)
}
