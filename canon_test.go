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
	// _OPzRaT0b5iCQUreB4HXLeNZ-zrAmZKwOs2e9AZyH5Q
	// AyVZoWUv_rJf7_KqoeRS5odr8g3MZwBzhtBdSZderxk
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
	// AHp3PJhm7UqvInSyNx979L0UclaHerXpr65_RQ
	// _OPzRaT0b5iCQUreB4HXLeNZ-zrAmZKwOs2e9AZyH5Q
	// 2nZLq6SkucLHoQ2uzWsaDxxHqtgUsQuYROh6gLfsHJG4zD3615TchJjx2s53-jF-
	// DmHMOL8rbl4WREEcI5vZSFmhRLX1doGpXI6ValNwzP8jorZJ3qki5xtFM_0pZOp7tE59I6MM5N8KtMANt7axQw
	// UOmfhfx7LB_1556F_gql1i7XxK69eZ7lCaVrBQ
	// 3vcvftqd4lI2bp8s6dEiLVI_M5_4_usBUb9lZFLpm1E
	// 6sACHpDGK47H8DMh0vt42OkAHoAXZ5lSic0ju1a3UFPt3TcUFWYWM7K62uOjE-zp
	// uRWmti8R6KnCbmRP3dhrVEsN7daLwvW6Jq21e14_4lnvq2p9futNxoOLW0rL2-1VCJ11SAOoxcBfIgmYv-LFXw
	// yfPYoA6MT_QoyQGp0DDGRrGP_EHfV4-sojVRtplZmdk
	// zM8SoNSsYdJoyxW_83tR5L0axAWYvyPkSRWnMDYkEsV9pfSgGhB05BPo5xJyslpCqxnZn9ETxgwuOLvCwBInpg
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
