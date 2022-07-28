# ⚠️ COZE IS IN ALPHA.  USE AT YOUR OWN RISK ⚠️
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/github.com/cyphrme/coze)](https://pkg.go.dev/github.com/cyphrme/coze)

![Coze](docs/img/coze_logo_zami_white_450x273.png)

[Presentation](https://docs.google.com/presentation/d/1bVojfkDs7K9hRwjr8zMW-AoHv5yAZjKL9Z3Bicz5Too)

# Coze 
Coze is a cryptographic JSON messaging specification designed for human
readability.

Play with Coze here: https://cyphr.me/coze_verifier

![coze_verifier](docs/img/Hello_World!.gif)

### Example Coze
```JSON
{
	"pay": {
		"msg": "Coze Rocks",
		"alg": "ES256",
		"iat": 1627518000,
		"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
		"typ": "cyphr.me/msg"
	},
	"sig": "ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"
}
```

### Coze Design Goals
1. Valid and idiomatic JSON. 
2. Human readable and writable.
3. Small in scope.
4. Cryptographic agility.


### Coze Names
Coze objects encapsulate a set of JSON name/value pairs.  Coze reserved names
are short, unique, and unlikely to require namespacing by applications.
Applications are permitted to use any name except reserved names. The Coze
objects `pay`, `key`, and `coze` have respective reserved names and all names
must be unique.

Binary values are encoded with URI safe base64 with padding omitted (b64ut).

Coze requires that:
- JSON objects must be valid JSON with unique fields names.
- JSON field names must be strings (which is standard JSON).
- JSON field names and values are case sensitive.

![Coze Reserved Fields](docs/img/coze_reserved_fields.png)

## Pay
`pay` may contain the standard fields `alg`, `iat`, `tmb`, and `typ` as
well as custom fields.  In the first example, `msg` is a custom field.

### `pay` Reserved Names
- `alg`  Specific signing algorithm.  E.g. `"ES256"`
- `iat`  The time when the message was signed. E.g. `1623132000`
- `tmb`  Thumbprint of the key used to sign the message.  E.g. `"0148F4..."`
- `typ`  Type of `pay`. `typ` may denote the canon of `pay`. E.g.
           `"cyphr.me/msg/create"` denotes the canon `["alg","iat","msg","tmb",
           "typ"]`

## Coze Key
### Example Public Coze Key
```JSON
{
	"alg":"ES256",
	"iat":1623132000,
	"kid":"Zami's Majuscule Key.",
	"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
	"x":"2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"
}
```

### Example Private Coze Key
```JSON
{
	"alg":"ES256",
	"iat":1623132000,
	"kid":"Zami's Majuscule Key.",
	"d":"bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVA",
	"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
	"x":"2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"
}
```

### Coze Key Standard Fields
Coze keys require `alg`, `iat`, `tmb`, and the public components according
to `alg`.  

- `key`  (Object) Key object. E.g. `"key":{"alg":"ES256", ...}`
- `alg`  (Algorithm) Signing algorithm.  E.g. `"ES256"`
- `d`    (Private) Private component.  E.g. `"bNstg4..."`
- `iat`  (Issued At) When the key was created. E.g. `1623132000`
- `kid`  (Key Identifier) Human readable, non-programmatic label for the key.
            E.g. `"kid":"My Cyphr.me Key"`. 
- `tmb`  (Thumbprint) Key thumbprint.  E.g. `"cLj8vs..."`
- `x`    (Public) .  E.g. `"2nTOaF..."`.
- `typ`  (Optional) Optional type.  
- `rvk`  (Optional) Time of key revocation.  See the `rvk` section.  

Note that `kid` must not be used programmatically.

Note that the private component `d` is not included in `tmb` generation.  `tmb`
is generated from `alg`'s thumbprint canon `["alg", "x"]` and `alg` 's hashing
algorithm, for example `SHA-256`. See the thumbprint section for more.



## `coze` 
The JSON name `coze`

### `coze` Standard Fields
- `coze` JSON label for a Coze object.  E.g. `{"coze":{"pay":..., sig:...}}`
- `can`  Canon for hashing over `pay`.  E.g. `["alg","iat","tmb","typ"]`
- `cad`  Canon digest.  The digest of `pay`.  E.g.: `"24F11D..."`
- `czd`  Coze digest, the digest over `["cad","sig"]`. `czd`'s hash must align
  with `alg` in `pay`.  
- `pay` Label for the pay object.  E.g. `"pay":{"alg":...}`
- `sig`  Signature over the bytes of `cad`, and `sig` does not rehash `cad`
  before signing.  E.g. `"sig":"CC3AD6..."`

Like `cad`, `czd` is calculated from brace to brace, including the braces.  

`cad` and `czd` are recalculatable and are recommended to be omitted, although
they may be useful for reference.  

## Wrapped Coze
The JSON name `coze` is used to wrap Coze objects.  For example:

```JSON
{
	"coze":{
		"pay": {
			"msg": "Coze Rocks",
			"alg": "ES256",
			"iat": 1627518000,
			"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
			"typ": "cyphr.me/msg"
		},
		"sig": "ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"
	}
}
```

We recommend against needlessly wrapping implicit Coze objects with `coze`. For
example, the JSON object `{"pay":{...},"sig":...}` doesn't need the labeled
`coze` if already implicitly known.

#### Example "full" `coze`  
The following, containing `pay`, `key`, `can`, `cad`, `czd`, and `sig`, expands
the first example and is largely redundant. `key` may be looked up based on
`tmb`. `can`, `cad`, and `czd` are recalculatable and generally should be
omitted.  The label `coze` may be inferred.

The following tautological coze:  

```JSON
{
	"coze": {
		"pay": {
			"msg": "Coze Rocks",
			"alg": "ES256",
			"iat": 1627518000,
			"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
			"typ": "cyphr.me/msg"
		},
		"key": {
			"alg":"ES256",
			"iat":1623132000,
			"kid":"Zami's Majuscule Key.",
			"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
			"x":"2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"
		},
		"can": ["alg","iat","msg","tmb","typ"],
		"cad": "LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto",
		"czd": "d0ygwQCGzuxqgUq1KsuAtJ8IBu0mkgAcKpUJzuX075M",
		"sig": "ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"
	}
}
```

Simplifies to:

```JSON
{
	"pay": {
		"msg": "Coze Rocks",
		"alg": "ES256",
		"iat": 1627518000,
		"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
		"typ": "cyphr.me/msg"
	},
	"sig": "ywctP6lEQ_HcYLhgpoecqhFrqNpBSyNPuAPOV94SThuztJek7x7H9mXFD0xTrlmQPg_WC7jwg70nzNoGn70JyA"
}
```

## Normal
## Canon
Coze JSON objects are canonicalized and hashed for creating digests, signing,
and verification. 

For JSON objects the canonical digest is generated by hashing the canonical
form.  For binary files the canonical digest is simply the digest of the file.

The hashing algorithm is specified by `alg`.  For example, the hashing algorithm
denoted by `"ES256"` is `"SHA-256"`.

## Canon for JSON 
The three steps in generating the canonical form of a JSON object:
 0. Order by canon and do not include fields excluded from canon.
 2. Remove insignificant whitespace and serialize.

The canonical digest of JSON objects is generated by hashing the canonical form.
The canonical digest of `pay` is `cad` and the canonical digest of `key` is
`tmb`.  

The canon for thumbprints is derived from `alg`. For example, the canon derived
from `"ES256"` is `["alg","x","y"]`.

`pay`'s canon is the present fields in the order as they appear, implicitly
known by applications, implicitly derived by "typ", or explicitly
specified by `can`.  Applications should specify canon expectations in API
documentation.  If a message is malformed the application must error. 

In the first example the `cad`  is `LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto`.

`czd` is generated from the canonical form of `coze` with the canon
["cad","sig"] which results in the JSON `{"cad":"...",sig:"..."}`. `czd` refers
to a particular signed message. `czd`'s hash must align with `alg` in `pay`. 

In the first example, the `czd` is `d0ygwQCGzuxqgUq1KsuAtJ8IBu0mkgAcKpUJzuX075M`.

Applications may ignore `can` in `coze` and under typical use it is recommended
that `can` be omitted.  


### Key Thumbprints
The key thumbprint, `tmb`, is the canonical digest of `key` with the canon 
`["alg","x"]`.  Using this canon, the canonical form of the key above is:

```JSON
{"alg":"ES256","x":"2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"}
```

Hashing results in the digest value of tmb: `cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk`.


### Canon for Binaries
The canonical digest, `cad`, of a binary file is simply the digest of the
file using the hash specified by `alg`. For example, an image
("Hello_World!.gif") may be referred to in a JSON object by its digest.

```JSON
{
	"alg":"SHA-256",
	"file_name":"Hello_World!.gif",
	"image":"rVOyJ144KwIQ3V2YJdatKAo_3QWAY4CpGLCDdnKOvAw"
}
```

As an application example, including a file's digest in a signed message,
denoted by `id`, may represent the authorization to upload a file to a user's
account:

```JSON
{
	"pay": {
		"alg": "ES256",
		"file_name": "Hello_World!.gif",
		"iat": 1657925839,
		"id": "rVOyJ144KwIQ3V2YJdatKAo_3QWAY4CpGLCDdnKOvAw",
		"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
		"typ": "cyphr.me/file/create"
	},
	"sig": "rY7XYc9sGYZX0jsNnVhOYvJGb_I6Z-xq8gdbOBw8K-M1uAD4J3V33yyw-FJMrihMeJr60wWgeHRXCKWFHb_SgA"
}
```


## Revoke
A Coze key is self-revoked by signing a self-revoke message.  A self-revoke
message has the field `rvk` with an integer value greater than `0`. The value of
`rvk` is a Unix timestamp.

### Example Self Revoke

```JSON
{
	"pay":{
		"alg":"ES256",
		"iat":1655924566,
		"msg":"Posted my private key on github",
		"rvk":1655924566,
		"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
		"typ":"cyphr.me/key/revoke"
		},
	"sig":"78Dp1YyArd19CMDDHFMUcLP3y876p1cpO6LTa94Pe8lKu8J3e2R93eK8EY3u2CaalJ6eV0O3b741atIDJ3uJgQ"
}
```



- `rvk` - Unix timestamp of when the key was revoked.  

Coze explicitly defines a self-revoke method so that third parties may revoke
leaked keys. Systems storing Coze keys should mark key as revoked when given a
self-revoke message.  Systems may use any non-zero value for `rvk` to denote key
revocation and the integer value "1" is suitable to denote revocation.

Key expiration policies, such as key rotation, are outside the scope of Coze.
Self revokes with future times must immediately be considered as revoked.  

## Supported Algorithms 
- ES224
- ES256
- ES384
- ES512
- Ed25519 
- Ed25519ph (planned)

### `alg` parameters:

"alg":"ES256"
- Genus: ESCDSA
- Family: EC
- Use: sig
- Sig.Size: 512
- Curve: P-256
- Hash: SHA-256
- Hash.Size:256 





## Coze Verifier
Cyphr.me has an online tool to sign and verify messages.  We hope in the near
future to release an open source, stand alone version of the webpage.  

Play with Coze here: https://cyphr.me/coze_verifier.

![coze_verifier](docs/img/Hello_World!.gif)


## Coze Implementations
 - [Go (This repo)](https://github.com/Cyphrme/coze)
 - [Javascript (Coze js)](https://github.com/Cyphrme/cozejs).



# FAQ




#### Pronunciation? What does "Coze" mean? 
We say "Co-zee" like a comfy cozy couch.  The English word Coze is pronounced
"kohz" and means "a friendly talk; a chat" which is the perfect name for a
messaging standard.

Jared suggested Coze because it looks like JOSE or COSE but it's funnier.


#### "Coze" vs "coze"?
We use upper case "Coze" to refer to the specification, and "coze" to refer to
coze messages and objects.


#### Why release pre-alpha on 2021/06/08?
Coze was released on 2021/06/08 (1623132000) since it's 30 years and one day
after the initial release of PGP 1.0.


#### Zero case
If `alg` and `tmb` are implicitly known, a zero case is legitimate. The
following is a valid coze.

```json
{
	"pay":{},
	"sig":"9iesKUSV7L1-xz5yd3A94vCkKLmdOAnrcPXTU3_qeKSuk4RMG7Qz0KyubpATy0XA_fXrcdaxJTvXg6saaQQcVQ"
}
```

#### Unicode/UTF-8?
Yes.  Unicode is a superset of ASCII and UTF-8 shares sorting order with
Unicode.  This results in broad, out of the box compatibility. Not that UTF-16
(Javascript) has some code points out of order. For these systems, a small
amount of additional logic is needed to correct the sort order.

#### Binary? Why not support binary payloads like JOSE?
JSON isn't well designed for binary.  Coze uses digests which we feel is an
acceptable compromise.  A binary file's digest is easily included in a coze,
while the binary itself should be transported outside of the coze. 

For example,  the digest of an image may be included in a coze.  The field `id`
below is the digest of an image. The coze includes other metadata.  

// TODO fix this
```JSON
{
 "pay": {
  "alg": "ES256",
  "ext": "png",
  "iat": 1623132000,
  "id": "D0AA358048454A7C564C42F8D066F0932A3985AAC3F71FF59A1BCD1372E58590",
  "parent": "96087E8C42DABC25EA5E33FDD4EFED4289F31F4B4826114088527E1F9A735647",
  "tmb": "0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD",
  "typ": "cyphr.me/ac/image/create"
 },
 "sig": "3456AB9532BDC0D9551ECF4D7D109EFA810921C5EB65A3912FD9988A58B3D14066326C7739778D787690BD8A7B88864DFA229A550DC9F39F7FD2729B2F4ED981"
}
```

This all said, there's nothing stopping an application from base64 encoding a
binary and transporting it that way, although we'd recommend against it.  

#### Is Coze versioned?
Our hope is Coze stays simple and stable enough to preclude versioning.  `alg`
refers to a specific set of parameters.  If a parameter needs changing, like the
hashing algorithm, `alg` would reflect that change.  


#### How can my API do versioning?
We suggest incorporating API versioning in `typ`, e.g.
`cyphr.me/v1/msg/create`, where v1 is the api version.


#### Why `pay` and not `head`, `payload`, or `body`?
`pay` is short and denotes inclusion of all fields that are cryptographically
signed.  

Excluding digests, Coze explicitly recommends against including binaries
in JSON messages. A minor concern was that "payload"/"pay" may denote including
arbitrary binary values.


#### Why does `pay` have  cryptographic components?
Coze's `pay` includes all payload information, a design we've dubbed a "fat
payload".  We consider single pass hashing critical for Coze's simple design.

Alternative schemes require a larger canon, `{"head":{...},"pay":{...}}`, or
concatenation like `digest(head) || digest(pay)`.  By hashing only `pay`, the
"head" label and encapsulating parenthathese are dropped, `pay:{...}`, and the
label `"pay"` may then be inferred, `{...}`.  `{...}` is better than
`{"head":{...},"pay":{...}}`.  

Verifying a coze already requires hashing `pay`.  Parsing `alg` from `pay` is a
small additional cost.  


#### I need to keep my JSON separate from Coze.  
We suggest encapsulating your JSON in "~", the last ASCII character.  We've
dubbed this a "tilde encapsulated payload". For example: 

```json
{
  "alg": "ES256",
  "iat": 1623132000,
  "tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
  "typ": "cyphr.me/msg/create",
  "~": {
   "msg": "tilde encapsulated payload"
  }
 }
 ```

#### `key.typ` vs `pay.typ`. 
For signed objects `typ` may be used to denote a canon.  For example, a `typ` with
value `cyphr.me/create/msg` has a canon of ["alg", "iat", "msg", "tmb", "typ"],
as defined by the service.  

For `tmb` canonical form, `typ` is ignored and a static canon is used. Like
`typ` in `pay`,  `typ` in `key` may be used to specify custom application
fields, e.g. "first_seen" or "account_id".  

#### ECDSA `x` and `sig` Bytes
For ECDSA , (X and Y) and (R and S) are concatenated for `x` and `sig`
respectively.  Padding is needed for ES512 because P-521 is rounded up to the
nearest byte before padding.  

#### Javascript Vs Golang Crypto.
Javascript's `SubtleCrypto.sign(algorithm, key, data)` always hashes a message
before signing while Go's ECDSA expects a digest to sign. This means that in
Javascript messages must be passed for signing, while in Go only a digest is
needed.  

See docs/developement.md for the Go development guide.

#### Why not PGP/OpenSSL/LibreSSL/libsodium/JOSE/COSE/etc...?
We have a lot of respect for existing projects. They're great at what they do.
Existing solutions were not meeting our particular needs. Coze is influenced
by many ideas and standards.  
 
See the `coze_vs..md` document for more. 


#### Does Coze have checksums?
For keys, the field `tmb` can the function of a checksum as it is the digest
over `alg` and the key's public components (such as `x` and `y`).  When given a
new public key, systems storing keys can recalculate `tmb`, compare the given
value, and error if values do not match. Alternatively, systems can verify a
signed message with the key.

For messages, `cad`, `czd`, or cryptographic verification may serve the function
of checksums.  

#### Performance hacks?
We did not see prudence in optimizing for long messages, but if early knowledge
of Coze standard fields is critical for application performance, put the Coze
standard fields first, e.g. `{"alg", "tmb", ...}`


#### Why is Coze's scope so limited?
Coze is intentionally scope limited.  It is easier to extend a limited standard
than to fix a large standard. Coze can be extended and customized for individual
applications. 


#### Where does the cryptography come from?
Much of this comes from NIST FIPS (See https://csrc.nist.gov/publications/fips)

For example, FIPS PUB 186-3 defines P-224, P-256, P-384, and P-521.


##### Unsupported Things?
The following are out of scope or redundant.  

- `ES192`, `P-192` - Not implemented anywhere and dropped from later FIPS.
- `SHA1`, `MD5` - Not considered secure for a long time.
- `kty` - "Key type". Redundant by `alg`. 
- `iss` - `tmb` fulfills this role.  Systems that need something like an issuer,
associating messages with people/systems, can look up "issuer" based on
thumbprint.  Associating thumbprints to issuers is the design we recommend.  
- `exp` - "Expiration". Outside the scope of Coze.  
- `nbf` - "Not before". Outside the scope of Coze.  
- `aud` - "Audience". Outside the scope of Coze, but consider denoting this with
  'typ'.
- `sub` - "Subject". Outside the scope of Coze, but consider denoting this with
  'typ'.
- `jti` - "Token ID/JWT ID". Redundant by `czd`, `cad`, or an application
  specified field.

#### Why are duplicate field names prohibited?
Coze explicitly requires that implementations disallow duplicate JSON names in
`coze`, `pay`, and `key`.  Douglas Crockford's Java implementation of JSON
errors on duplicate names. Other implementations use last-value-wins, and a few
support duplicate keys.  The [JSON
RFC](https://datatracker.ietf.org/doc/html/rfc8259#section-4) states that
implementations should not allow duplicate keys, notes the varying behavior of
existing implementations, and states that when names are not unique, "the
behavior of software that receives such an object is unpredictable."  

Duplicate fields is a security issue.  If multiple fields were allowed, for
example for `alg`, `tmb`, or `rvk`, this could be a source of bugs in
implementations and surprising behavior to users. See the article, "[An
Exploration of JSON Interoperability
Vulnerabilities](https://bishopfox.com/blog/json-interoperability-vulnerabilities)"

Javascript objects and Go structs already require unique names.  Since Coze
normalization requires implementations support ordered objects, and prohibiting
duplicates isn't much more complexity.  

# JSON Name, Key, Field Name, Member Name?
They're all synonyms.  A JSON name is a JSON key is a JSON field name is a JSON
member name.  In this document we use "field name" to avoid confusion with Coze
key.  The RFC prefers the terms name/member name, we prefer the term key

## JSON?
- (2017, Bray)      https://datatracker.ietf.org/doc/html/rfc8259
- (2014, Bray)      https://datatracker.ietf.org/doc/html/rfc7159
- (2013, Bray)      https://datatracker.ietf.org/doc/html/rfc7158
- (2006, Crockford) https://datatracker.ietf.org/doc/html/rfc4627

See also I-JSON
 - (2015, Bray)     https://datatracker.ietf.org/doc/html/rfc7493

## Who created Coze?
Coze was created by Cyphr.me.  

## Discussion
https://old.reddit.com/r/CozeJson


## Other Resources
Coze Table links: 
https://docs.google.com/document/d/15_1R7qwfCf-Y3rTamtYS_QXuoTSNrOwbIRopwmv4KOc


## Cyphr.me Online Coze Verifier
https://cyphr.me/coze_verifier



----------------------------------------------------------------------
# Attribution, Trademark notice, and License
Coze is released under The 3-Clause BSD License. 

"Cyphr.me" is a trademark of Cypherpunk, LLC. The Cyphr.me logo is all rights
reserved Cypherpunk, LLC and may not be used without permission.