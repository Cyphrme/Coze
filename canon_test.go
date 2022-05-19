package coze

import (
	"fmt"
	"testing"

	ce "github.com/cyphrme/coze/enum"
)

// TestCH tests CH, Canonical Hash.
func TestCH(t *testing.T) {
	//	The "canonical digest" of`head` is `cad`.
	cad, err := CH([]byte(Golden_Head_String), nil, ce.Sha256)

	if err != nil {
		t.Fatal(err)
	}
	if cad.String() != "0C359495353CD108BF5477F1084B4C2A656C565D2168D496A149B0990AE94286" {
		t.Fatal("canonical hash hex does not match")
	}
}

// ExampleCanonical tests CH (Canonical Hash), which calls Canon.
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

	//////////////////
	// struct (out of order) with struct (in order) canon
	/////////////////
	var cba = ABC{
		C: "c", B: "b", A: "a",
	}
	b, err = Canon(cba, new(ABC))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Input: struct; Canon: struct => %s\n", b)

	//////////////////
	// []byte (out of order) with nil canon.  Missing field  "b" should be omitted
	// from output.
	//////////////////
	ca, err := Marshal(map[string]string{"c": "c", "a": "a"})
	if err != nil {
		fmt.Println(err)
	}

	b, err = Canonical(ca, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Input: []byte; Canon: nil    => %s\n", b)

	//////////////////
	// []byte (out of order) with struct (in order) canon. Missing field "b"
	// should be omitted from output.
	//////////////////
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

func ExampleCanonS() {
	type structWithEmbedded struct {
		// Fields are out of order since CanonS must return in order.
		Name string `json:"name"`
		Head
	}

	var s = structWithEmbedded{
		Name: "Bob",
		Head: Head{
			Alg: ce.SEAlg(ce.ES256),
			Iat: 1626479633,
			Typ: "cyphr.me",
		},
	}

	marshaled, err := Marshal(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Marshaled: %s\n", marshaled)

	can, err := CanonS(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(can)

	// Example with an empty struct.  Since Typ tag is `json:omitempty`, it will
	// not br present in the canon.
	ss := structWithEmbedded{}
	can, err = CanonS(ss)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(can)
	// Output:
	// Marshaled: {"name":"Bob","alg":"ES256","iat":1626479633,"tmb":"","typ":"cyphr.me"}
	// [alg iat name tmb typ]
	// [alg iat name tmb]
}

func ExampleCanonB() {
	// Out of order since CanonB should return in order
	b := []byte(`{"z":"z","a":"a"}`)

	can, err := CanonB(b)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(can)
	// Output: [a z]
}
