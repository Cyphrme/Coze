# Coze vs JOSE
## Why JOSE is awesome.
- Has the goal of doing crypto in a somewhat human readable paradigm.
- Has the goal of updating old standards that are hard to use and sometimes
  require specific libraries, binaries, encodings, or other outside dependences.
- Defines a way to represent cryptographic keys in JSON.
- JSON crypto keys, both public and private, have thumbprints, which is like a
  PGP fingerprint or Ethereum address. Thumbprints universally address specific
  keys. 

## Why JOSE could be better.
- The "unencoded" option is still encoded, and was an afterthought.  (RFC 7797)
- Thumbprints were an afterthought, and defined in a later RFC.  
- Thumbprints have no way to signify hash algorithm (as of 2021/05/04).  Later,
  additional RFCs have followed this implicit requirement.  For example RFC
  8037 specifies that Ed25519 and Ed448, neither of which use SHA-256, use
  SHA-256 for their thumbprints. 
- Because headers are always transmitted encoded and not as strings, they increase 
  in size.  For example,
  `"eyJhbGciOiJIUzI1NiIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19"` is larger than
  the unencoded representation `{"alg":"HS256","b64":false,"crit":["b64"]}`
- Converts UTF-8 to b64ut and encodes that into ASCII bytes, and then then
  hashes/signs those bytes. That's at least one extra conversion step we
  consider unneeded.  
- Protected headers.  For example, "alg" is required but doesn't always have to
  appear in the "protected" header.  This makes "protected"/"unprotected"
  headers less meaningful. 

- Any string that is b64ut encoded grows in size, so normal JOSE objects, both
  the compact (like JWT) and JSON forms grow in size. 
- Using b64ut everywhere makes JOSE far less useful when small messaging is
  critical while also decreasing human readability. 
- There's no uniform method to represent digests. 


- Byte representation is defined RFC 4648 base64 URL safe truncated.  Since b64ut doesn't fit evenly into bytes, this makes
  conversion more difficult.  



JOSE base 64 url:
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
- Coze uses Hex and not b64ut.  b64ut is also a "bucket convert" encoding, which
  we also find cumbersome.  Hex is equivalent in "bucket convert" (treating
  strings as varying units of bytes) as well as "iterative divide by radix"
  (idbr), treating strings as numbers conversion. 
- JOSE doesn't give a good way to represent byte information.  You're on your
  own to figure it out although it's suggested to use a base64 url truncated.
  See section on base64 regarding JOSE comparison. 
- Canonicalization is used in JOSE, but it's only applied narrowly to
  thumbprints.  Coze uses the generalized "Canonical Hash" (CH) to thumbprint
  any JSON object or binary blob, including keys and messages.  
- Instead of "claims" inside of "payload" which is separate from head, Coze puts
  everything in head. 


## Why "RFC 4648 base64 URL Safe Truncated" isn't great.
b64ut is hard, and can't easily be used with arbitrary base conversion. 
- b64ut is not very efficient when transmitting/storing as a character per byte.
 If going through all the trouble of conversion, it should be 1. Efficient or
 2. easy to use. base64 is neither efficient nor easy to use.  It's a halfway
 measure. 
- Calculates padding and then doesn't use it.
- Not a generalized conversion method.
- Doesn't fit into bytes evenly.
- Reencoding a base64 string into base64 results in ever larger string sizes.




# Coze Key
## How a Coze Key is different from a JWK
A Coze Key is like a JOSE JWK, but it differs in a few significant ways. 
 
1. Binary data is represented in hex, not b64ut.  For example, the "x" parameter
   for a ECDSA key is in hex and not b64ut.
