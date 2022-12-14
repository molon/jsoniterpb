// Code generated by protoc-gen-gofuzz. DO NOT EDIT.
//
// Source: test/v1/test.proto

package testv1fuzz

import (
	gofuzz "github.com/google/gofuzz"
	v1 "github.com/molon/jsoniterpb/internal/gen/go/test/v1"
)

// FuzzMessage is a fuzz function.
// If can be registered using `Fuzzer.Funcs` function.
func FuzzMessage(x *v1.Message, f gofuzz.Continue) {
	f.Fuzz(&x.Id)
}

// FuzzCaseValue is a fuzz function.
// If can be registered using `Fuzzer.Funcs` function.
func FuzzCaseValue(x *v1.CaseValue, f gofuzz.Continue) {
	f.Fuzz(&x.V)
	f.Fuzz(&x.A)
}

// FuzzCase is a fuzz function.
// If can be registered using `Fuzzer.Funcs` function.
func FuzzCase(x *v1.Case, f gofuzz.Continue) {
	f.Fuzz(&x.WktI32A)
	f.Fuzz(&x.WktI32B)
	f.Fuzz(&x.OptI32A)
	f.Fuzz(&x.OptI32B)
	f.Fuzz(&x.OptWktI32A)
	f.Fuzz(&x.OptWktI32B)
	f.Fuzz(&x.RptWktI32)
	f.Fuzz(&x.MapWktI32)
	f.Fuzz(&x.B1)
	f.Fuzz(&x.B2)
	f.Fuzz(&x.OptB1)
	f.Fuzz(&x.OptB2)
	f.Fuzz(&x.RptB)
	f.Fuzz(&x.MapB)
	f.Fuzz(&x.WktB1)
	f.Fuzz(&x.WktB2)
	f.Fuzz(&x.OptWktB1)
	f.Fuzz(&x.OptWktB2)
	f.Fuzz(&x.RptWktB)
	f.Fuzz(&x.MapWktB)
	f.Fuzz(&x.RptMsg)
	f.Fuzz(&x.MapMsg)
	f.Fuzz(&x.MapEnum)
	f.Fuzz(&x.MapWktU64)
	f.Fuzz(&x.WktV)
	f.Fuzz(&x.WktS)
	f.Fuzz(&x.WktLv)
	f.Fuzz(&x.RptWktV)
	f.Fuzz(&x.RptWktS)
	f.Fuzz(&x.RptWktLv)
	switch f.Int31n(4) {
	case 0:
		var o v1.Case_OneofWktI32
		f.Fuzz(&o.OneofWktI32)
		x.OneOf = &o
	case 1:
		var o v1.Case_OneofB
		f.Fuzz(&o.OneofB)
		x.OneOf = &o
	case 2:
		var o v1.Case_OneofWktB
		f.Fuzz(&o.OneofWktB)
		x.OneOf = &o
	}
}

// FuzzAll is a fuzz function.
// If can be registered using `Fuzzer.Funcs` function.
func FuzzAll(x *v1.All, f gofuzz.Continue) {
	f.Fuzz(&x.R)
	f.Fuzz(&x.S)
	f.Fuzz(&x.OF)
	f.Fuzz(&x.E)
	f.Fuzz(&x.OWkt)
	f.Fuzz(&x.Wkt)
	f.Fuzz(&x.O)
	f.Fuzz(&x.RWkt)
	f.Fuzz(&x.N)
	f.Fuzz(&x.M)
	f.Fuzz(&x.SnakeCase)
	f.Fuzz(&x.LowerCamelCase)
	f.Fuzz(&x.UpwerCamelCase)
	f.Fuzz(&x.OptWkt)
}

