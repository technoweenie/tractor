package main

import (
	"fmt"
	"log"

	"github.com/getlantern/systray"
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

	agent, err := Open("")
	fatal(err)

	workspaces, err := agent.Workspaces()
	fatal(err)
	for _, ws := range workspaces {
		openItem := systray.AddMenuItem(ws.Name, "Open workspace")
		openItem.SetIcon(icons.Available)
		go func() {
			<-openItem.ClickedCh
			open.StartWith(ws.Path, "Visual Studio Code.app")
		}()
	}

	systray.AddSeparator()
	mQuitOrig := systray.AddMenuItem("Shutdown", "Quit and shutdown all workspaces")
	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
	}()

}
