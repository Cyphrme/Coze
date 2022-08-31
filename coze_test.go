package coze

import (
	"encoding/json"
	"fmt"
	"strings"
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
	// {"alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"}
}

// ExamplePay_jsonMarshalCustom demonstrates marshalling Pay with a custom
// structure.
func ExamplePay_jsonMarshalCustom() {
	customStruct := CustomStruct{
		Msg: "Coze Rocks",
	}

	inputPay := Pay{
		Alg:    SEAlg(ES256),
		Iat:    1627518000, // Static for demonstration.  Use time.Time.Unix().
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
	// {"alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg","msg":"Coze Rocks"}
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
	// {"alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"}
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
	// {"alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg","msg":"Coze Rocks"}
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
		Iat:    1627518000, // Static for demonstration.  Use time.Time.Unix().
		Tmb:    MustDecode("cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk"),
		Typ:    "cyphr.me/msg",
		Struct: customStruct,
	}

	fmt.Println(inputPay)

	// Output:
	// {"alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg","msg":"Coze Rocks"}
}

func Example_dup() {
	data := `{"a": "b", "a":true,"c":["field_3 string 1","field3 string2"], "d": {"e": 1, "e": 2}}`
	err := checkDuplicate(json.NewDecoder(strings.NewReader(data)))
	if err == ErrDuplicate {
		fmt.Println("found a duplicate")
	} else if err != nil {
		panic(err)
	}

	// Recursive check
	data = `{"a": "b", "c":"d", "d": {"e": 1, "e": 2}}`
	err = checkDuplicate(json.NewDecoder(strings.NewReader(data)))
	if err == ErrDuplicate {
		fmt.Println("found a duplicate")
	} else if err != nil {
		panic(err)
	}

	data = `{"a": "b", "c":"d", "d": {"e": 1, "f": 2}}`
	err = checkDuplicate(json.NewDecoder(strings.NewReader(data)))
	if err == nil {
		fmt.Println("no duplicate")
	} else {
		panic(err)
	}

	// Output:
	// found a duplicate
	// found a duplicate
	// no duplicate
}

// Example demonstrating that unmarshalling a `pay` that has duplicate field
// names results in an error.
func ExamplePay_UnmarshalJSON_duplicate() {
	h := &Pay{}
	msg := []byte(`{"alg":"ES256","alg":"ES384"}`)

	err := json.Unmarshal(msg, h)
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// duplicate
}

// Example demonstrating that unmarshalling a `coze` that has duplicate field
// names results in an error.
func ExampleCoze_UnmarshalJSON_duplicate() {
	h := &Pay{}
	msg := []byte(`{"pay":"ES256","pay":"ES384"}`)

	err := json.Unmarshal(msg, h)
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// duplicate
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
	// {"name":"Bob","coze":{"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"}}
}

func ExampleCoze_String() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(GoldenCoze), cz)
	if err != nil {
		panic(err)
	}
	fmt.Println(cz)
	// Output:
	// {"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"}
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
	//{"can":["msg","alg","iat","tmb","typ"],"cad":"LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto","czd":"d0ygwQCGzuxqgUq1KsuAtJ8IBu0mkgAcKpUJzuX075M","pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"}
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
	err = cz.MetaWithAlg(0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", cz)

	// Output:
	//{"can":["msg","alg","iat","tmb","typ"],"cad":"LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto","czd":"d0ygwQCGzuxqgUq1KsuAtJ8IBu0mkgAcKpUJzuX075M","pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"}
	//{"can":["msg","alg","iat","tmb","typ"],"cad":"LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto","czd":"d0ygwQCGzuxqgUq1KsuAtJ8IBu0mkgAcKpUJzuX075M","pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"}
}

// ExampleCoze_jsonUnMarshal tests unmarshalling a coze.
func ExampleCoze_jsonUnMarshal() {
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
	//{"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"}
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
	//{"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"}
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
	// {
	//     "pay": {
	//         "msg": "Coze Rocks",
	//         "alg": "ES256",
	//         "iat": 1627518000,
	//         "tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
	//         "typ": "cyphr.me/msg"
	//     },
	//     "sig": "ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"
	// }
}

// Example_jsonRawMessageMarshal demonstrates using empty string, quote
// characters with no other content, and nil for json.RawMessage.  When using
// RawMessage, it should always be valid JSON or nil or otherwise it will result
// in an error.
func Example_jsonRawMessageMarshal() {
	o := json.RawMessage([]byte("")) // empty string
	anon := struct {
		Obj *json.RawMessage `json:"obj,omitempty"`
	}{
		Obj: &o,
	}

	// Incorrect usage with pointer to a zero value string.
	// Pointer to empty string will fail Marshal since an empty string is not
	// valid JSON. (If it were the value `""` it would pass. )
	b, err := Marshal(anon) // fails
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", b) // Empty because of error.

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
