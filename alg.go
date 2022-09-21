package coze

import (
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"math"
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
	Alg    string // Alg is for all "SpcAlg" cryptographic algorithms.
	GenAlg string // Algorithm genus.
	FamAlg string // Algorithm family

	SigAlg SEAlg // Signing Algorithm
	EncAlg SEAlg // Encryption Algorithm
	HshAlg Alg   // Hashing Algorithm

	Use string // The only valid values are 'sig', 'enc', and 'hsh'.
	Crv string // Curve type.  Used for EC curves.

	// SEAlg is the Signing or Encryption alg. Super type of SigAlg and EncAlg and
	// is itself not a specific algorithm and is not included in `Alg`.  It is
	// useful for algorithms that need `x` and/or `d` and related functions.
	SEAlg Alg
)

////////////////
//  Alg
////////////////

const (
	UnknownAlg     Alg    = "UnknownAlg"
	UnknownSigAlg  SigAlg = "UnknownSigAlg"
	ES224          SigAlg = "ES224"
	ES256          SigAlg = "ES256"
	ES384          SigAlg = "ES384"
	ES512          SigAlg = "ES512"
	Ed25519        SigAlg = "Ed25519"
	Ed25519ph      SigAlg = "Ed25519ph"
	Ed448          SigAlg = "Ed448"
	UnknownEncAlg  EncAlg = "UnknownEncAlg"
	UnknownHashAlg HshAlg = "UnknownHashAlg"
	SHA224         HshAlg = "SHA-224"
	SHA256         HshAlg = "SHA-256"
	SHA384         HshAlg = "SHA-384"
	SHA512         HshAlg = "SHA-512"
	SHA3224        HshAlg = "SHA3-224"
	SHA3256        HshAlg = "SHA3-256"
	SHA3384        HshAlg = "SHA3-384"
	SHA3512        HshAlg = "SHA3-512"
	SHAKE128       HshAlg = "SHAKE128"
	SHAKE256       HshAlg = "SHAKE256"
)

// Algs is initialized in the init function to include all algs, including
// "UnknownAlg", SigAlg, EncAlg, and HshAlg.
var Algs = []Alg{
	UnknownAlg,
	Alg(UnknownSigAlg),
	Alg(ES224),
	Alg(ES256),
	Alg(ES384),
	Alg(ES512),
	Alg(Ed25519),
	Alg(Ed25519ph),
	Alg(Ed448),
	Alg(UnknownEncAlg),
	Alg(UnknownHashAlg),
	Alg(SHA224),
	Alg(SHA256),
	Alg(SHA384),
	Alg(SHA512),
	Alg(SHA3224),
	Alg(SHA3256),
	Alg(SHA3384),
	Alg(SHA3512),
	Alg(SHAKE128),
	Alg(SHAKE256),
}

var SigAlgs = []SigAlg{
	UnknownSigAlg,
	ES224,
	ES256,
	ES384,
	ES512,
	Ed25519,
	Ed25519ph,
	Ed448,
}

// Encryption algs.
var EncAlgs = []EncAlg{
	UnknownEncAlg,
	//// Placeholder for future.
	// "RS256",
	// "RS384",
	// "RS512",
}

// Hash algs.
var HshAlgs = []HshAlg{
	UnknownHashAlg,
	SHA224,
	SHA256,
	SHA384,
	SHA512,
	SHA3224,
	SHA3256,
	SHA3384,
	SHA3512,
	SHAKE128,
	SHAKE256,
}

func Parse(s string) Alg {
	a := new(Alg)
	a.Parse(s)
	return *a
}

