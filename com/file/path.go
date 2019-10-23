package file

import (
	"net/http"

	"github.com/manifold/tractor/pkg/manifold"
)

func init() {
	manifold.RegisterComponent(&Path{}, "")
}

type Path struct {
	Filepath string
}

func (c *Path) Open(name string) (http.File, error) {
	return http.Dir(c.Filepath).Open(name)
}
