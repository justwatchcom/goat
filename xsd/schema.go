package xsd

import (
	"encoding/xml"
	"fmt"
)

type InnerSchema struct {
	TargetNamespace    string        `xml:"targetNamespace,attr"`
	ElementFormDefault string        `xml:"elementFormDefault,attr"`
	Version            string        `xml:"version,attr"`
	ComplexTypes       []ComplexType `xml:"http://www.w3.org/2001/XMLSchema complexType"`
	SimpleTypes        []SimpleType  `xml:"http://www.w3.org/2001/XMLSchema simpleType"`
	Elements           []Element     `xml:"http://www.w3.org/2001/XMLSchema element"`
}

type Schema struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema schema"`
	Aliases map[string]string
	InnerSchema
}

func (self *Schema) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	err = d.DecodeElement(&self.InnerSchema, &start)
	if err != nil {
		return
	}

	self.XMLName = start.Name
	self.Aliases = map[string]string{}

	for _, attr := range start.Attr {
		self.Aliases[attr.Name.Local] = attr.Value
	}
	return
}

func (self *Schema) Namespace() string {
	return self.TargetNamespace
}

func (self *Schema) GetAlias(alias string) (space string) {
	return self.Aliases[alias]
}

func (self *Schema) EncodeElement(name string, enc *xml.Encoder, sr SchemaRepository, params map[string]interface{}, path ...string) error {
	for _, elem := range self.Elements {
		if elem.Name == name {
			return elem.Encode(enc, sr, self, params, path...)
		}
	}

	return fmt.Errorf("did not find element '%s'", name)
}

func (self *Schema) EncodeType(name string, enc *xml.Encoder, sr SchemaRepository, params map[string]interface{}, path ...string) error {
	for _, cmplx := range self.ComplexTypes {
		if cmplx.Name == name {
			return cmplx.Encode(enc, sr, self, params, path...)
		}
	}

	for _, smpl := range self.SimpleTypes {
		if smpl.Name == name {
			return smpl.Encode(enc, sr, self, params, path...)
		}
	}

	return fmt.Errorf("did not find type '%s'", name)
}
