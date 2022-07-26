# Coze Release Notes
## Major changes

- Added `alg` `Ed25519`
	- Like ECDSA, in Coze `Ed25519` only ever signs digests, not unhashed messages.
	Ed25519's digest for `cad` is SHA-512.  
	- Waiting on Go for `Ed25519ph` support. 
- Coze now uses base64 and not Hex.
- Normal
	- Normal is allowed to be variadic.  Normal is the generalization of Canon.  
	- The following are valid Normals:  [Canon, Only, Option, Need, Extra]
	- For Go, using Generics on Normal
- ECDSA x and y have been consolidated to x. 
	- All future designs will serialize all private components to `d`, and all
	public components to `x`. 
- `pay` takes the place of `head`
- `pay` no longer needs to be utf8 sorted and fields can appear in any order.
	- `cad` is still the hash over the compactified JSON representation.  
	- Specific applications can still require specific order of fields. 
- `coze` fields, except `sig`, are all "meta" fields.  
	- `coze.can` is no longer an input variable for creating a canonical pay, but
	rather is only used to describe the existing canon of pay.  
- Moved 'alg.go' to package `coze`, moved `cryptokey` to its own package, and
  removed `coze`'s dependency on `cryptokey`.  
- All `pay` fields are optional.  An empty `pay` is valid.  
- Removed redundant tests.
- Various bug fixes.
- Go Coze changes "Cy" struct to "Coze", removed the JSON encapsulator. Added an
  example to show Coze embedded in another struct.  
- Explicitly disallow duplicate JSON field names in Coze.  
	- Many JSON implementations, including Douglas Crockford's Java
	implementation, error on duplicate keys, some JSON implementations use
	last-value-wins, while others support duplicate keys.  Coze removes this
	ambiguity and requires that implementations fail on duplicate keys.  
	- Duplicate fields is a security issue.  If multiple fields were allowed, for
	example for alg, tmb, or rvk, this could be a source of bugs in
	implementations and surprising behavior to users.

## Additional Explicit design decisions:
- Coze follows the underlying algorithm's endianness.  ECDSA is big
  endian and EdDSA is little endian.  

## To do:
	Constant time Decode and MustDecode method for private keys. 

# Other thoughts:
- We may end up making a "core" Coze with decently trusted and widely adopted
algorithms, and an "extended"/"x" Coze with less commonly implemented
algorithms.  