package coze

import (
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"strings"

	"golang.org/x/crypto/sha3"
)

// Alg is a declarative abstraction for cryptographic functions
// for Coze in Go.
//
// See the main Coze README.
//
// Hierarchy for signing and hashing cryptographic functions.  Naming is inspired by taxonomic rank.
//
//  - Level 0 species - "SpcAlg"  (e.g.: ES256) (species)
//  - Level 1 genus   - "GenAlg"  (e.g.: ECDSA) (genus)
//  - Level 2 family  - "FamAlg"  (e.g.: EC)    (family)
//
// The value for a Coze `alg` is always specific (species) algorithm, e.g.
// "ES256", and never any other rank, e.g. "ECDSA".  The type `Alg` in this
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
// The integer value of the "enum" will change in the future. Use the string
// name for algos when storing information.
//
// See the main Coze README for Coze supported and unsupported things.
type Alg int    // Alg is for all cryptographic algorithms.  All levels included.
type GenAlg int // Algorithm genus.    Level 1.
type FamAlg int // Algorithm family    Level 2.

type HashAlg Alg  // Hashing Algorithm
type SigAlg SEAlg // Signing Algorithm
type EncAlg SEAlg // Encryption Algorithm

type Crv int    // Curve type.  Used for EC curves.
type KeyUse int // Key Use. 2021/05/19 Right now only "sig (Perhaps "enc" in the future)

// SEAlg is the Signing or Encryption alg.  Super type of SigAlg and EncAlg and
// is itself not a specific algorithm and is not included in Alg.
type SEAlg Alg

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

// GenAlg "Genus"
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

// FamAlg "Family"
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

//////////////////////
//       Alg        //
//////////////////////
const (
	UnknownAlg Alg = iota
)

func (h *Alg) UnmarshalJSON(b []byte) error {
	h.Parse(string(b))

	return nil
}

func (h Alg) MarshalJSON() ([]byte, error) {
	s := `"` + getString(int(h)) + `"`
	return []byte(s), nil
}

func (a *Alg) Parse(s string) {
	s = strings.Trim(s, `"`)

	switch s {
	default:
		*a = UnknownAlg
	case "UnknownAlg":
		*a = Alg(UnknownAlg)
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
		*a = Alg(Sha224)
	case "SHA-256":
		*a = Alg(Sha256)
	case "SHA-384":
		*a = Alg(Sha384)
	case "SHA-512":
		*a = Alg(Sha512)
	case "SHA3-224":
		*a = Alg(Sha3224)
	case "SHA3-256":
		*a = Alg(Sha3256)
	case "SHA3-384":
		*a = Alg(Sha3384)
	case "SHA3-512":
		*a = Alg(Sha3512)
	case "SHAKE128":
		*a = Alg(Shake128)
	case "SHAKE256":
		*a = Alg(Shake256)
	}
	return
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

func (s Alg) String() string {
	return getString(int(s))
}

func Parse(s string) (a *Alg) {
	a = new(Alg)
	a.Parse(s)
	return a
}

// Genus is for ECDSA, EdDSA, SHA-2, SHA-3
func (a Alg) Genus() GenAlg {
	switch a {
	default:
		return UnknownGenAlg
	case Alg(ES224), Alg(ES256), Alg(ES384), Alg(ES512):
		return Ecdsa
	case Alg(Ed25519), Alg(Ed25519ph), Alg(Ed448):
		return Eddsa
	case Alg(Sha224), Alg(Sha256), Alg(Sha384), Alg(Sha512):
		return SHA2
	case Alg(Sha3224), Alg(Sha3256), Alg(Sha3384), Alg(Sha3512), Alg(Shake128), Alg(Shake256):
		return SHA3
	}
}

// Family is for 	EC,	SHA, and	RSA
func (a Alg) Family() (f FamAlg) {
	switch a {
	default:
		f = UnknownFamAlg
	case Alg(ES224), Alg(ES256), Alg(ES384), Alg(ES512), Alg(Ed25519), Alg(Ed25519ph), Alg(Ed448):
		f = EC
	case Alg(Sha224), Alg(Sha256), Alg(Sha384), Alg(Sha512), Alg(Sha3224), Alg(Sha3256), Alg(Sha3384), Alg(Sha3512), Alg(Shake128), Alg(Shake256):
		f = SHA
	}
	return
}

// Hash returns respective hashing algorithm if specified.  If alg is a hashing
// algorithm, it returns itself.
func (a Alg) Hash() HashAlg {
	// Return itself if Alg is a HashAlg
	if a.Family() == SHA {
		return HashAlg(a)
	}
	// Assume Alg's hashing alg is defined by SEAlg.
	return SEAlg(a).Hash()
}

/////////////////////////////////////////
//  SEAlg (Signing or Encryption Alg)  //
/////////////////////////////////////////
const (
	SEAlgUnknown SEAlg = iota
)

func (s SEAlg) String() string {
	return getString(int(s))
}

func ParseSEAlg(s string) SEAlg {
	return (SEAlg)(*Parse(s))
}

func (s SEAlg) SigAlg() (sa SigAlg) {
	switch s {
	case SEAlg(ES224):
		sa = ES224
	case SEAlg(ES256):
		sa = ES256
	case SEAlg(ES384):
		sa = ES384
	case SEAlg(ES512):
		sa = ES512
	case SEAlg(Ed25519):
		sa = Ed25519
	case SEAlg(Ed25519ph):
		sa = Ed25519ph
	case SEAlg(Ed448):
		sa = Ed448
	}
	return sa
}

func (h *SEAlg) UnmarshalJSON(b []byte) error {
	h.Parse(string(b))
	return nil
}

func (h SEAlg) MarshalJSON() ([]byte, error) {
	s := `"` + h.String() + `"`
	return []byte(s), nil
}

func (se *SEAlg) Parse(s string) {
	*se = (SEAlg)(*Parse(s))
}

func (h SEAlg) Curve() Crv {
	return Alg(h).Curve()
}

func (h SEAlg) Genus() GenAlg {
	return Alg(h).Genus()
}

func (h SEAlg) Family() FamAlg {
	return Alg(h).Family()
}

// Hash returns respective hashing algorithm if specified.
func (a SEAlg) Hash() HashAlg {
	// Only SigAlgs support .Hash() at the moment.
	return a.SigAlg().Hash()
}

// XSize returns the byte size of `x`.  Returns 0 on error.
func (a SEAlg) XSize() int {
	// For ECDSA `x` is the concatenation of X and Y.
	switch SigAlg(a) {
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
	}
	return 0
}

// DSize returns the byte size of `d`.  Returns 0 on error.
func (a SEAlg) DSize() int {
	switch SigAlg(a) {
	case ES224:
		return 28
	case ES256, Ed25519, Ed25519ph:
		return 32
	case ES384:
		return 48
	case ES512:
		return 66 // Rounded up for P521
	}
	return 0
}

///////////////
//  Enc Alg  //
///////////////
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
	Sha224                 // SHA-2
	Sha256
	Sha384
	Sha512
	Sha3224 // SHA-3
	Sha3256
	Sha3384
	Sha3512
	Shake128 // Shake
	Shake256
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
	*h = (HashAlg)(*Parse(s))
}

