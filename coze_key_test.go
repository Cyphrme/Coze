package coze

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cyphrme/coze/enum"
)

var testDigest = []byte{112, 184, 252, 190, 198, 45, 48, 28, 24, 147, 58, 5, 85, 145, 193, 102, 142, 146, 52, 191, 48, 73, 208, 136, 140, 34, 128, 193, 115, 110, 132, 233}

var Golden_Key = CozeKey{
	Alg: enum.SEAlg(enum.ES256),
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
var Golden_Bad_Key = &CozeKey{
	Alg: enum.SEAlg(enum.ES256),
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

var Golden_Cad = "LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto"
var Golden_Sig = "ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"

var Golden_Cy = `{
	"pay":` + Golden_Pay + `,
	"sig": "` + Golden_Sig + `"
 }`

var Golden_Coze = `{
	"coze":` + Golden_Cy + `
}`

var Golden_Cy_W_Key = `{
	"pay": ` + Golden_Pay + `,
	"key": ` + Golden_Key_String + `
	,
	"sig": "` + Golden_Sig + `"
 }`

// See also ExampleCanonHash
func Example_genCad() {
	digest, err := CanonHash([]byte(Golden_Pay), nil, enum.ES256.Hash()) // compactify
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
	cozekey := new(CozeKey)
	err := json.Unmarshal([]byte(Golden_Key_String), cozekey)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", cozekey)
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
	Golden_Key.Tmb = []byte{} // set it to nil to ensure recalculation.
	Golden_Key.Thumbprint()
	h := Golden_Key.Tmb
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

	valid := Golden_Key.Verify(cad, sig)

	fmt.Printf("%v\n", valid)
	// Output: true
}

func ExampleCozeKey_SignCy() {
	cy := new(Cy)
	cy.Pay = []byte(Golden_Pay)

	err := Golden_Key.SignCy(cy, nil)
	if err != nil {
		fmt.Println(err)
	}

	v, err := Golden_Key.VerifyCy(cy)
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

func ExampleCozeKey_VerifyCy() {
	var cy = new(Cy)
	err := json.Unmarshal([]byte(Golden_Cy), cy)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Printf("%s", cy)
	v, err := Golden_Key.VerifyCy(cy)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(v)
	// Output: true
}

//  Tests valid on a good coze key and a bad coze key
func ExampleCozeKey_Valid() {
	valid := Golden_Key.Valid()
	fmt.Println(valid)

	valid = Golden_Bad_Key.Valid()
	fmt.Println(valid)

	// Output:
	// true
	// false
}

func ExampleNewKey_es256_valid() {
	ck, err := NewKey(enum.SEAlg(enum.ES256))
	if err != nil {
		fmt.Println(err)
	}

	v := ck.Valid()
	if v != true {
		fmt.Println("Invalid key")
	}

	fmt.Printf("%s %v\n", ck.Alg, v)

	// Output:
	// ES256 true
}

func TestGenKeys(t *testing.T) {
	algs := []enum.SigAlg{
		enum.ES224,
		enum.ES256,
		enum.ES384,
		enum.ES512,
		//enum.Ed25519 //TODO
	}

	for _, alg := range algs {
		cozeKey, err := NewKey(enum.SEAlg(alg))
		if err != nil {
			t.Fatal(err)
		}

		dig := Hash(enum.SEAlg(alg).Hash(), []byte("Test Message"))
		s, err := cozeKey.Sign(dig)
		if err != nil {
			t.Fatal(err)
		}
		v := cozeKey.Verify(dig, s)
		if v != true {
			t.Fatal("signature was not verified.  Alg: ", alg)
		}
	}
}

// BenchmarkNSV benchmarks several methods on a Coze Key. (NSV = New, Sign,
// Verify) It generates a new Coze Key, sign a message, and verifies the
// signature.
// go test -bench=.
// go test -bench=BenchmarkNSV -benchtime=30s
func BenchmarkNSV(b *testing.B) {
	// TODO Ed25519 Support:
	var algs = []enum.SigAlg{enum.ES224, enum.ES256, enum.ES384, enum.ES512, enum.Ed25519}

	for j := 0; j < b.N; j++ {
		for i := 0; i < len(algs); i++ {
			ck, err := NewKey(enum.SEAlg(algs[i]))
			if err != nil {
				b.Fatal("Could not generate Coze Key.")
			}

			sig, err := ck.Sign(testDigest)
			if err != nil {
				b.Fatal(err)
			}

			v := ck.Verify(testDigest, sig)
			if !v {
				b.Fatalf("The signature was invalid.  Alg: %s", ck.Alg)
			}
		}
	}
}

func Example_eS256_nsv() {
	ck, err := NewKey(enum.SEAlg(enum.ES256))
	// log.Printf("Alg: %+v, Key: %+v\n", ck.Alg, ck)
	if err != nil {
		panic("Could not generate Coze Key.")
	}

	sig, err := ck.Sign(testDigest)
	if err != nil {
		panic(err)
	}

	v := ck.Verify(testDigest, sig)

	fmt.Println(v)

	// Output:
	// true
}
