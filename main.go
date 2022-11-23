package tsrpc

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"text/template"
)

type TSConfig struct {
	Url string
}

func getTSSource(config TSConfig) string {
	var tsInfoData = TSInfo{}
	var tsSoucesData = TSSouces{}
	tsInfoData.Populate()
	tsSoucesData.Populate(tsInfoData)

	if len(tsSoucesData.Errors) != 0 {
		err := ""
		for _, v := range tsSoucesData.Errors {
			err += fmt.Sprintln(v)
		}
		exitOnError(errors.New(fmt.Sprintf("Some errors...\n %s", err)))
	}

	tsSource := ""
	tsSource += fmt.Sprintln("\n// Api Class")

	dat, err := os.ReadFile("typescript/sources/fetch.ts")
	if err == nil {
		t, err := template.New("tsRpc").Parse(string(dat))
		if err != nil {
			panic(err)
		}
		var result bytes.Buffer
		err = t.Execute(&result, config)
		if err != nil {
			panic(err)
		}
		tsSource += result.String()
	}

	tsSource += fmt.Sprintln("\n// Global Declarations ")
	for p := range tsSoucesData.Pakages {
		for _, v1 := range tsSoucesData.Pakages[p].GTypes {
			tsSource += fmt.Sprintln(v1)
		}
	}

	for p := range tsSoucesData.Pakages {
		tsSource += fmt.Sprintf("\n//\n// namespace %s\n//\n", p)
		tsSource += fmt.Sprintf("\nexport namespace %s {\n", p)
		for _, v1 := range tsSoucesData.Pakages[p].Structs {
			tsSource += fmt.Sprintln(v1)
		}
		for _, v1 := range tsSoucesData.Pakages[p].Types {
			tsSource += fmt.Sprintln(v1)
		}
		for _, v1 := range tsSoucesData.Pakages[p].Enums {
			tsSource += fmt.Sprintln(v1)
		}
		for _, v1 := range tsSoucesData.Pakages[p].Consts {
			tsSource += fmt.Sprintln(v1)
		}
		for _, v1 := range tsSoucesData.Pakages[p].Endpoints {
			tsSource += fmt.Sprintln(v1)
		}
		tsSource += fmt.Sprintf("}\n\n")
	}
	return tsSource
}
