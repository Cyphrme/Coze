package coze

import (
	"encoding/json"
	"fmt"
	"testing"

	ce "github.com/zamicol/cyphrme/src/coze/enum"
)

var TestMsg = []byte("Coze Key Test Message.")

var Golden_Key = CozeKey{
	Alg: ce.SEAlg(ce.ES256),
	Kid: "Zami's Majuscule Key.",
	Iat: 1623132000,
	X:   []byte{218, 116, 206, 104, 85, 102, 217, 2, 241, 153, 67, 191, 74, 56, 50, 177, 197, 71, 6, 219, 199, 17, 250, 54, 174, 174, 185, 50, 248, 13, 70, 51},
	Y:   []byte{145, 162, 58, 183, 244, 118, 170, 246, 181, 205, 198, 245, 241, 193, 182, 191, 94, 61, 5, 230, 246, 98, 108, 148, 119, 138, 192, 93, 57, 102, 232, 230},
	D:   []byte{108, 219, 45, 131, 143, 199, 222, 109, 210, 149, 19, 174, 127, 4, 82, 18, 8, 155, 46, 176, 110, 70, 175, 117, 215, 131, 175, 117, 170, 92, 165, 80},
	Tmb: []byte{1, 72, 244, 205, 144, 147, 201, 203, 227, 232, 191, 120, 211, 230, 201, 184, 36, 241, 29, 210, 242, 158, 43, 26, 99, 13, 209, 206, 30, 23, 108, 221},
}

// Jared's
var Golden_Key2 = &CozeKey{
	Alg: ce.SEAlg(ce.ES256),
	Kid: "Jared's Key.",
	Iat: 1623132000,
	D:   []byte{138, 166, 123, 144, 57, 109, 48, 60, 124, 218, 25, 55, 117, 228, 71, 156, 231, 49, 164, 6, 221, 47, 5, 239, 108, 127, 88, 125, 66, 109, 203, 30},
	X:   []byte{23, 235, 135, 95, 245, 37, 103, 157, 216, 10, 170, 131, 166, 158, 132, 168, 179, 96, 102, 64, 203, 39, 160, 215, 175, 167, 140, 71, 21, 95, 123, 82},
	Y:   []byte{254, 159, 7, 249, 143, 34, 222, 244, 176, 69, 209, 107, 250, 155, 36, 60, 32, 179, 229, 63, 70, 75, 60, 47, 17, 119, 103, 238, 104, 96, 239, 239},
	Tmb: []byte{184, 11, 46, 86, 35, 158, 230, 121, 1, 112, 133, 152, 244, 80, 224, 4, 54, 222, 172, 202, 23, 240, 182, 238, 117, 83, 147, 124, 85, 153, 149, 42},
}

const Golden_Key_String = `{
	"alg":"ES256",
	"d":"6CDB2D838FC7DE6DD29513AE7F045212089B2EB06E46AF75D783AF75AA5CA550",
	"iat":1623132000,
	"tmb":"0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD",
	"kid":"Zami's Majuscule Key.",
	"x":"DA74CE685566D902F19943BF4A3832B1C54706DBC711FA36AEAEB932F80D4633",
	"y":"91A23AB7F476AAF6B5CDC6F5F1C1B6BF5E3D05E6F6626C94778AC05D3966E8E6"
}`

const Golden_Key2_String = `{
	"use":"sig",
	"alg":"ES256",
	"d":"8AA67B90396D303C7CDA193775E4479CE731A406DD2F05EF6C7F587D426DCB1E",
	"x":"17EB875FF525679DD80AAA83A69E84A8B3606640CB27A0D7AFA78C47155F7B52",
	"y":"FE9F07F98F22DEF4B045D16BFA9B243C20B3E53F464B3C2F117767EE6860EFEF",
	"tmb":"B80B2E56239EE67901708598F450E00436DEACCA17F0B6EE7553937C5599952A",
	"iat":1623132000,
	"kid":"Jared's Key",
	"cyphrme_added":false,
	"cyphrme_accountid":"B80B2E56239EE67901708598F450E00436DEACCA17F0B6EE7553937C5599952A"
}`

