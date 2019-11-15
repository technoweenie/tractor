package agent

import (
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/manifold/tractor/pkg/agent/icons"
)

var BufferSize int64 = 1024 * 1024

type Status int

const (
	StatusAvailable = iota
	StatusPartially
	StatusUnavailable
)

func (s Status) Icon() []byte {
	switch int(s) {
	case 0:
		return icons.Available
	case 1:
		return icons.Partially
	default:
		return icons.Unavailable
	}
}

type Workspace struct {
	Name       string // base name of dir (~/.tractor/workspaces/{name})
	Path       string
	SocketPath string // absolute path to socket file (~/.tractor/sockets/{name}.sock)
	Status     Status
	pid        int
	bin        string
	mu         sync.Mutex
}

func NewWorkspace(a *Agent, name string) *Workspace {
	return &Workspace{
		Name:       name,
		Path:       filepath.Join(a.WorkspacesPath, name),
		SocketPath: filepath.Join(a.SocketsPath, fmt.Sprintf("%s.sock", name)),
		Status:     StatusUnavailable,
		bin:        a.bin,
	}
}

// Start starts the workspace daemon. creates the symlink to the path if it does
// not exist, using the path basename as the symlink name
func (w *Workspace) Start(out io.Writer) error {
	w.mu.Lock()
	w.Status = StatusPartially
	w.mu.Unlock()

	w.mu.Lock()
	defer w.mu.Unlock()

	time.Sleep(time.Second * 5)

	cmd := exec.Command(w.bin, "run", "workspace.go")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Dir = w.Path
	if out != nil {
		cmd.Stdout = out
		cmd.Stderr = out
	}

	if err := cmd.Start(); err != nil {
		w.Status = StatusUnavailable
		return err
	}

	go func(c *exec.Cmd, ws *Workspace) {
		c.Wait()
		ws.unavailable()
	}(cmd, w)

	w.pid = cmd.Process.Pid
	w.Status = StatusAvailable
	return nil
}

// Stop stops the workspace daemon, deleting the unix socket file.
func (w *Workspace) Stop() error {
	return syscall.Kill(-w.unavailable(), syscall.SIGTERM)
}

func (w *Workspace) unavailable() int {
	w.mu.Lock()
	pid := w.pid
	w.pid = 0
	w.Status = StatusUnavailable
	w.mu.Unlock()
	return pid
}
