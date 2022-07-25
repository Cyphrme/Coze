package coze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"golang.org/x/exp/slices"
)

// Normal - A normal is an arrays of fields specifying the normalization of
// payload.  Normals are implemented in Go as []string.  There are six types of
// normals, and a nil normal is valid.
//
//  canon       (can)
//  only        (ony)
//  need        (ned)
//  order       (ord)
//  option      (opt)
//  extra       (ext)
//  (nil)
//
// An nil normal matches all payloads.
//
// `canon` requires specified fields in the given order and no extra fields
// permitted.
//
// `only` specifies fields that are required to be present, does not specify any
// order, and no extra fields permitted.
//
// `need` specifies fields that are required to be present, but does not specify
// any order. Extra fields are permitted after the `need` fields.
//
// `order` requires specified fields in the given order and extra fields
// are permitted after the `order` fields.
//

//
// `extra` specifies optional fields in a given order and allows extra fields
// are permitted after the extra fields.
//
//
// A an empty canon of [] matches only an a empty (i.e. `{}`) payload.
//
// # Using with Option
// When a need is used with an option all fields are unordered.
//
// When an order is used with an option, all fields are ordered.
//
// ## Normal, Require, and Option
//
// `canon`, `only`, `need`, and `order` are valid `require` in that they specify
// required fields.  An option is distinct in that option specifies optional
// fields and precludes other optional fields.
//
//              ┌────────────────┐
//              │     Normal     │
//              └───────┬────────┘
//              ┌───────┴────────┐
//        ┌─────┴────┐     ┌─────┴──────┐
//        │ Require  │     │   Option   │
//        └──────────┘     └────────────┘
// Normal Hierarchy
// Normal
// 	Require
// 		canon
// 		only
// 		need
// 		order
// 	Option
// 		option
//
// Venn Diagram of Normal - Require "mixing" with Option
//
// Require | Both | Option
// ┌───────┬──────┬──────┐
// │       │      │      │
// │ Canon │ Need │Option│
// │ Only  │ Order│      │
// │       │      │      │
// └───────┴──────┴──────┘
//
type Normal string

type (
	Canon  []Normal
	Only   []Normal
	Need   []Normal
	Option []Normal
	Extra  []Normal
)

func Type[T ~[]Normal](norm T) string {
	switch any(norm).(type) {
	default:
		return ""
	case Canon:
		return "canon"
	case Only:
		return "only"
	case Option:
		return "option"
	case Need:
		return "need"
	case Extra:
		return "extra"
	}
}

// TODO
func normalMerge[T ~[]Normal](skip int, norms ...T) []T {
	if len(norms) < 2 || skip == len(norms) {
		return norms
	}
	norm := norms[skip]
	switch any(norm).(type) {
	default:
		fmt.Println("Warning: Default")
		return nil
	case Canon:

	case Only:

	case Option:

	case Need:

	case Extra:

	}

	// var out any    // [][]string
	// var merged any // []string
	// for i := 0; i < len(norms); i++ {

	// 	merged = norms[i]
	// 	if NormType(norms[i]) == NormType(norms[i+1]) {

	// 		v, ok := norms[i].([]string)
	// 		v2, ok2 := norms[i+1].([]string)

	// 		if !ok && !ok2 {
	// 			fmt.Println("Not okay")
	// 			return nil
	// 		}

	// 		merged = append(v, v2...)
	// 		continue
	// 	}
	// 	out = append(out, merged)
	// }

	return norms
}

// IsNormal checks if a Coze is normalized.  See notes on Normal.  Param opt may
// be nil.  If opt is considered invalid for Canon and Only and if opt is set
// for either type function returns false. If opt is not nil for Need or Order,
// no extra fields are allowed outside of what's specified by norm plus opt. If
// opt is nil, all extra fields are valid.
//
// Note that parameter norm must be typed as Canon, Only, Need, Order, or
// Option.  (TODO probably type norm as Normal.  There appears to be some Go
// issues typing this)
//
// Repeated keys between opt and norm is allowed.
func IsNormal[T ~[]Normal](pay json.RawMessage, norm ...T) bool {
	// fmt.Printf("IsNormal pay: %s, norm: %+v\n", pay, norm)
	ms := MapSlice{}
	err := json.Unmarshal(pay, &ms)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	return isNormal(ms, 0, norm...)
}

//
// if opt != nil {
// 	sort.Strings(opt)

// 	// Optional only.
// 	if norm == nil {
// 		norm = opt
// 	}
// }

