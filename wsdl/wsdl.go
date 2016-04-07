package wsdl

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/justwatchcom/goat/xsd"
)

type InnerDefinitions struct {
	TargetNamespace string    `xml:"targetNamespace,attr"`
	Types           Type      `xml:"types"`
	Messages        []Message `xml:"message"`
	PortType        PortType  `xml:"portType"`
	Binding         []Binding `xml:"binding"`
	Service         Service   `xml:"service"`
}

type Definitions struct {
	XMLName xml.Name `xml:"definitions"`
	Aliases map[string]string
	InnerDefinitions
}

func (self *Definitions) GetAlias(alias string) (space string) {
	return self.Aliases[alias]
}

func (self *Definitions) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	err = d.DecodeElement(&self.InnerDefinitions, &start)
	if err != nil {
		return
	}

	self.XMLName = start.Name
	self.Aliases = map[string]string{}

	self.Types.Schemas = xsd.SchemaMap{}
	for _, schema := range self.Types.Schemata {
		self.Types.Schemas[schema.TargetNamespace] = schema
	}

	for _, attr := range start.Attr {
		if _, ok := self.Aliases[attr.Name.Local]; !ok {
			self.Aliases[attr.Name.Local] = attr.Value
		}

		for k := range self.Types.Schemas {
			if _, ok := self.Types.Schemas[k].Aliases[attr.Name.Local]; !ok {
				self.Types.Schemas[k].Aliases[attr.Name.Local] = attr.Value
			}
		}
	}

	return
}

func copyMap(src map[string]interface{}) map[string]interface{} {
	dst := map[string]interface{}{}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (self *Definitions) WriteRequest(operation string, w io.Writer, headerParams, bodyParams map[string]interface{}) (err error) {
	headerParams = copyMap(headerParams)
	bodyParams = copyMap(bodyParams)

	var bndOp BindingOperation
	var ptOp PortTypeOperation
	bndOp, ptOp, err = self.getOperations(operation)
	if err != nil {
		return
	}

	var header, body xsd.Schema
	var headerElement, bodyElement string
	header, headerElement, err = self.getSchema(bndOp.Input.SoapHeader.PortTypeOperationMessage)
	if err != nil {
		return
	}

	body, bodyElement, err = self.getSchema(bndOp.Input.SoapBody.PortTypeOperationMessage, ptOp.Input)
	if err != nil {
		return
	}

	fmt.Fprint(w, xml.Header)
	enc := xml.NewEncoder(io.MultiWriter(w, os.Stdout))
	enc.Indent("", "  ")
	defer func() {
		if err == nil {
			err = enc.Flush()
		}
	}()

	envelope := xml.StartElement{
		Name: xml.Name{
			Space: "http://schemas.xmlsoap.org/soap/envelope/",
			Local: "Envelope",
		},
	}
	enc.EncodeToken(envelope)
	defer enc.EncodeToken(envelope.End())

	soapHeader := xml.StartElement{
		Name: xml.Name{
			Space: "http://schemas.xmlsoap.org/soap/envelope/",
			Local: "Header",
		},
	}
	enc.EncodeToken(soapHeader)

	err = header.EncodeElement(headerElement, enc, self.Types.Schemas, headerParams)
	if err != nil {
		return
	}
	enc.EncodeToken(soapHeader.End())

	soapBody := xml.StartElement{
		Name: xml.Name{
			Space: "http://schemas.xmlsoap.org/soap/envelope/",
			Local: "Body",
		},
	}
	enc.EncodeToken(soapBody)
	err = body.EncodeElement(bodyElement, enc, self.Types.Schemas, bodyParams)
	if err != nil {
		return
	}
	enc.EncodeToken(soapBody.End())

	return
}

func (self *Definitions) getSchema(msg ...PortTypeOperationMessage) (schema xsd.Schema, element string, err error) {
	for _, s := range msg {
		if s.Message == "" {
			continue
		}

		parts := strings.Split(s.Message, ":")
		if len(parts) == 2 {
			element = parts[1]
			var ok bool
			schema, ok = self.Types.Schemas[self.GetAlias(parts[0])]
			if ok {
				for _, m := range self.Messages {
					if m.Name == element {
						p := strings.Split(m.Part.Element, ":")
						if len(p) != 2 {
							err = fmt.Errorf("invalid message part element name '%s'", m.Part.Element)
							return
						}

						element = p[1]
						return
					}
				}

				err = fmt.Errorf("did not find message '%s'", element)
				return
			}
		} else {
			err = fmt.Errorf("invalid soapheader message format '%s'", s.Message)
		}
	}

	err = fmt.Errorf("did not find schema in %q", msg)
	return
}

func (self *Definitions) getOperations(operation string) (bndOp BindingOperation, ptOp PortTypeOperation, err error) {
	parts := strings.Split(self.Service.Port.Binding, ":")
	switch len(parts) {
	case 2:
		if self.GetAlias(parts[0]) != self.TargetNamespace {
			err = fmt.Errorf("have '%s', want '%s' as target namespace", parts[0], self.TargetNamespace)
			return
		}

		parts[0] = parts[1]
		fallthrough
	case 1:
		for _, bnd := range self.Binding {
			if bnd.Name == parts[0] {
				parts = strings.Split(bnd.Type, ":")
				switch len(parts) {
				case 2:
					if self.GetAlias(parts[0]) != self.TargetNamespace {
						err = fmt.Errorf("have '%s', want '%s' as target namespace in binding '%s'", parts[0], self.TargetNamespace, bnd.Name)
						return
					}

					parts[0] = parts[1]
					fallthrough
				case 1:
					if self.PortType.Name != parts[0] {
						err = fmt.Errorf("have '%s', want '%s' as target namespace in binding '%s'", parts[0], self.PortType.Name, bnd.Name)
						return
					}

					var found bool
					for _, ptOp = range self.PortType.Operations {
						found = ptOp.Name == operation
						if found {
							break
						}
					}

					if !found {
						err = fmt.Errorf("did not find porttype operation '%s' in binding '%s'", operation, bnd.Name)
						return
					}
				default:
					err = fmt.Errorf("malformed binding information '%s' in binding '%s'", bnd.Type, bnd.Name)
					return
				}

				for _, bndOp = range bnd.Operations {
					if bndOp.Name == operation {
						return
					}
				}

				err = fmt.Errorf("did not find operation '%s' in binding '%s'", operation, bnd.Name)
				return
			}
		}

		err = fmt.Errorf("did not find binding '%s'", parts[0])
	default:
		err = fmt.Errorf("malformed binding information: '%s'", self.Service.Port.Binding)
	}

	return
}
