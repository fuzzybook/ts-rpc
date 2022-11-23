package tsrpc

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
)

type TSEndpoint struct {
	Name     string
	Path     string
	Method   string
	Request  string
	Response string
	Source   string
	File     string
	Line     int
}

func ParseEndpoint(source string, file string, line int) TSEndpoint {
	p := strings.Index(source, "TSEndpoint=")
	s := strings.Trim(source[p+len("TSEndpoint="):], " ")
	a := strings.Split(s, ";")
	n := 0
	endpoint := TSEndpoint{}
	endpoint.Source = strings.Trim(source, "\t")
	endpoint.File = file
	endpoint.Line = line

	for _, v := range a {
		t := strings.Split(v, "=")
		if len(t) < 2 || strings.Trim(t[1], " ") == "" {
			exitOnError(errors.New(fmt.Sprintf("Worong Endpoint: %s", s)))
		}
		if len(t) == 2 {
			switch strings.Trim(t[0], " ") {
			case "path":
				n++
				endpoint.Path = strings.Trim(t[1], " ")
			case "method":
				n++
				endpoint.Method = strings.Trim(t[1], " ")
			case "name":
				n++
				endpoint.Name = strings.Trim(t[1], " ")
			case "request":
				n++
				endpoint.Request = strings.Trim(t[1], " ")
			case "response":
				n++
				endpoint.Response = strings.Trim(t[1], " ")
			}
		} else {
			exitOnError(errors.New(fmt.Sprintf("Worong Endpoint props: %s", s)))
		}
	}

	if endpoint.Method != "POST" && endpoint.Method != "GET" {
		exitOnError(errors.New(fmt.Sprintf("Worong Endpoint method: %s", s)))
	}

	if endpoint.Method == "GET" && n < 4 {
		exitOnError(errors.New(fmt.Sprintf("Worong Endpoint number of props: %s", s)))
	}
	if endpoint.Method == "POST" && n < 5 {
		exitOnError(errors.New(fmt.Sprintf("Worong Endpoint number of props: %s", s)))
	}

	return endpoint
}

type tplData struct {
	E      *TSEndpoint
	Path   string
	Params []string
}

func (e *TSEndpoint) VerifyTypes(info TSInfo, p string) {
	a := strings.Split(e.Request, ".")
	if e.Request != "" {
		if len(a) == 2 {
			if !info.find(a[0], a[1]) {
				exitOnError(errors.New(fmt.Sprintf("Worong Endpoint request: %s AT %s Line: %d ", e.Request, e.File, e.Line)))
			}
		}
		if len(a) == 1 && !isNativeType(a[0]) {
			if !info.find(p, a[0]) {
				exitOnError(errors.New(fmt.Sprintf("Worong Endpoint request: %s AT %s Line: %d ", e.Request, e.File, e.Line)))
			}
		}
	}
	a = strings.Split(e.Response, ".")
	if len(a) == 2 {
		if !info.find(a[0], a[1]) {
			exitOnError(errors.New(fmt.Sprintf("Worong Endpoint response: %s AT %s Line: %d ", e.Request, e.File, e.Line)))
		}
	}
	if len(a) == 1 && !isNativeType(a[0]) {
		if !info.find(p, a[0]) {
			exitOnError(errors.New(fmt.Sprintf("Worong Endpoint response: %s AT %s Line: %d ", e.Request, e.File, e.Line)))
		}
	}
}

func (e *TSEndpoint) ToTs(pkg string) string {
	data := tplData{E: e, Path: e.Path, Params: []string{}}
	tpl := `
{{ .E.Source }}
// {{ .E.File }} Line: {{ .E.Line }}
{{if eq .E.Method "GET"}}export const {{ .E.Name}} = async ({{range $v := .Params}}{{$v}}{{end}}):Promise<{ data:{{.E.Response }}; error: Nullable<string> }> => {
	return await api.GET({{ .Path}}) as { data: {{ .E.Response}}; error: Nullable<string> };
}{{end}}{{if eq .E.Method "POST"}}export const {{ .E.Name}} = async (data: {{ .E.Request}}):Promise<{ data:{{.E.Response }}; error: Nullable<string> }> => {
	return await api.POST("{{ .Path}}", data) as { data: {{ .E.Response}}; error: Nullable<string> };
}{{end}}`

	if e.Method == "GET" {
		a := strings.Split(e.Path, "/")
		c := ""
		f := false
		for _, v := range a {
			if len(v) == 0 {
				continue
			}
			prefix := v[0:1]

			if f {
				c = ", "
			}

			if prefix == ":" {
				f = true
				data.Params = append(data.Params, fmt.Sprintf("%s%s: string", c, v[1:]))
				data.Path = strings.Replace(data.Path, v, fmt.Sprintf("${%s}", v[1:]), 1)
			} else if prefix == "*" {
				f = true
				data.Params = append(data.Params, fmt.Sprintf("%s%s: Nullable<string>", c, v[1:]))
				data.Path = strings.Replace(data.Path, v, fmt.Sprintf("${%s}", v[1:]), 1)
			}
		}
		if len(data.Params) > 0 {
			data.Path = fmt.Sprintf("`%s`", data.Path)
		} else {
			data.Path = fmt.Sprintf("\"%s\"", data.Path)
		}
	}

	t, err := template.New("test").Parse(tpl)
	if err != nil {
		panic(err)
	}
	var result bytes.Buffer
	err = t.Execute(&result, data)
	if err != nil {
		panic(err)
	}

	return result.String()
}
