package coze

import (
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"strings"

	"golang.org/x/crypto/sha3"
)

type (
	// Alg is a declarative abstraction for cryptographic functions for Coze.
	// For more on Alg, see the main Coze README.
	//
	// The integer value of the "enum" will change in the future. Use the string
	// name for algos when storing information.
	//
	// Hierarchy for signing and hashing cryptographic functions. Naming is
	// inspired by taxonomic rank.
	//
	//  - Level 0 species - "SpcAlg"  (e.g.: ES256) (species)
	//  - Level 1 genus   - "GenAlg"  (e.g.: ECDSA) (genus)
	//  - Level 2 family  - "FamAlg"  (e.g.: EC)    (family)
	//
	// The value for a Coze `alg` is always a specific (species) algorithm, e.g.
	// "ES256", and never any other rank, e.g. "ECDSA". The type `Alg` in this
	// package may be any algorithm of any rank.
	//
	// Cryptographic Signature/Encryption/Hashing hierarchy
	//
	//  - EC
	//  -- ECDSA
	//  --- ES224
	//  --- ES256
	//  --- ES384
	//  --- ES512
	//  -- EdDSA
	//  --- Ed25519
	//  --- Ed25519ph
	//  --- Ed448
	//  - SHA
	//  -- SHA-2
	//  --- SHA-224
	//  --- SHA-256
	//  --- SHA-384
	//  --- SHA-512
	//  -- SHA-3
	//  --- SHA3-224
	//  --- SHA3-256
	//  --- SHA3-384
	//  --- SHA3-512
	//  --- SHAKE128
	//  --- SHAKE256
	//
	// Potential Future Support:
	//  - RSA
	//  -- RSASSA-PKCS1-v1_5
	//  --- RS256
	//  - Lattice-Based signatures
	//  - Other future broad types...
	//  -- ECDH
	//
	// "SE" (singing, encryption) is the super type of signing and encryption and
	// excludes hashing.
	//
	// See the main Coze README for Coze supported and unsupported things.
	Alg int // Alg is for all cryptographic algorithms. All levels included.

	GenAlg  int   // Algorithm genus.    Level 1.
	FamAlg  int   // Algorithm family    Level 2.
	HashAlg Alg   // Hashing Algorithm
	SigAlg  SEAlg // Signing Algorithm
	EncAlg  SEAlg // Encryption Algorithm
	Crv     int   // Curve type.  Used for EC curves.
	KeyUse  int   // Key Use. Right now only "sig".

	// SEAlg is the Signing or Encryption alg. Super type of SigAlg and EncAlg and
	// is itself not a specific algorithm and is not included in Alg.
	SEAlg Alg
)

// Params reports all relevant values for an `alg`. If values are not applicable
// for a particular `alg`, values may be populated with the Go zero value, e.g.
// for the hash alg "SHA-256" Curve's value is 0 and omitted from JSON
// marshaling.
type Params struct {
	Name     string
	Genus    GenAlg
	Family   FamAlg
	XSize    int     `json:"X.Size,omitempty"`
	DSize    int     `json:"D.Size,omitempty"`
	Hash     HashAlg `json:",omitempty"`
	HashSize int     `json:"Hash.Size,omitempty"`
	SigSize  int     `json:"Sig.Size,omitempty"`
	Curve    Crv     `json:",omitempty"`
	KeyUse   KeyUse  `json:"Use,omitempty"`
}

// Params sets and returns a Params struct. See struct definition.
func (a Alg) Params() Params {
	var p Params
	p.Name = a.String()
	p.Genus = a.Genus()
	p.Family = a.Family()
	p.XSize = SEAlg(a).XSize()
	p.DSize = SEAlg(a).DSize()
	p.Hash = a.Hash()
	p.HashSize = a.Hash().Size()
	p.SigSize = a.SigAlg().SigSize()
	p.Curve = a.Curve()
	p.KeyUse = a.KeyUse()

	return p
}

// GenAlg "Genus".
const (
	UnknownGenAlg GenAlg = iota
	Ecdsa
	Eddsa
	SHA2
	SHA3
)

func (g GenAlg) String() string {
	return []string{
		"UnknownGenAlg",
		"ECDSA",
		"EdDSA",
		"SHA2",
		"SHA3",
	}[g]
}

func (g GenAlg) MarshalJSON() ([]byte, error) {
	return []byte(`"` + g.String() + `"`), nil
}

// FamAlg "Family".
const (
	UnknownFamAlg FamAlg = iota
	EC
	SHA
	RSA
)

func (f FamAlg) String() string {
	return []string{
		"UnknownFamAlg",
		"EC",
		"SHA",
		"RSA",
	}[f]
}

