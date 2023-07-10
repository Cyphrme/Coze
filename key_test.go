package coze

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"testing"
)

var GoldenKey = Key{
	Alg: SEAlg(ES256),
	Kid: "Zami's Majuscule Key.",
	Iat: 1623132000,
	X:   MustDecode("2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"),
	D:   MustDecode("bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVA"),
	Tmb: MustDecode("cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk"),
}

const GoldenKeyString = `{
	"alg":"ES256",
	"iat":1623132000,
	"kid":"Zami's Majuscule Key.",
	"d":"bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVA",
	"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
	"x":"2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"
}`

// The last byte in D was changed from 80 (b64ut "VA"), to 81 (b64ut "VE),
// making it invalid. Note that Base64 needs to be to "E", not "B" through "D"
// for it to be effective as "B" through "D" is non-canonical base64 encoding
// and may decode to the same byte string.
var GoldenKeyBadD = Key{
	Alg: SEAlg(ES256),
	Kid: "GoldenKeyBadD",
	Iat: 1623132000,
	X:   MustDecode("2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"),
	D:   MustDecode("bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVE"),
	Tmb: MustDecode("cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk"),
}

// The last byte in X was changed from 230 (b64ut "5g") to 231 (b64ut
// "5w") , making it invalid. See documentation on GoldenKeyBadD.
var GoldenKeyBadX = Key{
	Alg: SEAlg(ES256),
	Kid: "GoldenKeyBadX",
	Iat: 1623132000,
	X:   MustDecode("2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5w"),
	D:   MustDecode("bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVE"),
	Tmb: MustDecode("cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk"),
}

var (
	GoldenTmb = "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk"
	GoldenCad = "Ie3xL77AsiCcb4r0pbnZJqMcfSBqg5Lk0npNJyJ9BC4"
	GoldenCzd = "TnRe4DRuGJlw280u3pGhMDOIYM7ii7J8_PhNuSScsIU"
	GoldenSig = "Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w"
)

var GoldenPay = `{
	"msg": "Coze Rocks",
	"alg": "ES256",
	"iat": 1623132000,
	"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
	"typ": "cyphr.me/msg"
 }`

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

var GoldenCozeEmpty = json.RawMessage(`{
	"pay":{},
	"sig":"9iesKUSV7L1-xz5yd3A94vCkKLmdOAnrcPXTU3_qeKRRbHuy5EvMMFNRkW_sNLo-vvEPO9BmeUkcNh-ok18I_A"
}`)

var GoldenCozeNoAlg = json.RawMessage(`{
	"pay": {
			"msg": "Coze Rocks",
			"iat": 1623132000,
			"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
			"typ": "cyphr.me/msg"
	},
	"sig": "reOiKUO--OwgTNlYpKN60_gZARnW5X6PmQw4zWYbz2QryetRg_qS4KvwEVe1aiSAsWlkVA3MqYuaIM5ihY_8NQ"
}`)

var GoldenPayNoAlg = json.RawMessage(`{
			"msg": "Coze Rocks",
			"iat": 1623132000,
			"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
			"typ": "cyphr.me/msg"
	}`)

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
	err := gk2.Thumbprint()
	fmt.Println(gk2.Tmb, err)

	// Output:
	// cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk <nil>
}

func ExampleThumbprint() {
	fmt.Println(Thumbprint(&GoldenKey))
	// Output: cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk <nil>
}

func ExampleKey_Sign() {
	// Manual signing of empty Coze, `{"pay":{},"sig":"9iesKU..."}`, is a valid
	// Coze.  In this case, it would be better to use SignPayJSON.
	d, err := Hash(GoldenKey.Alg.Hash(), []byte("{}"))
	if err != nil {
		panic(err)
	}
	sig, err := GoldenKey.Sign(d)
	if err != nil {
		panic(err)
	}
	fmt.Println(GoldenKey.Verify(d, sig))

	// Signing a previously known cad.
	cad := MustDecode(GoldenCad)
	sig, err = GoldenKey.Sign(cad)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", GoldenKey.Verify(cad, sig))

	// Output:
	// true
	// true
}

