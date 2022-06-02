package coze

import (
	"encoding/json"
	"fmt"
)

//ExampleHead_jsonUnmarshal tests unmarshalling a Head.
func ExampleHead_jsonUnmarshal() {
	h := &Head{}

	err := json.Unmarshal([]byte(Golden_Head), h)
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
