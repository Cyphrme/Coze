package coze

import (
	"fmt"
	"testing"

	"github.com/cyphrme/coze/enum"
)

// TestCanonHash.
func TestCanonHash(t *testing.T) {
	// With Canon
	canon := []string{"alg", "iat", "msg", "tmb", "typ"}
	cad, err := CanonHash([]byte(Golden_Pay), canon, enum.Sha256)

	if err != nil {
		t.Fatal(err)
	}
	if cad.String() != "aC2YKfNvovfnZOw_RVxSEW6NeaUq41DZXX0oeaOboRg" {
		t.Fatal("canonical hash does not match.  Got: " + cad.String())
	}

	// Without canon
	cad, err = CanonHash([]byte(Golden_Pay), nil, enum.Sha256)

	if err != nil {
		t.Fatal(err)
	}
	if cad.String() != "yRpc2dKE_Z7vhOeWiKbtRka9vQIDlPxZsdjsAJZ2ZRM" {
		t.Fatal("canonical hash does not match.  Got: " + cad.String())
	}
}

// ExampleCanonical.
func ExampleCanonical() {
	var b []byte
	var err error

	type CBA struct {
		// Fields must be out of order for testing.
		C string `json:"c"`
		B string `json:"b"`
		A string `json:"a"`
	}

	type ABC struct {
		A string `json:"a"`
		B string `json:"b,omitempty"`
		C string `json:"c"`
	}

	// struct (out of order) with struct (in order) canon
	var cba = ABC{
		C: "c", B: "b", A: "a",
	}
	b, err = CanonicalStruct(cba, new(ABC))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Input: struct; Canon: struct => %s\n", b)

	// []byte (out of order) with nil canon.  Missing field  "b" should be omitted
	// from output.
	ca, err := Marshal(map[string]string{"c": "c", "a": "a"})
	if err != nil {
		fmt.Println(err)
	}

	b, err = Canonical(ca, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Input: []byte; Canon: nil    => %s\n", b)

	// []byte (out of order) with struct (in order) canon. Missing field "b"
	// should be omitted from output.
	ca, err = Marshal(map[string]string{"c": "c", "a": "a"})
	if err != nil {
		fmt.Println(err)
	}
	b, err = Canonical(ca, new(ABC))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Input: []byte; Canon: struct => %s\n", b)

	//////////////////
	// []byte (out of order) with struct (in order) canon
	//////////////////
	byteJson := []byte(`{"c":"c","a":"a"}`)
	b, err = Canonical(byteJson, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Input: []byte; Canon: nil    => %s\n", b)

	// Output:
	// Input: struct; Canon: struct => {"a":"a","b":"b","c":"c"}
	// Input: []byte; Canon: nil    => {"a":"a","c":"c"}
	// Input: []byte; Canon: struct => {"a":"a","c":"c"}
	// Input: []byte; Canon: nil    => {"a":"a","c":"c"}
}

func ExampleCanonStruct() {
	type structWithEmbedded struct {
		Pay
		Name string `json:"name"`
	}

	var s = structWithEmbedded{
		Pay: Pay{
			Alg: enum.SEAlg(enum.ES256),
			Iat: 1626479633,
			Typ: "cyphr.me",
		},
		Name: "Bob",
	}

	marshaled, err := Marshal(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", marshaled)

	can, err := CanonStruct(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(can)

	// Example with an empty struct demonstrating the behavior of
	// `json:omitempty`.
	ss := structWithEmbedded{}
	can, err = CanonStruct(ss)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(can)
	// Output:
	// {"alg":"ES256","iat":1626479633,"typ":"cyphr.me","name":"Bob"}
	// [alg iat name typ]
	// [name]

}

func ExampleCanon() {
	// Out of order since Canon should return in order
	b := []byte(`{"z":"z","a":"a"}`)

	can, err := Canon(b)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(can)
	// Output: [a z]
}