// ExampleKey_SignPay demonstrates converting a custom data structure into a
// coze, signing it, and verifying the results.
func ExampleKey_SignPay() {
	customStruct := CustomStruct{
		Msg: "Coze Rocks",
	}

	pay := Pay{
		Alg:    SEAlg(ES256),
		Iat:    1623132000, // Static for demonstration.  Use time.Now().Unix().
		Tmb:    MustDecode("cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk"),
		Typ:    "cyphr.me/msg",
		Struct: customStruct,
	}

	coze, err := GoldenKey.SignPay(&pay)
	if err != nil {
		panic(err)
	}

	fmt.Println(GoldenKey.VerifyCoze(coze))

	// Output: true <nil>
}

func ExampleKey_SignPayJSON() {
	coze, err := GoldenKey.SignPayJSON(json.RawMessage(GoldenPay))
	if err != nil {
		panic(err)
	}
	fmt.Println(GoldenKey.VerifyCoze(coze))

	// Output: true <nil>
}

// ExampleKey_Sign_empty demonstrates signing of empty Coze,
// `{"pay":{},"sig":"9iesKU..."}`, is valid.
func ExampleKey_SignPayJSON_empty() {
	coze, err := GoldenKey.SignPayJSON([]byte("{}"))
	if err != nil {
		panic(err)
	}
	fmt.Println(GoldenKey.VerifyCoze(coze))

	// Output: true <nil>
}

// Example demonstrating the verification of the empty coze from the README.
func ExampleKey_Verify_empty() {
	cz := new(Coze)
	err := json.Unmarshal(GoldenCozeEmpty, cz)
	if err != nil {
		panic(err)
	}
	fmt.Println(GoldenKey.VerifyCoze(cz))

	// Output: true <nil>
}

// Example demonstrating that unmarshal generates `tmb` from x,
func ExampleKey_unmarshal() {
	var GoldenPukNoTmb = json.RawMessage(`{
		"alg":"ES256",
		"iat":1623132000,
		"kid":"Zami's Majuscule Key.",
		"x":"2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"
	}`)
	czk := new(Key)
	err := json.Unmarshal(GoldenPukNoTmb, czk)
	if err != nil {
		panic(err)
	}
	fmt.Println(czk)
	// Output:
	//  {"alg":"ES256","iat":1623132000,"kid":"Zami's Majuscule Key.","tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","x":"2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"}
}

// Example demonstrating that unmarshalling a `pay` that has duplicate field
// names results in an error.
func ExampleKey_UnmarshalJSON_duplicate() {
	k := &Key{}
	msg := []byte(`{"alg":"ES256","alg":"ES256"}`)
	err := json.Unmarshal(msg, k)
	fmt.Println(err)

	// Output:
	// Coze: JSON duplicate field "alg"
}

func ExampleKey_SignCoze() {
	cz := new(Coze)
	cz.Pay = json.RawMessage(GoldenPay)
	err := GoldenKey.SignCoze(cz)
	if err != nil {
		panic(err)
	}

	fmt.Println(GoldenKey.VerifyCoze(cz))

	// Output: true <nil>
}

func ExampleKey_Verify() {
	fmt.Println(GoldenKey.Verify(MustDecode(GoldenCad), MustDecode(GoldenSig)))

	// Output: true
}

func ExampleKey_VerifyCoze() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(GoldenCoze), cz)
	if err != nil {
		panic(err)
	}

	fmt.Println(GoldenKey.VerifyCoze(cz))

	// Output: true <nil>
}

// Tests valid on a good Coze key and a bad Coze key
func ExampleKey_Valid() {
	fmt.Println(GoldenKey.Valid(), GoldenKeyBadD.Valid())

	// Output: true false
}

func ExampleNewKey_valid() {
	ck, err := NewKey(SEAlg(ES256))
	if err != nil {
		panic(err)
	}
	fmt.Println(ck.Valid())

	// Output: true
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
	fmt.Println(NewKey(SEAlg(SHA256))) // Invalid signing alg, fails.

	// Output: <nil> NewKey: unsupported alg "SHA-256"
}

