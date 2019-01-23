package golang

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/visor-tax/firemodel"
	"github.com/visor-tax/firemodel/version"
)

func init() {
	firemodel.RegisterModeler("go", &GoModeler{})
}

const (
	fileExtension = ".firemodel.go"
)

type GoModeler struct {
	pkg         string
	clientNames []*ClientName
}

type ClientName struct {
	ModelName  string
	ClientName string
}

func (m *GoModeler) Model(schema *firemodel.Schema, sourceCoder firemodel.SourceCoder) error {
	m.pkg = schema.Options.Get("go")["package"]
	m.clientNames = []*ClientName{}
	for _, model := range schema.Models {
		if err := m.writeModel(model, sourceCoder); err != nil {
			return err
		}
	}
	for _, structType := range schema.Structs {
		if err := m.writeStruct(structType, sourceCoder); err != nil {
			return err
		}
	}
	for _, enum := range schema.Enums {
		if err := m.writeEnum(enum, sourceCoder); err != nil {
			return err
		}
	}

	if err := m.writeManifest(sourceCoder); err != nil {
		return err
	}

	return nil
}

func (m *GoModeler) writeManifest(sourceCoder firemodel.SourceCoder) error {
	f := jen.NewFile(m.packageName())
	f.HeaderComment(fmt.Sprintf("DO NOT EDIT - Code generated by firemodel %s.", version.Version))

	f.Type().Id("Client").StructFunc(func(g *jen.Group) {
		g.Id("Client").Op("*").Qual("cloud.google.com/go/firestore", "Client")
		for _, client := range m.clientNames {

			g.Id(client.ModelName).Op("*").Id(client.ClientName)
		}
	})

	f.Func().Id("NewClient").Params(jen.Id("client").Op("*").Qual("cloud.google.com/go/firestore", "Client")).Id("*Client").BlockFunc(func(g *jen.Group) {
		g.Id("temp").Op(":=").Op("&").Id("Client").ValuesFunc(func(g *jen.Group) {
			g.Id("Client").Op(":").Id("client")
		})
		for _, client := range m.clientNames {
			g.Id("temp." + client.ModelName).Op("=").Op("&").Id(client.ClientName).ValuesFunc(func(g *jen.Group) {
				g.Id("client").Op(":").Id("temp")
			})
		}
		g.Return(jen.Id("temp"))
	})

	w, err := sourceCoder.NewFile("module.go")
	if err != nil {
		return errors.Wrap(err, "firemodel/go: open source code file")
	}
	defer w.Close()

	if err := f.Render(w); err != nil {
		return err
	}
	return nil
}

