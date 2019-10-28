package http

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"sync"

	"github.com/manifold/tractor/pkg/manifold"
	frontend "github.com/manifold/tractor/pkg/session"
)

func init() {
	manifold.RegisterComponent(&APICaller{}, "")
}

type APICaller struct {
	Method       string
	URL          string
	header       map[string]string
	query        map[string]string
	lastResponse *http.Response `hash:"ignore"`

	mu   sync.Mutex     `hash:"ignore"`
	node *manifold.Node `hash:"ignore"`
}

func (c *APICaller) InitializeComponent(n *manifold.Node) {
	log.Println("api init component")
	c.node = n
}

func (c *APICaller) InspectorButtons() []frontend.Button {
	log.Println("api inspector buttons")
	return []frontend.Button{{
		Name: "Call",
	}}
}

func (c *APICaller) Call() {
	log.Println("api call")
	if req := c.buildReq(); req != nil {
		c.DoJSON(req)
	}
}

func (c *APICaller) buildReq() *http.Request {
	log.Println("api buildReq")
	log.Printf("HTTP %q %q", c.Method, c.URL)
	req, err := http.NewRequest(c.Method, c.URL, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	c.SetRequestQuery(req, c.query)
	c.SetRequestHeader(req, c.header)

	return req
}

func (c *APICaller) SetRequestQuery(req *http.Request, values map[string]string) {
	if len(values) == 0 {
		return
	}

	q := req.URL.Query()
	for k, v := range values {
		log.Printf("  ?%q=%q", k, v)
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()
}

func (c *APICaller) SetRequestHeader(req *http.Request, values map[string]string) {
	for k, v := range values {
		log.Printf("  %q=%q", k, v)
		req.Header.Set(k, v)
	}
}

func (c *APICaller) DoJSON(req *http.Request) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}

	c.saveResp(res)

	data := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&data)
	res.Body.Close()
	if err != nil {
		log.Print(err)
		return
	}

	log.Printf("JSON: %+v", data)

	for _, child := range c.node.Children {
		var jh JSONHandler
		child.Registry.ValueTo(reflect.ValueOf(&jh))
		if jh != nil {
			jh.OnJSON(data)
		}
	}
}

func (c *APICaller) saveResp(res *http.Response) {
	log.Println("api saveResp")
	log.Printf("HTTP %d", res.StatusCode)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastResponse = res
}

type JSONHandler interface {
	OnJSON(map[string]interface{})
}
