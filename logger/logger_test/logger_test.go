package logger_test

import (
	"launchpad/config"
	"launchpad/logger"
	"testing"
)

func init() {
	config.InitConfig("")
}

func TestLogger(t *testing.T) {
	_, err := logger.InitLog(config.Configure.Log)
	if err != nil {
		return
	}
	logger.GetEntry().Debugf("test logger.[num:%d]", 1)
}