func (m *GoModeler) writeModel(model *firemodel.SchemaModel, sourceCoder firemodel.SourceCoder) error {
	f := jen.NewFile(m.packageName())
	f.HeaderComment(fmt.Sprintf("DO NOT EDIT - Code generated by firemodel %s.", version.Version))

	if model.Comment == "" {
		f.Commentf("TODO: Add comment to %s in firemodel schema.", model.Name)
	} else {
		f.Comment(model.Comment)
	}

	format, args, err := model.Options.GetFirestorePath()
	if err != nil {
		return errors.Wrap(err, "firemodel/go: invalid firestore path")
	}
	if format != "" {
		commentargs := make([]interface{}, len(args))
		for idx, arg := range args {
			commentargs[idx] = fmt.Sprintf("{%s}", arg)
		}
		f.Comment("")
		f.Commentf("Firestore document location: /%s", fmt.Sprintf(format, commentargs...))
	}
	clientName := fmt.Sprint("client", model.Name)
	f.
		Type().
		Id(model.Name).
		StructFunc(func(g *jen.Group) {
			m.fields(model.Name, model.Fields, model.Options.GetAutoTimestamp())(g)

		})

	if format, args, err := model.Options.GetFirestorePath(); format != "" {
		f.
			Commentf("%s returns the path to a particular %s in Firestore.", fmt.Sprint(model.Name, "Path"), model.Name)
		f.
			Func().
			Id(fmt.Sprint(model.Name, "Path")).
			ParamsFunc(func(g *jen.Group) {
				if err != nil {
					panic(err)
				}
				for _, arg := range args {
					g.Id(strcase.ToLowerCamel(arg)).String()
				}
			}).
			String().
			Block(jen.ReturnFunc(func(g *jen.Group) {
				if err != nil {
					panic(err)
				}
				g.
					Qual("fmt", "Sprintf").
					CallFunc(func(g *jen.Group) {
						g.Lit(format)
						for _, arg := range args {
							g.Id(strcase.ToLowerCamel(arg))
						}
					})
			}))
		f.Commentf("%s is a regex that can be use to filter out firestore events of %s", fmt.Sprint(model.Name, "RegexPath"), model.Name)
		f.Var().Id(fmt.Sprint(model.Name, "RegexPath")).Op("=").Qual("regexp", "MustCompile").CallFunc(func(g *jen.Group) {
			regex := regexp.QuoteMeta(format)
			start := "^(?:projects/[^/]*/databases/[^/]*/documents/)?(?:/)?"
			g.Lit(fmt.Sprint(start, strings.Replace(regex, "%s", "([a-zA-Z0-9]+)", -1), "$"))
		})

		f.Commentf("%s is a named regex that can be use to filter out firestore events of %s", fmt.Sprint(model.Name, "RegexNamedPath"), model.Name)
		f.Var().Id(fmt.Sprint(model.Name, "RegexNamedPath")).Op("=").Qual("regexp", "MustCompile").CallFunc(func(g *jen.Group) {
			regex := regexp.QuoteMeta(format)
			start := "^(?:projects/[^/]*/databases/[^/]*/documents/)?(?:/)?"
			for _, arg := range args {
				repl := fmt.Sprint("(?P<", arg, ">[a-zA-Z0-9]+)")
				regex = strings.Replace(regex, "%s", repl, 1)
			}
			g.Lit(fmt.Sprint(start, regex, "$"))
		})

		pathStructName := fmt.Sprint(model.Name, "PathStruct")
		pathStructFunctionName := fmt.Sprint(model.Name, "PathToStruct")
		pathStructReverseFunctionName := fmt.Sprint(model.Name, "StructToPath")
		f.Commentf("%s is a struct that contains parts of a path of %s", pathStructName, model.Name)
		f.Type().Id(pathStructName).StructFunc(func(g *jen.Group) {
			for _, arg := range args {
				g.Id(strcase.ToCamel(arg)).String()
			}
		})

		f.Commentf("%s is a function that turns a firestore path into a PathStruct of %s", pathStructFunctionName, model.Name)
		f.
			Func().
			Id(pathStructFunctionName).Params(jen.Id("path").String()).Id("*" + pathStructName).BlockFunc(func(g *jen.Group) {
			g.Id("parsed").Op(":=").Id(fmt.Sprint(model.Name, "RegexPath")).Dot("FindStringSubmatch").Call(jen.Id("path"))
			g.Id("result").Op(":=").Op("&").Id(pathStructName).ValuesFunc(func(g *jen.Group) {
				for i, arg := range args {
					g.Id(strcase.ToCamel(arg)).Op(":").Id("parsed").Index(jen.Lit(i + 1))
				}
			})
			g.Return(jen.Id("result"))
		})

		f.Commentf("%s is a function that turns a PathStruct of %s into a firestore path", pathStructReverseFunctionName, model.Name)
		f.Func().Id(pathStructReverseFunctionName).Params(jen.Id("path").Id("*" + pathStructName)).String().BlockFunc(func(g *jen.Group) {
			g.Id("built").Op(":=").Qual("fmt", "Sprintf").CallFunc(func(g *jen.Group) {
				g.Lit(format)
				for _, arg := range args {
					g.Id("path").Dot(strcase.ToCamel(arg))
				}
			})
			g.Return(jen.Id("built"))
		})

		wrapperName := fmt.Sprint(model.Name, "Wrapper")
		f.Commentf("%s is a struct wrapper that contains a reference to the firemodel instance and the path", wrapperName)
		f.Type().Id(wrapperName).StructFunc(func(g *jen.Group) {
			g.Id("Data").Id("*" + model.Name)
			g.Id("Path").Id("*" + pathStructName)
			g.Id("PathStr").String()
			g.Comment("---- Internal Stuffs ----")
			g.Id("client").Op("*").Id(clientName)
			g.Id("pathStr").String()
			g.Id("ref").Op("*").Qual("cloud.google.com/go/firestore", "DocumentRef")
		})

		fromSnapshotName := fmt.Sprint(model.Name, "FromSnapshot")
		f.Commentf("%s is a function that will create an instance of the model from a document snapshot", fromSnapshotName)
		f.Func().
			Id(fromSnapshotName).
			Params(
				jen.Id("snapshot").
					Op("*").Qual("cloud.google.com/go/firestore", "DocumentSnapshot")).
			Params(
				jen.Id("*"+wrapperName),
				jen.Error()).
			BlockFunc(func(g *jen.Group) {
				g.Id("temp").Op(":=").Op("&").Id(model.Name).Values()
				g.Err().Op(":=").Id("snapshot").Dot("DataTo").Call(jen.Id("temp"))
				g.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err()))
				g.Id("path").Op(":=").Id(pathStructFunctionName).Call(jen.Id("snapshot.Ref.Path"))
				g.Id("pathStr").Op(":=").Id(pathStructReverseFunctionName).Call(jen.Id("path"))
				g.Id("wrapper").Op(":=").Id("&" + wrapperName).ValuesFunc(func(g *jen.Group) {
					g.Id("Path").Op(":").Id("path")
					g.Id("PathStr").Op(":").Id("pathStr")
					g.Id("pathStr").Op(":").Id("pathStr")
					g.Id("ref").Op(":").Id("snapshot.Ref")
					g.Id("Data").Op(":").Id("temp")
				})
				g.Return(jen.Id("wrapper"), jen.Nil())
			})

		m.clientNames = append(m.clientNames, &ClientName{ClientName: clientName, ModelName: model.Name})
		f.Type().Id(clientName).StructFunc(func(g *jen.Group) {
			g.Id("client").Op("*").Id("Client")
		})

		// Disable create for now, only set
		// f.Func().Params(jen.Id("c").Op("*").Id(clientName)).Id("Create").Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("path").String(), jen.Id("model").Op("*").Id(model.Name)).Params(jen.Op("*").Id(wrapperName), jen.Error()).BlockFunc(func(g *jen.Group) {
		// 	g.Id("ref").Op(":=").Id("c").Dot("client").Dot("Client").Dot("Doc").Call(jen.Id("path"))
		// 	g.Id("wrapper").Op(":=").Op("&").Id(wrapperName).ValuesFunc(func(g *jen.Group) {
		// 		g.Id("ref").Op(":").Id("ref")
		// 		g.Id("pathStr").Op(":").Id("path")
		// 		g.Id("PathStr").Op(":").Id("path")
		// 		g.Id("Path").Op(":").Id(pathStructFunctionName).Call(jen.Id("path"))
		// 		g.Id("client").Op(":").Id("c")
		// 		g.Id("Data").Op(":").Id("model")
		// 	})
		// 	g.Id("wrapper").Dot("Data").Dot("UpdatedAt").Op("=").Qual("time", "Now").Call()
		// 	g.Id("wrapper").Dot("Data").Dot("CreatedAt").Op("=").Qual("time", "Now").Call()
		// 	g.Err().Op(":=").Id("wrapper").Dot("Create").Call(jen.Id("ctx"))
		// 	g.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err()))
		// 	g.Return(jen.Id("wrapper"), jen.Nil())
		// })

		f.Func().Params(jen.Id("c").Op("*").Id(clientName)).Id("Set").Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("path").String(), jen.Id("model").Op("*").Id(model.Name)).Params(jen.Op("*").Id(wrapperName), jen.Error()).BlockFunc(func(g *jen.Group) {
			g.Id("ref").Op(":=").Id("c").Dot("client").Dot("Client").Dot("Doc").Call(jen.Id("path"))
			g.Id("snapshot").Op(",").Err().Op(":=").Id("ref").Dot("Get").Call(jen.Id("ctx"))
			g.If(jen.Id("snapshot").Dot("Exists").Call()).BlockFunc(func(g *jen.Group) {
				g.Id("temp").Op(",").Err().Op(":=").Id(fromSnapshotName).Call(jen.Id("snapshot"))
				g.If(jen.Err().Op("!=").Nil()).Block(jen.Comment("Don't do anything, just override")).Else().BlockFunc(func(g *jen.Group) {
					g.Id("model").Dot("CreatedAt").Op("=").Id("temp").Dot("Data").Dot("CreatedAt")
				})
			})
			g.Id("wrapper").Op(":=").Op("&").Id(wrapperName).ValuesFunc(func(g *jen.Group) {
				g.Id("ref").Op(":").Id("ref")
				g.Id("pathStr").Op(":").Id("path")
				g.Id("PathStr").Op(":").Id("path")
				g.Id("Path").Op(":").Id(pathStructFunctionName).Call(jen.Id("path"))
				g.Id("client").Op(":").Id("c")
				g.Id("Data").Op(":").Id("model")
			})
			g.Id("wrapper").Dot("Data").Dot("UpdatedAt").Op("=").Qual("time", "Now").Call()
			g.Err().Op("=").Id("wrapper").Dot("Set").Call(jen.Id("ctx"))
			g.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err()))
			g.Return(jen.Id("wrapper"), jen.Nil())
		})

		getCommandByPathName := fmt.Sprint("Get", "ByPath")

		f.Func().Params(jen.Id("c").Id("*"+clientName)).Id(getCommandByPathName).Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("path").String()).Params(
			jen.Id("*"+wrapperName),
			jen.Error()).
			BlockFunc(func(g *jen.Group) {
				g.Id("reference").Op(":=").Id("c").Dot("client").Dot("Client").Dot("Doc").Call(jen.Id("path"))
				g.Id("snapshot").Op(",").Err().Op(":=").Id("reference").Dot("Get").Call(jen.Id("ctx"))
				g.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err()))
				g.Id("wrapper").Op(",").Err().Op(":=").Id(fromSnapshotName).Call(jen.Id("snapshot"))
				g.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err()))
				g.Return(jen.Id("wrapper"), jen.Nil())
			})

		f.Func().Params(jen.Id("c").Id("*"+clientName)).Id(getCommandByPathName+"Tx").Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("tx").Op("*").Qual("cloud.google.com/go/firestore", "Transaction"), jen.Id("path").String()).Params(
			jen.Id("*"+wrapperName),
			jen.Error()).
			BlockFunc(func(g *jen.Group) {
				g.Id("reference").Op(":=").Id("c").Dot("client").Dot("Client").Dot("Doc").Call(jen.Id("path"))
				g.Id("snapshot").Op(",").Err().Op(":=").Id("tx").Dot("Get").Call(jen.Id("reference"))
				g.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err()))
				g.Id("wrapper").Op(",").Err().Op(":=").Id(fromSnapshotName).Call(jen.Id("snapshot"))
				g.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err()))
				g.Return(jen.Id("wrapper"), jen.Nil())
			})

		f.Func().Params(jen.Id("m").Id("*" + wrapperName)).Id("Set").Params(jen.Id("ctx").Qual("context", "Context")).Params(jen.Id("error")).BlockFunc(func(g *jen.Group) {
			g.If(jen.Id("m.ref").Op("==").Nil()).BlockFunc(func(g *jen.Group) {
				g.Return(jen.Qual("errors", "New").Call(jen.Lit("Cannot call set on a firemodel object that has no reference. Call `create` on the orm with this object instead")))
			})
			g.Id("_").Op(",").Err().Op(":=").Id("m").Dot("ref").Dot("Set").Call(jen.Id("ctx"), jen.Id("m").Dot("Data"))
			g.Return(jen.Err())
		})

		f.Func().Params(jen.Id("m").Id("*"+wrapperName)).Id("SetTx").Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("tx").Op("*").Qual("cloud.google.com/go/firestore", "Transaction")).Params(jen.Id("error")).BlockFunc(func(g *jen.Group) {
			g.If(jen.Id("m.ref").Op("==").Nil()).BlockFunc(func(g *jen.Group) {
				g.Return(jen.Qual("errors", "New").Call(jen.Lit("Cannot call set on a firemodel object that has no reference. Call `create` on the orm with this object instead")))
			})
			g.Err().Op(":=").Id("tx").Dot("Set").Call(jen.Id("m").Dot("ref"), jen.Id("m").Dot("Data"))
			g.Return(jen.Err())
		})

	}

	w, err := sourceCoder.NewFile(fmt.Sprint(strcase.ToSnake(model.Name), fileExtension))
	if err != nil {
		return errors.Wrap(err, "firemodel/go: open source code file")
	}
	defer w.Close()

	if err := f.Render(w); err != nil {
		return err
	}
	return nil
}

