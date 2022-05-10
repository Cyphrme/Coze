package coze

////////////////////////////////
// Keep package `basecnv` and `coze` hex_test.go in sync.
////////////////////////////////

import (
	"encoding/json"
	"fmt"
)

// Hex of nil is "" while Hex of 0 is "00".
func ExampleHex_zero_nil() {
	var b []byte
	b = nil
	fmt.Printf("Hex string nil: `%s`\n", Hex(b))

	b = []byte{0}
	fmt.Printf("Hex string zero: `%s`\n", Hex(b))

	// Output:
	// Hex string nil: ``
	// Hex string zero: `00`
}

func ExampleHex_marshalJSON() {
	h := Hex([]byte{0, 255})
	b, _ := h.MarshalJSON()
	stringed := string(b)
	fmt.Println(stringed)
	// Output: "00FF"
}

func ExampleHex_unmarshalJSON() {
	type Foo struct {
		Bar Hex
	}

	stringed := []byte(`{"Bar":"00FF"}`)
	f := new(Foo)

	err := json.Unmarshal(stringed, f)
	if err != nil {
		fmt.Println(err)
	}

	b, _ := json.Marshal(f)

	fmt.Printf("%s", b)
	// Output: {"Bar":"00FF"}
}

func ExampleHexEncodeString() {
	b := []byte{0, 255}
	fmt.Println(HexEncodeString(b))
	// Output: 00FF
}

// ExampleHexDecodeString decodes a string prints the Go string.
func ExampleHexDecodeString() {
	// Replace the string with what's wanting to be converted to bytes.
	b, err := HexDecodeString("064BC8ED150C7F0EED574688D5CE11E0F8B6E47CB0E247A882E1DCFBEDCF53AC")
	if err != nil {
		fmt.Println(err)
	}

	stringNum := "{"
	for _, byt := range b {
		stringNum += fmt.Sprintf("%d,", byt) // Go always has a trailing comma
	}
	stringNum += "}"

	fmt.Println(stringNum)
	// Output: {6,75,200,237,21,12,127,14,237,87,70,136,213,206,17,224,248,182,228,124,176,226,71,168,130,225,220,251,237,207,83,172,}
}

// ExampleMustHexDecodeString decomstrates use of Must.
func ExampleMustHexDecode() {
	// Replace the string with what's wanting to be converted to bytes.
	b := MustHexDecode("51E33CB2BF975D426FC349E04277E138AE4329EA2BD664E27D3AEA6DCB3AE199")

	fmt.Printf("%X\n", b)
	// Output: 51E33CB2BF975D426FC349E04277E138AE4329EA2BD664E27D3AEA6DCB3AE199
}

func ExampleHexDecodeString_odd() {
	// Input Hex is odd and should fail.
	_, err := HexDecodeString("00FFF")
	if err == nil {
		fmt.Println("error should not be nil")
	}

	b, _ := HexDecodeString("000FFF")
	fmt.Println(b)
	// Output: [0 15 255]
}
