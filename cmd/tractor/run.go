package main

import (
	"github.com/manifold/tractor/pkg/workspace/daemon"
	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Starts the tractor process",
		Long:  "Starts the tractor process",
		Run: func(cmd *cobra.Command, args []string) {
			daemon.Run()
		},
	}
}
