package logger

import (
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"os"
)

var log = logging.MustGetLogger("logger")

func init() {
	stdoutFormat := logging.MustStringFormatter(
		`%{color}%{time:2006-01-02T15:04:05.999Z-07:00} %{module} %{shortfunc} ▶ %{level:.4s} %{color:reset} %{message}`,
	)

	stdoutBackend := logging.NewLogBackend(os.Stdout, "", 0)
	stdoutFormatted := logging.NewBackendFormatter(stdoutBackend, stdoutFormat)
	stdoutLeveled := logging.AddModuleLevel(stdoutFormatted)
	stdoutLeveled.SetLevel(logging.WARNING, viper.GetString("logger.level"))
	logging.SetBackend(stdoutLeveled)

	fileFormat := logging.MustStringFormatter(
		`%{time:2006-01-02T15:04:05.999Z-07:00} %{module} %{shortfunc} - %{level:.4s} %{message}`,
	)

	logFile, err := os.OpenFile(viper.GetString("logger.file"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Warning("Cannot open log file:", err)
		log.Warning("All logs will be printed to stdout.")
		return
	}
	fileBackend := logging.NewLogBackend(logFile, "", 0)
	fileFormatted := logging.NewBackendFormatter(fileBackend, fileFormat)
	fileLeveled := logging.AddModuleLevel(fileFormatted)
	fileLeveled.SetLevel(logging.WARNING, viper.GetString("logger.level"))

	log.Info("Writing logs to: ", viper.GetString("logger.file"))

	logging.SetBackend(stdoutLeveled, fileLeveled)
}
