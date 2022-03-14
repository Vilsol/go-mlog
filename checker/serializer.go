package checker

import (
	"go/ast"
	"strings"
)

type serializedData struct {
	Type string `json:"type"`
}

type serializedFunction struct {
	serializedData
	Comments   []string          `json:"comments,omitempty"`
	Params     []serializedField `json:"params,omitempty"`
	Results    []serializedField `json:"results,omitempty"`
	TypeParams []serializedField `json:"type_params,omitempty"`
}

type serializedField struct {
	Name *string `json:"name,omitempty"`
	Type string  `json:"type"`
}

type serializedValue struct {
	serializedData
	Value    string   `json:"value"`
	Comments []string `json:"comments,omitempty"`
}

func GetSerializablePackages() map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{})

	for pack, files := range packages {
		result[pack] = make(map[string]interface{})
		for _, file := range files {
			for _, decl := range file.Decls {
				switch castDecl := decl.(type) {
				case *ast.FuncDecl:
					var f serializedFunction
					f.Type = "function"

					if castDecl.Doc != nil {
						f.Comments = make([]string, len(castDecl.Doc.List))
						for i, comment := range castDecl.Doc.List {
							f.Comments[i] = comment.Text
						}
					}

					funcType := castDecl.Type

					if funcType.Params != nil && funcType.Params.NumFields() > 0 {
						f.Params = make([]serializedField, 0, funcType.Params.NumFields())
						for _, field := range funcType.Params.List {
							f.Params = append(f.Params, serializeField(field)...)
						}
					}

					if funcType.Results != nil && funcType.Results.NumFields() > 0 {
						f.Results = make([]serializedField, 0, funcType.Results.NumFields())
						for _, field := range funcType.Results.List {
							f.Results = append(f.Results, serializeField(field)...)
						}
					}

					if funcType.TypeParams != nil && funcType.TypeParams.NumFields() > 0 {
						f.TypeParams = make([]serializedField, 0, funcType.TypeParams.NumFields())
						for _, field := range funcType.TypeParams.List {
							f.TypeParams = append(f.TypeParams, serializeField(field)...)
						}
					}

					result[pack][castDecl.Name.String()] = f
				case *ast.GenDecl:
					for _, spec := range castDecl.Specs {
						switch castSpec := spec.(type) {
						case *ast.ValueSpec:
							var comments []string
							if castSpec.Doc != nil {
								comments = make([]string, len(castSpec.Doc.List))
								for i, comment := range castSpec.Doc.List {
									comments[i] = comment.Text
								}
							}

							for i, name := range castSpec.Names {
								value := serializedValue{}
								value.Type = "value"
								value.Value = exprToName(castSpec.Values[i])
								value.Comments = comments
								result[pack][name.String()] = value
							}
						}
					}
				}
			}
		}
	}

	return result
}

func serializeField(field *ast.Field) []serializedField {
	if field.Names != nil && len(field.Names) > 0 {
		results := make([]serializedField, len(field.Names))

		for i, name := range field.Names {
			results[i] = serializedField{
				Name: &name.Name,
				Type: exprToName(field.Type),
			}
		}

		return results
	}

	return []serializedField{
		{
			Type: exprToName(field.Type),
		},
	}
}

func exprToName(expr ast.Expr) string {
	switch castType := expr.(type) {
	case *ast.Ident:
		return castType.String()
	case *ast.SelectorExpr:
		return castType.X.(*ast.Ident).String() + "." + castType.Sel.String()
	case *ast.BasicLit:
		return castType.Value
	case *ast.CallExpr:
		strArgs := make([]string, len(castType.Args))
		for i, arg := range castType.Args {
			strArgs[i] = exprToName(arg)
		}
		return exprToName(castType.Fun) + "(" + strings.Join(strArgs, ", ") + ")"
	}
	return ""
}
