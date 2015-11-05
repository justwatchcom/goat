package main

var (
	tmplService = `
package {{.PackageName}}

import (
	ws "github.com/pakohan/wsdl-go/webservice"
	"encoding/xml"
)

type Entry struct {
	Key string {{TagDelimiter}}xml:"key"{{TagDelimiter}}
	Value string {{TagDelimiter}}xml:"value"{{TagDelimiter}}
}

type {{.ServiceName}} struct {
	Url string
}

func New{{.ServiceName}}() *{{.ServiceName}}{
	s := {{.ServiceName}}{}
	s.Url = "{{.ServiceUrl}}"

	return &s
}
{{with $s := .}}
{{range .Types}}
type {{.Name}} struct {
	XMLNamespace string {{TagDelimiter}}xml:"xmlns,attr"{{TagDelimiter}}
	{{range .Fields}}{{.Name}} {{if StringHasValue .Type}}{{.Type}}{{end}} {{if StringHasValue .XMLName}}{{TagDelimiter}}xml:"{{.XMLName}}"{{TagDelimiter}}{{end}}
	{{end}}
}
{{end}}
{{range .Messages}}
type {{.Name}} struct {
	XMLName xml.Name        {{TagDelimiter}}xml:"http://webservice.auth.app.bsbr.altec.com/ {{.XMLName}}"{{TagDelimiter}}
	{{if .Input}}Action  string          {{TagDelimiter}}xml:"-"{{TagDelimiter}}{{end}}
	{{range .Params}}
	{{.ParamName}} {{.ParamType}} {{TagDelimiter}}xml:"{{.XMLParamName}}"{{TagDelimiter}}
	{{end}}
}

{{if .Input}}func (si {{.Name}}) GetAction() string {
	return si.Action
}{{end}}
{{end}}{{range .Methods}}
func (s *{{$s.ServiceName}}) {{.Name}}({{if .HasParams}}p{{end}} {{.InputType}}) (r *{{.OutputType}}, err error) {
	{{if .HasParams}}si := {{.MessageIn}}{}
	si.Action = "{{.Action}}"
	si.{{.ParamInName}} = p{{end}}

	sr, err := webservice.CallService({{if .HasParams}}si{{else}}nil{{end}}, s.Url)
	if err != nil {
		return nil, err
	}

	var so {{.MessageOut}}
	err = xml.Unmarshal([]byte(sr.Body.Content), &so)
	if err != nil {
		return nil, err
	}

	return &so.{{.ParamOutName}}, nil
}
{{end}}{{end}}
`
)
