// exportable typescript generated from golang
// Copyright (C) 2022  Fabio Prada

package tsrpc

import (
	"errors"
	"fmt"
	"strings"
)

type TSModule struct {
	Structs   map[string]string
	Types     map[string]string
	Enums     map[string]string
	Consts    map[string]string
	GTypes    map[string]string
	Endpoints map[string]string
}

type TSSouces struct {
	Pakages map[string]TSModule
	Errors  []string
}

func (ts *TSSouces) findStruct(p string, n string) bool {
	if _, ok := ts.Pakages[p]; ok {
		if _, ok := ts.Pakages[p].Structs[n]; ok {
			return true
		}
	}
	return false
}

func (ts *TSSouces) findType(p string, n string) bool {
	if _, ok := ts.Pakages[p]; ok {
		if _, ok := ts.Pakages[p].Types[n]; ok {
			return true
		}
	}
	return false
}

func (ts *TSSouces) find(p string, n string) bool {
	return ts.findType(p, n) || ts.findStruct(p, n)
}

func emptySrtuct(info TSInfo, p string, k string) bool {
	s := info.Packages[p].structs[k]
	fields := 0
	for _, v := range s.Fields {
		if !v.Json.Ignore {
			fields++
		}
	}
	return fields == 0
}

func structToTs(info TSInfo, p string, k string) (string, []string, error) {
	var result = ""
	var dependencies = []string{}
	fields := 0
	s := info.Packages[p].structs[k]
	if s.Name == "" {
		return "", []string{}, errors.New("Name not found")
	}
	if len(s.Fields) > 0 {
		result += fmt.Sprintf("\nexport interface %s {\n", k)

		for _, v := range info.Packages[p].structs[k].Fields {

			if v.Json.Ignore && !v.Ts.Expand {
				continue
			}
			var typeName = v.TsType
			if v.Ts.Type != "" {
				typeName = v.Ts.Type
			}
			if v.Ts.Expand {
				sp := strings.Split(v.Name, ".")
				pkg := p
				name := v.Name
				if len(sp) > 1 {
					pkg = sp[0]
					name = sp[1]
				}
				for _, v := range info.Packages[pkg].structs[name].Fields {
					result += fmt.Sprintf("\t%s: %s;\n", v.Name, v.TsType)
					fields++
				}
			} else {
				omitempty := ""
				if v.Json.OmitEmpty {
					omitempty = "?"
				}
				if !v.Json.Ignore {
					result += fmt.Sprintf("\t%s%s: %s;\n", v.Json.Name, omitempty, typeName)
					fields++
				}
			}
			if v.DependOn && v.Ts.Type == "" {
				dependencies = append(dependencies, v.TsType)
			}

		}
		result += fmt.Sprintf("}\n")
	}
	if fields == 0 {
		return result, dependencies, errors.New(fmt.Sprintf("Struct %s not export fields  AT: %s", k, info.Packages[p].structs[k].SourceInfo))
	}
	return result, dependencies, nil
}

func typeToTs(info TSInfo, p string, k string) (string, []string) {
	var result = ""
	// var dependencies = []string{}
	s := info.Packages[p].types[k]
	result = fmt.Sprintf("export type %s = %s\n", s.Name, s.TsType)
	return result, []string{}
}

func enumToTs(info TSInfo, p string, k string) string {
	var result = ""
	s := info.Packages[p].enums[k]
	result += fmt.Sprintf("export const Enum%s = {\n", s.Name)
	for _, v := range s.Info {
		result += fmt.Sprintf("%s: %s,\n", v.Key, v.Value)
	}
	result += fmt.Sprintf("} as const\n")
	return result
}

func constToTs(info TSInfo, p string, k string) string {
	var result = ""
	s := info.Packages[p].consts[k]
	result += fmt.Sprintf("export const %s = %s\n", s.Name, s.Value)
	return result
}

