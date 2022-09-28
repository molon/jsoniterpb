package jsoniterpb

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

// copy from extra/fuzzy_decode, with some changes

const maxUint = ^uint(0)
const maxInt = int(maxUint >> 1)
const minInt = -maxInt - 1

var fuzzyDecorateScalarDecoders = map[reflect.Kind]func(e *ProtoExtension, dec jsoniter.ValDecoder) jsoniter.ValDecoder{
	reflect.String: func(e *ProtoExtension, dec jsoniter.ValDecoder) jsoniter.ValDecoder {
		return &funcDecoder{
			fun: func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
				valueType := iter.WhatIsNext()
				switch valueType {
				case jsoniter.NumberValue:
					var number json.Number
					iter.ReadVal(&number)
					*((*string)(ptr)) = string(number)
				case jsoniter.NilValue:
					iter.Skip()
					*((*string)(ptr)) = ""
				default:
					dec.Decode(ptr, iter)
				}
			},
		}
	},
	reflect.Float32: func(e *ProtoExtension, dec jsoniter.ValDecoder) jsoniter.ValDecoder {
		return &funcDecoder{
			fun: func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
				valueType := iter.WhatIsNext()
				var str string
				switch valueType {
				case jsoniter.NumberValue:
					dec.Decode(ptr, iter)
				case jsoniter.StringValue:
					str = iter.ReadString()
					switch str {
					case "true":
						str = "1"
					case "false":
						str = "0"
					}
					newIter := iter.Pool().BorrowIterator([]byte(str))
					newIter.Attachment = iter.Attachment
					defer iter.Pool().ReturnIterator(newIter)
					*((*float32)(ptr)) = newIter.ReadFloat32()
					if newIter.Error != nil && newIter.Error != io.EOF {
						iter.Error = newIter.Error
					}
				case jsoniter.BoolValue:
					// support bool to float32
					if iter.ReadBool() {
						*((*float32)(ptr)) = 1
					} else {
						*((*float32)(ptr)) = 0
					}
				case jsoniter.NilValue:
					iter.Skip()
					*((*float32)(ptr)) = 0
				default:
					iter.ReportError("fuzzyFloat32Decoder", "not number or string")
				}
			},
		}
	},
	reflect.Float64: func(e *ProtoExtension, dec jsoniter.ValDecoder) jsoniter.ValDecoder {
		return &funcDecoder{
			fun: func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
				valueType := iter.WhatIsNext()
				var str string
				switch valueType {
				case jsoniter.NumberValue:
					dec.Decode(ptr, iter)
				case jsoniter.StringValue:
					str = iter.ReadString()
					switch str {
					case "true":
						str = "1"
					case "false":
						str = "0"
					}
					newIter := iter.Pool().BorrowIterator([]byte(str))
					newIter.Attachment = iter.Attachment
					defer iter.Pool().ReturnIterator(newIter)
					*((*float64)(ptr)) = newIter.ReadFloat64()
					if newIter.Error != nil && newIter.Error != io.EOF {
						iter.Error = newIter.Error
					}
				case jsoniter.BoolValue:
					// support bool to float64
					if iter.ReadBool() {
						*((*float64)(ptr)) = 1
					} else {
						*((*float64)(ptr)) = 0
					}
				case jsoniter.NilValue:
					iter.Skip()
					*((*float64)(ptr)) = 0
				default:
					iter.ReportError("fuzzyFloat64Decoder", "not number or string")
				}
			},
		}
	},
	reflect.Bool: func(e *ProtoExtension, dec jsoniter.ValDecoder) jsoniter.ValDecoder {
		return &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			var bint int
			if isFloat {
				val := iter.ReadFloat64()
				if val != 1 && val != 0 {
					iter.ReportError("fuzzy decode bool", "exceed range")
					return
				}
				bint = int(val)
			} else {
				bint = iter.ReadInt()
			}
			switch bint {
			case 1:
				*((*bool)(ptr)) = true
			case 0:
				*((*bool)(ptr)) = false
			default:
				iter.ReportError("fuzzy decode bool", fmt.Sprintf("invalid bool(%d)", bint))
				return
			}
		}}
	},
	reflect.Int: func(e *ProtoExtension, dec jsoniter.ValDecoder) jsoniter.ValDecoder {
		return &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			if isFloat {
				val := iter.ReadFloat64()
				if _, frac := math.Modf(val); frac != 0 {
					iter.ReportError("fuzzyIntegerDecoder", "found frac")
					return
				}
				if val > float64(maxInt) || val < float64(minInt) {
					iter.ReportError("fuzzy decode int", "exceed range")
					return
				}
				*((*int)(ptr)) = int(val)
			} else {
				dec.Decode(ptr, iter)
			}
		}}
	},
	reflect.Int8: func(e *ProtoExtension, dec jsoniter.ValDecoder) jsoniter.ValDecoder {
		return &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			if isFloat {
				val := iter.ReadFloat64()
				if _, frac := math.Modf(val); frac != 0 {
					iter.ReportError("fuzzyIntegerDecoder", "found frac")
					return
				}
				if val > float64(math.MaxInt8) || val < float64(math.MinInt8) {
					iter.ReportError("fuzzy decode int8", "exceed range")
					return
				}
				*((*int8)(ptr)) = int8(val)
			} else {
				dec.Decode(ptr, iter)
			}
		}}
	},
	reflect.Int16: func(e *ProtoExtension, dec jsoniter.ValDecoder) jsoniter.ValDecoder {
		return &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			if isFloat {
				val := iter.ReadFloat64()
				if _, frac := math.Modf(val); frac != 0 {
					iter.ReportError("fuzzyIntegerDecoder", "found frac")
					return
				}
				if val > float64(math.MaxInt16) || val < float64(math.MinInt16) {
					iter.ReportError("fuzzy decode int16", "exceed range")
					return
				}
				*((*int16)(ptr)) = int16(val)
			} else {
				dec.Decode(ptr, iter)
			}
		}}
	},
	reflect.Int32: func(e *ProtoExtension, dec jsoniter.ValDecoder) jsoniter.ValDecoder {
		return &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			if isFloat {
				val := iter.ReadFloat64()
				if _, frac := math.Modf(val); frac != 0 {
					iter.ReportError("fuzzyIntegerDecoder", "found frac")
					return
				}
				if val > float64(math.MaxInt32) || val < float64(math.MinInt32) {
					iter.ReportError("fuzzy decode int32", "exceed range")
					return
				}
				*((*int32)(ptr)) = int32(val)
			} else {
				dec.Decode(ptr, iter)
			}
		}}
	},
	reflect.Int64: func(e *ProtoExtension, dec jsoniter.ValDecoder) jsoniter.ValDecoder {
		return &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			if isFloat {
				val := iter.ReadFloat64()
				if _, frac := math.Modf(val); frac != 0 {
					iter.ReportError("fuzzyIntegerDecoder", "found frac")
					return
				}
				if val > float64(math.MaxInt64) || val < float64(math.MinInt64) {
					iter.ReportError("fuzzy decode int64", "exceed range")
					return
				}
				*((*int64)(ptr)) = int64(val)
			} else {
				dec.Decode(ptr, iter)
			}
		}}
	},
	reflect.Uint: func(e *ProtoExtension, dec jsoniter.ValDecoder) jsoniter.ValDecoder {
		return &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			if isFloat {
				val := iter.ReadFloat64()
				if _, frac := math.Modf(val); frac != 0 {
					iter.ReportError("fuzzyIntegerDecoder", "found frac")
					return
				}
				if val > float64(maxUint) || val < 0 {
					iter.ReportError("fuzzy decode uint", "exceed range")
					return
				}
				*((*uint)(ptr)) = uint(val)
			} else {
				dec.Decode(ptr, iter)
			}
		}}
	},
	reflect.Uint8: func(e *ProtoExtension, dec jsoniter.ValDecoder) jsoniter.ValDecoder {
		return &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			if isFloat {
				val := iter.ReadFloat64()
				if _, frac := math.Modf(val); frac != 0 {
					iter.ReportError("fuzzyIntegerDecoder", "found frac")
					return
				}
				if val > float64(math.MaxUint8) || val < 0 {
					iter.ReportError("fuzzy decode uint8", "exceed range")
					return
				}
				*((*uint8)(ptr)) = uint8(val)
			} else {
				dec.Decode(ptr, iter)
			}
		}}
	},
	reflect.Uint16: func(e *ProtoExtension, dec jsoniter.ValDecoder) jsoniter.ValDecoder {
		return &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			if isFloat {
				val := iter.ReadFloat64()
				if _, frac := math.Modf(val); frac != 0 {
					iter.ReportError("fuzzyIntegerDecoder", "found frac")
					return
				}
				if val > float64(math.MaxUint16) || val < 0 {
					iter.ReportError("fuzzy decode uint16", "exceed range")
					return
				}
				*((*uint16)(ptr)) = uint16(val)
			} else {
				dec.Decode(ptr, iter)
			}
		}}
	},
	reflect.Uint32: func(e *ProtoExtension, dec jsoniter.ValDecoder) jsoniter.ValDecoder {
		return &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			if isFloat {
				val := iter.ReadFloat64()
				if _, frac := math.Modf(val); frac != 0 {
					iter.ReportError("fuzzyIntegerDecoder", "found frac")
					return
				}
				if val > float64(math.MaxUint32) || val < 0 {
					iter.ReportError("fuzzy decode uint32", "exceed range")
					return
				}
				*((*uint32)(ptr)) = uint32(val)
			} else {
				dec.Decode(ptr, iter)
			}
		}}
	},
	reflect.Uint64: func(e *ProtoExtension, dec jsoniter.ValDecoder) jsoniter.ValDecoder {
		return &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			if isFloat {
				val := iter.ReadFloat64()
				if _, frac := math.Modf(val); frac != 0 {
					iter.ReportError("fuzzyIntegerDecoder", "found frac")
					return
				}
				if val > float64(math.MaxUint64) || val < 0 {
					iter.ReportError("fuzzy decode uint64", "exceed range")
					return
				}
				*((*uint64)(ptr)) = uint64(val)
			} else {
				dec.Decode(ptr, iter)
			}
		}}
	},
}

type fuzzyIntegerDecoder struct {
	fun func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator)
}

func (decoder *fuzzyIntegerDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	var str string
	switch valueType {
	case jsoniter.NumberValue:
		var number json.Number
		iter.ReadVal(&number)
		str = string(number)
	case jsoniter.StringValue:
		str = iter.ReadString()
		switch str {
		case "true":
			str = "1"
		case "false":
			str = "0"
		}
	case jsoniter.BoolValue:
		if iter.ReadBool() {
			str = "1"
		} else {
			str = "0"
		}
	case jsoniter.NilValue:
		iter.Skip()
		str = "0"
	default:
		iter.ReportError("fuzzyIntegerDecoder", "not number or string")
	}
	if len(str) == 0 {
		str = "0"
	}
	newIter := iter.Pool().BorrowIterator([]byte(str))
	newIter.Attachment = iter.Attachment
	defer iter.Pool().ReturnIterator(newIter)
	isFloat := strings.ContainsAny(str, ".eE")
	decoder.fun(isFloat, ptr, newIter)
	if newIter.Error != nil && newIter.Error != io.EOF {
		iter.Error = newIter.Error
	}
}
