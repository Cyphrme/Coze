# Coze Release Notes
## Major changes

- Added `alg` `Ed25519`
	- Like ECDSA, in Coze `Ed25519` only ever signs digests, not unhashed messages.
	Ed25519's digest for `cad` is SHA-512.  
	- Waiting on Go for `Ed25519ph` support. 
- Coze now uses base64 and not Hex.
- ECDSA x and y have been consolidated to x. 
	- All future designs will serialize all private components to `d`, and all
	public components to `x`. 
- `pay` takes the place of `head`
- `pay` no longer needs to be utf8 sorted and fields can appear in any order.
	- `cad` is still the hash over the compactified JSON representation.  
	- Specific applications can still require specific order of fields. 
- `cy` fields, except `sig`, are all "meta" fields.  
	- `cy.can` is no longer an input variable for creating a canonical pay, but
	rather is only used to describe the existing canon of pay.  
- Moved 'alg.go' to package `coze`, moved `cryptokey` to its own package, and
  removed `coze`'s dependency on `cryptokey`.  
- All `pay` fields are optional.  An empty `pay` is valid.  
- Removed redundant tests.
- Various bug fixes.
- Go Coze changes "Cy" struct to "Coze", removed the JSON encapsulator.

## Additional Explicit design decisions:
- Coze follows the underlying algorithm's endianness.  ECDSA is big
  endian and EdDSA is little endian.  

## To do:
	Constant time Decode and MustDecode method for private keys. 

# Other thoughts:
- We may end up making a "core" Coze with decently trusted and widely adopted
algorithms, and an "extended"/"x" Coze with less commonly implemented
algorithms.  