package coze

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

const hextable = "0123456789ABCDEF"

// Type Hex is a Coze addition to Go's hex package. Coze uses Hex (Majuscule
// Hex), not lower case hex.Hex is useful for marshaling and unmarshalling
// structs that have a Hex type. Hex's underlying type is []byte and is
// represented in JSON as upper case Hexadecimal.  Hex should always be padded
// to the two time the number of bytes; there's no truncation of left padding
// characters.
//
// When converting integers or other types, `nil` in Hex is "" and a zero type
// (for instance "0") in Hex is encoded as "00".
//
type Hex []byte

// UnmarshalJSON implements JSON.UnmarshalJSON.  It is a custom unmarshaler for
// binary data in JSON, which should always be represented as Hex.
func (h *Hex) UnmarshalJSON(b []byte) error {
	// JSON.Unmarshal gives b encapsulated in quote characters. Quotes characters
	// are invalid Hex and must be stripped.
	trimmed := strings.Trim(string(b), "\"")
	s, err := hex.DecodeString(trimmed)
	if err != nil {
		return err
	}

	*h = Hex(s)
	return nil
}

// MarshalJSON implements JSON.UnmarshalJSON. Converts bytes to Hex.  Error is
// always nil.
func (t Hex) MarshalJSON() ([]byte, error) {
	// JSON expects stings to be wrapped with double quote character.
	return []byte(fmt.Sprintf("\"%v\"", t)), nil
}

// String implements fmt.Stringer. Use `%s`, `%v`, `%+v` to get this form.
func (t Hex) String() string {
	return fmt.Sprintf("%X", []byte(t))
}

// GoString implements fmt.GoString. Use `%#v` to get this form (not %s or %+v).
func (t Hex) GoString() string {
	// // Base256 representation
	// return fmt.Sprintf("%s", []byte(t))
	// Base64 Representation
	return base64.StdEncoding.EncodeToString([]byte(t))
}

// HexEncode converts bytes to Hex.
func HexEncode(b []byte) string {
	return Hex(b).String()
}

// HexDecode converts a Hex string to bytes.
func HexDecode(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

// MustHexDecode panics on error.
func MustHexDecode(s string) []byte {
	b, err := HexDecode(s)
	if err != nil {
		panic(err)
	}
	return b
}
