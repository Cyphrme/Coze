package coze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	ce "github.com/zamicol/cyphrme/src/coze/enum"
)

///////////////////
///////////////////
// Golden_Cy
///////////////////
///////////////////

var Golden_Cy = &Cy{
	Head: Golden_Cy_Head,
	Sch: Head{ // Don't do this in practice.  This only works because there's no extra fields.  Use Cy.Head (which is rawjson) instead.
		Alg: ce.SEAlg(ce.ES256),
		Iat: 1623132000,
		Tmb: []byte{1, 72, 244, 205, 144, 147, 201, 203, 227, 232, 191, 120, 211, 230, 201, 184, 36, 241, 29, 210, 242, 158, 43, 26, 99, 13, 209, 206, 30, 23, 108, 221},
		Typ: "cyphr.me",
	},
	Key: &CozeKey{
		Alg: ce.SEAlg(ce.ES256),
		Kid: "Zami's Majuscule Key.",
		Iat: 1623132000,
		Tmb: []byte{1, 72, 244, 205, 144, 147, 201, 203, 227, 232, 191, 120, 211, 230, 201, 184, 36, 241, 29, 210, 242, 158, 43, 26, 99, 13, 209, 206, 30, 23, 108, 221},
		X:   []byte{218, 116, 206, 104, 85, 102, 217, 2, 241, 153, 67, 191, 74, 56, 50, 177, 197, 71, 6, 219, 199, 17, 250, 54, 174, 174, 185, 50, 248, 13, 70, 51},
		Y:   []byte{145, 162, 58, 183, 244, 118, 170, 246, 181, 205, 198, 245, 241, 193, 182, 191, 94, 61, 5, 230, 246, 98, 108, 148, 119, 138, 192, 93, 57, 102, 232, 230},
	},
	Sig: []byte{255, 153, 117, 161, 218, 44, 8, 28, 193, 37, 222, 27, 161, 173, 214, 217, 165, 212, 77, 22, 0, 207, 142, 54, 84, 10, 32, 79, 96, 213, 32, 219, 102, 228, 153, 105, 5, 237, 176, 252, 143, 47, 142, 237, 118, 134, 183, 44, 103, 183, 95, 209, 231, 10, 248, 18, 126, 7, 239, 235, 9, 141, 195, 47},
}

var Golden_Cy_Head = []byte(`{
	 "alg": "ES256",
	 "iat": 1623132000,
	 "tmb": "0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD",
	 "typ": "cyphr.me"
	}`)

// head's `tmb` and key's `alg` are purposely out of order for testing.
const Golden_Cy_String = `{
	"head": {
	 "alg": "ES256",
	 "iat": 1623132000,
	 "typ": "cyphr.me",
	 "tmb": "0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD"
	},
	"key": {
	 "iat": 1623132000,
	 "kid": "Zami's Majuscule Key.",
	 "tmb": "0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD",
	 "x": "DA74CE685566D902F19943BF4A3832B1C54706DBC711FA36AEAEB932F80D4633",
	 "y": "91A23AB7F476AAF6B5CDC6F5F1C1B6BF5E3D05E6F6626C94778AC05D3966E8E6",
	 "alg": "ES256"
	},
	"sig": "FF9975A1DA2C081CC125DE1BA1ADD6D9A5D44D1600CF8E36540A204F60D520DB66E4996905EDB0FC8F2F8EED7686B72C67B75FD1E70AF8127E07EFEB098DC32F"
 }`

///////////////////
///////////////////
// Golden_Msg
///////////////////
///////////////////

var Golden_Msg_Head = []byte(`{
	"msg": "Coze Rocks",
	"alg": "ES256",
	"iat": 1627518000,
	"tmb": "0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD",
	"typ": "cyphr.me/msg"
 }`)

// cad: BCE8938B3CD933036C60C9DE003DB07245A55EE1C8E711A61BC10C59E0AA0A24
const Golden_Msg_String = `{
	"head": {
		"msg": "Coze Rocks",
		"alg": "ES256",
		"iat": 1627518000,
		"tmb": "0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD",
		"typ": "cyphr.me/msg"
	},
	"sig": "B40F147AC98726737FF1FBF64B33F6A2BA9EDEEBAC140491465D6EC460DC6595A9DF4593D6E908042C5F4B7F83FFB4948C6285CDAE421741B79E15AB0C7D8CEC"
 }`

const Golden_Msg_Cy_String = `{
	"cy":{
		"head": {
			"msg": "Coze Rocks",
			"alg": "ES256",
			"iat": 1627518000,
			"tmb": "0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD",
			"typ": "cyphr.me/msg"
		},
		"sig": "B40F147AC98726737FF1FBF64B33F6A2BA9EDEEBAC140491465D6EC460DC6595A9DF4593D6E908042C5F4B7F83FFB4948C6285CDAE421741B79E15AB0C7D8CEC"
	}
}`