func (ts *TSSouces) AddDependencies(info TSInfo, p string, s string, dependencies []string) {
	if len(dependencies) > 0 {
		for _, v := range dependencies {
			pk := p
			st := string(v)
			if strings.Contains(st, ".") {
				sp := strings.Split(st, ".")
				if len(sp) == 2 {
					pk = sp[0]
					st = sp[1]
				}
			}

			if info.findStruct(pk, st) {
				if emptySrtuct(info, pk, st) {
					ts.Errors = append(ts.Errors, fmt.Sprintf("Empty struct %s.%s AT: %s", pk, st, info.Packages[p].structs[s].SourceInfo))
				}
				s, d, err := structToTs(info, pk, st)
				if err != nil {
					ts.Errors = append(ts.Errors, err.Error())
				}
				if _, ok := ts.Pakages[pk]; !ok {
					ts.Pakages[pk] = TSModule{Structs: make(map[string]string), Types: make(map[string]string)}
				}
				ts.Pakages[pk].Structs[st] = s
				if len(d) > 0 {
					ts.AddDependencies(info, pk, st, d)
				}
			} else if info.findType(pk, st) {
				s, _ := typeToTs(info, pk, st)
				if _, ok := ts.Pakages[pk]; !ok {
					ts.Pakages[pk] = TSModule{Structs: make(map[string]string), Types: make(map[string]string)}
				}
				ts.Pakages[pk].Types[st] = s
			} else {
				ts.Errors = append(ts.Errors, fmt.Sprintf("Dipendence not found %s.%s AT: %s", pk, st, info.Packages[p].structs[s].SourceInfo))
			}
		}
	}

}

func (ts *TSSouces) Populate(info TSInfo) {
	ts.Pakages = make(map[string]TSModule)
	ts.Errors = []string{}
	for p, _ := range info.Packages {

		for _, st := range info.Packages[p].structs {
			if st.Typescript {
				if len(st.Fields) == 0 {
					ts.Errors = append(ts.Errors, fmt.Sprintf("Empty struct %s.%s AT: %s", p, st.Name, info.Packages[p].structs[st.Name].SourceInfo))
				}
				s, dependencies, err := structToTs(info, p, st.Name)
				if err != nil {
					ts.Errors = append(ts.Errors, err.Error())
				}
				if _, ok := ts.Pakages[p]; !ok {
					ts.Pakages[p] = TSModule{Structs: make(map[string]string), Types: make(map[string]string), Enums: make(map[string]string), Consts: make(map[string]string), GTypes: make(map[string]string), Endpoints: make(map[string]string)}
				}
				ts.Pakages[p].Structs[st.Name] = s
				ts.AddDependencies(info, p, st.Name, dependencies)
			}
		}

		for _, t := range info.Packages[p].decs {
			ts.Pakages[p].GTypes[t.Name] = fmt.Sprintf("export type %s = %s\n", t.Name, t.Value)
		}

		for _, e := range info.Packages[p].consts {
			ts.Pakages[p].Consts[e.Name] = fmt.Sprint(constToTs(info, p, e.Name))
		}

		for _, e := range info.Packages[p].enums {
			ts.Pakages[p].Enums[e.Name] = fmt.Sprint(enumToTs(info, p, e.Name))
		}

		for _, t := range info.Packages[p].types {
			if t.Typescript {
				ts.Pakages[p].Types[t.Name] = fmt.Sprintf("export type %s = %s\n", t.Name, t.TsType)
			}
		}

		for _, e := range info.Packages[p].endpoints {

			e.VerifyTypes(info, p)

			endpoint := e.ToTs(p)
			if _, ok := ts.Pakages[p]; !ok {
				ts.Pakages[p] = TSModule{Structs: make(map[string]string), Types: make(map[string]string), Enums: make(map[string]string), Consts: make(map[string]string), GTypes: make(map[string]string), Endpoints: make(map[string]string)}
			}
			ts.Pakages[p].Endpoints[e.Name] = endpoint
		}
	}
}
