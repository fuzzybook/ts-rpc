package tsrpc

import (
	"errors"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type TSInfoPakage struct {
	structs   map[string]TSStruct
	types     map[string]TSType
	enums     map[string]TSEnum
	consts    map[string]TSConst
	decs      map[string]TSDec
	endpoints map[string]TSEndpoint
}

type TSDec struct {
	Name       string
	Value      string
	SourceInfo string
}

type TSType struct {
	Name       string
	Type       string
	TsType     string
	Typescript bool
	dependOn   bool
	SourceInfo string
}

type TSInfo struct {
	Packages map[string]TSInfoPakage
}

type TSConst struct {
	Name  string
	Value string
}

type TSEnum struct {
	Name string
	Info []TSEnumInfo
}

type TSEnumInfo struct {
	Key   string
	Value string
}

/// exporter
type TSModuleInfo struct {
	structs map[string]string
	types   map[string]string
}

type TSSourceLine struct {
	Pos    int
	End    int
	Line   int
	Source string
}

type TSSourceFile struct {
	Source string
	Name   string
	Lines  []TSSourceLine
	Len    int
}

func (ts *TSInfo) findStruct(p string, n string) bool {
	if _, ok := ts.Packages[p]; ok {
		if _, ok := ts.Packages[p].structs[n]; ok {
			return true
		}
	}
	return false
}

func (ts *TSInfo) findType(p string, n string) bool {
	if _, ok := ts.Packages[p]; ok {
		if _, ok := ts.Packages[p].types[n]; ok {
			return true
		}
	}
	return false
}

func (ts TSInfo) find(p string, n string) bool {
	return ts.findType(p, n) || ts.findStruct(p, n)
}

// popola TsInfo con tutte le definizioni dei tipi

func (i *TSInfo) getConst(p string, c *doc.Value, src []TSSourceFile) {
	var isTypescript = strings.HasPrefix(c.Doc, "Typescript:")
	if isTypescript {
		command := strings.TrimPrefix(c.Doc, "Typescript:")
		command = strings.TrimSpace(command)
		command = strings.Trim(command, "\n")

		if strings.Contains(command, "enum=") {
			enumName := strings.TrimPrefix(command, "enum=")
			enum := TSEnum{
				Name: enumName,
				Info: []TSEnumInfo{},
			}
			d := c.Decl
			iota := false
			iotaValue := 0
			for _, s := range d.Specs {
				v := s.(*ast.ValueSpec) // safe because decl.Tok == token.CONST
				if len(v.Values) > 0 {
					be, ok := v.Values[0].(*ast.BinaryExpr)
					if ok {
						x := be.X.(*ast.BasicLit)
						exitOnError(errors.New(fmt.Sprintf("Enum Binary Expression Not implemented %s %s %s AT: %s\n", x.Value, be.Op.String(), be.Y, getSourceInfo(int(x.ValuePos), src))))
					}
					ident, ok := v.Values[0].(*ast.Ident)
					if ok {
						if ident.Name == "iota" {
							iota = true
							iotaValue = v.Names[0].Obj.Data.(int)
							enum.Info = append(enum.Info, TSEnumInfo{Key: v.Names[0].Name, Value: fmt.Sprintf("%d", iotaValue)})
						}
					}
					list, ok := v.Values[0].(*ast.BasicLit)
					if ok {
						enum.Info = append(enum.Info, TSEnumInfo{Key: v.Names[0].Name, Value: fmt.Sprintf("%s", list.Value)})
					}
				} else {
					for _, name := range v.Names {
						if iota {
							iotaValue++
							enum.Info = append(enum.Info, TSEnumInfo{Key: name.Name, Value: fmt.Sprintf("%d", iotaValue)})
						}

					}
				}
			}
			i.Packages[p].enums[enumName] = enum
			t1 := TSType{
				Name:       enumName,
				Typescript: true,
				Type:       "",
				TsType:     fmt.Sprintf("typeof Enum%s[keyof typeof Enum%s] ", enumName, enumName), //getFieldTsInfo(expr.Type),
				dependOn:   false,
				SourceInfo: "",
			}
			i.Packages[p].types[enumName] = t1
		}

		if strings.Contains(command, "const") {
			d := c.Decl
			for _, s := range d.Specs {
				v := s.(*ast.ValueSpec)
				if len(v.Names) == 0 || len(v.Values) == 0 {
					continue
				}
				c := TSConst{
					Name:  v.Names[0].Name,
					Value: v.Values[0].(*ast.BasicLit).Value,
				}
				fmt.Println("const", v.Names, v.Values, c)
				i.Packages[p].consts[c.Name] = c
			}
		}

	}
}

func (i *TSInfo) getType(p string, t *doc.Type, src []TSSourceFile) {
	var isTypescript = strings.HasPrefix(t.Doc, "Typescript:")
	command := ""
	param := ""
	if isTypescript {
		command = strings.TrimPrefix(t.Doc, "Typescript:")
		command = strings.TrimSpace(command)
		command = strings.Trim(command, "\n")

		if strings.Contains(command, "=") {
			a := strings.Split(command, "=")
			if len(a) == 2 {
				param = a[1]
			}
			command = strings.Trim(a[0], " ")
		}
	}
	for _, spec := range t.Decl.Specs {
		if len(t.Consts) > 0 {
			i.getConst(p, t.Consts[0], src)
			continue
		}
		switch spec.(type) {
		case *ast.TypeSpec:
			typeSpec := spec.(*ast.TypeSpec)
			switch typeSpec.Type.(type) {
			case *ast.StructType:
				if isTypescript && command != "interface" {
					exitOnError(errors.New(fmt.Sprintf("\nMismatch delaration for interface %s AT: %s\n", t.Doc, getSourceInfo(int(typeSpec.Name.NamePos), src))))
				}
				v := TSStruct{
					Name:       typeSpec.Name.Name,
					Typescript: isTypescript,
					Fields:     []TSSField{},
					SourceInfo: getSourceInfo(int(typeSpec.Name.NamePos), src),
				}
				v.getStruct(typeSpec, src)
				i.Packages[p].structs[typeSpec.Name.Name] = v
			default:
				if isTypescript && command != "type" {
					exitOnError(errors.New(fmt.Sprintf("\nMismatch delaration for type %s AT: %s\n", t.Doc, getSourceInfo(int(typeSpec.Name.NamePos), src))))
				}
				tsInfo := getFieldTsInfo(typeSpec.Type)
				if command == "type" && param != "" {
					tsInfo = param
				}
				t := TSType{
					Name:       typeSpec.Name.Name,
					Typescript: isTypescript,
					Type:       getFieldInfo(typeSpec.Type),
					TsType:     tsInfo,
					dependOn:   toBeImported(typeSpec.Type.(ast.Expr)),
					SourceInfo: getSourceInfo(int(typeSpec.Name.NamePos), src),
				}
				i.Packages[p].types[typeSpec.Name.Name] = t
			}
		}
	}
}

func (i *TSInfo) Populate() {
	i.Packages = make(map[string]TSInfoPakage)
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				fset := token.NewFileSet()
				packages, err := parser.ParseDir(fset, path, nil, parser.ParseComments)

				if err != nil {
					exitOnError(err)
				}

				for pkg, f := range packages {
					if _, ok := i.Packages[pkg]; !ok {
						i.Packages[pkg] = TSInfoPakage{structs: make(map[string]TSStruct), types: make(map[string]TSType), enums: make(map[string]TSEnum), consts: make(map[string]TSConst), decs: make(map[string]TSDec), endpoints: make(map[string]TSEndpoint)}
					}
					if pkg == "typescript" {
						continue
					}

					var src = []TSSourceFile{}
					for n := range f.Files {

						dat, err := os.ReadFile(n)
						if err == nil {
							lines := []TSSourceLine{}
							pos := 0
							line := 1
							for p, k := range dat {
								if string(k) == "\n" {
									l := TSSourceLine{
										Pos:    pos,
										End:    p,
										Line:   line,
										Source: string(dat[pos:p]),
									}
									lines = append(lines, l)
									pos = p + 1
									line++
									if strings.Contains(l.Source, "// Typescript:") {

										if strings.Contains(l.Source, "TStype=") {
											p := strings.Index(l.Source, "TStype=")
											s := strings.Trim(l.Source[p+len("TStype="):], " ")
											a := strings.Split(s, "=")
											if len(a) == 2 {
												t := TSType{
													Name:       strings.Trim(a[0], " "),
													Typescript: true,
													Type:       "UserDefined",
													TsType:     strings.Trim(a[1], " "),
													dependOn:   false,
													SourceInfo: getSourceInfo(int(l.Pos), src),
												}
												i.Packages[pkg].types[strings.Trim(a[0], " ")] = t
											}

										}
										if strings.Contains(l.Source, "TSDeclaration=") {
											p := strings.Index(l.Source, "TSDeclaration=")
											s := strings.Trim(l.Source[p+len("TSDeclaration="):], " ")
											a := strings.Split(s, "=")
											if len(a) == 2 {
												t := TSDec{
													Name:       strings.Trim(a[0], " "),
													Value:      strings.Trim(a[1], " "),
													SourceInfo: getSourceInfo(int(l.Pos), src),
												}
												i.Packages[pkg].decs[strings.Trim(a[0], " ")] = t
											}
										}

										if strings.Contains(l.Source, "TSEndpoint= ") {
											e := ParseEndpoint(l.Source, n, l.Line)
											if _, ok := i.Packages[pkg].endpoints[e.Name]; ok {
												exitOnError(errors.New(fmt.Sprintf("Enpoint name %s allready in use: %s", e.Name, l.Source)))
											}

											i.Packages[pkg].endpoints[e.Name] = e
										}

									}
								}
							}
							s := TSSourceFile{
								Name:   n,
								Source: string(dat),
								Len:    len(dat),
								Lines:  lines,
							}
							src = append(src, s)
						}
					}

					p := doc.New(f, "./", 0)

					for _, t := range p.Types {
						i.getType(pkg, t, src)
					}

					for _, c := range p.Consts {
						i.getConst(pkg, c, src)
					}
				}
			}
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}
}
