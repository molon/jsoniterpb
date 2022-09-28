package jsoniterpb

import (
	"reflect"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func (e *ProtoExtension) updateMapEncoderConstructorForSortMapKeys(v *jsoniter.MapEncoderConstructor) {
	if e.SortMapKeysAsString {
		return
	}

	v.DecorateFunc = func(mapEncoder jsoniter.ValEncoder) jsoniter.ValEncoder {
		enc, ok := mapEncoder.(*jsoniter.SortKeysMapEncoder)
		if !ok {
			return mapEncoder
		}

		// protobuf-go GenericKeyOrder
		switch v.MapType.Key().Kind() {
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:
			enc.KeyLess = func(x, y string) bool {
				xi, _ := strconv.ParseInt(x, 10, 64)
				yi, _ := strconv.ParseInt(y, 10, 64)
				return xi < yi
			}
		case reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			enc.KeyLess = func(x, y string) bool {
				xi, _ := strconv.ParseUint(x, 10, 64)
				yi, _ := strconv.ParseUint(y, 10, 64)
				return xi < yi
			}
		}
		return enc
	}
}