func (a *Alg) Parse(s string) {
	switch strings.Trim(s, `"`) {
	default:
		*a = UnknownAlg
	case string(UnknownAlg):
		*a = UnknownAlg
	case string(UnknownSigAlg):
		*a = Alg(UnknownSigAlg)
	case string(ES224):
		*a = Alg(ES224)
	case string(ES256):
		*a = Alg(ES256)
	case string(ES384):
		*a = Alg(ES384)
	case string(ES512):
		*a = Alg(ES512)
	case string(Ed25519):
		*a = Alg(Ed25519)
	case string(Ed25519ph):
		*a = Alg(Ed25519ph)
	case string(Ed448):
		*a = Alg(Ed448)
	case string(UnknownEncAlg):
		*a = Alg(UnknownEncAlg)
	case string(UnknownHashAlg):
		*a = Alg(UnknownHashAlg)
	case string(SHA224):
		*a = Alg(SHA224)
	case string(SHA256):
		*a = Alg(SHA256)
	case string(SHA384):
		*a = Alg(SHA384)
	case string(SHA512):
		*a = Alg(SHA512)
	case string(SHA3224):
		*a = Alg(SHA3224)
	case string(SHA3256):
		*a = Alg(SHA3256)
	case string(SHA3384):
		*a = Alg(SHA3384)
	case string(SHA3512):
		*a = Alg(SHA3512)
	case string(SHAKE128):
		*a = Alg(SHAKE128)
	case string(SHAKE256):
		*a = Alg(SHAKE256)
	}
}

