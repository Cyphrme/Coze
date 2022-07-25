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

// func ExampleIsNormal() {
// 	b := []byte(`{"z":"z","a":"a"}`)
// 	by := []byte(`{"z":"z","y":"y","a":"a"}`)

// 	////////////////////
// 	// Canon
// 	////////////////////
// 	fmt.Println("Canon")
// 	// In order, pass
// 	canon := []string{"z", "a"}
// 	// Canon in order, pass
// 	v := IsNormal(b, Canon(canon), nil)
// 	fmt.Println(v)

// 	// Canon Out of order, fail
// 	canon = []string{"a", "z"}
// 	v = IsNormal(b, Canon(canon), nil)
// 	fmt.Println(v)

// 	////////////////////
// 	// Only
// 	////////////////////
// 	fmt.Println("\nOnly")
// 	// Only with extra field, fail
// 	only := []string{"a", "y", "z"}
// 	v = IsNormal(b, Only(only), nil)
// 	fmt.Println(v)

// 	// Only Out of order, pass
// 	only = []string{"a", "z"}
// 	v = IsNormal(b, Only(only), nil)
// 	fmt.Println(v)

// 	// In order only , pass
// 	only = []string{"z", "a"}
// 	v = IsNormal(b, Only(only), nil)
// 	fmt.Println(v)

// 	////////////////////
// 	// Need
// 	////////////////////
// 	fmt.Println("\nNeed")
// 	// Need missing field, fail
// 	need := []string{"a", "y", "z"}
// 	v = IsNormal(b, Need(need), nil)
// 	fmt.Println(v)

// 	// Need out of order, pass
// 	need = []string{"a", "z"}
// 	v = IsNormal(b, Need(need), nil)
// 	fmt.Println(v)

// 	// Need with option missing, pass
// 	need = []string{"a", "z"}
// 	opt := []string{"y"}
// 	v = IsNormal(b, Need(need), Option(opt))
// 	fmt.Println(v)

// 	// Need with option present, pass
// 	v = IsNormal(by, Need(need), Option(opt))
// 	fmt.Println(v)

// 	// Need with option and extra field, fail
// 	need = []string{"a"}
// 	opt = []string{"z"}
// 	v = IsNormal(by, Need(need), opt)
// 	fmt.Println(v)

// 	// Breaking test
// 	pay := `{"alg":"ES256","iat":1658248498,"tmb":"L0SS81e5QKSUSu-17LTQsvwKpUhBxe6ZZIEnSRV73o8","typ":"cyphr.me/comment/update","id":"G1VvNnwglup8HYNYlSgLwgnTha0dYWkVta8x7HzoPg8","root":"LKMACi2SyJNR3rLYWe9jtW1EzNk6T1cnj_NBlTooFXc","text":"Comment updated","edit":true}`
// 	need = []string{"alg", "iat", "tmb", "typ", "id", "root"}
// 	opt = []string{"text", "counterfeit_report", "rating", "rating_denom", "purchase_location", "purchased", "edit"}
// 	v = IsNormal([]byte(pay), Need(need), Option(opt))
// 	fmt.Println(v)

// 	// ////////////////////
// 	// // Order
// 	// ////////////////////
// 	// fmt.Println("\nOrder")
// 	// // Order missing field, fail
// 	// order := []string{"z", "a", "y"}
// 	// v = IsNormal(b, Order(order), nil)
// 	// fmt.Println(v)

// 	// // Order out of order, fail
// 	// order = []string{"a", "z"}
// 	// v = IsNormal(b, Order(order), nil)
// 	// fmt.Println(v)

// 	// // Order extra field, pass
// 	// order = []string{"z", "y"}
// 	// v = IsNormal(by, Order(order), nil)
// 	// fmt.Println(v)

// 	// // Order with option missing, pass
// 	// order = []string{"z", "a"}
// 	// opt = []string{"y"}
// 	// v = IsNormal(b, Order(order), Option(opt))
// 	// fmt.Println(v)

// 	// // Order with option present, pass
// 	// order = []string{"z", "y"}
// 	// opt = []string{"a"}
// 	// v = IsNormal(by, Order(order), Option(opt))
// 	// fmt.Println(v)

// 	// // Order with option and extra field in the middle, fail
// 	// order = []string{"a"}
// 	// opt = []string{"y"}
// 	// v = IsNormal(by, Order(order), opt)
// 	// fmt.Println(v)

// 	// // Order with option and extra field at the end, fail
// 	// order = []string{"a"}
// 	// opt = []string{"z"}
// 	// v = IsNormal(by, Order(order), opt)
// 	// fmt.Println(v)

