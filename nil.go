package jsoniterpb

import (
	"reflect"
	"sync"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"google.golang.org/protobuf/types/known/structpb"
)

// https://github.com/golang/protobuf/issues/1487

var doubleQuote = []byte{'"', '"'}

func (e *ProtoExtension) decorateEncoderForNilCollection(typ reflect2.Type, encoder jsoniter.ValEncoder) jsoniter.ValEncoder {
	// - marshal nil []byte to ""
	// - marshal nil slice to []
	// - marshal nil map to {}
	isList := typ.Kind() == reflect.Slice
	isLikeMap := typ.Kind() == reflect.Map
	if isList || isLikeMap {
		isBytes := typ.Kind() == reflect.Slice && typ.(reflect2.SliceType).Elem().Kind() == reflect.Uint8
		return &funcEncoder{
			fun: func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
				if *((*unsafe.Pointer)(ptr)) == nil {
					if isBytes {
						stream.Write(doubleQuote)
						return
					}
					if isList {
						stream.WriteEmptyArray()
						return
					}
					if isLikeMap {
						stream.WriteEmptyObject()
						return
					}
				}
				encoder.Encode(ptr, stream)
			},
			isEmptyFunc: func(ptr unsafe.Pointer) bool {
				return encoder.IsEmpty(ptr)
			},
		}
	}
	return nil
}

var wktValuePtrType = reflect2.TypeOfPtr((*structpb.Value)(nil))

func (e *ProtoExtension) decorateDecoderForNil(typ reflect2.Type, dec jsoniter.ValDecoder) jsoniter.ValDecoder {
	// - unmarshal null to NULL value
	if typ == wktValuePtrType {
		return &funcDecoder{
			fun: func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
				if iter.ReadNil() {
					v := structpb.NewNullValue()
					*((**structpb.Value)(ptr)) = v
				} else {
					dec.Decode(ptr, iter)
				}
			},
		}
	}
	return nil
}

type lazyValue struct {
	once    sync.Once
	ret     interface{}
	creator func() interface{}
}

func newLazyValue(creator func() interface{}) *lazyValue {
	return &lazyValue{
		creator: creator,
	}
}

func (v *lazyValue) Get() interface{} {
	v.once.Do(func() {
		v.ret = v.creator()
	})
	return v.ret
}

var lazyPtrWithZeroValueMap sync.Map

// - marshal []type{a,nil,c} to [a,zero,c]
// - marshal map[string]type to {"a":"valueA",b:zero,c:"valueC"}
func noNullElemEncoderForCollection(valueType reflect2.Type, encoder jsoniter.ValEncoder) jsoniter.ValEncoder {
	return &funcEncoder{
		fun: func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			if valueType.Kind() == reflect.Ptr {
				if *(*unsafe.Pointer)(ptr) == nil {
					v, _ := lazyPtrWithZeroValueMap.LoadOrStore(valueType, newLazyValue(func() interface{} {
						ptrType := valueType.(reflect2.PtrType)
						elemType := ptrType.Elem()
						elemPtr := elemType.UnsafeNew()

						// record first
						newPtr := ptrType.UnsafeNew()
						*(*unsafe.Pointer)(newPtr) = elemPtr

						for elemType.Kind() == reflect.Ptr {
							ptrType = elemType.(reflect2.PtrType)
							elemType = ptrType.Elem()
							newElemPtr := elemType.UnsafeNew()
							*(*unsafe.Pointer)(elemPtr) = newElemPtr
							elemPtr = newElemPtr
						}
						return newPtr
					}))
					ptr = v.(*lazyValue).Get().(unsafe.Pointer)
				}
			}
			encoder.Encode(ptr, stream)
		},
		isEmptyFunc: func(ptr unsafe.Pointer) bool {
			return encoder.IsEmpty(ptr)
		},
	}
}

func (e *ProtoExtension) updateMapEncoderConstructorForNonNull(v *jsoniter.MapEncoderConstructor) {
	v.ElemEncoder = noNullElemEncoderForCollection(v.MapType.Elem(), v.ElemEncoder)
}

func (e *ProtoExtension) updateSliceEncoderConstructorForNonNull(v *jsoniter.SliceEncoderConstructor) {
	v.ElemEncoder = noNullElemEncoderForCollection(v.SliceType.Elem(), v.ElemEncoder)
}

func (e *ProtoExtension) updateArrayEncoderConstructorForNonNull(v *jsoniter.ArrayEncoderConstructor) {
	v.ElemEncoder = noNullElemEncoderForCollection(v.ArrayType.Elem(), v.ElemEncoder)
}
