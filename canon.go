package coze

import (
	"bytes"
	"encoding/json"
)

// Canon returns the current canon from raw JSON.
//
// It returns only top level fields with no recursion or promotion of embedded
// fields.
func Canon(raw json.RawMessage) (can []string, err error) {
	ms := MapSlice{}
	err = json.Unmarshal(raw, &ms)
	if err != nil {
		return nil, err
	}
	can = make([]string, len(ms))
	for i, v := range ms {
		can[i] = v.Key
	}
	return can, nil
}

// Canonical returns the canonical form. Input canon is optional and may be nil.
// If canon is nil, input JSON is only compactified.
//
// Interface "canon" may be `[]string`, `struct“, or `nil`.  If "canon" is a
// struct or slice it must be properly ordered.  If canon is nil, json.Unmarshal
// will place the input into a UTF-8 ordered map.
//
// In the Go version of Coze, the canonical form of a struct is (currently)
// achieved by unmarshalling and remarshaling.
func Canonical(input []byte, canon any) (b []byte, err error) {
	if canon == nil {
		return compact(input)
	}

	s, ok := canon.([]string)
	if ok {
		// The only datastructure that can unmarshal arbitrary JSON is map, but
		// json.Marshal will unmarshal *all* elements and there is no way to specify
		// unmarshalling to only the given fields.  Solution: unmarshal into new
		// map, and transfer appropriate fields to a second map.
		m := make(map[string]any)
		err = json.Unmarshal(input, &m)
		if err != nil {
			return nil, err
		}

		mm := make(map[string]any)
		for i := 0; i < len(s); i++ {
			mm[s[i]] = m[s[i]]
		}

		return Marshal(mm)
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
	return Hash(hash, input)
}

// Compact compactifies JSON.
func compact(msg json.RawMessage) ([]byte, error) {
	var b bytes.Buffer
	err := json.Compact(&b, msg)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
