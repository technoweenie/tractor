package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/manifold/tractor/pkg/manifold"
)

func init() {
	manifold.RegisterComponent(&WeatherUnlockedAPI{
		APICaller: APICaller{},
	}, "")
}

type WeatherUnlockedAPI struct {
	Latitude  string
	Longitude string
	AppID     string
	AppKey    string
	APICaller
}

func (c *WeatherUnlockedAPI) Call() {
	log.Println("weather call")
	if req := c.buildReq(); req != nil {
		c.DoJSON(req)
	}
}

func (c *WeatherUnlockedAPI) buildReq() *http.Request {
	log.Println("weather buildReq")
	rawurl := fmt.Sprintf(
		"http://api.weatherunlocked.com/api/current/%s,%s", c.Latitude, c.Longitude)
	req, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	c.SetRequestQuery(req, map[string]string{
		"app_id":  c.AppID,
		"app_key": c.AppKey,
	})

	log.Printf("HTTP %q %q", req.Method, req.URL)
	return req
}
