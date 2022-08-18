# ⚠️ COZE IS IN ALPHA.  USE AT YOUR OWN RISK ⚠️
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/github.com/cyphrme/coze)](https://pkg.go.dev/github.com/cyphrme/coze)

![Coze](docs/img/coze_logo_zami_white_450x273.png)

[Presentation](https://docs.google.com/presentation/d/1bVojfkDs7K9hRwjr8zMW-AoHv5yAZjKL9Z3Bicz5Too)

# Coze 
**Coze** is a cryptographic JSON messaging specification designed for human
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

### Coze Fields
Coze objects encapsulate a set of JSON name/value fields.  Coze JSON objects are
case sensitive, must be valid JSON, and must contain unique field names. Coze
**reserved fields** must be used according to Coze.  Applications are permitted
to use additional fields as desired.  All reserved fields are optional, but
omitting standard fields may limit compatibility among Coze supporting
applications.  Binary values are encoded as RFC 4648 base64 URI with padding
omitted.  The Coze objects `pay`, `key`, and `coze` have respective reserved
fields.

![Coze Reserved Fields](docs/img/coze_reserved_fields.png)

## Pay
`pay` may contain the standard fields `alg`, `iat`, `tmb`, and `typ` and
additional fields.  In the first example, `msg` is additional.

### `pay` Reserved Names
- `alg` - Specific cryptographic algorithm.  E.g. `"ES256"`
- `iat` - The time when the message was signed. E.g. `1623132000`
- `tmb` - Thumbprint of the key used to sign the message.  E.g. `"0148F4..."`
- `typ` - Type of `pay`.

`typ`'s value may be used by applications as desired.  The value is recommended
to denote API information such as versioning, expected fields, and/or other
application defined programmatic functions.  In the first example,
`"typ":"cyphr.me/msg"` denotes a `pay` with the fields
`["msg","alg","iat","tmb","typ"]` as defined by a hypothetical application.  

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

### `key` Reserved Names
- `key` - Key object.  E.g. `"key":{"alg":"ES256", ...}`
- `alg` - Algorithm.  E.g. `"ES256"`
- `d`   - Private component.  E.g. `"bNstg4..."`
- `iat` - "Issued at", When the key was created.  E.g. `1623132000`
- `kid` - "Key identifier", Human readable, non-programmatic label.  E.g. `"kid":"My Cyphr.me Key"`. 
- `tmb` - Thumbprint.  E.g. `"cLj8vs..."`
- `x`   - Public component.  E.g. `"2nTOaF..."`.
- `typ` - "Type", Additional application information.  E.g. `"cyphr.me/msg"`
- `rvk` - "Revoke", time of key revocation.  See the `rvk` section.  E.g. `1655924566`

Note that the private component `d` is not included in `tmb` generation.   Also
note that `kid` must not be used programmatically while`typ` may be used
programmatically. 


## `coze` Reserved Names
- `coze` - JSON name for Coze objects.  E.g. `{"coze":{"pay":..., sig:...}}`
- `can` - "Canon" of `pay`.  E.g. `["alg","iat","tmb","typ"]`
- `cad` - "Canon digest", the digest of `pay`.  E.g.: `"LSgWE4v..."`
- `czd` - "Coze digest", the digest over `["cad","sig"]`.  E.g. `d0ygwQ...`
- `pay` - Label for the pay object.  E.g. `"pay":{"alg":...}`
- `sig` - Signature over `cad`.  E.g. `"sig":"ywctP6..."`

`sig` is the signature of the bytes represented by `cad` and `cad` is not
rehashed before signing. `czd`'s hashing algorithm must align with `alg` in
`pay`.  `czd` refers to a particular signed message. Like `cad`, `czd` is
calculated from brace to brace, including the braces. `cad` and `czd` are
recalculatable and are recommended to be omitted, although they may be useful
for reference.  


## Coze Labels
The JSON name `coze` may be used to wrap Coze objects.  For example:

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

It is recommend to not needlessly wrap Coze objects with labels. For example,
the JSON object `{"pay":{...},"sig":...}` doesn't need the labeled `coze` if
implicitly known by applications.

The following coze expands the first example by adding the labels `key`, `can`,
`cad`, and `czd` that should generally be omitted unless needed by applications.
`key` may be looked up by applications by using `tmb`, `can`, `cad`, and `czd`
are recalculatable, and the label `coze` may be inferred.  

The tautological coze

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

simplifies to

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

### Canon
A **canon** is a list of fields used for normalization, e.g. `["alg","x"]`.  Some
Coze objects are canonicalized for creating digests, signing, verification, and
reference. Using a canon, the **canonical form** of an object is generated by
removing fields not appearing in the canon, ordering remaining field by
appearance in the canon, and eliding unnecessary whitespace.

Generation steps for the canonical form:
 0. Omit fields not present in canon.
 1. Order fields by canon.
 2. Omit insignificant whitespace.

The following Coze fields have predefined canons:  
- `tmb`'s canon is `["alg","x"]`.
- `cad`'s canon is `pay`'s fields in order of appearance.
- `czd`'s canon is `["cad","sig"]`.

A **canonical digest** is generated by hashing the canonical form using the
hashing algorithm specified by `alg`.  For example,`"ES256"`'s hashing algorithm
is `"SHA-256"`.

The key thumbprint, `tmb`, is the canonical digest of `key` using the canon
`["alg","x"]` and hashing algorithm specified by `key.alg`.  For example, a key
`alg` of `ES256` corresponds to the hashing algorithm `SHA-256`. The canonical
form of the example key is:

```JSON
{"alg":"ES256","x":"2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g"}
```

Hashing this canonical form results in the following digest, which is `tmb`:
`cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk`. 

`czd` is the canonical digest of `coze` with the canon `["cad","sig"]`, which
results in the JSON `{"cad":"...",sig:"..."}`.  `czd`'s hash must align with
`alg` in `pay`. 

The canonical digest of 
 - `key` is `tmb`, 
 - `pay` is `cad`, 
 - `["cad","sig"]` is `czd`.

Using the first example, the following canonical digests are calculated:
- `tmb` is `cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk`
- `cad` is `LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto`.
- `czd` is `d0ygwQCGzuxqgUq1KsuAtJ8IBu0mkgAcKpUJzuX075M`.


### Coze and Binaries
The canonical digest of binary files may simply be the digest of the file. The
hashing algorithm and any other metadata may be denoted by an accompanying coze.
For example, an image ("Hello_World!.gif") may be referred to in a JSON object
by its digest. 

```JSON
{
	"alg":"SHA-256",
	"file_name":"Hello_World!.gif",
	"image":"rVOyJ144KwIQ3V2YJdatKAo_3QWAY4CpGLCDdnKOvAw"
}
```

For example, including a file's digest in a signed message, denoted by `id`, may
represent the authorization to upload a file to a user's account:

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
A Coze key may be revoked by signing a self-revoke coze.  A self-revoke coze has
the field `rvk` with any value other than `0`.  For example, the integer value
`1` is suitable to denote revocation.  A Unix timestamp of the time of
revocation is the suggested value for `rvk`.

### Example Self Revoke

```JSON
{
	"pay": {
		"alg": "ES256",
		"iat": 1655924566,
		"msg": "Posted my private key on github",
		"rvk": 1655924566,
		"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
		"typ": "cyphr.me/key/revoke"
	},
	"sig": "y3wpVXpBeaJNnUn8Q_3j9WOZH4gey78naDrP14TEToio0tloGP-6mNrXGQdWsvMvVYgg09EoxJYC9mE4PEuMXg"
}
```

- `rvk` - Unix timestamp of when the key was revoked.  

Coze explicitly defines a self-revoke method so that third parties may revoke
leaked keys. Systems storing Coze keys should provide an interface permitting a
given Coze key to be mark as revoked by receiving a self-revoke message.
Self-revokes with future times must immediately be considered as revoked.  Coze
suggests rvk to be a 32-bit unsigned integer with a maximum value of
4,294,967,295.

Key expiration policies, such as key rotation, are outside the scope of Coze.


## Supported Algorithms 
- ES224
- ES256
- ES384
- ES512
- Ed25519 
- Ed25519ph (planned)

### `alg` parameters:
`alg` is a single source of truth for Coze cryptographic operations.  Other
parameters are derived from the value of `alg`. For example:  

"alg":"ES256"
- Genus: ECDSA
- Family: EC
- Use: sig
- Sig.Size: 512
- Curve: P-256
- Hash: SHA-256
- Hash.Size:256 


## Coze Verifier
Cyphr.me provides an online tool for signing and verify Coze messages and plans
to release an open source, stand alone version of the webpage.  

Play with Coze here: https://cyphr.me/coze_verifier.

![coze_verifier](docs/img/Hello_World!.gif)


## Coze Implementations
 - [Go Coze (this repo)](https://github.com/Cyphrme/coze)
 - [Coze js (Javascript)](https://github.com/Cyphrme/cozejs)



# Standard Coze
The sections above are defined as "Core Coze".  Further expansions on Coze may
be included in "Coze Standard".  Further draft, proposals, and extended
algorithm support are planned in "Coze Experimental".

See `normal.go` for an example of a Coze Standard feature not included in Core
Coze.  


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
`alg` refers to a specific set of parameters for all operations.  If a parameter
needs changing, like switching out a hashing algorithm, `alg` must reflect that
change.  

Our hope is "Coze Core" stays simple and stable enough to preclude versioning.
Instead, Coze Core "versioning" will be accomplished by noting specific
algorithm support.  

Versioning by feature also permits Coze implementations to support a subset of
features while remaining Coze compliant.  In this way libraries may remain
spartan and avoid feature bloat.

Further expansions on Coze may be included in "Coze Standard".  Further draft,
proposals, and extended algorithm support are planned in "Coze Experimental".


#### How can my API do versioning?
API versioning may be handled an application however desired.  A suggested way
of incorporating API versioning in Coze is to use `typ`, e.g.
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


#### Can my application use Canon/Canonicalization?
Applications may find it useful to have messages in a specific normalized form.
Core Coze has canonicalization features that may be used, or for more expressive
capabilities, see Normal in Coze Standard.  

Canon may be implicitly known by applications, implicitly derived by "typ", or
explicitly specified by `can`. Applications may specify canon expectations in
API documentation.  If a message is malformed, applications must error. 

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
For signed objects `typ` may be used to denote a canon.  For example, a `typ`
with value `cyphr.me/create/msg` has a canon of ["alg", "iat", "msg", "tmb",
"typ"], as defined by the service.  

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
Coze is not optimized for long messages, but if early knowledge of Coze standard
fields is critical for application performance, put the Coze standard fields
first, e.g. `{"alg", "tmb", ...}`


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

#### JSON Name, Key, Field Name, Member Name?
They're all synonyms.  A JSON name is a JSON key is a JSON field name is a JSON
member name.  In this document we use "field name" to avoid confusion with Coze
key.  The RFC prefers the terms name/member name, we prefer the term key


#### Cryptographic Agility?
The moral is the need for cryptographic agility. It’s not enough to implement a
single standard; it’s vital that our systems be able to easily swap in new
algorithms when required. We’ve learned the hard way how algorithms can get so
entrenched in systems that it can take many years to update them: in the
transition from DES to AES, and the transition from MD4 and MD5 to SHA, SHA-1,
and then SHA-3.

- https://www.schneier.com/blog/archives/2022/08/nists-post-quantum-cryptography-standards.html


#### JSON?
- (2017, Bray)      https://datatracker.ietf.org/doc/html/rfc8259
- (2014, Bray)      https://datatracker.ietf.org/doc/html/rfc7159
- (2013, Bray)      https://datatracker.ietf.org/doc/html/rfc7158
- (2006, Crockford) https://datatracker.ietf.org/doc/html/rfc4627

See also I-JSON
 - (2015, Bray)     https://datatracker.ietf.org/doc/html/rfc7493

#### Who created Coze?
Coze was created by Cyphr.me.  

#### Discussion?
https://old.reddit.com/r/CozeJson

PM zamicol for our telegram group.  


#### Other Resources
Coze Table links: 
https://docs.google.com/document/d/15_1R7qwfCf-Y3rTamtYS_QXuoTSNrOwbIRopwmv4KOc




----------------------------------------------------------------------
# Attribution, Trademark notice, and License
Coze is released under The 3-Clause BSD License. 

"Cyphr.me" is a trademark of Cypherpunk, LLC. The Cyphr.me logo is all rights
reserved Cypherpunk, LLC and may not be used without permission.