func ParseHashAlg(s string) HashAlg {
	return (HashAlg)(*Parse(s))
}

// goHash returns a Go hash.Hash from the hashing algo.
//
// SHAKE does not satisfy Go's hash.Hash and uses sha3.SkakeHash.
func (ha *HashAlg) goHash() (h hash.Hash) {
	switch *ha {
	case Sha224:
		h = sha256.New224() // There is no 224 package.  224 is in the 256 package.
	case Sha256:
		h = sha256.New()
	case Sha384:
		h = sha512.New384() // There is no 384 package.  384 is in the 512 package.
	case Sha512:
		h = sha512.New()
	case Sha3224:
		h = sha3.New224()
	case Sha3256:
		h = sha3.New256()
	case Sha3384:
		h = sha3.New384()
	case Sha3512:
		h = sha3.New512()
	default:
		return nil
	}

	return h
}

// HashSize returns the digest size in bytes for the given hashing algorithm.
//
// SHAKE128 has 128 bits of pre-collision resistance and a capacity of 256,
// although it has arbitrary output size. SHAKE256 has 256 bits of pre-collision
// resistance and a capacity of 512, although it has arbitrary output size.
func (h HashAlg) Size() int {
	switch h {
	case Sha224, Sha3224:
		return 28
	case Sha256, Sha3256, Shake128:
		return 32
	case Sha384, Sha3384:
		return 48
	case Sha512, Sha3512, Shake256:
		return 64
	}
	return 0
}

//////////////
//  SigAlg  //
//////////////
const (
	// Must be in order according to Alg.Parse()
	UnknownSignAlg SigAlg = iota + 1
	ES224
	ES256
	ES384
	ES512

	Ed25519
	Ed25519ph
	Ed448

	// // Not implemented:
	// RS256
	// RS384
	// RS512
)

func (a Alg) SigAlg() SigAlg {
	return SigAlg(a)
}

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
	var h HashAlg
	switch s {
	case ES224:
		h = Sha224
	case ES256:
		h = Sha256
	case ES384:
		h = Sha384
	case ES512, Ed25519, Ed25519ph:
		h = Sha512
	}
	return h
}

// SigSize returns the signature size for the given Crypto Algorithm.
//
// Ed25519's SigSize is from RFC8032_5.1.6.6.
func (s SigAlg) SigSize() int {
	switch s {
	case ES224:
		return 56
	case ES256, Ed25519, Ed25519ph:
		return 64
	case ES384:
		return 96
	case Ed448:
		return 114
	case ES512:
		// Curve P-521 uses 521 bits.  This is then padded up the the nearest
		// byte (528) for R and S. 132 = (528*2)/8
		return 132
	}
	return 0
}

///////////////
//  Key Use //
//////////////
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
	case "sig":
		*u = SigUse
	case "enc":
		*u = EncUse
	default:
		*u = KeyUseUnknown
	}
	return
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

///////////////////
//  Crv (Curve)  //
///////////////////

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
func (a Alg) Curve() (c Crv) {
	switch SigAlg(a) {
	default:
		c = UnknownCrv
	case ES224:
		c = P224
	case ES256:
		c = P256
	case ES384:
		c = P384
	case ES512:
		c = P521 // The curve != the alg
	case Ed25519, Ed25519ph:
		c = Curve25519
	case Ed448:
		c = Curve448
	}
	return
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
	default:
		*c = UnknownCrv
	}

	return
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
