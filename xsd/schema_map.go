package xsd

import (
	"fmt"
	"log"
)

type SchemaMap map[string]Schema

func (self SchemaMap) GetSchema(space string) (s Schemaer, err error) {
	switch space {
	case "http://www.w3.org/2001/XMLSchema":
		s = baseSchema{}
	default:
		if ss, ok := self[space]; !ok {
			err = fmt.Errorf("namespace not found: '%s'", space)
		} else {
			s = &ss
		}
	}

	return
}

func (self SchemaMap) GetElement(space, name string) *Element {
	schema, ok := self[space]
	if !ok {
		log.Printf("namespace not found: '%s'", space)
		return nil
	}

	for _, elem := range schema.Elements {
		if elem.Name == name {
			return &elem
		}
	}

	return nil
}