// The very last byte in D was purposely changed from 80, to 81, which should make it invalid.
var Golden_Bad_Key = &CozeKey{
	Alg: ce.SEAlg(ce.ES256),
	Kid: "Zami's Majuscule Key.",
	Iat: 1623132000,
	X:   []byte{218, 116, 206, 104, 85, 102, 217, 2, 241, 153, 67, 191, 74, 56, 50, 177, 197, 71, 6, 219, 199, 17, 250, 54, 174, 174, 185, 50, 248, 13, 70, 51},
	Y:   []byte{145, 162, 58, 183, 244, 118, 170, 246, 181, 205, 198, 245, 241, 193, 182, 191, 94, 61, 5, 230, 246, 98, 108, 148, 119, 138, 192, 93, 57, 102, 232, 230},
	D:   []byte{108, 219, 45, 131, 143, 199, 222, 109, 210, 149, 19, 174, 127, 4, 82, 18, 8, 155, 46, 176, 110, 70, 175, 117, 215, 131, 175, 117, 170, 92, 165, 81},
	Tmb: []byte{1, 72, 244, 205, 144, 147, 201, 203, 227, 232, 191, 120, 211, 230, 201, 184, 36, 241, 29, 210, 242, 158, 43, 26, 99, 13, 209, 206, 30, 23, 108, 221},
}

//ExampleCyUnmarshal tests unmarshalling a `cy`.
func ExampleCozeKey_jsonUnmarshal() {
	cozekey := new(CozeKey)
	err := json.Unmarshal([]byte(Golden_Key_String), cozekey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", cozekey)
	// Output:
	// {"alg":"ES256","kid":"Zami's Majuscule Key.","iat":1623132000,"tmb":"0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD","d":"6CDB2D838FC7DE6DD29513AE7F045212089B2EB06E46AF75D783AF75AA5CA550","x":"DA74CE685566D902F19943BF4A3832B1C54706DBC711FA36AEAEB932F80D4633","y":"91A23AB7F476AAF6B5CDC6F5F1C1B6BF5E3D05E6F6626C94778AC05D3966E8E6"}
}

func ExampleCozeKey_jsonMarshal() {
	b, err := Marshal(Golden_Key)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s\n", string(b))
	// Output: {"alg":"ES256","kid":"Zami's Majuscule Key.","iat":1623132000,"tmb":"0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD","d":"6CDB2D838FC7DE6DD29513AE7F045212089B2EB06E46AF75D783AF75AA5CA550","x":"DA74CE685566D902F19943BF4A3832B1C54706DBC711FA36AEAEB932F80D4633","y":"91A23AB7F476AAF6B5CDC6F5F1C1B6BF5E3D05E6F6626C94778AC05D3966E8E6"}
}

func ExampleCozeKey_String() {
	var gk2 = Golden_Key // Make a copy

	fmt.Printf("%s\n", &gk2)
	// Output:
	// {"alg":"ES256","kid":"Zami's Majuscule Key.","iat":1623132000,"tmb":"0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD","d":"6CDB2D838FC7DE6DD29513AE7F045212089B2EB06E46AF75D783AF75AA5CA550","x":"DA74CE685566D902F19943BF4A3832B1C54706DBC711FA36AEAEB932F80D4633","y":"91A23AB7F476AAF6B5CDC6F5F1C1B6BF5E3D05E6F6626C94778AC05D3966E8E6"}
}

func ExampleCozeKey_Thumbprint() {
	Golden_Key.Tmb = []byte{} // set it to nil to ensure recalc
	Golden_Key.Thumbprint()

	h := Golden_Key.Tmb

	fmt.Println(h)
	// Output:
	// 0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD
}

func ExampleThumbprint() {
	h, err := Thumbprint(&Golden_Key)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(h)
	// Output:
	// 0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD
}

// This test will take a good coze key, and a bad coze key, and sign a message with them.
// The good coze key's signature should be verified, and the bad key's sig
// should be invalid.
func ExampleCozeKey_Valid() {
	valid := Golden_Key.Valid()
	if !valid {
		fmt.Println("Coze Key is invalid")
	}
	fmt.Println(valid)

	valid = Golden_Bad_Key.Valid()
	if valid {
		fmt.Println("Invalid Coze key is valid")
	}

	// Output: true
}

func ExampleCozeKey_SignRaw() {
	sig, err := Golden_Key.SignRaw(TestMsg)
	if err != nil {
		fmt.Println(err)
	}

	valid, err := Golden_Key.VerifyRaw(TestMsg, sig)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%v\n", valid)
	// Output: true
}

func TestCozeKey_SignRaw(t *testing.T) {
	sig, err := Golden_Key.SignRaw(TestMsg)
	if err != nil {
		fmt.Println(err)
	}

	valid, err := Golden_Key.VerifyRaw(TestMsg, sig)

	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("Signature is not valid.")
	}
}

func ExampleCozeKey_SignCy() {
	cy := new(Cy)
	b := []byte(Golden_Head_String)
	cy.Head = b

	err := Golden_Key.SignCy(cy, nil)
	if err != nil {
		fmt.Println(err)
	}
	cy.Sig = nil // nil sig for marshal since for ecdsa it's non-deterministic.
	b, err = Marshal(cy)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s", b)
	// Output:
	//{"head":{"alg":"ES256","iat":1623132000,"tmb":"0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD","typ":"cyphr.me"}}
}

