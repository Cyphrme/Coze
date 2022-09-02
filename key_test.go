package coze

import (
	"encoding/json"
	"fmt"
	"testing"
)

var GoldenKey = Key{
	Alg: SEAlg(ES256),
	Kid: "Zami's Majuscule Key.",
	Iat: 1623132000,
	X:   []byte{218, 116, 206, 104, 85, 102, 217, 2, 241, 153, 67, 191, 74, 56, 50, 177, 197, 71, 6, 219, 199, 17, 250, 54, 174, 174, 185, 50, 248, 13, 70, 51, 145, 162, 58, 183, 244, 118, 170, 246, 181, 205, 198, 245, 241, 193, 182, 191, 94, 61, 5, 230, 246, 98, 108, 148, 119, 138, 192, 93, 57, 102, 232, 230},
	D:   []byte{108, 219, 45, 131, 143, 199, 222, 109, 210, 149, 19, 174, 127, 4, 82, 18, 8, 155, 46, 176, 110, 70, 175, 117, 215, 131, 175, 117, 170, 92, 165, 80},
	Tmb: []byte{112, 184, 252, 190, 198, 45, 48, 28, 24, 147, 58, 5, 85, 145, 193, 102, 142, 146, 52, 191, 48, 73, 208, 136, 140, 34, 128, 193, 115, 110, 132, 233},
}

const GoldenKeyString = `{
	"alg":"ES256",
	"iat":1623132000,
	"kid":"Zami's Majuscule Key.",
	"d":"bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVA",
	"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
	"x":"2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"
}`

// The very last byte in D was changed from 80, to 81, making it invalid.
// Base64 needs to be to "E", not B-D for it to be effective: bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVE
var GoldenKeyBad = Key{
	Alg: SEAlg(ES256),
	Kid: "Zami's Majuscule Key.",
	Iat: 1623132000,
	X:   []byte{218, 116, 206, 104, 85, 102, 217, 2, 241, 153, 67, 191, 74, 56, 50, 177, 197, 71, 6, 219, 199, 17, 250, 54, 174, 174, 185, 50, 248, 13, 70, 51, 145, 162, 58, 183, 244, 118, 170, 246, 181, 205, 198, 245, 241, 193, 182, 191, 94, 61, 5, 230, 246, 98, 108, 148, 119, 138, 192, 93, 57, 102, 232, 230},
	D:   []byte{108, 219, 45, 131, 143, 199, 222, 109, 210, 149, 19, 174, 127, 4, 82, 18, 8, 155, 46, 176, 110, 70, 175, 117, 215, 131, 175, 117, 170, 92, 165, 81},
	Tmb: []byte{112, 184, 252, 190, 198, 45, 48, 28, 24, 147, 58, 5, 85, 145, 193, 102, 142, 146, 52, 191, 48, 73, 208, 136, 140, 34, 128, 193, 115, 110, 132, 233},
}

var GoldenPay = `{
	"msg": "Coze Rocks",
	"alg": "ES256",
	"iat": 1627518000,
	"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
	"typ": "cyphr.me/msg"
 }`

var (
	GoldenTmb = "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk"
	GoldenCad = "LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto"
	GoldenCzd = "d0ygwQCGzuxqgUq1KsuAtJ8IBu0mkgAcKpUJzuX075M"
	GoldenSig = "ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"
)

var GoldenCoze = `{
	"pay":` + GoldenPay + `,
	"sig": "` + GoldenSig + `"
 }`

// Encapsulated coze.
var GoldenCozeE = `{
	"coze":` + GoldenCoze + `
}`

var GoldenCozeWKey = `{
	"pay": ` + GoldenPay + `,
	"key": ` + GoldenKeyString + `,
	"sig": "` + GoldenSig + `"
 }`

// CustomStruct is for examples demonstrating Pay/Coze with custom structs.
type CustomStruct struct {
	Msg string `json:"msg,omitempty"`
}

func ExampleKey_String() {
	fmt.Printf("%s\n", GoldenKey)
	// Output:
	// {"alg":"ES256","d":"bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVA","iat":1623132000,"kid":"Zami's Majuscule Key.","tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","x":"2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"}
}

