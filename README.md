# jsoniterpb
Replacement of Protojson over Jsoniter Extension

#### Features
- Most of the features are the same as protojson
- Support more fuzzy decode methods
- Support for objects that do not implement `proto.Message`
- better performance

#### Warns
- Only `proto3` is supported, `proto2` is not supported
- `protojson` marshal nil `proto.Message` as zero value **if it is root**. but `jsoniterpb` will marshal it to `null`

- View [internal/protojson/tests/jsoniterpb_decode_test.go](internal/protojson/tests/jsoniterpb_decode_test.go)
  - Support more fuzzy decode methods => Search `FuzzyDecode`
  - Some error messages are not the same => Search `ErrMsgNotSame`
  - Some error check are not supported => Search `NotSupport`

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

#### Benchmark
```
BenchmarkWrite/protojson
BenchmarkWrite/protojson-8         	    4759	    251567 ns/op	  121746 B/op	    2380 allocs/op
BenchmarkWrite/jsoniter
BenchmarkWrite/jsoniter-8          	    6506	    177602 ns/op	   88752 B/op	    2322 allocs/op
BenchmarkWrite/jsoniter-fast
BenchmarkWrite/jsoniter-fast-8     	   10000	    112298 ns/op	   48103 B/op	    1150 allocs/op
```
```
BenchmarkRead/protojson
BenchmarkRead/protojson-8                   3120            384538 ns/op          115852 B/op       4137 allocs/op
BenchmarkRead/jsoniter
BenchmarkRead/jsoniter-8                    5943            200014 ns/op           87369 B/op       2787 allocs/op
BenchmarkRead/jsoniter-fast
BenchmarkRead/jsoniter-fast-8               6291            193630 ns/op           87441 B/op       2781 allocs/op
BenchmarkRead/jsoniter-nofuzzydecode
BenchmarkRead/jsoniter-nofuzzydecode-8              6241            174049 ns/op           74271 B/op       2167 allocs/op
BenchmarkRead/jsoniter-fast-nofuzzydecode
BenchmarkRead/jsoniter-fast-nofuzzydecode-8         6708            166749 ns/op           74963 B/op       2172 allocs/op
```