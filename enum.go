package jsoniterpb

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	NullValue_enum_fullname = "google.protobuf.NullValue"
)

var (
	nullValuePtrType = reflect2.TypeOfPtr((*structpb.NullValue)(nil))
	protoEnumType    = reflect2.TypeOfPtr((*protoreflect.Enum)(nil)).Elem()
)

func (e *ProtoExtension) createProtoEnumEncoder(typ reflect2.Type) (xret jsoniter.ValEncoder) {
	if !e.UseEnumNumbers {
		if typ.Implements(protoEnumType) && typ.Kind() != reflect.Ptr {
			return &protoEnumEncoder{
				valueType: typ,
			}
		}
	}
	return nil
}

func (e *ProtoExtension) createProtoEnumDecoder(typ reflect2.Type) (xret jsoniter.ValDecoder) {
	// we want fuzzy decode, so does not need to check e.UseEnumNumbers
	if typ.Implements(protoEnumType) {
		if typ.Kind() != reflect.Ptr {
			return &protoEnumDecoder{
				valueType: typ,
			}
		}

		if typ == nullValuePtrType {
			return &funcDecoder{
				fun: func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
					if iter.ReadNil() {
						v := structpb.NullValue_NULL_VALUE
						*((**structpb.NullValue)(ptr)) = &v
					} else {
						iter.ReportError("protobuf", fmt.Sprintf("%v: invalid value %v", NullValue_enum_fullname, string(iter.SkipAndReturnBytes())))
					}
				},
			}
		}
	}
	return nil
}

type protoEnumEncoder struct {
	valueType reflect2.Type
	once      sync.Once
	enumDesc  protoreflect.EnumDescriptor
}

func (enc *protoEnumEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	x, ok := enc.valueType.UnsafeIndirect(ptr).(protoreflect.Enum)
	if !ok {
		stream.WriteVal(protoreflect.EnumNumber(0))
		return
	}
	enc.once.Do(func() {
		enc.enumDesc = x.Descriptor()
	})
	if enc.enumDesc.FullName() == NullValue_enum_fullname {
		stream.WriteNil()
		return
	}
	n := x.Number()
	ev := enc.enumDesc.Values().ByNumber(n)
	if ev != nil {
		stream.WriteVal(string(ev.Name()))
	} else {
		stream.WriteVal(n)
	}
}

func (enc *protoEnumEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*protoreflect.EnumNumber)(ptr)) == 0
}

type protoEnumDecoder struct {
	valueType    reflect2.Type
	once         sync.Once
	enumValDescs protoreflect.EnumValueDescriptors
}

func (dec *protoEnumDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	switch valueType {
	case jsoniter.NumberValue:
		num := iter.ReadInt32()
		*((*protoreflect.EnumNumber)(ptr)) = protoreflect.EnumNumber(num)
	case jsoniter.StringValue:
		var name string
		iter.ReadVal(&name)
		dec.once.Do(func() {
			x := dec.valueType.UnsafeIndirect(ptr).(protoreflect.Enum)
			dec.enumValDescs = x.Descriptor().Values()
		})
		ev := dec.enumValDescs.ByName(protoreflect.Name(name))
		if ev != nil {
			*((*protoreflect.EnumNumber)(ptr)) = ev.Number()
		} else {
			// is "num"?
			num, err := strconv.ParseInt(name, 10, 32)
			if err == nil {
				*((*protoreflect.EnumNumber)(ptr)) = protoreflect.EnumNumber(num)
			} else {
				iter.ReportError("protobuf", fmt.Sprintf(
					"error decode from string for type %s",
					dec.valueType,
				))
			}
		}
	case jsoniter.NilValue:
		iter.Skip()
		*((*protoreflect.EnumNumber)(ptr)) = 0
	default:
		iter.ReportError("protobuf", fmt.Sprintf(
			"error decode for type %s",
			dec.valueType,
		))
	}
}
