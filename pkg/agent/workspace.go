package agent

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/armon/circbuf"
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
	Command    *exec.Cmd
	bin        string
	buf        *circbuf.Buffer
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

func (w *Workspace) Bytes() []byte {
	if w == nil || w.buf == nil {
		return nil
	}
	return w.buf.Bytes()
}

// Start starts the workspace daemon. creates the symlink to the path if it does
// not exist, using the path basename as the symlink name
func (w *Workspace) Start() error {
	w.mu.Lock()
	w.Status = StatusPartially
	w.mu.Unlock()

	w.mu.Lock()
	defer w.mu.Unlock()

	time.Sleep(time.Second * 5)

	var err error
	w.buf, err = circbuf.NewBuffer(BufferSize)
	if err != nil {
		w.Status = StatusUnavailable
		return err
	}

	fmt.Println(w.bin, w.Path)
	w.Command = exec.Command(w.bin, "run", "workspace.go")
	w.Command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	w.Command.Dir = w.Path
	w.Command.Stdout = w.buf
	w.Command.Stderr = w.buf

	if err := w.Command.Run(); err != nil {
		w.Status = StatusUnavailable
		return err
	}

	w.Status = StatusAvailable
	return nil
}

// Stop stops the workspace daemon, deleting the unix socket file.
func (w *Workspace) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	return syscall.Kill(-w.Command.Process.Pid, syscall.SIGTERM)
}
