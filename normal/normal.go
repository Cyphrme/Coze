package coze

import (
	"encoding/json"
	"fmt"

	"github.com/cyphrme/coze"
	"golang.org/x/exp/slices"
)

// Normal - A normal is an arrays of fields specifying the normalization of a
// payload. Normals may be chained to represent various combinations of
// normalization.  Normals are implemented in Go as []string.  There are five
// types of normals plus a nil normal.
//
//  canon       (can)
//  only        (ony)
//  option      (opt)
//  need        (ned)
//  extra       (ext)
//  (nil)
//
// Canon requires specified fields in the given order and prohibits extra fields.
//
// Only requires specified fields in any order and prohibits extra fields.
//
// Option permits specified fields in any order and prohibits extra fields.
//
// Need requires specified fields in any order and permits extra fields.
//
// Extra denotes extra fields are permitted in a location in the normal
// chain.
//
// A nil normal is valid and matches all payloads.
//
// Normal Chaining
//
// Normals may be chained:
//
// - A Need in a chain is equivalent to a [Need, Extra].
// - Option order may be given by chaining options together.
// - An Extra containing fields has no addition meaning over an empty
// Extra.
//
//
// Repeated keys are allowed among normals, but Coze itself prohibits duplicate
// object keys.
//
//
// Notable Combinations:
// - A an empty Canon or Only ("[]") matches only an a empty (i.e. `{}`) payload.
// - An empty Need or Option does nothing.
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

// TODO make append generic
func Append(n, m []Normal) []Normal {
	return append(n, m...)
}

// Merge merges the given normals.
func Merge[T ~[]Normal](norms ...T) any {
	n := norms[0]
	for i := 1; i < len(norms); i++ {
		n = Append(n, norms[i])
	}
	return n
}

// IsNormal checks if a Coze is normalized.  See notes on Normal.  Parameters
// may be nil.
func IsNormal(pay json.RawMessage, norm ...Normaler) bool {
	// fmt.Printf("IsNormal pay: %s, norm: %+v\n", pay, norm)

	ms := coze.MapSlice{}
	err := json.Unmarshal(pay, &ms)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	return isNormal(ms, 0, 0, false, norm...)
}

// isNormal checks if datastructure conforms to the given normal chain. See docs
// on IsNormal.
//
// TODO consider pointers for r and norms.
//
// Params:
//  r          (Records) The fields being checked if conforming to normal chain.
//  rSkip      Record pointer - First field that has not yet been checked.
//  nSkip      Normal pointer - First Normal that has not been processed.
//  extraFlag  When Extra has been evoked. Is disabled when Norm is not an
//    Extra, and enabled when Norm is an Extra.  When chained, allows any
//    fields until first record of next Normal.
//  norms - The Normal chain, the full slice of normals.
//
func isNormal(r coze.MapSlice, rSkip int, nSkip int, extraFlag bool, norms ...Normaler) bool {
	if nSkip >= len(norms) {
		return true
	}
	norm := norms[nSkip]

	if extraFlag { // Progress record pointer to first match.
		n := norm.Normal()
		keys := r[rSkip:].KeysString()
		var i int
		for i = 0; i < len(keys); i++ {
			if slices.Contains(n, Normal(keys[i])) {
				rSkip = i + rSkip
				break
			}
		}
		// TODO this might be written better:
		if i+1 >= len(r) { // If option is missing after an extra, return true.
			if Type(norm) == "option" {
				return true
			}
		}
	}

	switch norm.(type) {
	case nil, Need:
		// Nil or an empty Need progresses norm pointer and nothing else.
		if norm == nil || norm.Len() == 0 {
			return isNormal(r, rSkip, nSkip+1, false, norms...)
		}
	case Extra: // Extra flag.  In this function, from here out, Extra is ignored.
		return isNormal(r, rSkip, nSkip+1, true, norms...)
	case Canon, Only, Option:
		// Last norm does not allow extra records
		if nSkip+1 == len(norms) && norm.Len() < len(r)-rSkip {
			return false
		}
	}

	//fmt.Printf("isNormal{r: %s, rSkip %d, nSkip %d, GoType/Type: %T/%s norm: %v, norm.Len(): %d, norms len: %d, norms: %v}\n", r, rSkip, nSkip, norm, Type(norm), norm, norm.Len(), len(norms), norms)

	passedRecs := 0
	switch v := norm.(type) {
	case Canon:
		if norm.Len() > len(r)-rSkip {
			return false
		}
		for i, n := range v {
			if n != Normal(r[rSkip+i].Key) {
				return false
			}
			passedRecs++
		}
	case Only:
		if norm.Len() > len(r)-rSkip {
			return false
		}
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
		keys := r[rSkip:].KeysString()
		for i, n := range keys {
			if !slices.Contains(v, Normal(n)) {
				if nSkip+1 == len(norms) { // last norm
					// Extras are not allowed after Option.
					return false
				}
				// Progress record pointer to position of first non-match for passing to
				// next norm.
				passedRecs = i
				break
			}
			passedRecs++
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
		// Progress record pointer up the last match of Need, and turn on extraFlag
		// to progress record pointer to first match of next Normal.
		return isNormal(r, rSkip+i, nSkip+1, true, norms...)
	}

	//fmt.Printf("rSkip %d, nSkip %d, passedRecs: %d\n", rSkip+passedRecs, nSkip+1, passedRecs)
	return isNormal(r, rSkip+passedRecs, nSkip+1, false, norms...)
}