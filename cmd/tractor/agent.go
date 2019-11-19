package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/getlantern/systray"
	"github.com/manifold/tractor/pkg/agent"
	"github.com/manifold/tractor/pkg/agent/icons"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Starts the agent systray app",
	Long:  "Starts the agent systray app.",
	Run:   runAgent,
}

func runAgent(cmd *cobra.Command, args []string) {
	ag, err := agent.Open("")
	fatal(err)

	go func(a *agent.Agent) {
		fatal(agent.ListenAndServe(a))
	}(ag)

	systray.Run(onReady(ag), ag.Shutdown)
}

func onReady(ag *agent.Agent) func() {
	return func() { buildSystray(ag) }
}

func buildSystray(ag *agent.Agent) {
	systray.SetIcon(icons.Tractor)
	systray.SetTooltip("Tractor")

	workspaces, err := ag.Workspaces()
	fatal(err)

	for _, ws := range workspaces {
		openItem := systray.AddMenuItem(ws.Name, "Open workspace")

		ws.OnStatusChange(func(ws *agent.Workspace) {
			openItem.SetIcon(ws.Status.Icon())
		})

		go func(mi *systray.MenuItem, ws *agent.Workspace) {
			for {
				<-openItem.ClickedCh
				open.StartWith(ws.Path, "Visual Studio Code.app")
			}
		}(openItem, ws)
	}

	systray.AddSeparator()
	mQuitOrig := systray.AddMenuItem("Shutdown", "Quit and shutdown all workspaces")
	go func(mi *systray.MenuItem) {
		<-mi.ClickedCh
		systray.Quit()
	}(mQuitOrig)

	notifySig()
}

func notifySig() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		systray.Quit()
	}()
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
