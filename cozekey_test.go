package coze

import (
	"encoding/json"
	"fmt"
	"testing"
)

var testDigest = []byte{112, 184, 252, 190, 198, 45, 48, 28, 24, 147, 58, 5, 85, 145, 193, 102, 142, 146, 52, 191, 48, 73, 208, 136, 140, 34, 128, 193, 115, 110, 132, 233}

var Golden_Key = CozeKey{
	Alg: SEAlg(ES256),
	Kid: "Zami's Majuscule Key.",
	Iat: 1623132000,
	X:   []byte{218, 116, 206, 104, 85, 102, 217, 2, 241, 153, 67, 191, 74, 56, 50, 177, 197, 71, 6, 219, 199, 17, 250, 54, 174, 174, 185, 50, 248, 13, 70, 51, 145, 162, 58, 183, 244, 118, 170, 246, 181, 205, 198, 245, 241, 193, 182, 191, 94, 61, 5, 230, 246, 98, 108, 148, 119, 138, 192, 93, 57, 102, 232, 230},
	D:   []byte{108, 219, 45, 131, 143, 199, 222, 109, 210, 149, 19, 174, 127, 4, 82, 18, 8, 155, 46, 176, 110, 70, 175, 117, 215, 131, 175, 117, 170, 92, 165, 80},
	Tmb: []byte{112, 184, 252, 190, 198, 45, 48, 28, 24, 147, 58, 5, 85, 145, 193, 102, 142, 146, 52, 191, 48, 73, 208, 136, 140, 34, 128, 193, 115, 110, 132, 233},
}

const Golden_Key_String = `{
	"alg":"ES256",
	"iat":1623132000,
	"kid":"Zami's Majuscule Key.",
	"d":"bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVA",
	"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
	"x":"2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"
}`

// The very last byte in D was changed from 80, to 81, making it invalid.
// Base64 needs to be to "E", not B-D for it to be effective: bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVE
var Golden_Bad_Key = CozeKey{
	Alg: SEAlg(ES256),
	Kid: "Zami's Majuscule Key.",
	Iat: 1623132000,
	X:   []byte{218, 116, 206, 104, 85, 102, 217, 2, 241, 153, 67, 191, 74, 56, 50, 177, 197, 71, 6, 219, 199, 17, 250, 54, 174, 174, 185, 50, 248, 13, 70, 51, 145, 162, 58, 183, 244, 118, 170, 246, 181, 205, 198, 245, 241, 193, 182, 191, 94, 61, 5, 230, 246, 98, 108, 148, 119, 138, 192, 93, 57, 102, 232, 230},
	D:   []byte{108, 219, 45, 131, 143, 199, 222, 109, 210, 149, 19, 174, 127, 4, 82, 18, 8, 155, 46, 176, 110, 70, 175, 117, 215, 131, 175, 117, 170, 92, 165, 81},
	Tmb: []byte{112, 184, 252, 190, 198, 45, 48, 28, 24, 147, 58, 5, 85, 145, 193, 102, 142, 146, 52, 191, 48, 73, 208, 136, 140, 34, 128, 193, 115, 110, 132, 233},
}

var Golden_Pay = `{
	"msg": "Coze Rocks",
	"alg": "ES256",
	"iat": 1627518000,
	"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
	"typ": "cyphr.me/msg"
 }`

var Golden_Tmb = "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk"
var Golden_Cad = "LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto"
var Golden_Czd = "d0ygwQCGzuxqgUq1KsuAtJ8IBu0mkgAcKpUJzuX075M"
var Golden_Sig = "ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"

var Golden_Coze = `{
	"pay":` + Golden_Pay + `,
	"sig": "` + Golden_Sig + `"
 }`

var Golden_CozeE = `{
	"coze":` + Golden_Coze + `
}`

var Golden_Coze_W_Key = `{
	"pay": ` + Golden_Pay + `,
	"key": ` + Golden_Key_String + `,
	"sig": "` + Golden_Sig + `"
 }`

// See also ExampleCanonicalHash
func Example_genCad() {
	digest, err := CanonicalHash([]byte(Golden_Pay), nil, ES256.Hash()) // compactify
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(digest)
	// Output:
	// LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto
}

func ExampleCozeKey_String() {
	var gk2 = Golden_Key // Make a copy
	fmt.Printf("%s\n", &gk2)
	// Output:
	// {"alg":"ES256","d":"bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVA","iat":1623132000,"kid":"Zami's Majuscule Key.","tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","x":"2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"}
}

