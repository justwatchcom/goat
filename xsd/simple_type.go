package xsd

import "encoding/xml"

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
