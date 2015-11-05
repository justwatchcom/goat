package xsd

import "encoding/xml"

type ComplexContent struct {
	XMLName   xml.Name  `xml:"http://www.w3.org/2001/XMLSchema complexContent"`
	Extension Extension `xml:"http://www.w3.org/2001/XMLSchema extension"`
}

type Extension struct {
	XMLName  xml.Name  `xml:"http://www.w3.org/2001/XMLSchema extension"`
	Base     string    `xml:"base,attr"`
	Sequence []Element `xml:"sequence>element"`
}