// 	////////////////////
// 	// Option
// 	////////////////////
// 	fmt.Println("\nOption")
// 	// Option with one missing, pass
// 	option := []string{"z", "a", "y"}
// 	v = IsNormal(b, Option(option), nil)
// 	fmt.Println(v)

// 	// Option with one missing, fail
// 	option = []string{"z", "y"}
// 	v = IsNormal(by, Option(option), nil)
// 	fmt.Println(v)

// 	//Output:
// 	// Canon
// 	// true
// 	// false
// 	//
// 	// Only
// 	// false
// 	// true
// 	// true
// 	//
// 	// Need
// 	// false
// 	// true
// 	// true
// 	// true
// 	// false
// 	//
// 	// Order
// 	// false
// 	// false
// 	// true
// 	// true
// 	// true
// 	// false
// 	// false
// 	//
// 	// Option
// 	// true
// 	// false
// }

//func ExampleIsNormal() {
func TestIsNormal(t *testing.T) {
	az := []byte(`{"a":"a","z":"z"}`)
	//ayz := []byte(`{"a":"a","y":"y","z":"z"}`)
	var v bool

	// // Nil case matches everything, true.
	// v = IsNormal(az, Canon{})
	// fmt.Println(v)

	// Test if ending nil is valid.
	// Canon in order variadic with ending nil, true.
	v = IsNormal(az, Canon{"a"}, Canon{"z"}, nil)
	fmt.Println(v)

	////////////////////
	// Canon
	////////////////////
	fmt.Println("\nCanon")

	// Canon in order, true.
	v = IsNormal(az, Canon{"a", "z"})
	fmt.Println(v)

	// Canon out of order, false.
	v = IsNormal(az, Canon{"z", "a"})
	fmt.Println(v)

	// Canon in order variadic, true.
	v = IsNormal(az, Canon{"a"}, Canon{"z"})
	fmt.Println(v)

	////////////////////
	// Only
	////////////////////
	fmt.Println("\nOnly")

	// Only in order, true.
	v = IsNormal(az, Only{"a", "z"})
	fmt.Println(v)

	// Only out of order, true.
	v = IsNormal(az, Only{"z", "a"})
	fmt.Println(v)

	// Only with extra field, false.
	v = IsNormal(az, Only{"a", "y", "z"})
	fmt.Println(v)

	// Only in order variadic, true.
	v = IsNormal(az, Only{"a"}, Only{"z"})
	fmt.Println(v)

	// ////////////////////
	// // Option
	// ////////////////////
	// fmt.Println("\nOption")

	// // Option with optional field missing, true.
	// v = IsNormal(az, Option{"a", "z", "x"})
	// fmt.Println(v)

	// // Option in order with optional field missing and variadic, true.
	// v = IsNormal(az, Option{"a"}, Option{"z", "x"})
	// fmt.Println(v)

	// // Option out of order with optional field missing and variadic, false.
	// v = IsNormal(az, Option{"z"}, Option{"x", "a"})
	// fmt.Println(v)

	// // Option with extra pay field, false.
	// v = IsNormal(ayz, Option{"a", "y"})
	// fmt.Println(v)

	// ////////////////////
	// // Need
	// ////////////////////
	// fmt.Println("\nNeed")

	// // Need in of order, true.
	// v = IsNormal(az, Need{"a", "z"})
	// fmt.Println(v)

	// // Need out of order, true.
	// v = IsNormal(az, Need{"a", "z"})
	// fmt.Println(v)

	// // Need missing field, false.
	// v = IsNormal(az, Need{"a", "y", "z"})
	// fmt.Println(v)

	// // Need with option missing, true.
	// v = IsNormal(az, Need{"a", "z"}, Option{"y"})
	// fmt.Println(v)

	// // Need with option present, true.
	// v = IsNormal(ayz, Need{"a", "z"}, Option{"y"})
	// fmt.Println(v)

	// // Need, extra field, then option, false.
	// v = IsNormal(ayz, Need{"a"}, Option{"z"})
	// fmt.Println(v)

	// // Need, option,then extra field, false.
	// v = IsNormal(ayz, Need{"a"}, Option{"y"})
	// fmt.Println(v)

	//
	// Only
	// false
	// true
	// true
	//
	// Need
	// false
	// true
	// true
	// true
	// false
	//
	// Order
	// false
	// false
	// true
	// true
	// true
	// false
	// false
	//
	// Option
	// true
	// false

	// Output:
	// Canon
	// true
	// false
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
