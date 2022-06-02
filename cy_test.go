package coze

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	ce "github.com/cyphrme/coze/enum"
)

var testMsg = []byte("Coze Rock")

var Golden_Head = `{
	"msg": "Coze Rocks",
	"alg": "ES256",
	"iat": 1627518000,
	"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
	"typ": "cyphr.me/msg"
 }`

var Golden_Sig = "Z8yK1AuBWdfGzwmXK_xwlZizlFsxFkK7bKJ8FEDoNEA1IFJECjaK0ZLPLDIFhLX6kD8jis-9tCKlB1Qzb-mEzg"

var Golden_Cy = `{
	"head":` + Golden_Head + `,
	"sig": "` + Golden_Sig + `"
 }`

var Golden_Cy_En = `{
	"coze":` + Golden_Cy + `
}`

var Golden_Cy_W_Key = `{
	"head": ` + Golden_Head + `,
	"key": ` + Golden_Key_String + `
	,
	"sig": "` + Golden_Sig + `"
 }`

//ExampleCy_jsonUnMarshal tests unmarshalling a `coze`.
func ExampleCy_jsonUnMarshal() {
	cy := new(Cy)
	err := json.Unmarshal([]byte(Golden_Cy), cy)
	if err != nil {
		fmt.Println(err)
	}

	// remarshal for comparison
	b, err := Marshal(cy)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
	// Output:
	//{"head":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"Z8yK1AuBWdfGzwmXK_xwlZizlFsxFkK7bKJ8FEDoNEA1IFJECjaK0ZLPLDIFhLX6kD8jis-9tCKlB1Qzb-mEzg"}
}

func ExampleCoze_jsonMarshal() {
	var err error

	cye := new(Coze)

	err = json.Unmarshal([]byte(Golden_Cy_En), cye)
	if err != nil {
		fmt.Println(err)
	}

	b, err := Marshal(cye)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+s\n", b)
	// Output:
	//{"coze":{"head":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"Z8yK1AuBWdfGzwmXK_xwlZizlFsxFkK7bKJ8FEDoNEA1IFJECjaK0ZLPLDIFhLX6kD8jis-9tCKlB1Qzb-mEzg"}}
}

// TestVerifyCy
func TestVerifyCy(t *testing.T) {
	var cy *Cy
	err := json.Unmarshal([]byte(Golden_Cy), &cy)
	if err != nil {
		t.Fatal(err)
	}

	cozekey := Golden_Key

	ck, err := cozekey.ToCryptoKey()
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Printf("ck: %+v\n", ck)
	// fmt.Printf("\n Coze Key in testing: %+v \n", cozekey)
	// fmt.Printf("\n Crypto Key: %+v \n", ck)
	// fmt.Printf("\n Crypto Key Private: %+v \n", *ck.Private)
	// fmt.Printf("\n Crypto Key Public: %+v \n", *ck.Public)

	ch, err := CanonHash(cy.Head, nil, ce.Sha256)
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Printf("Canonical Hash (cad): %+v\n", ch)
	// fmt.Printf("Head String: %+v\n", string(headB))
	// fmt.Printf("CAD Hex: %X\n", ch)
	// fmt.Printf("Sig Hex: %s\n", cy.Sig)

	valid, err := ck.Verify(ch, cy.Sig)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("Not a valid cy.sig.")
	}
}

