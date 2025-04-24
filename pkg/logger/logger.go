package logger

import (
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerConfig struct {
	FilePath string `json:"filePath"`
}

var Logger *log.Logger

func Init(workDir string, loggerConfig *LoggerConfig) error {
	Logger = log.New()

	Logger.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	logFilePath := workDir + loggerConfig.FilePath
	dir := filepath.Dir(logFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Println("Error creating directory:", err)
			return err
		}
	} else if err != nil {
		log.Println("Error checking directory existence:", err)
		return err
	}
	lumberjack := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxBackups: 14,
		LocalTime:  true,
	}
	Logger.SetOutput(io.MultiWriter(os.Stdout, lumberjack))

	Logger.SetLevel(log.DebugLevel)
	Logger.SetReportCaller(true)
	return nil
}
