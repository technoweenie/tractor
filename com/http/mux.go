package http

import (
	"net/http"
	"reflect"

	"github.com/manifold/tractor/pkg/manifold"
)

func init() {
	manifold.RegisterComponent(&Mux{}, "")
}

type Mux struct {
	node *manifold.Node `hash:"ignore"`
}

func (c *Mux) InitializeComponent(n *manifold.Node) {
	c.node = n
}

func (c *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	for _, child := range c.node.Children {
		var handler http.Handler
		child.Registry.ValueTo(reflect.ValueOf(&handler))
		if handler != nil {
			mux.Handle("/"+child.Name+"/", handler)
		}
	}
	mux.ServeHTTP(w, r)
}
