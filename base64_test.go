package coze

import (
	"encoding/json"
	"fmt"
)

// B64 of nil is "" while B64 of 0 is "AA".
func ExampleB64_zero_nil() {
	var b []byte
	b = nil
	fmt.Printf("B64 string nil: `%s`\n", B64(b))

	b = []byte{0}
	fmt.Printf("B64 string zero: `%s`\n", B64(b))

	// Output:
	// B64 string nil: ``
	// B64 string zero: `AA`
}

func ExampleB64_marshalJSON() {
	h := B64([]byte{0, 255})
	b, err := h.MarshalJSON()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
	// Output:
	// "AP8"
}

func ExampleB64_unmarshalJSON() {
	type Foo struct {
		Bar B64
	}

	f := new(Foo)
	err := json.Unmarshal([]byte(`{"Bar":"AP8"}`), f)
	if err != nil {
		panic(err)
	}

	b, err := Marshal(f)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n%#v\n", b, B64(b))
	// Output:
	// {"Bar":"AP8"}
	// eyJCYXIiOiJBUDgifQ
}

func ExampleDecode() {
	b, err := Decode(GoldenTmb)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
	// Output: cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk
}

func ExampleMustDecode() {
	fmt.Println(MustDecode(GoldenTmb))
	// Output: cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk
}

// Demonstrates that Coze Go will error on non-canonical base 64 encoding.  See
// https://github.com/Cyphrme/Coze/issues/18. The last three characters of
// example `tmb` is `hOk`, but `hOl` also decodes to the same byte value (in
// Hex, `84E9`) even though they are different UTF-8 values. Tool for decoding
// [hOk](https://convert.zamicol.com/#?inAlph=base64&in=hOk&outAlph=Hex) and
// [hOl](https://convert.zamicol.com/#?inAlph=base64&in=hOl&outAlph=Hex).
func ExampleB64_non_strict_decode() {
	type Foo struct {
		Bar B64
	}

	// Canonical
	f := new(Foo)
	err := json.Unmarshal([]byte(`{"Bar":"hOk"}`), f)
	if err != nil {
		panic(err)
	}
	fmt.Println(f)

	// Non-canonical
	f2 := new(Foo)
	err = json.Unmarshal([]byte(`{"Bar":"hOl"}`), f2)
	if err != nil { // should error, but doesn't
		fmt.Println("unmarshalling error: ", err)
		return
	}

	// Output:
	// &{hOk}
	// unmarshalling error:  illegal base64 data at input byte 2
}
