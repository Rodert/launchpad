package pgtest_test

import (
	"launchpad/config"
	repository "launchpad/repository/pg"
	"testing"
)

func init() {
	config.InitConfig("")
}

func TestInitDB(t *testing.T) {
	repository.InitPg(config.Configure)
}
