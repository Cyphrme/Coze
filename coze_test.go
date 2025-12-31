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
	// {"pay":{"alg":"ES256","tmb":"U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg","DisplayName":"Coze","FirstName":"Foo","LastName":"Bar"}}
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
	// {"alg":"ES256","now":1623132000,"tmb":"U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg","typ":"cyphr.me/msg"}
}

// ExamplePay_jsonMarshalCustom demonstrates marshalling Pay with a custom
// structure.
func ExamplePay_jsonMarshalCustom() {
	customStruct := CustomStruct{
		Msg: "Coze Rocks",
	}

	inputPay := Pay{
		Alg:    SEAlg(ES256),
		Now:    1623132000, // Static for demonstration.  Use time.Now().Unix().
		Tmb:    MustDecode("U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg"),
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
	// {"alg":"ES256","now":1623132000,"tmb":"U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg","typ":"cyphr.me/msg","msg":"Coze Rocks"}
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
	// {"alg":"ES256","now":1623132000,"tmb":"U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg","typ":"cyphr.me/msg"}
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
	// {"alg":"ES256","now":1623132000,"tmb":"U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg","typ":"cyphr.me/msg","msg":"Coze Rocks"}
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
		Now:    1623132000, // Static for demonstration.  Use time.Now().Unix().
		Tmb:    MustDecode("U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg"),
		Typ:    "cyphr.me/msg",
		Struct: customStruct,
	}
	fmt.Println(inputPay)

	// Output:
	// {"alg":"ES256","now":1623132000,"tmb":"U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg","typ":"cyphr.me/msg","msg":"Coze Rocks"}
}

// Example demonstrating that unmarshalling a `pay` that has duplicate field
// names results in an error.
func ExamplePay_UnmarshalJSON_duplicate() {
	h := &Pay{}
	msg := []byte(`{"alg":"ES256","alg":"ES384"}`)
	err := json.Unmarshal(msg, h)
	fmt.Println(err)

	// Output:
	// Coze: JSON duplicate field "alg"
}

// Example demonstrating that unmarshalling a `coze` that has duplicate field
// names results in an error.
func ExampleCoze_UnmarshalJSON_duplicate() {
	h := &Pay{}
	msg := []byte(`{"coze":{"pay":"ES256","pay":"ES384"}}`)
	err := json.Unmarshal(msg, h)
	fmt.Println(err)

	// Output:
	// Coze: JSON duplicate field "pay"
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
	// {"name":"Bob","coze":{"pay":{"msg":"Coze Rocks","alg":"ES256","now":1623132000,"tmb":"U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg","typ":"cyphr.me/msg"},"sig":"bbO49APro9TGzAxDWvyT0a41l2sEFMpYWqC-hvDlJukyXKZ_0TRNsrJNcTIso3b8kh5wbLL2KLvOO4zfsHplwA"}}
}

func ExampleCoze_String() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(GoldenCoze), cz)
	if err != nil {
		panic(err)
	}
	fmt.Println(cz)

	// Output:
	// {"pay":{"msg":"Coze Rocks","alg":"ES256","now":1623132000,"tmb":"U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg","typ":"cyphr.me/msg"},"sig":"bbO49APro9TGzAxDWvyT0a41l2sEFMpYWqC-hvDlJukyXKZ_0TRNsrJNcTIso3b8kh5wbLL2KLvOO4zfsHplwA"}
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
	//{"pay":{"msg":"Coze Rocks","alg":"ES256","now":1623132000,"tmb":"U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg","typ":"cyphr.me/msg"},"can":["msg","alg","now","tmb","typ"],"cad":"AyVZoWUv_rJf7_KqoeRS5odr8g3MZwBzhtBdSZderxk","sig":"bbO49APro9TGzAxDWvyT0a41l2sEFMpYWqC-hvDlJukyXKZ_0TRNsrJNcTIso3b8kh5wbLL2KLvOO4zfsHplwA","czd":"DHEHV1BZPYMMzZs2auqF5vlvCySOdiOWdPleWHy3Ypg"}
}

