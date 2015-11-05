package xsd

import (
	"encoding/xml"
	"strings"
)

type Schemaer interface {
	EncodeElement(name string, enc *xml.Encoder, sr SchemaRepository, params map[string]interface{}, path ...string) (err error)
	EncodeType(name string, enc *xml.Encoder, sr SchemaRepository, params map[string]interface{}, path ...string) (err error)
}

type GetAliaser interface {
	GetAlias(string) string
	Namespace() string
}

type SchemaRepository interface {
	GetSchema(space string) (Schemaer, error)
}

func hasPrefix(m map[string]interface{}, prefix string) (ok bool) {
	for k := range m {
		ok = strings.HasPrefix(k, prefix)
		if ok {
			return
		}
	}

	return
}

func MakePath(path []string) string {
	return strings.Join(path, "/")
}
