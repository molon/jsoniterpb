package jsoniterpb

import (
	"fmt"
	"strings"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

const (
	FieldMask_Paths_field_fullname protoreflect.FullName = "google.protobuf.FieldMask.paths"
)

var wktFieldmaskCodec = (&ProtoCodec{}).
	SetElemEncodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, stream *jsoniter.Stream) {
		v := ((*fieldmaskpb.FieldMask)(ptr))
		paths := make([]string, 0, len(v.GetPaths()))
		for _, s := range v.GetPaths() {
			if !protoreflect.FullName(s).IsValid() {
				stream.Error = fmt.Errorf("%s contains invalid path: %q", FieldMask_Paths_field_fullname, s)
				return
			}
			// Return error if conversion to camelCase is not reversible.
			cc := JSONCamelCase(s)
			if s != JSONSnakeCase(cc) {
				stream.Error = fmt.Errorf("%s contains irreversible value %q", FieldMask_Paths_field_fullname, s)
				return
			}
			paths = append(paths, cc)
		}
		stream.WriteVal(strings.Join(paths, ","))
	}).
	SetElemDecodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		var str string
		iter.ReadVal(&str)
		if str == "" {
			(*fieldmaskpb.FieldMask)(ptr).Paths = []string{}
			return
		}
		paths := strings.Split(str, ",")
		for idx, s0 := range paths {
			s := JSONSnakeCase(s0)
			if strings.Contains(s0, "_") || !protoreflect.FullName(s).IsValid() {
				iter.ReportError("protobuf", fmt.Sprintf("%v contains invalid path: %q", FieldMask_Paths_field_fullname, s0))
				return
			}
			paths[idx] = s
		}
		(*fieldmaskpb.FieldMask)(ptr).Paths = paths
	})

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}
func isASCIIUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// JSONCamelCase converts a snake_case identifier to a camelCase identifier,
// according to the protobuf JSON specification.
func JSONCamelCase(s string) string {
	var b []byte
	var wasUnderscore bool
	for i := 0; i < len(s); i++ { // proto identifiers are always ASCII
		c := s[i]
		if c != '_' {
			if wasUnderscore && isASCIILower(c) {
				c -= 'a' - 'A' // convert to uppercase
			}
			b = append(b, c)
		}
		wasUnderscore = c == '_'
	}
	return string(b)
}

// JSONSnakeCase converts a camelCase identifier to a snake_case identifier,
// according to the protobuf JSON specification.
func JSONSnakeCase(s string) string {
	var b []byte
	for i := 0; i < len(s); i++ { // proto identifiers are always ASCII
		c := s[i]
		if isASCIIUpper(c) {
			b = append(b, '_')
			c += 'a' - 'A' // convert to lowercase
		}
		b = append(b, c)
	}
	return string(b)
}
