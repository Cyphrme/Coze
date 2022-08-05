package coze

import (
	"fmt"
	"testing"
)

func ExampleCanon() {
	b := []byte(`{"z":"z", "a":"a"}`)

	can, err := GetCanon(b)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(can)
	// Output: [z a]
}

// ExampleCanonicalHash. See also Example_genCad
func ExampleCanonicalHash() {
	canon := []string{"alg", "iat", "msg", "tmb", "typ"}
	cad, err := CanonicalHash([]byte(GoldenPay), canon, SHA256)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cad.String())

	// Without canon
	cad, err = CanonicalHash([]byte(GoldenPay), nil, SHA256)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cad.String())

	// Output:
	// aC2YKfNvovfnZOw_RVxSEW6NeaUq41DZXX0oeaOboRg
	// LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto
}

// Example CanonicalHash for all hashing algos.
func ExampleCanonicalHash_permutations() {
	canon := []string{"alg", "iat", "msg", "tmb", "typ"}
	algs := []string{"SHA-224", "SHA-256", "SHA-384", "SHA-512", "SHA3-224", "SHA3-256", "SHA3-384", "SHA3-512", "SHAKE128", "SHAKE256"}
	for _, alg := range algs {
		cad, err := CanonicalHash([]byte(GoldenPay), canon, ParseHashAlg(alg))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(cad.String())
	}

	// Output:
	// nrGQqKYvFKeVDFlOMIusP2A2AWn4DX-2XLNfJA
	// aC2YKfNvovfnZOw_RVxSEW6NeaUq41DZXX0oeaOboRg
	// KC5UHPxzNl567oONOphhWY6cHhuXoSiyOLiNGSmTIcvA8XDtWQf-fr4xNPCLfzCo
	// oX2NMgJ_QRW9rf59N5VOSMILg6mzVHld5CqRaOatLCbRWVRh1Y6Rq4tRZGzZNvNKEM0qbYBlWk6y9BcnuRzczA
	// GAYOBAxW2x7MvHVYpRLnjX3rUKcuhvDOCVVK3Q
	// UfFl2lw4KHc2-0GX-mnqtfpScM1Qf7L_IaTGojR6_Go
	// pVo43tSAG8apVs26QLOFG0Cbh3ScrbHd_VGjaFAIQtlCLiXcsgdmsGwOyXoK4zBz
	// IA8Xv6tt32B49THtWOzN9AyKtnG5a0v93DSF4IShHsT6S2lWKQl1H2yuyMAYocVKBkMF5dp0miKB58NXROqAMg
	// muWDwpDGlR-jwGPOQlj6A6B5FYA_U5nFq2KtwV8B-Uw
	// QPfIPjKmLO4qLmiClA6GjYQKBO6MI2wBZUhi9uVTVr0WGP3LgOTQRup6l5Caxz6GtiUnNeQe6JMdVSvhdvLW-Q
}

