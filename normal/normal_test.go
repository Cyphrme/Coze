package normal

import "fmt"

func ExampleAppend_canon() {
	fmt.Printf("%v\n", Append(Canon{"a", "b"}, Canon{"c", "d"}))

	// Output:
	// [a b c d]
}

func ExampleNormaler_Len() {
	fmt.Printf("%d %d %d %d %d %d\n",
		Canon{"a", "b"}.Len(),
		Only{"a", "b"}.Len(),
		Option{"a", "b"}.Len(),
		Need{"a", "b"}.Len(),
		Extra{"a", "b"}.Len(),
		Normaler(Canon{"a", "b"}).Len(),
	)

	// Output:
	// 2 2 2 2 2 2
}

func ExampleMerge() {
	fmt.Printf("%v\n", Merge(Canon{"a", "b"}, Canon{"c", "d"}, Canon{"e", "f"}))

	// When merging with Normals of different type, all type need to be the same
	// type.  The following casts Only as a Canon.
	fmt.Printf("%+v", Merge(Canon{"a", "b"}, Canon{"c", "d"}, Canon{"e", "f"}, Canon(Only{"g", "h"})))

	// Output:
	// [a b c d e f]
	// [a b c d e f g h]
}

func ExampleType() {
	fmt.Println(Type(Canon{}))
	fmt.Println(Type(Only{}))
	fmt.Println(Type(Option{}))
	fmt.Println(Type(Need{}))
	fmt.Println(Type(Extra{}))

	// Output:
	// canon
	// only
	// option
	// need
	// extra
}

var (
	az  = []byte(`{"a":"a","z":"z"}`)
	ayz = []byte(`{"a":"a","y":"y","z":"z"}`)
	v   bool
)

// ExampleIsNormal_nil shows the nil and zero case.
func ExampleIsNormal_nil() {
	v, _ = IsNormal(nil, nil)
	fmt.Println(v)

	// Nil matches empty JSON, true.
	v, _ = IsNormal([]byte(`{}`), nil)
	fmt.Println(v)

	// Nil Normal matches everything, true.
	v, _ = IsNormal(az, nil)
	fmt.Println(v)

	// Output:
	// false
	// true
	// true
}

func ExampleIsNormal_canon() {
	fmt.Println("\nCanon")

	// Canon empty with empty records, true.
	v, _ = IsNormal([]byte(`{}`), Canon{})
	fmt.Println(v)

	// Canon in order, Canon in order, ending nil with no record (variadic), true.
	v, _ = IsNormal(az, Canon{"a"}, Canon{"z"}, nil)
	fmt.Println(v)

	// Canon in order, Canon in order, ending nil with record (variadic), true.
	v, _ = IsNormal(ayz, Canon{"a"}, Canon{"y"}, nil)
	fmt.Println(v)

	// Canon in order, true.
	v, _ = IsNormal(az, Canon{"a", "z"})
	fmt.Println(v)

	// Canon in order variadic, true.
	v, _ = IsNormal(az, Canon{"a"}, Canon{"z"})
	fmt.Println(v)

	// Canon in order with Only in order (variadic), true.
	v, _ = IsNormal(az, Canon{"a"}, Only{"z"})
	fmt.Println(v)

	// Canon in order with Extra (variadic), true.
	v, _ = IsNormal(az, Canon{"a"}, Extra{})
	fmt.Println(v)

	// Canon in order with Option missing (variadic), true.
	v, _ = IsNormal(az, Canon{"a", "z"}, Option{"b"})
	fmt.Println(v)

	// Canon with Extra (not present) and Canon (variadic), true.
	v, _ = IsNormal(az, Canon{"a"}, Extra{}, Canon{"z"})
	fmt.Println(v)

	// Canon with Extra not present and Canon (variadic), true.
	v, _ = IsNormal(ayz, Canon{"a"}, Extra{}, Canon{"y", "z"})
	fmt.Println(v)

	// Canon with Extra present and Canon (variadic), true.
	v, _ = IsNormal(ayz, Canon{"a"}, Extra{}, Canon{"z"})
	fmt.Println(v)

	// Canon empty with records, false.
	v, _ = IsNormal(az, Canon{})
	fmt.Println(v)

	// Canon in order, Canon in order, extra field, false.
	v, _ = IsNormal(ayz, Canon{"a"}, Canon{"y"})
	fmt.Println(v)

	// Canon out of order, false.
	v, _ = IsNormal(az, Canon{"z", "a"})
	fmt.Println(v)

	// Canon (correct) succeeded by extra (incorrect), false.
	v, _ = IsNormal(az, Canon{"a"})
	fmt.Println(v)

	// Canon in order (correct) with Only missing (incorrect) (variadic), false.
	v, _ = IsNormal(az, Canon{"a"}, Only{"b"})
	fmt.Println(v)

	// Canon amd Canon with extra field inbetween (variadic), false.
	v, _ = IsNormal(ayz, Canon{"a"}, Canon{"z"})
	fmt.Println(v)

	// Canon with Extra (not present) and Canon and with succeeding extra (variadic), false.
	v, _ = IsNormal(ayz, Canon{"a"}, Extra{}, Canon{"y"})
	fmt.Println(v)

	// Canon with extra (not present, incorrect) and Extra (variadic)(Checks for panic on out of bounds), false.
	v, _ = IsNormal(az, Canon{"a", "z", "y"}, Extra{})
	fmt.Println(v)

	// Output:
	// Canon
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// false
	// false
	// false
	// false
	// false
	// false
	// false
	// false
}