func ExampleCoze_MetaWithAlg() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(GoldenCoze), cz)
	if err != nil {
		panic(err)
	}

	// coze.pay.alg given and parameter alg given.
	err = cz.MetaWithAlg(SEAlg(ES256))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", cz)

	// coze.pay.alg given and parameter alg not given.  (Alg is parsed from pay).
	err = cz.MetaWithAlg("")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", cz)

	// Test mismatch alg, which must error.
	err = cz.MetaWithAlg(SEAlg(ES224))
	if err == nil {
		fmt.Println("Test must error")
	}

	// Test no coze.pay.alg or alg given, which must error.
	// Will error with Hash: invalid HshAlg "UnknownHshAlg")
	cz2 := new(Coze)
	err = json.Unmarshal(GoldenCozeNoAlg, cz2)
	if err != nil {
		panic(err)
	}
	err = cz2.Meta()
	if err == nil {
		fmt.Println("Test must error")
	}

	// Test no coze.pay.alg but alg is given (contextual coze)
	err = cz2.MetaWithAlg(SEAlg(ES256))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", cz2)

	// Test no coze.pay.alg or coze.sig, so czd should not be calculated
	cz3 := new(Coze)
	err = json.Unmarshal(GoldenPayNoAlg, &cz3.Pay)
	if err != nil {
		panic(err)
	}
	err = cz3.MetaWithAlg(SEAlg(ES256))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", cz3)

	// Output:
	// {"pay":{"msg":"Coze Rocks","alg":"ES256","now":1623132000,"tmb":"U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg","typ":"cyphr.me/msg"},"can":["msg","alg","now","tmb","typ"],"cad":"AyVZoWUv_rJf7_KqoeRS5odr8g3MZwBzhtBdSZderxk","sig":"bbO49APro9TGzAxDWvyT0a41l2sEFMpYWqC-hvDlJukyXKZ_0TRNsrJNcTIso3b8kh5wbLL2KLvOO4zfsHplwA","czd":"DHEHV1BZPYMMzZs2auqF5vlvCySOdiOWdPleWHy3Ypg"}
	// {"pay":{"msg":"Coze Rocks","alg":"ES256","now":1623132000,"tmb":"U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg","typ":"cyphr.me/msg"},"can":["msg","alg","now","tmb","typ"],"cad":"AyVZoWUv_rJf7_KqoeRS5odr8g3MZwBzhtBdSZderxk","sig":"bbO49APro9TGzAxDWvyT0a41l2sEFMpYWqC-hvDlJukyXKZ_0TRNsrJNcTIso3b8kh5wbLL2KLvOO4zfsHplwA","czd":"DHEHV1BZPYMMzZs2auqF5vlvCySOdiOWdPleWHy3Ypg"}
	// {"pay":{"msg":"Coze Rocks","now":1623132000,"tmb":"U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg","typ":"cyphr.me/msg"},"can":["msg","now","tmb","typ"],"cad":"ywD6jd0t00XECuIG873VWZKXrobsAQz9tHTT8_tOHtg","sig":"rVVX9Px9ZVdU-YQdWTHK-hrgjQZVngztqJq7QlPBw1o9XUhN7GzWRV_0u2s-gP7Z9MHRCicq9j7InhUrg8LNjg","czd":"aDGGb6KlUhcMdufV3lyUxmg9MBbDfe1SvANU5fAME8c"}
	// {"pay":{"msg":"Coze Rocks","now":1623132000,"tmb":"U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg","typ":"cyphr.me/msg"},"can":["msg","now","tmb","typ"],"cad":"ywD6jd0t00XECuIG873VWZKXrobsAQz9tHTT8_tOHtg"}
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
	// {"pay":{"msg":"Coze Rocks","alg":"ES256","now":1623132000,"tmb":"U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg","typ":"cyphr.me/msg"},"can":["msg","alg","now","tmb","typ"],"cad":"AyVZoWUv_rJf7_KqoeRS5odr8g3MZwBzhtBdSZderxk"}
	// {"pay":{},"cad":"RBNvo1WzZ4oRRq0W9-hknpT7T8If536DEMBg9hyq_4o","sig":"UG0KP-cElD3mPoN8LRVd4_uoNzMwmpUm3pKxt-iy6So8f1JxmxMcO9JFzsmecFXyt5PjsOTZdUKyV6eZRNl-hg","czd":"nib9RLKirNz50PA2Sv6uZnA03_wdMmA1dyAoUi0OhVY"}
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
	//{"pay":{"msg":"Coze Rocks","alg":"ES256","now":1623132000,"tmb":"U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg","typ":"cyphr.me/msg"},"sig":"bbO49APro9TGzAxDWvyT0a41l2sEFMpYWqC-hvDlJukyXKZ_0TRNsrJNcTIso3b8kh5wbLL2KLvOO4zfsHplwA"}
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
	//{"pay":{"msg":"Coze Rocks","alg":"ES256","now":1623132000,"tmb":"U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg","typ":"cyphr.me/msg"},"sig":"bbO49APro9TGzAxDWvyT0a41l2sEFMpYWqC-hvDlJukyXKZ_0TRNsrJNcTIso3b8kh5wbLL2KLvOO4zfsHplwA"}
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
	//         "now": 1623132000,
	//         "tmb": "U5XUZots-WmQYcQWmsO751Xk0yeVi9XUKWQ2mGz6Aqg",
	//         "typ": "cyphr.me/msg"
	//     },
	//     "sig": "bbO49APro9TGzAxDWvyT0a41l2sEFMpYWqC-hvDlJukyXKZ_0TRNsrJNcTIso3b8kh5wbLL2KLvOO4zfsHplwA"
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
	// Happy path; no duplicate.  Should not error.
	data := `{"a": "b", "c":"d", "d": {"e": 1, "f": 2}}`
	err := checkDuplicate(json.NewDecoder(strings.NewReader(data)))
	if err != nil {
		t.Fatal(err)
	}

	// Duplicate, should error.
	data = `{"a": "aValue", "a":true,"c":["field_3 string 1","field3 string2"], "d": {"e": 1, "e": 2}}`
	err = checkDuplicate(json.NewDecoder(strings.NewReader(data)))
	if err == nil {
		t.Fatal("Should have found duplicate.")
	}

	// Recursive check with duplicate in inner struct.  Should error.
	data = `{"a": "aValue", "c":"cValue", "d": {"e": 1, "e": 2}}`
	err = checkDuplicate(json.NewDecoder(strings.NewReader(data)))
	if err == nil {
		t.Fatal("Recursive check should have found duplicate.")
	}
}

// Demonstrates expectations for values that are non-integer, negative, or too
// large for rvk.
func Example_now_rvk_too_big() {
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
	// Pay.UnmarshalJSON: values for now and rvk must be between 0 and 2^53 - 1
	// Pay.UnmarshalJSON: values for now and rvk must be between 0 and 2^53 - 1
	// {"rvk":9007199254740991}
}