func (m *GoModeler) writeEnum(enum *firemodel.SchemaEnum, sourceCoder firemodel.SourceCoder) error {
	enumName := strcase.ToCamel(enum.Name)
	f := jen.NewFile(m.packageName())
	f.HeaderComment(fmt.Sprintf("DO NOT EDIT - Code generated by firemodel %s.", version.Version))

	if enum.Comment == "" {
		f.Commentf("TODO: Add comment to %s in firemodel schema.", enumName)
	} else {
		f.Comment(enum.Comment)
	}
	f.Type().Id(enumName).String()

	f.Const().DefsFunc(func(g *jen.Group) {
		for _, val := range enum.Values {
			enumValName := fmt.Sprintf("%s_%s", enumName, strcase.ToScreamingSnake(val.Name))
			if enum.Comment == "" {
				g.Commentf("TODO: Add comment to %s in firemodel schema.", enumValName)
			} else {
				g.Comment(val.Comment)
			}
			g.
				Id(enumValName).
				Id(enumName).
				Op("=").
				Lit(strcase.ToScreamingSnake(val.Name))
		}
	})
	f.Var().Id(enumName + "_Strings").Op("=").Map(jen.Id(enumName)).String().ValuesFunc(func(g *jen.Group) {
		for _, val := range enum.Values {
			enumValName := fmt.Sprintf("%s_%s", enumName, strcase.ToScreamingSnake(val.Name))
			g.Id(enumValName).Op(":").Lit(enumValName)
		}
	})
	f.Var().Id(enumName + "_Values").Op("=").Map(jen.String()).Id(enumName).ValuesFunc(func(g *jen.Group) {
		for _, val := range enum.Values {
			enumValName := fmt.Sprintf("%s_%s", enumName, strcase.ToScreamingSnake(val.Name))
			g.Lit(enumValName).Op(":").Id(enumValName)
		}
	})
	f.Func().Params(jen.Id("e").Id(enumName)).Id("String").Params().String().Block(jen.Return(jen.Id(enumName + "_Strings").Index(jen.Id("e"))))
	w, err := sourceCoder.NewFile(fmt.Sprint(strcase.ToSnake(enum.Name), fileExtension))
	if err != nil {
		return errors.Wrap(err, "firemodel/go: open source code file")
	}

	defer w.Close()

	if err := f.Render(w); err != nil {
		return err
	}
	return nil
}

