package logger

import (
	"fmt"
	"os"

	"github.com/labstack/gommon/log"
)

var (
	logger      *log.Logger
	logFilePath string
	writer      *LogWriter
)

type LogWriter struct{}

func NewLogWriter() *LogWriter {
	logFilePath = "/usr/hfc"
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		err = os.MkdirAll(logFilePath, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
	}
	logFilePath += "/logs"

	return &LogWriter{}
}

func (w *LogWriter) Write(p []byte) (n int, err error) {
	if w == nil {
		return 0, nil
	}

	logToForward := p
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Error creating or opening file:", err)
		return
	}
	defer file.Close()
	_, err = file.Write(p)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	return os.Stdout.Write(logToForward)
}

func Init() {
	writer = NewLogWriter()

	logger = log.New("hfc")
	logger.SetOutput(writer)
	logger.DisableColor()
	logger.SetLevel(getLogLevel())
}

func Log() *log.Logger {
	return logger
}

func getLogLevel() log.Lvl {
	logLevel := os.Getenv("LOGLEVEL")
	switch logLevel {
	case "DEBUG":
		return log.DEBUG
	case "INFO":
		return log.INFO
	case "WARN":
		return log.WARN
	case "ERROR":
		return log.ERROR
	default:
		logger.Warn("No loglevel found, defaulting to INFO ...")
		return log.INFO
	}
}