// ExampleKey_Correct demonstrates the expectations from Correct() when
// different key fields appear.  Note that some calls to Correct() pass on
// _invalid_ keys depending on given fields.
func ExampleKey_Correct() {
	// helper print function
	tf := func(err ...error) {
		for i, e := range err {
			if e != nil {
				fmt.Printf("%t", false)
			} else {
				fmt.Printf("%t", true)
			}
			if i < len(err)-1 {
				fmt.Printf(", ")
			}
		}
	}

	keys := []Key{GoldenKeyBadD, GoldenKeyBadX, GoldenKey}
	// Test new keys.  These keys should pass every test.
	algs := []string{"ES224", "ES256", "ES384", "ES512", "Ed25519"}
	for _, alg := range algs {
		key, err := NewKey(SEAlg(Parse(alg)))
		if err != nil {
			panic(err)
		}
		keys = append(keys, *key)
	}

	for _, k := range keys {
		gk2 := k // Make a copy

		// Key with with [alg,d,tmb,x]
		p1 := gk2.Correct()

		// A key with [alg,tmb,d]
		gk2 = k
		gk2.X = []byte{}
		p2 := gk2.Correct()

		// Key with [alg,d].
		gk2 = k
		gk2.X = []byte{}
		gk2.Tmb = []byte{}
		p3 := gk2.Correct()

		// A key with [alg,x,d].
		gk2 = k
		gk2.Tmb = []byte{}
		p4 := gk2.Correct()

		// A key with [alg,x,tmb]
		gk2 = k
		gk2.D = []byte{}
		p5 := gk2.Correct()

		// Key with [alg,tmb]
		gk2 = k
		gk2.D = []byte{}
		gk2.X = []byte{}
		p6 := gk2.Correct()

		tf(p1, p2, p3, p4, p5, p6)
		fmt.Printf("\n")
	}

	// Output:
	// false, false, true, false, true, true
	// false, false, true, false, false, true
	// true, true, true, true, true, true
	// true, true, true, true, true, true
	// true, true, true, true, true, true
	// true, true, true, true, true, true
	// true, true, true, true, true, true
	// true, true, true, true, true, true
}

// See also ExampleCanonicalHash.
func ExampleCanonicalHash_genCad() {
	fmt.Println(CanonicalHash([]byte(GoldenPay), nil, GoldenKey.Alg.Hash()))

	// Output: Ie3xL77AsiCcb4r0pbnZJqMcfSBqg5Lk0npNJyJ9BC4 <nil>
}

// BenchmarkNSV benchmarks several methods on a Coze key. (NSV = New, Sign,
// Verify) It generates a new Coze key, signs a message, and verifies the
// signature.
// go test -bench=.
// go test -bench=BenchmarkNSV -benchtime=30s
func BenchmarkNSV(b *testing.B) {
	algs := []SigAlg{ES224, ES256, ES384, ES512, Ed25519}
	for j := 0; j < b.N; j++ {
		for _, alg := range algs {
			ck, err := NewKey(SEAlg(alg))
			if err != nil {
				b.Fatal("Could not generate Coze Key.")
			}

			if !ck.Valid() {
				b.Fatalf("The signature was invalid.  Alg: %s", ck.Alg)
			}
		}
	}
}

func ExampleKey_IsRevoked() {
	gk2 := GoldenKey // Make a copy
	fmt.Println(gk2.IsRevoked())
	coze, err := gk2.Revoke()
	if err != nil {
		panic(err)
	}

	pay := new(Pay)
	err = pay.UnmarshalJSON(coze.Pay)
	if err != nil {
		panic(err)
	}
	// Both the revoke coze and the key should be interpreted as revoked.
	fmt.Println(pay.IsRevoke())
	fmt.Println(gk2.IsRevoked())

	// Output:
	// false
	// true
	// true
}