//ExampleCy_jsonUnMarshal tests unmarshalling a `cy`.
func ExampleCy_jsonUnMarshal() {
	cy := new(Cy)
	err := json.Unmarshal([]byte(Golden_Cy_String), cy)
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
	// {"head":{"alg":"ES256","iat":1623132000,"tmb":"0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD","typ":"cyphr.me"},"key":{"alg":"ES256","iat":1623132000,"kid":"Zami's Majuscule Key.","tmb":"0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD","x":"DA74CE685566D902F19943BF4A3832B1C54706DBC711FA36AEAEB932F80D4633","y":"91A23AB7F476AAF6B5CDC6F5F1C1B6BF5E3D05E6F6626C94778AC05D3966E8E6"},"sig":"FF9975A1DA2C081CC125DE1BA1ADD6D9A5D44D1600CF8E36540A204F60D520DB66E4996905EDB0FC8F2F8EED7686B72C67B75FD1E70AF8127E07EFEB098DC32F"}
}

func ExampleCy_msgJsonUnMarshal() {
	cy := new(Cy)
	err := json.Unmarshal([]byte(Golden_Msg_String), cy)
	if err != nil {
		fmt.Println(err)
	}

	// remarshal for comparison
	b, err := Marshal(cy)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("Marshaled: %s\n", b)

	fmt.Println(string(b))
	// Output:
	//{"head":{"alg":"ES256","iat":1627518000,"msg":"Coze Rocks","tmb":"0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD","typ":"cyphr.me/msg"},"sig":"B40F147AC98726737FF1FBF64B33F6A2BA9EDEEBAC140491465D6EC460DC6595A9DF4593D6E908042C5F4B7F83FFB4948C6285CDAE421741B79E15AB0C7D8CEC"}
}

func ExampleCyEn_jsonMarshal() {
	var err error

	cye := new(CyEn)

	err = json.Unmarshal([]byte(Golden_Msg_Cy_String), cye)
	if err != nil {
		fmt.Println(err)
	}

	b, err := Marshal(cye)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+s\n", b)
	// Output:
	//{"cy":{"head":{"alg":"ES256","iat":1627518000,"msg":"Coze Rocks","tmb":"0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD","typ":"cyphr.me/msg"},"sig":"B40F147AC98726737FF1FBF64B33F6A2BA9EDEEBAC140491465D6EC460DC6595A9DF4593D6E908042C5F4B7F83FFB4948C6285CDAE421741B79E15AB0C7D8CEC"}}
}

// ExampleCy_jsonMarshal tests unmarshalling a `cy`.
func ExampleCy_jsonMarshal() {
	var err error

	cy := new(Cy)
	err = json.Unmarshal([]byte(Golden_Cy_String), cy)
	if err != nil {
		fmt.Println(err)
	}

	b, err := Marshal(cy)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+s\n", b)
	// Output:
	//{"head":{"alg":"ES256","iat":1623132000,"tmb":"0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD","typ":"cyphr.me"},"key":{"alg":"ES256","iat":1623132000,"kid":"Zami's Majuscule Key.","tmb":"0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD","x":"DA74CE685566D902F19943BF4A3832B1C54706DBC711FA36AEAEB932F80D4633","y":"91A23AB7F476AAF6B5CDC6F5F1C1B6BF5E3D05E6F6626C94778AC05D3966E8E6"},"sig":"FF9975A1DA2C081CC125DE1BA1ADD6D9A5D44D1600CF8E36540A204F60D520DB66E4996905EDB0FC8F2F8EED7686B72C67B75FD1E70AF8127E07EFEB098DC32F"}
}

// TestVerifyCy
func TestVerifyCy(t *testing.T) {
	var cy *Cy
	err := json.Unmarshal([]byte(Golden_Cy_String), &cy)
	if err != nil {
		t.Fatal(err)
	}

	cozekey := Golden_Key

	ck, err := cozekey.ToCryptoKey()
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Printf("ck: %+v\n", ck)
	// fmt.Printf("MinCy: %+v\n", cy)
	// fmt.Printf("\n Coze Key in testing: %+v \n", cozekey)
	// fmt.Printf("\n Crypto Key: %+v \n", ck)
	// fmt.Printf("\n Crypto Key Private: %+v \n", *ck.Private)
	// fmt.Printf("\n Crypto Key Public: %+v \n", *ck.Public)

	ch, err := CH(cy.Head, nil, ce.Sha256)
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Printf("Canonical Hash (cad): %+v\n", ch)
	// fmt.Printf("Head String: %+v\n", string(headB))
	// fmt.Printf("CAD Hex: %X\n", ch)
	// fmt.Printf("Sig Hex: %s\n", cy.Sig)

	valid, err := ck.VerifyDigest(ch, cy.Sig)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("Not a valid cy.sig.")
	}
}

