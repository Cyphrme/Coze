package coze

import (
	"fmt"
)

func ExampleCanon() {
	b := []byte(`{"z":"z", "a":"a"}`)

	can, err := Canon(b)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(can)
	// Output: [z a]
}

// ExampleCanonHash. See also Example_genCad
func ExampleCanonHash() {
	canon := []string{"alg", "iat", "msg", "tmb", "typ"}
	cad, err := CanonHash([]byte(Golden_Pay), canon, SHA256)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cad.String())

	// Without canon
	cad, err = CanonHash([]byte(Golden_Pay), nil, SHA256)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cad.String())

	// Output:
	// aC2YKfNvovfnZOw_RVxSEW6NeaUq41DZXX0oeaOboRg
	// LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto

}

// ExampleCanonical.
func ExampleCanonical() {
	var b []byte
	var err error

	type ABC struct {
		A string `json:"a"`
		B string `json:"b,omitempty"`
		C string `json:"c"`
	}

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
	byteJson := []byte(`{"c":"c", "a": "a"}`)
	b, err = Canonical(byteJson, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Input: []byte; Canon: nil    => %s\n", b)

	// Output:
	// Input: []byte; Canon: nil    => {"a":"a","c":"c"}
	// Input: []byte; Canon: struct => {"a":"a","c":"c"}
	// Input: []byte; Canon: nil    => {"c":"c","a":"a"}
}
