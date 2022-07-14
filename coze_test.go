package coze

import (
	"encoding/json"
	"fmt"
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
		Alg:    Golden_Key.Alg,
		Tmb:    Golden_Key.Tmb,
		Struct: &user,
	}

	coze, err := Golden_Key.SignPay(&pay)
	if err != nil {
		fmt.Println(err)
	}

	v, err := Golden_Key.VerifyCoze(coze)
	if err != nil {
		fmt.Println(err)
	}

	// Set sig to nil for deterministic printout
	coze.Sig = nil
	fmt.Println(v)
	fmt.Printf("%+v\n", coze)

	// Output:
	// true
	// {"pay":{"alg":"ES256","tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","DisplayName":"Coze","FirstName":"Foo","LastName":"Bar"}}
}

//ExamplePay_jsonUnmarshal tests unmarshalling a Pay.
func ExamplePay_jsonUnmarshal() {
	h := &Pay{}

	err := json.Unmarshal([]byte(Golden_Pay), h)
	if err != nil {
		fmt.Println(err)
	}

	out, err := Marshal(h)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s\n", out)
	// Output:
	// {"alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"}
}

//ExamplePay_jsonMarshalCustom demonstrates marshalling Pay with a custom
//structure.
func ExamplePay_jsonMarshalCustom() {
	var customStruct = CustomStruct{
		Msg: "Coze Rocks",
	}

	inputPay := Pay{
		Alg:    SEAlg(ES256),
		Iat:    1627518000, // Static for demonstration.  Use time.Time.Unix().
		Tmb:    MustDecode("cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk"),
		Typ:    "cyphr.me/msg",
		Struct: customStruct,
	}

	// May also call inputPay.MarshalJSON() instead.
	s, err := json.Marshal(&inputPay)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(s))

	// Output:
	// {"alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg","msg":"Coze Rocks"}
}

//ExamplePay_String_custom demonstrates fmt.Stringer on Pay with a custom
//structure.
func ExamplePay_String_custom() {
	var customStruct = CustomStruct{
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

// ExampleCoze_embed demonstrates how to embed a JSON `coze` into a third party
// JSON structure.
func ExampleCoze_embed() {
	type Outer struct {
		Name string
		Coze Coze // Embed a Coze into a larger, application defined JSON structure.
	}
	cz := new(Coze)
	err := json.Unmarshal([]byte(Golden_Coze), cz)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v", Outer{Name: "Bob", Coze: *cz})
	// Output:
	// {Name:Bob Coze:{"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"}}
}

func ExampleCoze_String() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(Golden_Coze), cz)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cz)
	// Output:
	// {"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"}
}

func ExampleCoze_Meta() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(Golden_Coze), cz)
	if err != nil {
		fmt.Println(err)
	}

	err = cz.Meta()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", cz)

	// Output:
	//{"can":["msg","alg","iat","tmb","typ"],"cad":"LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto","czd":"d0ygwQCGzuxqgUq1KsuAtJ8IBu0mkgAcKpUJzuX075M","pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"}
}

func ExampleCoze_MetaWithAlg() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(Golden_Coze), cz)
	if err != nil {
		fmt.Println(err)
	}

	// Test mismatch alg, which should error.
	err = cz.MetaWithAlg(SEAlg(ES224))
	if err == nil {
		fmt.Println("Test should error")
	}

	// Test with correct alg, no error.
	err = cz.MetaWithAlg(SEAlg(ES256))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", cz)

	// No alg given.
	err = cz.MetaWithAlg(0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", cz)

	// No alg given. // TODO
	// err = cz.MetaWithAlg(0)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("%s\n", cz)

	// Output:
	//{"can":["msg","alg","iat","tmb","typ"],"cad":"LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto","czd":"d0ygwQCGzuxqgUq1KsuAtJ8IBu0mkgAcKpUJzuX075M","pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"}
	//{"can":["msg","alg","iat","tmb","typ"],"cad":"LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto","czd":"d0ygwQCGzuxqgUq1KsuAtJ8IBu0mkgAcKpUJzuX075M","pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"}
}

//ExampleCoze_jsonUnMarshal tests unmarshalling a coze.
func ExampleCoze_jsonUnMarshal() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(Golden_Coze), cz)
	if err != nil {
		fmt.Println(err)
	}

	// remarshal for comparison
	b, err := Marshal(cz)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
	// Output:
	//{"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"}
}

func ExampleCoze_jsonMarshal() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(Golden_Coze), cz)
	if err != nil {
		fmt.Println(err)
	}

	b, err := Marshal(cz)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+s\n", b)
	// Output:
	//{"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"}
}

func ExampleCoze_jsonMarshalPretty() {
	cz := new(Coze)
	err := json.Unmarshal([]byte(Golden_Coze), cz)
	if err != nil {
		fmt.Println(err)
	}

	b, err := MarshalPretty(cz)
	if err != nil {
		fmt.Println(err)
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

	// Pointer to empty string will fail Marshal since an empty string is not
	// valid JSON. (If it were the value `""`, it would pass. )
	// Incorrect usage with pointer to a zero value string.
	b, err := Marshal(anon)
	if err != nil {
		fmt.Println(err)
	}

	// Correct usage with quotes characters.
	quotes := []byte("\"\"") // string with two quote characters
	anon.Obj = (*json.RawMessage)(&quotes)
	b, err = Marshal(anon)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", b)

	// Correct usage with with `nil`, which prints as the JSON "null".
	o = nil
	anon.Obj = &o
	b, err = Marshal(anon)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", b)

	// Output:
	// json: error calling MarshalJSON for type *json.RawMessage: unexpected end of JSON input
	// {"obj":""}
	// {"obj":null}
}
