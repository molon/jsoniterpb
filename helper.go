package jsoniterpb

import (
	"fmt"
	"reflect"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

type stringModeNumberEncoder struct {
	elemEncoder jsoniter.ValEncoder
}

var singleQuote = []byte{'"'}

func (encoder *stringModeNumberEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.Write(singleQuote)
	encoder.elemEncoder.Encode(ptr, stream)
	stream.Write(singleQuote)
}

func (encoder *stringModeNumberEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.elemEncoder.IsEmpty(ptr)
}

type stringModeNumberDecoder struct {
	elemDecoder jsoniter.ValDecoder
}

func (decoder *stringModeNumberDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	if iter.WhatIsNext() == jsoniter.StringValue {
		iter.NextToken()
		decoder.elemDecoder.Decode(ptr, iter)
		iter.NextToken()
		return
	}
	decoder.elemDecoder.Decode(ptr, iter)
}

type funcEncoder struct {
	fun         jsoniter.EncoderFunc
	isEmptyFunc func(ptr unsafe.Pointer) bool
}

func (encoder *funcEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	encoder.fun(ptr, stream)
}

func (encoder *funcEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	if encoder.isEmptyFunc == nil {
		return false
	}
	return encoder.isEmptyFunc(ptr)
}

type funcDecoder struct {
	fun jsoniter.DecoderFunc
}

func (decoder *funcDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	decoder.fun(ptr, iter)
}

type dynamicEncoder struct {
	valType reflect2.Type
}

func (encoder *dynamicEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	obj := encoder.valType.UnsafeIndirect(ptr)
	stream.WriteVal(obj)
}

func (encoder *dynamicEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.valType.UnsafeIndirect(ptr) == nil
}

type OptionalEncoder struct {
	ValueEncoder jsoniter.ValEncoder
	IfNil        func(stream *jsoniter.Stream)
}

func (encoder *OptionalEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	if *((*unsafe.Pointer)(ptr)) == nil {
		if encoder.IfNil != nil {
			encoder.IfNil(stream)
		} else {
			stream.WriteNil()
		}
	} else {
		encoder.ValueEncoder.Encode(*((*unsafe.Pointer)(ptr)), stream)
	}
}

func (encoder *OptionalEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*unsafe.Pointer)(ptr)) == nil
}

type OptionalDecoder struct {
	ValueType    reflect2.Type
	ValueDecoder jsoniter.ValDecoder
	IfNil        func(ptr unsafe.Pointer)
}

func (decoder *OptionalDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	if iter.ReadNil() {
		if decoder.IfNil != nil {
			decoder.IfNil(ptr)
			return
		}
		*((*unsafe.Pointer)(ptr)) = nil
	} else {
		if *((*unsafe.Pointer)(ptr)) == nil {
			//pointer to null, we have to allocate memory to hold the value
			newPtr := decoder.ValueType.UnsafeNew()
			decoder.ValueDecoder.Decode(newPtr, iter)
			*((*unsafe.Pointer)(ptr)) = newPtr
		} else {
			//reuse existing instance
			decoder.ValueDecoder.Decode(*((*unsafe.Pointer)(ptr)), iter)
		}
	}
}

func WrapElemEncoder(typ reflect2.Type, enc jsoniter.ValEncoder, ifNil func(stream *jsoniter.Stream)) jsoniter.ValEncoder {
	if typ.Kind() == reflect.Ptr {
		if typ.(reflect2.PtrType).Elem().Kind() == reflect.Struct {
			return &OptionalEncoder{
				ValueEncoder: enc,
				IfNil:        ifNil,
			}
		}
		panic(fmt.Sprintf("WrapElemEncoderWithIfNil does not support type %v", typ))
	}
	return enc
}

func WrapElemDecoder(typ reflect2.Type, dec jsoniter.ValDecoder, ifNil func(ptr unsafe.Pointer)) jsoniter.ValDecoder {
	if typ.Kind() == reflect.Ptr {
		elemType := typ.(reflect2.PtrType).Elem()
		if elemType.Kind() == reflect.Struct {
			return &OptionalDecoder{
				ValueType:    elemType,
				ValueDecoder: dec,
				IfNil:        ifNil,
			}
		}
		panic(fmt.Sprintf("WrapElemDecoderIfNil does not support type %v", typ))
	}
	return dec
}
