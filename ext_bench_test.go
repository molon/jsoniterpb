package jsoniterpb_test

import (
	"testing"

	gofuzz "github.com/google/gofuzz"
	jsoniter "github.com/json-iterator/go"
	"github.com/molon/jsoniterpb"
	testv1 "github.com/molon/jsoniterpb/internal/gen/go/test/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

func BenchmarkWrite(b *testing.B) {
	f := appendFuzzFuncs(gofuzz.New())
	var ms []*testv1.All
	for i := 0; i < 10000; i++ {
		var all testv1.All
		f.Fuzz(&all)
		ms = append(ms, &all)
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.Run("protojson", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m := ms[i%len(ms)]
			_, err := protojson.Marshal(m)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	cfg := jsoniter.Config{SortMapKeys: true, DisallowUnknownFields: true}.Froze()
	cfg.RegisterExtension(&jsoniterpb.ProtoExtension{})
	b.Run("jsoniter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m := ms[i%len(ms)]
			_, err := cfg.Marshal(m)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	cfg = jsoniter.Config{SortMapKeys: false, DisallowUnknownFields: false}.Froze()
	cfg.RegisterExtension(&jsoniterpb.ProtoExtension{PermitInvalidUTF8: true})
	b.Run("jsoniter-fast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m := ms[i%len(ms)]
			_, err := cfg.Marshal(m)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	cfg = jsoniter.Config{SortMapKeys: false, DisallowUnknownFields: false}.Froze()
	b.Run("jsoniter-noext", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m := ms[i%len(ms)]
			_, err := cfg.Marshal(m)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkRead(b *testing.B) {
	f := appendFuzzFuncs(gofuzz.New())
	var buffers [][]byte
	for i := 0; i < 10000; i++ {
		var all testv1.All
		f.Fuzz(&all)
		buffer, _ := protojson.Marshal(&all)
		buffers = append(buffers, buffer)
	}

	var all testv1.All
	b.ReportAllocs()
	b.ResetTimer()
	b.Run("protojson", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buffer := buffers[i%len(buffers)]
			err := protojson.Unmarshal(buffer, &all)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	cfg := jsoniter.Config{SortMapKeys: true, DisallowUnknownFields: true}.Froze()
	cfg.RegisterExtension(&jsoniterpb.ProtoExtension{})
	b.Run("jsoniter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buffer := buffers[i%len(buffers)]
			err := cfg.Unmarshal(buffer, &all)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	cfg = jsoniter.Config{SortMapKeys: false, DisallowUnknownFields: false}.Froze()
	cfg.RegisterExtension(&jsoniterpb.ProtoExtension{PermitInvalidUTF8: true})
	b.Run("jsoniter-fast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buffer := buffers[i%len(buffers)]
			err := cfg.Unmarshal(buffer, &all)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	cfg = jsoniter.Config{SortMapKeys: true, DisallowUnknownFields: true}.Froze()
	cfg.RegisterExtension(&jsoniterpb.ProtoExtension{DisableFuzzyDecode: true})
	b.Run("jsoniter-nofuzzydecode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buffer := buffers[i%len(buffers)]
			err := cfg.Unmarshal(buffer, &all)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	cfg = jsoniter.Config{SortMapKeys: false, DisallowUnknownFields: false}.Froze()
	cfg.RegisterExtension(&jsoniterpb.ProtoExtension{PermitInvalidUTF8: true, DisableFuzzyDecode: true})
	b.Run("jsoniter-fast-nofuzzydecode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buffer := buffers[i%len(buffers)]
			err := cfg.Unmarshal(buffer, &all)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
