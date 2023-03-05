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
replace github.com/json-iterator/go => github.com/molon/jsoniter v0.0.0-20230305181513-eac2ab4f5edf
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
BenchmarkWrite/protojson-8                  4699            256354 ns/op          118876 B/op       2371 allocs/op
BenchmarkWrite/jsoniter
BenchmarkWrite/jsoniter-8                   6052            185370 ns/op           88997 B/op       2327 allocs/op
BenchmarkWrite/jsoniter-fast
BenchmarkWrite/jsoniter-fast-8             10000            115706 ns/op           48134 B/op       1150 allocs/op
BenchmarkWrite/jsoniter-noext
BenchmarkWrite/jsoniter-noext-8            13857             85909 ns/op           41894 B/op        517 allocs/op
```
```
goos: darwin
goarch: arm64
pkg: github.com/molon/jsoniterpb
BenchmarkRead
BenchmarkRead/protojson
BenchmarkRead/protojson-8                   3168            380545 ns/op          116098 B/op       4146 allocs/op
BenchmarkRead/jsoniter
BenchmarkRead/jsoniter-8                    5809            203376 ns/op           87670 B/op       2793 allocs/op
BenchmarkRead/jsoniter-fast
BenchmarkRead/jsoniter-fast-8               6332            193698 ns/op           87154 B/op       2780 allocs/op
BenchmarkRead/jsoniter-nofuzzydecode
BenchmarkRead/jsoniter-nofuzzydecode-8              6721            177525 ns/op           75400 B/op       2189 allocs/op
BenchmarkRead/jsoniter-fast-nofuzzydecode
BenchmarkRead/jsoniter-fast-nofuzzydecode-8         7183            166601 ns/op           75145 B/op       2178 allocs/op
```
