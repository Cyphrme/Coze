package coze

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	ce "github.com/cyphrme/coze/enum"
)

var TestMsg = []byte("Coze Key Test Message.")

var Golden_Key = CozeKey{
	Alg: ce.SEAlg(ce.ES256),
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

// The very last byte in D was purposely changed from 80, to 81, which should make it invalid.
var Golden_Bad_Key = &CozeKey{
	Alg: ce.SEAlg(ce.ES256),
	Kid: "Zami's Majuscule Key.",
	Iat: 1623132000,
	X:   []byte{218, 116, 206, 104, 85, 102, 217, 2, 241, 153, 67, 191, 74, 56, 50, 177, 197, 71, 6, 219, 199, 17, 250, 54, 174, 174, 185, 50, 248, 13, 70, 51, 145, 162, 58, 183, 244, 118, 170, 246, 181, 205, 198, 245, 241, 193, 182, 191, 94, 61, 5, 230, 246, 98, 108, 148, 119, 138, 192, 93, 57, 102, 232, 230},
	D:   []byte{108, 219, 45, 131, 143, 199, 222, 109, 210, 149, 19, 174, 127, 4, 82, 18, 8, 155, 46, 176, 110, 70, 175, 117, 215, 131, 175, 117, 170, 92, 165, 81},
	Tmb: []byte{112, 184, 252, 190, 198, 45, 48, 28, 24, 147, 58, 5, 85, 145, 193, 102, 142, 146, 52, 191, 48, 73, 208, 136, 140, 34, 128, 193, 115, 110, 132, 233},
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
	Golden_Key.Tmb = []byte{} // set it to nil to ensure recalc
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
	sig, err := Golden_Key.SignMsg(TestMsg)
	if err != nil {
		fmt.Println(err)
	}

	valid, err := Golden_Key.VerifyMsg(TestMsg, sig)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%v\n", valid)
	// Output: true
}

func ExampleCozeKey_SignCy() {
	cy := new(Cy)
	cy.Head = []byte(Golden_Head)

	err := Golden_Key.SignCy(cy, nil)
	if err != nil {
		fmt.Println(err)
	}

	// To print sig, use:
	//fmt.Printf("%+s\n %s\n, %v\n", cy.Head, cy.Sig, []byte(cy.Sig))
	v, err := cy.Verify(&Golden_Key, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v)

	cy.Sig = nil // nil sig for marshal since for ecdsa it's non-deterministic.
	b, err := Marshal(cy)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s", b)

	// Output:
	// true
	//{"head":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"}}
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

func ExampleCozeKey_SignRaw_ed25519() {
	ck, err := NewKey(ce.SEAlg(ce.Ed25519))
	if err != nil {
		fmt.Println(err)
	}

	s, err := ck.SignMsg(TestMsg)
	if err != nil {
		fmt.Println(err)
	}
	v, err := ck.VerifyMsg(TestMsg, s)
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
	var sb []byte
	sb, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString("ILgt54kVj4r9pU1m3VJWTu4BcwZCHIAmqvnhqNploc9uiAA2EFpJLN65PrQ39PAt5WF41NtNPS4gvxIITU7rsw")
	if err != nil {
		fmt.Println(err)
	}

	valid, err := Golden_Key.VerifyMsg(TestMsg, sb)

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

		s, err := cozeKey.SignMsg(TestMsg)
		if err != nil {
			t.Fatal(err)
		}
		v, err := cozeKey.VerifyMsg(TestMsg, s)
		if err != nil {
			t.Fatal(err)
		}
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
	var algs = []ce.SigAlg{ce.ES224, ce.ES256, ce.ES384, ce.ES512, ce.Ed25519}

	for j := 0; j < b.N; j++ {
		for i := 0; i < len(algs); i++ {
			ck, err := NewKey(ce.SEAlg(algs[i]))
			if err != nil {
				b.Fatal("Could not generate Coze Key.")
			}

			sig, err := ck.SignMsg(TestMsg)
			if err != nil {
				b.Fatal(err)
			}

			valid, err := ck.VerifyMsg(TestMsg, sig)
			if err != nil {
				b.Fatal(err)
			}
			if !valid {
				b.Fatalf("The signature was invalid.  Alg: %s", ck.Alg)
			}
		}
	}

}

func Example_eS256_nsv() {
	msg := []byte("Test message.")

	ck, err := NewKey(ce.SEAlg(ce.ES256))
	// log.Printf("Alg: %+v, Key: %+v\n", ck.Alg, ck)
	if err != nil {
		panic("Could not generate Coze Key.")
	}

	sig, err := ck.SignMsg(msg)
	if err != nil {
		panic(err)
	}

	valid, _ := ck.VerifyMsg(msg, sig)
	if !valid {
		panic("The signature was invalid")
	}

	fmt.Println(valid)

	// Output:
	// true
}
