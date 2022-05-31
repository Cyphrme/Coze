package coze

import (
	"encoding/json"
	"sort"

	ce "github.com/cyphrme/coze/enum"
)

// Canon returns a canon from raw JSON.
//
// It returns only top level fields with no recursion or promotion of embedded
// fields.
func Canon(raw json.RawMessage) (can []string, err error) {
	// In Go, map order is unspecified and package json cannot unmarshal into
	// array/slice.  Unmarshal into a map, put keys into a slice, and sort.
	var m map[string]interface{}
	err = json.Unmarshal(raw, &m)
	if err != nil {
		return nil, err
	}

	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys, nil
}

// CanonStruct generates a canon from the given struct.
//
//It returns only top level fields with no recursion and no promoting
// of embedded fields.
func CanonStruct(structure interface{}) (can []string, err error) {
	m, err := Marshal(structure)
	if err != nil {
		return nil, err
	}

	return Canon(m)
}

// Canonical returns the compactified and/or canonical form. Input canon may be
// nil. If canon is nil, JSON is only compactified.
//
// Interface "canon" may be any valid type for json.Unmarshal, including
// `[]string`, `struct``, and `nil`.  If canon is nil, json.Unmarshal will place
// the input into a UTF-8 ordered map.
//
// If "canon" is a struct the struct must be properly ordered. Go's JSON package
// orders struct fields according to their struct position.
//
// In the Go version of Coze, the canonical form of a struct is (currently)
// achieved by unmarshalling and remarshaling.
func Canonical(input []byte, canon interface{}) (b []byte, err error) {
	// Unmarshal the given bytes into the given canonical format.
	err = json.Unmarshal(input, &canon)
	if err != nil {
		return nil, err
	}

	return Marshal(canon)
}

// CanonicalStruct returns the canonical form of a struct.  See notes on Canonical.
func CanonicalStruct(structure interface{}, canon interface{}) (b []byte, err error) {
	m, err := Marshal(structure)
	if err != nil {
		return nil, err
	}

	return Canonical(m, canon)
}

// CanonHash accepts []byte and an optional canon, and returns digest
// of the canonical form.
//
// If input is already in canonical form, enum.Hash() can be called instead.
func CanonHash(input []byte, canon interface{}, hash ce.HashAlg) (digest B64, err error) {
	b, err := Canonical(input, canon)
	if err != nil {
		return nil, err
	}

	return ce.Hash(hash, b), nil
}
