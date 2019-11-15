package agent

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/manifold/tractor/pkg/agent/icons"
)

type WorkplaceStatus int

const (
	StatusAvailable = iota
	StatusPartially
	StatusUnavailable
)

func (s WorkplaceStatus) Icon() []byte {
	switch int(s) {
	case 0:
		return icons.Available
	case 1:
		return icons.Partially
	default:
		return icons.Unavailable
	}
}

func (s WorkplaceStatus) String() string {
	switch int(s) {
	case 0:
		return "Available"
	case 1:
		return "Partially"
	default:
		return "Unavailable"
	}
}

type Workspace struct {
	Name      string // base name of dir (~/.tractor/workspaces/{name})
	Path      string
	Socket    string // absolute path to socket file (~/.tractor/sockets/{name}.sock)
	Status    WorkplaceStatus
	buf       *Buffer
	callbacks []func(*Workspace)
	pid       int
	bin       string
	mu        sync.Mutex
}

func NewWorkspace(a *Agent, name string) *Workspace {
	return &Workspace{
		Name:      name,
		Path:      filepath.Join(a.WorkspacesPath, name),
		Socket:    filepath.Join(a.SocketsPath, fmt.Sprintf("%s.sock", name)),
		Status:    StatusPartially,
		bin:       a.bin,
		callbacks: make([]func(*Workspace), 0),
	}
}

func (w *Workspace) Connect() (io.ReadCloser, error) {
	w.mu.Lock()
	if w.buf != nil {
		w.setStatus(StatusAvailable)
		out := w.buf.Pipe()
		w.mu.Unlock()

		return out, nil
	}

	out, err := w.start()
	w.mu.Unlock()
	return out, err
}

// Start starts the workspace daemon. creates the symlink to the path if it does
// not exist, using the path basename as the symlink name
func (w *Workspace) Start() (io.ReadCloser, error) {
	w.mu.Lock()
	out, err := w.start()
	w.mu.Unlock()
	return out, err
}

// must run this when the w.mu mutex is locked
func (w *Workspace) start() (io.ReadCloser, error) {
	buf, err := NewBuffer(1024 * 1024)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(w.bin, "run", "workspace.go",
		"-proto", "unix", "-addr", w.Socket)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Dir = w.Path
	cmd.Stdout = buf
	cmd.Stderr = buf

	if err := cmd.Start(); err != nil {
		w.setStatus(StatusUnavailable)
		return nil, err
	}

	w.buf = buf
	w.pid = cmd.Process.Pid
	w.setStatus(StatusAvailable)

	go func(c *exec.Cmd, ws *Workspace) {
		c.Wait()
		ws.unavailable()
	}(cmd, w)

	return buf.Pipe(), nil
}

// Stop stops the workspace daemon, deleting the unix socket file.
func (w *Workspace) Stop() error {
	if pid := w.unavailable(); pid > 0 {
		return syscall.Kill(-pid, syscall.SIGTERM)
	}
	return nil
}

func (w *Workspace) OnStatusChange(cb func(*Workspace)) {
	cb(w)
	w.mu.Lock()
	w.callbacks = append(w.callbacks, cb)
	w.mu.Unlock()
}

func (w *Workspace) unavailable() int {
	w.mu.Lock()
	if w.buf != nil {
		w.buf.Close()
		w.buf = nil
	}
	pid := w.pid
	w.pid = 0
	w.setStatus(StatusUnavailable)
	w.mu.Unlock()
	return pid
}

// always run when w.mu mutex is locked
func (w *Workspace) setStatus(s WorkplaceStatus) {
	if w.Status == s {
		return
	}

	log.Println("[workspace]", w.Name, "state:", w.Status, "=>", s)
	w.Status = s
	for _, cb := range w.callbacks {
		cb(w)
	}
}