// Useful:
// if opt != nil {
// 	v = append(v, opt...) // merge
// }

// case Order:
// 	i := 0
// 	value := ""
// 	keys := ms.KeysString()

// 	if len(v) > len(keys) {
// 		return false
// 	}

// 	for i, value = range v {
// 		if value != keys[i] {
// 			return false
// 		}
// 	}
// 	if opt != nil {
// 		after := keys[i+1:]
// 		sort.Strings(opt)
// 		for _, value = range after {
// 			if !contains(opt, value) {
// 				return false
// 			}
// 		}
// 	}

// When `canon`, `only`, and `option` are succeeded by another norm, look ony at
// the next x number of records and ignore all other preceding.
//
// When `canon`, `only`, and `option` are preceded by another norm, matching starts at the first record.
// the next x number of records and ignore all other preceding.
//
// Norms touching adjacent types of the same type are merged.
//
// TODO think about merging canon, only, opt, when touching
func isNormal[T ~[]Normal](ms MapSlice, skip int, norms ...T) bool {
	if skip >= len(norms) {
		return true
	}

	fmt.Printf("isNormal ms: %s, skip %d, skipNorm: %v, norm: %v, normlen: %d\n", ms, skip, norms[skip], norms, len(norms))

	norm := norms[skip]
	if norm == nil || len(norm) == 0 {
		return true // A nil or zero norm matches everything.
	}

	lastNorm := false // If the norm is the last variadic norm
	if skip+1 == len(norms) {
		lastNorm = true
	}

	// Canon, Only, Option do not allow any preceding fields when last norm.
	if lastNorm {
		fmt.Println("Last Norm")
		// On switch with generic, there is a Go "bug"
		// https://github.com/golang/go/issues/45380#issuecomment-1014950980
		switch any(norm).(type) {
		case Canon, Only:
			if len(norm) != len(ms) {
				return false
			}
		case Option:
			if len(ms) > len(norm) {
				return false
			}
		}
	}

	passedRecs := 0

	switch v := any(norm).(type) {
	case Canon:
		fmt.Printf("can\n")
		for i, n := range v {
			if n != Normal(ms[i].Key) {
				fmt.Println(ms[i].Key, n)
				return false
			}
			passedRecs++
		}
	case Only:
		fmt.Println("ony")
		keys := ms[:len(norm)].KeysString()
		slices.Sort(keys)
		slices.Sort(v)

		for i := range v {
			if v[i] != Normal(keys[i]) {
				return false
			}
			passedRecs++
		}
	// case Option:
	// 	fmt.Println("Opt")

	// 	length := min(len(norm), len(ms))
	// 	fmt.Printf("keylen %d, lenght %d\n", len(ms), length)
	// 	keys := ms[:length].KeysString()
	// 	sort.Strings(keys)
	// 	sort.Strings(v)

	// 	// If keys contains any extra key, return false.

	// 	for _, n := range keys {
	// 		if !contains(v, n) {
	// 			return false
	// 		}
	// 		passedRecs++
	// 	}
	// case Need:
	// 	fmt.Println("ned")
	// 	keys := ms.KeysString()
	// 	sort.Strings(keys)
	// 	sort.Strings(v)
	// 	matches := 0
	// 	// optMatches := 0
	// 	for _, value := range keys {
	// 		// TODO fix len norms
	// 		if matches == len(norms) || v[matches] != value {
	// 			// if opt != nil {
	// 			// 	if optMatches == len(opt) || opt[optMatches] != value {
	// 			// 		return false
	// 			// 	}
	// 			// 	optMatches++
	// 			// }
	// 			continue
	// 		} else {
	// 			matches++
	// 		}
	// 	}
	// 	// Bookends
	// 	if matches != len(v) {
	// 		return false
	// 	}
	// 	// if opt != nil && optMatches+matches != len(keys) {
	// 	// 	return false
	// 	// }

	case Extra:
		// Do nothing
	}

	if passedRecs >= len(ms) {
		return true
	}

	fmt.Printf("lastNorm: %v, skip %d, normLength: %d, passedRecs: %d\n", lastNorm, skip, len(norms), passedRecs)
	return isNormal(ms[passedRecs:], skip+1, norms...)
}

func contains(s []string, search string) bool {
	i := sort.SearchStrings(s, search)
	return i < len(s) && s[i] == search
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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
	if canon == nil { // only compactify
		var b bytes.Buffer
		err = json.Compact(&b, input)
		if err != nil {
			return nil, err
		}
		return b.Bytes(), nil
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
