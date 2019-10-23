package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"path"

	"github.com/getlantern/systray"
	"github.com/manifold/tractor/agent/icons"
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

	usr, err := user.Current()
	fatal(err)
	workspacesPath := path.Join(usr.HomeDir, ".tractor", "workspaces")
	files, err := ioutil.ReadDir(workspacesPath)
	fatal(err)
	for _, file := range files {
		if !file.IsDir() {
			openItem := systray.AddMenuItem(file.Name(), "Open workspace")
			openItem.SetIcon(icons.Available)
			go func() {
				<-openItem.ClickedCh
				open.StartWith(path.Join(workspacesPath, file.Name()), "Visual Studio Code.app")
			}()
		}
	}

	systray.AddSeparator()
	mQuitOrig := systray.AddMenuItem("Shutdown", "Quit and shutdown all workspaces")
	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
	}()

}
