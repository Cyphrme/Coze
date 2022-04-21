package enum

import (
	"fmt"
	"testing"
)

// BenchmarkNSV (New, Sign, Verify) will generate a new Crypto Key, sign a
// message with that key. verify the signature, and return the results.  It will
// also test verify digest.
// `go test -bench=.`
func BenchmarkNSV(b *testing.B) {
	var passCount = 0

	msg := []byte("Test message.")
	// TODO Ed25519 Support: Ed25519
	var algs = []SigAlg{ES224, ES256, ES384, ES512}

	for j := 0; j < b.N; j++ {
		for i := 0; i < len(algs); i++ {
			// log.Printf("Alg: %s\n", algs[i])
			cryptoKey, err := NewCryptoKey(SEAlg(algs[i]))
			//log.Printf("Alg: %+v, Key: %+v\n", cryptoKey.Alg, *cryptoKey.Private)
			if err != nil {
				panic("Could not generate a new valid Crypto Key.")
			}

			sig, err := cryptoKey.Sign(msg)
			if err != nil {
				panic(err)
			}

			valid, err := cryptoKey.Verify(msg, sig)
			if err != nil {
				panic(err)
			}
			if !valid {
				panic("The signature was invalid")
			}

			// Test VerifyDigest
			msgDigest := Hash(SigAlg(algs[i]).Hash(), msg)
			valid, _ = cryptoKey.VerifyDigest(msgDigest, sig)
			if !valid {
				panic("The signature was invalid")
			}

			passCount++
		}
	}

	fmt.Printf("TestCryptoKeyNSV Pass Count: %+v \n", passCount)
}

func ExampleHashAlg_print() {
	h := Sha256
	fmt.Println(h)
	// Output: SHA-256
}

func ExampleAlg_jsonMarshal() {
	type zStruct struct {
		A Alg `json:"alg"`
	}

	z := zStruct{A: Alg(ES256)}

	jm, err := Marshal(z)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s\n", jm)
	// Output: {"alg":"ES256"}
}

func ExampleHashAlg_jsonMarshal() {
	type zstruct = struct {
		H HashAlg `json:"hashAlg"`
	}
	z := zstruct{H: Sha256}

	jm, err := Marshal(z)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+s\n", jm)
	// Output: {"hashAlg":"SHA-256"}
}

func ExampleHashAlg_Parse() {
	h := new(HashAlg)

	hashes := []string{
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

	for _, hash := range hashes {
		h.Parse(hash)
		fmt.Println(h)
	}

	// Output:
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

func ExampleAlg_Params() {
	algs := []Alg{
		Alg(ES224), Alg(ES256), Alg(ES384), Alg(ES512), Alg(Ed25519),
		Alg(Sha224), Alg(Sha256), Alg(Sha384), Alg(Sha512),
		Alg(Sha3224), Alg(Sha3256), Alg(Sha3384), Alg(Sha3512),
		Alg(Shake128), Alg(Shake256),
	}
	fmt.Println(algs)

	for _, a := range algs {
		params := a.Params()

		b, _ := Marshal(params)
		fmt.Printf("%s\n", b)
	}

	// Output:
	// [ES224 ES256 ES384 ES512 Ed25519 SHA-224 SHA-256 SHA-384 SHA-512 SHA3-224 SHA3-256 SHA3-384 SHA3-512 SHAKE128 SHAKE256]
	// {"Name":"ES224","Genus":"ECDSA","Family":"EC","Hash":"SHA-224","Hash.Size":28,"Sig.Size":56,"Curve":"P-224","KeyUse":"sig"}
	// {"Name":"ES256","Genus":"ECDSA","Family":"EC","Hash":"SHA-256","Hash.Size":32,"Sig.Size":64,"Curve":"P-256","KeyUse":"sig"}
	// {"Name":"ES384","Genus":"ECDSA","Family":"EC","Hash":"SHA-384","Hash.Size":48,"Sig.Size":96,"Curve":"P-384","KeyUse":"sig"}
	// {"Name":"ES512","Genus":"ECDSA","Family":"EC","Hash":"SHA-512","Hash.Size":64,"Sig.Size":132,"Curve":"P-521","KeyUse":"sig"}
	// {"Name":"Ed25519","Genus":"EdDSA","Family":"EC","Hash":"SHA-512","Hash.Size":64,"Sig.Size":64,"Curve":"Curve25519","KeyUse":"sig"}
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
