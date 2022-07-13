package coze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// Normal - All normals are arrays of fields.  In Go this is implemented at a
// slice []string.
//
//  canon      (can)
//  only       (ony)
//  need       (ned)
//  order      (ord)
//  option     (opt)
//
// `canon` requires specified fields in the given order and no extra fields
// permitted.
//
// `only` specifies fields that are required to be present, does not specify any
// order, and no extra fields permitted.
//
// `need` specifies fields that are required to be present, but does not specify
// any order. Additional fields are permitted.
//
// `order` requires specified fields in the given order and additional fields
// are permitted after the order fields.
//
// `option` specifies permissable optional fields and may be used alone or used
// with `need` or `order`. All fields not in `option` and the respective `need`
// or `order` are invalid. If option is nil, all extra fields are valid.
//
// ## Normal, Require, and Option
//
// `canon`, `only`, `need`, and `order` are a `require` in that they specify
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
//
// Normal
// 	Requires:
// 		canon
// 		only
// 		need
// 		order
// 	Option:
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
type Normal []string

type Canon Normal
type Only Normal
type Need Normal
type Order Normal
type Option Normal

// IsNormal checks if a Coze is normalized.  Param opt may be nil.  If opt is
// considered invalid for Canon and Only and if opt is set for either type
// function returns false. If opt is not nil for Need or Order, no extra fields
// are allowed outside of what's specified by norm plus opt.
//
// Repeated keys between opt and norm is allowed.
//
// TODO write Normalize()
func IsNormal(pay json.RawMessage, norm any, opt Option) bool {
	var ms = MapSlice{}
	err := json.Unmarshal(pay, &ms)
	if err != nil {
		return false
	}

	if opt != nil {
		sort.Strings(opt)

		// Optional only.
		if norm == nil {
			norm = opt
		}
	}

	switch v := norm.(type) {
	case Canon:
		if len(v) != len(ms) {
			return false
		}
		for i, mi := range ms {
			if mi.Key != v[i] {
				return false
			}
		}
	case Only:
		if len(v) != len(ms) {
			return false
		}

		keys := ms.Keys()
		sort.Strings(keys)
		sort.Strings(v)
		for i := range v {
			if v[i] != keys[i] {
				return false
			}
		}

	case Need:
		keys := ms.Keys()
		sort.Strings(keys)
		sort.Strings(v)
		matches := 0
		optMatches := 0
		for _, value := range keys {
			if matches == len(v) || v[matches] != value {
				if opt != nil {
					if optMatches == len(opt) || opt[optMatches] != value {
						return false
					}
					optMatches++
				}
				continue
			} else {
				matches++
			}
		}
		// Bookends
		if matches != len(v) {
			return false
		}

		//fmt.Println(optMatches, matches, len(keys))
		if opt != nil && optMatches+matches != len(keys) {
			return false
		}
	case Order:
		i := 0
		value := ""
		keys := ms.Keys()

		if len(v) > len(keys) {
			return false
		}

		for i, value = range v {
			if value != keys[i] {
				return false
			}
		}
		if opt != nil {
			after := keys[i+1:]
			sort.Strings(opt)
			for _, value = range after {
				if !contains(opt, value) {
					return false
				}
			}
		}
	case Option:
		if opt != nil {
			v = append(v, opt...) //merge
		}
		sort.Strings(v)
		keys := ms.Keys()
		for _, value := range keys {
			if !contains(v, value) {
				return false
			}
		}
	default:
		return false
	}
	return true
}

func contains(s []string, search string) bool {
	i := sort.SearchStrings(s, search)
	return i < len(s) && s[i] == search
}

// Canon returns the current canon from raw JSON.
//
// It returns only top level fields with no recursion or promotion of embedded
// fields.
func GetCanon(raw json.RawMessage) (can []string, err error) {
	var ms = MapSlice{}
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
