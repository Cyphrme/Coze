package coz

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"testing"
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

func ExampleDecode() {
	b, err := Decode(GoldenTmb)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
	// Output: U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg
}

func ExampleMustDecode() {
	fmt.Println(MustDecode(GoldenTmb))
	// Output: U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg
}

type B64Struct struct {
	B B64
}

func ExampleB64_unmarshalJSON() {
	f := new(B64Struct)
	err := json.Unmarshal([]byte(`{"B":"AP8"}`), f)
	if err != nil {
		panic(err)
	}

	b, err := Marshal(f)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s,%#v\n", b, B64(b))
	// Output:
	// {"B":"AP8"},eyJCIjoiQVA4In0
}

// Demonstrates that Coze Go will error on non-canonical base 64 encoding.  See
// https://github.com/Cyphrme/Coze/issues/18. The last three characters of
// example `tmb` is `hOk`, but `hOl` also decodes to the same byte value (in
// Hex, `84E9`) even though they are different UTF-8 values. Tool for decoding
// [hOk](https://convert.zamicol.com/#?inAlph=base64&in=hOk&outAlph=Hex) and
// [hOl](https://convert.zamicol.com/#?inAlph=base64&in=hOl&outAlph=Hex).
//
// As an added concern, Go's base64 ignores new line and carriage return.
// Thankfully, JSON unmarshal does not, making Coze's interpretation of base 64
// non-malleable since Coze is JSON.
func ExampleB64_non_strict_decode() {
	// Canonical
	f := new(B64Struct)
	err := json.Unmarshal([]byte(`{"B":"hOk"}`), f)
	if err != nil {
		panic(err)
	}
	fmt.Println(f)

	// Non-canonical (hOk and hOl will decode to the same bytes when non-canonical
	// is permitted.)
	f2 := new(B64Struct)
	err = json.Unmarshal([]byte(`{"B":"hOl"}`), f2)
	if err != nil { // Correctly errors
		fmt.Println("unmarshalling error: ", err)
	}

	// Print Unicode to show that Go is interpreting the string below correctly.
	b1 := []byte(fmt.Sprintf(`{"B":"hOk"}`))
	b2 := []byte(fmt.Sprintf("{\"B\":\"hOk\n\"}")) // Unicode U+000A is line feed.
	b3 := []byte(fmt.Sprintf("{\"B\":\"hOk\r\"}")) // Unicode U+000D is line feed.

	fmt.Printf("%U\n", b1)
	fmt.Printf("%U\n", b2)
	fmt.Printf("%U\n", b3)

	fb1 := new(B64Struct)
	err = json.Unmarshal(b1, fb1) // Will not error
	if err != nil {
		fmt.Println(err)
	}
	fb2 := new(B64Struct)
	err = json.Unmarshal(b2, fb2) // Correctly errors.
	if err != nil {
		fmt.Println(err)
	}
	fb3 := new(B64Struct)
	err = json.Unmarshal(b3, fb3) // Correctly errors.
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// &{hOk}
	// unmarshalling error:  illegal base64 data at input byte 2
	// [U+007B U+0022 U+0042 U+0022 U+003A U+0022 U+0068 U+004F U+006B U+0022 U+007D]
	// [U+007B U+0022 U+0042 U+0022 U+003A U+0022 U+0068 U+004F U+006B U+000A U+0022 U+007D]
	// [U+007B U+0022 U+0042 U+0022 U+003A U+0022 U+0068 U+004F U+006B U+000D U+0022 U+007D]
	// invalid character '\n' in string literal
	// invalid character '\r' in string literal
}

// FuzzCastB64ToString ensures that casting to and from B64 and string does not
// cause unexpected issues (issues like replacing bytes with the unicode
// replacement character).
// https://go.dev/security/fuzz/
// https://go.dev/doc/tutorial/fuzz
func FuzzCastB64ToString(f *testing.F) {
	f.Add(100)
	f.Fuzz(func(t *testing.T, a int) {
		var b B64
		var err error
		for i := 0; i < a; i++ {
			b = make([]byte, 32)
			_, err = rand.Read(b)
			if err != nil {
				t.Fatal(err)
			}

			s := string(b)
			bb := B64(s)
			if !bytes.Equal(b, bb) {
				t.Fatalf("Casting to string: %s failed when converting back B64.", s)
			}
		}
	})
}
