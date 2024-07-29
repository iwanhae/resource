package openapi3

import (
	"crypto/sha256"
	"encoding"
	"encoding/hex"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

var (
	textMarshalerType   = reflect.TypeFor[encoding.TextMarshaler]()
	textUnmarshalerType = reflect.TypeFor[encoding.TextUnmarshaler]()
)

func NewBuilder() *builder {
	return &builder{
		schemas: make(openapi3.Schemas),
	}
}

type builder struct {
	schemas openapi3.Schemas
}

func (b *builder) Build() openapi3.T {
	return openapi3.T{
		Components: &openapi3.Components{
			Schemas: b.schemas,
		},
	}
}

func (b *builder) Register(t reflect.Type) string {
	name := camelCase(t.Name())
	if name == "" {
		b := sha256.Sum256([]byte(t.String()))
		name = hex.EncodeToString(b[:8])
	}
	schema := openapi3.NewObjectSchema()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.IsExported() {
			jsonTag := field.Tag.Get("json")
			jsonTag = strings.Split(jsonTag, ",")[0]
			if jsonTag != "" && jsonTag != "-" {
				fieldType, isPtr := derefType(field.Type)
				if !isPtr {
					schema.Required = append(schema.Required, jsonTag)
				}
				ref := b.schemaRefFor(fieldType)
				schema.Properties[jsonTag] = ref
			}
		}
	}

	b.schemas[name] = openapi3.NewSchemaRef("", schema)
	return fmt.Sprintf("#/components/schemas/%s", name)
}

func (b *builder) schemaRefFor(t reflect.Type) *openapi3.SchemaRef {
	var schemaRef *openapi3.SchemaRef

	// if t implements encoding.TextMarshaler and encoding.TextUnmarshaler
	if t.Implements(textMarshalerType) && reflect.PointerTo(t).Implements(textUnmarshalerType) {
		schemaRef = openapi3.NewSchemaRef("", openapi3.NewStringSchema())
		if t.Name() == "Time" && t.PkgPath() == "time" {
			schemaRef.Value.Format = "date-time"
		}
	} else {
		switch t.Kind() {
		case reflect.Bool:
			schemaRef = openapi3.NewSchemaRef("", openapi3.NewBoolSchema())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			schemaRef = openapi3.NewSchemaRef("", openapi3.NewIntegerSchema())
		case reflect.Float32, reflect.Float64:
			schemaRef = openapi3.NewSchemaRef("", openapi3.NewFloat64Schema())
		case reflect.String:
			schemaRef = openapi3.NewSchemaRef("", openapi3.NewStringSchema())
		case reflect.Struct:
			ref := b.Register(t)
			schemaRef = openapi3.NewSchemaRef(ref, nil)
		case reflect.Array, reflect.Slice:
			schema := openapi3.NewArraySchema()
			itemType, _ := dearrType(t)
			ref := b.schemaRefFor(itemType)
			schema.Items = ref
			schemaRef = openapi3.NewSchemaRef("", schema)
		}
	}

	return schemaRef
}

func derefType(t reflect.Type) (deref reflect.Type, isPtr bool) {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
		isPtr = true
	}
	return t, isPtr
}
func dearrType(t reflect.Type) (item reflect.Type, isArr bool) {
	for t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		t = t.Elem()
		isArr = true
	}
	return t, isArr
}

func camelCase(s string) string {

	// Remove all characters that are not alphanumeric or spaces or underscores
	s = regexp.MustCompile("[^a-zA-Z0-9_ ]+").ReplaceAllString(s, "")

	// Replace all underscores with spaces
	s = strings.ReplaceAll(s, "_", " ")

	// Remove all spaces
	s = strings.ReplaceAll(s, " ", "")

	// Lowercase the first letter
	if len(s) > 0 {
		s = strings.ToLower(s[:1]) + s[1:]
	}

	return s
}
