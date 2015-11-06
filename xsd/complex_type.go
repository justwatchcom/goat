package xsd

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type ComplexType struct {
	XMLName  xml.Name        `xml:"http://www.w3.org/2001/XMLSchema complexType"`
	Name     string          `xml:"name,attr"`
	Abstract bool            `xml:"abstract,attr"`
	Sequence []Element       `xml:"sequence>element"`
	Content  *ComplexContent `xml:"http://www.w3.org/2001/XMLSchema complexContent"`
}

type ComplexContent struct {
	XMLName   xml.Name  `xml:"http://www.w3.org/2001/XMLSchema complexContent"`
	Extension Extension `xml:"http://www.w3.org/2001/XMLSchema extension"`
}

type Extension struct {
	XMLName  xml.Name  `xml:"http://www.w3.org/2001/XMLSchema extension"`
	Base     string    `xml:"base,attr"`
	Sequence []Element `xml:"sequence>element"`
}

func (self *ComplexType) Encode(enc *xml.Encoder, sr SchemaRepository, ga GetAliaser, params map[string]interface{}, path ...string) (err error) {
	for _, e := range self.Sequence {
		err = e.Encode(enc, sr, ga, params, path...)
		if err != nil {
			return
		}
	}

	if self.Content != nil {
		parts := strings.Split(self.Content.Extension.Base, ":")
		switch len(parts) {
		case 2:
			var schema Schemaer
			schema, err = sr.GetSchema(ga.GetAlias(parts[0]))
			if err != nil {
				return
			}

			err = schema.EncodeType(parts[1], enc, sr, params, path...)
			if err != nil {
				return
			}
		default:
			err = fmt.Errorf("malformed base '%s' in path %q", self.Content.Extension.Base, path)
			return
		}

		for _, e := range self.Content.Extension.Sequence {
			err = e.Encode(enc, sr, ga, params, path...)
			if err != nil {
				return
			}
		}
	}

	return
}
