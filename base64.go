package coze

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Type B64 is a Coze addition to Go's base64. B64 is useful for marshaling and
// unmarshalling structs. B64's underlying type is []byte and is represented in
// JSON as base64 URI truncated (b64ut).
//
// When converting integers or other types, `nil` in B64 is "" and the zero is
// encoded as "AA".
type B64 []byte

// UnmarshalJSON implements JSON.UnmarshalJSON. It is a custom unmarshaler for
// binary data in JSON, which should always be represented as Hex.
func (t *B64) UnmarshalJSON(b []byte) error {
	// JSON.Unmarshal gives b encapsulated in quote characters. Quotes characters
	// are invalid base64 and must be stripped.
	s, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(strings.Trim(string(b), "\""))
	if err != nil {
		return err
	}
	*t = B64(s)
	return nil
}

// MarshalJSON implements JSON.UnmarshalJSON. Converts bytes to Hex.  Error is
// always nil.
func (t B64) MarshalJSON() ([]byte, error) {
	// JSON expects stings to be wrapped with double quote character.
	return []byte(fmt.Sprintf("\"%v\"", t)), nil
}

// String implements fmt.Stringer. Use with `%s`, `%v`, `%+v`.
func (t B64) String() string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(t))
}

// GoString implements fmt.GoStringer. Use with `%#v` (not %s or %+v).
func (t B64) GoString() string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(t))
}

// Decode decodes a base64 string to B64.
func Decode(b64 string) (B64, error) {
	return base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(b64)
}

// MustDecode decodes a base64 string to B64.  Will panic on error.
func MustDecode(b64 string) B64 {
	b, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(b64)
	if err != nil {
		panic(err)
	}
	return b
}
