package goat

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ResponseEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Header  struct {
		XMLName xml.Name `xml:"Header"`
		Data    []byte   `xml:",innerxml"`
	}
	Body struct {
		XMLName xml.Name `xml:"Body"`
		Data    []byte   `xml:",innerxml"`
	}
}

func (self *Webservice) Do(service, method string, res interface{}, params map[string]interface{}) (err error) {
	s := self.services[service]
	if s == nil {
		err = fmt.Errorf("no such service '%s'", service)
		return
	}

	buf := new(bytes.Buffer)
	err = s.WriteRequest(method, buf, self.header, params)
	if err != nil {
		return
	}

	log.Println(buf.String())

	var resp *http.Response
	resp, err = self.Client.Post(s.Service.Port.Address.Location, "application/soap+xml", io.TeeReader(buf, os.Stdout))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var b []byte
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		err = errors.New(string(b))
		return
	}

	e := new(ResponseEnvelope)
	err = xml.NewDecoder(resp.Body).Decode(e)
	if err != nil {
		return
	}

	err = xml.Unmarshal(e.Body.Data, res)
	return
}
