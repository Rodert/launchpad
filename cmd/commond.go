package cmd

import (
	"github.com/spf13/cobra"
)

func Command(root *cobra.Command) *cobra.Command {
	// run
	root.AddCommand(
		RunTestCmd(),
	)

	root.AddCommand(
		RunAPICmd(),
	)

	// conf
	flagSet := root.PersistentFlags()
	{
		flagSet.StringVarP(&ConfigPathFlag, "configure", "c", "./conf.json", "configure file path")
	}

	return root
}
