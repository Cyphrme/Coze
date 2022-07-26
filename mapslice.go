package coze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// MapSlice of map items.
//
// MapSlice and MapItem was originally taken from:
// https://github.com/golang/go/issues/27179#issuecomment-587528269,
// https://github.com/ake-persson/mapslice-json. Thanks ake-persson!  This has
// been modified from ake-persson/mapslice-json to support only things that JSON
// or Coze needs, chiefly, Key is now type string instead of type any.
type MapSlice []MapItem

func (ms MapSlice) Len() int           { return len(ms) }
func (ms MapSlice) Less(i, j int) bool { return ms[i].index < ms[j].index }
func (ms MapSlice) Swap(i, j int)      { ms[i], ms[j] = ms[j], ms[i] }

var indexCounter uint64

func nextIndex() uint64 {
	indexCounter++
	return indexCounter
}

// Keys returns a MapSlice's Keys in a slice.
func (ms MapSlice) Keys() []any {
	s := make([]any, len(ms))
	for i, k := range ms {
		s[i] = k.Key
	}
	return s
}

// KeysString returns MapSlice's Keys as []string.
func (ms MapSlice) KeysString() []string {
	s := make([]string, len(ms))
	for i, k := range ms {
		s[i] = fmt.Sprintf("%s", k.Key)
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
	buf.Write([]byte{'{'})
	for i, mi := range ms {
		b, err := json.Marshal(&mi.Value)
		if err != nil {
			return nil, err
		}
		buf.WriteString(fmt.Sprintf("%q:", fmt.Sprintf("%v", mi.Key)))
		buf.Write(b)
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
