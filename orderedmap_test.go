// The MIT License (MIT)
//
// Copyright (c) 2023 Cypherpunk LLC and contributors
// Copyright (c) 2017 Ian Coleman
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, Subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or Substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package coze

import (
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestOrderedMap(t *testing.T) {
	o := newOrderedMap()
	o.Set("number", 3) // number
	v, _ := o.Get("number")
	if v.(int) != 3 {
		t.Error("Set number")
	}
	o.Set("string", "x") // string
	v, _ = o.Get("string")
	if v.(string) != "x" {
		t.Error("Set string")
	}
	o.Set("strings", []string{ // string slice
		"t",
		"u",
	})
	v, _ = o.Get("strings")
	if v.([]string)[0] != "t" {
		t.Error("Set strings first index")
	}
	if v.([]string)[1] != "u" {
		t.Error("Set strings second index")
	}
	o.Set("mixed", []any{ // mixed slice
		1,
		"1",
	})
	v, _ = o.Get("mixed")
	if v.([]any)[0].(int) != 1 {
		t.Error("Set mixed int")
	}
	if v.([]any)[1].(string) != "1" {
		t.Error("Set mixed string")
	}

	// overriding existing key
	o.Set("number", 4)
	v, _ = o.Get("number")
	if v.(int) != 4 {
		t.Error("Override existing key")
	}

	// Keys method
	keys := o.Keys()
	expectedKeys := []string{
		"number",
		"string",
		"strings",
		"mixed",
	}
	for i, key := range keys {
		if key != expectedKeys[i] {
			t.Error("Keys method", key, "!=", expectedKeys[i])
		}
	}

	// Values method
	values := o.Values()
	expectedValues := []any{
		4, // 3 is overwritten
		"x",
		[]string{"t", "u"},
		[]any{1, "1"},
	}
	for i, val := range values {
		if !reflect.DeepEqual(expectedValues[i], val) {
			t.Error("Values method", expectedValues[i], "!=", val)
		}
	}
}

func TestOrderedMapDelete(t *testing.T) {
	o := newOrderedMap()
	o.Set("strings", "stringValue")
	o.Delete("strings")
	o.Delete("not a key being used")
	_, ok := o.Get("strings")
	if ok || len(o.Keys()) != 0 { // Demonstrates that delete isn't set to nil, but is removed.
		t.Error("Delete did not remove 'strings' key")
	}
}

func TestBlankMarshalJSON(t *testing.T) {
	o := newOrderedMap()
	// blank map
	b, err := json.Marshal(o)
	if err != nil {
		t.Error("Marshalling blank map to json", err)
	}
	s := string(b)
	// check json is correctly ordered
	if s != `{}` {
		t.Error("JSON Marshaling blank map value is incorrect", s)
	}
	// convert to indented json
	bi, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		t.Error("Marshalling indented json for blank map", err)
	}

	ei := `{}`
	if string(bi) != ei {
		t.Error("JSON MarshalIndent blank map value is incorrect", ei)
	}
}

func TestMarshalJSON(t *testing.T) {
	o := newOrderedMap()
	o.Set("number", 3)                    // number
	o.Set("string", "x")                  // string
	o.Set("specialstring", "\\.<>[]{}_-") // string
	o.Set("number", 4)                    // new value keeps key in old position
	// keys not sorted alphabetically
	o.Set("z", 1)
	o.Set("a", 2)
	o.Set("b", 3)
	o.Set("slice", []any{ // slice
		"1",
		1,
	})

	v := newOrderedMap()
	v.Set("e", 1)
	v.Set("a", 2)
	o.Set("orderedmap", v)
	o.Set("test\n\r\t\\\"ing", 9) // escape key

	// convert to indented json
	bi, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		t.Error("Marshalling indented json", err)
	}
	ei := `{
  "number": 4,
  "string": "x",
  "specialstring": "\\.\u003c\u003e[]{}_-",
  "z": 1,
  "a": 2,
  "b": 3,
  "slice": [
    "1",
    1
  ],
  "orderedmap": {
    "e": 1,
    "a": 2
  },
  "test\n\r\t\\\"ing": 9
}`
	if ei != string(bi) {
		t.Error("JSON MarshalIndent value is incorrect", ei)
	}
}

func TestMarshalJSONNoEscapeHTML(t *testing.T) {
	o := newOrderedMap()
	// string special characters
	o.Set("specialstring", "\\.<>[]{}_-")
	// convert to json
	b, err := o.MarshalJSON()
	if err != nil {
		t.Error("Marshalling json", err)
	}
	s := strings.Replace(string(b), "\n", "", -1)
	// check json is correctly ordered
	if s != `{"specialstring":"\\.<>[]{}_-"}` {
		t.Error("JSON Marshal value is incorrect", s)
	}
}