func (m *GoModeler) writeStruct(structType *firemodel.SchemaStruct, sourceCoder firemodel.SourceCoder) error {
	structName := strcase.ToCamel(structType.Name)
	f := jen.NewFile(m.packageName())
	f.HeaderComment(fmt.Sprintf("DO NOT EDIT - Code generated by firemodel %s.", version.Version))

	if structType.Comment == "" {
		f.Commentf("TODO: Add comment to %s in firemodel schema.", structType.Name)
	} else {
		f.Comment(structType.Comment)
	}
	f.Type().Id(structName).StructFunc(m.fields(structName, structType.Fields, false))

	w, err := sourceCoder.NewFile(fmt.Sprint(strcase.ToSnake(structType.Name), fileExtension))
	if err != nil {
		return errors.Wrap(err, "firemodel/go: open source code file")
	}

	defer w.Close()

	if err := f.Render(w); err != nil {
		return err
	}
	return nil
}

func (m *GoModeler) packageName() string {
	if m.pkg == "" {
		return "firemodel"
	}
	return m.pkg
}

func (m *GoModeler) fieldTags(field *firemodel.SchemaField) string {
	switch field.Type.(type) {
	// "false" and "0" should be written
	case *firemodel.Boolean,
		*firemodel.Integer,
		*firemodel.Double:
		return strcase.ToLowerCamel(field.Name)

	default:
		return strcase.ToLowerCamel(field.Name) + ",omitempty"
	}

}

