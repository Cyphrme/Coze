# Coze 
Adorable and deplorable. 
(Pronounced "Co-zee" like a comfy cozy couch.)

JOSE was proving cumbersome.  Coze is another design iteration on the basic
ideas of JOSE.
 
Coze Table links: 
https://docs.google.com/document/d/15_1R7qwfCf-Y3rTamtYS_QXuoTSNrOwbIRopwmv4KOc

We decided to release this on 2021/06/08 since it's been 30 years and one day
since the release of PGP 1.0.  

## What does "Coze" mean? 
COSE was already taken so we couldn't name it "Cyphr OSE" (COSE, Cyphr Object
Signing and Encryption) and CJOSE seemed verbose.  Jared suggested Coze because
it sounds like JOSE but it's funnier. 
 
Also, "coze" means "a friendly talk; a chat." which we found perfect for the
name of a messaging standard. 
 
## Why Hex? 
Hex support is universal and easy to implement. Coze requires that byte
information is always represented as lower case hex. This allows for easy
storing of information in the database as bytes. 
 
# UTF-8
Coze uses UTF-8 encoding for everything.  This includes how keys in objects are
ordered.  

# Coze Object
## Simple Signed Coze Object:
```json
{
 "head":{
   "act":"/ac/upload/",
   "alg":"ES256",  
   "cd":"C32365bb3c70eb559dfffc902849442f343bd431d9f322d6a50a43e863686d2bd", 
   "iat":1621893376,
   "sth":"SHA-256",  
   "std":"32365bb3c70eb559dfffc902849442f343bd431d9f322d6a50a43e863686d2bd",
   "typ":"cyphr.me" 
 },
 "sig":"32365bb3c70eb559dfffc902849442f343bd431d9f322d6a50a43e863686d2bd32365bb3c70eb559dfffc902849442f343bd431d9f322d6a50a43e863686d2bd"
}
```


## Advanced Signed Coze Object
```json
{
  "head": {
    "act": "/cyphr/sign/message",
    "alg": "ES256",
    "iat": 1623194836589,
    "msg": "My First Coze Message!",
    "sth": "SHA-256",
    "std": "252ba5916f8665ced714121d9d5f4d3c971e46ed5ea86e04762c1a5be5ffe660",
    "typ": "cyphr.me"
  },
  "pubkey": {
    "x": "e50b9a5b3b4318ab21be393903f69b21ea473d9e0b09071a26b37a935e8b3f05",
    "y": "68cbe808ee968dba53f59d6047729ee3196f84cc04c04dfb6cf695778e8c3ca4",
    "use": "sig",
    "alg": "ES256",
    "iat": 1623170221,
    "tmb": "252ba5916f8665ced714121d9d5f4d3c971e46ed5ea86e04762c1a5be5ffe660",
    "kid": "My Cyphr.me Key.",
    "th": "SHA-256"
  },
  "sig": "8d1e74ef5700759a5561941c4170346c12023e43df57024d187b473def54425e9146c622818f9c1deffc8062f39f65997baccc9b1dbcb7cb31e5f3336eea27c8"
}
```


## Coze Object Reserved Keys
A Coze object is a set of Key:Values.  Keys may be anything the application uses, except for the reserved keys defined in the Coze standard.  

For application specific keys, we suggest namespacing for application specific useage.  (Example: "myapp_UserName")

`*` is for existing JOSE keys, although their usages may be different from JOSE's. 
 
```
*alg: (required) Algo of signing/verification/encryption/key algo.  "alg" is also required in keys.  Example: "ES256". 

 
std:  (required) "Signing Thumbprint Digest" thumbprint of the key used to sign the message.  Example: "ba9ed05fa4e8b4c3f911427914582fb481e4803a940e954890dff87cb5d86cfe"
sth:  (required) "Signing Thumbprint Hashing Algorithm" Hashing algorithm used to create the std.  Example: "SHA-256"
 
*iat: (required) "Issued at" // In place of "created"
*typ: (Optional)"cyphr.me/c" // Custom Cyphr.me headers.  "typ" specifies what headers should appear in the object.  

can   (optional) Canon for hashing over object.  Example: ["alg","iat","std","sth","typ"]

# For use outside the "head" object
sig    Sig in bytes  Example: sig:"60c74ce[...]"
sigs  (Encapsulation) One or many sigs with crypto variables included.  Example: sigs[...]" // TODO
cy (Encapsulation) Encapsulation for signed Coze object.  Example: cy:{"head":{"alg":"ES256", [...]}, sig:"60c74ce[...]"}
```

