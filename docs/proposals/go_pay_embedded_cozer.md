Decided against this design pattern in favor of embedding structs into Pay.


Problems with this design:
- Cannot embed Pay directly into structs in Go, which is ugly.   
	- {Foo:bar, Pay: Pay{...}}
- JSON tags have to be correct on Coze types, e.g. `Alg` vs `alg`
	- Burdens downstream with Coze tags.   
- Pollutes 3rd party structs

 By embedding 3rd party structs in Pay, these problems are avoided.  

```Go
// Cozer is an interface for any type that returns a Coze.
type Cozer interface {
	Coze() Coze
}

// Example of how to implement the Coze method in the Cozeer interface on a
// custom data structure.
type User struct {
	Alg         SEAlg
	Tmb         B64
	DisplayName string
	FirstName   string
	LastName    string
	Email       string `json:",omitempty"`
}

// Implements the Cozeer interface.
func (u *User) Coze() (coz *Coze, err error) {
	coz = new(Coze)
	coz.Pay, err = Marshal(u)
	return coz, err
}

func ExampleCozeer_Coze() {
	ucoz := User{
		Alg:         SEAlg(ES256),
		Tmb:         MustDecode("cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk"),
		DisplayName: "Coze",
		FirstName:   "Foo",
		LastName:    "Bar",
	}

	c, err := ucoz.Coze()
	if err != nil {
		fmt.Println(err)
	}

	err = Golden_Key.SignCoze(c, nil)
	if err != nil {
		fmt.Println(err)
	}

	v, err := Golden_Key.VerifyCoze(c)
	if err != nil {
		fmt.Println(err)
	}

	// Set sig to nil for deterministic printout
	c.Sig = nil

	fmt.Printf("%+v\n", c)

	// Output:
	// {"pay":{"alg":"ES256","tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","DisplayName":"Coze","FirstName":"Foo","LastName":"Bar"}}
}
```