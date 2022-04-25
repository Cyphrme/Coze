package coze

import (
	"encoding/json"
	"sort"

	ce "github.com/cyphrme/coze/enum"
)

// Canon returns the canonical form of input. See docs on Canonical().
func Canon(input interface{}, canon interface{}) (b []byte, err error) {
	b, err = Marshal(input)
	if err != nil {
		return nil, err
	}

	return Canonical(b, canon)
}

// Canonical returns the canonical form of input.
//
// Interface "canon" may be any valid type for json.Unmarshal, including
// `[]string` and `nil`.  If canon is nil, json.Unmarshal will place the input
// into a UTF-8 ordered map.
//
// If "canon" is a struct the struct must be properly ordered. Go's
// JSON package orders struct fields according to their struct position.
//
// In the Go version of Coze, the canonical form of a struct is achieved by
// unmarshalling and remarshaling.
func Canonical(input []byte, canon interface{}) (b []byte, err error) {
	// Unmarshal the given bytes into the given canonical format.
	err = json.Unmarshal(input, &canon)
	if err != nil {
		return nil, err
	}

	return Marshal(canon)
}

// CanonS generates a canon, the UTF-8 sorted field names, from the given
// struct. It returns only top level fields with no recursion and no promoting
// of embedded fields.
func CanonS(structure interface{}) (can []string, err error) {
	m, err := Marshal(structure)
	if err != nil {
		return nil, err
	}

	return CanonB(m)
}

// CanonB returns a sorted canon from a byte slice.  It returns only top level
// fields with no recursion and no promoting of embedded fields.
func CanonB(b []byte) (can []string, err error) {
	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	v := make([]string, 0, len(m))
	for s := range m {
		v = append(v, s)
	}
	sort.Strings(v)

	return v, nil
}

// CH (Canonical Hash) accepts []byte and an optional canon, and returns digest
// of the canonical form.
// If input is already in canonical form, cozeenum.Hash() can be called instead.
func CH(input []byte, canon interface{}, hash ce.HashAlg) (digest Hex, err error) {
	b, err := Canonical(input, canon)
	if err != nil {
		return nil, err
	}

	digest = ce.Hash(hash, b)
	return
}
