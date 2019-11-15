package agent

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

type Agent struct {
	Path           string // ~/.tractor
	AgentPath      string // ~/.tractor/agent.sock
	WorkspacesPath string // ~/.tractor/workspaces
	SocketsPath    string // ~/.tractor/sockets
	bin            string
	workspaces     map[string]*Workspace
}

func Open(path string) (*Agent, error) {
	if len(path) == 0 {
		p, err := defaultPath()
		if err != nil {
			return nil, err
		}
		path = p
	}

	bin, err := exec.LookPath("go")
	if err != nil {
		return nil, err
	}

	return &Agent{
		Path:           path,
		AgentPath:      filepath.Join(path, "agent.sock"),
		WorkspacesPath: filepath.Join(path, "workspaces"),
		SocketsPath:    filepath.Join(path, "sockets"),
		bin:            bin,
		workspaces:     make(map[string]*Workspace),
	}, nil
}

func (a *Agent) Workspace(path string) *Workspace {
	return a.workspaces[path]
}

func (a *Agent) Shutdown() {
	fmt.Println("shutdown")
	for _, ws := range a.workspaces {
		fmt.Println("shutting down", ws.Name)
		ws.Stop()
	}
}

func (a *Agent) Workspaces() ([]*Workspace, error) {
	entries, err := ioutil.ReadDir(a.WorkspacesPath)
	if err != nil {
		return nil, err
	}

	workspaces := make([]*Workspace, 0, len(entries))
	for _, entry := range entries {
		if !a.isWorkspacePath(entry) {
			continue
		}

		n := entry.Name()
		ws := a.workspaces[n]
		if ws == nil {
			ws = NewWorkspace(a, n)
			a.workspaces[n] = ws
		}
		workspaces = append(workspaces, ws)
	}
	return workspaces, nil
}

func (a *Agent) isWorkspacePath(fi os.FileInfo) bool {
	if fi.IsDir() {
		return true
	}

	path := filepath.Join(a.WorkspacesPath, fi.Name())
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		log.Println(err)
		return false
	}

	if resolved == path {
		return false
	}

	rfi, err := os.Lstat(resolved)
	if err != nil {
		log.Println(err)
		return false
	}

	return rfi.IsDir()
}
func defaultPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, ".tractor"), nil
}
