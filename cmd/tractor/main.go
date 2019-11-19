package main

import (
	"log"

	"github.com/manifold/tractor/pkg/workspace/daemon"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tractor",
	Short: "Tractor",
	Long:  "Tractor",
	Run: func(cmd *cobra.Command, args []string) {
		daemon.Run()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