func (f FamAlg) MarshalJSON() ([]byte, error) {
	return []byte(`"` + f.String() + `"`), nil
}

const (
	UnknownAlg Alg = iota
)

func (a *Alg) UnmarshalJSON(b []byte) error {
	a.Parse(string(b))

	return nil
}

func (a Alg) MarshalJSON() ([]byte, error) {
	s := `"` + getString(int(a)) + `"`
	return []byte(s), nil
}

func (a *Alg) Parse(s string) {
	s = strings.Trim(s, `"`)

	switch s {
	default:
		*a = UnknownAlg
	case "UnknownAlg":
		*a = UnknownAlg
	case "UnknownSigAlg":
		*a = Alg(UnknownSignAlg)
	case "ES224":
		*a = Alg(ES224)
	case "ES256":
		*a = Alg(ES256)
	case "ES384":
		*a = Alg(ES384)
	case "ES512":
		*a = Alg(ES512)
	case "Ed25519":
		*a = Alg(Ed25519)
	case "Ed25519ph":
		*a = Alg(Ed25519ph)
	case "Ed448":
		*a = Alg(Ed448)
	//	Placeholder for future.
	// case "RS256":
	// 	*a = Alg(RS256)
	// case "RS384":
	// 	*a = Alg(RS384)
	// case "RS512":
	// 	*a = Alg(RS512)
	case "UnknownEncAlg":
		*a = Alg(UnknownEncAlg)
	case "UnknownHashAlg":
		*a = Alg(UnknownHashAlg)
	case "SHA-224":
		*a = Alg(SHA224)
	case "SHA-256":
		*a = Alg(SHA256)
	case "SHA-384":
		*a = Alg(SHA384)
	case "SHA-512":
		*a = Alg(SHA512)
	case "SHA3-224":
		*a = Alg(SHA3224)
	case "SHA3-256":
		*a = Alg(SHA3256)
	case "SHA3-384":
		*a = Alg(SHA3384)
	case "SHA3-512":
		*a = Alg(SHA3512)
	case "SHAKE128":
		*a = Alg(SHAKE128)
	case "SHAKE256":
		*a = Alg(SHAKE256)
	}
}

// getString must follow the same order as Alg's Parse.
func getString(i int) (s string) {
	return []string{
		"UnknownAlg",
		"UnknownSigAlg",
		"ES224",
		"ES256",
		"ES384",
		"ES512",
		"Ed25519",
		"Ed25519ph",
		"Ed448",
		"RS256", // Placeholder for future.
		"RS384",
		"RS512",
		"UnknownEncAlg",
		"UnknownHashAlg",
		"SHA-224",
		"SHA-256",
		"SHA-384",
		"SHA-512",
		"SHA3-224",
		"SHA3-256",
		"SHA3-384",
		"SHA3-512",
		"SHAKE128",
		"SHAKE256",
	}[i]
}

func (a Alg) String() string {
	return getString(int(a))
}

func Parse(s string) (a *Alg) {
	a = new(Alg)
	a.Parse(s)
	return a
}

// Genus is for ECDSA, EdDSA, SHA-2, SHA-3.
func (a Alg) Genus() GenAlg {
	switch a {
	default:
		return UnknownGenAlg
	case Alg(ES224), Alg(ES256), Alg(ES384), Alg(ES512):
		return Ecdsa
	case Alg(Ed25519), Alg(Ed25519ph), Alg(Ed448):
		return Eddsa
	case Alg(SHA224), Alg(SHA256), Alg(SHA384), Alg(SHA512):
		return SHA2
	case Alg(SHA3224), Alg(SHA3256), Alg(SHA3384), Alg(SHA3512), Alg(SHAKE128), Alg(SHAKE256):
		return SHA3
	}
}

// Family is for EC, SHA, and RSA.
func (a Alg) Family() (f FamAlg) {
	switch a {
	default:
		f = UnknownFamAlg
	case Alg(ES224), Alg(ES256), Alg(ES384), Alg(ES512), Alg(Ed25519), Alg(Ed25519ph), Alg(Ed448):
		f = EC
	case Alg(SHA224), Alg(SHA256), Alg(SHA384), Alg(SHA512), Alg(SHA3224), Alg(SHA3256), Alg(SHA3384), Alg(SHA3512), Alg(SHAKE128), Alg(SHAKE256):
		f = SHA
	}
	return
}

// Hash returns respective hashing algorithm if specified. If alg is a hashing
// algorithm, it returns itself.
func (a Alg) Hash() HashAlg {
	// Return itself if Alg is a HashAlg
	if a.Family() == SHA {
		return HashAlg(a)
	}
	// Assume Alg's hashing alg is defined by SEAlg.
	return SEAlg(a).Hash()
}

