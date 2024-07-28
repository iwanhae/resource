package openapi3

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func RegisterStructToSchemas(t reflect.Type, s openapi3.Schemas) string {
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
			if jsonTag != "" && jsonTag != "-" {
				fieldType, isPtr := derefType(field.Type)
				if !isPtr {
					schema.Required = append(schema.Required, jsonTag)
				}
				switch fieldType.Kind() {
				case reflect.Bool:
					t := openapi3.NewBoolSchema()
					schema.Properties[jsonTag] = openapi3.NewSchemaRef("", t)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
					reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					t := openapi3.NewIntegerSchema()
					schema.Properties[jsonTag] = openapi3.NewSchemaRef("", t)
				case reflect.Float32, reflect.Float64:
					t := openapi3.NewFloat64Schema()
					schema.Properties[jsonTag] = openapi3.NewSchemaRef("", t)
				case reflect.String:
					t := openapi3.NewStringSchema()
					schema.Properties[jsonTag] = openapi3.NewSchemaRef("", t)
				case reflect.Struct:
					ref := RegisterStructToSchemas(fieldType, s)
					schema.Properties[jsonTag] = openapi3.NewSchemaRef(ref, nil)
				default:
					break
				}
			}
		}
	}

	s[name] = openapi3.NewSchemaRef("", schema)
	return fmt.Sprintf("#/components/schemas/%s", name)
}

func derefType(t reflect.Type) (deref reflect.Type, isPtr bool) {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
		isPtr = true
	}
	return t, isPtr
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