## Signature Bytes
The signature bytes must be the concated, and if needed padded, bytes of the signature.  



# Coze Key
# Example Coze Key
```json
"key":{
 "alg": "ES256",
 "d":   "e2bb2f49dd86dc2796eeb45d6a0d5520cf076a604487e7868c6398877c82282d",
 "kid": "Example Coze Key",
 "iat": 1623132000,
 "tmb": "926b006acd3c2687a652d0f68e8d4929331dea1d57315592c8353c89d87fbe29",
 "th":  "SHA-256",
 "use": "sig",
 "x":   "7bc4803828260c74ce05c83cce3539dd681697624a523b50a14f736c13dbad4c",
 "y":   "9cfff8d73cef3eb6edaf217a248550159a35e3f865eef5b4e07acde5ab1331c2",
}
```
 
Coze keys must include `alg`, `iat`, `tmb`, `th`, `use`.  For specific `alg`'s,
other keys may be required.  

- `tmb` is generated from the canon fields.  
- `kid` is optional and is the human readable label for the key. 
- `crv` and `kty` are optional and redundant parameters. 
- If `crv` is used, it must be a specific curve (Curve25519) and not an algorithm (Ed25519).  

## Coze Key Thumbprint Canon:
For ecdsa: ["alg","x","y"]
 
## Coze Key Required Fields:
For ecdsa public key:  ["alg", "iat", "tmb", "th", "use", "x", "y"]
For ecdsa private key: ["alg", "d", "iat", "tmb", "th", "use", "x", "y"]

## Coze Key Keys
 * is for existing JOSE keys, although their usages may be different from JOSE. 

```
*alg: (required) Specific algorithm for signing/verification/encryption/key.  "alg" is required in keys.  Example: "ES256" or "Ed25519".
*d    (required) Private component of ECDSA.  Example: "f3acb690c2e35b517c012135114b88e22b73042ac12e0229e94738679782145a".
*kid  (optional) Human Readable identifier for the key. Example: "My Cyphr.me Key". 
*tmb  (required) Thumbprint of a key inside of a key object.  When referring to a key outside of a key object, use "std" instead.  Example:
         "ba9ed05fa4e8b4c3f911427914582fb481e4803a940e954890dff87cb5d86cfe"
th    (required) Hash algorithm used used to create the thumbprint.  Example: "SHA-256". 
*use  (required) As of 2021/05/13, only "sig" or "enc" are valid. Example: "sig".
*x    (required) X coordinate.  Example: "a2f637239406c8528eff61047fb232e902dc4c513e01544c9358f27d18c0bdd8".
*y    (required) Y coordinate.  Example: "4b70d24c1eca9060609ddd39fa76870aa7ea2822bb0475f3516af2f8a5290ba7".

# Optional And Redundant
key   (Encapsulation) Encapsulation of a key object.   Its value must be an object. 
        Example: "key:{"alg": "ES256", ...}
kpt   (Optional)  Only "pubkey" and "pairkey" are currently valid.  Value 
        "pubkey"  Must not have private components.  Must be the public key.  
        Value "pairkey" must be the full keypair.   Must have the whole private and public components. 

*kty (optional and redundant) Key Type.  For the elliptic curve this is "EC" and for RSA this is "RSA". See notes about "kty".  Example: "EC"
*crv (optional and redundant) Named curve used in the elliptic curve scheme.  See notes about "crv". Example: "P-256" or "Curve25519"

# Forbidden
key_ops (forbidden)  Use "use" instead.  (Forbidden because of poorly defined behaviour in JWK.)
usages  (forbidden)  Use "use" instead.  (Forbidden because Javascript conflicts with "use" and "key_ops" and it's confusing.)
```
 



