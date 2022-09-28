package jsoniterpb

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"github.com/molon/jsoniterpb/extra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var protoMessageType = reflect2.TypeOfPtr((*proto.Message)(nil)).Elem()

func (e *ProtoExtension) updateStructDescriptorConstructorForOneOf(c *jsoniter.StructDescriptorConstructor) {
	if !reflect2.PtrTo(c.Type).Implements(protoMessageType) {
		return
	}
	if c.Type == wktValueType {
		return
	}

	newBindings := make([]*jsoniter.Binding, 0, len(c.Bindings))
	defer func() {
		c.Bindings = newBindings
	}()

	var pb proto.Message
	var pbReflect protoreflect.Message
	for _, binding := range c.Bindings {
		field := binding.Field

		if field.Type().Kind() == reflect.Interface {
			oneofsTag, hasOneofsTag := field.Tag().Lookup("protobuf_oneof")
			if hasOneofsTag {
				if pb == nil {
					pb = c.Type.New().(proto.Message)
					pbReflect = pb.ProtoReflect()
				}
				fieldType := field.Type()
				fieldPtr := field.UnsafeGet(reflect2.PtrOf(pb))
				od := pbReflect.Descriptor().Oneofs().ByName(protoreflect.Name(oneofsTag))
				if !od.IsSynthetic() { // ignore optional
					fds := od.Fields()
					for j := 0; j < fds.Len(); j++ {
						fd := fds.Get(j)
						value := pbReflect.NewField(fd)
						pbReflect.Set(fd, value)

						fTyp := reflect2.TypeOf(fieldType.UnsafeIndirect(fieldPtr))
						if fTyp.Kind() == reflect.Ptr {
							wrapPtrType := fTyp.(*reflect2.UnsafePtrType)
							if wrapPtrType.Elem().Kind() == reflect.Struct {
								structDescriptor := c.DescribeStructFunc(wrapPtrType.Elem())
								for _, b := range structDescriptor.Fields {
									b.Levels = append([]int{binding.Levels[0], j}, b.Levels...)
									omitempty := b.Encoder.(*jsoniter.StructFieldEncoder).OmitEmpty
									b.Encoder = &protoOneofWrapperEncoder{wrapPtrType, b.Field, b.Encoder}
									b.Encoder = &jsoniter.StructFieldEncoder{field, b.Encoder, omitempty}
									b.Decoder = &protoOneofWrapperDecoder{field.Type(), wrapPtrType, wrapPtrType.Elem(), b.Field, b.Decoder}
									b.Decoder = &jsoniter.StructFieldDecoder{field, b.Decoder}
									c.EmbeddedBindings = append(c.EmbeddedBindings, b)
								}
								continue
							}
						}
					}
					continue
				}
			}
		}

		newBindings = append(newBindings, binding)

		if len(binding.FromNames) <= 0 { // simple check should exported
			continue
		}

		tag, hastag := binding.Field.Tag().Lookup("protobuf")
		if !hastag {
			continue
		}

		var name string
		tagParts := strings.Split(tag, ",")
		for _, part := range tagParts {
			colons := strings.SplitN(part, "=", 2)
			if len(colons) == 2 {
				if strings.TrimSpace(colons[0]) == "name" {
					name = strings.TrimSpace(colons[1])
				}
				continue
			}
			if strings.TrimSpace(part) == "oneof" {
				if pb == nil {
					pb = c.Type.New().(proto.Message)
					pbReflect = pb.ProtoReflect()
				}
				od := pbReflect.Descriptor().Fields().ByName(protoreflect.Name(name))
				if od != nil {
					oneof := od.ContainingOneof()
					// IsSynthetic OneOf (optional keyword)
					if oneof != nil && oneof.IsSynthetic() {
						binding.Encoder = &extra.ImmunityEmitEmptyEncoder{
							&protoOptionalEncoder{binding.Encoder},
						}
					}
				}
			}
		}
	}
}

type protoOptionalEncoder struct {
	jsoniter.ValEncoder
}

func (enc *protoOptionalEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	enc.ValEncoder.Encode(ptr, stream)
}

func (enc *protoOptionalEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*unsafe.Pointer)(ptr)) == nil
}

type protoOneofWrapperEncoder struct {
	wrapperPtrType reflect2.Type
	valueField     reflect2.StructField
	valueEncoder   jsoniter.ValEncoder
}

func (encoder *protoOneofWrapperEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	if *((*unsafe.Pointer)(ptr)) == nil {
		stream.WriteNil()
		return
	}
	val := reflect2.IFaceToEFace(ptr)
	if reflect2.TypeOf(val).RType() != encoder.wrapperPtrType.RType() {
		stream.WriteNil()
		return
	}
	encoder.valueEncoder.Encode(reflect2.PtrOf(val), stream)
	if stream.Error != nil && stream.Error != io.EOF {
		stream.Error = fmt.Errorf("%s: %s", encoder.valueField.Name(), stream.Error.Error())
	}
}

func (encoder *protoOneofWrapperEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	if *((*unsafe.Pointer)(ptr)) == nil {
		return true
	}
	val := reflect2.IFaceToEFace(ptr)
	if reflect2.TypeOf(val).RType() != encoder.wrapperPtrType.RType() {
		return true
	}
	return encoder.valueEncoder.IsEmpty(reflect2.PtrOf(val))
}

func (encoder *protoOneofWrapperEncoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	if *((*unsafe.Pointer)(ptr)) == nil {
		return true
	}
	val := reflect2.IFaceToEFace(ptr)
	if reflect2.TypeOf(val).RType() != encoder.wrapperPtrType.RType() {
		return true
	}
	isEmbeddedPtrNil, converted := encoder.valueEncoder.(jsoniter.IsEmbeddedPtrNil)
	if !converted {
		return false
	}
	return isEmbeddedPtrNil.IsEmbeddedPtrNil(reflect2.PtrOf(val))
}

type protoOneofWrapperDecoder struct {
	wrapperIfaceType reflect2.Type
	wrapperPtrType   reflect2.Type
	wrapperElemType  reflect2.Type
	valueField       reflect2.StructField
	valueDecoder     jsoniter.ValDecoder
}

func (decoder *protoOneofWrapperDecoder) Decode(fieldPtr unsafe.Pointer, iter *jsoniter.Iterator) {
	var elem interface{}

	// reuse it if type match
	if *((*unsafe.Pointer)(fieldPtr)) != nil {
		elem = reflect2.IFaceToEFace(fieldPtr)
		if reflect2.TypeOf(elem).RType() != decoder.wrapperPtrType.RType() {
			elem = nil
		}
	}
	if elem == nil {
		elem = decoder.wrapperElemType.New()
	}

	decoder.valueDecoder.Decode(reflect2.PtrOf(elem), iter)
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%s: %s", decoder.valueField.Name(), iter.Error.Error())
		return
	}

	rval := reflect.ValueOf(decoder.wrapperIfaceType.PackEFace(fieldPtr))
	rval.Elem().Set(reflect.ValueOf(elem))
}
