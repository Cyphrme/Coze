package coze

import "fmt"

func ExampleCanon_Append() {
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
	m := Merge(Canon{"a", "b"}, Canon{"c", "d"}, Canon{"e", "f"}, Canon(Only{"g", "h"}))
	fmt.Printf("%+v", m)

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

var az = []byte(`{"a":"a","z":"z"}`)
var ayz = []byte(`{"a":"a","y":"y","z":"z"}`)
var v bool

func ExampleIsNormal_nil() {
	fmt.Println("Nil")
	// Nil matches empty JSON, true.
	v = IsNormal([]byte(`{}`), nil)
	fmt.Println(v)

	// Nil Normal matches everything, true.
	v = IsNormal(az, nil)
	fmt.Println(v)

	// Output:
	// Nil
	// true
	// true
}

func ExampleIsNormal_canon() {
	fmt.Println("\nCanon")

	// Canon empty with empty records, true.
	v = IsNormal([]byte(`{}`), Canon{})
	fmt.Println(v)

	// Canon in order, Canon in order, ending nil with no record (variadic), true.
	v = IsNormal(az, Canon{"a"}, Canon{"z"}, nil)
	fmt.Println(v)

	// Canon in order, Canon in order, ending nil with record (variadic), true.
	v = IsNormal(ayz, Canon{"a"}, Canon{"y"}, nil)
	fmt.Println(v)

	// Canon in order, true.
	v = IsNormal(az, Canon{"a", "z"})
	fmt.Println(v)

	// Canon in order variadic, true.
	v = IsNormal(az, Canon{"a"}, Canon{"z"})
	fmt.Println(v)

	// Canon in order with Only in order (variadic), true.
	v = IsNormal(az, Canon{"a"}, Only{"z"})
	fmt.Println(v)

	// Canon in order with Extra (variadic), true.
	v = IsNormal(az, Canon{"a"}, Extra{})
	fmt.Println(v)

	// Canon in order with Option missing (variadic), true.
	v = IsNormal(az, Canon{"a", "z"}, Option{"b"})
	fmt.Println(v)

	// Canon with Extra (not present) and Canon (variadic), true.
	v = IsNormal(az, Canon{"a"}, Extra{}, Canon{"z"})
	fmt.Println(v)

	// Canon with Extra not present and Canon (variadic), true.
	v = IsNormal(ayz, Canon{"a"}, Extra{}, Canon{"y", "z"})
	fmt.Println(v)

	// Canon with Extra present and Canon (variadic), true.
	v = IsNormal(ayz, Canon{"a"}, Extra{}, Canon{"z"})
	fmt.Println(v)

	// Canon empty with records, false.
	v = IsNormal(az, Canon{})
	fmt.Println(v)

	// Canon in order, Canon in order, extra field, false.
	v = IsNormal(ayz, Canon{"a"}, Canon{"y"})
	fmt.Println(v)

	// Canon out of order, false.
	v = IsNormal(az, Canon{"z", "a"})
	fmt.Println(v)

	// Canon (correct) succeeded by extra (incorrect), false.
	v = IsNormal(az, Canon{"a"})
	fmt.Println(v)

	// Canon in order (correct) with Only missing (incorrect) (variadic), false.
	v = IsNormal(az, Canon{"a"}, Only{"b"})
	fmt.Println(v)

	// Canon amd Canon with extra field inbetween (variadic), false.
	v = IsNormal(ayz, Canon{"a"}, Canon{"z"})
	fmt.Println(v)

	// Canon with Extra (not present) and Canon and with succeeding extra (variadic), false.
	v = IsNormal(ayz, Canon{"a"}, Extra{}, Canon{"y"})
	fmt.Println(v)

	// Canon with extra (not present, incorrect) and Extra (variadic)(Checks for panic on out of bounds), false.
	v = IsNormal(az, Canon{"a", "z", "y"}, Extra{})
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
	v = IsNormal([]byte(`{}`), Only{})
	fmt.Println(v)

	// Only in order, true.
	v = IsNormal(az, Only{"a", "z"})
	fmt.Println(v)

	// Only out of order, true.
	v = IsNormal(az, Only{"z", "a"})
	fmt.Println(v)

	// Only in order variadic, true.
	v = IsNormal(az, Only{"a"}, Only{"z"})
	fmt.Println(v)

	// Only empty with records, false.
	v = IsNormal(az, Only{})
	fmt.Println(v)

	// Only with extra field, false.
	v = IsNormal(az, Only{"a", "y", "z"})
	fmt.Println(v)

	//Output:
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
	v = IsNormal([]byte(`{}`), Option{})
	fmt.Println(v)

	// Option with optional one field missing, true.
	v = IsNormal(az, Option{"a", "z", "x"})
	fmt.Println(v)

	// Two Options, true.
	v = IsNormal(az, Option{"a"}, Option{"z"})
	fmt.Println(v)

	// Three Options with last missing, true.
	v = IsNormal(az, Option{"a"}, Option{"z"}, Option{"x"}) //TODO
	fmt.Println(v)

	// Option with field missing and Extra, true.
	v = IsNormal(az, Option{"b"}, Extra{})
	fmt.Println(v)

	// Option (field missing) with canon present (variadic), true.
	v = IsNormal(ayz, Option{"b"}, Canon{"a", "y", "z"})
	fmt.Println(v)

	// Option in order with optional field missing and variadic, true.
	v = IsNormal(az, Option{"a"}, Option{"z", "x"})
	fmt.Println(v)

	// Option with canon present (variadic), true.
	v = IsNormal(ayz, Option{"a"}, Canon{"y", "z"})
	fmt.Println(v)

	// Option missing with Canon (variadic), true.
	v = IsNormal(az, Option{"b"}, Canon{"a", "z"})
	fmt.Println(v)

	// Need with option missing, true.
	v = IsNormal(ayz, Need{"a"}, Option{"b"})
	fmt.Println(v)

	// Option empty with records, false.
	v = IsNormal(az, Option{})
	fmt.Println(v)

	// Option with extra field, false.
	v = IsNormal(az, Option{"a"})
	fmt.Println(v)

	// Extra field then option, false.
	v = IsNormal(az, Option{"z"})
	fmt.Println(v)

	// Option out of order with optional field missing and variadic, false.
	v = IsNormal(az, Option{"z"}, Option{"x", "a"})
	fmt.Println(v)

	// Option with extra pay field, false.
	v = IsNormal(ayz, Option{"a", "y"})
	fmt.Println(v)

	// Need, option,then extra field, false.
	v = IsNormal(ayz, Need{"a"}, Option{"y"})
	fmt.Println(v)

	//Output:
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
	v = IsNormal([]byte(`{}`), Need{})
	fmt.Println(v)

	// Need empty with records, true.
	v = IsNormal(az, Need{})
	fmt.Println(v)

	// Need with extra second field, true.
	v = IsNormal(az, Need{"a"})
	fmt.Println(v)

	// Need with extra first field, true.
	v = IsNormal(az, Need{"z"})
	fmt.Println(v)

	// Need in order, true.
	v = IsNormal(az, Need{"a", "z"})
	fmt.Println(v)

	// Need out of order, true.
	v = IsNormal(az, Need{"a", "z"})
	fmt.Println(v)

	// Need in order with extra, true.
	v = IsNormal(ayz, Need{"a", "y"})
	fmt.Println(v)

	// Need in order with extra and Canon, true.
	v = IsNormal(ayz, Need{"a"}, Canon{"z"})
	fmt.Println(v)

	// Need with option present, true.
	v = IsNormal(ayz, Need{"a", "y"}, Option{"z"})
	fmt.Println(v)

	// Need, extra field, then option, true.
	v = IsNormal(ayz, Need{"a"}, Option{"z"})
	fmt.Println(v)

	// Need missing field, false.
	v = IsNormal(az, Need{"a", "y", "z"})
	fmt.Println(v)

	// Option present, Need repeated, false.
	v = IsNormal(az, Option{"a"}, Need{"a"})
	fmt.Println(v)

	// Need with and Canon and extra, false.
	v = IsNormal(ayz, Need{"a"}, Canon{"y"})
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
