package agent

import (
	"fmt"
	"io"
	"log"
	"os"

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
