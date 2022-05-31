// See `coze.md` for details.
package coze

import (
	"bytes"
	"encoding/json"

	ce "github.com/cyphrme/coze/enum"
)

// Head contains the standard fields in a signed Coze object.
// TODO prehaps create a RVK struct.
type Head struct {
	Alg ce.SEAlg `json:"alg,omitempty"` // e.g. "ES256"
	Iat int64    `json:"iat,omitempty"` // e.g. 1623132000
	Tmb B64      `json:"tmb,omitempty"` // e.g. "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk"
	Typ string   `json:"typ,omitempty"` // e.g. "cyphr.me/msg/create"
}

// CozeMarshaler is a UTF-8 marshaler for Go structs.  Go's `json.Marshal`
// removes the valid characters "&". "<", ">".  See note on Marshal.
type CozeMarshaler interface {
	CozeMarshal() ([]byte, error)
}

// Marshal is a UTF-8 friendly marshaler.  Go's json.Marshal is not UTF-8
// friendly because it replaces the valid UTF-8 and JSON characters "&". "<",
// ">" with the "slash u" unicode escaped forms (e.g. \u0026).  It preemptively
// escapes for HTML friendliness.  Where JSON may include these characters,
// json.Marshal should not be used. Playground of Go breaking a book title:
// https://play.golang.org/p/o2hiX0c62oN
func Marshal(i interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(i)
	if err != nil {
		return nil, err
	}
	return bytes.TrimRight(buffer.Bytes(), "\n"), nil
}

// MarshalPretty is the pretty version of Marshal. It uses 4 spaces for each
// level.  Spaces instead of tabs because some applications use 8 spaces per
// tab, which is excessive.
func MarshalPretty(i interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "    ")
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(i)
	if err != nil {
		return nil, err
	}
	return bytes.TrimRight(buffer.Bytes(), "\n"), nil
}