// Genus is for ECDSA, EdDSA, SHA-2, SHA-3.
func (a Alg) Genus() GenAlg {
	switch a {
	default:
		return UnknownGenAlg
	case Alg(ES224), Alg(ES256), Alg(ES384), Alg(ES512):
		return ECDSA
	case Alg(Ed25519), Alg(Ed25519ph), Alg(Ed448):
		return EdDSA
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
func (a Alg) Hash() HshAlg {
	// Return itself if type HshAlg // TODO actually do a type def here.
	if a.Family() == SHA {
		return HshAlg(a)
	}
	// Assume Alg's hashing alg is defined by SEAlg.
	return SEAlg(a).Hash()
}

func (a Alg) SigAlg() SigAlg {
	return SigAlg(a)
}

// Params reports all relevant parameters for an `alg`. If a parameter is not
// applicable for a particular `alg`, its value is be populated with the Go
// zero value, e.g. for the hash alg "SHA-256" Curve's value is 0 and omitted
// from JSON marshaling.
type Params struct {
	Name        string
	Genus       GenAlg `json:"Genus"`
	Family      FamAlg `json:"Family"`
	Use         Use    `json:"Use,omitempty"`
	Hash        HshAlg `json:"Hash,omitempty"` // Hash
	HashSize    int    `json:"HashSize,omitempty"`
	HashSizeB64 int    `json:"HashSizeB64,omitempty"`
	XSize       int    `json:"XSize,omitempty"` // Key
	XSizeB64    int    `json:"XSizeB64,omitempty"`
	DSize       int    `json:"DSize,omitempty"`
	DSizeB64    int    `json:"DSizeB64,omitempty"`
	Curve       Crv    `json:"Curve,omitempty"`
	SigSize     int    `json:"SigSize,omitempty"` // Sig
	SigSizeB64  int    `json:"SigSizeB64,omitempty"`
}

// Params sets and returns a Params struct. See struct definition.
func (a Alg) Params() Params {
	p := Params{
		Name:     string(a),
		Genus:    a.Genus(),
		Family:   a.Family(),
		Use:      a.Use(),
		Hash:     a.Hash(),
		HashSize: a.Hash().Size(),
		XSize:    SEAlg(a).XSize(),
		DSize:    SEAlg(a).DSize(),
		Curve:    a.Curve(),
		SigSize:  a.SigAlg().SigSize(),
	}

	toB64 := func(sizeInBytes int) int {
		return int(math.Ceil(float64(4*sizeInBytes) / 3))
	}
	p.HashSizeB64 = toB64(p.HashSize)
	p.XSizeB64 = toB64(p.XSize)
	p.DSizeB64 = toB64(p.DSize)
	p.SigSizeB64 = toB64(p.SigSize)
	return p
}

// GenAlg "Genus".
const (
	UnknownGenAlg GenAlg = "UnknownGenAlg"
	ECDSA         GenAlg = "ECDSA"
	EdDSA         GenAlg = "EdDSA"
	SHA2          GenAlg = "SHA2"
	SHA3          GenAlg = "SHA3"
)

// FamAlg "Family".
const (
	UnknownFamAlg FamAlg = "UnknownFamAlg"
	EC            FamAlg = "EC"
	SHA           FamAlg = "SHA"
	RSA           FamAlg = "RSA"
)

////////////////
//  SEAlg
////////////////

const (
	SEAlgUnknown SEAlg = "UnknownSEAlg"
)

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

func (se *SEAlg) Parse(s string) {
	*se = SEAlg(Parse(s))
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
func (se SEAlg) Hash() HshAlg {
	// Only SigAlgs support .Hash() at the moment.
	return se.SigAlg().Hash()
}

// XSize returns the byte size of `x`.  Returns 0 on invalid algorithm.
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

// DSize returns the byte size of `d`. Returns 0 on invalid algorithm.
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

////////////////
//  SigAlg
////////////////

func (s SigAlg) FamAlg() FamAlg {
	switch s {
	default:
		return UnknownFamAlg
	case ES224, ES256, ES384, ES512, Ed25519, Ed25519ph, Ed448:
		return EC
	}
}

func (s SigAlg) Genus() GenAlg {
	switch s {
	default:
		return UnknownGenAlg
	case ES224, ES256, ES384, ES512:
		return ECDSA
	case Ed25519, Ed25519ph, Ed448:
		return EdDSA
	}
}

// Hash returns respective hashing algorithm if specified.
func (s SigAlg) Hash() HshAlg {
	switch s {
	default:
		return UnknownHashAlg
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

////////////////
//  EncAlg
////////////////

////////////////
//  HashAlg
////////////////

func (h *HshAlg) Parse(s string) {
	*h = HshAlg(Parse(s))
}

func ParseHashAlg(s string) HshAlg {
	return HshAlg(Parse(s))
}

// goHash returns a Go hash.Hash from the hashing algo.
//
// SHAKE does not satisfy Go's hash.Hash and uses sha3.SkakeHash.
func (h *HshAlg) goHash() hash.Hash {
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
// For SHAKE128 and SHAKE256, this function returns the static sizes, 32 and 64 respectively, although the algorithm permits any larger arbitrary output size.
// SHAKE128 has 128 bits of pre-collision resistance and a capacity of 256,
// although it has arbitrary output size. SHAKE256 has 256 bits of pre-collision
// resistance and a capacity of 512, although it has arbitrary output size.
func (h HshAlg) Size() int {
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

////////////////
//  Use
////////////////

const (
	UseUnknown Use = "UnknownUse"
	SigUse     Use = "sig" // "Signing Use"
	EncUse     Use = "enc" // "Encryption Use"
	HshUse     Use = "hsh" // "Hash Use"
)

// Use returns the Use.
func (a Alg) Use() Use {
	switch a.Family() {
	default:
		return UseUnknown
	case EC:
		return SigUse
	case SHA:
		return HshUse
	}
}

func (u *Use) Parse(s string) {
	switch Use(strings.Trim(s, "\"")) {
	default:
		*u = UseUnknown
	case SigUse:
		*u = SigUse
	case EncUse:
		*u = EncUse
	case HshUse:
		*u = HshUse
	}
}

func ParseUse(s string) Use {
	var u Use
	u.Parse(s)
	return u
}

////////////////
//  Curve
////////////////

const (
	UnknownCrv Crv = "UnknownCrv"
	P224       Crv = "P-224"
	P256       Crv = "P-256"
	P384       Crv = "P-384"
	P521       Crv = "P-521"
	Curve25519 Crv = "Curve25519"
	Curve448   Crv = "Curve448"
)

// Curve returns the curve for the given alg.  Returns empty if alg does not
// have a curve.
func (a Alg) Curve() Crv {
	switch SigAlg(a) {
	default:
		return ""
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

func (c *Crv) Parse(s string) {
	switch Crv(s) {
	default:
		*c = UnknownCrv
	case P224:
		*c = P224
	case P256:
		*c = P256
	case P384:
		*c = P384
	case P521:
		*c = P521
	case Curve25519:
		*c = Curve25519
	case Curve448:
		*c = Curve448
	}
}

// Curve returns Go's elliptic.Curve for the given crv.  Returns nil if there is
// no matching  `elliptic.Curve`.
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
