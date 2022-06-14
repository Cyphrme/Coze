package coze

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Type B64 is a Coze addition to Go's base64.  B64 is useful for marshaling and
// unmarshalling structs. B64's underlying type is []byte and is represented in
// JSON as base64 URI truncated (b64ut).
//
// When converting integers or other types, `nil` in B64 is "" and the zero is
// encoded as "AA".
//
// To unencode from string, use the base64 package:
// base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(b64String)
type B64 []byte

// UnmarshalJSON implements JSON.UnmarshalJSON.  It is a custom unmarshaler for
// binary data in JSON, which should always be represented as Hex.
func (t *B64) UnmarshalJSON(b []byte) error {
	// JSON.Unmarshal gives b encapsulated in quote characters. Quotes characters
	// are invalid base64 and must be stripped.
	trimmed := strings.Trim(string(b), "\"")
	s, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(trimmed)
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

// String implements fmt.Stringer. Use `%s`, `%v`, `%+v` to get this form.
func (t B64) String() string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(t))

}

// GoString implements fmt.GoString. Use `%#v` to get this form (not %s or %+v).
func (t B64) GoString() string {
	// Base256 representation
	// return fmt.Sprintf("%s", []byte(t))
	return fmt.Sprintf("%X", []byte(t))
}

// MustDecode decodes a base64 string to B64.  Will panic on error.
func MustDecode(b64 string) B64 {
	b, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(b64)
	if err != nil {
		panic(err)
	}
	return b
}
