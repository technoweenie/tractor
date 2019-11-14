package agent

import (
	"fmt"
	"path/filepath"

	"github.com/manifold/tractor/pkg/agent/icons"
)

type Status struct {
	Name string
	Icon []byte
}

var (
	StatusAvailable   = Status{Name: "Available", Icon: icons.Available}
	StatusPartially   = Status{Name: "Partially", Icon: icons.Partially}
	StatusUnavailable = Status{Name: "Unavailable", Icon: icons.Unavailable}
)

type Workspace struct {
	Name       string
	Path       string
	SocketPath string
	Status     Status
}

func NewWorkspace(a *Agent, name string) *Workspace {
	return &Workspace{
		Name:       name,
		Path:       filepath.Join(a.WorkspacesPath, name),
		SocketPath: filepath.Join(a.SocketsPath, fmt.Sprintf("%s.sock", name)),
		Status:     StatusUnavailable,
	}
}
