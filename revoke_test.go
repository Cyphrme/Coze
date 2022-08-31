package coze

import (
	"encoding/json"
	"fmt"
)

func ExampleRevoke_MarshalJSON() {
	r := Revoke{
		Rvk: 1,
	}
	fmt.Printf("%s\n", r)

	r = Revoke{
		Rvk: 1,
		Msg: "Test",
	}
	fmt.Printf("%s\n", r)

	p := Pay{
		Iat: 1627518000,
	}
	r = Revoke{
		Rvk: 1,
		Msg: "Test",
		Pay: p,
	}
	fmt.Printf("%s\n", r)

	// Output:
	// {"rvk":1}
	// {"rvk":1,"msg":"Test"}
	// {"rvk":1,"msg":"Test","iat":1627518000}
}

func ExampleRevoke_unmarshalJSON() {
	r := new(Revoke)

	err := json.Unmarshal([]byte(`{"rvk":1}`), r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", r)

	err = json.Unmarshal([]byte(`{"rvk":1,"msg":"Test","iat":1627518000}`), r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", r)

	// Output:
	// {"rvk":1}
	// {"rvk":1,"msg":"Test","iat":1627518000}
}

func ExampleKey_Revoke() {
	gk2 := GoldenKey // Make a copy
	fmt.Println(gk2.IsRevoked())
	coze, err := gk2.Revoke("Posted my private key online.")
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%s\n", coze.Pay)

	// Both the revoke coze and the key should be interpreted as revoked.
	fmt.Println(IsRevoked(coze.Pay))
	fmt.Println(gk2.IsRevoked())

	// Manually set rvk to 1 (Revoked)
	gk2.Rvk = 1
	fmt.Println(gk2.IsRevoked())

	// Manually set rvk to 0 (Not revoked)
	gk2.Rvk = 0
	fmt.Println(gk2.IsRevoked())

	// Output:
	// false
	// true
	// true
	// true
	// false
}

func ExampleIsRevoked() {
	gk2 := GoldenKey // Make a copy
	fmt.Println(IsRevoked([]byte(gk2.String())))
	coze, err := gk2.Revoke("Posted my private key online.")
	if err != nil {
		panic(err)
	}

	// Both the revoke coze and the key should be interpreted as revoked.
	fmt.Println(IsRevoked(coze.Pay))
	fmt.Println(IsRevoked([]byte(gk2.String())))

	// Output:
	// false
	// true
	// true
}