func ExampleIsNormal_only() {
	fmt.Println("\nOnly")

	// Only empty with empty records, true.
	v, _ = IsNormal([]byte(`{}`), Only{})
	fmt.Println(v)

	// Only in order, true.
	v, _ = IsNormal(az, Only{"a", "z"})
	fmt.Println(v)

	// Only out of order, true.
	v, _ = IsNormal(az, Only{"z", "a"})
	fmt.Println(v)

	// Only in order variadic, true.
	v, _ = IsNormal(az, Only{"a"}, Only{"z"})
	fmt.Println(v)

	// Only empty with records, false.
	v, _ = IsNormal(az, Only{})
	fmt.Println(v)

	// Only with extra field, false.
	v, _ = IsNormal(az, Only{"a", "y", "z"})
	fmt.Println(v)

	// Output:
	// Only
	// true
	// true
	// true
	// true
	// false
	// false
}

func ExampleIsNormal_option() {
	fmt.Println("\nOption")

	// Option empty with empty records, true.
	v, _ = IsNormal([]byte(`{}`), Option{})
	fmt.Println(v)

	// Option with optional one field missing, true.
	v, _ = IsNormal(az, Option{"a", "z", "x"})
	fmt.Println(v)

	// Two Options, true.
	v, _ = IsNormal(az, Option{"a"}, Option{"z"})
	fmt.Println(v)

	// Three Options with last missing, true.
	v, _ = IsNormal(az, Option{"a"}, Option{"z"}, Option{"x"}) // TODO
	fmt.Println(v)

	// Option with field missing and Extra, true.
	v, _ = IsNormal(az, Option{"b"}, Extra{})
	fmt.Println(v)

	// Option (field missing) with canon present (variadic), true.
	v, _ = IsNormal(ayz, Option{"b"}, Canon{"a", "y", "z"})
	fmt.Println(v)

	// Option in order with optional field missing and variadic, true.
	v, _ = IsNormal(az, Option{"a"}, Option{"z", "x"})
	fmt.Println(v)

	// Option with canon present (variadic), true.
	v, _ = IsNormal(ayz, Option{"a"}, Canon{"y", "z"})
	fmt.Println(v)

	// Option missing with Canon (variadic), true.
	v, _ = IsNormal(az, Option{"b"}, Canon{"a", "z"})
	fmt.Println(v)

	// Need with option missing, true.
	v, _ = IsNormal(ayz, Need{"a"}, Option{"b"})
	fmt.Println(v)

	// Option empty with records, false.
	v, _ = IsNormal(az, Option{})
	fmt.Println(v)

	// Option with extra field, false.
	v, _ = IsNormal(az, Option{"a"})
	fmt.Println(v)

	// Extra field then option, false.
	v, _ = IsNormal(az, Option{"z"})
	fmt.Println(v)

	// Option out of order with optional field missing and variadic, false.
	v, _ = IsNormal(az, Option{"z"}, Option{"x", "a"})
	fmt.Println(v)

	// Option with extra pay field, false.
	v, _ = IsNormal(ayz, Option{"a", "y"})
	fmt.Println(v)

	// Need, option,then extra field, false.
	v, _ = IsNormal(ayz, Need{"a"}, Option{"y"})
	fmt.Println(v)

	// Output:
	// Option
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// false
	// false
	// false
	// false
	// false
	// false
}