// ExampleKey_jsonUnmarshal tests unmarshalling a Coze key.
func ExampleKey_jsonUnmarshal() {
	Key := new(Key)
	err := json.Unmarshal([]byte(GoldenKeyString), Key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", Key)
	// Output:
	//{"alg":"ES256","d":"bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVA","iat":1623132000,"kid":"Zami's Majuscule Key.","tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","x":"2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"}
}

func ExampleKey_jsonMarshal() {
	b, err := Marshal(GoldenKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(b))
	// Output:
	//{"alg":"ES256","d":"bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVA","iat":1623132000,"kid":"Zami's Majuscule Key.","tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","x":"2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"}
}

func ExampleKey_Thumbprint() {
	gk2 := GoldenKey   // Make a copy.
	gk2.Tmb = []byte{} // Set to empty to ensure recalculation.
	gk2.Thumbprint()
	fmt.Println(gk2.Tmb)
	// Output:
	// cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk
}

func ExampleThumbprint() {
	h, err := Thumbprint(&GoldenKey)
	if err != nil {
		panic(err)
	}
	fmt.Println(h)
	// Output:
	// cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk
}

func ExampleKey_Sign() {
	cad := MustDecode(GoldenCad)
	sig, err := GoldenKey.Sign(cad)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", GoldenKey.Verify(cad, sig))
	// Output: true
}

// ExampleKey_Sign_empty demonstrates that signing an empty Coze,
// `{"pay":{},"sig":"9iesKU..."}`, is valid.
func ExampleKey_Sign_empty() {
	d, err := Hash(GoldenKey.Alg.Hash(), []byte("{}"))
	if err != nil {
		panic(err)
	}
	sig, err := GoldenKey.Sign(d)
	if err != nil {
		panic(err)
	}

	fmt.Println(GoldenKey.Verify(d, sig))
	// Output: true
}

// ExampleKey_SignPay demonstrates converting a custom data structure into a
// coze, signing it, and verifying the results.
func ExampleKey_SignPay() {
	customStruct := CustomStruct{
		Msg: "Coze Rocks",
	}

	pay := Pay{
		Alg:    SEAlg(ES256),
		Iat:    1627518000, // Static for demonstration.  Use time.Time.Unix().
		Tmb:    MustDecode("cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk"),
		Typ:    "cyphr.me/msg",
		Struct: customStruct,
	}

	coze, err := GoldenKey.SignPay(&pay)
	if err != nil {
		panic(err)
	}

	v, err := GoldenKey.VerifyCoze(coze)
	if err != nil {
		panic(err)
	}
	fmt.Println(v)

	// Output:
	// true
}

func ExampleKey_SignPayJSON() {
	coze := new(Coze)
	coze.Pay = []byte(GoldenPay)

	var err error
	coze.Sig, err = GoldenKey.SignPayJSON(coze.Pay)
	if err != nil {
		panic(err)
	}

	v, err := GoldenKey.VerifyCoze(coze)
	if err != nil {
		panic(err)
	}
	fmt.Println(v)

	// Output:
	// true
}

func ExampleKey_SignCoze() {
	cz := new(Coze)
	cz.Pay = []byte(GoldenPay)

	err := GoldenKey.SignCoze(cz)
	if err != nil {
		panic(err)
	}

	v, err := GoldenKey.VerifyCoze(cz)
	if err != nil {
		panic(err)
	}
	fmt.Println(v)

	// Output:
	// true
}

func ExampleKey_Verify() {
	v := GoldenKey.Verify(MustDecode(GoldenCad), MustDecode(GoldenSig))
	fmt.Println(v)
	// Output: true
}

func ExampleKey_VerifyCoze() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(GoldenCoze), cz)
	if err != nil {
		panic(err)
	}

	v, err := GoldenKey.VerifyCoze(cz)
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: true
}

//  Tests valid on a good Coze key and a bad Coze key
func ExampleKey_Valid() {
	fmt.Println(GoldenKey.Valid())
	fmt.Println(GoldenKeyBad.Valid())

	// Output:
	// true
	// false
}

func ExampleNewKey_valid() {
	ck, err := NewKey(SEAlg(ES256))
	if err != nil {
		panic(err)
	}
	fmt.Println(ck.Valid())
	// Output:
	// true
}

func ExampleNewKey() {
	algs := []SigAlg{
		ES224,
		ES256,
		ES384,
		ES512,
		Ed25519,
	}

	for _, alg := range algs {
		Key, err := NewKey(SEAlg(alg))
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s, %t\n", Key.Alg, Key.Valid())
	}

	// Output:
	// ES224, true
	// ES256, true
	// ES384, true
	// ES512, true
	// Ed25519, true
}

func ExampleNewKey_bad() {
	_, err := NewKey(SEAlg(SHA256)) // Invalid signing alg, fails.
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// NewKey: unsupported alg: SHA-256
}

func ExampleKey_Correct() {
	// Note that some calls to Correct() pass on **invalid** keys depending on
	// given fields. Second static key is valid so all field combinations pass.
	keys := []Key{GoldenKeyBad, GoldenKey}

	// Test new keys.  These keys should pass every test.
	algs := []string{"ES224", "ES256", "ES384", "ES512", "Ed25519"}
	for _, alg := range algs {
		key, err := NewKey(ParseSEAlg(alg))
		if err != nil {
			panic(err)
		}
		keys = append(keys, *key)
	}

	for _, k := range keys {
		gk2 := k // Make a copy

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

// See also ExampleCanonicalHash.
func ExampleCanonicalHash_genCad() {
	digest, err := CanonicalHash([]byte(GoldenPay), nil, GoldenKey.Alg.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Println(digest)
	// Output:
	// LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto
}

// BenchmarkNSV benchmarks several methods on a Coze key. (NSV = New, Sign,
// Verify) It generates a new Coze key, signs a message, and verifies the
// signature.
// go test -bench=.
// go test -bench=BenchmarkNSV -benchtime=30s
func BenchmarkNSV(b *testing.B) {
	algs := []SigAlg{ES224, ES256, ES384, ES512, Ed25519}
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
