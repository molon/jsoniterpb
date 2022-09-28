package jsoniterpb

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
	"sync"
	"unicode/utf8"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

func (e *ProtoExtension) updateMapEncoderConstructorForScalar(v *jsoniter.MapEncoderConstructor) {
	// handle 64bit integer key, avoid quote it repeatedly
	if e.Encode64BitAsInteger {
		return
	}
	typ := v.MapType.Key()
	switch typ.Kind() {
	case reflect.Int64, reflect.Uint64:
		v.KeyEncoder = &dynamicEncoder{v.MapType.Key()}
	}
}

func (e *ProtoExtension) decorateEncoderForScalar(typ reflect2.Type, enc jsoniter.ValEncoder) jsoniter.ValEncoder {
	var bitSize int
	switch typ.Kind() {
	case reflect.String:
		if e.PermitInvalidUTF8 {
			return enc
		}
		return &protoStringEncoder{}
	case reflect.Int64, reflect.Uint64:
		// https://developers.google.com/protocol-buffers/docs/proto3 int64, fixed64, uint64 should be string
		// https://github.com/protocolbuffers/protobuf-go/blob/e62d8edb7570c986a51e541c161a0c93bbaf9253/encoding/protojson/encode.go#L274-L277
		// https://github.com/protocolbuffers/protobuf-go/pull/14
		// https://github.com/golang/protobuf/issues/1414
		if e.Encode64BitAsInteger {
			return enc
		}
		return &stringModeNumberEncoder{enc}
	case reflect.Float32:
		bitSize = 32
	case reflect.Float64:
		bitSize = 64
	}

	if bitSize <= 0 {
		return enc
	}

	return &funcEncoder{
		fun: func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			var n float64
			if bitSize == 32 {
				n = float64(*((*float32)(ptr)))
			} else {
				n = *((*float64)(ptr))
			}
			switch {
			case math.IsNaN(n):
				stream.WriteRaw(`"NaN"`)
			case math.IsInf(n, +1):
				stream.WriteRaw(`"Infinity"`)
			case math.IsInf(n, -1):
				stream.WriteRaw(`"-Infinity"`)
			default:
				enc.Encode(ptr, stream)
			}
		},
		isEmptyFunc: func(ptr unsafe.Pointer) bool {
			return enc.IsEmpty(ptr)
		},
	}
}

var (
	nanBytes           = []byte(`"NaN"`)
	infBytes           = []byte(`"Infinity"`)
	ninfBytes          = []byte(`"-Infinity"`)
	jsonNumberElemType = reflect2.TypeOfPtr((*json.Number)(nil)).Elem()
)

func (e *ProtoExtension) decorateDecoderForScalar(typ reflect2.Type, dec jsoniter.ValDecoder) jsoniter.ValDecoder {
	if typ.Implements(protoEnumType) || typ == jsonNumberElemType {
		return dec
	}

	// []byte
	if typ.Kind() == reflect.Slice && typ.(reflect2.SliceType).Elem().Kind() == reflect.Uint8 {
		return &funcDecoder{
			fun: func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
				if iter.WhatIsNext() == jsoniter.StringValue {
					s := iter.ReadString()
					// copy from protobuf-go
					enc := base64.StdEncoding
					if strings.ContainsAny(s, "-_") {
						enc = base64.URLEncoding
					}
					if len(s)%4 != 0 {
						enc = enc.WithPadding(base64.NoPadding)
					}

					dst, err := enc.DecodeString(s)
					if err != nil {
						iter.ReportError("decode base64", err.Error())
					} else {
						typ.UnsafeSet(ptr, unsafe.Pointer(&dst))
					}
					return
				}
				dec.Decode(ptr, iter)
			},
		}
	}

	if !e.DisableFuzzyDecode {
		if ddec, ok := fuzzyDecorateScalarDecoders[typ.Kind()]; ok {
			dec = ddec(e, dec)
		}
	}

	var bitSize int
	switch typ.Kind() {
	case reflect.Int64, reflect.Uint64:
		return &stringModeNumberDecoder{elemDecoder: dec}
	case reflect.String:
		return &funcDecoder{
			fun: func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
				dec.Decode(ptr, iter)
				if iter.Error == nil {
					if !e.PermitInvalidUTF8 {
						if !utf8.ValidString(*((*string)(ptr))) {
							iter.Error = errInvalidUTF8
						}
					}
				}
			},
		}
	case reflect.Float32:
		bitSize = 32
	case reflect.Float64:
		bitSize = 64
	}

	if bitSize <= 0 {
		return dec
	}

	return &funcDecoder{
		fun: func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			if iter.WhatIsNext() == jsoniter.StringValue {
				b := iter.SkipAndReturnBytes()
				if bytes.Equal(b, nanBytes) {
					if bitSize == 32 {
						*((*float32)(ptr)) = float32(math.NaN())
					} else {
						*((*float64)(ptr)) = math.NaN()
					}
				} else if bytes.Equal(b, infBytes) {
					if bitSize == 32 {
						*((*float32)(ptr)) = float32(math.Inf(+1))
					} else {
						*((*float64)(ptr)) = math.Inf(+1)
					}
				} else if bytes.Equal(b, ninfBytes) {
					if bitSize == 32 {
						*((*float32)(ptr)) = float32(math.Inf(-1))
					} else {
						*((*float64)(ptr)) = math.Inf(-1)
					}
				} else {
					subIter := iter.Pool().BorrowIterator(b)
					subIter.Attachment = iter.Attachment
					defer iter.API().ReturnIterator(subIter)
					dec.Decode(ptr, subIter)
					if subIter.Error != nil && subIter.Error != io.EOF && iter.Error == nil {
						iter.Error = subIter.Error
					}
				}
				return
			}
			dec.Decode(ptr, iter)
		},
	}
}