// Example_ECDSAToLowSSig demonstrates converting non-coze compliant high S
// signatures to the canonicalized, coze compliant low-S form.
func ExampleECDSAToLowSSig() {
	highSCozies := []string{
		`{"pay":{},"sig":"9iesKUSV7L1-xz5yd3A94vCkKLmdOAnrcPXTU3_qeKSuk4RMG7Qz0KyubpATy0XA_fXrcdaxJTvXg6saaQQcVQ"}`,
		`{"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"mVw8N6ZncWcObVGvnwUMRIC6m2fbX3Sr1LlHMbj_tZ3ji1rNL-00pVaB12_fmlK3d_BVDipNQUsaRyIlGJudtg"}`,
		`{"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"cn6KNl4VQlk5MzmhYFVyyJoTOU57O5Bq-8r-yXXR6Ojfs0-6LFGd8j1Y6wiJAQrGpWj_RptsiEg49v95FsVWMQ"}`,
		`{"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"9KvWfOSIZUjW8Ie0jbdVdu9UlIP4TT4MXz3YyNW3fCTWXHnO1MPROwcXvfNZN_icOvMAK3vfsr2w-CeBozS81w"}`,
	}

	for _, s := range highSCozies {
		cz := new(Coze)
		err := json.Unmarshal([]byte(s), cz)
		if err != nil {
			panic(err)
		}
		v, _ := GoldenKey.VerifyCoze(cz)
		if v {
			panic("High S coze should not validate.")
		}

		err = ECDSAToLowSSig(&GoldenKey, cz)
		if err != nil {
			panic(err)
		}

		v, _ = GoldenKey.VerifyCoze(cz)
		if !v {
			panic("low-S coze should validate.")
		}

		fmt.Printf("%s\n", cz)
	}

	// Output:
	// {"pay":{},"sig":"9iesKUSV7L1-xz5yd3A94vCkKLmdOAnrcPXTU3_qeKRRbHuy5EvMMFNRkW_sNLo-vvEPO9BmeUkcNh-ok18I_A"}
	// {"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"mVw8N6ZncWcObVGvnwUMRIC6m2fbX3Sr1LlHMbj_tZ0cdKUx0BLLW6l-KJAgZa1IRPaln3zKXTnZcqid48eHmw"}
	// {"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"cn6KNl4VQlk5MzmhYFVyyJoTOU57O5Bq-8r-yXXR6OggTLBE065iDsKnFPd2_vU5F337ZwurFjy6wstJ5Z3PIA"}
	// {"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"9KvWfOSIZUjW8Ie0jbdVdu9UlIP4TT4MXz3YyNW3fCQpo4YwKzwuxfjoQgymyAdjgfP6gis368dCwaNBWS5oeg"}
}

// Test_lowS tests to make sure generated ECDSA keys are low-s and not high-s.
func Test_lowS(t *testing.T) {
	d, err := Hash(SHA512, []byte("7AtyaCHO2BAG06z0W1tOQlZFWbhxGgqej4k9-HWP3DE-zshRbrE-69DIfgY704_FDYez7h_rEI1WQVKhv5Hd5Q"))
	if err != nil {
		panic(err)
	}

	// Tests runs 512 times (4 * 128)
	algs := []SigAlg{ES224, ES256, ES384, ES512}
	for i := 0; i < 128; i++ {
		for _, alg := range algs {
			ck, err := NewKey(SEAlg(alg))
			if err != nil {
				panic(err)
			}
			sig, err := ck.Sign(d)
			if err != nil {
				panic(err)
			}

			size := ck.Alg.SigAlg().SigSize() / 2
			s := big.NewInt(0).SetBytes(sig[size:])
			ls, err := IsLowS(ck, s)
			if err != nil {
				panic(err)
			}
			if !ls {
				t.Errorf("Signature was not low-s")
			}
		}
	}
}