# Canonical Hash (CH)
CH is used to hash over JSON objects and binary blobs. For binary blobs, the CH
of the binary is just the hash over the bytes of the binary. 
 
For objects, standard JSON is used.  For example, strings must be in quotes and
keys are always strings. 
 
For binary blobs (.png, .txt, .pdf, etc...) is just the hash over the bytes of
the binary blob.  This thumbprint in hex may then be used as a value in a JSON
object. 
 
In this way a Coze object may use CH at many levels.  An image blob
(example.png) may be referred to in a JSON object
({"Sub":"Bob","Image":"32365bb3c70eb559dfffc902849442f343bd431d9f322d6a50a43e863686d2bd"})
Then the CH of the Coze object may be used to refer to object.
 
CH can also accept a defined `canon`, a schema with the fields that will be
thumbprinted over. 
 
## How to create a CH Thumbprint
1. Create a json object with key value pairs.  Keys in JSON are always strings.
2. Order elements by key according to their UTF-8 encoding.
3. Remove insignificant spaces (named "compactify" by some).
4. Convert stringified JSON object to bytes (by taking the UTF-8 byte
   representation of the string)
5. Using a hashing algorithm of choice, hash over bytes. 

CH returns the digest as bytes, and those bytes can be represented in any base
paradigm such as hex or base 10.
 
## Canonical Helpers:
The following are recommended Canonical Hash helper functions. 
 
- Canon - Removes extra elements in the object according to the "canon" schema.
  Accepts JSON/Object/Struct/Array as canon and removes "extra" fields.  Returns
  Canonicalized JSON object. (Go doesn't need this, Javascript does.)
- Canonical - Accepts JSON object and optional canon.   Returns compactified
  JSON string/bytes.   
- CH - "CanonicalHash" - Accepts objects and optional canon.  Returns digest
  bytes. 
- CHH - "CanonicalHashHex" Accepts a JSON object and optional canon.  Returns
  digest hex. 
- CHSH - "CanonicalHashSignHex" Accepts JSON, Coze key, hashing algorithm, and
  optional canon.  Returns digest hex. 

 
# Coze Constraints
- JSON keys and values are case sensitive.
- If is reccomended that any short values for strings is no longer than 150
  characters.  

## Constraints on values. 
- `kid`    (Key ID)    Soft limit of 150 characters. 
- `typ`    (Type)      Soft limit of 150 characters.

 

# Future Use Reserved
- Support trailing commas
 
```

pubkey    Encapsulation of a public key object.  Its value must be an object that
                 *does not contain* any private key components.  Example: "pubkey:{"alg": "ES256", ...}
prikey    Encapsulation of a private key object   Its value must be an object.  It may contain public key components.  Example: "prikey:{"alg": "ES256", ...}

pd: "pay digest",   The digest of the pay.
ph "pay hash"       The hashing algorithm used for the pay digest.  
 
mrd: merkle root digest []Bytes
mrh: merkle root hash algo  (specific hash (SHA-256))
mra: merkle root algo  (For different types of merkle tree algos)
 
rout: thumprint of the core indexer
root: //thumprint of the core identity used to sign. (Maybe iss?  "Issuer")
oad: "Owner Address Digest"   (ie 8ed1ff182a964df1c9b0ec846e259b0ff217685de5fa7006581c33c6f5fa7a4b)
oah: "Owner Address Hashing algo",  (ie SHA256)
 
ha:   (optional) Hashing algo Global default.  More specific can be different.  Example: "SHA-256"
hd: "head digest"
hh: "Head Digest Hashing Algorithm"

// Other ideas:
dig: Digest. 
 
la: List algorithm.  (Hash as hex, hash as Base189, hash as Bytes, Merkle Patricia trie)
lh: List hashing alog (Sha-256)
ld: List Hex of bytes[32]
td   (redundant) Alternative to "tmb"

// multisig
*iss: (if multisig, this is the thumb of the account that's issuing it.)
sh: (Signature Header) A sig header object separate from the sig if the header needs to be signed. 
sd: "sig or sigs Digest"
 
```
 




