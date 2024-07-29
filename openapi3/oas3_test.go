package openapi3_test

import (
	"reflect"
	"testing"

	"github.com/iwanhae/resource/openapi3"
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

	Users []StructCommon `json:"users"`
	User  StructCommon   `json:"user"`
}

type StructCommon struct {
	UUID string `json:"uuid"`
}

func TestSchema(t *testing.T) {
	b := openapi3.NewBuilder()

	b.Register(reflect.TypeFor[Category]())
	b.Register(reflect.TypeFor[Order]())
	b.Register(reflect.TypeFor[Pet]())
	b.Register(reflect.TypeFor[Tag]())
	b.Register(reflect.TypeFor[User]())

	result := b.Build()
	s := result.Components.Schemas

	testCases := []struct {
		name   string
		got    interface{}
		expect interface{}
	}{
		{
			name:   "category should be object",
			got:    s["category"].Value.Type.Is("object"),
			expect: true,
		},
		{
			name:   "category.id should be integer",
			got:    s["category"].Value.Properties["id"].Value.Type.Is("integer"),
			expect: true,
		},
		{
			name:   "category.name should be integer",
			got:    s["category"].Value.Properties["name"].Value.Type.Is("string"),
			expect: true,
		},
		{
			name:   "category does not has required field",
			got:    len(s["category"].Value.Required),
			expect: 0,
		},
	}
	for _, tc := range testCases {
		if !reflect.DeepEqual(tc.got, tc.expect) {
			t.Errorf("%s: Expect %q but got %q", tc.name, tc.expect, tc.got)
		}
	}
}