// Test_ed25519Malleability demonstrates that the Go Ed25519 implementation
// is not malleable.  The malleable form was generated using
// https://slowli.github.io/ed25519-quirks/malleability.
func Test_ed25519Malleability(t *testing.T) {
	msg := []byte("Hello, world!")
	// Public key b64ut: 6ySOZK7GdGzY-AJvpRwNyWOD4RtWGB4rDJCD0MLCG-M
	// Public key ub64p: 6ySOZK7GdGzY+AJvpRwNyWOD4RtWGB4rDJCD0MLCG+M=
	// Public key Hex: 85B3AC41C8A2F1D1CD1E684169E67253F26D1862605ACB615D2D4E0CC44941AA
	//
	// seed || pub key form:
	// Hex: 85B3AC41C8A2F1D1CD1E684169E67253F26D1862605ACB615D2D4E0CC44941AAEB248E64AEC6746CD8F8026FA51C0DC96383E11B56181E2B0C9083D0C2C21BE3
	// ub64p: hbOsQcii8dHNHmhBaeZyU/JtGGJgWsthXS1ODMRJQarrJI5krsZ0bNj4Am+lHA3JY4PhG1YYHisMkIPQwsIb4w==
	//
	// Seed:
	// ub64p: hbOsQcii8dHNHmhBaeZyU/JtGGJgWsthXS1ODMRJQao=
	pri := ed25519.NewKeyFromSeed(MustDecode("hbOsQcii8dHNHmhBaeZyU_JtGGJgWsthXS1ODMRJQao"))
	p := pri.Public()
	pub, ok := p.(ed25519.PublicKey)
	if !ok {
		return
	}

	// ub64p: WA8oFP3rnGa/Fbcei89ztetTEJ921iOLgPlUbww2ZbyHq3pYD/ZN5mpUC7iBXMJdzM7zV1nbi0TSzbFpAk4ACA==
	valid := ed25519.Verify(pub, msg, MustDecode("WA8oFP3rnGa_Fbcei89ztetTEJ921iOLgPlUbww2ZbyHq3pYD_ZN5mpUC7iBXMJdzM7zV1nbi0TSzbFpAk4ACA"))
	if !valid {
		t.Errorf("Example valid b64ut canonical (non-malleable) signature is not valid.")
	}

	// Alternative, malleable form.
	valid = ed25519.Verify(pub, msg, MustDecode("WA8oFP3rnGa_Fbcei89ztetTEJ921iOLgPlUbww2Zbx0f3C1KVlgPkHxAltgVqFyzM7zV1nbi0TSzbFpAk4AGA"))
	if valid {
		t.Errorf("Example invalid b64ut non-canonical (non-malleable) signature is not valid.")
	}
}

// Test_curveOrder tests if the curve order values are correct
func Test_curveOrder(t *testing.T) {
	algs := []SigAlg{
		ES224,
		ES256,
		ES384,
		ES512,
	}

	s := ""
	for _, a := range algs {
		hexSize := Alg(a).Params().XSize
		s += fmt.Sprintf("%0"+strconv.Itoa(hexSize)+"X\n", curveOrders[a])
	}
	s += "\n"
	for _, a := range algs {
		hexSize := Alg(a).Params().XSize
		s += fmt.Sprintf("%0"+strconv.Itoa(hexSize)+"X\n", curveHalfOrders[a])
	}

	golden := `FFFFFFFFFFFFFFFFFFFFFFFFFFFF16A2E0B8F03E13DD29455C5C2A3D
FFFFFFFF00000000FFFFFFFFFFFFFFFFBCE6FAADA7179E84F3B9CAC2FC632551
FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFC7634D81F4372DDF581A0DB248B0A77AECEC196ACCC52973
01FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFA51868783BF2F966B7FCC0148F709A5D03BB5C9B8899C47AEBB6FB71E91386409

7FFFFFFFFFFFFFFFFFFFFFFFFFFF8B51705C781F09EE94A2AE2E151E
7FFFFFFF800000007FFFFFFFFFFFFFFFDE737D56D38BCF4279DCE5617E3192A8
7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFE3B1A6C0FA1B96EFAC0D06D9245853BD76760CB5666294B9
00FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFD28C343C1DF97CB35BFE600A47B84D2E81DDAE4DC44CE23D75DB7DB8F489C3204
`
	if s != golden {
		t.Errorf("incorrect curve order values")
	}
}

func TestKeyTmb_nilX(t *testing.T) {
	kb := []byte(`{
		"alg":"ES256",
		"iat":1647357960,
		"kid":"Cyphr.me Dev Test Key 2",
		"tmb":"e02u-nce-Wdc_xH4-7WRp-4Fr6pKn2oY_8SX3wdT41U"
	}`)
	key := new(Key)
	err := json.Unmarshal(kb, key)
	if err != nil {
		t.Fatal(err)
	}

	// Test nil x.
	err = key.Thumbprint() // Should error since x is nil.
	if err == nil {
		t.Fatal("key.Thumbprint() must error when x is nil.")
	}
	if len(key.Tmb) != 0 {
		t.Fatal("key.Thumbprint() must set thumbprint to nil on error.")
	}

	// Test malformed x.
	key.X = MustDecode("e02u")   // incorrect x for alg ES256 (known by length)
	key.Tmb = MustDecode("e02u") // incorrect.  Should be set to nil on Thumbprint error.
	err = key.Thumbprint()       // Should error
	if err == nil {
		t.Fatal("key.Thumbprint() must error when x is incorrect length.")
	}
}
