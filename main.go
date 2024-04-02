package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"launchpad/cmd"
	"launchpad/internal/version"
)

// launchpad
func main() {
	root := &cobra.Command{Use: "launchpad", Version: version.VersionInfo()}
	if err := cmd.Command(root).Execute(); err != nil {
		logrus.Error(err)
	}
}