func (a Alg) SigAlg() SigAlg {
	return SigAlg(a)
}

const (
	SEAlgUnknown SEAlg = iota
)

func (se SEAlg) String() string {
	return getString(int(se))
}

func ParseSEAlg(s string) SEAlg {
	return SEAlg(*Parse(s))
}

func (se SEAlg) SigAlg() SigAlg {
	switch SigAlg(se) {
	default:
		return UnknownSignAlg
	case ES224:
		return ES224
	case ES256:
		return ES256
	case ES384:
		return ES384
	case ES512:
		return ES512
	case Ed25519:
		return Ed25519
	case Ed25519ph:
		return Ed25519ph
	case Ed448:
		return Ed448
	}
}

func (se *SEAlg) UnmarshalJSON(b []byte) error {
	se.Parse(string(b))
	return nil
}

func (se SEAlg) MarshalJSON() ([]byte, error) {
	s := `"` + se.String() + `"`
	return []byte(s), nil
}

func (se *SEAlg) Parse(s string) {
	*se = SEAlg(*Parse(s))
}

func (se SEAlg) Curve() Crv {
	return Alg(se).Curve()
}

func (se SEAlg) Genus() GenAlg {
	return Alg(se).Genus()
}

func (se SEAlg) Family() FamAlg {
	return Alg(se).Family()
}

// Hash returns respective hashing algorithm if specified.
func (se SEAlg) Hash() HashAlg {
	// Only SigAlgs support .Hash() at the moment.
	return se.SigAlg().Hash()
}

// XSize returns the byte size of `x`.  Returns 0 on error.
//
//For ECDSA `x` is the concatenation of X and Y.
func (se SEAlg) XSize() int {
	switch SigAlg(se) {
	default:
		return 0
	case ES224:
		return 56
	case ES256:
		return 64
	case ES384:
		return 96
	case ES512:
		return 132 // X and Y are 66 bytes (Rounded up for P521)
	case Ed25519, Ed25519ph:
		return 32
	case Ed448:
		return 57
	}
}

// DSize returns the byte size of `d`. Returns 0 on error.
func (se SEAlg) DSize() int {
	switch SigAlg(se) {
	default:
		return 0
	case ES224:
		return 28
	case ES256, Ed25519, Ed25519ph:
		return 32
	case ES384:
		return 48
	case Ed448:
		return 57
	case ES512:
		return 66 // Rounded up for P521
	}
}

const (
	UnknownEncAlg EncAlg = iota + 10
)

////////////////
//  HashAlg  //
////////////////

// HashAlg is a hashing algorithm. See also https://golang.org/pkg/crypto/Hash
const (
	// HashAlg is after Alg, SigAlg, and EncAlg.
	UnknownHashAlg HashAlg = iota + 13
	SHA224                 // SHA-2
	SHA256
	SHA384
	SHA512
	SHA3224 // SHA-3
	SHA3256
	SHA3384
	SHA3512
	SHAKE128 // Shake
	SHAKE256
)

func (h HashAlg) String() string {
	return getString(int(h))
}

func (h *HashAlg) UnmarshalJSON(b []byte) error {
	h.Parse(string(b))

	return nil
}

func (h HashAlg) MarshalJSON() ([]byte, error) {
	s := `"` + getString(int(h)) + `"`
	return []byte(s), nil
}

func (h *HashAlg) Parse(s string) {
	*h = HashAlg(*Parse(s))
}

func ParseHashAlg(s string) HashAlg {
	return HashAlg(*Parse(s))
}

// goHash returns a Go hash.Hash from the hashing algo.
//
// SHAKE does not satisfy Go's hash.Hash and uses sha3.SkakeHash.
func (h *HashAlg) goHash() hash.Hash {
	switch *h {
	default:
		return nil
	case SHA224:
		return sha256.New224() // There is no 224 package. 224 is in the 256 package.
	case SHA256:
		return sha256.New()
	case SHA384:
		return sha512.New384() // There is no 384 package. 384 is in the 512 package.
	case SHA512:
		return sha512.New()
	case SHA3224:
		return sha3.New224()
	case SHA3256:
		return sha3.New256()
	case SHA3384:
		return sha3.New384()
	case SHA3512:
		return sha3.New512()
	}
}

// HashSize returns the digest size in bytes for the given hashing algorithm.
//
// SHAKE128 has 128 bits of pre-collision resistance and a capacity of 256,
// although it has arbitrary output size. SHAKE256 has 256 bits of pre-collision
// resistance and a capacity of 512, although it has arbitrary output size.
func (h HashAlg) Size() int {
	switch h {
	default:
		return 0
	case SHA224, SHA3224:
		return 28
	case SHA256, SHA3256, SHAKE128:
		return 32
	case SHA384, SHA3384:
		return 48
	case SHA512, SHA3512, SHAKE256:
		return 64
	}
}

