package main

import (
	"fmt"
	"log"

	"github.com/getlantern/systray"
	"github.com/manifold/tractor/pkg/agent"
	"github.com/manifold/tractor/pkg/agent/icons"
	"github.com/skratchdot/open-golang/open"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	ag, err := agent.Open("")
	fatal(err)

	go func(a *agent.Agent) {
		log.Fatal(agent.ListenAndServe(a, ":8081"))
	}(ag)

	systray.Run(onReady(ag), func() {
		fmt.Println("Shutting down...")
		ag.Shutdown()
	})
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
}
