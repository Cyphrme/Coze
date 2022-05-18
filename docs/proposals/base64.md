# base64: Support base64 encoded values.

## Problem: Hex Coze messages are longer than JWT

Coze messages are larger than they would be with base64 values. Coze messages
can be (much) smaller than JWTs if base64 encoded.  

For example, when compactified the following is 298 characters.  
```json
{
"head": {
 "alg": "ES256",
 "iat": 1623132000,
 "msg": "Coze Rocks",
 "tmb": "0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD",
 "typ": "cyphr.me/msg/create"
},
"sig": "E848D97CA3A1BAE8C1AE6ACBAE1E73B7C23C9A74581003CAEB4FCBA4EF39EC8B07996B4F52F5D5925C48A793C54495A3B89DD9A8B55D29E72B8B9DF599E0A734"
}
```

A comparable JWT is smaller at 280 characters.  

```jwt
eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJtc2ciOiJDb3plIFJvY2tzIiwiaWF0IjoxNjI3NTE4MDAwLCJ0bWIiOiJyYkxYM1NzV0xXQkpvSXFNUENOVUZ1VURScFZVX28tMFNERms4WWxURXU4IiwidHlwIjoiY3lwaHIubWUvbXNnL2NyZWF0ZSJ9.7uLr31zS5_I-UeJWj4Olrufu9C7sr2-2DB4dDyKY4yf3g6Jr30JSLS3wfyMEWUbW1OVAzsB1wYhaWbUz0VWtGA
```



Coze would be smaller if base64 was used.  The following is 235 characters.  

```json
{
"head": {
 "alg": "ES256",
 "iat": 1623132000,
 "msg": "Coze Rocks",
 "tmb": "AUj0zZCTycvj6L940-bJuCTxHdLynisaYw3Rzh4XbN0",
 "typ": "cyphr.me/msg/create"
},
"sig": "6EjZfKOhuujBrmrLrh5zt8I8mnRYEAPK60_LpO857IsHmWtPUvXVklxIp5PFRJWjuJ3ZqLVdKecri531meCnNA"
}
```

# Solution
Since `alg` already defines sizes and values are padded, encoding is discernible
from `alg` and string length.

For example, for alg: "ES256", "tmb" is 64 in Hex and 43 in base64.

# Considerations

## Consideration for using length in determining length
This proposal uses string length instead of an addition identifier for base64.

Currently, if the lengths are wrong an error should be thrown.  That logic will
be the same with the added inclusion of base64 lengths, which can be calculated
by (4/3 * number of bytes) rounded up to the nearest integer.

Single byte Hex, e.g. with two values `FF`, would be indiscernable from the
base64 `_w` based on length.  All other longer values are discernable:  the
Hex `FF01` converts to base64`_wE`.

This situation is considered unlikely and is considered to be an acceptable
cost.  Since Coze specifies that values are always padded, if this ever was a
problem, padding to two bytes would be suggested, e.g. the Hex "00FF" converts
to base64 "AP8".

## Consideration for RFC 4648
RFC 4648 base64 URI truncated (b64ut) is the selected encoding (padding
characters removed and use the base64url alphabet). 

## Consideration for processing messages and keys
The encoding isn't relevant. Hex and base64 are still processed as text (UTF-8).
Hashing can still be started once `alg` is parsed since 

## Consideration for mixed encoding
Mixed encoding is unsupported and should result in an error. This is simpler to
implement.  All relevant lengths can be checked and if any do not match any
operation should error.  All standard Coze messages and keys should have an
encoded payload, such as "tmb", to perform the switching over.  

Alternatively, each string's length would need to be checked individually.  

If a system wants to sign something from another encoding or mixed encoding,
convert to the desired encoding before signing.  

## Consideration for the online tool.
The online tool will have an additional added button Hex to base64 and base64 to
Hex.  

## Considered Alternatives
1. Use an identifier for base64 encoding.  "_64"

## Future considerations
Perhaps base 85, basE91, or another efficient encoding can be supported,
especially since JSON itself is not URI safe. 

Since we want to support more efficient encodings than base64, supporting
multiple encodings, (Hex and base64) appears to be a logical first step.  

## Why not deprecate Hex?
1. Hex is human readable.  
 A. Small alphabet. 
 B. No doppelgängers (O0I1l)
2. Hex is phonetic
	A. Many people know how to pronounce Hex. 
	B. "-" and "_" is difficult for some people to pronounce.  
3. We want to support even more efficient encodings in the future.  


## Appendix

### JWK used to create valid JWT/JWS
{"kty":"EC","d":"MAqyJgK2kB5YsouHtLbaiowPqGCQj15hG3B5f65IG9w","crv":"P-256","x":"FD7r1byJlBcvsBXTi47ltG4JRbMVHgxp91Ds2efiM_g","y":"5fec2TDH4QJ1fBmIogTgHB6y00H9lJiVQ7MoqB6Tidg"}

### JOSE Invalid hypothetical unencoded JWT:
 A JWT in the invalid "unencoded" representation, with the header and
payload unencoded, is 240 characters.  JOSE always encodes headers as base64,
**even in the "unencoded" option**, so this efficient hypothetical isn't
valid in JOSE.

```jwt
{
  "alg": "ES256",
  "typ": "JWT"
}.{
  "msg": "Coze Rocks",
  "iat": 1627518000,
  "tmb": "rbLX3SsWLWBJoIqMPCNUFuUDRpVU_o-0SDFk8YlTEu8",
  "typ": "cyphr.me/msg/create"
}.7uLr31zS5_I-UeJWj4Olrufu9C7sr2-2DB4dDyKY4yf3g6Jr30JSLS3wfyMEWUbW1OVAzsB1wYhaWbUz0VWtGA
```
