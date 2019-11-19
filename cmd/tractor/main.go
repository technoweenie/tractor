package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tractor",
		Short: "Tractor",
		Long:  "Tractor",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// context that cancels when an os signal to quit the app has been received.
	sigQuit context.Context
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(agentCmd())
	rootCmd.AddCommand(runCmd())

	ct, cancelFunc := context.WithCancel(context.Background())
	sigQuit = ct

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func(c <-chan os.Signal) {
		<-c
		cancelFunc()
	}(c)
}
