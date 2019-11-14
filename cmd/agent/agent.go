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
	systray.Run(onReady, func() {
		fmt.Println("Shutting down...")
	})
}

func onReady() {
	systray.SetIcon(icons.Tractor)
	systray.SetTooltip("Tractor")

	ag, err := agent.Open("")
	fatal(err)

	workspaces, err := ag.Workspaces()
	fatal(err)
	for _, ws := range workspaces {
		openItem := systray.AddMenuItem(ws.Name, "Open workspace")
		openItem.SetIcon(icons.Available)
		go func(mi *systray.MenuItem, ws *agent.Workspace) {
			<-openItem.ClickedCh
			open.StartWith(ws.Path, "Visual Studio Code.app")
		}(openItem, ws)
	}

	systray.AddSeparator()
	mQuitOrig := systray.AddMenuItem("Shutdown", "Quit and shutdown all workspaces")
	go func(mi *systray.MenuItem) {
		<-mi.ClickedCh
		systray.Quit()
	}(mQuitOrig)

}
