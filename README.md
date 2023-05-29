# jsoniterpb
Replacement of Protojson over Jsoniter Extension

### Features
- Handle any type of object, not just `proto.Message`, to get a consistent format even with nested uses
- All features of `protojson`: `ProtobufWellKnownType/Oneof/JsonName/64IntToStr/SortMapKeysByRealValue/CheckUT8/...`
- Support more fuzzy decode methods
- Better performance

### Compatibility test
`cd ./internal/protojson && go run ./gen.go`, it will download the latest tests file from `protocolbuffers/protobuf-go` and make it available to `jsoniterpb`

### Warns
Some differences with `protojson`
- Only `proto3` is supported, `proto2` is not supported
- `protojson` marshal nil `proto.Message` as zero value **if it is root**. but `jsoniterpb` will marshal it to `null`
- View [internal/protojson/tests/jsoniterpb_decode_test.go](internal/protojson/tests/jsoniterpb_decode_test.go)
  - Support more fuzzy decode methods => Search `FuzzyDecode`
  - Most error messages are not the same => Search `ErrMsgNotSame`
  - Some error check are not supported => Search `NotSupport`

### Usage
Since the current extensibility of [json-iterator/go](https://github.com/json-iterator/go) is not enough to complete this project, it needs to be replaced with another version.
```
// go.mod 
// go get github.com/molon/jsoniter@jsoniterpb
replace github.com/json-iterator/go => github.com/molon/jsoniter v0.0.0-20230529062209-e42e40bd8588
```

```
// protojson.MarshalOptions{} equals
cfg := jsoniter.Config{SortMapKeys: true, DisallowUnknownFields: true}.Froze()
cfg.RegisterExtension(&jsoniterpb.ProtoExtension{})

// protojson.MarshalOptions{EmitUnpopulated: true} equals
cfg := jsoniter.Config{SortMapKeys: true, DisallowUnknownFields: true}.Froze()
cfg.RegisterExtension(&jsoniterpb.ProtoExtension{EmitUnpopulated: true})

// protojson.UnmarshalOptions{DiscardUnknown: true} equals
cfg := jsoniter.Config{SortMapKeys: true, DisallowUnknownFields: false}.Froze()
cfg.RegisterExtension(&jsoniterpb.ProtoExtension{})
```

### Benchmark
```
goos: darwin
goarch: arm64
pkg: github.com/molon/jsoniterpb
BenchmarkWrite
BenchmarkWrite/protojson
BenchmarkWrite/protojson-8         	    4590	    252336 ns/op	  120711 B/op	    2360 allocs/op
BenchmarkWrite/jsoniter
BenchmarkWrite/jsoniter-8          	    6589	    188076 ns/op	   88383 B/op	    2313 allocs/op
BenchmarkWrite/jsoniter-fast
BenchmarkWrite/jsoniter-fast-8     	   10000	    119375 ns/op	   47925 B/op	    1146 allocs/op
```
```
goos: darwin
goarch: arm64
pkg: github.com/molon/jsoniterpb
BenchmarkRead
BenchmarkRead/protojson
BenchmarkRead/protojson-8         	    3328	    371155 ns/op	  113408 B/op	    4021 allocs/op
BenchmarkRead/jsoniter
BenchmarkRead/jsoniter-8          	    5720	    204019 ns/op	   87803 B/op	    2790 allocs/op
BenchmarkRead/jsoniter-nofuzzydecode
BenchmarkRead/jsoniter-nofuzzydecode-8         	    6697	    177609 ns/op	   75491 B/op	    2185 allocs/op
```
