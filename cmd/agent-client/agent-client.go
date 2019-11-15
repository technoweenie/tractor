package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/manifold/qtalk/libmux/mux"
	"github.com/manifold/qtalk/qrpc"
)

var addr = flag.String("addr", ":8081", "qrpc server address")

var commands = map[string]bool{
	"connect": true,
	"start":   true,
	"stop":    true,
}

func main() {
	flag.Parse()
	cmd := strings.ToLower(flag.Arg(0))
	if !commands[cmd] {
		cmd = "connect"
	}

	// connect client to server, call echo
	sess, err := mux.DialWebsocket(*addr)
	if err != nil {
		log.Fatal(err)
	}

	client := &qrpc.Client{Session: sess}
	resp, err := client.Call(cmd, flag.Arg(1), nil)
	if err != nil {
		log.Fatal(err)
	}

	if resp.Hijacked {
		io.Copy(os.Stdout, resp.Channel)
	} else {
		fmt.Println("not hijacked")
	}
}
