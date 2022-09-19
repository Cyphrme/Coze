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
	//    - ECDSA
	//      - ES224
	//      - ES256
	//      - ES384
	//      - ES512
	//    - EdDSA
	//      - Ed25519
	//      - Ed25519ph
	//      - Ed448
	//  - SHA
	//    - SHA-2
	//      - SHA-224
	//      - SHA-256
	//      - SHA-384
	//      - SHA-512
	//    - SHA-3
	//      - SHA3-224
	//      - SHA3-256
	//      - SHA3-384
	//      - SHA3-512
	//      - SHAKE128
	//      - SHAKE256
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
	Use     int   // The only valid values are 'sig', 'enc', and 'dig'.

	// SEAlg is the Signing or Encryption alg. Super type of SigAlg and EncAlg and
	// is itself not a specific algorithm and is not included in `Alg`.
	SEAlg Alg
)

// Params reports all relevant values for an `alg`. If values are not applicable
// for a particular `alg`, values may be populated with the Go zero value, e.g.
// for the hash alg "SHA-256" Curve's value is 0 and omitted from JSON
// marshaling.
type Params struct {
	Name     string
	Genus    GenAlg  `json:"Genus"`
	Family   FamAlg  `json:"Family"`
	XSize    int     `json:"X.Size,omitempty"`
	DSize    int     `json:"D.Size,omitempty"`
	Hash     HashAlg `json:"Hash,omitempty"`
	HashSize int     `json:"Hash.Size,omitempty"`
	SigSize  int     `json:"Sig.Size,omitempty"`
	Curve    Crv     `json:"Curve,omitempty"`
	Use      Use     `json:"Use,omitempty"`
}

