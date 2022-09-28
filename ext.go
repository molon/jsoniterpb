package jsoniterpb

import (
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"github.com/molon/jsoniterpb/extra"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type ProtoExtension struct {
	jsoniter.DummyExtension

	EmitUnpopulated bool
	UseEnumNumbers  bool
	UseProtoNames   bool
	Resolver        interface {
		protoregistry.MessageTypeResolver
		// TIPS: does not support it now
		protoregistry.ExtensionTypeResolver
	}

	Encode64BitAsInteger bool
	SortMapKeysAsString  bool
	PermitInvalidUTF8    bool
	DisableFuzzyDecode   bool
}

func (e *ProtoExtension) GetResolver() interface {
	protoregistry.MessageTypeResolver
	protoregistry.ExtensionTypeResolver
} {
	if e.Resolver != nil {
		return e.Resolver
	}
	return protoregistry.GlobalTypes
}

func (e *ProtoExtension) CreateEncoder(typ reflect2.Type) jsoniter.ValEncoder {
	if enc := e.createProtoEncoder(typ); enc != nil {
		return enc
	}
	if enc := e.createProtoEnumEncoder(typ); enc != nil {
		return enc
	}
	return nil
}

func (e *ProtoExtension) CreateDecoder(typ reflect2.Type) jsoniter.ValDecoder {
	if dec := e.createProtoDecoder(typ); dec != nil {
		return dec
	}
	if dec := e.createProtoEnumDecoder(typ); dec != nil {
		return dec
	}
	return nil
}

func (e *ProtoExtension) UpdateMapEncoderConstructor(v *jsoniter.MapEncoderConstructor) {
	e.updateMapEncoderConstructorForNonNull(v)
	e.updateMapEncoderConstructorForSortMapKeys(v)
	e.updateMapEncoderConstructorForScalar(v)
}

func (e *ProtoExtension) UpdateSliceEncoderConstructor(v *jsoniter.SliceEncoderConstructor) {
	e.updateSliceEncoderConstructorForNonNull(v)
}

func (e *ProtoExtension) UpdateArrayEncoderConstructor(v *jsoniter.ArrayEncoderConstructor) {
	e.updateArrayEncoderConstructorForNonNull(v)
}

func (e *ProtoExtension) DecorateEncoder(typ reflect2.Type, encoder jsoniter.ValEncoder) jsoniter.ValEncoder {
	if enc := e.decorateEncoderForNilCollection(typ, encoder); enc != nil {
		encoder = enc
	}
	if enc := e.decorateEncoderForScalar(typ, encoder); enc != nil {
		encoder = enc
	}
	return encoder
}

func (e *ProtoExtension) DecorateDecoder(typ reflect2.Type, decoder jsoniter.ValDecoder) jsoniter.ValDecoder {
	if dec := e.decorateDecoderForNil(typ, decoder); dec != nil {
		decoder = dec
	}
	if dec := e.decorateDecoderForScalar(typ, decoder); dec != nil {
		decoder = dec
	}
	return decoder
}

func (e *ProtoExtension) UpdateStructDescriptorConstructor(c *jsoniter.StructDescriptorConstructor) {
	e.updateStructDescriptorConstructorForOneOf(c)
}

// Handle EmitUnpopulated and UseProtoNames
func (e *ProtoExtension) UpdateStructDescriptor(desc *jsoniter.StructDescriptor) {
	for _, binding := range desc.Fields {
		if len(binding.FromNames) <= 0 { // simple check should exported
			continue
		}

		// Because oneof wrapper does not satisfy proto.Message, we can only check with tag instead of protoreflect here
		tag, hastag := binding.Field.Tag().Lookup("protobuf")
		if !hastag {
			continue
		}

		if e.EmitUnpopulated {
			binding.Encoder = &extra.EmitEmptyEncoder{binding.Encoder}
		}

		var name, jsonName string
		tagParts := strings.Split(tag, ",")
		for _, part := range tagParts {
			colons := strings.SplitN(part, "=", 2)
			if len(colons) == 2 {
				switch strings.TrimSpace(colons[0]) {
				case "name":
					name = strings.TrimSpace(colons[1])
				case "json":
					jsonName = strings.TrimSpace(colons[1])
				}
				continue
			}
		}
		if jsonName == "" {
			jsonName = name
		}
		if name != "" {
			if e.UseProtoNames {
				binding.FromNames = []string{name}
				// fuzzy
				if jsonName != name {
					binding.FromNames = append(binding.FromNames, jsonName)
				}
				binding.ToNames = []string{name}
			} else {
				binding.FromNames = []string{jsonName}
				// fuzzy
				if name != jsonName {
					binding.FromNames = append(binding.FromNames, name)
				}
				binding.ToNames = []string{jsonName}
			}
		}
	}
}
