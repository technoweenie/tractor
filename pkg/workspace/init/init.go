package init

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/gliderlabs/com/objects"
	"github.com/gliderlabs/stdcom/daemon"
	"github.com/gliderlabs/stdcom/log/std"
	"github.com/manifold/tractor/pkg/manifold"
	frontend "github.com/manifold/tractor/pkg/session"
	"github.com/manifold/tractor/pkg/workspace"

	_ "github.com/manifold/tractor/com/file"
	_ "github.com/manifold/tractor/com/http"
	_ "github.com/manifold/tractor/com/net"
)

var (
	addr  = flag.String("addr", "localhost:4243", "server listener address")
	proto = flag.String("proto", "websocket", "server listener protocol")
)

type PreInitializer interface {
	PreInitialize()
}

type Initializer interface {
	Initialize()
}

func init() {
	flag.Parse()

	var err error
	manifold.Root, err = workspace.LoadHierarchy()
	if err != nil {
		panic(err)
	}

	registry := &objects.Registry{}

	manifold.Walk(manifold.Root, func(n *manifold.Node) {
		for _, com := range n.Components {
			registry.Register(objects.New(com.Ref, ""))
			initializer, ok := com.Ref.(PreInitializer)
			if ok {
				initializer.PreInitialize()
			}
		}
	})
	std.Register(registry)
	registry.Reload()

	manifold.Walk(manifold.Root, func(n *manifold.Node) {
		for _, com := range n.Components {
			initializer, ok := com.Ref.(Initializer)
			if ok {
				initializer.Initialize()
			}
		}
	})

	manifold.Root.Observe(&manifold.NodeObserver{
		OnChange: func(node *manifold.Node, path string, old, new interface{}) {
			if path == "Name" && node.Dir != "" {
				newDir := filepath.Join(filepath.Dir(node.Dir), new.(string))
				if node.Dir != newDir {
					// TODO: do not break abstraction, have workspace handle this
					if err := os.Rename(node.Dir, newDir); err != nil {
						log.Fatal(err)
					}
				}
			}
			err := workspace.SaveHierarchy(manifold.Root)
			if err != nil {
				panic(err)
				//log.Println(err)
			}
		},
	})

	go func() {
		daemon.Run(registry, "app")
	}()

	frontend.ListenAndServe(manifold.Root, *proto, *addr)
}
