package coze

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func ExamplePay_embedded() {
	// Example custom struct.
	type User struct {
		DisplayName string
		FirstName   string
		LastName    string
		Email       string `json:",omitempty"` // Example of non-required field.
	}

	user := User{
		DisplayName: "Coze",
		FirstName:   "Foo",
		LastName:    "Bar",
	}

	// Example of converting a custom struct to a coze.
	pay := Pay{
		Alg:    GoldenKey.Alg,
		Tmb:    GoldenKey.Tmb,
		Struct: &user,
	}

	coze, err := GoldenKey.SignPay(&pay)
	if err != nil {
		panic(err)
	}

	v, err := GoldenKey.VerifyCoze(coze)
	if err != nil {
		panic(err)
	}

	// Set sig to nil for deterministic printout
	coze.Sig = nil
	fmt.Println(v)
	fmt.Printf("%+v\n", coze)

	// Output:
	// true
	// {"pay":{"alg":"ES256","tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","DisplayName":"Coze","FirstName":"Foo","LastName":"Bar"}}
}

// ExamplePay_jsonUnmarshal tests unmarshalling a Pay.
func ExamplePay_jsonUnmarshal() {
	h := &Pay{}

	err := json.Unmarshal([]byte(GoldenPay), h)
	if err != nil {
		panic(err)
	}

	out, err := Marshal(h)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", out)

	// Output:
	// {"alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"}
}

// ExamplePay_jsonMarshalCustom demonstrates marshalling Pay with a custom
// structure.
func ExamplePay_jsonMarshalCustom() {
	customStruct := CustomStruct{
		Msg: "Coze Rocks",
	}

	inputPay := Pay{
		Alg:    SEAlg(ES256),
		Iat:    1623132000, // Static for demonstration.  Use time.Now().Unix().
		Tmb:    MustDecode("cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk"),
		Typ:    "cyphr.me/msg",
		Struct: customStruct,
	}

	// May also call inputPay.MarshalJSON() or Marshal(&inputPay) instead.
	s, err := Marshal(&inputPay)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(s))

	// Output:
	// {"alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg","msg":"Coze Rocks"}
}

// ExamplePay_jsonUnmarshalCustomManual demonstrates "manually" unmarshalling
// Pay with a custom structure.
func ExamplePay_jsonUnmarshalCustomManual() {
	var pay Pay
	err := json.Unmarshal([]byte(GoldenPay), &pay)
	if err != nil {
		panic(err)
	}
	fmt.Println(pay)

	var custom CustomStruct
	err = json.Unmarshal([]byte(GoldenPay), &custom)
	if err != nil {
		panic(err)
	}
	fmt.Println(custom)

	// Output:
	// {"alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"}
	// {Coze Rocks}
}

// ExamplePay_jsonUnmarshalCustom demonstrates unmarshalling Pay with a custom
// structure.
func ExamplePay_jsonUnmarshalCustom() {
	pay := new(Pay)
	var emptyCustomStruct CustomStruct
	pay.Struct = &emptyCustomStruct
	err := json.Unmarshal([]byte(GoldenPay), &pay)
	if err != nil {
		fmt.Printf("Unmarshal error: %s\n", err)
	}
	fmt.Println(pay)
	fmt.Println(pay.Struct)

	// Output:
	// {"alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg","msg":"Coze Rocks"}
	// &{Coze Rocks}
}

// ExamplePay_String_custom demonstrates fmt.Stringer on Pay with a custom
// structure.
func ExamplePay_String_custom() {
	customStruct := CustomStruct{
		Msg: "Coze Rocks",
	}

	inputPay := Pay{
		Alg:    SEAlg(ES256),
		Iat:    1623132000, // Static for demonstration.  Use time.Now().Unix().
		Tmb:    MustDecode("cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk"),
		Typ:    "cyphr.me/msg",
		Struct: customStruct,
	}
	fmt.Println(inputPay)

	// Output:
	// {"alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg","msg":"Coze Rocks"}
}

// Example demonstrating that unmarshalling a `pay` that has duplicate field
// names results in an error.
func ExamplePay_UnmarshalJSON_duplicate() {
	h := &Pay{}
	msg := []byte(`{"alg":"ES256","alg":"ES384"}`)

	err := json.Unmarshal(msg, h)
	fmt.Println(err)

	// Output:
	// Coze: JSON duplicate field name
}