//ExampleCozeKey_jsonUnmarshal tests unmarshalling a Coze key.
func ExampleCozeKey_jsonUnmarshal() {
	cozeKey := new(CozeKey)
	err := json.Unmarshal([]byte(Golden_Key_String), cozeKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", cozeKey)
	// Output:
	//{"alg":"ES256","d":"bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVA","iat":1623132000,"kid":"Zami's Majuscule Key.","tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","x":"2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"}
}

func ExampleCozeKey_jsonMarshal() {
	b, err := Marshal(Golden_Key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", string(b))
	// Output:
	//{"alg":"ES256","d":"bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVA","iat":1623132000,"kid":"Zami's Majuscule Key.","tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","x":"2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"}
}

func ExampleCozeKey_Thumbprint() {
	var gk2 = Golden_Key // Make a copy
	gk2.Tmb = []byte{}   // set it to nil to ensure recalculation.
	gk2.Thumbprint()
	h := gk2.Tmb
	fmt.Println(h)
	// Output:
	// cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk
}

func ExampleThumbprint() {
	h, err := Thumbprint(&Golden_Key)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(h)
	// Output:
	// cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk
}

func ExampleCozeKey_Sign() {
	cad := MustDecode(Golden_Cad)
	sig, err := Golden_Key.Sign(cad)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", Golden_Key.Verify(cad, sig))
	// Output: true
}

// ExampleCozeKey_Sign_empty shows signing an empty Coze is valid:
//
// {"pay":{},"sig":"9iesKUSV7L1-xz5yd3A94vCkKLmdOAnrcPXTU3_qeKSuk4RMG7Qz0KyubpATy0XA_fXrcdaxJTvXg6saaQQcVQ"}
//
// Where `alg`` and `key` are already implicitly known by the application.
func ExampleCozeKey_Sign_empty() {
	input := "{}"
	dig := Hash(Golden_Key.Alg.Hash(), []byte(input))

	sig, err := Golden_Key.Sign(dig)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%v\n", Golden_Key.Verify(dig, sig))
	// Output: true
}

// TODO test SignPay and SignPayJSON

func ExampleCozeKey_SignCoze() {
	cz := new(Coze)
	cz.Pay = []byte(Golden_Pay)

	err := Golden_Key.SignCoze(cz, nil)
	if err != nil {
		fmt.Println(err)
	}

	v, err := Golden_Key.VerifyCoze(cz)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v)

	// Output:
	// true
}

func ExampleCozeKey_Verify() {
	v := Golden_Key.Verify(MustDecode(Golden_Cad), MustDecode(Golden_Sig))
	fmt.Println(v)
	// Output: true
}

func ExampleCozeKey_VerifyCoze() {
	var cz = new(Coze)
	err := json.Unmarshal([]byte(Golden_Coze), cz)
	if err != nil {
		fmt.Println(err)
	}

	v, err := Golden_Key.VerifyCoze(cz)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(v)
	// Output: true
}

//  Tests valid on a good Coze key and a bad Coze key
func ExampleCozeKey_Valid() {
	fmt.Println(Golden_Key.Valid())
	fmt.Println(Golden_Bad_Key.Valid())

	// Output:
	// true
	// false
}

func ExampleNewKey_valid() {
	ck, err := NewKey(SEAlg(ES256))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ck.Valid())
	// Output:
	// true
}

func ExampleNewKey() {
	algs := []SigAlg{
		SigAlg(SHA256), // Invalid alg
		ES224,
		ES256,
		ES384,
		ES512,
		Ed25519,
	}

	for _, alg := range algs {
		cozeKey, err := NewKey(SEAlg(alg))
		if err != nil {
			fmt.Println(err)
			continue
		}

		if cozeKey.Valid() != true {
			fmt.Printf("signature was not verified.  Alg: %s\n", alg)
		}
	}
	fmt.Println("Done")
	// Output:
	// NewKey: unsupported alg: SHA-256
	// Done
}

func ExampleCozeKey_Correct() {

	// First key is "bad", all other keys should pass every test.
	var keys = []CozeKey{Golden_Bad_Key, Golden_Key}

	// Test new keys
	algs := []string{"ES224", "ES256", "ES384", "ES512", "Ed25519"}
	for _, alg := range algs {
		key, err := NewKey(ParseSEAlg(alg))
		if err != nil {
			fmt.Println(err)
		}
		keys = append(keys, *key)
	}

	for _, k := range keys {
		var gk2 = k // Make a copy

		// Key with with [alg,d,tmb,x]
		p1, _ := gk2.Correct()

		// A key with [alg,tmb,d]
		gk2.X = []byte{}
		p2, _ := gk2.Correct()

		// Key with [alg,d].
		gk2.Tmb = []byte{}
		p3, _ := gk2.Correct()

		// A key with [alg,x,d].
		gk2.X = k.X
		p4, _ := gk2.Correct()

		// A key with [alg,x,tmb]
		gk2.D = []byte{}
		gk2.Tmb = k.Tmb
		p5, _ := gk2.Correct()

		// Key with [alg,tmb]
		gk2.X = []byte{}
		p6, _ := gk2.Correct()

		fmt.Printf("%t, %t, %t, %t, %t, %t\n", p1, p2, p3, p4, p5, p6)
	}
	// Output:
	// false, false, true, false, true, true
	// true, true, true, true, true, true
	// true, true, true, true, true, true
	// true, true, true, true, true, true
	// true, true, true, true, true, true
	// true, true, true, true, true, true
	// true, true, true, true, true, true
}

func ExampleCozeKey_Revoke() {
	var gk2 = Golden_Key // Make a copy
	fmt.Println(gk2.IsRevoked())
	coze, err := gk2.Revoke("Posted my private key on github")
	if err != nil {
		fmt.Println(err)
	}
	v, err := gk2.VerifyCoze(coze)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n%+v\n", v, gk2.IsRevoked())
	// Output:
	// false
	// true
	// true
}

// BenchmarkNSV benchmarks several methods on a Coze Key. (NSV = New, Sign,
// Verify) It generates a new Coze Key, sign a message, and verifies the
// signature.
// go test -bench=.
// go test -bench=BenchmarkNSV -benchtime=30s
func BenchmarkNSV(b *testing.B) {
	var algs = []SigAlg{ES224, ES256, ES384, ES512, Ed25519}
	for j := 0; j < b.N; j++ {
		for i := 0; i < len(algs); i++ {
			ck, err := NewKey(SEAlg(algs[i]))
			if err != nil {
				b.Fatal("Could not generate Coze Key.")
			}

			if !ck.Valid() {
				b.Fatalf("The signature was invalid.  Alg: %s", ck.Alg)
			}
		}
	}
}
