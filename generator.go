package main

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	structDefinition = `
	type %s struct {
		%s
	}
`
	fieldDefinition   = "%s %s %s\n"
	jsonTagDefinition = "`json:\"%s\"`"
	enumDefinition    = `
	type %s int32
	const (
		%s
	)
`
	serviceDefinition = `
	type %s struct {}
`
	handlerDefinition   = "func %s(ctx context.Context, params json.RawMessage) (interface{}, error) { return nil, nil }\n"
	interfaceDefinition = `
	type %s interface {
		%s
	}
`
	interfaceMethodDefinition = "%s(%s) (%s)\n"
)

type jsonRpcGenerator struct {
	plugin *protogen.Plugin
}

func (g *jsonRpcGenerator) generate() error {
	for _, file := range g.plugin.Files {
		var buf bytes.Buffer
		pkg := fmt.Sprintf("package %s\n\n", file.GoPackageName)
		buf.Write([]byte(pkg))

		// create services
		for _, service := range file.Services {
			buf.WriteString(makeServiceDefinition(service))
			buf.WriteString(makeInterface(service))
		}

		// create Go types and constants for top-level enums
		for _, enum := range file.Enums {
			goEnum := makeEnum(enum)
			buf.WriteString(goEnum + "\n\n")
		}

		// create Go structs and enums from messages
		enums, structs, err := makeCustomTypes(file.Messages)
		if err != nil {
			return fmt.Errorf("cannot make structs: %w", err)
		}
		buf.WriteString(enums)
		buf.WriteString(structs)

		filename := file.GeneratedFilenamePrefix + ".pjson.go"
		file := g.plugin.NewGeneratedFile(filename, ".")
		_, err = file.Write(buf.Bytes())
		if err != nil {
			return fmt.Errorf("write error: %w", err)
		}
	}

	return nil
}

// makeCustomTypes creates Go structs and enums from all message definitions.
func makeCustomTypes(messages []*protogen.Message) (string, string, error) {
	enums := strings.Builder{}
	structs := strings.Builder{}
	for _, msg := range messages {
		// skip map entries
		if msg.Desc.IsMapEntry() {
			continue
		}
		// create nested enums
		for _, enum := range msg.Enums {
			enumStr := makeEnum(enum)
			enums.WriteString(enumStr)
		}

		// create struct for current message
		str, err := makeStruct(msg)
		if err != nil {
			return "", "", fmt.Errorf("cannot make struct: %w", err)
		}

		structs.WriteString(str)

		// create structs for nested message declarations
		en, str, err := makeCustomTypes(msg.Messages)
		enums.WriteString(en)
		structs.WriteString(str)
	}

	return enums.String(), structs.String(), nil
}

// makeServiceDefinition creates Go struct which represents a service.
func makeServiceDefinition(service *protogen.Service) string {
	name := service.GoName
	name = string(unicode.ToLower(rune(name[0]))) + name[1:]
	return fmt.Sprintf(serviceDefinition, name)
}

// makeInterface creates Go interface definition from protobuf service.
func makeInterface(service *protogen.Service) string {
	interfaceMethods := strings.Builder{}
	for _, method := range service.Methods {
		interfaceMethods.WriteString(fmt.Sprintf(interfaceMethodDefinition, method.GoName, method.Input.Desc.Name(), method.Output.Desc.Name()))
	}

	return fmt.Sprintf(interfaceDefinition, service.GoName, interfaceMethods.String())
}

// makeEnum creates Go type from int32 and constants.
func makeEnum(enum *protogen.Enum) string {
	values := strings.Builder{}
	for _, val := range enum.Values {
		values.WriteString(fmt.Sprintf("%s_%s %s = %d\n", enum.Desc.Name(), val.Desc.Name(), enum.Desc.Name(), val.Desc.Number()))
	}

	return fmt.Sprintf(enumDefinition, enum.Desc.Name(), values.String())
}

