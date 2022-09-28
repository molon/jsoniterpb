package extra

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

type ImmunityEmitEmptyEncoder struct {
	jsoniter.ValEncoder
}

func (enc *ImmunityEmitEmptyEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	enc.ValEncoder.Encode(ptr, stream)
}

func (enc *ImmunityEmitEmptyEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return enc.ValEncoder.IsEmpty(ptr)
}

func (enc *ImmunityEmitEmptyEncoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	isEmbeddedPtrNil, converted := enc.ValEncoder.(jsoniter.IsEmbeddedPtrNil)
	if !converted {
		return false
	}
	return isEmbeddedPtrNil.IsEmbeddedPtrNil(ptr)
}

type EmitEmptyEncoder struct {
	jsoniter.ValEncoder
}

func (enc *EmitEmptyEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	enc.ValEncoder.Encode(ptr, stream)
}

func (enc *EmitEmptyEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	if _, ok := enc.ValEncoder.(*ImmunityEmitEmptyEncoder); ok {
		return enc.ValEncoder.IsEmpty(ptr)
	}
	return false
}

func (enc *EmitEmptyEncoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	isEmbeddedPtrNil, converted := enc.ValEncoder.(jsoniter.IsEmbeddedPtrNil)
	if !converted {
		return false
	}
	return isEmbeddedPtrNil.IsEmbeddedPtrNil(ptr)
}

type EmitEmptyExtension struct {
	jsoniter.DummyExtension
	Filter func(binding *jsoniter.Binding) bool
}

func (e *EmitEmptyExtension) UpdateStructDescriptor(desc *jsoniter.StructDescriptor) {
	for _, binding := range desc.Fields {
		if binding.Encoder != nil {
			if e.Filter == nil || e.Filter(binding) {
				if _, ok := binding.Encoder.(*ImmunityEmitEmptyEncoder); ok {
					continue
				}
				binding.Encoder = &EmitEmptyEncoder{ValEncoder: binding.Encoder}
			}
		}
	}
}