// FuzzRepeated is a fuzz function.
// If can be registered using `Fuzzer.Funcs` function.
func FuzzRepeated(x *v1.Repeated, f gofuzz.Continue) {
	f.Fuzz(&x.S)
	f.Fuzz(&x.I32)
	f.Fuzz(&x.I64)
	f.Fuzz(&x.U32)
	f.Fuzz(&x.U64)
	f.Fuzz(&x.F32)
	f.Fuzz(&x.F64)
	f.Fuzz(&x.Si32)
	f.Fuzz(&x.Si64)
	f.Fuzz(&x.Fi32)
	f.Fuzz(&x.Fi64)
	f.Fuzz(&x.Sfi32)
	f.Fuzz(&x.Sfi64)
	f.Fuzz(&x.Bl)
	f.Fuzz(&x.By)
	f.Fuzz(&x.E)
	f.Fuzz(&x.Msg)
}

// FuzzOptionals is a fuzz function.
// If can be registered using `Fuzzer.Funcs` function.
func FuzzOptionals(x *v1.Optionals, f gofuzz.Continue) {
	f.Fuzz(&x.Id)
	f.Fuzz(&x.I32)
	f.Fuzz(&x.I64)
	f.Fuzz(&x.U32)
	f.Fuzz(&x.U64)
	f.Fuzz(&x.F32)
	f.Fuzz(&x.F64)
	f.Fuzz(&x.Si32)
	f.Fuzz(&x.Si64)
	f.Fuzz(&x.Fi32)
	f.Fuzz(&x.Fi64)
	f.Fuzz(&x.Sfi32)
	f.Fuzz(&x.Sfi64)
	f.Fuzz(&x.Bl)
	f.Fuzz(&x.By)
	f.Fuzz(&x.S)
	f.Fuzz(&x.E)
}

// FuzzWKTs is a fuzz function.
// If can be registered using `Fuzzer.Funcs` function.
func FuzzWKTs(x *v1.WKTs, f gofuzz.Continue) {
	f.Fuzz(&x.A)
	f.Fuzz(&x.D)
	f.Fuzz(&x.T)
	f.Fuzz(&x.St)
	f.Fuzz(&x.I32)
	f.Fuzz(&x.Ui32)
	f.Fuzz(&x.I64)
	f.Fuzz(&x.U64)
	f.Fuzz(&x.F32)
	f.Fuzz(&x.F64)
	f.Fuzz(&x.B)
	f.Fuzz(&x.S)
	f.Fuzz(&x.By)
	f.Fuzz(&x.Fm)
	f.Fuzz(&x.Em)
	f.Fuzz(&x.Nu)
	f.Fuzz(&x.V)
	f.Fuzz(&x.Lv)
}

// FuzzRepeatedWKTs is a fuzz function.
// If can be registered using `Fuzzer.Funcs` function.
func FuzzRepeatedWKTs(x *v1.RepeatedWKTs, f gofuzz.Continue) {
	f.Fuzz(&x.A)
	f.Fuzz(&x.D)
	f.Fuzz(&x.T)
	f.Fuzz(&x.St)
	f.Fuzz(&x.I32)
	f.Fuzz(&x.Ui32)
	f.Fuzz(&x.I64)
	f.Fuzz(&x.U64)
	f.Fuzz(&x.F32)
	f.Fuzz(&x.F64)
	f.Fuzz(&x.B)
	f.Fuzz(&x.S)
	f.Fuzz(&x.By)
	f.Fuzz(&x.Fm)
	f.Fuzz(&x.Em)
	f.Fuzz(&x.Nu)
	f.Fuzz(&x.V)
	f.Fuzz(&x.Lv)
}

// FuzzOptionalWKTs is a fuzz function.
// If can be registered using `Fuzzer.Funcs` function.
func FuzzOptionalWKTs(x *v1.OptionalWKTs, f gofuzz.Continue) {
	f.Fuzz(&x.A)
	f.Fuzz(&x.D)
	f.Fuzz(&x.T)
	f.Fuzz(&x.St)
	f.Fuzz(&x.I32)
	f.Fuzz(&x.Ui32)
	f.Fuzz(&x.I64)
	f.Fuzz(&x.U64)
	f.Fuzz(&x.F32)
	f.Fuzz(&x.F64)
	f.Fuzz(&x.B)
	f.Fuzz(&x.S)
	f.Fuzz(&x.By)
	f.Fuzz(&x.Fm)
	f.Fuzz(&x.Em)
	f.Fuzz(&x.Nu)
	f.Fuzz(&x.V)
	f.Fuzz(&x.Lv)
}

