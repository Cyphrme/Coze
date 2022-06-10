package coze

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cyphrme/coze/enum"
)

var testMsg = []byte("Coze Rock")

var Golden_Pay = `{
	"msg": "Coze Rocks",
	"alg": "ES256",
	"iat": 1627518000,
	"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
	"typ": "cyphr.me/msg"
 }`

var Golden_Sig = "Z8yK1AuBWdfGzwmXK_xwlZizlFsxFkK7bKJ8FEDoNEA1IFJECjaK0ZLPLDIFhLX6kD8jis-9tCKlB1Qzb-mEzg"

var Golden_Cy = `{
	"pay":` + Golden_Pay + `,
	"sig": "` + Golden_Sig + `"
 }`

var Golden_Cy_En = `{
	"coze":` + Golden_Cy + `
}`

var Golden_Cy_W_Key = `{
	"pay": ` + Golden_Pay + `,
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
	//{"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"Z8yK1AuBWdfGzwmXK_xwlZizlFsxFkK7bKJ8FEDoNEA1IFJECjaK0ZLPLDIFhLX6kD8jis-9tCKlB1Qzb-mEzg"}
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
	//{"coze":{"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"Z8yK1AuBWdfGzwmXK_xwlZizlFsxFkK7bKJ8FEDoNEA1IFJECjaK0ZLPLDIFhLX6kD8jis-9tCKlB1Qzb-mEzg"}}
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

	ch, err := CanonHash(cy.Pay, nil, enum.Sha256)
	if err != nil {
		t.Fatal(err)
	}

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

	// Unmarshal does not normalize pay bytes.  SetMeta() does.
	err = cy.SetMeta()
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Printf("Cy: %+v\n", cy)
	// fmt.Printf("Pay: %s\n", cy.Pay)

	valid, err := Golden_Key.VerifyMsg(cy.Pay, cy.Sig)

	if err != nil {
		t.Fatal(err)
	}

	if !valid {
		t.Fatal("Not a valid cy.sig.")
	}

	// Test again with the manually calculated digest.
	ch, err := CanonHash(cy.Pay, nil, enum.Sha256)
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
	fmt.Printf("Cy.Pay: %s\n,, Cy.Pay: %+v\n", cy.Pay, cy.Pay)

	if bytes.Compare(cy.Parsed.Tmb, Golden_Key.Tmb) != 0 {
		fmt.Println("coze: key thumbprints do not match")
	}

	b, err := Marshal(cy.Pay)
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
	// {"can":["alg","iat","msg","tmb","typ"],"cad":"aC2YKfNvovfnZOw_RVxSEW6NeaUq41DZXX0oeaOboRg","cyd":"D9riVnxvV5qxoJLFbq4pzhoetcOaKASHvln_C7aip3I","pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"Z8yK1AuBWdfGzwmXK_xwlZizlFsxFkK7bKJ8FEDoNEA1IFJECjaK0ZLPLDIFhLX6kD8jis-9tCKlB1Qzb-mEzg"}
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
	cy.Cad, err = CanonHash(cy.Pay, nil, cy.Parsed.Alg.Hash())
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
