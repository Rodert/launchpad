package cmd

import (
	"github.com/spf13/cobra"
	"launchpad/config"
	"launchpad/logger"
	"launchpad/web"
)

var ConfigPathFlag string

// task
func RunTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "run test",
		Run: func(cmd *cobra.Command, args []string) {
			if err := Init(); err != nil {
				return
			}
		},
	}
}

// api
func RunAPICmd() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "run api",
		Run: func(cmd *cobra.Command, args []string) {
			if err := Init(); err != nil {
				return
			}
			web.Run()
		},
	}
}

func Init() error {
	// init config
	config.InitConfig(ConfigPathFlag)

	// init logger
	_, err := logger.InitLog(config.Configure.Log)
	if err != nil {
		return err
	}
	// db
	//_ = repository.InitPg(config.Configure)

	return err
}

func Init2() error {
	// init config
	config.InitConfig(ConfigPathFlag)

	// init logger
	_, err := logger.InitLog(config.Configure.Log)
	if err != nil {
		return err
	}

	return err
}
