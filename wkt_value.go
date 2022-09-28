package jsoniterpb

import (
	"errors"
	"fmt"
	"math"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
)

var wktValueType = reflect2.TypeOfPtr((*structpb.Value)(nil)).Elem()

const (
	Value_message_fullname           protoreflect.FullName = "google.protobuf.Value"
	Value_NumberValue_field_fullname protoreflect.FullName = "google.protobuf.Value.number_value"
)

var wktValueCodec = (&ProtoCodec{}).
	SetElemEncodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, stream *jsoniter.Stream) {
		x := ((*structpb.Value)(ptr))
		err := marshalWktValue(x, stream)
		if err != nil {
			stream.Error = fmt.Errorf("%s: %w", Value_message_fullname, err)
			return
		}
	}).
	SetElemDecodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		err := unmarshalWktValue(((*structpb.Value)(ptr)), iter)
		if err != nil {
			iter.ReportError("protobuf", fmt.Sprintf("%s: %v", Value_message_fullname, err))
		}
	})

func marshalWktValue(x *structpb.Value, stream *jsoniter.Stream) error {
	switch v := x.GetKind().(type) {
	case *structpb.Value_NullValue:
		if v != nil {
			stream.WriteNil()
			return nil
		}
	case *structpb.Value_NumberValue:
		if v != nil {
			if math.IsNaN(v.NumberValue) || math.IsInf(v.NumberValue, 0) {
				return fmt.Errorf("%s: invalid %v value", Value_NumberValue_field_fullname, v)
			}
			stream.WriteVal(v.NumberValue)
			return nil
		}
	case *structpb.Value_StringValue:
		if v != nil {
			stream.WriteVal(v.StringValue)
			return nil
		}
	case *structpb.Value_BoolValue:
		if v != nil {
			stream.WriteVal(v.BoolValue)
			return nil
		}
	case *structpb.Value_StructValue:
		if v != nil {
			if v.StructValue == nil {
				stream.WriteEmptyObject()
				return nil
			}
			stream.WriteVal(v.StructValue)
			return nil
		}
	case *structpb.Value_ListValue:
		if v != nil {
			if v.ListValue == nil {
				stream.WriteEmptyArray()
				return nil
			}
			stream.WriteVal(v.ListValue)
			return nil
		}
	}
	// https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#google.protobuf.Value
	// A producer of value is expected to set one of that variants, absence of any variant indicates an error.
	return errors.New("none of the oneof fields is set")
}

func unmarshalWktValue(x *structpb.Value, iter *jsoniter.Iterator) error {
	valueType := iter.WhatIsNext()
	switch valueType {
	case jsoniter.NilValue:
		iter.ReadNil()
		x.Kind = &structpb.Value_NullValue{
			NullValue: structpb.NullValue_NULL_VALUE,
		}
	case jsoniter.BoolValue:
		var val bool
		iter.ReadVal(&val)
		x.Kind = &structpb.Value_BoolValue{
			BoolValue: val,
		}
	case jsoniter.NumberValue:
		var val float64
		iter.ReadVal(&val)
		x.Kind = &structpb.Value_NumberValue{
			NumberValue: val,
		}
	case jsoniter.StringValue:
		var str string
		iter.ReadVal(&str)
		x.Kind = &structpb.Value_StringValue{
			StringValue: str,
		}
	case jsoniter.ObjectValue:
		v := &structpb.Struct{}
		iter.ReadVal(v)
		x.Kind = &structpb.Value_StructValue{
			StructValue: v,
		}
	case jsoniter.ArrayValue:
		v := &structpb.ListValue{}
		iter.ReadVal(v)
		x.Kind = &structpb.Value_ListValue{
			ListValue: v,
		}
	default:
		return errors.New("not number or string or object")
	}
	return nil
}
