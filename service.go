package goat

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"github.com/justwatchcom/goat/wsdl"
)

type Webservice struct {
	services map[string]*wsdl.Definitions
	Client   *http.Client
	header   map[string]interface{}
}

func NewWebservice(c *http.Client, header map[string]interface{}) Webservice {
	if c == nil {
		c = http.DefaultClient
	}

	return Webservice{
		services: map[string]*wsdl.Definitions{},
		Client:   c,
		header:   header,
	}
}

func (self *Webservice) AddServices(urls ...string) (err error) {
	for _, u := range urls {
		var resp *http.Response
		resp, err = http.Get(u)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		s := new(wsdl.Definitions)
		err = xml.NewDecoder(resp.Body).Decode(s)
		if err != nil {
			return
		}

		if s.Service.Name == "" {
			err = fmt.Errorf("invalid service name '%s' for url '%s'", s.Service.Name, u)
			return
		}

		self.services[s.Service.Name] = s
		log.Printf("adding service '%s' from '%s'", s.Service.Name, u)
	}

	return
}