// makeStruct makes Go struct from message.
func makeStruct(msg *protogen.Message) (string, error) {
	fields := strings.Builder{}
	for _, field := range msg.Fields {
		// create Go struct field from each proto field
		fieldType, err := makeType(field.Desc)
		if err != nil {
			return "", fmt.Errorf("cannot make type: %w", err)
		}
		// write field comment
		fields.WriteString(field.Comments.Leading.String())
		// write Go field with json tag
		fields.WriteString(fmt.Sprintf(fieldDefinition, field.GoName, fieldType, fmt.Sprintf(jsonTagDefinition, field.Desc.JSONName())))
	}

	structString := fmt.Sprintf(structDefinition, msg.Desc.Name(), fields.String())
	return fmt.Sprintf("%s%s", strings.TrimSuffix(msg.Comments.Leading.String(), "\n"), structString), nil
}

// makeType returns type of field for struct definition.
func makeType(descriptor protoreflect.FieldDescriptor) (string, error) {
	// maps
	if descriptor.IsMap() {
		// get key type
		keyType, err := makeType(descriptor.MapKey())
		if err != nil {
			return "", fmt.Errorf("cannot make type for map key: %w", err)
		}

		// get value type
		valueType, err := makeType(descriptor.MapValue())
		if err != nil {
			return "", fmt.Errorf("cannot make type for map key: %w", err)
		}

		return fmt.Sprintf("map[%s]%s", keyType, valueType), nil
	}

	// messages
	if descriptor.Kind() == protoreflect.MessageKind {
		// repeated
		if descriptor.IsList() {
			return fmt.Sprintf("[]*%s", descriptor.Message().Name()), nil
		}
		// ordinary
		return fmt.Sprintf("*%s", descriptor.Message().Name()), nil
	}

	// enums
	if descriptor.Kind() == protoreflect.EnumKind {
		// repeated
		if descriptor.IsList() {
			return fmt.Sprintf("[]%s", descriptor.Enum().Name()), nil
		}
		// ordinary
		return string(descriptor.Enum().Name()), nil
	}

	// repeated scalar types
	if descriptor.IsList() {
		slice, err := makeSliceOfScalarType(descriptor)
		if err != nil {
			return "", fmt.Errorf("cannot make slice: %w", err)
		}

		return slice, nil
	}

	// ordinary scalar types
	scalar, err := makeScalarType(descriptor)
	if err != nil {
		return "", fmt.Errorf("cannot make scalar type: %w", err)
	}

	return scalar, nil
}

// makeSliceOfScalarType returns slice of probuf types.
func makeSliceOfScalarType(descriptor protoreflect.FieldDescriptor) (string, error) {
	scalarType, err := makeScalarType(descriptor)
	if err != nil {
		return "", fmt.Errorf("cannot make scalar type: %w", err)
	}

	return fmt.Sprintf("[]%s", scalarType), nil
}

// makeScalarType returns Go type for scalar profobuf types.
func makeScalarType(descriptor protoreflect.FieldDescriptor) (string, error) {
	var goType string
	goInt32 := "int32"
	goInt64 := "int64"
	switch descriptor.Kind() {
	case protoreflect.BoolKind:
		goType = "bool"
	case protoreflect.Int32Kind:
		goType = goInt32
	case protoreflect.Sint32Kind:
		goType = goInt32
	case protoreflect.Uint32Kind:
		goType = "uint32"
	case protoreflect.Int64Kind:
		goType = goInt64
	case protoreflect.Sint64Kind:
		goType = goInt64
	case protoreflect.Uint64Kind:
		goType = "uint64"
	case protoreflect.Sfixed32Kind:
		goType = goInt32
	case protoreflect.Fixed32Kind:
		goType = "uint32"
	case protoreflect.FloatKind:
		goType = "float32"
	case protoreflect.Sfixed64Kind:
		goType = goInt64
	case protoreflect.Fixed64Kind:
		goType = "uint64"
	case protoreflect.DoubleKind:
		goType = "float64"
	case protoreflect.StringKind:
		goType = "string"
	case protoreflect.BytesKind:
		goType = "[]byte"
	default:
		return "", fmt.Errorf("unknown kind %s", descriptor.Kind())
	}

	return goType, nil
}