func ExampleCy_Verify() {
	var cy Cy
	err := json.Unmarshal([]byte(Golden_Cy), &cy)
	if err != nil {
		fmt.Println(err)
	}
	v, err := cy.Verify(&Golden_Key, nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(v)
	// Output: true
}

func ExampleVerify() {
	var cy = new(Cy)
	err := json.Unmarshal([]byte(Golden_Cy), cy)
	if err != nil {
		fmt.Println(err)
	}
	cy.Sig, err = base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(Golden_Sig)
	if err != nil {
		fmt.Println(err)
	}

	v, err := cy.Verify(&Golden_Key, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v)
	// Output: true
}

func TestVerifyCyMsg(t *testing.T) {
	cy := new(Cy)
	err := json.Unmarshal([]byte(Golden_Cy), cy)
	if err != nil {
		t.Fatal(err)
	}

	// Unmarshal does not normalize head bytes.  SetMeta() does.
	err = cy.SetMeta()
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Printf("Cy: %+v\n", cy)
	// fmt.Printf("Head: %s\n", cy.Head)

	valid, err := Golden_Key.VerifyMsg(cy.Head, cy.Sig)

	if err != nil {
		t.Fatal(err)
	}

	if !valid {
		t.Fatal("Not a valid cy.sig.")
	}

	// Test again with the manually calculated digest.
	ch, err := CanonHash(cy.Head, nil, ce.Sha256)
	if err != nil {
		t.Fatal(err)
	}

	valid, err = Golden_Key.Verify(ch, cy.Sig)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("Not a valid cy.sig.")
	}
}

// Example of how to manually verify a Cy if not using the Cy.Verify method.
func ExampleCy_verify_manual() {
	cy := new(Cy)
	err := json.Unmarshal([]byte(Golden_Cy), cy)
	if err != nil {
		fmt.Println(err)
	}
	err = cy.SetMeta()
	if err != nil {
		fmt.Println(err)
	}
	// debugging
	// fmt.Printf("Cy.Head: %s\n,, Cy.Head: %+v\n", cy.Head, cy.Head)

	if bytes.Compare(cy.Parsed.Tmb, Golden_Key.Tmb) != 0 {
		fmt.Println("coze: key thumbprints do not match")
	}

	b, err := Marshal(cy.Head)
	if err != nil {
		fmt.Println(err)
	}

	t, err := Golden_Key.VerifyMsg(b, cy.Sig)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(t)
	// Output: true
}

func ExampleCy_SetMeta() {
	cy := new(Cy)
	err := json.Unmarshal([]byte(Golden_Cy), cy)
	if err != nil {
		fmt.Println(err)
	}

	err = cy.SetMeta()
	if err != nil {
		fmt.Println(err)
	}

	cyb, err := Marshal(cy)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", cyb)

	// Output:
	// {"cad":"aC2YKfNvovfnZOw_RVxSEW6NeaUq41DZXX0oeaOboRg","can":["alg","iat","msg","tmb","typ"],"cyd":"D9riVnxvV5qxoJLFbq4pzhoetcOaKASHvln_C7aip3I","head":{"alg":"ES256","iat":1627518000,"msg":"Coze Rocks","tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"Z8yK1AuBWdfGzwmXK_xwlZizlFsxFkK7bKJ8FEDoNEA1IFJECjaK0ZLPLDIFhLX6kD8jis-9tCKlB1Qzb-mEzg"}
}

func ExampleGenCyd_manual() {
	cy := new(Cy)
	err := json.Unmarshal([]byte(Golden_Cy), cy)
	if err != nil {
		fmt.Println(err)
	}

	cy.SetMeta()
	fmt.Println(cy.Cad)
	fmt.Println(cy.Cyd)

	// Set `cyd` and `cad` empty for test.
	cy.Cyd = nil
	cy.Cad = nil

	// "Manually" recalculate values.
	cy.Cad, err = CanonHash(cy.Head, nil, cy.Parsed.Alg.Hash())
	if err != nil {
		fmt.Println(err)
	}

	cy.Cyd = GenCyd(cy.Parsed.Alg.Hash(), cy.Cad, cy.Sig)

	fmt.Println(cy.Cad)
	fmt.Println(cy.Cyd)

	// Output:
	// aC2YKfNvovfnZOw_RVxSEW6NeaUq41DZXX0oeaOboRg
	// D9riVnxvV5qxoJLFbq4pzhoetcOaKASHvln_C7aip3I
	// aC2YKfNvovfnZOw_RVxSEW6NeaUq41DZXX0oeaOboRg
	// D9riVnxvV5qxoJLFbq4pzhoetcOaKASHvln_C7aip3I
}
