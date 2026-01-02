package coz

import (
	"fmt"
)

func ExampleAlg_jsonMarshal() {
	type algStruct struct {
		A Alg `json:"alg"`
	}

	b, err := Marshal(algStruct{A: Alg(ES256)})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)

	type seAlgStruct struct {
		A SEAlg `json:"alg"`
	}

	b, err = Marshal(seAlgStruct{A: SEAlg(ES256)})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)

	// Output:
	// {"alg":"ES256"}
	// {"alg":"ES256"}
}

func ExampleAlg_Parse() {
	algs := []string{
		"",
		"foo",
		"UnknownAlg",
		"UnknownSigAlg",
		"ES224",
		"ES256",
		"ES384",
		"ES512",
		"Ed25519",
		"Ed25519ph",
		"Ed448",
		"UnknownEncAlg",
		"UnknownHshAlg",
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

	var a Alg
	for _, alg := range algs {
		a.Parse(alg) // Call as method
		fmt.Println(a)
	}

	// Output:
	// UnknownAlg
	// UnknownAlg
	// UnknownAlg
	// UnknownSigAlg
	// ES224
	// ES256
	// ES384
	// ES512
	// Ed25519
	// Ed25519ph
	// Ed448
	// UnknownEncAlg
	// UnknownHshAlg
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

func ExampleHshAlg_print() {
	h := SHA256
	fmt.Println(h)

	// Output:
	// SHA-256
}

func ExampleHshAlg_jsonMarshal() {
	type testStruct = struct {
		H HshAlg `json:"hshAlg"`
	}
	z := testStruct{H: SHA256}
	jm, err := Marshal(z)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+s\n", jm)

	// Output:
	// {"hshAlg":"SHA-256"}
}

func ExampleCrv_Parse() {
	var c Crv

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

func ExampleUse_Parse() {
	var u Use
	uses := []string{
		"sig",
		"enc",
		"hsh",
	}
	for _, use := range uses {
		u.Parse(use)
		fmt.Println(u)
	}

	// Output:
	// sig
	// enc
	// hsh
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
	// {"Name":"ES224","Genus":"ECDSA","Family":"EC","Use":"sig","Hash":"SHA-224","HashSize":28,"HashSizeB64":38,"PubSize":56,"PubSizeB64":75,"PrvSize":28,"PrvSizeB64":38,"Curve":"P-224","SigSize":56,"SigSizeB64":75}
	// {"Name":"ES256","Genus":"ECDSA","Family":"EC","Use":"sig","Hash":"SHA-256","HashSize":32,"HashSizeB64":43,"PubSize":64,"PubSizeB64":86,"PrvSize":32,"PrvSizeB64":43,"Curve":"P-256","SigSize":64,"SigSizeB64":86}
	// {"Name":"ES384","Genus":"ECDSA","Family":"EC","Use":"sig","Hash":"SHA-384","HashSize":48,"HashSizeB64":64,"PubSize":96,"PubSizeB64":128,"PrvSize":48,"PrvSizeB64":64,"Curve":"P-384","SigSize":96,"SigSizeB64":128}
	// {"Name":"ES512","Genus":"ECDSA","Family":"EC","Use":"sig","Hash":"SHA-512","HashSize":64,"HashSizeB64":86,"PubSize":132,"PubSizeB64":176,"PrvSize":66,"PrvSizeB64":88,"Curve":"P-521","SigSize":132,"SigSizeB64":176}
	// {"Name":"Ed25519","Genus":"EdDSA","Family":"EC","Use":"sig","Hash":"SHA-512","HashSize":64,"HashSizeB64":86,"PubSize":32,"PubSizeB64":43,"PrvSize":32,"PrvSizeB64":43,"Curve":"Curve25519","SigSize":64,"SigSizeB64":86}
	// {"Name":"Ed25519ph","Genus":"EdDSA","Family":"EC","Use":"sig","Hash":"SHA-512","HashSize":64,"HashSizeB64":86,"PubSize":32,"PubSizeB64":43,"PrvSize":32,"PrvSizeB64":43,"Curve":"Curve25519","SigSize":64,"SigSizeB64":86}
	// {"Name":"Ed448","Genus":"EdDSA","Family":"EC","Use":"sig","Hash":"SHAKE256","HashSize":64,"HashSizeB64":86,"PubSize":57,"PubSizeB64":76,"PrvSize":57,"PrvSizeB64":76,"Curve":"Curve448","SigSize":114,"SigSizeB64":152}
	// {"Name":"SHA-224","Genus":"SHA2","Family":"SHA","Use":"hsh","Hash":"SHA-224","HashSize":28,"HashSizeB64":38}
	// {"Name":"SHA-256","Genus":"SHA2","Family":"SHA","Use":"hsh","Hash":"SHA-256","HashSize":32,"HashSizeB64":43}
	// {"Name":"SHA-384","Genus":"SHA2","Family":"SHA","Use":"hsh","Hash":"SHA-384","HashSize":48,"HashSizeB64":64}
	// {"Name":"SHA-512","Genus":"SHA2","Family":"SHA","Use":"hsh","Hash":"SHA-512","HashSize":64,"HashSizeB64":86}
	// {"Name":"SHA3-224","Genus":"SHA3","Family":"SHA","Use":"hsh","Hash":"SHA3-224","HashSize":28,"HashSizeB64":38}
	// {"Name":"SHA3-256","Genus":"SHA3","Family":"SHA","Use":"hsh","Hash":"SHA3-256","HashSize":32,"HashSizeB64":43}
	// {"Name":"SHA3-384","Genus":"SHA3","Family":"SHA","Use":"hsh","Hash":"SHA3-384","HashSize":48,"HashSizeB64":64}
	// {"Name":"SHA3-512","Genus":"SHA3","Family":"SHA","Use":"hsh","Hash":"SHA3-512","HashSize":64,"HashSizeB64":86}
	// {"Name":"SHAKE128","Genus":"SHA3","Family":"SHA","Use":"hsh","Hash":"SHAKE128","HashSize":32,"HashSizeB64":43}
	// {"Name":"SHAKE256","Genus":"SHA3","Family":"SHA","Use":"hsh","Hash":"SHAKE256","HashSize":64,"HashSizeB64":86}
}
