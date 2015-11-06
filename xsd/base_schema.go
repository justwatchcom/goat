package xsd

import (
	"encoding/xml"
	"fmt"
	"reflect"
)

type mapping struct {
	xsdSchema []string
	kinds     []reflect.Kind
	format    string
}

// These mappings are used to map between a xsd type which has a base like
// '<restriction base="string"/>'.
var mappings = []mapping{
	{
		xsdSchema: []string{"boolean"},
		kinds:     []reflect.Kind{reflect.Bool},
		format:    "%t",
	},
	{
		xsdSchema: []string{"int", "long"},
		kinds:     []reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr},
		format:    "%d",
	},
	{
		xsdSchema: []string{"float"},
		kinds:     []reflect.Kind{reflect.Float32, reflect.Float64},
		format:    "%f",
	},
	{
		xsdSchema: []string{"string"},
		kinds:     []reflect.Kind{reflect.String},
		format:    "%s",
	},
}

// baseSchema is the Schema implementation of http://www.w3.org/2001/XMLSchema
// (can be found at http://www.w3.org/2001/XMLSchema-datatypes).
// It is not a valid Schema, since it does not tell how to implement the
// simpleType 'string', it just shows string is a simpleType whose base is 'string' itself.
type baseSchema struct{}

// http://www.w3.org/2001/XMLSchema-datatypes does not have elements.
func (baseSchema) EncodeElement(name string, enc *xml.Encoder, sr SchemaRepository, params map[string]interface{}, path ...string) error {
	return fmt.Errorf("not implemented")
}

func (baseSchema) EncodeType(name string, enc *xml.Encoder, sr SchemaRepository, params map[string]interface{}, path ...string) (err error) {
	v, ok := params[MakePath(path)]
	if !ok {
		err = fmt.Errorf("did not find data '%s'", MakePath(path))
		return
	}

	var del bool
	var newVal interface{}
	del, newVal, err = encodeInterfaceType(name, enc, v)
	if err != nil {
		return
	}

	if newVal != nil {
		params[MakePath(path)] = newVal
	}

	if del {
		delete(params, MakePath(path))
	}
	return
}

func encodeInterfaceType(name string, enc *xml.Encoder, v interface{}) (del bool, newVal interface{}, err error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Slice {
		if val.Len() == 0 {
			del = true
			return
		}

		var d bool
		if d, newVal, err = encodeInterfaceType(name, enc, val.Index(0).Interface()); err != nil {
			return
		} else if d {
			newVal = val.Slice(1, val.Len()).Interface()
		}

		if reflect.ValueOf(newVal).Len() == 0 {
			del = true
		}
		return
	}

	for _, m := range mappings {
		for _, n := range m.xsdSchema {
			if n == name {
				for _, t := range m.kinds {
					if t == val.Kind() {
						del = true
						err = enc.EncodeToken(xml.CharData(fmt.Sprintf(m.format, v)))
						return
					}
				}
			}
		}
	}

	err = fmt.Errorf("no mapping found for xsd base type %s and kind %s", name, val.Kind())
	return
}
