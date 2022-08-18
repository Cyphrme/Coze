# Coze Vs "X" Disclaimer
We have a lot of respect for the various projects in the space.  Many projects
have noble goals and we're thankful they exists.  We also don't think it's cool
to "take a dump" all over someone else's work when the authors have worked so
hard to bring value, freely, to everyone.  

We also think it's important to give specific reason why Coze's design is
different from other projects.  In this document, we attempt to give specific
reasons why Coze was needed.  

# Why not "x"?

## signify (OpenBSD):
 - Not JSON.  
 - No browser implementations. 
 - No algorithm agility.  
 - No real plan to expand its use.  



# Coze vs JOSE

We have a lot of respect for JOSE.  We think its goals are noble
and we're glad it exists.  

## Why JOSE is awesome.
- Has the goal of doing crypto in a somewhat human readable paradigm.
- Has the goal of updating old standards that are hard to use and sometimes
  require specific libraries, binaries, encodings, or other outside dependencies.
- Defines a way to represent cryptographic keys in JSON.
- JSON crypto keys, both public and private, have thumbprints, which is like a
  PGP fingerprint or Ethereum address. Thumbprints universally address specific
  keys. 
- Permits algorithm agility.

## How JOSE could be better.
- JWT is not JSON (despite the name).  JWT is not JSON in both encoded and
  decoded form.
- The "unencoded" option is still encoded, and was added to the standard later.
  (RFC 7797)
- Thumbprints have no way to signify hash algorithm (as of 2021/05/04) and it
  appears to be always assumed to be SHA-256, even for ES384 and ES512.  Later,
  additional RFCs have followed this implicit requirement.  For example RFC 8037
  specifies that Ed25519 and Ed448, neither of which use SHA-256, use SHA-256
  for their thumbprints. 
- Payers are always transmitted encoded and as base64 and they increase in
  size.  For example,
  `"eyJhbGciOiJIUzI1NiIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19"` is larger than
  the unencoded representation `{"alg":"HS256","b64":false,"crit":["b64"]}`.
- Converts UTF-8 to b64ut and encodes that into ASCII bytes, and then 
  hashes/signs those bytes. That's at least one extra conversion.
- JOSE's double encoding of some base64 values is inefficient.  
- Protected headers.  For example, "alg" is required but doesn't always have to
  appear in the "protected" header.  This makes "protected"/"unprotected"
  headers less meaningful. 
- Any string that re-encodes b64ut grows in size. normal JOSE objects, both the
  compact (like JWT) and JSON forms grow in size.


JOSE:
https://tools.ietf.org/html/rfc7515#appendix-A.4.1
 
Example JWS aesthetic:
```
eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGF
tcGxlLmNvbS9pc19yb290Ijp0cnVlfQ",
"signatures":[{"protected":"eyJhbGciOiJSUzI1NiJ9", sign over:
UTF8(eyJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ)
```
 
# Coze
## Coze Vs JOSE
### Key Differentiators from JOSE to Coze.
- Canonicalization is used in JOSE, but it's only applied narrowly to
  thumbprints.  JWS and JWTs can be out of order and not canonicalized.  
		-The JWT MUST conform to either the [JWS] or [JWE] specification.  Note that
		whitespace is explicitly allowed in the  representation and no
		canonicalization need be performed before encoding. [...] [A]pplication[s]
		may need to define a convention for the canonical case [...] if more than
		one party might need to produce the same value so that they can be compared.
	- Coze uses the generalized "Canonical Hash" (CH) to thumbprint any JSON
  object or binary blob, including keys and messages.  
- Instead of "claims" inside of "payload" which is separate from head, Coze puts
  everything in pay. 



# Coze Key
## How a Coze Key is different from a JWK
A Coze Key is like a JOSE JWK, but it differs in a few significant ways. 
 
 Coze requires:
