package agent

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/manifold/qtalk/libmux/mux"
	"github.com/manifold/qtalk/qrpc"
)

func ListenAndServe(a *Agent) error {
	api := qrpc.NewAPI()
	api.HandleFunc("connect", func(r qrpc.Responder, c *qrpc.Call) {
		ws, err := findWorkspace(a, c)
		if err != nil {
			r.Return(err)
			return
		}

		if err := streamWorkspaceOutput(a, r, ws.Connect); err != nil {
			r.Return(err)
			return
		}
	})

	api.HandleFunc("start", func(r qrpc.Responder, c *qrpc.Call) {
		ws, err := findWorkspace(a, c)
		if err != nil {
			r.Return(err)
			return
		}

		if err := streamWorkspaceOutput(a, r, ws.Start); err != nil {
			r.Return(err)
			return
		}
	})

	api.HandleFunc("stop", func(r qrpc.Responder, c *qrpc.Call) {
		ws, err := findWorkspace(a, c)
		if err != nil {
			r.Return(err)
			return
		}
		ws.Stop()

		r.Return(fmt.Sprintf("workspace %q stopped", ws.Name))
	})

	go func() {
		var lastMsg string
		for {
			time.Sleep(time.Second * 3)
			msg, err := wsStatus(a)
			if err != nil {
				log.Println("[workspaces]", err)
			}
			if lastMsg != msg && len(msg) > 0 {
				log.Println("[workspaces]", msg)
			}
			lastMsg = msg
		}
	}()

	server := &qrpc.Server{}
	l, err := mux.ListenUnix(a.AgentSocket)
	if err != nil {
		return err
	}

	log.Println("unix server listening at", a.AgentSocket)
	err = server.Serve(l, api)
	os.Remove(a.AgentSocket)
	return err
}

func wsStatus(a *Agent) (string, error) {
	workspaces, err := a.Workspaces()
	if err != nil || len(workspaces) == 0 {
		return "", err
	}

	pairs := make([]string, len(workspaces))
	for i, ws := range workspaces {
		p, w := ws.BufferStatus()
		pairs[i] = fmt.Sprintf("%s=%s (%d pipe(s), %d written)",
			ws.Name, ws.Status, p, w)
	}
	return strings.Join(pairs, ", "), nil
}

type workspaceFunc func() (io.ReadCloser, error)

func streamWorkspaceOutput(a *Agent, r qrpc.Responder, fn workspaceFunc) error {
	out, err := fn()
	if err != nil {
		return err
	}

	defer out.Close()

	ch, err := r.Hijack("how am i alive?")
	if err != nil {
		return err
	}

	if _, err := io.Copy(ch, out); err != nil {
		if err == io.ErrClosedPipe {
			return ch.Close()
		}
		ch.Close()
		return err
	}
	return ch.Close()
}

func findWorkspace(a *Agent, call *qrpc.Call) (*Workspace, error) {
	var workspacePath string
	if err := call.Decode(&workspacePath); err != nil {
		return nil, err
	}
	log.Println("[qrpc]", call.Destination, workspacePath)

	if ws := a.Workspace(workspacePath); ws != nil {
		return ws, nil
	}

	return nil, fmt.Errorf("no workspace found for %q", workspacePath)
}
