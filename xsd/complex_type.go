package xsd

import "encoding/xml"

type ComplexType struct {
	XMLName  xml.Name        `xml:"http://www.w3.org/2001/XMLSchema complexType"`
	Name     string          `xml:"name,attr"`
	Abstract bool            `xml:"abstract,attr"`
	Sequence []Element       `xml:"sequence>element"`
	Content  *ComplexContent `xml:"http://www.w3.org/2001/XMLSchema complexContent"`
}

func (self *ComplexType) Encode(enc *xml.Encoder, sr SchemaRepository, ga GetAliaser, params map[string]interface{}, path ...string) (err error) {
	for _, e := range self.Sequence {
		err = e.Encode(enc, sr, ga, params, path...)
		if err != nil {
			return
		}
	}

	return
}
