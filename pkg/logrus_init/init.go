package logging

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// Logrus 日志相关命令行标志
type LogrusFlags struct {
	LogLevel  string
	LogOutput string
	LogFormat string
	LogCaller bool
	LogColors bool

	LogRotateMaxSize  int
	LogRotateBackups  int
	LogRotateMaxAge   int
	LogRotateCompress bool
}

// 添加命令行标志
func AddFlags(flags *LogrusFlags) {
	pflag.StringVar(&flags.LogLevel, "log-level", "info", "日志级别:[debug, info, warn, error, fatal]")
	pflag.StringVar(&flags.LogOutput, "log-output", "", "日志输出位置，不填默认标准输出 stdout")
	pflag.StringVar(&flags.LogFormat, "log-format", "text", "日志输出格式: [text, json]")
	pflag.BoolVar(&flags.LogCaller, "log-caller", false, "是否输出调用者信息，包括函数名、文件名、行号等")
	pflag.BoolVar(&flags.LogColors, "log-colors", false, "是否关闭颜色输出。对于 text 格式，该标志还会影响输出格式")

	pflag.IntVar(&flags.LogRotateMaxSize, "log-rotate-max-size", 5, "日志文件最大 size ，单位为 MiB")
	pflag.IntVar(&flags.LogRotateBackups, "log-rotate-backups", 5, "日志文件最大备份数")
	pflag.IntVar(&flags.LogRotateMaxAge, "log-rotate-max-age", 1, "日志文件最大保留天数")
	pflag.BoolVar(&flags.LogRotateCompress, "log-rotate-compress", true, "是否压缩日志文件，转为 .gz 格式")
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
		if err := os.MkdirAll(filepath.Dir(logFlags.LogOutput), 0755); err != nil {
			return err
		}
		// f, err := os.OpenFile(logFlags.LogOutput, os.O_WRONLY|os.O_CREATE, 0755)
		// if err != nil {
		// 	return err
		// }
		// logrus.SetOutput(f)
		logrus.SetOutput(&lumberjack.Logger{
			Filename:   logFlags.LogOutput,
			MaxSize:    logFlags.LogRotateMaxSize,
			MaxBackups: logFlags.LogRotateBackups,
			MaxAge:     logFlags.LogRotateMaxAge,
			Compress:   logFlags.LogRotateCompress,
		})
	}

	return nil
}
