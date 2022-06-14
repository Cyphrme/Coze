package coze

import (
	"encoding/json"
	"fmt"
)

func ExampleCy_Meta() {
	cy := new(Cy)
	err := json.Unmarshal([]byte(Golden_Cy), cy)
	if err != nil {
		fmt.Println(err)
	}

	err = cy.Meta()
	if err != nil {
		fmt.Println(err)
	}

	cyb, err := Marshal(cy)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", cyb)

	// Output:
	//{"can":["msg","alg","iat","tmb","typ"],"cad":"LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto","cyd":"d0ygwQCGzuxqgUq1KsuAtJ8IBu0mkgAcKpUJzuX075M","pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"}
}

//ExampleCy_jsonUnMarshal tests unmarshalling a Cy.
func ExampleCy_jsonUnMarshal() {
	cy := new(Cy)
	err := json.Unmarshal([]byte(Golden_Cy), cy)
	if err != nil {
		fmt.Println(err)
	}

	// remarshal for comparison
	b, err := Marshal(cy)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
	// Output:
	//{"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"}
}

func ExampleCoze_jsonMarshal() {
	var err error

	cye := new(Coze)

	err = json.Unmarshal([]byte(Golden_Coze), cye)
	if err != nil {
		fmt.Println(err)
	}

	b, err := Marshal(cye)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+s\n", b)
	// Output:
	//{"coze":{"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1627518000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"}}
}
