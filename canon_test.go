package coz

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
	// RbiD-Xsi0UoKI0YyDRysM6n7_hRHEOMHcE3WbmgWPxg
	// XzrXMGnY0QFwAKkr43Hh-Ku3yUS8NVE0BdzSlMLSuTU
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
	// srWf7I1Pbi8seVq2kMqit400wFrUxKTynfGyJg
	// RbiD-Xsi0UoKI0YyDRysM6n7_hRHEOMHcE3WbmgWPxg
	// WjLurqzLxPCLtjfTY2GXoCO9YnaMQv-1_6hzMkWN1UNULxLC1jidFqBIFupItTZf
	// 42XMUM7-vA2qfFdpf3aQ8Qf1d0pf3mCWtxDJHcrgu9JunN-KQjWhmhTd9LzEl4APMSE0_OvZZa2c80mWAzGr1w
	// 0DSloPImruVlJ1j7_f99HdtyTzqzQ-QOXmHcVw
	// QSxrw98yTsbWJkl2oB_n3dlxfBLhkazFnlwoLdLtjGA
	// G6Jpl9JmOlyXe_-oOJxXexvblftekqeU0fEl-7Q9JB-96JgpjXqprHjVlE0TJnaH
	// XuaKEaN9dcfLkMBqcbhdHCPST5jeykrACQ01uWvgjF0AIaE71wcDyA1DVoOcdiolpITPiNc_aSMDe7vP82cBnw
	// KlK05pQJJJVRFbEL6ux3Z3BfAXJ7932gjJwlE_GPqNE
	// bRiwjhCdHUmbnpujKtCT9QqJmKPjmK9J44UACMIgSf2lwqw_qxsxD-Xyir7IoWjuBlN24NzpvfTvV_NRSYgO3Q
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
	// KeyCanon is the canonical form of a Coz key in struct form.
	type KeyCanonStruct struct {
		Alg string `json:"alg"`
		Pub B64    `json:"pub"`
	}
	kcs := new(KeyCanonStruct)

	dig, err := CanonicalHash([]byte(GoldenKeyString), kcs, GoldenKey.Alg.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Println(dig)

	// Output:
	// U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg
}

func ExampleCanonical_slice() {
	dig, err := CanonicalHash([]byte(GoldenKeyString), KeyCanon, GoldenKey.Alg.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Println(dig)

	// Output: U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg
}
