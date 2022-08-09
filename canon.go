package coze

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Canon returns the current canon from raw JSON.
//
// It returns only top level fields with no recursion or promotion of embedded
// fields.
func GetCanon(raw json.RawMessage) (can []string, err error) {
	ms := MapSlice{}
	err = json.Unmarshal(raw, &ms)
	if err != nil {
		return nil, err
	}

	keys := make([]string, len(ms))
	for i, v := range ms {
		keys[i] = fmt.Sprintf("%v", v.Key)
	}
	return keys, nil
}

// Canonical returns the canonical form. Input canon may be nil. If canon is
// nil, JSON is only compactified.
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
func Canonical(input []byte, canon any) (b []byte, err error) {
	if canon == nil {
		return compact(input)
	}

	// Unmarshal the given bytes into the given canonical format.
	err = json.Unmarshal(input, &canon)
	if err != nil {
		return nil, err
	}
	return Marshal(canon)
}

// CanonicalHash accepts []byte and optional canon and returns digest.
//
// If input is already in canonical form, Hash() may also be called instead.
func CanonicalHash(input []byte, canon any, hash HashAlg) (digest B64, err error) {
	input, err = Canonical(input, canon)
	if err != nil {
		return nil, err
	}

	return Hash(hash, input), nil
}

func compact(msg json.RawMessage) ([]byte, error) {
	var b bytes.Buffer
	err := json.Compact(&b, msg)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
