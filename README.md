[![pkg.go.dev][GoBadge]][GoDoc]
#### ⚠️ Coze is in alpha.  We appreciate feedback and contributions.  Use at your own risk.

![Coze][CozeLogo]


# Coze 
**Coze** is a cryptographic JSON messaging specification.

[Try Coze out!](https://cyphr.me/coze)

[Presentation][Presentation]

### Example Coze
```JSON
{
	"pay": {
		"msg": "Coze Rocks",
		"alg": "ES256",
		"iat": 1623132000,
		"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
		"typ": "cyphr.me/msg"
	},
	"sig": "Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w"
}
```

### Coze Design Goals
1. Idiomatic JSON
2. Human readable
3. Limited scope
4. Providing defined cipher suites

See also [the Coze philosophy](#the-coze-philosophy-of-abstraction)

### Coze Fields
Coze defines standard fields for the objects pay, key, and coze. Applications
may include additional fields as desired. While all fields are optional,
omitting standard fields may limit compatibility. Binary values are encoded as
[RFC 4648 base 64 URI canonical with padding truncated][RFC4648] (b64ut). JSON
components are serialized into UTF-8 for signing, verification, and hashing. All
JSON fields must be unique, and unmarshalling JSON with duplicate fields must
result in an error. All timestamp values are suggested to be Unix time (UTC).

#### All Coze Standard Fields
![Coze Standard Fields](docs/img/coze_standard_fields.png)

## Pay
`pay` contains the fields `alg`, `iat`, `tmb`, and `typ` and optionally any
additional application fields.  In the first example `msg` is additional.

### `pay` Standard Fields
- `alg` - Specific cryptographic algorithm.   E.g. `"ES256"`
- `iat` - Unix time of message signature.     E.g. `1623132000`
- `tmb` - Thumbprint of the signature's key.  E.g. `"cLj8vs..."`
- `typ` - Type of `pay`.                      E.g. `"cyphr.me/msg"`

`typ`'s value may be used by applications as desired.  The value is recommended
to denote API information such as versioning, expected fields, and/or other
application defined programmatic functions.  In the first example,
`"typ":"cyphr.me/msg"` denotes a `pay` with the fields
`["msg","alg","iat","tmb","typ"]` as defined by an application.


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

### `key` Standard Fields
- `key` - Key object.                                     E.g. `"key":{"alg":"ES256", ...}`
- `alg` - Algorithm.                                      E.g. `"ES256"`
- `d`   - Private component.                              E.g. `"bNstg4..."`
- `iat` - "Issued at", Key creation Unix time.            E.g. `1623132000`
- `kid` - "Key identifier", Non-programmatic label.       E.g. `"kid":"My Cyphr.me Key"`. 
- `tmb` - Thumbprint.                                     E.g. `"cLj8vs..."`
- `x`   - Public component.                               E.g. `"2nTOaF..."`.
- `typ` - "Type", Application defined label.              E.g. `"cyphr.me/key"`
- `rvk` - "Revoke", Key revocation Unix time.             E.g. `1623132000`

Note that the private component `d` is not included in `tmb` generation.   Also
note that `kid` must not be used programmatically while `typ` may be used
programmatically.


## Coze object
The JSON name `coze` may be used to wrap a coze.

```JSON
{
	"coze":{
		"pay": {
			"msg": "Coze Rocks",
			"alg": "ES256",
			"iat": 1623132000,
			"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
			"typ": "cyphr.me/msg"
		},
		"sig": "Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w"
	}
}
```

### `coze` Standard Fields
- `coze` "Coze"         Coze object.                E.g. `{"coze":{"pay":..., sig:...}}`
- `can`  "Canon"        Canon of `pay`.             E.g. `["alg","iat","tmb","typ"]`
- `cad`  "Canon digest" Digest of `pay`.            E.g. `"LSgWE4v..."`
- `czd`  "Coze digest"  Digest of `["cad","sig"]`.  E.g. `d0ygwQ...`
- `pay`  "Payload"      Signed payload.             E.g. `"pay":{"alg":...}`
- `sig`  "Signature"    Signature over `cad`.       E.g. `"sig":"ywctP6..."`

`sig` is the signature over the bytes of `cad`.  `cad` is not rehashed before
signing. `czd`'s hashing algorithm must align with `alg` in `pay`.  `czd` refers
to a particular signed message just as `cad` refers to a particular payload.
`cad` and `czd` are calculated from brace to brace, including the braces. `cad`
and `czd` are recalculatable and are recommended to be omitted from cozies,
although they may be useful for reference.

As an added technical constraint, because `sig` and `czd` are used as
identifiers, `sig` must be non-malleable. Malleable schemes like ECDSA must
perform signature canonicalization that constrains signatures to a non-malleable
form.

### Verbose `coze`
Including unnecessary labels is not recommended. For example, the JSON object
`{"pay":{...},"sig":...}` doesn't need the label `coze` if implicitly known by
applications. The following should generally be omitted: `key` may be looked up
by applications by using `tmb`, the fields `can`, `cad`, and `czd` are
recalculatable, and the label `coze` may be inferred.

A tautologic coze:

```JSON
{
	"coze": {
		"pay": {
			"msg": "Coze Rocks",
			"alg": "ES256",
			"iat": 1623132000,
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
		"can": ["msg","alg","iat","tmb","typ"],
		"cad": "Ie3xL77AsiCcb4r0pbnZJqMcfSBqg5Lk0npNJyJ9BC4",
		"czd": "TnRe4DRuGJlw280u3pGhMDOIYM7ii7J8_PhNuSScsIU",
		"sig": "Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w"
	}
}
```

Simplified:

```JSON
{
	"pay": {
		"msg": "Coze Rocks",
		"alg": "ES256",
		"iat": 1623132000,
		"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
		"typ": "cyphr.me/msg"
	},
	"sig": "Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w"
}
```


## Canon
A **canon** is a list of fields used for normalization, e.g. `["alg","x"]`. Coze
objects are canonicalized for creating digests, signing, and verification. The
canon of `pay` is the currently present fields in order of appearance. The
following Coze fields have predefined canons:

- `cad`'s canon is `pay`'s canon.
- `tmb`'s canon is `["alg","x"]`.
- `czd`'s canon is `["cad","sig"]`.

Using a canon, the **canonical form** of an object is generated by removing
fields not appearing in the canon, ordering remaining fields by appearance in
the canon, and eliding unnecessary whitespace.  The canonical form is serialized
into UTF-8 for signing, verification, and hashing.

Canonical form generation steps:

 - Omit fields not present in canon.
 - Order fields by canon.
 - Omit insignificant whitespace.

A **canonical digest** is generated by hashing the UTF-8 serialized canonical
form using the hashing algorithm specified by `alg`.  For example,`"ES256"`'s
hashing algorithm is `"SHA-256"`.

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
 - `pay` is `cad`, 
 - `["alg","x"]` is `tmb`, 
 - `["cad","sig"]` is `czd`.

Using the first example, the following canonical digests are calculated:
- `tmb` is `cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk`
- `cad` is `LSgWE4vEfyxJZUTFaRaB2JdEclORdZcm4UVH9D8vVto`.
- `czd` is `d0ygwQCGzuxqgUq1KsuAtJ8IBu0mkgAcKpUJzuX075M`.

Signing and verification functions must not mutate `pay`.  Any mutation of `pay`
via `can` must occur by canon related functions. Note that's since `pay`'s canon
is the present fields, no fields are removed when canonicalizing `pay`.

### Coze and Binaries
The canonical digest of a binary file may simply be the digest of the file. The
hashing algorithm and any other metadata may be denoted by an accompanying coze.
For example, an image ("coze_logo_icon_256.png") may be referred to by its
digest. 

```JSON
{
	"alg":"SHA-256",
	"file_name":"coze_logo_icon_256.png",
	"id":"oDBDAg4xplHQby6iQ2lZMS1Jz4Op0bNoD5LK3KxEUZo"
}
```

For example, a file's digest, denoted by `id`, may represent the authorization
to upload a file to a user's account.

```JSON
{
 "pay": {
  "alg": "ES256",
  "file_name": "coze_logo_icon_256.png",
  "id": "oDBDAg4xplHQby6iQ2lZMS1Jz4Op0bNoD5LK3KxEUZo",
  "iat": 1623132000,
  "tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
  "typ": "cyphr.me/file/create"
 },
 "sig": "DgJb6Qb81uhC-ulZJlIIj8ahi0b5rAbtnkQhiEH1FB0HeNiACVh_Deo6a22OkK2tr0UcDOiIRY1X-BUriw03Mg"
}
```


## Revoke
A Coze key may be revoked by signing a coze containing the field `rvk` with an
integer value greater than `0`. The integer value `1` is suitable to denote
revocation and the current Unix timestamp is the suggested value.

### Example Self Revoke
```JSON
{
 "pay": {
  "alg": "ES256",
  "iat": 1623132000,
  "msg": "Posted my private key online",
  "rvk": 1623132000,
  "tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
  "typ": "cyphr.me/key/revoke"
 },
 "sig": "KVjPjMVHoL828WyAH5biqIOt-IOaQ5EBtN_7eQifP2w3agUHu6KfqO40_oqQ5GE_BShgXvhbK0O6Z2h5YPNAcw"
}
```

- `rvk` - Unix timestamp of key expiry.

Coze explicitly defines a self-revoke method so that third parties may revoke
leaked keys. Systems storing Coze keys should provide an interface permitting a
given Coze key to be marked as expired by receiving a self-revoke message.
Self-revokes with future times must immediately be considered as expired.

`rvk` and `iat` must be a positive integer less than 2^53 – 1
(9,007,199,254,740,991), which is the integer precision limit specified by
IEEE754 minus one. Revoke checks must error if `rvk` is not an integer or larger
than 2^53 - 1.

Key expiration policies, key rotation, backdating, and alternative revocation
methods are outside the scope of Coze.


## Alg
`alg` specifies a parameter set and is a single source of truth for Coze
cryptographic operations. 

### Example - "alg":"ES256"
- Name:         ES256
- Genus:        ECDSA
- Family:       EC
- Use:          sig
- Hash:         SHA-256
- HashSize:     32 
- HashSizeB64:  43
- XSize:        64
- XSizeB64:     86
- DSize:        32
- DSizeB64:     43
- Curve:        P-256
- SigSize:      64
- SigSizeB64:   86

### Supported Algorithms 
- ES224
- ES256
- ES384
- ES512
- Ed25519 
- Ed25519ph (planned)


## Coze Verifier
The Coze verifier is an in-browser tool for signing and verifying.

[Coze Verifier][Verifier]

![coze_verifier](docs/img/Hello_World!.gif)

There is also the [Simple Coze
Verifier][Verifier_simple] that has the minimal
amount of code needed for a basic Coze application.
Its [codebase is in the Cozejs repo][CozeJSVerifier] and may be locally hosted.


## Coze Implementations
 - [Coze CLI repository][CozeCLI].
 - [Go Coze (this repo)][Coze]
 - [CozeJS (Javascript)][CozeJS]

See [`docs/development.md`](docs/development.md) for the Go development guide.


## Coze Core and Coze X
The sections above are defined as the main Coze specification, Coze core. There
are no plans to increase Coze's scope or features in core other than additional
algorithm support. This will be especially true after Coze is out of Alpha/Beta.
(At the moment, we would like more time for feedback before casting the
specification into stone.)

Coze x (Coze extended) includes additional documentation, extra features,
drafts, proposals, early new algorithms support that's not yet adopted in Coze
core, and extended algorithm support.

See [Coze_go_x/normal][Normal] for an example of a Coze x feature not included in
Coze core.

Repository structure:
- [Coze][Coze]              Main specification (core) and the Go Coze reference implementation.
- [CozeX][CozeX]          Coze extended. Additional documents, discussion, and new algorithms (Not a code repository).
- [CozeGoX][CozeGoX]    Go implementation of extended features.
- [CozeJS][CozeJS]        Javascript implementation of Coze core.
- [CozeJSX][CozeJSX]    Javascript implementation of extended.
- etc...


# FAQ
#### Pronunciation? What does "Coze" mean? 
We say "Co-zee" like a comfy cozy couch.  Jared suggested Coze because it's
funny. The English word Coze is pronounced "kohz" and means "a friendly talk; a
chat" which is the perfect name for a messaging standard.

#### "Coze" vs "coze"?
We use upper case "Coze" to refer to the specification, and "coze"/"cozies" to
refer to messages.

### What is Coze useful for?
Coze's applications are endless as Coze is useful for anything needing
cryptographic signing. Coze is deployed in various applications such as user
authentication (user login), authorization, product tracking, user comments,
user votes, chain of custody, Internet of things (IoT), sessions, and cookies.

As a timely example the CEO of Reddit (reddit.com/u/spez) [edited people's
comments.](https://www.theverge.com/2016/11/23/13739026/reddit-ceo-steve-huffman-edit-comments)
Messages signed by Coze prevents tampering by third parties.


## The Coze Philosophy of Abstraction
Providing a cryptographic abstraction layer is a key feature of Coze. Coze
provides gentle standardization that increases compatibility across various
systems. In this way, Coze is like a simple cryptographic programming language,
allowing projects to "speak the same language". A much larger strategic benefit
is that Coze decouples projects from underlying cryptographic primitives. If
problems are discovered with particular primitives, decoupled architecture is
simple to change while tightly coupled architecture can be extremely difficult
to alter.

A prevalent limitation many cryptographic projects face is that they are
rigidly tied to single primitives. It's hard to overstate the significance of
this problem. 

Coze takes [Bruce Schneier's](https://www.schneier.com/blog/archives/2022/08/nists-post-quantum-cryptography-standards.html) advice seriously:

> It's not enough to implement a single standard; it's vital that our systems be
able to easily swap in new algorithms when required. We've learned the hard way
how algorithms can get so entrenched in systems that it can take many years to
update them: in the transition from DES to AES, and the transition from MD4 and
MD5 to SHA, SHA-1, and then SHA-3. 

The consequences of rigidly tying primitives for cryptographic projects are:

 - **Compatibility** - Projects are frequently incompatible because they use
	 different primitives.
 - **Standardization** - Even when projects use the same primitives,
	 implementations vary widely across the industry, leading to incompatibility
	 between systems.
- **Implementation** - Projects frequently implement primitives directly
	resulting in redundant codebases.
 - **Decoupling** - It is very difficult, if not impossible, for projects to
	 replace broken or problematic primitives.

Coze provides a standardized abstraction layer, eliminating significant
redundant effort of implementation, resolving compatibility conflicts,
establishing implementation standards, and providing the flexibility to use
various primitives as needed.

##### The Debt of Inflexibility
The cost of tight coupling to cryptographic primitives has been demonstrated
repeatedly throughout the industry.  

**Git** - Git was tightly coupled to the questionably secure SHA1 in 2005.
Upgrading to SHA2 has been a herculean effort requiring many years to implement,
which is still ongoing in 2025, 20 years later.

**SSL/TLS** prior to TLS 1.2 - When MD5 and SHA1 weaknesses were discovered, it
required major protocol revisions rather than simple primitive substitution.

**Bitcoin** - Bitcoin is tightly coupled to SHA256, ES256k, and RIPEMD for its
proof-of-work and transaction verification. While ES256k remains secure, if
quantum computers eventually threaten it, transitioning would be extremely
challenging.

**PGP/GPG** - When MD5 was broken, transitioning to newer algorithms required
significant protocol changes and caused compatibility issues between versions.

**DNSSEC** - Transitioning from RSA and SHA1 to newer algorithms like ECDSA and
SHA2 required protocol extensions and complex transition mechanisms.

**WPA2** - WPA2 was tightly coupled to RC4 and AES in specific modes. When
vulnerabilities were found (like KRACK), updating the protocol was complex.

Tight coupling to cryptographic primitives creates significant technical debt
that becomes increasingly expensive to address over time. Quantum computing and
other security threats require the ability to adapt cryptographic systems
quickly. Coze aims to, as an intentional byproduct of standardization, prevent
projects from accumulating this kind of debt, ensuring they remain adaptable to
future cryptographic needs.


#### Binary? Why not support binary payloads?
JSON isn't well designed for large binary payloads.  Instead, Coze suggests
including the digest of a binary file in a coze message while transporting the
binary separately. There's nothing stopping an application from base 64 encoding
a binary for transport, although it's not recommended.

#### How Should I Handle Large Text Messages?
JSON is not ideal for arbitrary text due to escaping, which increases the
size of arbitrary text and makes human readability difficult. Instead of signing
text embedded in a JSON field like `msg`, consider hashing the text, signing its
digest, and transporting the message separately. We consider signing digests
instead of large messages a best practice, although it may add some complexity
to services.  

The following example, which signs a portion of this README, isn't ergonomic.

```JSON
{
"alg": "ES256",
"iat": 1623132000,
"msg": "# Coze \n**Coze** is a cryptographic JSON messaging specification.\n\n[Try Coze out!](https://cyphr.me/coze)",
"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
"typ": "cyphr.me/msg/create"
}
```

The SHA-256 digest (which aligns with an `alg` of ES256) of the message gives
`4FO2pB9yGxo8BBW2whULqbL5m7eAfUWOkvgQu7-9h08`, which is then signed.

```JSON
{
"alg": "ES256",
"iat": 1623132000,
"dig": "4FO2pB9yGxo8BBW2whULqbL5m7eAfUWOkvgQu7-9h08",
"tmb": "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk",
"typ": "cyphr.me/msg/dig/create"
}
```

Like binary, the text is transported independently of the JSON Coze payload.  

```md
# Coze 
**Coze** is a cryptographic JSON messaging specification.

[Try Coze out!](https://cyphr.me/coze)
```

The coze not only digitally signs the message, but also the digest integrity
protects its value.  An additional benefit of this method is that coze
signatures are guaranteed to remain a small, constant size regardless of the
size of the input text.  

#### Why is Coze's scope so limited?
Coze is intentionally scope limited.  It is easier to extend a limited standard
than to fix a large standard. Coze can be extended and customized for individual
applications. 

#### Is Coze versioned?
`alg` refers to a specific set of parameters for all operations and Coze Core
"versioning" is accomplished by noting specific algorithm support.  If an
operation needs a different parameter set, `alg` itself must denote the
difference.  `alg` permits Coze implementations to support a subset of features
while remaining Coze compliant.  The specification hopes to stay simple and
stable enough to preclude versioning, however we suspect further tweaks are
probably warranted, so a long alpha and beta time is planned.  Extension to Coze
are defined by [CozeX][Cozex] so implementations avoid feature bloat.
Implementation releases themselves are versioned. 

#### Why does `pay` have cryptographic components?
Coze's `pay` includes all payload information, a design we've dubbed a "fat
payload".  We consider single pass hashing critical for Coze's simple design.

Alternative schemes require a larger canon, `{"head":{...},"pay":{...}}`, or
concatenation like `digest(head) || digest(pay)`.  By hashing only `pay`, the
"head" label and encapsulating braces are dropped, `pay:{...}`, and the label
`"pay"` may then be inferred, `{...}`.  `{...}` is better than
`{"head":{...},"pay":{...}}`.

Verifying a coze already requires hashing `pay`.  Parsing `alg` from `pay` is a
small additional cost.

#### JSON APIs?  Can my API do versioning?
Coze is well suited for JSON APIs. API versioning may be handled by applications
however desired.  A suggested way of incorporating API versioning in Coze is to
use `typ`, e.g. `"typ":"cyphr.me/v1/msg/create"`, where "v1" is the api version.

#### Can my application use Canon/Canonicalization?
Yes, canon is suitable for general purpose application.  Applications may
specify canon expectations in API documentation, if using Coze denoted by "typ"
or explicitly specified by `can`, or implicitly known and pre-established.  Coze
Core contains simple canonicalization functions, or for more expressive
capabilities see [Normal][Normal].

#### `pay.typ` vs `key.typ`. 
For applications, `pay.typ` may denote a canon.  For example, a `typ` with value
`cyphr.me/msg/create` has a canon, as defined by the service, of ["alg", "iat",
"msg", "tmb", "typ"].  The service may reject a coze that's not canonicalized as
expected.  For example, the service might reject cozies missing `iat`.

Like `typ` in `pay`, applications may use `key.typ` to specify custom fields
(e.g., "first_seen" or "account_id") and field order.

`Key.tmb` ignores `key.typ` because `alg` serves as the key's `typ` so the
static canon, `["alg","x"]`, is sufficient.  Using `alg` in the generation of
`tmb` ensures the impossibility of algorithms producing colliding thumbprints
(where one algorithm could produce `x` values colliding with other algorithms).  

#### ECDSA `x` and `sig` Bytes.
For ECDSA , (X and Y) and (R and S) are concatenated for `x` and `sig`
respectively.  For ES512, which unlike the other ECDSA algorithms uses the odd
numbered P-521, X, Y, R, and S are padded before concatenation.

#### Why use `tmb` and not `x` for references in messages?
Coze places no limit on public key size, which can be very large. For example,
GeMSS128 public keys are 352,188 bytes, compared to Ed25519's 32 bytes.  Using
`tmb` instead of `x` generalizes Coze for present and future algorithm use.
Additionally, `x` may be cryptographically significant for key security while
`tmb` is not.

####  Required Coze Fields, Contextual Cozies, and the Empty Coze.
The standard fields provide Coze and applications fields with known types since
JSON has limited type identifiers. Coze has no required fields, however omitting
standard fields limits interoperability among applications, so it is suggested
to include standard fields appropriately.

Cozies that are missing the fields `pay.alg` and/or `pay.tmb` are **contextual
cozies**, denoting that additional information is needed for verification.
Caution is urged when deploying contextual cozies as including the standard
fields `pay.alg` and `pay.tmb` is preferred.

An **empty coze**, which has an empty `pay` and populated `sig`, is legitimate.
It may be verified if `key` is known. The following empty coze was signed with
the example key "cLj8vs".

```json
{
	"pay":{},
	"sig":"9iesKUSV7L1-xz5yd3A94vCkKLmdOAnrcPXTU3_qeKSuk4RMG7Qz0KyubpATy0XA_fXrcdaxJTvXg6saaQQcVQ"
}
```

#### UTF-8 and b64ut (RFC base 64 URI canonical truncated) Encoding
[Canonical base 64][RFC6468Canonical] (sometimes called "strict") encoding is
required and non-strict encoding of both b64ut and UTF-8 must error.   For the
initial reason for why Coze uses b64ut see [base64.md][base64.md].

#### Why not PGP/OpenSSL/LibreSSL/SSHSIG/libsodium/JOSE(JWT)/COSE/etc...?  How does Coze compare with prior arts?
We respect the various projects in the space.  Other projects have noble goals
and we're thankful they exist.  Coze is influenced by ideas from many others.
However existing solutions were not meeting our particular needs so we created
Coze.
 
See [coze_vs.md][coze_vs] and the [introduction
presentation](https://docs.google.com/presentation/d/1bVojfkDs7K9hRwjr8zMW-AoHv5yAZjKL9Z3Bicz5Too/edit#slide=id.g1367bc4eb0f_0_6)
for more. 

#### Does Coze have checksums?
`x`, `tmb`,`cad`, `czd`, and `sig` may be used for integrity checking.

Systems may use `sig` as an integrity check via cryptographic verification.  If
`cad` and/or `czd` are included they may be recalculated and error on mismatch.

For keys, `x` and/or `tmb` may be recalculated and error on mismatch.Coze keys
cannot be integrity checked when `d`, `x`, or `tmb` are presented alone.  In
situations needing integrity checking, we recommend including at least two
components.  See [checksums.md][checksums] for more.

#### Performance hacks?
Coze is not optimized for long messages, but if early knowledge of Coze standard
fields is critical for application performance, put the Coze standard fields
first, e.g. `{"alg", "tmb", ...}`

#### I need to keep my JSON separate but inside a coze.
If appending custom fields after the standard Coze fields isn't sufficient, we
suggest encapsulating custom JSON in "~", the last ASCII character.  We've
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

#### ASCII/Unicode/UTF-8/UTF-16 and Ordering?
Even though Javascript uses UTF-16 and JSON was designed in a Javascript
context, JSON implementations rejected the problematic UTF-16, which has some
code points out of order, in favor of UTF-8.  Requiring JSON UTF-8 encoding was
formalized by the [JSON RFC 8259 section 8.1][RFC8259-8.1].  Unicode, ASCII, and
UTF-8 all share sorting order.  

Although JSON arrays are defined as ordered, JSON objects are defined as
unordered.  How is `pay`, an unordered JSON object, signed when signing requires
a static representation? [UTF-8 is the explicitly defined serialization for
JSON][RFC8259-8.1]. Coze's signing and verification operations are not over
abstract JSON, but rather the concrete UTF-8.  Coze marshals JSON into UTF-8
before signing, and Coze verifies UTF-8 before unmarshalling into JSON.  

Additionally, object field order may be denoted by `can`, [chaining
normals][Normal], or communicate via other means. 

#### Where does the cryptography come from?
Much of this comes from [NIST FIPS][FIPS].

For example, FIPS PUB 186-3 defines P-224, P-256, P-384, and P-521.

To learn more see this [walkthrough of ECDSA](https://learnmeabitcoin.com/technical/ecdsa). 

#### Unsupported Things?
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
  `typ`.
- `sub` - "Subject". Outside the scope of Coze, but consider denoting this with
  `typ`.
- `jti` - "Token ID/JWT ID". Redundant by `czd`, `cad`, or an application
  specified field.

#### Encryption?
Coze does not currently support encryption.  If or when it ever does it would be
similar to or simply complement [age](https://github.com/FiloSottile/age).

#### Why define algorithms?
Coze's design is generalized and not overly coupled to any single primitive.
Because of this, applications that use Coze can easy upgrade cryptographic
primitives. Using a single primitive is perfectly fine, but tightly coupling
systems to a single primitive is not.  Simultaneous support for multiple
primitives is a secondary, and optional, perk.

#### JSON "Name", "Key", "Field Name", "Member Name"?
They're all synonyms.  A JSON name is a JSON key is a JSON field name is a JSON
member name.  In this document we use "field name" to avoid confusion with Coze
key.

#### Why are duplicate field names prohibited?
Coze explicitly requires that implementations disallow duplicate field names in
`coze`, `pay`, and `key`.  Existing JSON implementations have varying behavior.
Douglas Crockford, JSON's inventor, [tried to fix this but it was decided it
was too late](https://esdiscuss.org/topic/json-duplicate-keys).

Although Douglas Crockford couldn't change the spec forcing all implementations
to error on duplicate, his Java JSON implementation errors on duplicate names.
Others use `last-value-wins`, support duplicate keys, or other non-standard
behavior. The [JSON
RFC](https://datatracker.ietf.org/doc/html/rfc8259#section-4) states that
implementations should not allow duplicate keys, notes the varying behavior
of existing implementations, and states that when names are not unique, "the
behavior of software that receives such an object is unpredictable."  Also note
that Javascript objects (ES6) and Go structs already require unique names.

Duplicate fields are a security issue, a source of bugs, and a surprising
behavior to users. See the article, "[An Exploration of JSON Interoperability
Vulnerabilities](https://bishopfox.com/blog/json-interoperability-vulnerabilities)"

Disallowing duplicates conforms to the small I-JSON RFC. The author of I-JSON,
Tim Bray, is also the author of current JSON specification ([RFC
8259][RFC8259]).  See also https://github.com/json5/json5-spec/issues/38.

#### Why is human readability a goal?
Although humans cannot verify a signature without the assistance of tools,
readability allows humans to visually verify what a message does.

We saw the need for JSON-centric cryptography and idiomatic JSON is human
readable.  JSON is not a binary format; it is a human readable format and any
framwork built on JSON should embrace its human readability.  If human
readability is unneeded, JSON is entirely the wrong message format to employ.
All else being equal, human readability is better than non-human readability.

#### JSON?
- [RFC 8259 (2017, Bray)][RFC8259]
- [RFC 7159 (2014, Bray)][RFC7159]
- [RFC 7158 (2013, Bray)][RFC7158]
- [RFC 4627 (2006, Crockford)][RFC4627]

See also I-JSON and JSON5
 - [RFC 7493 (2015, Bray)][RFC7493]
 - [JSON5][JSON5]

#### HTTP?  HTTP Cookies?  HTTP Headers?
When using Coze with HTTP cookies, Coze messages should be JSON minified.  For
example, we've encountered no issues using the first example as a cookie:

```
token={"pay":{"msg":"Coze Rocks","alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w"}; Path=/;  Secure; Max-Age=999999999; SameSite=None
```

For more considerations see [http_headers.md][http_headers]

#### Why release pre-alpha on 2021/06/08?
Coze was released on 2021/06/08 (1623132000) since it's 30 years and one day
after the initial release of PGP 1.0.  We wrote a blog with more [details of
Coze's
genesis](https://cyphr.me/md/xlA-MSFwPxmWED4ZcNdxUy8OA22UPiLWlGQUQik8DwY).

#### Signature Malleability?
Coze prohibits signature malleability.  See
[malleability_low_s.md][low_s].

#### Who created Coze?
Coze was created by [Cyphr.me](https://cyphr.me).

#### Discussion?  Social Media?
 - We have a bridged Matrix and Telegram chat room.  (This is where we are the most active) 
	- Matrix: https://app.element.io/#/room/#cyphrmepub:matrix.org
	- PM zamicol for our bridged Telegram group.
 - We also hang out in the Go rooms:
   - https://app.element.io/#/room/#go-lang:matrix.org
   - https://t.me/+TgkdqZw0Q-jAkGWS
 - https://twitter.com/CozeJSON
 - https://old.reddit.com/r/CozeJson

#### Other Resources
 - This README as a page: https://cyphrme.github.io/Coze/
 - [Coze go.pkg.dev](https://pkg.go.dev/github.com/cyphrme/coze#section-readme)
 - CozeJSON.com (which is currently pointed to the [Coze verifier](https://cyphr.me/coze))
 - Coze Table links: https://docs.google.com/document/d/15_1R7qwfCf-Y3rTamtYS_QXuoTSNrOwbIRopwmv4KOc

#### Keywords
Coze JSON alg iat tmb typ rvk kid d x coze pay key can cad czd sig cryptography crypto authentication auth login hash digest signature Cypherpunk Cyphrme Ed25519 Ed25519ph ES224 ES256 ES384 ES512 SHA-224 SHA-256 SHA-384 SHA512 JOSE JWS JWE JWK JWT PASETO PASERK signify ssh SSHSIG PGP Bitcoin Ethereum base64 b64ut SQRL


----------------------------------------------------------------------
# Attribution, Trademark Notice, and License
Coze is released under The 3-Clause BSD License.

"Cyphr.me" is a trademark of Cypherpunk, LLC. The Cyphr.me logo is all rights
reserved Cypherpunk, LLC and may not be used without permission.


[GoBadge]:  https://pkg.go.dev/badge/github.com/github.com/cyphrme/coze
[GoDoc]:    https://pkg.go.dev/github.com/cyphrme/coze
[CozeLogo]: docs/img/coze_logo_zami_white_450x273.png


[Presentation]:       https://docs.google.com/presentation/d/1bVojfkDs7K9hRwjr8zMW-AoHv5yAZjKL9Z3Bicz5Too

[Verifier]:           https://cyphr.me/coze
[Verifier_simple]:    https://cyphr.me/coze_verifier_simple/coze.html
[CozeJSVerifier]:     https://github.com/Cyphrme/Cozejs/tree/master/verifier
[GithubCozeVerifier]: https://cyphrme.github.io/Cozejs/verifier/coze.html

[Coze]:      https://github.com/Cyphrme/Coze
[CozeCLI]:   https://github.com/Cyphrme/CozeCLI
[CozeX]:     https://github.com/Cyphrme/CozeX
[CozeGoX]:   https://github.com/Cyphrme/CozeGoX
[CozeJS]:    https://github.com/Cyphrme/CozeJS
[CozeJSX]:   https://github.com/Cyphrme/CozeJS

[Normal]:           https://github.com/Cyphrme/Coze_go_x/tree/master/normal
[checksums]:        https://github.com/Cyphrme/Coze_x/blob/master/proposal/checksum.md
[coze_vs]:          https://github.com/Cyphrme/Coze_x/blob/master/coze_vs.md
[http_headers]:     https://github.com/Cyphrme/Coze_x/blob/master/http_headers.md
[low_s]:            https://github.com/Cyphrme/Coze_x/blob/master/implemented/malleability_low_s.md
[base64.md]:        https://github.com/Cyphrme/Coze_x/blob/master/implemented/base64.md
[RFC6468Canonical]: https://datatracker.ietf.org/doc/html/rfc4648#section-3.5


[RFC4648]: https://datatracker.ietf.org/doc/html/rfc4648

[RFC8259]:     https://datatracker.ietf.org/doc/html/rfc8259
[RFC8259-8.1]: https://datatracker.ietf.org/doc/html/rfc8259#section-8.1
[RFC7159]:     https://datatracker.ietf.org/doc/html/rfc7159
[RFC7158]:     https://datatracker.ietf.org/doc/html/rfc7158
[RFC4627]:     https://datatracker.ietf.org/doc/html/rfc4627
[RFC7493]:     https://datatracker.ietf.org/doc/html/rfc7493
[JSON5]:       https://github.com/json5/json5-spec
[FIPS]:        https://csrc.nist.gov/publications/fips





