package logging

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func TestLogrusInit(t *testing.T) {
	var logFlags LogrusFlags
	AddFlags(&logFlags)
	logFlags.LogOutput = "test/test_log_dir/test2"
	if err := LogrusInit(&logFlags); err != nil {
		logrus.Fatal("初始化日志失败", err)
	}
	pflag.Parse()
}
