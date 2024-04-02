package config_test

import (
	"launchpad/config"
	"testing"
)

func TestInitConfig(t *testing.T) {
	config.InitConfig("conf.json")
}
