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
		fmt.Println(err)
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
		fmt.Println(err)
	}

	b, err := Marshal(f)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s\n%#v\n", b, B64(b))
	// Output:
	// {"Bar":"AP8"}
	// eyJCYXIiOiJBUDgifQ
}

func ExampleDecode() {
	b, err := Decode(Golden_Tmb)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
	// Output: cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk
}

func ExampleMustDecode() {
	fmt.Println(MustDecode(Golden_Tmb))
	// Output: cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk
}
