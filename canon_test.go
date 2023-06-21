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
	canon := []string{"alg", "iat", "msg", "tmb", "typ"}
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
	// 4bmwgjkxQJIG2jiLiqq6eKptwTs97lYAFUtS25Rc3DU
	// Ie3xL77AsiCcb4r0pbnZJqMcfSBqg5Lk0npNJyJ9BC4
}

// Demonstrates expected behavior for invalid HashAlgs.
func ExampleCanonicalHash_invalidAlg() {
	_, err := CanonicalHash([]byte(GoldenPay), nil, "")
	fmt.Println(err)
	_, err = CanonicalHash([]byte(GoldenPay), nil, "test")
	fmt.Println(err)
	// Output:
	// Hash: invalid HashAlg ""
	// Hash: invalid HashAlg "test"
}

// Example CanonicalHash for all hashing algos.
func ExampleCanonicalHash_permutations() {
	canon := []string{"alg", "iat", "msg", "tmb", "typ"}
	algs := []string{"SHA-224", "SHA-256", "SHA-384", "SHA-512", "SHA3-224", "SHA3-256", "SHA3-384", "SHA3-512", "SHAKE128", "SHAKE256"}
	for _, alg := range algs {
		cad, err := CanonicalHash([]byte(GoldenPay), canon, ParseHashAlg(alg))
		if err != nil {
			panic(err)
		}
		fmt.Println(cad.String())
	}

	// Output:
	// cGCQ6FHj0fjAyYbvxS_8sfC0qTaSLJtcu0Xkhw
	// 4bmwgjkxQJIG2jiLiqq6eKptwTs97lYAFUtS25Rc3DU
	// WQiyyY5Ye2Y8vKcbANlmiXJkU-SVEgboYJg-wnrOKJ3v8PcI5XvQu_-C4yyGFrbW
	// irByY6uGnp6DrPvInvggL0ibo2p5yNvcuMVx1GiZoOArVIp4cGkAfB2FvknV5DyzKMHH-tV6vW8TyW7LZOyVFw
	// 9YyKIbtFYbSNqdwAXcwV0lwLb-X65k6zTBTWeQ
	// 8P9aSEJjC8tRKzfLNYBQTTXK-9E-DPlNaH_ikFkYUHQ
	// suqhBt29HS7c_wwDpcp943h0HlSI_FQdOkiz-Tjf9R_Wegil2pXHVxIFXkpOaceP
	// dAzMJWHLnGw9kjeo4RbVhzAAL6bwGasQbLFLZ1kHhdhGNNQm5nMib0cAQAAoIwdnKf0L8RADELg1XSFd8aJKww
	// QzbJ9ONj21KF3Zno1ctdIHfpGqmFGm11tinsAJUkYOg
	// MFYCTNhmavKZmFk_JNcttN9ccm4MAuYN3T868B2q0olpJ_6po2l98-617RfjnxkVuY2J--JjKt-KGi1S2RL4Bw
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