func ExampleIsNormal_need() {
	fmt.Println("\nNeed")

	// Need empty with empty records, true.
	v, _ = IsNormal([]byte(`{}`), Need{})
	fmt.Println(v)

	// Need empty with records, true.
	v, _ = IsNormal(az, Need{})
	fmt.Println(v)

	// Need with extra second field, true.
	v, _ = IsNormal(az, Need{"a"})
	fmt.Println(v)

	// Need with extra first field, true.
	v, _ = IsNormal(az, Need{"z"})
	fmt.Println(v)

	// Need in order, true.
	v, _ = IsNormal(az, Need{"a", "z"})
	fmt.Println(v)

	// Need out of order, true.
	v, _ = IsNormal(az, Need{"a", "z"})
	fmt.Println(v)

	// Need in order with extra, true.
	v, _ = IsNormal(ayz, Need{"a", "y"})
	fmt.Println(v)

	// Need in order with extra and Canon, true.
	v, _ = IsNormal(ayz, Need{"a"}, Canon{"z"})
	fmt.Println(v)

	// Need with option present, true.
	v, _ = IsNormal(ayz, Need{"a", "y"}, Option{"z"})
	fmt.Println(v)

	// Need, extra field, then option, true.
	v, _ = IsNormal(ayz, Need{"a"}, Option{"z"})
	fmt.Println(v)

	// Need missing field, false.
	v, _ = IsNormal(az, Need{"a", "y", "z"})
	fmt.Println(v)

	// Option present, Need repeated, false.
	v, _ = IsNormal(az, Option{"a"}, Need{"a"})
	fmt.Println(v)

	// Need with and Canon and extra, false.
	v, _ = IsNormal(ayz, Need{"a"}, Canon{"y"})
	fmt.Println(v)

	// Output:
	// Need
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// false
	// false
	// false
}

func ExampleIsNormalUnchained() {
	v, _ = IsNormalUnchained(az, Need{"a"}, Need{"z"})
	fmt.Println(v)
	v, _ = IsNormalUnchained(ayz, Need{"a"}, Need{"z"}, Need{"y"})
	fmt.Println(v)
	v, _ = IsNormalUnchained(az, Need{"a"}, Option{"z"})
	fmt.Println(v)

	// Output:
	// true
	// true
	// false
}

func ExampleIsNormalNeedOption() {
	standard := `{
		"alg": "ES256",
		"iat": 1647357960,
		"tmb": "L0SS81e5QKSUSu-17LTQsvwKpUhBxe6ZZIEnSRV73o8",
		"typ": "cyphr.me/user/profile/update",
		`
	required := `"id": "L0SS81e5QKSUSu-17LTQsvwKpUhBxe6ZZIEnSRV73o8",`
	optional := `
	"city": "Pueblo",
	"country": "ISO 3166-2:US",
	"display_name": "Mr. Dev",
	"first_name": "Dev Test",
	"last_name": "1"
 }`
	need := Need{"alg", "iat", "tmb", "typ", "id"}
	option := Option{"display_name", "first_name", "last_name", "email", "address_1", "address_2", "phone_1", "phone_2", "city", "state", "zip", "country"}
	v, _ = IsNormalNeedOption([]byte(standard+required+optional), need, option)
	fmt.Println(v)
	// Missing required field 'id'.
	v, _ = IsNormalNeedOption([]byte(standard+optional), need, option)
	fmt.Println(v)

	// Output:
	// true
	// false
}
