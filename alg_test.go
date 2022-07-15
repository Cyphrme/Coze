package coze

import (
	"fmt"
)

// Signing and verifying tests are done in package Coze.

func ExampleHashAlg_print() {
	h := SHA256
	fmt.Println(h)
	// Output: SHA-256
}

func ExampleAlg_jsonMarshal() {
	type zStruct struct {
		A Alg `json:"alg"`
	}

	jm, err := Marshal(zStruct{A: Alg(ES256)})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s\n", jm)
	// Output: {"alg":"ES256"}
}

func ExampleHashAlg_jsonMarshal() {
	type testStruct = struct {
		H HashAlg `json:"hashAlg"`
	}
	z := testStruct{H: SHA256}
	jm, err := Marshal(z)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+s\n", jm)
	// Output: {"hashAlg":"SHA-256"}
}

func ExampleAlg_Parse() {
	a := new(Alg)

	algs := []string{
		// SEAlgs
		"ES224",
		"ES256",
		"ES384",
		"ES512",
		"Ed25519",
		"Ed25519ph",
		"Ed448",
		// Hash algs
		"SHA-224",
		"SHA-256",
		"SHA-384",
		"SHA-512",
		"SHA3-224",
		"SHA3-256",
		"SHA3-384",
		"SHA3-512",
		"SHAKE128",
		"SHAKE256",
	}

	for _, alg := range algs {
		a.Parse(alg) // Call as method
		fmt.Println(a)
	}

	// Output:
	// ES224
	// ES256
	// ES384
	// ES512
	// Ed25519
	// Ed25519ph
	// Ed448
	// SHA-224
	// SHA-256
	// SHA-384
	// SHA-512
	// SHA3-224
	// SHA3-256
	// SHA3-384
	// SHA3-512
	// SHAKE128
	// SHAKE256
}

func ExampleCrv_Parse() {
	c := new(Crv)

	crvs := []string{
		"P-224",
		"P-256",
		"P-384",
		"P-521",
		"Curve25519",
		"Curve448",
	}

	for _, crv := range crvs {
		c.Parse(crv)
		fmt.Println(c)
	}

	// Output:
	// P-224
	// P-256
	// P-384
	// P-521
	// Curve25519
	// Curve448
}

func ExampleKeyUse_Parse() {
	u := new(KeyUse)

	uses := []string{
		"sig",
		"enc",
	}

	for _, use := range uses {
		u.Parse(use)
		fmt.Println(u)
	}

	// Output:
	// sig
	// enc
}

func ExampleAlg_Params() {
	algs := []Alg{
		Alg(ES224), Alg(ES256), Alg(ES384), Alg(ES512), Alg(Ed25519), Alg(Ed25519ph),
		Alg(Ed448), Alg(SHA224), Alg(SHA256), Alg(SHA384), Alg(SHA512), Alg(SHA3224),
		Alg(SHA3256), Alg(SHA3384), Alg(SHA3512), Alg(SHAKE128), Alg(SHAKE256),
	}
	fmt.Println(algs)

	for _, a := range algs {
		params := a.Params()

		b, _ := Marshal(params)
		fmt.Printf("%s\n", b)
	}

	// Output:
	// [ES224 ES256 ES384 ES512 Ed25519 Ed25519ph Ed448 SHA-224 SHA-256 SHA-384 SHA-512 SHA3-224 SHA3-256 SHA3-384 SHA3-512 SHAKE128 SHAKE256]
	// {"Name":"ES224","Genus":"ECDSA","Family":"EC","X.Size":56,"D.Size":28,"Hash":"SHA-224","Hash.Size":28,"Sig.Size":56,"Curve":"P-224","Use":"sig"}
	// {"Name":"ES256","Genus":"ECDSA","Family":"EC","X.Size":64,"D.Size":32,"Hash":"SHA-256","Hash.Size":32,"Sig.Size":64,"Curve":"P-256","Use":"sig"}
	// {"Name":"ES384","Genus":"ECDSA","Family":"EC","X.Size":96,"D.Size":48,"Hash":"SHA-384","Hash.Size":48,"Sig.Size":96,"Curve":"P-384","Use":"sig"}
	// {"Name":"ES512","Genus":"ECDSA","Family":"EC","X.Size":132,"D.Size":66,"Hash":"SHA-512","Hash.Size":64,"Sig.Size":132,"Curve":"P-521","Use":"sig"}
	// {"Name":"Ed25519","Genus":"EdDSA","Family":"EC","X.Size":32,"D.Size":32,"Hash":"SHA-512","Hash.Size":64,"Sig.Size":64,"Curve":"Curve25519","Use":"sig"}
	// {"Name":"Ed25519ph","Genus":"EdDSA","Family":"EC","X.Size":32,"D.Size":32,"Hash":"SHA-512","Hash.Size":64,"Sig.Size":64,"Curve":"Curve25519","Use":"sig"}
	// {"Name":"Ed448","Genus":"EdDSA","Family":"EC","X.Size":57,"D.Size":57,"Hash":"SHAKE256","Hash.Size":64,"Sig.Size":114,"Curve":"Curve448","Use":"sig"}
	// {"Name":"SHA-224","Genus":"SHA2","Family":"SHA","Hash":"SHA-224","Hash.Size":28}
	// {"Name":"SHA-256","Genus":"SHA2","Family":"SHA","Hash":"SHA-256","Hash.Size":32}
	// {"Name":"SHA-384","Genus":"SHA2","Family":"SHA","Hash":"SHA-384","Hash.Size":48}
	// {"Name":"SHA-512","Genus":"SHA2","Family":"SHA","Hash":"SHA-512","Hash.Size":64}
	// {"Name":"SHA3-224","Genus":"SHA3","Family":"SHA","Hash":"SHA3-224","Hash.Size":28}
	// {"Name":"SHA3-256","Genus":"SHA3","Family":"SHA","Hash":"SHA3-256","Hash.Size":32}
	// {"Name":"SHA3-384","Genus":"SHA3","Family":"SHA","Hash":"SHA3-384","Hash.Size":48}
	// {"Name":"SHA3-512","Genus":"SHA3","Family":"SHA","Hash":"SHA3-512","Hash.Size":64}
	// {"Name":"SHAKE128","Genus":"SHA3","Family":"SHA","Hash":"SHAKE128","Hash.Size":32}
	// {"Name":"SHAKE256","Genus":"SHA3","Family":"SHA","Hash":"SHAKE256","Hash.Size":64}
}
