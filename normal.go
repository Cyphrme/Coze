package coze

import (
	"bytes"
	"encoding/json"
	"fmt"

	"golang.org/x/exp/slices"
)

// Normal - A normal is an arrays of fields specifying the normalization of
// payload and Normals may be chained to represent various combinations of
// normalization.  Normals are implemented in Go as []string.  There are five
// types of normals.
//
//  canon       (can)
//  only        (ony)
//  option      (opt)
//  need        (ned)
//  extra       (ext)
//  (nil)
//
// An nil normal is valid and matches all payloads.
//
// Canon requires specified fields in the given order and no extra fields
// permitted.
//
// Only specifies fields that are required to be present, does not specify any
// order, and no extra fields permitted.
//
// Option permits the presence of the given fields and excludes the presence of
// extra fields.  Option does not specify order, but order may be given by
// chaining options together.
//
// Need specifies fields that are required to be present, but does not specify
// any order and does not exclude the presence of other fields.
//
// Extra specifies extra fields are permitted in it's location in the normal
// chain.  An extra containing fields has no addition meaning over an empty
// Extra.
//
//
// Normal Chaining
//
// Normals may be chained:
//
// Canon, Only, Option, and Extra have meaning when chained
// Need's position in a has no chain meaning.  Needs position in a chain is irrelevant.
//
// Repeated keys between (Canon or Option or Only) and (Need) is allowed.
//
// If opt is considered invalid for Canon and Only and if opt is set
// for either type function returns false. If opt is not nil for Need or Order,
// no extra fields are allowed outside of what's specified by norm plus opt. If
// opt is nil, all extra fields are valid.
//
// Note that parameter norm must be typed as Canon, Only, Need, Order, or
// Option.  (TODO probably type norm as Normal.  There appears to be some Go
// issues typing this)
//
// Canon, Only, Option and Extra have meaning when chained.  Need has no meaning
// when chained, and may appear anywhere in the chain without changing the
// meaning of the Normal chain.
//
// Any names in Extra are ignored since they are equivalent to an empty Extra.
//
//
//
// Interesting Combinations:
// - A an empty Canon or Only ("[]") matches only an a empty (i.e. `{}`) payload.
// - An Empty Need or Option does nothing. // TODO?
// - If need can appear before or after another normal, call IsNormal twice: a IsNormal(r, Need{a}), IsNormal(r, Canon{"b","c"}})
//

//
//            ┌────────────────┐
//            │     Normal     │
//            └───────┬────────┘
//            ┌───────┴────────┐
//      ┌─────┴─────┐    ┌─────┴────────┐
//      │ Exclusive │    │ Permissive   │
//      └───────────┘    └──────────────┘
// Normal Hierarchy
// Normal
// 	Exclusive
// 		canon
// 		only
// 		option
// 	Permissive
// 		need
//    extra
//
type Normal string

type (
	Canon  []Normal
	Only   []Normal
	Need   []Normal
	Option []Normal
	Extra  []Normal
)

type Normaler interface {
	Len() int
	Normal() []Normal
}

func (n Canon) Len() int {
	return len(n)
}
func (n Only) Len() int {
	return len(n)
}
func (n Need) Len() int {
	return len(n)
}
func (n Option) Len() int {
	return len(n)
}
func (n Extra) Len() int {
	return len(n)
}

func (n Canon) Normal() []Normal {
	return n
}
func (n Only) Normal() []Normal {
	return n
}
func (n Need) Normal() []Normal {
	return n
}
func (n Option) Normal() []Normal {
	return n
}
func (n Extra) Normal() []Normal {
	return n
}

// func (c *Canon) Append(norms []Normal) {
// 	*c = append(*c, norms...)
// }
// func (o *Only) Append(norms []Normal) {
// 	*o = append(*o, norms...)
// }
// func (n *Need) Append(norms []Normal) {
// 	*n = append(*n, norms...)
// }
// func (o *Option) Append(norms []Normal) {
// 	*o = append(*o, norms...)
// }
// func (e *Extra) Append(norms []Normal) {
// 	*e = append(*e, norms...)
// }

func Append(n, m []Normal) []Normal {
	return append(n, m...)
}

