package jsoniterpb

import (
	"reflect"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// https://github.com/protocolbuffers/protobuf-go/blob/master/encoding/protojson/well_known_types.go
var WellKnownTypes = map[reflect2.Type]bool{
	reflect2.TypeOfPtr((*anypb.Any)(nil)).Elem():              true,
	reflect2.TypeOfPtr((*timestamppb.Timestamp)(nil)).Elem():  true,
	reflect2.TypeOfPtr((*durationpb.Duration)(nil)).Elem():    true,
	reflect2.TypeOfPtr((*wrapperspb.BoolValue)(nil)).Elem():   true,
	reflect2.TypeOfPtr((*wrapperspb.Int32Value)(nil)).Elem():  true,
	reflect2.TypeOfPtr((*wrapperspb.Int64Value)(nil)).Elem():  true,
	reflect2.TypeOfPtr((*wrapperspb.UInt32Value)(nil)).Elem(): true,
	reflect2.TypeOfPtr((*wrapperspb.UInt64Value)(nil)).Elem(): true,
	reflect2.TypeOfPtr((*wrapperspb.FloatValue)(nil)).Elem():  true,
	reflect2.TypeOfPtr((*wrapperspb.DoubleValue)(nil)).Elem(): true,
	reflect2.TypeOfPtr((*wrapperspb.StringValue)(nil)).Elem(): true,
	reflect2.TypeOfPtr((*wrapperspb.BytesValue)(nil)).Elem():  true,
	reflect2.TypeOfPtr((*structpb.Struct)(nil)).Elem():        true,
	reflect2.TypeOfPtr((*structpb.ListValue)(nil)).Elem():     true,
	reflect2.TypeOfPtr((*structpb.Value)(nil)).Elem():         true,
	reflect2.TypeOfPtr((*fieldmaskpb.FieldMask)(nil)).Elem():  true,
	reflect2.TypeOfPtr((*emptypb.Empty)(nil)).Elem():          true,
}

func IsWellKnownType(typ reflect2.Type) bool {
	for typ.Kind() == reflect.Ptr {
		typ = typ.(reflect2.PtrType).Elem()
	}
	return WellKnownTypes[typ]
}

var WktProtoCodecs = map[reflect2.Type]*ProtoCodec{
	reflect2.TypeOfPtr((*anypb.Any)(nil)).Elem(): wktAnyCodec,

	reflect2.TypeOfPtr((*timestamppb.Timestamp)(nil)).Elem(): wktTimestampCodec,
	reflect2.TypeOfPtr((*durationpb.Duration)(nil)).Elem():   wktDurationCodec,

	reflect2.TypeOfPtr((*wrapperspb.BoolValue)(nil)).Elem(): (&ProtoCodec{}).
		SetElemEncodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, stream *jsoniter.Stream) {
			stream.WriteVal(((*wrapperspb.BoolValue)(ptr)).GetValue())
		}).
		SetElemDecodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			iter.ReadVal(&((*wrapperspb.BoolValue)(ptr).Value))
		}),
	reflect2.TypeOfPtr((*wrapperspb.Int32Value)(nil)).Elem(): (&ProtoCodec{}).
		SetElemEncodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, stream *jsoniter.Stream) {
			stream.WriteVal(((*wrapperspb.Int32Value)(ptr)).GetValue())
		}).
		SetElemDecodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			iter.ReadVal(&((*wrapperspb.Int32Value)(ptr).Value))
		}),
	reflect2.TypeOfPtr((*wrapperspb.Int64Value)(nil)).Elem(): (&ProtoCodec{}).
		SetElemEncodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, stream *jsoniter.Stream) {
			stream.WriteVal(((*wrapperspb.Int64Value)(ptr)).GetValue())
		}).
		SetElemDecodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			iter.ReadVal(&((*wrapperspb.Int64Value)(ptr).Value))
		}),
	reflect2.TypeOfPtr((*wrapperspb.UInt32Value)(nil)).Elem(): (&ProtoCodec{}).
		SetElemEncodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, stream *jsoniter.Stream) {
			stream.WriteVal(((*wrapperspb.UInt32Value)(ptr)).GetValue())
		}).
		SetElemDecodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			iter.ReadVal(&((*wrapperspb.UInt32Value)(ptr).Value))
		}),
	reflect2.TypeOfPtr((*wrapperspb.UInt64Value)(nil)).Elem(): (&ProtoCodec{}).
		SetElemEncodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, stream *jsoniter.Stream) {
			stream.WriteVal(((*wrapperspb.UInt64Value)(ptr)).GetValue())
		}).
		SetElemDecodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			iter.ReadVal(&((*wrapperspb.UInt64Value)(ptr).Value))
		}),
	reflect2.TypeOfPtr((*wrapperspb.FloatValue)(nil)).Elem(): (&ProtoCodec{}).
		SetElemEncodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, stream *jsoniter.Stream) {
			stream.WriteVal(((*wrapperspb.FloatValue)(ptr)).GetValue())
		}).
		SetElemDecodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			iter.ReadVal(&((*wrapperspb.FloatValue)(ptr).Value))
		}),
	reflect2.TypeOfPtr((*wrapperspb.DoubleValue)(nil)).Elem(): (&ProtoCodec{}).
		SetElemEncodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, stream *jsoniter.Stream) {
			stream.WriteVal(((*wrapperspb.DoubleValue)(ptr)).GetValue())
		}).
		SetElemDecodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			iter.ReadVal(&((*wrapperspb.DoubleValue)(ptr).Value))
		}),
	reflect2.TypeOfPtr((*wrapperspb.StringValue)(nil)).Elem(): (&ProtoCodec{}).
		SetElemEncodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, stream *jsoniter.Stream) {
			stream.WriteVal(((*wrapperspb.StringValue)(ptr)).GetValue())
		}).
		SetElemDecodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			iter.ReadVal(&((*wrapperspb.StringValue)(ptr).Value))
		}),
	reflect2.TypeOfPtr((*wrapperspb.BytesValue)(nil)).Elem(): (&ProtoCodec{}).
		SetElemEncodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, stream *jsoniter.Stream) {
			stream.WriteVal(((*wrapperspb.BytesValue)(ptr)).GetValue())
		}).
		SetElemDecodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			iter.ReadVal(&((*wrapperspb.BytesValue)(ptr).Value))
		}),

	// Because the following three implement json.Marshaler/Unmarshaler, we must also set the codec of its pointer type to override
	reflect2.TypeOfPtr((*structpb.Struct)(nil)).Elem():    wktStructCodec,
	reflect2.TypeOfPtr((*structpb.Struct)(nil)):           wktStructCodec,
	reflect2.TypeOfPtr((*structpb.ListValue)(nil)).Elem(): wktListValueCodec,
	reflect2.TypeOfPtr((*structpb.ListValue)(nil)):        wktListValueCodec,
	reflect2.TypeOfPtr((*structpb.Value)(nil)).Elem():     wktValueCodec,
	reflect2.TypeOfPtr((*structpb.Value)(nil)):            wktValueCodec,

	reflect2.TypeOfPtr((*fieldmaskpb.FieldMask)(nil)).Elem(): wktFieldmaskCodec,
	// reflect2.TypeOfPtr((*emptypb.Empty)(nil)).Elem(): // No special handling required
}

func init() {
	for k, v := range WktProtoCodecs {
		ProtoCodecs[k] = v
	}
}