func ExampleNewKey_es256_valid() {
	ck, err := NewKey(ce.SEAlg(ce.ES256))
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

func ExampleNewKey_ed25519_sign() {
	ck, err := NewKey(ce.SEAlg(ce.Ed25519))
	if err != nil {
		fmt.Println(err)
	}

	s, err := ck.SignRaw(TestMsg)
	if err != nil {
		fmt.Println(err)
	}
	v, err := ck.VerifyRaw(TestMsg, s)
	if err != nil {
		fmt.Println(err)
	}
	if v != true {
		fmt.Println("signature was not verified")
	}
	fmt.Printf("%s %v\n", ck.Alg, v)
	// Output:
	// Ed25519 true
}

// ExampleCozeKey_Verify verifies a signature.
func ExampleCozeKey_VerifyRaw() {
	// String was signed by Golden_Key
	sb, err := HexDecodeString("2A3E94B9501165FC70CCD9CEABFAD985B1FE71F29E7EABC8B09D2B13A10C362BD09D93FD473E5599960D4607FFB2C8F99ABFE8805210EBA604705A5A6F9AD0F4")
	if err != nil {
		fmt.Println(err)
	}

	valid, err := Golden_Key.VerifyRaw(TestMsg, sb)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(valid)
	// Output: true
}

func TestGenKeys(t *testing.T) {
	algs := []ce.SigAlg{
		ce.ES224,
		ce.ES256,
		ce.ES384,
		ce.ES512,
		//ce.Ed25519 //TODO
	}

	for _, alg := range algs {
		cozeKey, err := NewKey(ce.SEAlg(alg))
		if err != nil {
			t.Fatal(err)
		}

		s, err := cozeKey.SignRaw(TestMsg)
		if err != nil {
			t.Fatal(err)
		}
		v, err := cozeKey.VerifyRaw(TestMsg, s)
		if err != nil {
			t.Fatal(err)
		}
		if v != true {
			t.Fatal("signature was not verified")
		}
	}
}

// go test -bench=.
func BenchmarkGenKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n := b.N % 5
		algs := []ce.SigAlg{
			ce.ES224,
			ce.ES256,
			ce.ES384,
			ce.ES512,
			ce.Ed25519,
		}

		cozeKey, err := NewKey(ce.SEAlg(algs[n]))
		if err != nil {
			panic(err)
		}

		s, err := cozeKey.SignRaw(TestMsg)
		if err != nil {
			panic(err)
		}
		v, err := cozeKey.VerifyRaw(TestMsg, s)
		if err != nil {
			panic(err)
		}
		if v != true {
			panic("signature was not verified")
		}
	}
}

// BenchmarkNSV benchmarks several methods on a Coze Key. (NSV = New, Sign,
// Verify) It generatea a new Coze Key, sign a message, and verifies the
// signature.
// `go test -bench=.`
func BenchmarkNSV(b *testing.B) {
	var passCount = 0

	// TODO Ed25519 Support:
	var algs = []ce.SigAlg{ce.ES224, ce.ES256, ce.ES384, ce.ES512}

	for j := 0; j < b.N; j++ {
		for i := 0; i < len(algs); i++ {
			// log.Printf("Alg: %s\n", algs[i])
			ck, err := NewKey(ce.SEAlg(algs[i]))
			//log.Printf("Alg: %+v, Key: %+v\n", ck.Alg, ck)
			if err != nil {
				panic("Could not generate Coze Key.")
			}

			sig, err := ck.SignRaw(TestMsg)
			if err != nil {
				panic(err)
			}

			valid, err := ck.VerifyRaw(TestMsg, sig)
			if err != nil {
				panic(err)
			}
			if !valid {
				panic("The signature was invalid")
			}

			passCount++
		}
	}

	fmt.Printf("TestCryptoKeyNSV Pass Count: %+v \n", passCount)
}

func Example_es256_nsv() {
	msg := []byte("Test message.")

	ck, err := NewKey(ce.SEAlg(ce.ES256))
	// log.Printf("Alg: %+v, Key: %+v\n", ck.Alg, ck)
	if err != nil {
		panic("Could not generate Coze Key.")
	}

	sig, err := ck.SignRaw(msg)
	if err != nil {
		panic(err)
	}

	valid, err := ck.VerifyRaw(msg, sig)
	if err != nil {
		panic(err)
	}
	if !valid {
		panic("The signature was invalid")
	}

	fmt.Println(valid)

	// Output:
	// true
}