2. "iat" (issued at) is required for messages and keys. 
3. "tmb" may be included in the Coze key.  "tmb" is deterministic digest from
   the key's canonical form and uses the hashing algorithm specified by `alg`. 
  - For JOSE, ["Selection of Hash
    Function"](https://tools.ietf.org/html/rfc7638#section-3.4) isn't well
    defined.  Coze explicitly defines how this is done. 
4. For ECDSA, the Coze Key thumbprint canon is {"alg","x","y"}.  For EdDSA the
   thumbprint canon is {"alg","x"}
5. "alg" (algorithm) is required and must refer to a specific cryptographic
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
6. "kty" and "crv" are optional **and redundant**.  "alg" takes the place of
  "kty" and "crv". 
  - Here's our thinking: Instead of "EC" for the value of "kty" for ECDSA or
    ES256, the value of "kty" should have been "ECDSA" and then "ES256" for alg.
    JOSE could have even included "kta" for "key type algorithm" and set that to
    ES256 if they didn't want to reuse "alg". The authors of this
    document see no conflict in using "alg" for keys as well, and makes the
    standard simpler, more descriptive, more consistent, and easier to
    understand.  
7. "kid" ("Key ID") is an optional human readable label for the key.  "kid" must
   not be used for anything programmatic. 
 - We think JWK's use of `kid` was a bad idea in because it says it "is a hint
   indicating which key was used".  What is the key hint?  We think this should
   have be explicitly structured.  This is why Coze specifies `tmb` , which is
   explicitly structured and used to identify the key used for siging.  We
   consider `kid` useless for programmatic function, so we reuse it for human
   readable key labeling. 
8.  We use "use" and _NOT_ "key_ops".  "usages" (which is used by
   Javascript implementations) and "key_ops" are both prohibited. 
9. "Ed25519" and "Ed448" is an algorithm ("alg"), not a curve ("crv"). An
   example of a curve would be "Curve25519".  The authors of Coze consider this
   to be one of the more head scratching JWK/IANA decisions.  For Coze, Ed25519
   is instantiated with specific key parameters, for example, "Ed25519" has the
   hashing algorithm  SHA-512. "Ed25519" is a sufficient identifier for both the
   key and the signing algorithm. 
10. The hex representation is not the value.  This is a small, and for most
   applications, irrelevant distinction. However, for the authors of Coze, we
   found this allowed more efficient references and storage of information. 
- For example, a "tmb" of a key is the bytes and not the hex representation.
    This allows thumbprints to be stored in databases as bytes and not strings.
    This also allows lookups from arbitrary bases or encodings, since the thing
    itself is defined to be the byte values. 
 - The only time hex is used as the value for representation is when signing
   JSON object strings which include hex values.  In that case the values are
   the UTF-8 byte representation of the hex.  We consider this an okay tradeoff
   since these byte values are encapsulated in larger stringed datastructres
   anyway. 
 
 
# Coze Key and Javascript's JWK implementation. 
Note on Javascript's Subtle.Crypto
 
## usages != use != key_opts  
Example of silliness:
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



 ## Example JWK (Coze does not use JWK, this is just for comparison):
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


# More on Base Conversion  Waste
One of the reason why we rejected JOSE's base64 paradigm is that it's very
inefficient.  

## Encoding Waste Example

The example string, "Potatoes,"

is 9 characters, and is encoded in UTF-8 as 9 bytes.  

Encoded into base64, this string is: UG90YXRvZXMs

Which is 12 characters.  All strings are still encoded as UTF-8 in JOSE,
including base64, which is 12 bytes. Base64 is only 75% efficient in the byte
space. 

25% waste doesn't sound bad, until you realize that normal english plus URL
characters uses a about 98 characters out of the potential 256 for byte
encoding.  

98 common characters: (Allowing for space, tab, line feed, and carriage return
at the end)

0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_~!"$'()*+,:;<=>?@[\]^`{|}#%&/.  

If Base98 was encoded efficiently, it should use 6.61 bits per character (2^x =
98  which is x ~= 6.61 bits per character.)

The string "Potatoes," would efficiently be encoded into ~60 bits.  The base64
representation of this string is 96 bits.  This means that the string is only
about 63% efficient.  

base64 forces a minimum amount of waste onto messages.  By allowing strings to
remain strings, more efficient encoding is possible.  

Please see the document "Efficient Barcode-URI Encoding for Arbitrary Payloads"
for more notes on efficient encoding methods.  