package extra

import (
	"reflect"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"github.com/stretchr/testify/assert"
)

func TestEmitEmptyExtension(t *testing.T) {
	type Message struct {
		A string `json:"a,omitempty"`
		B int    `json:"b,omitempty"`
		C *int   `json:"c,omitempty"`
	}
	m := Message{}

	cfg := jsoniter.Config{}.Froze()

	jsn, err := cfg.MarshalToString(m)
	assert.Nil(t, err)
	assert.Equal(t, `{}`, jsn)

	m.A = "strA"
	jsn, err = cfg.MarshalToString(m)
	assert.Nil(t, err)
	assert.Equal(t, `{"a":"strA"}`, jsn)

	cfg = jsoniter.Config{}.Froze()
	cfg.RegisterExtension(&EmitEmptyExtension{})

	jsn, err = cfg.MarshalToString(m)
	assert.Nil(t, err)
	assert.Equal(t, `{"a":"strA","b":0,"c":null}`, jsn)

	m.A = ""
	jsn, err = cfg.MarshalToString(m)
	assert.Nil(t, err)
	assert.Equal(t, `{"a":"","b":0,"c":null}`, jsn)

	cfg = jsoniter.Config{}.Froze()
	cfg.RegisterExtension(&EmitEmptyExtension{Filter: func(binding *jsoniter.Binding) bool {
		typ := binding.Field.Type()
		return typ.Kind() == reflect.Int || (typ.Kind() == reflect.Ptr && typ.(reflect2.PtrType).Elem().Kind() == reflect.Int)
	}})
	jsn, err = cfg.MarshalToString(m)
	assert.Nil(t, err)
	assert.Equal(t, `{"b":0,"c":null}`, jsn)

	cfg = jsoniter.Config{}.Froze()
	cfg.RegisterExtension(&EmitEmptyExtension{Filter: func(binding *jsoniter.Binding) bool {
		return binding.Field.Name() == "C"
	}})
	m.A = "strA"
	jsn, err = cfg.MarshalToString(m)
	assert.Nil(t, err)
	assert.Equal(t, `{"a":"strA","c":null}`, jsn)

	// embedded
	cfg = jsoniter.Config{}.Froze()
	cfg.RegisterExtension(&EmitEmptyExtension{})

	type Embedded struct {
		FirstName string `json:"firstName"`
		LaseName  string `json:"laseName,omitempty"`
	}
	type OutA struct {
		Age int `json:"age,omitempty"`
		*Embedded
	}
	a := &OutA{}
	jsn, err = cfg.MarshalToString(a)
	assert.Nil(t, err)
	assert.Equal(t, `{"age":0}`, jsn)

	type OutB struct {
		Age int `json:"age,omitempty"`
		Embedded
	}
	b := &OutB{}
	jsn, err = cfg.MarshalToString(b)
	assert.Nil(t, err)
	assert.Equal(t, `{"age":0,"firstName":"","laseName":""}`, jsn)

}
