package openapi3_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	oas3 "github.com/iwanhae/resource/openapi3"
)

type StructSimple struct {
	Name     string  `json:"name"`
	ID       int64   `json:"id"`
	IsGood   bool    `json:"is_good"`
	Height   float32 `json:"height"`
	Optional *string `json:"optional"`

	Test struct {
		Nested string `json:"nested"`
	} `json:"hello"`

	Users StructCommon `json:"users"`
}

type StructCommon struct {
	UUID string `json:"uuid"`
}

func TestSchema(t *testing.T) {
	schemas := make(openapi3.Schemas)
	oas3.RegisterStructToSchemas(reflect.TypeOf(StructSimple{}), schemas)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(schemas)
}