// Example demonstrating that unmarshalling a `coze` that has duplicate field
// names results in an error.
func ExampleCoze_UnmarshalJSON_duplicate() {
	h := &Pay{}
	msg := []byte(`{"coze":{"pay":"ES256","pay":"ES384"}}`)
	err := json.Unmarshal(msg, h)
	fmt.Println(err)

	// Output:
	// Coze: JSON duplicate field name
}

// ExampleCoze_embed demonstrates how to embed a JSON `coze` into a third party
// JSON structure.
func ExampleCoze_embed() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(GoldenCoze), cz)
	if err != nil {
		panic(err)
	}

	type Outer struct {
		Name string `json:"name"`
		Coze Coze   `json:"coze"` // Embed a Coze into a larger, application defined JSON structure.
	}
	b, _ := json.Marshal(Outer{Name: "Bob", Coze: *cz})
	fmt.Printf("%s", b)

	// Output:
	// {"name":"Bob","coze":{"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w"}}
}

func ExampleCoze_String() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(GoldenCoze), cz)
	if err != nil {
		panic(err)
	}
	fmt.Println(cz)

	// Output:
	// {"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w"}
}

func ExampleCoze_Meta() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(GoldenCoze), cz)
	if err != nil {
		panic(err)
	}

	err = cz.Meta()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", cz)

	// Output:
	//{"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"can":["msg","alg","iat","tmb","typ"],"cad":"Ie3xL77AsiCcb4r0pbnZJqMcfSBqg5Lk0npNJyJ9BC4","sig":"Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w","czd":"TnRe4DRuGJlw280u3pGhMDOIYM7ii7J8_PhNuSScsIU"}
}

func ExampleCoze_MetaWithAlg() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(GoldenCoze), cz)
	if err != nil {
		panic(err)
	}

	// Test mismatch alg, which should error.
	err = cz.MetaWithAlg(SEAlg(ES224))
	if err == nil {
		fmt.Println("Test should error")
	}

	// Test with correct alg, no error.
	err = cz.MetaWithAlg(SEAlg(ES256))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", cz)

	// No alg given.  Alg is parsed from pay.
	err = cz.MetaWithAlg("")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", cz)

	// Output:
	// {"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"can":["msg","alg","iat","tmb","typ"],"cad":"Ie3xL77AsiCcb4r0pbnZJqMcfSBqg5Lk0npNJyJ9BC4","sig":"Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w","czd":"TnRe4DRuGJlw280u3pGhMDOIYM7ii7J8_PhNuSScsIU"}
	// {"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"can":["msg","alg","iat","tmb","typ"],"cad":"Ie3xL77AsiCcb4r0pbnZJqMcfSBqg5Lk0npNJyJ9BC4","sig":"Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w","czd":"TnRe4DRuGJlw280u3pGhMDOIYM7ii7J8_PhNuSScsIU"}
}

func ExampleCoze_MetaWithAlg_contextual() {
	// Test MetaWithAlg using no sig, which should calc what it can.
	cz := new(Coze)
	err := json.Unmarshal([]byte(GoldenCoze), cz)
	if err != nil {
		panic(err)
	}
	cz.Sig = []byte{} // set sig to nothing.
	err = cz.MetaWithAlg("")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", cz)

	// Empty coze with coze.parsed.alg
	cz = new(Coze)
	err = json.Unmarshal(GoldenCozeEmpty, cz)
	if err != nil {
		panic(err)
	}
	err = cz.MetaWithAlg(GoldenKey.Alg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", cz)

	// Output:
	// {"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"can":["msg","alg","iat","tmb","typ"],"cad":"Ie3xL77AsiCcb4r0pbnZJqMcfSBqg5Lk0npNJyJ9BC4"}
	// {"pay":{},"cad":"RBNvo1WzZ4oRRq0W9-hknpT7T8If536DEMBg9hyq_4o","sig":"9iesKUSV7L1-xz5yd3A94vCkKLmdOAnrcPXTU3_qeKRRbHuy5EvMMFNRkW_sNLo-vvEPO9BmeUkcNh-ok18I_A","czd":"zU7xRwp8XU_VmdOLNBlMBualhoyHiM_cGhib6LPwWlc"}
}

// ExampleCoze_jsonUnmarshal tests unmarshalling a coze.
func ExampleCoze_jsonUnmarshal() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(GoldenCoze), cz)
	if err != nil {
		panic(err)
	}

	// remarshal for comparison
	b, err := Marshal(cz)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	// Output:
	//{"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w"}
}

