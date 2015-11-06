package xsd

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type SimpleType struct {
	XMLName     xml.Name              `xml:"http://www.w3.org/2001/XMLSchema simpleType"`
	Name        string                `xml:"name,attr"`
	Restriction SimpleTypeRestriction `xml:"restriction"`
}

type SimpleTypeRestriction struct {
	XMLName      xml.Name      `xml:"http://www.w3.org/2001/XMLSchema restriction"`
	Base         string        `xml:"base,attr"`
	Enumerations []Enumeration `xml:"http://www.w3.org/2001/XMLSchema enumeration"`
}

type Enumeration struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema enumeration"`
	Value   string   `xml:"value,attr"`
}

func (self *SimpleType) Encode(enc *xml.Encoder, sr SchemaRepository, ga GetAliaser, params map[string]interface{}, path ...string) (err error) {
	name := self.Restriction.Base
	parts := strings.Split(name, ":")
	switch len(parts) {
	case 2:
		var schema Schemaer
		schema, err = sr.GetSchema(ga.GetAlias(parts[0]))
		if err != nil {
			return
		}

		err = schema.EncodeType(parts[1], enc, sr, params, path...)
	default:
		err = fmt.Errorf("invalid restriction format '%s'", name)
	}

	return
}
