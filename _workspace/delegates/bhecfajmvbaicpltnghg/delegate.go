package node

import (
	"io"

	"github.com/gliderlabs/ssh"
	"github.com/manifold/tractor/server/manifold"
)

func init() {
	manifold.RegisterDelegate(&Delegate{}, "bhecfajmvbaicpltnghg")
}

type Delegate struct {
	Message string
}

func (d *Delegate) HandleSSH(sess ssh.Session) {
	io.WriteString(sess, d.Message+"\n")
}