// FuzzOneOf is a fuzz function.
// If can be registered using `Fuzzer.Funcs` function.
func FuzzOneOf(x *v1.OneOf, f gofuzz.Continue) {
	f.Fuzz(&x.Extra)
	switch f.Int31n(18) {
	case 0:
		var o v1.OneOf_E
		f.Fuzz(&o.E)
		x.OneOf = &o
	case 1:
		var o v1.OneOf_STr
		f.Fuzz(&o.STr)
		x.OneOf = &o
	case 2:
		var o v1.OneOf_I32
		f.Fuzz(&o.I32)
		x.OneOf = &o
	case 3:
		var o v1.OneOf_I64
		f.Fuzz(&o.I64)
		x.OneOf = &o
	case 4:
		var o v1.OneOf_U32
		f.Fuzz(&o.U32)
		x.OneOf = &o
	case 5:
		var o v1.OneOf_U64
		f.Fuzz(&o.U64)
		x.OneOf = &o
	case 6:
		var o v1.OneOf_F32
		f.Fuzz(&o.F32)
		x.OneOf = &o
	case 7:
		var o v1.OneOf_F64
		f.Fuzz(&o.F64)
		x.OneOf = &o
	case 8:
		var o v1.OneOf_Si32
		f.Fuzz(&o.Si32)
		x.OneOf = &o
	case 9:
		var o v1.OneOf_Si64
		f.Fuzz(&o.Si64)
		x.OneOf = &o
	case 10:
		var o v1.OneOf_Fi32
		f.Fuzz(&o.Fi32)
		x.OneOf = &o
	case 11:
		var o v1.OneOf_Fi64
		f.Fuzz(&o.Fi64)
		x.OneOf = &o
	case 12:
		var o v1.OneOf_Sfi32
		f.Fuzz(&o.Sfi32)
		x.OneOf = &o
	case 13:
		var o v1.OneOf_Sfi64
		f.Fuzz(&o.Sfi64)
		x.OneOf = &o
	case 14:
		var o v1.OneOf_Bl
		f.Fuzz(&o.Bl)
		x.OneOf = &o
	case 15:
		var o v1.OneOf_By
		f.Fuzz(&o.By)
		x.OneOf = &o
	case 16:
		var o v1.OneOf_Msg
		f.Fuzz(&o.Msg)
		x.OneOf = &o
	}
}

// FuzzOneOfWKT is a fuzz function.
// If can be registered using `Fuzzer.Funcs` function.
func FuzzOneOfWKT(x *v1.OneOfWKT, f gofuzz.Continue) {
	switch f.Int31n(19) {
	case 0:
		var o v1.OneOfWKT_A
		f.Fuzz(&o.A)
		x.OneOf = &o
	case 1:
		var o v1.OneOfWKT_D
		f.Fuzz(&o.D)
		x.OneOf = &o
	case 2:
		var o v1.OneOfWKT_T
		f.Fuzz(&o.T)
		x.OneOf = &o
	case 3:
		var o v1.OneOfWKT_St
		f.Fuzz(&o.St)
		x.OneOf = &o
	case 4:
		var o v1.OneOfWKT_I32
		f.Fuzz(&o.I32)
		x.OneOf = &o
	case 5:
		var o v1.OneOfWKT_Ui32
		f.Fuzz(&o.Ui32)
		x.OneOf = &o
	case 6:
		var o v1.OneOfWKT_I64
		f.Fuzz(&o.I64)
		x.OneOf = &o
	case 7:
		var o v1.OneOfWKT_U64
		f.Fuzz(&o.U64)
		x.OneOf = &o
	case 8:
		var o v1.OneOfWKT_F32
		f.Fuzz(&o.F32)
		x.OneOf = &o
	case 9:
		var o v1.OneOfWKT_F64
		f.Fuzz(&o.F64)
		x.OneOf = &o
	case 10:
		var o v1.OneOfWKT_B
		f.Fuzz(&o.B)
		x.OneOf = &o
	case 11:
		var o v1.OneOfWKT_S
		f.Fuzz(&o.S)
		x.OneOf = &o
	case 12:
		var o v1.OneOfWKT_By
		f.Fuzz(&o.By)
		x.OneOf = &o
	case 13:
		var o v1.OneOfWKT_Fm
		f.Fuzz(&o.Fm)
		x.OneOf = &o
	case 14:
		var o v1.OneOfWKT_Em
		f.Fuzz(&o.Em)
		x.OneOf = &o
	case 15:
		var o v1.OneOfWKT_Nu
		f.Fuzz(&o.Nu)
		x.OneOf = &o
	case 16:
		var o v1.OneOfWKT_V
		f.Fuzz(&o.V)
		x.OneOf = &o
	case 17:
		var o v1.OneOfWKT_Lv
		f.Fuzz(&o.Lv)
		x.OneOf = &o
	}
}