func TestMarshalJSONNoEscapeHTMLRecursive(t *testing.T) {
	src := `{"x":"<>","y":[{"z":["<>"]}]}`
	o := newOrderedMap()
	err := json.Unmarshal([]byte(src), &o)
	if err != nil {
		t.Error("JSON Unmarshal error with special chars", err)
	}
	b, err := o.MarshalJSON()
	if err != nil {
		t.Error("Marshalling json", err)
	}
	s := strings.Replace(string(b), "\n", "", -1)
	if s != src {
		t.Error("JSON Marshal value is incorrect", s)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	s := `{
  "number": 4,
  "string": "x",
  "z": 1,
  "a": "should not break with unclosed { character in value",
  "b": 3,
  "slice": [
    "1",
    1
  ],
  "orderedmap": {
    "e": 1,
    "a { nested key with brace": "with a }}}} }} {{{ brace value",
	"after": {
		"link": "test {{{ with even deeper nested braces }"
	}
  },
  "test\"ing": 9,
  "after": 1,
  "multitype_array": [
    "test",
	1,
	{ "map": "obj", "it" : 5, ":colon in key": "colon: in value" },
	[{"inner": "map"}]
  ],
  "should not break with { character in key": 1
}`
	o := newOrderedMap()
	err := json.Unmarshal([]byte(s), &o)
	if err != nil {
		t.Error("JSON Unmarshal error", err)
	}
	// Check the root keys
	expectedKeys := []string{
		"number",
		"string",
		"z",
		"a",
		"b",
		"slice",
		"orderedmap",
		"test\"ing",
		"after",
		"multitype_array",
		"should not break with { character in key",
	}
	k := o.Keys()
	for i := range k {
		if k[i] != expectedKeys[i] {
			t.Error("Unmarshal root key order", i, k[i], "!=", expectedKeys[i])
		}
	}
	// Check nested maps are converted to orderedmaps
	// nested 1 level deep
	expectedKeys = []string{
		"e",
		"a { nested key with brace",
		"after",
	}
	vi, ok := o.Get("orderedmap")
	if !ok {
		t.Error("Missing key for nested map 1 deep")
	}
	v := vi.(orderedMap) // panics if not correct type
	k = v.Keys()
	for i := range k {
		if k[i] != expectedKeys[i] {
			t.Error("Key order for nested map 1 deep ", i, k[i], "!=", expectedKeys[i])
		}
	}
	// nested 2 levels deep
	expectedKeys = []string{
		"link",
	}
	vi, ok = v.Get("after")
	if !ok {
		t.Error("Missing key for nested map 2 deep")
	}
	v = vi.(orderedMap) // panics if not correct type
	k = v.Keys()
	for i := range k {
		if k[i] != expectedKeys[i] {
			t.Error("Key order for nested map 2 deep", i, k[i], "!=", expectedKeys[i])
		}
	}
	// multitype array
	expectedKeys = []string{
		"map",
		"it",
		":colon in key",
	}
	vislice, ok := o.Get("multitype_array")
	if !ok {
		t.Error("Missing key for multitype array")
	}
	vslice := vislice.([]any) // panics if not correct type
	vmap := vslice[2].(orderedMap)
	k = vmap.Keys()
	for i := range k {
		if k[i] != expectedKeys[i] {
			t.Error("Key order for nested map 2 deep", i, k[i], "!=", expectedKeys[i])
		}
	}
	// nested map 3 deep
	vislice, _ = o.Get("multitype_array")
	vslice = vislice.([]any) // panics if not correct type
	expectedKeys = []string{"inner"}
	vinnerslice := vslice[3].([]any)
	vinnermap := vinnerslice[0].(orderedMap)
	k = vinnermap.Keys()
	for i := range k {
		if k[i] != expectedKeys[i] {
			t.Error("Key order for nested map 3 deep", i, k[i], "!=", expectedKeys[i])
		}
	}
}

// TestOrderedMapUnmarshalJSONDuplicate tests that orderedMap errors on
// duplicate JSON fields.
func TestOrderedMapUnmarshalJSONDuplicate(t *testing.T) {
	s := `{
		"a": [{}, []],
		"b": {"x":[1]},
		"c": "x",
		"d": {"x":1},
		"b": [{"x":[]}]
	}`
	o := newOrderedMap()
	err := json.Unmarshal([]byte(s), &o)
	if err == nil {
		t.Errorf("orderedMap unmarshal did not error on duplicate key")
	}
}

func TestUnmarshalJSONSpecialChars(t *testing.T) {
	s := `{ " \u0041\n\r\t\\\\\\\\\\\\ "  : { "\\\\\\" : "\\\\\"\\" }, "\\":  " \\\\ test ", "\n": "\r" }`
	o := newOrderedMap()
	err := json.Unmarshal([]byte(s), &o)
	if err != nil {
		t.Error("JSON Unmarshal error with special chars", err)
	}
	expectedKeys := []string{
		" \u0041\n\r\t\\\\\\\\\\\\ ",
		"\\",
		"\n",
	}
	keys := o.Keys()
	if len(keys) != len(expectedKeys) {
		t.Error("Unmarshal key count", len(keys), "!=", len(expectedKeys))
	}
	for i, key := range keys {
		if key != expectedKeys[i] {
			t.Errorf("Unmarshal root key order: %d, %q != %q", i, key, expectedKeys[i])
		}
	}
}

func TestUnmarshalJSONArrayOfMaps(t *testing.T) {
	s := `
{
  "name": "test",
  "percent": 6,
  "breakdown": [
    {
      "name": "a",
      "percent": 0.9
    },
    {
      "name": "b",
      "percent": 0.9
    },
    {
      "name": "d",
      "percent": 0.4
    },
    {
      "name": "e",
      "percent": 2.7
    }
  ]
}
`
	o := newOrderedMap()
	err := json.Unmarshal([]byte(s), &o)
	if err != nil {
		t.Error("JSON Unmarshal error", err)
	}
	// Check the root keys
	expectedKeys := []string{
		"name",
		"percent",
		"breakdown",
	}
	k := o.Keys()
	for i := range k {
		if k[i] != expectedKeys[i] {
			t.Error("Unmarshal root key order", i, k[i], "!=", expectedKeys[i])
		}
	}
	// Check nested maps are converted to orderedMaps
	// nested 1 level deep
	expectedKeys = []string{
		"name",
		"percent",
	}
	vi, ok := o.Get("breakdown")
	if !ok {
		t.Error("Missing key for nested map 1 deep")
	}
	vs := vi.([]any)
	for _, vInterface := range vs {
		v := vInterface.(orderedMap)
		k = v.Keys()
		for i := range k {
			if k[i] != expectedKeys[i] {
				t.Error("Key order for nested map 1 deep ", i, k[i], "!=", expectedKeys[i])
			}
		}
	}
}

func TestUnmarshalJSONStruct(t *testing.T) {
	var v struct {
		Data *orderedMap `json:"data"`
	}

	err := json.Unmarshal([]byte(`{ "data": { "x": 1 } }`), &v)
	if err != nil {
		t.Fatalf("JSON unmarshal error: %v", err)
	}

	x, ok := v.Data.Get("x")
	if !ok {
		t.Errorf("missing expected key")
	} else if x != float64(1) {
		t.Errorf("unexpected value: %#v", x)
	}
}

func TestOrderedMap_SortKeys(t *testing.T) {
	s := `
{
  "b": 2,
  "a": 1,
  "c": 3
}
`
	o := newOrderedMap()
	json.Unmarshal([]byte(s), &o)

	o.SortKeys(sort.Strings)

	// Check the root keys
	expectedKeys := []string{
		"a",
		"b",
		"c",
	}
	k := o.Keys()
	for i := range k {
		if k[i] != expectedKeys[i] {
			t.Error("SortKeys root key order", i, k[i], "!=", expectedKeys[i])
		}
	}
}

func TestOrderedMap_Sort(t *testing.T) {
	s := `
{
  "b": 2,
  "a": 1,
  "c": 3
}
`
	o := newOrderedMap()
	json.Unmarshal([]byte(s), &o)
	o.Sort(func(a *pair, b *pair) bool {
		return a.value.(float64) > b.value.(float64)
	})

	// Check the root keys
	expectedKeys := []string{
		"c",
		"b",
		"a",
	}
	k := o.Keys()
	for i := range k {
		if k[i] != expectedKeys[i] {
			t.Error("Sort root key order", i, k[i], "!=", expectedKeys[i])
		}
	}
}

// https://github.com/iancoleman/orderedmap/issues/11
func TestOrderedMap_empty_array(t *testing.T) {
	srcStr := `{"x":[]}`
	src := []byte(srcStr)
	om := newOrderedMap()
	json.Unmarshal(src, om)
	bs, _ := json.Marshal(om)
	marshalledStr := string(bs)
	if marshalledStr != srcStr {
		t.Error("Empty array does not serialize to json correctly")
		t.Error("Expect", srcStr)
		t.Error("Got", marshalledStr)
	}
}

// Inspired by
// https://github.com/iancoleman/orderedmap/issues/11
// but using empty maps instead of empty slices
func TestOrderedMap_empty_map(t *testing.T) {
	srcStr := `{"x":{}}`
	src := []byte(srcStr)
	om := newOrderedMap()
	json.Unmarshal(src, om)
	bs, _ := json.Marshal(om)
	marshalledStr := string(bs)
	if marshalledStr != srcStr {
		t.Error("Empty map does not serialize to json correctly")
		t.Error("Expect", srcStr)
		t.Error("Got", marshalledStr)
	}
}
