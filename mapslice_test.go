package coze

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestMarshal(t *testing.T) {
	ms := MapSlice{
		MapItem{Key: "abc", Value: 123},
		MapItem{Key: "def", Value: 456},
		MapItem{Key: "ghi", Value: 789},
	}

	b, err := json.Marshal(ms)
	if err != nil {
		t.Fatal(err)
	}

	e := "{\"abc\":123,\"def\":456,\"ghi\":789}"
	r := string(b)

	if r != e {
		t.Errorf("expected: %s\ngot: %s", e, r)
	}
}

func TestMarshalError(t *testing.T) {
	ms := MapSlice{
		MapItem{Key: "abc", Value: make(chan int)},
	}

	e := "json: error calling MarshalJSON for type coze.MapSlice: json: unsupported type: chan int"
	if _, err := json.Marshal(ms); err != nil && e != err.Error() {
		t.Errorf("expected: %s\ngot: %v", e, err)
	}
}

func TestUnmarshal(t *testing.T) {
	ms := MapSlice{}
	if err := json.Unmarshal([]byte("{\"abc\":123,\"def\":456,\"ghi\":789}"), &ms); err != nil {
		t.Fatal(err)
	}

	e := "[{abc 123} {def 456} {ghi 789}]"
	r := fmt.Sprintf("%v", ms)

	if r != e {
		t.Errorf("expected: %s\ngot: %s", e, r)
	}
}

func ExampleMapSlice_MarshalJSON() {
	ms := MapSlice{
		MapItem{"abc", 123, 0},
		MapItem{"def", 456, 0},
		MapItem{"ghi", 789, 0},
	}

	b, err := Marshal(ms)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", b)

	// Output:
	// {"abc":123,"def":456,"ghi":789}
}

func ExampleMapSlice_UnmarshalJSON() {
	ms := MapSlice{}

	if err := json.Unmarshal([]byte(`{"abc":123,"def":456,"ghi":789}`), &ms); err != nil {
		log.Fatal(err)
	}

	fmt.Println(ms)
	// Output:
	// [{abc 123} {def 456} {ghi 789}]
}