func ExampleCoze_jsonMarshal() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(GoldenCoze), cz)
	if err != nil {
		panic(err)
	}

	b, err := Marshal(cz)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+s\n", b)

	// Output:
	//{"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w"}
}

func ExampleCoze_jsonMarshalPretty() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(GoldenCoze), cz)
	if err != nil {
		panic(err)
	}

	b, err := MarshalPretty(cz)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+s\n", b)

	// Output:
	// 	{
	//     "pay": {
	//         "msg": "Coze Rocks",
	//         "alg": "ES256",
	//         "iat": 1623132000,
	//         "tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
	//         "typ": "cyphr.me/msg"
	//     },
	//     "sig": "Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w"
	// }
}

// ExampleMarshal_jsonRawMessage demonstrates using empty string, two quote
// characters, and nil for json.RawMessage.  When using json.RawMessage, it
// should always be valid JSON or nil or otherwise it will result in an error.
func ExampleMarshal_jsonRawMessage() {
	o := json.RawMessage([]byte("")) // empty string
	anon := struct {
		Obj *json.RawMessage `json:"obj,omitempty"`
	}{
		Obj: &o,
	}

	// Incorrect usage with pointer to a zero value string.
	// Pointer to empty string will fail Marshal since an empty string is not
	// valid JSON., while the value `""` will pass.
	b, err := Marshal(anon)         // fails
	fmt.Printf("%s\n%+v\n", err, b) // Error is populated and b is empty because of error.

	// Correct usage with quotes characters.
	quotes := []byte("\"\"") // string with two quote characters
	anon.Obj = (*json.RawMessage)(&quotes)
	b, err = Marshal(anon)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)

	// Correct usage with with `nil`, which prints as the JSON "null".
	o = nil
	anon.Obj = &o
	b, err = Marshal(anon)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)

	// Output:
	// json: error calling MarshalJSON for type *json.RawMessage: unexpected end of JSON input
	// []
	// {"obj":""}
	// {"obj":null}
}

func Test_checkDuplicate(t *testing.T) {
	// Duplicate, should error.
	data := `{"a": "b", "a":true,"c":["field_3 string 1","field3 string2"], "d": {"e": 1, "e": 2}}`
	err := checkDuplicate(json.NewDecoder(strings.NewReader(data)))
	if err != ErrJSONDuplicate {
		t.Fatal("Should have found duplicate.")
	}

	// Recursive check with duplicate in inner struct.  Should error.
	data = `{"a": "b", "c":"d", "d": {"e": 1, "e": 2}}`
	err = checkDuplicate(json.NewDecoder(strings.NewReader(data)))
	if err != ErrJSONDuplicate {
		t.Fatal("Should have found duplicate.")
	}

	// No duplicate.  Should not error.
	data = `{"a": "b", "c":"d", "d": {"e": 1, "f": 2}}`
	err = checkDuplicate(json.NewDecoder(strings.NewReader(data)))
	if err != nil {
		t.Fatal(err)
	}
}

// Demonstrates expectations for values that are non-integer, negative, or too
// large for rvk.
func Example_iat_rvk_too_big() {
	p := &Pay{}

	//  2^53 - 1 as a string which must error.
	err := json.Unmarshal([]byte(`{"rvk":"9007199254740991"}`), p)
	if err != nil {
		fmt.Println(err)
	}

	// 2^53
	err = json.Unmarshal([]byte(`{"rvk":9007199254740992}`), p)
	if err != nil {
		fmt.Println(err)
	}

	//  Negative 2^53 + 1 must error as rvk must be positive.
	err = json.Unmarshal([]byte(`{"rvk":-9007199254740991}`), p)
	if err != nil {
		fmt.Println(err)
	}

	// Finally, 2^53 - 1 as an integer which is okay.
	err = json.Unmarshal([]byte(`{"rvk":9007199254740991}`), p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)

	// Output:
	// json: cannot unmarshal string into Go struct field pay2.rvk of type int64
	// Pay.UnmarshalJSON: values for iat and rvk must be between 0 and 2^53 - 1
	// Pay.UnmarshalJSON: values for iat and rvk must be between 0 and 2^53 - 1
	// {"rvk":9007199254740991}
}