const (
	// Must be in order according to Alg.Parse().
	UnknownSignAlg SigAlg = iota + 1
	ES224
	ES256
	ES384
	ES512

	Ed25519
	Ed25519ph
	Ed448

	// Not implemented [RS256, RS384, RS512].
)

func (s SigAlg) FamAlg() FamAlg {
	switch s {
	default:
		return UnknownFamAlg
	case ES224, ES256, ES384, ES512, Ed25519, Ed25519ph, Ed448:
		return EC
		// // Not implemented:
		// case RS256, RS384, RS512:
		// 	return RSA
	}
}

func (s SigAlg) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

func (s SigAlg) Genus() GenAlg {
	switch s {
	default:
		return UnknownGenAlg
	case ES224, ES256, ES384, ES512:
		return Ecdsa
	case Ed25519, Ed25519ph, Ed448:
		return Eddsa
	}
}

func (s SigAlg) String() string {
	return getString(int(s))
}

// Hash returns respective hashing algorithm if specified.
func (s SigAlg) Hash() HashAlg {
	switch s {
	default:
		return 0
	case ES224:
		return SHA224
	case ES256:
		return SHA256
	case ES384:
		return SHA384
	case ES512, Ed25519, Ed25519ph:
		return SHA512
	case Ed448:
		return SHAKE256
	}
}

// SigSize returns the signature size for the given Crypto Algorithm.
//
// Ed25519's SigSize is from RFC8032_5.1.6.6.
func (s SigAlg) SigSize() int {
	switch s {
	default:
		return 0
	case ES224:
		return 56
	case ES256, Ed25519, Ed25519ph:
		return 64
	case ES384:
		return 96
	case Ed448:
		return 114
	case ES512:
		// Curve P-521 uses 521 bits. This is then padded up the the nearest
		// byte (528) for R and S. 132 = (528*2)/8
		return 132
	}
}

const (
	KeyUseUnknown KeyUse = iota
	SigUse               // "Signing Use"
	EncUse               // "Encryption Use"
)

func (u *KeyUse) UnmarshalJSON(b []byte) error {
	u.Parse(string(b))
	return nil
}

func (u KeyUse) MarshalJSON() ([]byte, error) {
	s := "\"" + u.String() + "\""
	return []byte(s), nil
}

// KeyUse returns the KeyUse.
func (a Alg) KeyUse() KeyUse {
	switch a.Genus() {
	default:
		return KeyUseUnknown
	case Eddsa, Ecdsa:
		return SigUse
	}
}

func (u *KeyUse) Parse(s string) {
	s = strings.Trim(s, "\"")
	switch s {
	default:
		*u = KeyUseUnknown
	case "sig":
		*u = SigUse
	case "enc":
		*u = EncUse
	}
}

func ParseKeyUse(s string) KeyUse {
	u := new(KeyUse)
	u.Parse(s)
	return *u
}

func (u KeyUse) String() string {
	return []string{
		"UnknownKeyUse",
		"sig",
		"enc",
	}[u]
}

const (
	UnknownCrv Crv = iota
	P224
	P256
	P384
	P521
	Curve25519
	Curve448
)

// Curve returns the curve for the given alg, if it has one.
func (a Alg) Curve() Crv {
	switch SigAlg(a) {
	default:
		return UnknownCrv
	case ES224:
		return P224
	case ES256:
		return P256
	case ES384:
		return P384
	case ES512:
		return P521 // The curve != the alg
	case Ed25519, Ed25519ph:
		return Curve25519
	case Ed448:
		return Curve448
	}
}

func (c Crv) String() string {
	return []string{
		"UnknownCrv",
		"P-224",
		"P-256",
		"P-384",
		"P-521",
		"Curve25519",
		"Curve448",
	}[c]
}

func (c Crv) MarshalJSON() ([]byte, error) {
	return []byte(`"` + c.String() + `"`), nil
}

func (c *Crv) Parse(s string) {
	switch s {
	default:
		*c = UnknownCrv
	case "P-224":
		*c = P224
	case "P-256":
		*c = P256
	case "P-384":
		*c = P384
	case "P-521":
		*c = P521
	case "Curve25519":
		*c = Curve25519
	case "Curve448":
		*c = Curve448
	}
}

// Curve returns Go's elliptic.Curve for the given crv.
func (c Crv) EllipticCurve() elliptic.Curve {
	switch c {
	default:
		return nil
	case P224:
		return elliptic.P224()
	case P256:
		return elliptic.P256()
	case P384:
		return elliptic.P384()
	case P521:
		return elliptic.P521()
	}
}
