package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/manifold/qtalk/libmux/mux"
	"github.com/manifold/qtalk/qrpc"
)

var sock = flag.String("sock", "", "qrpc server unix socket")

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

	socket := *sock
	if len(socket) == 0 {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		socket = filepath.Join(usr.HomeDir, ".tractor", "agent.sock")
	}

	// connect client to server, call echo
	sess, err := mux.DialUnix(socket)
	if err != nil {
		log.Fatal(err)
	}

	client := &qrpc.Client{Session: sess}
	start := time.Now()
	var msg string
	resp, err := client.Call(cmd, flag.Arg(1), &msg)
	if err != nil {
		log.Fatal(err)
	}

	if resp.Hijacked {
		notifySig(func() { resp.Channel.Close() })

		_, err = io.Copy(os.Stdout, resp.Channel)
		resp.Channel.Close()
		if err != nil && err != io.EOF {
			fmt.Println(err)
		}
		fmt.Println()
	} else {
		fmt.Println(msg)
	}

	fmt.Printf("qrpc: %s(%q) %s\n", cmd, flag.Arg(1), time.Since(start))
}

func notifySig(fn func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		fn()
	}()
}
