package logging

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Logrus 日志相关命令行标志
type LogrusFlags struct {
	LogLevel  string
	LogOutput string
	LogFormat string
	LogCaller bool
	LogColors bool
}

// LogrusInit 日志功能初始化，若指定了 log-output 命令行标志，则将日志写入到文件中
func LogrusInit(logFlags *LogrusFlags) error {
	logrus.SetReportCaller(logFlags.LogCaller)
	callerPrettyfier := func(frame *runtime.Frame) (function string, file string) {
		fileName := path.Base(frame.File)
		return fmt.Sprintf("%v:%v", frame.Func.Name(), frame.Line), fileName
	}

	switch logFlags.LogFormat {
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableColors:    logFlags.LogColors,
			FullTimestamp:    true,
			TimestampFormat:  "2006-01-02 15:04:05",
			CallerPrettyfier: callerPrettyfier,
		})
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:   "2006-01-02 15:04:05",
			DisableTimestamp:  false,
			DisableHTMLEscape: false,
			DataKey:           "",
			// FieldMap:          map[logrus.fieldKey]string{},
			CallerPrettyfier: callerPrettyfier,
			PrettyPrint:      false,
		})
	}

	logLevel, err := logrus.ParseLevel(logFlags.LogLevel)
	if err != nil {
		return err
	}
	logrus.SetLevel(logLevel)

	if logFlags.LogOutput != "" {
		f, err := os.OpenFile(logFlags.LogOutput, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			return err
		}
		logrus.SetOutput(f)
	}

	return nil
}