func (m *GoModeler) fields(structName string, fields []*firemodel.SchemaField, addTimestampFields bool) func(g *jen.Group) {
	return func(g *jen.Group) {
		for _, field := range fields {
			if field.Comment == "" {
				g.Commentf("TODO: Add comment to %s.%s.", structName, field.Name)
			} else {
				g.Comment(field.Comment)
			}

			g.
				Id(strcase.ToCamel(field.Name)).
				Do(m.goType(field.Type)).
				Tag(map[string]string{"firestore": m.fieldTags(field)})
		}
		if addTimestampFields {
			g.Line()
			g.Comment("Creation timestamp.")
			g.
				Id("CreatedAt").
				Qual("time", "Time").
				Tag(map[string]string{"firestore": "createdAt"})

			g.Comment("Update timestamp.")
			g.
				Id("UpdatedAt").
				Qual("time", "Time").
				Tag(map[string]string{"firestore": "updatedAt"})
		}
	}
}

func (m *GoModeler) goType(firetype firemodel.SchemaFieldType) func(s *jen.Statement) {
	switch firetype := firetype.(type) {
	case *firemodel.Boolean:
		return func(s *jen.Statement) { s.Bool() }
	case *firemodel.Integer:
		return func(s *jen.Statement) { s.Int64() }
	case *firemodel.Double:
		return func(s *jen.Statement) { s.Float64() }
	case *firemodel.Timestamp:
		return func(s *jen.Statement) { s.Qual("time", "Time") }
	case *firemodel.String:
		return func(s *jen.Statement) { s.String() }
	case *firemodel.URL:
		return func(s *jen.Statement) { s.Qual("github.com/visor-tax/firemodel/runtime", "URL") }
	case *firemodel.Enum:
		return func(s *jen.Statement) { s.Id(firetype.T.Name) }
	case *firemodel.Bytes:
		return func(s *jen.Statement) { s.Index().Byte() }
	case *firemodel.Reference:
		return func(s *jen.Statement) { s.Op("*").Qual("cloud.google.com/go/firestore", "DocumentRef") }
	case *firemodel.GeoPoint:
		return func(s *jen.Statement) { s.Op("*").Qual("google.golang.org/genproto/googleapis/type/latlng", "LatLng") }
	case *firemodel.Struct:
		return func(s *jen.Statement) { s.Op("*").Id(firetype.T.Name) }
	case *firemodel.Array:
		if firetype.T != nil {
			return func(s *jen.Statement) { s.Index().Do(m.goType(firetype.T)) }
		}
		return func(s *jen.Statement) { s.Index().Interface() }
	case *firemodel.File:
		return func(s *jen.Statement) { s.Op("*").Qual("github.com/visor-tax/firemodel/runtime", "File") }
	case *firemodel.Map:
		if firetype.T != nil {
			return func(s *jen.Statement) { s.Map(jen.String()).Do(m.goType(firetype.T)) }
		}
		return func(s *jen.Statement) { s.Map(jen.String()).Interface() }
	default:
		err := errors.Errorf("firemodel/go: unknown type %s", firetype)
		panic(err)
	}
}