// FuzzSingular is a fuzz function.
// If can be registered using `Fuzzer.Funcs` function.
func FuzzSingular(x *v1.Singular, f gofuzz.Continue) {
	f.Fuzz(&x.E)
	f.Fuzz(&x.S)
	f.Fuzz(&x.I32)
	f.Fuzz(&x.I64)
	f.Fuzz(&x.U32)
	f.Fuzz(&x.U64)
	f.Fuzz(&x.F32)
	f.Fuzz(&x.F64)
	f.Fuzz(&x.Si32)
	f.Fuzz(&x.Si64)
	f.Fuzz(&x.Fi32)
	f.Fuzz(&x.Fi64)
	f.Fuzz(&x.Sfi32)
	f.Fuzz(&x.Sfi64)
	f.Fuzz(&x.Bl)
	f.Fuzz(&x.By)
	f.Fuzz(&x.Msg)
}

// FuzzMap is a fuzz function.
// If can be registered using `Fuzzer.Funcs` function.
func FuzzMap(x *v1.Map, f gofuzz.Continue) {
	f.Fuzz(&x.En)
	f.Fuzz(&x.Msg)
	f.Fuzz(&x.Str)
	f.Fuzz(&x.By)
	f.Fuzz(&x.Bo)
	f.Fuzz(&x.An)
	f.Fuzz(&x.Bn)
}

// FuzzNested is a fuzz function.
// If can be registered using `Fuzzer.Funcs` function.
func FuzzNested(x *v1.Nested, f gofuzz.Continue) {
	f.Fuzz(&x.N)
}

// FuzzNested_NestedMessage is a fuzz function.
// If can be registered using `Fuzzer.Funcs` function.
func FuzzNested_NestedMessage(x *v1.Nested_NestedMessage, f gofuzz.Continue) {
	f.Fuzz(&x.E)
}

// FuzzNested_NestedMessage_NestedEnum is a fuzz function.
// If can be registered using `Fuzzer.Funcs` function.
func FuzzNested_NestedMessage_NestedEnum(x *v1.Nested_NestedMessage_NestedEnum, f gofuzz.Continue) {
	switch f.Int31n(1) {
	case 0:
		*x = v1.Nested_NestedMessage_NESTED_ENUM_UNSPECIFIED
	}
}

// FuzzJsonEnum is a fuzz function.
// If can be registered using `Fuzzer.Funcs` function.
func FuzzJsonEnum(x *v1.JsonEnum, f gofuzz.Continue) {
	switch f.Int31n(2) {
	case 0:
		*x = v1.JsonEnum_JSON_ENUM_UNSPECIFIED
	case 1:
		*x = v1.JsonEnum_JSON_ENUM_SOME
	}
}
