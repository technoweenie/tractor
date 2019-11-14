package agent

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type Agent struct {
	Path           string // ~/.tractor
	AgentPath      string // ~/.tractor/agent.sock
	WorkspacesPath string // ~/.tractor/workspaces
	SocketsPath    string // ~/.tractor/sockets
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

	return &Agent{
		Path:           path,
		AgentPath:      filepath.Join(path, "agent.sock"),
		WorkspacesPath: filepath.Join(path, "workspaces"),
		SocketsPath:    filepath.Join(path, "sockets"),
		workspaces:     make(map[string]*Workspace),
	}, nil
}

func (a *Agent) Workspaces() ([]*Workspace, error) {
	names, err := a.workspaceNames()
	if err != nil {
		return nil, err
	}

	workspaces := make([]*Workspace, len(names))
	for i, n := range names {
		ws := a.workspaces[n]
		if ws == nil {
			ws = NewWorkspace(a, n)
			a.workspaces[n] = ws
		}
		workspaces[i] = ws
	}
	return workspaces, nil
}

func (a *Agent) workspaceNames() ([]string, error) {
	entries, err := ioutil.ReadDir(a.WorkspacesPath)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !a.entryIsWorkspace(entry) {
			continue
		}
		names = append(names, entry.Name())
	}

	return names, nil
}

func (a *Agent) entryIsWorkspace(fi os.FileInfo) bool {
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
