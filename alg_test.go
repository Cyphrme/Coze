package coze

import (
	"fmt"
	"strconv"
	"testing"

	"golang.org/x/exp/slices"
)

func ExampleAlgs_print() {
	fmt.Println(Algs)

	// Output:
	// [UnknownAlg UnknownSigAlg ES224 ES256 ES384 ES512 Ed25519 Ed25519ph Ed448 UnknownEncAlg UnknownHashAlg SHA-224 SHA-256 SHA-384 SHA-512 SHA3-224 SHA3-256 SHA3-384 SHA3-512 SHAKE128 SHAKE256]
}

// Tests to make sure the alg enums are in order according to Parse.
func TestParse(t *testing.T) {
	algs := []Alg{}
	for _, a := range Algs {
		algs = append(algs, Parse(string(a)))
	}
	if !slices.Equal(Algs, algs) {
		t.Fatal("slices not equal. ")
	}
}

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
		"UnknownHashAlg",
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
	// UnknownHashAlg
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
		H HshAlg `json:"hashAlg"`
	}
	z := testStruct{H: SHA256}
	jm, err := Marshal(z)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+s\n", jm)

	// Output:
	// {"hashAlg":"SHA-256"}
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
	// {"Name":"ES224","Genus":"ECDSA","Family":"EC","Use":"sig","Hash":"SHA-224","HashSize":28,"HashSizeB64":38,"XSize":56,"XSizeB64":75,"DSize":28,"DSizeB64":38,"Curve":"P-224","SigSize":56,"SigSizeB64":75}
	// {"Name":"ES256","Genus":"ECDSA","Family":"EC","Use":"sig","Hash":"SHA-256","HashSize":32,"HashSizeB64":43,"XSize":64,"XSizeB64":86,"DSize":32,"DSizeB64":43,"Curve":"P-256","SigSize":64,"SigSizeB64":86}
	// {"Name":"ES384","Genus":"ECDSA","Family":"EC","Use":"sig","Hash":"SHA-384","HashSize":48,"HashSizeB64":64,"XSize":96,"XSizeB64":128,"DSize":48,"DSizeB64":64,"Curve":"P-384","SigSize":96,"SigSizeB64":128}
	// {"Name":"ES512","Genus":"ECDSA","Family":"EC","Use":"sig","Hash":"SHA-512","HashSize":64,"HashSizeB64":86,"XSize":132,"XSizeB64":176,"DSize":66,"DSizeB64":88,"Curve":"P-521","SigSize":132,"SigSizeB64":176}
	// {"Name":"Ed25519","Genus":"EdDSA","Family":"EC","Use":"sig","Hash":"SHA-512","HashSize":64,"HashSizeB64":86,"XSize":32,"XSizeB64":43,"DSize":32,"DSizeB64":43,"Curve":"Curve25519","SigSize":64,"SigSizeB64":86}
	// {"Name":"Ed25519ph","Genus":"EdDSA","Family":"EC","Use":"sig","Hash":"SHA-512","HashSize":64,"HashSizeB64":86,"XSize":32,"XSizeB64":43,"DSize":32,"DSizeB64":43,"Curve":"Curve25519","SigSize":64,"SigSizeB64":86}
	// {"Name":"Ed448","Genus":"EdDSA","Family":"EC","Use":"sig","Hash":"SHAKE256","HashSize":64,"HashSizeB64":86,"XSize":57,"XSizeB64":76,"DSize":57,"DSizeB64":76,"Curve":"Curve448","SigSize":114,"SigSizeB64":152}
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

func Example_curveOrder() {
	algs := []SigAlg{
		ES224,
		ES256,
		ES384,
		ES512,
	}

	for _, a := range algs {
		hexSize := Alg(a).Params().XSize
		fmt.Printf("%0"+strconv.Itoa(hexSize)+"X\n", curveOrders[a])
	}
	fmt.Println()
	for _, a := range algs {

		hexSize := Alg(a).Params().XSize
		fmt.Printf("%0"+strconv.Itoa(hexSize)+"X\n", curveHalfOrders[a])
	}

	// Output:
	//
	// FFFFFFFFFFFFFFFFFFFFFFFFFFFF16A2E0B8F03E13DD29455C5C2A3D
	// FFFFFFFF00000000FFFFFFFFFFFFFFFFBCE6FAADA7179E84F3B9CAC2FC632551
	// FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFC7634D81F4372DDF581A0DB248B0A77AECEC196ACCC52973
	// 01FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFA51868783BF2F966B7FCC0148F709A5D03BB5C9B8899C47AEBB6FB71E91386409
	//
	// 7FFFFFFFFFFFFFFFFFFFFFFFFFFF8B51705C781F09EE94A2AE2E151E
	// 7FFFFFFF800000007FFFFFFFFFFFFFFFDE737D56D38BCF4279DCE5617E3192A8
	// 7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFE3B1A6C0FA1B96EFAC0D06D9245853BD76760CB5666294B9
	// 00FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFD28C343C1DF97CB35BFE600A47B84D2E81DDAE4DC44CE23D75DB7DB8F489C3204

}
