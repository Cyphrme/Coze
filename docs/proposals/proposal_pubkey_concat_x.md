
EdDSA, Ed25519, concatenates public components in the JSON label "x" (the RFC names it "B")


https://datatracker.ietf.org/doc/html/rfc8032#section-5.1

```
B     | (X(P),Y(P)) of edwards25519 in [RFC7748] (i.e., (1511 |
   |           | 22213495354007725011514095885315114540126930418572060 |
   |           | 46113283949847762202, 4631683569492647816942839400347 |
   |           | 516314130799386625622561578303360316525185596
```

- FIPS defines the ECDSA public key as "Q", private as "d"
- Wikipedia names the ECDSA public key as "Q_A"
- JOSE names the public parts "x" and "y". 

Current Form:
```JSON
{
	"alg":"ES256",
	"iat":1624472390,
	"kid":"Zami's Majuscule Key.","tmb":"0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD","x":"DA74CE685566D902F19943BF4A3832B1C54706DBC711FA36AEAEB932F80D4633",
	"y":"91A23AB7F476AAF6B5CDC6F5F1C1B6BF5E3D05E6F6626C94778AC05D3966E8E6"
}
```

Proposed Key form:
```JSON
{
	"alg":"ES256",
	"iat":1623132000,
	"tmb":"2C65C2CE62249683E37FEF7933539C7A1B364F80BC4552908A6C9DE8BFDEAB01",
	"kid":"Zami's Majuscule Key.",
	"x":"DA74CE685566D902F19943BF4A3832B1C54706DBC711FA36AEAEB932F80D463391A23AB7F476AAF6B5CDC6F5F1C1B6BF5E3D05E6F6626C94778AC05D3966E8E6",
}
```

Key in canonical form:
```JSON
{"alg":"ES256","x":"DA74CE685566D902F19943BF4A3832B1C54706DBC711FA36AEAEB932F80D463391A23AB7F476AAF6B5CDC6F5F1C1B6BF5E3D05E6F6626C94778AC05D3966E8E6"}
```

Resulting in the thumbprint and cad of:
`2C65C2CE62249683E37FEF7933539C7A1B364F80BC4552908A6C9DE8BFDEAB01`

Neutral:
1. Already have padding code if needed for future `algs`.
2. Componets are already in Hex which is padded and ready for concatenation. 

Pros:
1. Public key will be the same size as sig.  
2. Consistent
3. Easier to implement, no special logic for "y".  

Cons:
1. Cyphr.me has to rebuild everything
2. Very minor: To interact with libraries, public key will need to be split into to parts.


Conclusion:

Move to concated (when time permits).  