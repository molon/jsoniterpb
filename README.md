# jsoniterpb
Replacement of Protojson over Jsoniter Extension

#### Features
- Most of the features are the same as protojson.
- Support more fuzzy decode.
- Support for objects that do not implement `proto.Message`
- 2x performance

#### Warns
Some features of protojson are not supported
- Only `proto3` is supported, `proto2` is not supported
- View [internal/protojson/tests/jsoniterpb_decode_test.go](internal/protojson/tests/jsoniterpb_decode_test.go) and search `NotSupport`

#### Usage
```
// protojson.MarshalOptions{}
cfg := jsoniter.Config{SortMapKeys: true, DisallowUnknownFields: true}.Froze()
cfg.RegisterExtension(&jsoniterpb.ProtoExtension{})

// protojson.MarshalOptions{EmitUnpopulated: true}
cfg := jsoniter.Config{SortMapKeys: true, DisallowUnknownFields: true}.Froze()
cfg.RegisterExtension(&jsoniterpb.ProtoExtension{EmitUnpopulated: true})

// protojson.UnmarshalOptions{DiscardUnknown: true}
cfg := jsoniter.Config{SortMapKeys: true, DisallowUnknownFields: false}.Froze()
cfg.RegisterExtension(&jsoniterpb.ProtoExtension{})
```