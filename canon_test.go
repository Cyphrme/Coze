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

// Example CanonHash for all hashing algos.
func ExampleCanonHash_AllAlgs() {
	canon := []string{"alg", "iat", "msg", "tmb", "typ"}
	algs := []string{"SHA-224", "SHA-256", "SHA-384", "SHA-512", "SHA3-224", "SHA3-256", "SHA3-384", "SHA3-512", "SHAKE128", "SHAKE256"}
	for _, alg := range algs {
		cad, err := CanonHash([]byte(Golden_Pay), canon, ParseHashAlg(alg))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(cad.String())
	}

	// Output:
	// nrGQqKYvFKeVDFlOMIusP2A2AWn4DX-2XLNfJA
	// aC2YKfNvovfnZOw_RVxSEW6NeaUq41DZXX0oeaOboRg
	// KC5UHPxzNl567oONOphhWY6cHhuXoSiyOLiNGSmTIcvA8XDtWQf-fr4xNPCLfzCo
	// oX2NMgJ_QRW9rf59N5VOSMILg6mzVHld5CqRaOatLCbRWVRh1Y6Rq4tRZGzZNvNKEM0qbYBlWk6y9BcnuRzczA
	// GAYOBAxW2x7MvHVYpRLnjX3rUKcuhvDOCVVK3Q
	// UfFl2lw4KHc2-0GX-mnqtfpScM1Qf7L_IaTGojR6_Go
	// pVo43tSAG8apVs26QLOFG0Cbh3ScrbHd_VGjaFAIQtlCLiXcsgdmsGwOyXoK4zBz
	// IA8Xv6tt32B49THtWOzN9AyKtnG5a0v93DSF4IShHsT6S2lWKQl1H2yuyMAYocVKBkMF5dp0miKB58NXROqAMg
	// muWDwpDGlR-jwGPOQlj6A6B5FYA_U5nFq2KtwV8B-Uw
	// QPfIPjKmLO4qLmiClA6GjYQKBO6MI2wBZUhi9uVTVr0WGP3LgOTQRup6l5Caxz6GtiUnNeQe6JMdVSvhdvLW-Q
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