1. "iat" (issued at) is required for messages and keys. 
2. "tmb" may be included in the Coze key.  "tmb" is deterministic digest from
   the key's canonical form and uses the hashing algorithm specified by `alg`. 
  - For JOSE, ["Selection of Hash
    Function"](https://tools.ietf.org/html/rfc7638#section-3.4) isn't well
    defined.  Coze explicitly defines how this is done. 
3. The Coze Key thumbprint canon is {"alg","x"}.  
4. "alg" (algorithm) is required and must refer to a specific cryptographic
  algorithm.  "alg" should be descriptive of any parameter information needed
  about the key's signing algorithm.  For example, for an ecdsa key, "alg"
  should be descriptive of the type of signing algorithm (ECDSA), the curve
  (P-256), and the hash (SHA-256), which "ES256" is fully descriptive. 
  - Coze does not allow keys to interchange signing or key parameters by
    designed.  For example, a key designed to be used with ES256 must only ever
    use the same ES256 parameters (such as the curve, hashing algo/design,
    ect...) and only ever be used with ES256 signatures. 
 - Note: "EC" or "ECDSA" is insufficient for the value of the "alg" parameter
   since they are not descriptive of a specific cryptographic algorithm. 
5. "kty" and "crv" are optional and redundant.  "alg" takes the place of
  "kty" and "crv". 
  -  In JOSE, instead of "EC" for the value of "kty" for ECDSA or ES256, the value of
    "kty" could have been "ECDSA" and then "ES256" for alg. JOSE could have even
    included "kta" for "key type algorithm" and set that to ES256 if saw
    conflict in reusing "alg". For Coze we saw no conflict in using "alg" for
		keys as well, and makes the standard simpler, more descriptive, more
    consistent, and easier to understand.  
6. "kid" ("Key ID") is an optional human readable label for the key.  "kid" must
   not be used for anything programmatic. 
 - JOSE says that "kid" "is a hint indicating which key was used".  What is the
   key hint?  `tmb` is better explicitly structured.  This is why Coze specifies
   `tmb` , which is explicitly structured and used to identify the key used for
   singing.  Since `kid` isn't ideal for programmatic function, we use it as
   human readable key labeling. 
7. "use" and "key_ops" are redundant.  "usages" (which is used by
   Javascript implementations) and "key_ops" are both absent in Coze. 
8. For Coze, "Ed25519" and "Ed448" is an algorithm ("alg"). An example of a
   curve would be "Curve25519".  In JOSE, `crv` is "Ed25519" and is combined
	 with a key type of OKP.
	 (https://datatracker.ietf.org/doc/html/rfc8037#appendix-A.3)
		- For Coze, Ed25519 is instantiated with specific key parameters, for
   example, "Ed25519" has the hashing algorithm SHA-512. "Ed25519" is a
   sufficient identifier for both the key and the signing algorithm. 

 
 
# Coze Key and Javascript's JWK implementation. 
Note on Javascript's Subtle.Crypto
 
## usages != use != key_opts  
Example of JOSE:
"use":"sig",
"key_ops": "["sign", "verify"]"
 
- Javascript includes "key_ops" but the RFC says "key_ops" should
 not be used with "use" (2021/05/27).  "use" is far more clear.  (See
 https://datatracker.ietf.org/doc/html/rfc7517#section-4.3 for where the RFC
 clearly says "key_ops" should not be used with "use".) Further, Javascript
 uses "usages" which is confusing with the RFC's "use".  Eliminating "key_ops"
 and Javascript's "usages" makes the key's "usages" clear. 
 
- Javascript's cryptoKey.usages doesn't allow for verification.
 - A Crypto key with "usages" of `["verify","sign"]` cannot be used to verify
   (2021/05/27). Chrome throws an error and there's no docs as to why.  This is
   also counter to the JWK RFC about "use".  What's the point of `"use"` in the
   JWK if the browser doesn't even use it?
 
-Javascript's CryptoKey does not use `"use"`, and does not have, JWK's `"use"`
as of 2021/05/13.
- As a further pain point, "use" is a single string, "key_ops", as
 required by javascript, is required to be an array.  (Example: "[sign]")



## Reference Example JWK:
```json
{
 "crv": "P-256",
 "d":   "bJnCQX7Ogd91FTIkmKtXeYFfjUfN4sQ3YXz2hLIbxJQ",
 "use": "sig",
 "kty": "EC",
 "x":   "JxnHyqkG9J4gygj9jBhooRIOmGNcHTdplNt3ODhEtmo",
 "y":   "zueErjY0awFg9-7bt3NRnUFj1ZrL8MNc8kIYM1AQFwI"
}
```



## Encoding Waste Example
The example string, "Potatoes,"

is 9 characters, and is encoded in UTF-8 as 9 bytes.  

Encoded into base64, this string is: UG90YXRvZXMs

Which is 12 characters.  All strings are still encoded as UTF-8 in JOSE,
including base64, which is 12 bytes. Base64 is only 75% efficient in the byte
space. 

Normal english plus URL characters uses a about 98 characters out of the
potential 256 for byte encoding.  

98 common characters: (Allowing for space, tab, line feed, and carriage return
at the end)

```
0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_~!"$'()*+,:;<=>?@[\]^`{|}#%&/.
```  

If Base98 was encoded efficiently, it should use 6.61 bits per character (2^x =
98  which is x ~= 6.61 bits per character.)

The string "Potatoes," would efficiently be encoded into ~60 bits.  The base64
representation of this string is 96 bits.  This means that the string is only
about 63% efficient.  

Please see the document "Efficient Barcode-URI Encoding for Arbitrary Payloads"
for more notes on efficient encoding methods.  