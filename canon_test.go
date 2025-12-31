package coze

import (
	"fmt"
)

func ExampleCanon() {
	b := []byte(`{"z":"z", "a":"a"}`)

	can, err := Canon(b)
	if err != nil {
		panic(err)
	}

	fmt.Println(can)
	// Output:
	// [z a]
}

// ExampleCanonicalHash. See also Example_genCad
func ExampleCanonicalHash() {
	canon := []string{"alg", "now", "msg", "tmb", "typ"}
	cad, err := CanonicalHash([]byte(GoldenPay), canon, SHA256)
	if err != nil {
		panic(err)
	}
	fmt.Println(cad.String())

	// Without canon
	cad, err = CanonicalHash([]byte(GoldenPay), nil, SHA256)
	if err != nil {
		panic(err)
	}
	fmt.Println(cad.String())

	// Output:
	// r8hLiXxSW0I8U-me0duNdodgyVgHPPTatxVAYJYg8BM
	// Y6e186hFUbPDDjjJxgXuEvsJmX6-lvkeqhzTf1MDQS4
}

// Demonstrates expected behavior for invalid HshAlgs.
func ExampleCanonicalHash_invalidAlg() {
	_, err := CanonicalHash([]byte(GoldenPay), nil, "")
	fmt.Println(err)
	_, err = CanonicalHash([]byte(GoldenPay), nil, "test")
	fmt.Println(err)
	// Output:
	// Hash: invalid HshAlg ""
	// Hash: invalid HshAlg "test"
}

// Example CanonicalHash for all hashing algos.
func ExampleCanonicalHash_permutations() {
	canon := []string{"alg", "now", "msg", "tmb", "typ"}
	algs := []string{"SHA-224", "SHA-256", "SHA-384", "SHA-512", "SHA3-224", "SHA3-256", "SHA3-384", "SHA3-512", "SHAKE128", "SHAKE256"}
	for _, alg := range algs {
		cad, err := CanonicalHash([]byte(GoldenPay), canon, ParseHashAlg(alg))
		if err != nil {
			panic(err)
		}
		fmt.Println(cad.String())
	}

	// Output:
	// XscoLNnxTsPcit0fVgLT90XgDgkxnVBmMEs02w
	// r8hLiXxSW0I8U-me0duNdodgyVgHPPTatxVAYJYg8BM
	// VT-K_sJplYTMv2Gi7QyUE0ja9Usgj9YHfGCmDvMqPBDzYDPlGD-55MsLbnPY1pi3
	// QY0ybderCLD0Lu7w2vXNq4XEtfNVPexphJPoiBet9Ly2vXZibo59fITNG39tL3fIwtjVWqDFB34GOajkeW1g5w
	// vD-BBNYrBznF6oCL03XZFM3HM4uGZ36KxRolGg
	// 2sHfdx0Gg83_5vM3WJugtWBLDJdnR2zVAeRZPF6_v-8
	// g-4N-CIuMcfM8aZIxGVzNmo6QDRVv_NHs-S-8WfNKSyCTjMPE14T02Iyzf52kLdf
	// ZmUwWY2kyuWJxJLxJHYVDG3Zrt6lcOErVuXKoYxz21AUDqpcGAFzQ_WZZLXIAoQhtuF1XtiIJpmS89KQxWon5A
	// S6du7_QxmI9_Oqp-iXNKhTST-VdqX4dBZSKp7tl4ktA
	// fm38qGtplvLccc8_pvHv1wVy16C9mnsj05tsL53Gd7TC89nh6S2NykepY6BeyLMPeMorcOH1X5E5fqkpIjjjvA
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
		panic(err)
	}

	b, err = Canonical(ca, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Input: []byte; Canon: nil    => %s\n", b)

	// []byte (out of order) with struct (in order) canon. Missing field "b"
	// should be omitted from output.
	ca, err = Marshal(map[string]string{"c": "c", "a": "a"})
	if err != nil {
		panic(err)
	}
	b, err = Canonical(ca, new(ABC))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Input: []byte; Canon: struct => %s\n", b)

	// []byte (out of order) with struct (in order) canon.
	byteJSON := []byte(`{"c":"c", "a": "a"}`)
	b, err = Canonical(byteJSON, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Input: []byte; Canon: nil    => %s\n", b)

	// Output:
	// Input: []byte; Canon: nil    => {"a":"a","c":"c"}
	// Input: []byte; Canon: struct => {"a":"a","c":"c"}
	// Input: []byte; Canon: nil    => {"c":"c","a":"a"}
}

// ExampleCanonical_struct demonstrates using a given struct as a canon.
func ExampleCanonical_struct() {
	// KeyCanon is the canonical form of a Coze key in struct form.
	type KeyCanonStruct struct {
		Alg string `json:"alg"`
		X   B64    `json:"x"`
	}
	kcs := new(KeyCanonStruct)

	dig, err := CanonicalHash([]byte(GoldenKeyString), kcs, GoldenKey.Alg.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Println(dig)

	// Output:
	// cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk
}

func ExampleCanonical_slice() {
	dig, err := CanonicalHash([]byte(GoldenKeyString), KeyCanon, GoldenKey.Alg.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Println(dig)

	// Output: cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk
}
