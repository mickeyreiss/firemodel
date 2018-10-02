package golang

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/mickeyreiss/firemodel"
	"github.com/mickeyreiss/firemodel/version"
	"github.com/pkg/errors"
)

func init() {
	firemodel.RegisterModeler("go", &GoModeler{})
}

const (
	fileExtension = ".firemodel.go"
)

type GoModeler struct {
	pkg string
}

func (m *GoModeler) Model(schema *firemodel.Schema, sourceCoder firemodel.SourceCoder) error {
	m.pkg = schema.Options.Get("go")["package"]
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

	f.
		Type().
		Id(model.Name).
		StructFunc(m.fields(model.Name, model.Fields, model.Options.GetAutoTimestamp()))

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
			enumValName := fmt.Sprintf("%s_%s", enumName, val.Name)
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
				Tag(map[string]string{"firestore": strcase.ToLowerCamel(field.Name)})
		}
		if addTimestampFields {
			g.Line()
			g.Comment("Creation timestamp.")
			g.
				Id("CreatedAt").
				Qual("time", "Time").
				Tag(map[string]string{"firestore": "createdAt,serverTimestamp"})

			g.Comment("Update timestamp.")
			g.
				Id("UpdatedAt").
				Qual("time", "Time").
				Tag(map[string]string{"firestore": "updatedAt,serverTimestamp"})
		}
		// if model.Options.GetAutoVersion() {
		// 	g.Line()
		// 	g.Comment("Version number")
		// 	g.
		// 		Id("Version").
		// 		Qual("time", "Time").                                          // What's this?
		// 		Tag(map[string]string{"firestore": "version,serverTimestamp"}) // ?
		//
		// 	g.Comment("Tombstone")
		// 	g.
		// 		Id("Tombstone").
		// 		Qual("time", "Time").                                            // What's this?
		// 		Tag(map[string]string{"firestore": "tombstone,serverTimestamp"}) // ?
		// }
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
		return func(s *jen.Statement) { s.Qual("github.com/mickeyreiss/firemodel/runtime", "URL") }
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
		return func(s *jen.Statement) { s.Op("*").Qual("github.com/mickeyreiss/firemodel/runtime", "File") }
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
