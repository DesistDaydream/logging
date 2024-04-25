package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	logging "github.com/DesistDaydream/logging/pkg/logrus_init"
)

func main() {
	var logFlags logging.LogrusFlags
	logging.AddFlags(&logFlags)
	logFlags.LogOutput = "test/test_log_dir/test2"
	pflag.Parse()

	if err := logging.LogrusInit(&logFlags); err != nil {
		logrus.Fatal("初始化日志失败", err)
	}

	for {
		// 快速写入 1M 大小的内容到日志文件中
		logrus.Info("开始写入日志文件")
		for i := 0; i < 10000; i++ {
			logrus.Info("这是一条测试日志")
		}

		time.Sleep(2 * time.Second)
	}
}
