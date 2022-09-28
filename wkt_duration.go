package jsoniterpb

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	secondsInNanos                                  = 999999999
	maxSecondsInDuration                            = 315576000000
	Duration_message_fullname protoreflect.FullName = "google.protobuf.Duration"
)

var wktDurationCodec = (&ProtoCodec{}).
	SetElemEncodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, stream *jsoniter.Stream) {
		s, err := marshalWktDuration(((*durationpb.Duration)(ptr)))
		if err != nil {
			stream.Error = err
			return
		}
		stream.WriteVal(s)
	}).
	SetElemDecodeFunc(func(e *ProtoExtension, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		s := iter.ReadString()
		if err := unmarshalWktDuration(s, (*durationpb.Duration)(ptr)); err != nil {
			iter.ReportError("protobuf", err.Error())
			return
		}
	})

func marshalWktDuration(m *durationpb.Duration) (string, error) {
	secs := m.GetSeconds()
	nanos := m.GetNanos()
	if secs < -maxSecondsInDuration || secs > maxSecondsInDuration {
		return "", fmt.Errorf("%s: seconds out of range %v", Duration_message_fullname, secs)
	}
	if nanos < -secondsInNanos || nanos > secondsInNanos {
		return "", fmt.Errorf("%s: nanos out of range %v", Duration_message_fullname, nanos)
	}
	if (secs > 0 && nanos < 0) || (secs < 0 && nanos > 0) {
		return "", fmt.Errorf("%s: signs of seconds and nanos do not match", Duration_message_fullname)
	}
	// Generated output always contains 0, 3, 6, or 9 fractional digits,
	// depending on required precision, followed by the suffix "s".
	var sign string
	if secs < 0 || nanos < 0 {
		sign, secs, nanos = "-", -1*secs, -1*nanos
	}
	x := fmt.Sprintf("%s%d.%09d", sign, secs, nanos)
	x = strings.TrimSuffix(x, "000")
	x = strings.TrimSuffix(x, "000")
	x = strings.TrimSuffix(x, ".000")
	return x + "s", nil
}

func unmarshalWktDuration(s string, m *durationpb.Duration) error {
	secs, nanos, ok := parseDuration(s)
	if !ok {
		return fmt.Errorf("invalid %v value %q", Duration_message_fullname, s)
	}
	// Validate seconds. No need to validate nanos because parseDuration would
	// have covered that already.
	if secs < -maxSecondsInDuration || secs > maxSecondsInDuration {
		return fmt.Errorf("%v value out of range: %q", Duration_message_fullname, s)
	}

	m.Seconds = secs
	m.Nanos = nanos
	return nil
}

// parseDuration parses the given input string for seconds and nanoseconds value
// for the Duration JSON format. The format is a decimal number with a suffix
// 's'. It can have optional plus/minus sign. There needs to be at least an
// integer or fractional part. Fractional part is limited to 9 digits only for
// nanoseconds precision, regardless of whether there are trailing zero digits.
// Example values are 1s, 0.1s, 1.s, .1s, +1s, -1s, -.1s.
func parseDuration(input string) (int64, int32, bool) {
	b := []byte(input)
	size := len(b)
	if size < 2 {
		return 0, 0, false
	}
	if b[size-1] != 's' {
		return 0, 0, false
	}
	b = b[:size-1]

	// Read optional plus/minus symbol.
	var neg bool
	switch b[0] {
	case '-':
		neg = true
		b = b[1:]
	case '+':
		b = b[1:]
	}
	if len(b) == 0 {
		return 0, 0, false
	}

	// Read the integer part.
	var intp []byte
	switch {
	case b[0] == '0':
		b = b[1:]

	case '1' <= b[0] && b[0] <= '9':
		intp = b[0:]
		b = b[1:]
		n := 1
		for len(b) > 0 && '0' <= b[0] && b[0] <= '9' {
			n++
			b = b[1:]
		}
		intp = intp[:n]

	case b[0] == '.':
		// Continue below.

	default:
		return 0, 0, false
	}

	hasFrac := false
	var frac [9]byte
	if len(b) > 0 {
		if b[0] != '.' {
			return 0, 0, false
		}
		// Read the fractional part.
		b = b[1:]
		n := 0
		for len(b) > 0 && n < 9 && '0' <= b[0] && b[0] <= '9' {
			frac[n] = b[0]
			n++
			b = b[1:]
		}
		// It is not valid if there are more bytes left.
		if len(b) > 0 {
			return 0, 0, false
		}
		// Pad fractional part with 0s.
		for i := n; i < 9; i++ {
			frac[i] = '0'
		}
		hasFrac = true
	}

	var secs int64
	if len(intp) > 0 {
		var err error
		secs, err = strconv.ParseInt(string(intp), 10, 64)
		if err != nil {
			return 0, 0, false
		}
	}

	var nanos int64
	if hasFrac {
		nanob := bytes.TrimLeft(frac[:], "0")
		if len(nanob) > 0 {
			var err error
			nanos, err = strconv.ParseInt(string(nanob), 10, 32)
			if err != nil {
				return 0, 0, false
			}
		}
	}

	if neg {
		if secs > 0 {
			secs = -secs
		}
		if nanos > 0 {
			nanos = -nanos
		}
	}
	return secs, int32(nanos), true
}
