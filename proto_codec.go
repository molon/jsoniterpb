package jsoniterpb

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

var ProtoCodecs = map[reflect2.Type]*ProtoCodec{}

func (e *ProtoExtension) createProtoEncoder(typ reflect2.Type) (xret jsoniter.ValEncoder) {
	codec, ok := ProtoCodecs[typ]
	if ok && codec != nil && codec.EncoderCreator != nil {
		return codec.EncoderCreator(e, typ)
	}
	return nil
}

func (e *ProtoExtension) createProtoDecoder(typ reflect2.Type) (xret jsoniter.ValDecoder) {
	codec, ok := ProtoCodecs[typ]
	if ok && codec != nil && codec.DecoderCreator != nil {
		return codec.DecoderCreator(e, typ)
	}
	return nil
}

type ProtoCodec struct {
	EncoderCreator func(e *ProtoExtension, typ reflect2.Type) jsoniter.ValEncoder
	DecoderCreator func(e *ProtoExtension, typ reflect2.Type) jsoniter.ValDecoder
}

func (codec *ProtoCodec) SetElemEncodeFunc(encodeFunc func(e *ProtoExtension, ptr unsafe.Pointer, stream *jsoniter.Stream)) *ProtoCodec {
	codec.EncoderCreator = func(e *ProtoExtension, typ reflect2.Type) jsoniter.ValEncoder {
		return WrapElemEncoder(typ, &funcEncoder{
			fun: func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
				encodeFunc(e, ptr, stream)
			},
		}, nil)
	}
	return codec
}

func (codec *ProtoCodec) SetElemDecodeFunc(decodeFunc func(e *ProtoExtension, ptr unsafe.Pointer, iter *jsoniter.Iterator)) *ProtoCodec {
	codec.DecoderCreator = func(e *ProtoExtension, typ reflect2.Type) jsoniter.ValDecoder {
		return WrapElemDecoder(typ, &funcDecoder{
			fun: func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
				decodeFunc(e, ptr, iter)
			},
		}, nil)
	}
	return codec
}
