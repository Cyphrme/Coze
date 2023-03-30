// # Normal
//
// A normal is an arrays of fields specifying the normalization of a
// payload. Normals may be chained to represent various combinations of
// normalization.  Normals are implemented in Go as []string.  There are five
// types of normals plus a nil normal.
//
//	canon       (can)
//	only        (ony)
//	option      (opt)
//	need        (ned)
//	extra       (ext)
//	(nil)
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
// Repeated field names are allowed among normals and normal chains, but Coze
// itself prohibits duplicates.
//
// # Normal Chaining
//
// Normals may be chained.  A chained normal moves a record pointer up.
//
//   - A Need in a chain is equivalent to a [Need, Extra].
//   - Options in order may be given by chaining options together.
//   - An Extra containing fields has no addition meaning over an empty
//     Extra.
//
// Notable Combinations:
// - A an empty Canon or Only ("[]") matches only an a empty (i.e. `{}`) payload.
// - An empty Need or Option does nothing.
// - If need can appear before or after another normal, call IsNormal twice: a IsNormal(r, Need{a}), IsNormal(r, Canon{"b","c"}})
//
// Normals are in two groups, exclusive and permissive.  Exclusive allows only
// listed fields.  Permissive allows fields other than listed.
//
//	      ┌────────────────┐
//	      │     Normal     │
//	      └───────┬────────┘
//	      ┌───────┴────────┐
//	┌─────┴─────┐    ┌─────┴──────┐
//	│ Exclusive │    │ Permissive │
//	└───────────┘    └────────────┘
//
// Grouping
//
//	-Exclusive
//	  -canon
//	  -only
//	  -option
//	-Permissive
//	  -need
//	  -extra
package normal

import (
	"encoding/json"

	"github.com/cyphrme/coze"
	"golang.org/x/exp/slices"
)

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

// Standard is the standard coze.pay fields. Custom fields may be appended after
// standard. e.g. `normID := append(standard, "id")`.
var Standard = []Normal{"alg", "iat", "tmb", "typ"}

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

// IsNormal checks if a Coze is normalized according to the given normal chain.
// Normals are interpreted as a chain that progress a record pointer based on
// normal rules.  See notes on Normal.  Parameters may be nil.
func IsNormal(pay json.RawMessage, norm ...Normaler) (bool, error) {
	ms := coze.MapSlice{}
	err := json.Unmarshal(pay, &ms)
	if err != nil {
		return false, err
	}
	return isNormal(ms, 0, 0, false, norm...), nil
}

// isNormal checks if datastructure conforms to the given normal chain. See docs
// on IsNormal.
//
// TODO consider pointers for r and norms.
//
// Params:
//
//	r         - (Records) The fields being checked if conforming to normal chain.
//	rSkip     - Record pointer - First field that has not yet been checked.
//	nSkip     - Normal pointer - First Normal that has not been processed.
//	extraFlag - Is set to false when Norm is not an Extra and set to true when
//	            Norm is an Extra.  When set to true, it moves the record pointer
//	            to the first field matching the following Normal.
//	norms -     The Normal chain, the full slice of normals.
func isNormal(r coze.MapSlice, rSkip int, nSkip int, extraFlag bool, norms ...Normaler) bool {
	if nSkip >= len(norms) {
		return true
	}
	norm := norms[nSkip]

	if extraFlag { // Progress record pointer to first match.
		keys := r[rSkip:].Keys()
		var i int
		for i = 0; i < len(keys); i++ {
			if slices.Contains(norm.Normal(), Normal(keys[i])) {
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

	// fmt.Printf("isNormal{r: %s, rSkip %d, nSkip %d, GoType/Type: %T/%s norm: %v, norm.Len(): %d, norms len: %d, norms: %v}\n", r, rSkip, nSkip, norm, Type(norm), norm, norm.Len(), len(norms), norms)

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
		keys := r[rSkip : v.Len()+rSkip].Keys()
		slices.Sort(keys)
		slices.Sort(v)
		for i := range v {
			if v[i] != Normal(keys[i]) {
				return false
			}
			passedRecs++
		}
	case Option:
		keys := r[rSkip:].Keys()
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
		keys := r[rSkip:].Keys()
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

	// fmt.Printf("rSkip %d, nSkip %d, passedRecs: %d\n", rSkip+passedRecs, nSkip+1, passedRecs)
	return isNormal(r, rSkip+passedRecs, nSkip+1, false, norms...)
}

// IsNormalUnchained is a helper to run a slice of normals individually and not
// as a chain.
// Passing an 'option' will return false unless the given pay only has one
// field, and it is the 'option'.
func IsNormalUnchained(pay json.RawMessage, norm ...Normaler) (bool, error) {
	for _, n := range norm {
		v, err := IsNormal(pay, n)
		if err != nil || !v {
			return false, err
		}
	}
	return true, nil
}

// IsNormalNeedOption is a helper for a special case.
//
// If a need's fields and an option's fields can be intermixed, the need is
// checked first and matching fields subtracted from records.  Then the option
// is called with the subset.
//
// Another (alternative) method that's not implemented in this function:
// IsNormal may be called twice.  Once with the need(s), and a second time with
// the option(s) concatenated with the need(s).  This is logically equivalent to
// subtracting the need.
func IsNormalNeedOption(pay json.RawMessage, need Need, option Option) (bool, error) {
	ms := coze.MapSlice{}
	err := json.Unmarshal(pay, &ms)
	if err != nil {
		return false, err
	}

	if !isNormal(ms, 0, 0, false, need) {
		return false, nil
	}

	// TODO add function "delete" to map slice.
	ms2 := coze.MapSlice{}
	for _, mi := range ms {
		if !slices.Contains(need, Normal(mi.Key)) {
			ms2 = append(ms2, coze.MapItem{Key: mi.Key, Value: mi.Value})
		}
	}
	return isNormal(ms2, 0, 0, false, option), nil
}
