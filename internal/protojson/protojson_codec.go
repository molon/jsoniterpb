package main

import (
	"fmt"
	"io"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"github.com/molon/jsoniterpb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// TIPS: does not use this now
type ProtojsonEncoder struct {
	ElemType    reflect2.Type
	MarshalOpts protojson.MarshalOptions
}

func (enc *ProtojsonEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	// TODO: indent?
	data, err := enc.MarshalOpts.Marshal(enc.ElemType.PackEFace(ptr).(proto.Message))
	if err != nil {
		stream.Error = fmt.Errorf("error calling protojson.Marshal for type %s: %w", reflect2.PtrTo(enc.ElemType), err)
		return
	}
	_, stream.Error = stream.Write(data)
}

func (enc *ProtojsonEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	// protojson will not omit zero value, only omit zero pointer, we stay compatible,
	return false
}

type ProtojsonDecoder struct {
	ElemType      reflect2.Type
	UnmarshalOpts protojson.UnmarshalOptions
}

func (dec *ProtojsonDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	bytes := iter.SkipAndReturnBytes()
	if iter.Error != nil && iter.Error != io.EOF {
		return
	}

	err := dec.UnmarshalOpts.Unmarshal(bytes, dec.ElemType.PackEFace(ptr).(proto.Message))
	if err != nil {
		iter.ReportError("protobuf", fmt.Sprintf(
			"error calling protojson.Unmarshal for type %s: %s",
			reflect2.PtrTo(dec.ElemType), err,
		))
	}
}

var ProtojsonEncoderCreator = func(e *jsoniterpb.ProtoExtension, typ reflect2.Type) jsoniter.ValEncoder {
	return jsoniterpb.WrapElemEncoder(typ, &ProtojsonEncoder{
		ElemType: typ,
		MarshalOpts: protojson.MarshalOptions{
			EmitUnpopulated: e.EmitUnpopulated,
			UseEnumNumbers:  e.UseEnumNumbers,
			UseProtoNames:   e.UseProtoNames,
			Resolver:        e.Resolver,
		},
	}, nil)
}

var ProtojsonDecoderCreator = func(e *jsoniterpb.ProtoExtension, typ reflect2.Type) jsoniter.ValDecoder {
	return jsoniterpb.WrapElemDecoder(typ, &ProtojsonDecoder{
		ElemType: typ,
		UnmarshalOpts: protojson.UnmarshalOptions{
			Resolver:       e.Resolver,
			DiscardUnknown: true, // TODO: ???
		},
	}, nil)
}

func NewProtojsonCodec() *jsoniterpb.ProtoCodec {
	return &jsoniterpb.ProtoCodec{
		EncoderCreator: ProtojsonEncoderCreator,
		DecoderCreator: ProtojsonDecoderCreator,
	}
}
