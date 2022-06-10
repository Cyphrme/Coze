package coze

import (
	"encoding/json"
	"fmt"
	"log"
)

func ExampleMapSlice_MarshalJSON_unmarshal() {
	ms := MapSlice{
		MapItem{"abc", 123, 0},
		MapItem{"def", 456, 0},
		MapItem{"ghi", 789, 0},
	}

	b, err := Marshal(ms)
	if err != nil {
		log.Fatal(err)
	}

	ms = MapSlice{}
	if err := json.Unmarshal(b, &ms); err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
	fmt.Println(ms)
	// Output:
	// {"abc":123,"def":456,"ghi":789}
	// [{abc 123 1} {def 456 2} {ghi 789 3}]
}

func ExampleMapSlice_UnmarshalJSON() {

	var ms = MapSlice{}
	if err := json.Unmarshal([]byte(`{"abc":123,"def":456,"ghi":789}`), &ms); err != nil {
		log.Fatal(err)
	}

	fmt.Println(ms)
	// Output:
	// [{abc 123 1} {def 456 2} {ghi 789 3}]
}
