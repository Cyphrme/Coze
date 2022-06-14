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
	stringed := string(b)
	fmt.Println(stringed)
	// Output: "AP8"
}

func ExampleB64_unmarshalJSON() {
	type Foo struct {
		Bar B64
	}

	stringed := []byte(`{"Bar":"AP8"}`)
	f := new(Foo)

	err := json.Unmarshal(stringed, f)
	if err != nil {
		fmt.Println(err)
	}

	b, err := Marshal(f)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s", b)
	// Output: {"Bar":"AP8"}
}