// Params sets and returns a Params struct. See struct definition.
func (a Alg) Params() Params {
	return Params{
		Name:     a.String(),
		Genus:    a.Genus(),
		Family:   a.Family(),
		XSize:    SEAlg(a).XSize(),
		DSize:    SEAlg(a).DSize(),
		Hash:     a.Hash(),
		HashSize: a.Hash().Size(),
		SigSize:  a.SigAlg().SigSize(),
		Curve:    a.Curve(),
		Use:      a.Use(),
	}
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

//////////////////////
// SEAlg
/////////////////////

const (
	SEAlgUnknown SEAlg = iota
)

func (se SEAlg) String() string {
	return Alg(se).String()
}

func ParseSEAlg(s string) SEAlg {
	return SEAlg(*Parse(s))
}

func (se SEAlg) SigAlg() SigAlg {
	switch SigAlg(se) {
	default:
		return UnknownSigAlg
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
	return []byte(`"` + se.String() + `"`), nil
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
// For ECDSA `x` is the concatenation of X and Y.
func (se SEAlg) XSize() int {
	switch SigAlg(se) {
	default:
		return 0
	case Ed25519, Ed25519ph:
		return 32
	case ES224:
		return 56
	case Ed448:
		return 57
	case ES256:
		return 64
	case ES384:
		return 96
	case ES512:
		return 132 // X and Y are 66 bytes (Rounded up for P521)
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

////////////////////
// Alg
///////////////////

// algs is the "default" algs, which at the moment is only "UnknownAlg".  It
// remains unmodified, unlike variable `Algs`.
var algs = []string{
	"UnknownAlg",
}

// Algs is initialized in the init function to include all algs, including
// "UnknownAlg", SigAlg, EncAlg, and DigAlg.
var Algs = algs
var SigAlgs = []string{
	"UnknownSigAlg",
	"ES224",
	"ES256",
	"ES384",
	"ES512",
	"Ed25519",
	"Ed25519ph",
	"Ed448",
}

// Encryption algs.
var EncAlgs = []string{
	"UnknownEncAlg",
	//// Placeholder for future.
	// "RS256",
	// "RS384",
	// "RS512",
}

// Digest/Hash algs.
var DigAlgs = []string{
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
}

func init() {
	// Algs must be in order according to Alg.Parse().
	Algs = append(Algs, SigAlgs...)
	Algs = append(Algs, EncAlgs...)
	Algs = append(Algs, DigAlgs...)
}

const (
	UnknownAlg Alg = iota
)

func (a *Alg) UnmarshalJSON(b []byte) error {
	a.Parse(string(b))
	return nil
}

func (a Alg) MarshalJSON() ([]byte, error) {
	return []byte(`"` + a.String() + `"`), nil
}

func (a Alg) String() string {
	return Algs[int(a)]
}

func Parse(s string) (a *Alg) {
	a = new(Alg)
	a.Parse(s)
	return a
}

func (a *Alg) Parse(s string) {
	switch strings.Trim(s, `"`) {
	default:
		*a = UnknownAlg
	case "UnknownAlg":
		*a = UnknownAlg
	case "UnknownSigAlg":
		*a = Alg(UnknownSigAlg)
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
func (a Alg) Family() FamAlg {
	switch a {
	default:
		return UnknownFamAlg
	case Alg(ES224), Alg(ES256), Alg(ES384), Alg(ES512), Alg(Ed25519), Alg(Ed25519ph), Alg(Ed448):
		return EC
	case Alg(SHA224), Alg(SHA256), Alg(SHA384), Alg(SHA512), Alg(SHA3224), Alg(SHA3256), Alg(SHA3384), Alg(SHA3512), Alg(SHAKE128), Alg(SHAKE256):
		return SHA
	}
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

//////////////////////
// SigAlg
/////////////////////

const (
	// SigAlg appears in `Algs` after algs.
	UnknownSigAlg SigAlg = iota + SigAlg(UnknownAlg) + 1
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
	return Alg(s).String()
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

//////////////////////
// EncAlg
/////////////////////

const (
	// EncAlg appears in `Algs` after SigAlgs.
	UnknownEncAlg EncAlg = iota + EncAlg(Ed448) + 1
)

////////////////
//  HashAlg  //
////////////////

// HashAlg is a hashing algorithm. See also https://golang.org/pkg/crypto/Hash
const (
	// HashAlg appears in `Algs` after EncAlgs.
	UnknownHashAlg HashAlg = iota + HashAlg(UnknownEncAlg) + 1
	// SHA-2
	SHA224
	SHA256
	SHA384
	SHA512
	// SHA-3
	SHA3224
	SHA3256
	SHA3384
	SHA3512
	// SHAKE
	SHAKE128
	SHAKE256
)

func (h HashAlg) String() string {
	return Alg(h).String()
}

func (h *HashAlg) UnmarshalJSON(b []byte) error {
	h.Parse(string(b))

	return nil
}

func (h HashAlg) MarshalJSON() ([]byte, error) {
	return []byte(`"` + Alg(h).String() + `"`), nil
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
	UseUnknown Use = iota
	SigUse         // "Signing Use"
	EncUse         // "Encryption Use"
	DigUse         // "Digest Use"
)

func (u *Use) UnmarshalJSON(b []byte) error {
	u.Parse(string(b))
	return nil
}

func (u Use) MarshalJSON() ([]byte, error) {
	return []byte("\"" + u.String() + "\""), nil
}

// Use returns the Use.
func (a Alg) Use() Use {
	switch a.Genus() {
	default:
		return UseUnknown
	case Eddsa, Ecdsa:
		return SigUse
	case SHA2, SHA3:
		return DigUse
	}
}

func (u *Use) Parse(s string) {
	switch strings.Trim(s, "\"") {
	default:
		*u = UseUnknown
	case "sig":
		*u = SigUse
	case "enc":
		*u = EncUse
	case "dig":
		*u = DigUse
	}
}

func ParseUse(s string) Use {
	u := new(Use)
	u.Parse(s)
	return *u
}

func (u Use) String() string {
	return []string{
		"UnknownUse",
		"sig",
		"enc",
		"dig",
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
