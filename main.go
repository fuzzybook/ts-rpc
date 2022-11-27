// exportable typescript generated from golang
// Copyright (C) 2022  Fabio Prada

package tsrpc

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"text/template"
	"time"
)

// configuration

type TSConfig struct {
	Url   string
	TsApi *string
	Path  string
}

type tsTemplateData struct {
	Url       string
	CreatedOn time.Time
}

var tsConfig TSConfig

func GetTSSource(config TSConfig) string {
	tsConfig = config
	var tsInfoData = TSInfo{}
	var tsSoucesData = TSSouces{}
	tsInfoData.Populate(tsConfig.Path)
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
	data := ""
	if tsConfig.TsApi == nil {
		data = TsApiTemplate
	} else {
		d, err := os.ReadFile(*tsConfig.TsApi)
		if err != nil {
			panic(err)
		}
		data = string(d)
	}

	t, err := template.New("tsRpc").Parse(data)
	if err != nil {
		panic(err)
	}
	var templateData = tsTemplateData{
		Url:       config.Url,
		CreatedOn: time.Now(),
	}
	var result bytes.Buffer
	err = t.Execute(&result, templateData)
	if err != nil {
		panic(err)
	}
	tsSource += result.String()

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
