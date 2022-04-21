package coze

import (
	"encoding/json"
	"fmt"
)

const Golden_Head_String = `{
	"alg": "ES256",
	"iat": 1623132000,
	"tmb": "0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD",
	"typ": "cyphr.me"
	}`

//ExampleCyUnmarshal tests unmarshalling a `cy`.
func ExampleHead_jsonUnmarshal() {
	h := &Head{}

	err := json.Unmarshal([]byte(Golden_Head_String), h)
	if err != nil {
		fmt.Println(err)
	}

	out, err := Marshal(h)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s\n", out)
	// Output:
	// {"alg":"ES256","iat":1623132000,"tmb":"0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD","typ":"cyphr.me"}
}

// Example_jsonMarshalRaw demonstrates using nil for RawMessage.  RawMessage
// should always be valid JSON as it is marshaled as is.
func Example_jsonMarshalRaw() {
	o := json.RawMessage([]byte(""))
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
	quotes := []byte("\"\"")
	anon.Obj = (*json.RawMessage)(&quotes)
	b, err = Marshal(anon)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", b)

	// Correct usage with with `nil`.
	o = nil
	anon.Obj = &o
	b, err = Marshal(anon)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", b)

	// Correct usage with with `nil`.
	o = nil
	anon.Obj = &o
	b, err = json.Marshal(anon)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", b)

	// Output:
	// json: error calling MarshalJSON for type *json.RawMessage: unexpected end of JSON input
	// {"obj":""}
	// {"obj":null}
	// {"obj":null}
}
