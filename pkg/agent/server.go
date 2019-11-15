package agent

import (
	"fmt"
	"log"

	"github.com/manifold/qtalk/libmux/mux"
	"github.com/manifold/qtalk/qrpc"
)

func ListenAndServe(a *Agent, addr string) error {
	api := qrpc.NewAPI()
	api.HandleFunc("connect", func(r qrpc.Responder, c *qrpc.Call) {
		ws, err := findWorkspace(a, c)
		if err != nil {
			r.Return(err)
			return
		}

		fmt.Printf("connect: %+v => %+v\n", c, ws)
		r.Return("connect")
	})

	api.HandleFunc("start", func(r qrpc.Responder, c *qrpc.Call) {
		var workspacePath string
		if err := c.Decode(&workspacePath); err != nil {
			r.Return(err)
			return
		}

		r.Return("start")
	})

	api.HandleFunc("stop", func(r qrpc.Responder, c *qrpc.Call) {
		var workspacePath string
		if err := c.Decode(&workspacePath); err != nil {
			r.Return(err)
			return
		}

		r.Return("stop")
	})

	server := &qrpc.Server{}
	l, err := mux.ListenWebsocket(addr)
	if err != nil {
		return err
	}
	log.Println("websocket server listening at", addr)
	return server.Serve(l, api)
}

func findWorkspace(a *Agent, call *qrpc.Call) (*Workspace, error) {
	var workspacePath string
	if err := call.Decode(&workspacePath); err != nil {
		return nil, err
	}

	if ws := a.Workspace(workspacePath); ws != nil {
		return ws, nil
	}

	return nil, fmt.Errorf("no workspace found for %q", workspacePath)
}
