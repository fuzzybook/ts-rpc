package tsrpc

import (
	"errors"
	"fmt"
	"go/ast"
)

type TSSField struct {
	Name       string
	Type       string
	TsType     string
	Json       TSTagJson
	Ts         TSTagTs
	DependOn   bool
	SourceInfo string
}

type TSStruct struct {
	Name       string
	Typescript bool
	Fields     []TSSField
	SourceInfo string
}

func isNativeType(t string) bool {
	switch t {
	case "uint8", "uint16", "uint32", "uint64", "uint",
		"int8", "int16", "int32", "int64", "int",
		"float32", "float64":
		return true
	case "bool":
		return true
	case "string":
		return true
	}
	return false
}

func toBeImported(t ast.Expr) bool {
	switch ft := t.(type) {
	case *ast.Ident:
		return !isNativeType(ft.Name)
	case *ast.SelectorExpr:
		return true

	}
	return false
}

func typeToTypescript(k string) string {
	switch k {
	case "uint8", "uint16", "uint32", "uint64", "uint",
		"int8", "int16", "int32", "int64", "int",
		"float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "string":
		return "string"
	}
	return k
}

func getFieldInfo(t ast.Expr) string {
	result := ""
	switch ft := t.(type) {
	case *ast.Ident:
		result = ft.Name
	case *ast.SelectorExpr:
		result = fmt.Sprintf("%s.%s", ft.X, ft.Sel)
	case *ast.ArrayType:
		result = fmt.Sprintf("[]%s", getFieldInfo(ft.Elt))
	case *ast.StarExpr:
		result = fmt.Sprintf("*%s", getFieldInfo(ft.X))
	case *ast.MapType:
		result = fmt.Sprintf("map[%s]%s", ft.Key, getFieldInfo(ft.Value))
	case *ast.InterfaceType:
		result = "interface{}"
	default:
		exitOnError(errors.New(fmt.Sprintf("this go type: %T is not evaluated!\n", ft)))
	}
	return result
}

func getFieldTsInfo(t ast.Expr) string {
	result := ""
	switch ft := t.(type) {
	case *ast.Ident:
		result = typeToTypescript(ft.Name)
	case *ast.SelectorExpr:
		result = typeToTypescript(fmt.Sprintf("%s.%s", ft.X, ft.Sel))
	case *ast.ArrayType:
		result = fmt.Sprintf("%s[]", typeToTypescript(getFieldTsInfo(ft.Elt)))
	case *ast.StarExpr:
		result = fmt.Sprintf("Nullable<%s>", typeToTypescript(getFieldTsInfo(ft.X)))
	case *ast.MapType:
		result = fmt.Sprintf("Record<%s , %s>", typeToTypescript(fmt.Sprintf("%s", ft.Key)), typeToTypescript(getFieldTsInfo(ft.Value)))
	case *ast.InterfaceType:
		result = "unknown"
	default:
		exitOnError(errors.New(fmt.Sprintf("this typescript type: %T is not evaluated!\n", ft)))
	}
	return result
}

func getSourceInfo(pos int, src []TSSourceFile) string {
	for _, v := range src {
		for _, l := range v.Lines {
			if pos >= l.Pos && pos <= l.End {
				return fmt.Sprintf("%s Line: %d", v.Name, l.Line)
			}
		}
	}
	return ""
}

func (s *TSStruct) getStruct(ts *ast.TypeSpec, src []TSSourceFile) {
	if st, ok := ts.Type.(*ast.StructType); ok {
		for _, field := range st.Fields.List {

			tag := ""
			if field.Tag != nil {
				tag = field.Tag.Value
			}

			tagJson := TSTagJson{}
			tagJson.parse(tag)

			tagTs := TSTagTs{}
			tagTs.parse(tag)

			tsType := ""

			if len(field.Names) > 0 {
				tsType = getFieldTsInfo(field.Type.(ast.Expr))
				var f = TSSField{
					Name:       field.Names[0].String(),
					Json:       tagJson,
					Ts:         tagTs,
					Type:       getFieldInfo(field.Type.(ast.Expr)),
					TsType:     tsType,
					DependOn:   toBeImported(field.Type.(ast.Expr)),
					SourceInfo: getSourceInfo(int(field.Type.Pos()), src),
				}
				s.Fields = append(s.Fields, f)
			} else {
				if se, ok := field.Type.(*ast.SelectorExpr); ok {
					var f = TSSField{
						Name:       fmt.Sprintf("%s.%s", se.X, se.Sel),
						Json:       tagJson,
						Ts:         tagTs,
						Type:       getFieldInfo(field.Type.(ast.Expr)),
						TsType:     tsType,
						DependOn:   toBeImported(field.Type.(ast.Expr)),
						SourceInfo: getSourceInfo(int(field.Type.Pos()), src),
					}
					s.Fields = append(s.Fields, f)
				} else {
					exitOnError(errors.New(fmt.Sprintf("this typescript type: %T is not evaluated!\n", field.Type)))
				}
			}
		}
	}
}