func ExampleCy_Verify() {
	v, err := Golden_Cy.Verify(&Golden_Key, nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(v)
	// Output: true
}

func ExampleVerify() {
	v, err := Verify(Golden_Cy, &Golden_Key, Golden_Cy.Sig, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v)
	// Output: true
}

func TestVerifyCyMsg(t *testing.T) {
	var cy *Cy
	err := json.Unmarshal([]byte(Golden_Cy_String), &cy)
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

	valid, err := Golden_Key.VerifyRaw(cy.Head, cy.Sig)

	if err != nil {
		t.Fatal(err)
	}

	if !valid {
		t.Fatal("Not a valid cy.sig.")
	}

	// Test again with the manually calculated digest.
	ch, err := CH(cy.Head, nil, ce.Sha256)
	if err != nil {
		t.Fatal(err)
	}

	valid, err = Golden_Key.VerifyDigest(ch, cy.Sig)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("Not a valid cy.sig.")
	}
}

// This function takes the Golden Majuscule Key and signs over a message.
// The signature is then verified.
func TestValidateGoldenKeySig(t *testing.T) {
	msg := []byte("This is a test sign for coze.")
	// Golden_Key is in coze_test
	sig, err := Golden_Key.SignRaw(msg)
	if err != nil {
		t.Fatal(err)
	}

	valid, err := Golden_Key.VerifyRaw(msg, sig)
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
	err := json.Unmarshal([]byte(Golden_Cy_String), &cy)
	if err != nil {
		fmt.Println(err)
	}
	err = cy.SetMeta()
	if err != nil {
		fmt.Println(err)
	}
	// debugging
	// fmt.Printf("Cy.Head: %s\n,, Cy.Head: %+v\n", cy.Head, cy.Head)

	if bytes.Compare(cy.Sch.Tmb, Golden_Key.Tmb) != 0 {
		fmt.Println("coze: key thumbprints do not match")
	}

	b, err := Marshal(cy.Head)
	if err != nil {
		fmt.Println(err)
	}

	t, err := Golden_Key.VerifyRaw(b, cy.Sig)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(t)
	// Output: true
}

func ExampleCy_SetMeta() {
	cy := new(Cy)
	err := json.Unmarshal([]byte(Golden_Cy_String), cy)
	if err != nil {
		fmt.Println(err)
	}

	err = cy.SetMeta()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Can: %s\n", cy.Can)
	cy.Can = nil // for testing, don't print memory address
	fmt.Printf("Head: %s\n", cy.Head)
	cy.Head = nil
	fmt.Printf("Cy: %+v\n", cy)

	// Output:
	// Can: [alg iat tmb typ]
	// Head: {"alg":"ES256","iat":1623132000,"tmb":"0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD","typ":"cyphr.me"}
	// Cy: &{Cad:0C359495353CD108BF5477F1084B4C2A656C565D2168D496A149B0990AE94286 Can:[] Cyd:47859B79167738FB1EB2390840BD53C596E7083157EBFA8FA1167A8919B680B8 Head:[] Key:{"alg":"ES256","kid":"Zami's Majuscule Key.","iat":1623132000,"tmb":"0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD","x":"DA74CE685566D902F19943BF4A3832B1C54706DBC711FA36AEAEB932F80D4633","y":"91A23AB7F476AAF6B5CDC6F5F1C1B6BF5E3D05E6F6626C94778AC05D3966E8E6"} Sig:FF9975A1DA2C081CC125DE1BA1ADD6D9A5D44D1600CF8E36540A204F60D520DB66E4996905EDB0FC8F2F8EED7686B72C67B75FD1E70AF8127E07EFEB098DC32F Sigs:[] Sch:{Alg:ES256 Iat:1623132000 Tmb:0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD Typ:cyphr.me}}
}

func ExampleGenCyd() {
	Golden_Cy.Cyd = nil // ensure cyd is empty for test.
	cad, err := CH(Golden_Cy.Head, nil, Golden_Cy.Sch.Alg.Hash())
	if err != nil {
		fmt.Println(err)
	}
	Golden_Cy.Cad = cad
	Golden_Cy.Cyd = GenCyd(Golden_Cy.Sch.Alg.Hash(), Golden_Cy.Cad, Golden_Cy.Sig)

	fmt.Println(Golden_Cy.Cad)
	fmt.Println(Golden_Cy.Cyd)

	// Output:
	//0C359495353CD108BF5477F1084B4C2A656C565D2168D496A149B0990AE94286
	//47859B79167738FB1EB2390840BD53C596E7083157EBFA8FA1167A8919B680B8
}