// ExampleCanonical.
func ExampleCanonical() {
	var b []byte
	var err error

	type ABC struct {
		A string `json:"a"`
		B string `json:"b,omitempty"`
		C string `json:"c"`
	}

	// []byte (out of order) with nil canon.  Missing field  "b" should be omitted
	// from output.
	ca, err := Marshal(map[string]string{"c": "c", "a": "a"})
	if err != nil {
		fmt.Println(err)
	}

	b, err = Canonical(ca, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Input: []byte; Canon: nil    => %s\n", b)

	// []byte (out of order) with struct (in order) canon. Missing field "b"
	// should be omitted from output.
	ca, err = Marshal(map[string]string{"c": "c", "a": "a"})
	if err != nil {
		fmt.Println(err)
	}
	b, err = Canonical(ca, new(ABC))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Input: []byte; Canon: struct => %s\n", b)

	// []byte (out of order) with struct (in order) canon.
	byteJSON := []byte(`{"c":"c", "a": "a"}`)
	b, err = Canonical(byteJSON, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Input: []byte; Canon: nil    => %s\n", b)

	// Output:
	// Input: []byte; Canon: nil    => {"a":"a","c":"c"}
	// Input: []byte; Canon: struct => {"a":"a","c":"c"}
	// Input: []byte; Canon: nil    => {"c":"c","a":"a"}
}

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

// func ExampleUnion() {
// 	fmt.Printf("%v\n", Union(0, Canon{"a", "b"}, Canon{"c", "d"}, Canon{"e", "f"}))
// 	fmt.Printf("%v\n", Union[[]Normal](0, Canon{"a", "b"}, Canon{"c", "d"}, Canon{"e", "f"}, Only{"g", "h"}))
// 	// Output:
// 	// [a b c d e f]
// }

//func ExampleIsNormal() {
func TestIsNormal(t *testing.T) {
	az := []byte(`{"a":"a","z":"z"}`)
	ayz := []byte(`{"a":"a","y":"y","z":"z"}`)
	_ = ayz
	_ = az
	var v bool

	fmt.Println("Nil")
	// Nil matches empty JSON, true.
	v = IsNormal([]byte(`{}`), nil)
	fmt.Println(v)

	// Nil Normal matches everything, true.
	v = IsNormal(az, nil)
	fmt.Println(v)

	////////////////////
	// Canon
	////////////////////
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

	////////////////////
	// Only
	////////////////////
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

	////////////////////
	// Option
	////////////////////
	fmt.Println("\nOption")

	// Option empty with empty records, true.
	v = IsNormal([]byte(`{}`), Option{})
	fmt.Println(v)

	// Option with optional one field missing, true.
	v = IsNormal(az, Option{"a", "z", "x"})
	fmt.Println(v)

	// Option with field missing and extra, true.
	v = IsNormal(az, Option{"b"}, Extra{})
	fmt.Println(v)

	// Option (field missing) with canon present (variadic), true.
	v = IsNormal(ayz, Option{"b"}, Canon{"a", "y", "z"})
	fmt.Println(v)

	// Option in order with optional field missing and variadic, true.
	v = IsNormal(az, Option{"a"}, Option{"z", "x"})
	fmt.Println(v)

	// Option (field missing) with canon present (variadic), true.
	v = IsNormal(ayz, Option{"a"}, Canon{"y", "z"})
	fmt.Println(v)

	// Option empty with records, false. // TODO
	v = IsNormal(az, Option{})
	fmt.Println(v)

	// Option out of order with optional field missing and variadic, false. // TODO
	v = IsNormal(az, Option{"z"}, Option{"x", "a"})
	fmt.Println(v)

	// Option with extra pay field, false.
	v = IsNormal(ayz, Option{"a", "y"})
	fmt.Println(v)

	////////////////////
	// Need
	////////////////////
	fmt.Println("\nNeed")

	// Need empty with empty records, true.
	v = IsNormal([]byte(`{}`), Need{})
	fmt.Println(v)

	// Need empty with records, true.
	v = IsNormal(az, Need{})
	fmt.Println(v)

	// Need in of order, true.
	v = IsNormal(az, Need{"a", "z"})
	fmt.Println(v)

	// Need out of order, true.
	v = IsNormal(az, Need{"a", "z"})
	fmt.Println(v)

	// Need with option missing, true.
	v = IsNormal(az, Need{"a", "z"}, Option{"y"})
	fmt.Println(v)

	// Need with option present, true.
	v = IsNormal(ayz, Need{"a", "z"}, Option{"y"})
	fmt.Println(v)

	// Need missing field, false.
	v = IsNormal(az, Need{"a", "y", "z"})
	fmt.Println(v)

	//Need, extra field, then option, false.
	v = IsNormal(ayz, Need{"a"}, Option{"z"})
	fmt.Println(v)

	// Need, option,then extra field, false.
	v = IsNormal(ayz, Need{"a"}, Option{"y"})
	fmt.Println(v)

	// Output:
	// Nil
	// true
	// true
	//
	// Canon
	// true
	// true
	// true
	// true
	// true
	// true
	// false
	// true
	// true
	// true
	// false
	// true
	// false
	// true
	// true
	// false
	// true
	// false
	//
	// Only
	// true
	// true
	// false
	// true
	//
	// Option
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	//
	// Need
	// true
	// true
	// false
	// true
	// true
	// true
	// true
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