// Type returns the type for a given Normaler including a case for []Normal.
func Type(norm Normaler) string {
	switch norm.(type) {
	default:
		return "invalid normal"
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

// Merge merges the given normals.
func Merge[T ~[]Normal](norms ...T) any {
	n := norms[0]
	for i := 1; i < len(norms); i++ {
		n = Append(n, norms[i])
	}
	return n
}

// // Union merges adjacent normals of the same type.
// func Union(skip int, norms ...Normaler) any {
// 	fmt.Printf("Union norms: %v\n", norms)
// 	if len(norms) < 2 || skip == len(norms) {
// 		return norms
// 	}

// 	normType := Type(norms[skip])
// 	mergeTo := 0

// 	// Get repeated normals
// 	for i := skip + 1; i < len(norms); i++ {
// 		fmt.Printf("Type: %s\n", Type(norms[i]))
// 		if normType == Type(norms[i]) {
// 			mergeTo = i
// 			continue
// 		}
// 		break
// 	}

// 	var n []Normal
// 	//Finally, merge repeated normals
// 	if mergeTo > 0 {
// 		toMerge := norms[skip+1 : mergeTo+1]
// 		for _, v := range toMerge {
// 			for _, y := range v {
// 				n = append(n, Normal(y))
// 			}
// 		}

// 		fmt.Printf("mergeTo: %d, ToMerge: %v, n: %v\n", mergeTo, toMerge, n)
// 		norms[skip] = append(norms[skip], n...)

// 		fmt.Printf("Before Delete: %v, skip %d, mergeto %d \n", norms, skip, mergeTo)
// 		norms = slices.Delete(norms, skip+1, mergeTo+1)
// 		fmt.Printf("After Delete: %v\n", norms)
// 	}
// 	return norms
// }

// IsNormal checks if a Coze is normalized.  See notes on Normal.  Parameters may
// be nil.
func IsNormal(pay json.RawMessage, norm ...Normaler) bool {
	// fmt.Printf("IsNormal pay: %s, norm: %+v\n", pay, norm)

	ms := MapSlice{}
	err := json.Unmarshal(pay, &ms)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	return isNormal(ms, 0, 0, false, norm...)
}

// When `canon`, `only`, and `option` are succeeded by another norm, look ony at
// the next x number of records and ignore all other preceding.
//
// When `canon`, `only`, and `option` are preceded by another norm, matching starts at the first record.
// the next x number of records and ignore all other preceding.
//
// Norms touching adjacent types of the same type are merged.

// isNormal checks if datastructure conforms to the given normal chain. See docs
// on IsNormal.
//
// Params:
//  r          (Records) The fields being checked if conforming to normal chain.
//  rSkip      Record Pointer - First field that has not yet beet checked.
//  nSkip      Normal Pointer - First Normal that has not been processed.
//  extraFlag  When Extra has been evoked. Is disabled when Norm is not an
//    Extra, and enabled when Norm is an Extra.  When chained, allows any
//    fields until next Normal.
//  norms - The Normal chain, the full slice of normals.
//
func isNormal(r MapSlice, rSkip int, nSkip int, extraFlag bool, norms ...Normaler) bool {
	if nSkip >= len(norms) {
		return true
	}
	norm := norms[nSkip]
	switch norm.(type) {
	case nil: // Nil normal matches everything.
		return isNormal(r, rSkip, nSkip+1, false, norms...)
	case Canon, Only:
		// []Norm cannot be is greater than the number of remaining records, false.
		// (Prevent out of bounds later on)
		if norm.Len() > len(r)-rSkip {
			return false
		}
	case Option, Need:
		// An empty Option or Need progresses norm pointer and nothing else.
		if norm.Len() == 0 {
			return isNormal(r, rSkip, nSkip+1, false, norms...)
		}
	case Extra: // Extra flag.
		return isNormal(r, rSkip, nSkip+1, true, norms...)
	}

	//fmt.Printf("isNormal{r: %s, rSkip %d, nSkip %d, GoType/Type: %T/%s norm: %v, norm.Len(): %d, norms len: %d, norms: %v}\n", r, rSkip, nSkip, norm, Type(norm), norm, norm.Len(), len(norms), norms)

	// Progress pointer to first match for Canon/Only/Option
	if extraFlag {
		switch v := norm.(type) {
		case Canon:
			for i := rSkip; i < len(r); i++ {
				if v[0] == Normal(r[i].Key) {
					rSkip = i
					break
				}
			}
		case Option, Only, Need:
			keys := r[rSkip:].KeysString()
			n := v.Normal()
			for i, key := range keys {
				if !slices.Contains(n, Normal(key)) {
					continue
				}
				rSkip = rSkip + i
				break
			}
		}
	}

	switch norm.(type) {
	case Canon, Only, Option: // last norm does not allow extra records
		if nSkip+1 == len(norms) && norm.Len() < len(r)-rSkip {
			return false
		}
	}

	passedRecs := 0
	switch v := norm.(type) {
	default:
		return false
	case Canon:
		for i, n := range v {
			if n != Normal(r[rSkip+i].Key) {
				return false
			}
			passedRecs++
		}
	case Only:
		keys := r[rSkip : v.Len()+rSkip].KeysString()
		slices.Sort(keys)
		slices.Sort(v)
		for i := range v {
			if v[i] != Normal(keys[i]) {
				return false
			}
			passedRecs++
		}
	case Option:
		if nSkip+1 == len(norms) { // last norm
			min := min(norm.Len(), len(r))
			keys := r[rSkip : rSkip+min-1].KeysString()
			slices.Sort(keys)
			slices.Sort(v)
			// If keys contains any extra key and last norm, return false.
			for _, n := range keys {
				if !slices.Contains(v, Normal(n)) { // TODO sort and contains should be used together
					return false
				}
				passedRecs++
			}
		} else {
			// Progress record pointer to position of first non-match.
			keys := r[rSkip:].KeysString()
			for i, n := range keys {
				if !slices.Contains(v, Normal(n)) {
					passedRecs = i
					break
				}
			}
		}
	case Need:
		i := 0
		key := ""
		keys := r[rSkip:].KeysString()
		for i, key = range keys {
			if passedRecs == v.Len() {
				break
			}
			if !slices.Contains(v, Normal(key)) {
				continue
			}
			passedRecs++
		}

		if passedRecs != v.Len() {
			return false
		}
		// Progress record pointer up the last match of Need, and turn on extraFlag to
		// progress record pointer to first match of next Normal.
		isNormal(r, rSkip+i, nSkip+1, true, norms...)
	}

	//fmt.Printf("rSkip %d, nSkip %d, passedRecs: %d\n", rSkip+passedRecs, nSkip+1, passedRecs)
	return isNormal(r, rSkip+passedRecs, nSkip+1, false, norms...)
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
