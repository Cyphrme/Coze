package coze

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Type B64 is a Coze addition to Go's base64. B64 is useful for marshaling and
// unmarshalling structs. B64's underlying type is []byte and is represented in
// JSON as "RFC 4648 base 64 URI canonical with padding truncated" (b64ut).
//
// When converting integers or other types to B64, `nil` is encoded as "" and
// zero is encoded as "AA".
type B64 []byte

// UnmarshalJSON implements JSON.UnmarshalJSON so B64 is encoded as b64ut.
func (t *B64) UnmarshalJSON(b []byte) error {
	// JSON.Unmarshal returns b encapsulated in quotes which is invalid base64 characters.
	s, err := base64.URLEncoding.Strict().WithPadding(base64.NoPadding).DecodeString(strings.Trim(string(b), "\""))
	if err != nil {
		return err
	}
	*t = B64(s)
	return nil
}

// MarshalJSON implements JSON.MarshalJSON so B64 is decoded from b64ut. Error
// is always nil.
func (t B64) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", t)), nil
}

// String implements fmt.Stringer. Use with `%s`, `%v`, `%+v`.
func (t B64) String() string {
	return base64.URLEncoding.Strict().WithPadding(base64.NoPadding).EncodeToString([]byte(t))
}

// GoString implements fmt.GoStringer. Use with `%#v` (not %s or %+v).
func (t B64) GoString() string {
	return base64.URLEncoding.Strict().WithPadding(base64.NoPadding).EncodeToString([]byte(t))
}

// Decode decodes a b64ut string.
func Decode(b64 string) (B64, error) {
	return base64.URLEncoding.Strict().WithPadding(base64.NoPadding).DecodeString(b64)
}

// MustDecode decodes b64ut and panics on error.
func MustDecode(b64 string) B64 {
	b, err := base64.URLEncoding.Strict().WithPadding(base64.NoPadding).DecodeString(b64)
	if err != nil {
		panic(err)
	}
	return b
}

// sB64 is useful for map keys since currently in Go []byte is not allowed to be
// a map key but string is.  The type really should be []byte and not string,
// but Go does not yet support this. See https://github.com/golang/go/issues/283
// and https://github.com/google/go-cmp/issues/67.  SB64 will be deprecated
// if/when Go supports []byte keys.
//
// From https://go.dev/blog/strings >[A] string holds arbitrary bytes. It is not
// required to hold Unicode text, UTF-8 text, or any other predefined format. As
// far as the content of a string is concerned, it is exactly equivalent to a
// slice of bytes.
type SB64 string

// String implements fmt.Stringer
func (b SB64) String() string {
	return B64(b).String()
}

// GoString implements fmt.GoString
func (b SB64) GoString() string {
	return b.String()
}
