package coze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// `encoding/json` has no way to preserve the order of map keys.
// See https://github.com/golang/go/issues/27179

// MapSlice is not safe for concurrent use.  See https://github.com/Cyphrme/Coze/issues/10#issuecomment-1498196356

// MapSlice is for setting of a given order for maps map items.  Go maps have an
// arbitrary order that cannot be (easily) set.
//
// This implementation of MapSlice is designed for JSON/Coze. Key is now type
// string (instead of type any), and has extra methods, "Keys()",
// "KeyStrings()", "Values()".
//
// We may publish this as a standalone package.
//
// go-yaml, a "dead"/done project, had the same problem and was solved 8 years
// ago:
// https://github.com/go-yaml/yaml/blob/7649d4548cb53a614db133b2a8ac1f31859dda8c/yaml.go#L20
//
// MapSlice and MapItem was originally inspired from (3 years ago):
// https://github.com/golang/go/issues/27179#issuecomment-587528269,
// https://github.com/ake-persson/mapslice-json, but now appears to be dead.
// https://github.com/ake-persson/mapslice-json/pull/1.
// https://github.com/ake-persson/mapslice-json/pull/3
type MapSlice []MapItem

// Implements "sort.Interface"
func (ms MapSlice) Len() int           { return len(ms) }
func (ms MapSlice) Less(i, j int) bool { return ms[i].index < ms[j].index }
func (ms MapSlice) Swap(i, j int)      { ms[i], ms[j] = ms[j], ms[i] }

var indexCounter uint64

func nextIndex() uint64 {
	indexCounter++
	return indexCounter
}

// Keys returns MapSlice's Keys as []string.
func (ms MapSlice) Keys() []string {
	s := make([]string, len(ms))
	for i, k := range ms {
		s[i] = k.Key
	}
	return s
}

// Values returns a MapSlice's values in a slice.
func (ms MapSlice) Values() []any {
	s := make([]any, len(ms))
	for i, k := range ms {
		s[i] = k.Value
	}
	return s
}

// MarshalJSON for map slice.
func (ms MapSlice) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.Write([]byte{'{'}) // Write err is always nil.
	for i, mi := range ms {
		k, err := Marshal(&mi.Key)
		if err != nil {
			return nil, err
		}
		buf.Write(k)
		buf.Write([]byte{':'})

		v, err := Marshal(&mi.Value)
		if err != nil {
			return nil, err
		}

		buf.Write(v)

		if i < len(ms)-1 {
			buf.Write([]byte{','})
		}
	}
	buf.Write([]byte{'}'})
	return buf.Bytes(), nil
}

// UnmarshalJSON for map slice.
func (ms *MapSlice) UnmarshalJSON(b []byte) error {
	m := map[string]MapItem{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	for k, v := range m {
		*ms = append(*ms, MapItem{Key: k, Value: v.Value, index: v.index})
	}
	sort.Sort(*ms)
	return nil
}

// MapItem representation of one map item.
type MapItem struct {
	Key   string
	Value any
	index uint64
}

// MapItem as a string.
func (mi MapItem) String() string {
	return fmt.Sprintf("{%v %v}", mi.Key, mi.Value)
}

// UnmarshalJSON for map item.
func (mi *MapItem) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	mi.Value = v
	mi.index = nextIndex()
	return nil
}
