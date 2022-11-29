TODO consider using compressed ECDSA.  Should be able to decrease public key
sizes to 33-ish bytes, depending on how it is implemented.  

Y^2 = X^3 + aX + b

See section 2.3.3
https://www.secg.org/sec1-v2.pdf

Golang:
https://pkg.go.dev/crypto/elliptic#MarshalCompressed

Javascript:
There doesn't appear to be a standard library that can be used. 

Someone wrote it themselves here: 
https://stackoverflow.com/questions/17171542/algorithm-for-elliptic-curve-point-compression

The did group has a discussion on the issue:
https://github.com/w3c-ccg/did-method-key/issues/32


Good little primer: https://medium.com/asecuritysite-when-bob-met-alice/02-03-or-04-so-what-are-compressed-and-uncompressed-public-keys-6abcb57efeb6


JOSE does not support point compression.  
>SEC1 [SEC1] point compression is not supported for any of these three curves.
https://www.rfc-editor.org/rfc/rfc7518#section-6.2.1.1



# Thumbprint considerations

Thumbprints should continue being produced as they currently are, e.g. "x":"x concated y"

This allow Coze to support systems that do not support compression, for whatever reason.  

# Support
Coze should support 

1. No SEC compression awareness (known by length, xsize).  
2. Two modes of short SEC compression (02, 03) (Known by (xsize/2) + 1 )
3. Should not support SEC mode 04, since that is the default.


# Example
Using the example key: 

    "x":"2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"

x and y:

    2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjO
    Rojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g

y in decimal:
    31903811644451988338668220999084907216815044169343737858363203889617182630808

Which is even.  So x gets prepended with byte "02"

Atp0zmhVZtkC8ZlDv0o4MrHFRwbbxxH6Nq6uuTL4DUYz


The full resulting key:

```json
{
	"alg":"ES256",
	"iat":1623132000,
	"kid":"Zami's Majuscule Key.",
	"d":"bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVA",
	"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
	"x":"Atp0zmhVZtkC8ZlDv0o4MrHFRwbbxxH6Nq6uuTL4DUYz"
}
```




The uncompressed key, uncompressed aware is (87 b64ut characters):

```
BNp0zmhVZtkC8ZlDv0o4MrHFRwbbxxH6Nq6uuTL4DUYzkaI6t_R2qva1zcb18cG2v149Beb2YmyUd4rAXTlm6OY
```

(to Hex: DA74CE685566D902F19943BF4A3832B1C54706DBC711FA36AEAEB932F80D463391A23AB7F476AAF6B5CDC6F5F1C1B6BF5E3D05E6F6626C94778AC05D3966E8E6
Add byte 04 : 04DA74CE685566D902F19943BF4A3832B1C54706DBC711FA36AEAEB932F80D463391A23AB7F476AAF6B5CDC6F5F1C1B6BF5E3D05E6F6626C94778AC05D3966E8E6)

This seems redundant and is suggested to not be used.    