type protoStringEncoder struct {
	once       sync.Once
	escapeHTML bool
}

func (encoder *protoStringEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	encoder.once.Do(func() {
		if fcfg, ok := stream.API().(interface {
			GetConfig() jsoniter.Config
		}); ok {
			encoder.escapeHTML = fcfg.GetConfig().EscapeHTML
		}
	})

	str := *((*string)(ptr))
	var buf []byte
	var err error
	if encoder.escapeHTML {
		buf, err = QuoteValidUTF8StringWithHTMLEscaped(str)
	} else {
		buf, err = QuoteValidUTF8String(str)
	}
	if err != nil {
		stream.Error = fmt.Errorf("ProtoStringEncoder: %w", err)
		return
	}
	stream.Write(buf)
}

func (encoder *protoStringEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*string)(ptr)) == ""
}

func QuoteValidUTF8String(s string) ([]byte, error) {
	valLen := len(s)

	buf := []byte{'"'}
	// write string, the fast path, without utf8 and escape support
	i := 0
	for ; i < valLen; i++ {
		c := s[i]
		if c < utf8.RuneSelf && safeSet[c] {
			buf = append(buf, c)
		} else {
			break
		}
	}
	if i == valLen {
		buf = append(buf, '"')
		return buf, nil
	}
	return appendStringSlowPath(buf, i, s, valLen)
}

var errInvalidUTF8 = errors.New("invalid UTF-8")

func appendStringSlowPath(buf []byte, i int, s string, valLen int) ([]byte, error) {
	start := i
	// for the remaining parts, we process them char by char
	for i < valLen {
		if b := s[i]; b < utf8.RuneSelf {
			if safeSet[b] {
				i++
				continue
			}
			if start < i {
				buf = append(buf, s[start:i]...)
			}
			switch b {
			case '\\', '"':
				buf = append(buf, '\\', b)
			case '\b':
				buf = append(buf, '\\', 'b')
			case '\f':
				buf = append(buf, '\\', 'f')
			case '\n':
				buf = append(buf, '\\', 'n')
			case '\r':
				buf = append(buf, '\\', 'r')
			case '\t':
				buf = append(buf, '\\', 't')
			default:
				buf = append(buf, `\u00`...)
				buf = append(buf, hex[b>>4], hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			return buf, errInvalidUTF8
		}
		i += size
	}
	if start < len(s) {
		buf = append(buf, s[start:]...)
	}
	buf = append(buf, '"')
	return buf, nil
}

func QuoteValidUTF8StringWithHTMLEscaped(s string) ([]byte, error) {
	valLen := len(s)
	buf := []byte{'"'}
	// write string, the fast path, without utf8 and escape support
	i := 0
	for ; i < valLen; i++ {
		c := s[i]
		if c < utf8.RuneSelf && htmlSafeSet[c] {
			buf = append(buf, c)
		} else {
			break
		}
	}
	if i == valLen {
		buf = append(buf, '"')
		return buf, nil
	}
	return appendStringSlowPathWithHTMLEscaped(buf, i, s, valLen)
}

func appendStringSlowPathWithHTMLEscaped(buf []byte, i int, s string, valLen int) ([]byte, error) {
	start := i
	// for the remaining parts, we process them char by char
	for i < valLen {
		if b := s[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] {
				i++
				continue
			}
			if start < i {
				buf = append(buf, s[start:i]...)
			}
			switch b {
			case '\\', '"':
				buf = append(buf, '\\', b)
			case '\b':
				buf = append(buf, '\\', 'b')
			case '\f':
				buf = append(buf, '\\', 'f')
			case '\n':
				buf = append(buf, '\\', 'n')
			case '\r':
				buf = append(buf, '\\', 'r')
			case '\t':
				buf = append(buf, '\\', 't')
			default:
				buf = append(buf, `\u00`...)
				buf = append(buf, hex[b>>4], hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			return buf, errInvalidUTF8
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				buf = append(buf, s[start:i]...)
			}
			buf = append(buf, `\u202`...)
			buf = append(buf, hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		buf = append(buf, s[start:]...)
	}
	buf = append(buf, '"')
	return buf, nil
}

var hex = "0123456789abcdef"

var safeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

var htmlSafeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      false,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      false,
	'=':      true,
	'>':      false,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}
