package coze

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const hextable = "0123456789ABCDEF"

// Type Hex is a Coze addition to Go's hex package. Hex is useful for marshaling
// and unmarshalling structs that have a Hex type. Hex underlying type is
// []byte, and represented in JSON as upper case Hexadecimal.  Hex should always
// be padded to the "full size" byte representation, meaning there's no
// truncation of left "padding characters"
//
// When converting integers or other types, `nil`` in Hex is "" and a zero type
// (for instance "0") in Hex is encoded as "00".
//
// Cyphr.me uses Majuscule Hex, not lower case hex.
type Hex []byte

// UnmarshalJSON implements JSON.UnmarshalJSON.  It is a custom unmarshaler for
// binary data in JSON, which should always be represented as upper case hex.
func (h *Hex) UnmarshalJSON(b []byte) error {
	// JSON.Unmarshal will send b encapsulated in quote characters. Quotes
	// characters are invalid hex and need to be stripped.
	trimmed := strings.Trim(string(b), "\"")
	s, err := hex.DecodeString(trimmed)
	if err != nil {
		return err
	}

	*h = Hex(s)
	return nil
}

// MarshalJSON implements JSON.UnmarshalJSON. Converts bytes to upper case Hex
// string.  []byte(string(Hex))) error is always nil.
func (t Hex) MarshalJSON() ([]byte, error) {
	// JSON expects stings to be wrapped with double quote character.
	return []byte(fmt.Sprintf("\"%v\"", t)), nil
}

// String implements fmt.Stringer.
// use `%s`, `%v`, `%+v` to get this form.
func (t Hex) String() string {
	return fmt.Sprintf("%X", []byte(t))
}

// String implements fmt.GoString.
// Use `%#v` to get this form (not %s or %+v).
func (t Hex) GoString() string {
	// // Base256 representation
	// return fmt.Sprintf("%s", []byte(t))
	// Base64 Representation
	return base64.StdEncoding.EncodeToString([]byte(t))
}

// HexEncodeString converts bytes to Hex string.
func HexEncodeString(b []byte) string {
	return Hex(b).String()
}

// HexDecodeString is a convenience function
func HexDecodeString(s string) ([]byte, error) {
	// if len(s)%2 != 0 {
	// 	s = "0" + s
	// }
	return hex.DecodeString(s)
}

// MustHexDecode panics on error.
func MustHexDecode(s string) []byte {
	b, err := HexDecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

// HexEncodeStringPadded returns the padded hexadecimal encoding of src with
// given length of pad.
func HexEncodeStringPadded(src []byte, pad int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(pad)+"s", HexEncodeString(src))
}
