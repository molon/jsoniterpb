package protojson_test

import (
	"bytes"
	"encoding/json"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	testv1 "github.com/molon/jsoniterpb/internal/gen/go/test/v1"
	pb3 "github.com/molon/jsoniterpb/internal/protojson/textpb3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func pMarshalToStringWithOpts(opts protojson.MarshalOptions, m proto.Message) (string, error) {
	by, err := opts.Marshal(m)
	if err != nil {
		return "", err
	}
	// https://github.com/golang/protobuf/issues/1121
	var out bytes.Buffer
	err = json.Compact(&out, by)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func pMarshalToString(m proto.Message) (string, error) {
	return pMarshalToStringWithOpts(protojson.MarshalOptions{}, m)
}

func pUnmarshalFromString(s string, m proto.Message) error {
	return protojson.Unmarshal([]byte(s), m)
}

// https://github.com/golang/protobuf/issues/1121
func TestPjCompactIssue(t *testing.T) {
	m := &testv1.WKTs{}
	bb, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(m)
	assert.Nil(t, err)

	var out bytes.Buffer
	err = json.Compact(&out, bb)
	assert.Nil(t, err)
	// assert.Equal(t, out.String(), string(bb)) // maybe false sometimes
}

// https://github.com/golang/protobuf/issues/1487
// https://github.com/golang/protobuf/issues/1258
// https://github.com/golang/protobuf/issues/1313
func TestPjNilValue(t *testing.T) {
	euOpts := protojson.MarshalOptions{EmitUnpopulated: true}
	var jsn string
	var err error
	var m proto.Message

	// marshal nil proto.Message to zero value if it is root (// TIPS: jsoniterpb does not support this feature)
	jsn, err = pMarshalToStringWithOpts(euOpts, (*wrapperspb.Int32Value)(nil))
	assert.Nil(t, err)
	assert.Equal(t, `0`, jsn)
	jsn, err = pMarshalToStringWithOpts(euOpts, (*testv1.Message)(nil))
	assert.Nil(t, err)
	assert.Equal(t, `{"id":""}`, jsn)

	// but marshal nil to null if it is not root, just like standard encoding/json
	m = &testv1.WKTs{}
	jsn, err = pMarshalToStringWithOpts(euOpts, m)
	assert.Nil(t, err)
	assert.Equal(t, `{"a":null,"d":null,"t":null,"st":null,"i32":null,"ui32":null,"i64":null,"u64":null,"f32":null,"f64":null,"b":null,"s":null,"by":null,"fm":null,"em":null,"nu":null,"v":null,"lv":null}`, jsn)

	// marshal nil list/map to zero value
	jsn, err = pMarshalToStringWithOpts(euOpts, &testv1.RepeatedWKTs{})
	assert.Nil(t, err)
	assert.Equal(t, `{"a":[],"d":[],"t":[],"st":[],"i32":[],"ui32":[],"i64":[],"u64":[],"f32":[],"f64":[],"b":[],"s":[],"by":[],"fm":[],"em":[],"nu":[],"v":[],"lv":[]}`, jsn)
	jsn, err = pMarshalToStringWithOpts(euOpts, &testv1.Map{})
	assert.Nil(t, err)
	assert.Equal(t, `{"en":{},"msg":{},"str":{},"by":{},"bo":{},"an":{},"bn":{}}`, jsn)

	// marshal nil []byte to ""
	m = &testv1.Singular{
		By: []byte(nil),
	}
	jsn, err = pMarshalToStringWithOpts(euOpts, m)
	assert.Nil(t, err)
	assert.Equal(t, `{"e":"JSON_ENUM_UNSPECIFIED","s":"","i32":0,"i64":"0","u32":0,"u64":"0","f32":0,"f64":0,"si32":0,"si64":"0","fi32":0,"fi64":"0","sfi32":0,"sfi64":"0","bl":false,"by":"","msg":null}`, jsn)

	// - marshal elem in (list or map) to zero value
	// - sorts false before true, numeric keys in ascending order, and strings in lexicographical ordering according to UTF-8 codepoints.
	m = &testv1.All{
		RWkt: &testv1.RepeatedWKTs{
			D: []*durationpb.Duration{nil},
		},
		M: &testv1.Map{
			Msg: map[int32]*testv1.Nested{1: nil},
			By:  map[bool][]byte{false: nil},
			Bn: map[uint64]*wrapperspb.UInt64Value{
				181818: wrapperspb.UInt64(123),
				2:      nil,
				181817: wrapperspb.UInt64(223),
			},
		},
	}
	jsn, err = pMarshalToString(m)
	assert.Nil(t, err)
	assert.Equal(t, `{"rWkt":{"d":["0s"]},"m":{"msg":{"1":{}},"by":{"false":""},"bn":{"2":"0","181817":"223","181818":"123"}}}`, jsn)
}

// https://github.com/golang/protobuf/issues/1487#issuecomment-1251803426
// https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#google.protobuf.Value
// A producer of value is expected to set one of that variants, absence of any variant indicates an error.
func TestPjWktValue(t *testing.T) {
	euOpts := protojson.MarshalOptions{EmitUnpopulated: true}
	var jsn string
	var err error
	var m, m2 proto.Message

	// avoid the zero value of *structpb.Value at any time, otherwise you will get an error
	jsn, err = pMarshalToStringWithOpts(euOpts, &testv1.RepeatedWKTs{
		V: []*structpb.Value{
			structpb.NewNullValue(),
			structpb.NewBoolValue(false),
			nil, // error
			structpb.NewStringValue("str"),
		},
	})
	assert.Contains(t, err.Error(), "google.protobuf.Value: none of the oneof fields is set")
	jsn, err = pMarshalToStringWithOpts(euOpts, (*structpb.Value)(nil))
	assert.Contains(t, err.Error(), "google.protobuf.Value: none of the oneof fields is set")

	// unmarshal null to *structpb.Value, will results `structpb.NewNullValue()`
	// TIPS: In this case, m and m2 will not be equal, so try to avoid using (*structpb.Value)(nil)
	m = &testv1.WKTs{
		V: nil,
	}
	jsn, err = pMarshalToStringWithOpts(euOpts, m)
	assert.Nil(t, err)
	assert.Equal(t, `{"a":null,"d":null,"t":null,"st":null,"i32":null,"ui32":null,"i64":null,"u64":null,"f32":null,"f64":null,"b":null,"s":null,"by":null,"fm":null,"em":null,"nu":null,"v":null,"lv":null}`, jsn)
	m2 = &testv1.WKTs{}
	err = pUnmarshalFromString(jsn, m2)
	assert.Nil(t, err)
	assert.True(t, proto.Equal(structpb.NewNullValue(), m2.(*testv1.WKTs).V))
	assert.Nil(t, m.(*testv1.WKTs).V)
	assert.NotNil(t, m2.(*testv1.WKTs).V)
	assert.False(t, cmp.Diff(m, m2, protocmp.Transform()) == "")
}

// marshal nan inf float to string
func TestPjInfNaN(t *testing.T) {
	jsn, err := pMarshalToString(&testv1.Singular{
		F32: float32(math.NaN()),
		F64: math.NaN(),
	})
	assert.Nil(t, err)
	assert.Equal(t, `{"f32":"NaN","f64":"NaN"}`, jsn)

	jsn, err = pMarshalToString(&testv1.Singular{
		F32: float32(math.Inf(+1)),
		F64: math.Inf(-1),
	})
	assert.Nil(t, err)
	assert.Equal(t, `{"f32":"Infinity","f64":"-Infinity"}`, jsn)
}

// marshal bit 64 to string
func TestMarshal64BitInteger(t *testing.T) {
	var jsn string
	var err error

	m := &testv1.Singular{I64: 123, U64: 234}
	jsn, err = pMarshalToString(m)
	assert.Nil(t, err)
	assert.Equal(t, `{"i64":"123","u64":"234"}`, jsn)
}

// fuzzy decode num
func TestFuzzyDecode(t *testing.T) {
	m := &pb3.Scalars{}
	err := pUnmarshalFromString(`{
		"sInt32": 1234,
		"sInt64": -1234,
		"sUint32": 1e2,
		"sUint64": 100E-2,
		"sSint32": 1.0,
		"sSint64": -1.0,
		"sFixed32": 1.234e+5,
		"sFixed64": 1200E-2,
		"sSfixed32": -1.234e05,
		"sSfixed64": -1200e-02,
		"sDouble": "123"
	  }`, m)
	assert.Nil(t, err)
	assert.True(t, proto.Equal(m, &pb3.Scalars{
		SInt32:    1234,
		SInt64:    -1234,
		SUint32:   100,
		SUint64:   1,
		SSint32:   1,
		SSint64:   -1,
		SFixed32:  123400,
		SFixed64:  12,
		SSfixed32: -123400,
		SSfixed64: -12,
		SDouble:   123,
	}))
}

func TestEmitUnpopulatedWithOptional(t *testing.T) {
	var jsn string
	var err error

	m := &pb3.Proto3Optional{
		OptInt64: proto.Int64(0),
		OptInt32: proto.Int32(1),
	}
	// if opt is not nil, means it is not unpopulated although zero value
	jsn, err = pMarshalToStringWithOpts(protojson.MarshalOptions{EmitUnpopulated: false}, m)
	assert.Nil(t, err)
	assert.Equal(t, `{"optInt32":1,"optInt64":"0"}`, jsn)

	// dont marshal nil optional although EmitUnpopulated:true
	jsn, err = pMarshalToStringWithOpts(protojson.MarshalOptions{EmitUnpopulated: true}, m)
	assert.Nil(t, err)
	assert.Equal(t, `{"optInt32":1,"optInt64":"0"}`, jsn)
}

func TestInvalidUTF8(t *testing.T) {
	var m proto.Message
	var jsn string
	var err error

	m = &testv1.Singular{S: "\xff"}

	// can not marshal invalid utf8
	jsn, err = pMarshalToString(m)
	assert.Contains(t, err.Error(), "invalid UTF-8")

	// can not unmarshal invalid utf8
	jsn = `{"s":"` + "abc\xff" + `"}`
	err = pUnmarshalFromString(jsn, m)
	assert.Contains(t, err.Error(), "invalid UTF-8")
}

// proto.Equal cant handle any.Any which contains map
// https://github.com/golang/protobuf/issues/455
// reason => https://github.com/golang/protobuf/commit/efcaa340c1a788c79e1ca31217d66aa41c405a51
func TestProtoEqualIssue(t *testing.T) {
	var m, m2 proto.Message
	var jsn string
	var err error

	s, _ := structpb.NewStruct(map[string]interface{}{
		"keyA": "valueA",
		"keyB": "valueB",
		"keyC": "valueC",
	})
	a, _ := anypb.New(s)
	m = &testv1.WKTs{A: a}

	jsn, err = pMarshalToString(m)
	assert.Nil(t, err)
	m2 = &testv1.WKTs{}
	err = pUnmarshalFromString(jsn, m2)
	assert.Nil(t, err)
	// assert.True(t, proto.Equal(m, m2)) // maybe false sometimes
	assert.Equal(t, "", cmp.Diff(m, m2, protocmp.Transform(), cmpopts.EquateNaNs()))
}
