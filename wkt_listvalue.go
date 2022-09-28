package jsoniterpb

import (
	"io"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/types/known/structpb"
)

var wktListValueCodec = (&ProtoCodec{}).
	SetElemEncodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, stream *jsoniter.Stream) {
		x := ((*structpb.ListValue)(ptr))
		if x.Values == nil {
			stream.WriteEmptyArray()
			return
		}
		stream.WriteVal(x.Values)
	}).
	SetElemDecodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		x := ((*structpb.ListValue)(ptr))

		values := []*structpb.Value{}
		iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
			v := &structpb.Value{}
			iter.ReadVal(v)
			if iter.Error != nil && iter.Error != io.EOF {
				return false
			}
			values = append(values, v)
			return true
		})

		x.Values = values
	})
