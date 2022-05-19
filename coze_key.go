package coze

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/cyphrme/coze/enum"
	ce "github.com/cyphrme/coze/enum"
)

// CozeKeyArrayCanon is the canonical form of a Coze Key in array form.
var CozeKeyArrayCanon = []string{"alg", "x"}

// CozeKeyCanon is the canonical form of a Coze Key in struct form.
type CozeKeyCanon struct {
	Alg string `json:"alg"`
	X   B64    `json:"x"`
}

// CozeKey is a Coze key.
//
// See `README.md` for details.
//
// Required Fields (Plus any `alg` specific fields.)
//	- `alg` - Specific algorithm of the key. E.g. "ES256" or "Ed25519".
//	- `iat` - Unix time of when the key was created. E.g. 1626069600.
//	- `tmb` - Key's thumbprint. E.g. "0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD".
//
// Recommended Fields:
//	- `kid` - Human readable label and must not be used programatically. E.g. "My Cyphr.me key".
//
// Key Fields
//	- `d` - Private component.
//	- `x` - Public component. Without JSON omitempty since x is always set, public, private, and all algs.
//
// Revoked
//  - `rvk` - Unix time of key revocation. See docs on `rvk`. E.g. 1626069601.
//
// Optional standard fields
//	-`typ` - The key's type.  "coze/key".
//
type CozeKey struct {
	Alg ce.SEAlg `json:"alg"`
	Kid string   `json:"kid,omitempty"`
	Iat int64    `json:"iat"`
	Tmb B64      `json:"tmb"`

	// ECDSA/EdDSA parameters
	D B64 `json:"d,omitempty"`
	X B64 `json:"x"` // No omitempty since X is always set.

	// Revoked
	Rvk int64 `json:"rvk,omitempty"`

	// Optional parameters
	Typ string `json:"typ,omitempty"`
}

// String returns the stringified Coze key.
//
// Example ECDSA Coze key:
// {
// 	"alg":"ES256",
// 	"d":"95DE8C5F50A71B392417AE0E5D60CD63AFF967FC6DA8060DECC031B0E63B3280"
// 	"kid":"Example Coze Key",
// 	"iat":1623132000,
// 	"tmb":"9EC680EEDE972F334D9B1F6775D0E61B510884DD663F982DD8323EC07D2E3FB6",
// 	"x":"C8E9E522BE0CD40B20DB86DE972B9158C227EDBE99DD2C280544C23D20728A645FE39DD3B1DDBEEA9C80A400C7CF6D2E43FFE40F660873688AAB1D676020ACBD",
// }
func (c *CozeKey) String() string {
	b, err := Marshal(c)
	if err != nil {
		return ""
	}
	return string(b)
}

// NewKey generates a new Coze Key.
func NewKey(alg ce.SEAlg) (c *CozeKey, err error) {
	c = new(CozeKey)
	c.Alg = alg

	if c.Alg.SigAlg().Genus() == ce.Ecdsa {
		eck := new(ecdsa.PrivateKey)
		switch ce.SigAlg(alg) {
		case ce.ES224:
			eck, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
		case ce.ES256:
			eck, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		case ce.ES384:
			eck, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		case ce.ES512:
			eck, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		}
		c.D = eck.D.Bytes()
		c.X = append(eck.X.Bytes()[:], eck.Y.Bytes()[:]...)
	} else if c.Alg == ce.SEAlg(ce.Ed25519) {
		var pub, pri []byte
		pub, pri, err = ed25519.GenerateKey(rand.Reader)
		c.D = pri
		c.X = pub
	} else {
		return nil, errors.New("coze.NewKey:unsupported alg")
	}
	if err != nil {
		return nil, err
	}

	c.Kid = "My Coze Key"
	c.Iat = time.Now().Unix()
	err = c.Thumbprint()

	return c, err
}

// Thumbprint generates and sets the Coze key thumbprint (`tmb`).
func (c *CozeKey) Thumbprint() (err error) {
	tmb, err := Thumbprint(c)
	if err != nil {
		return err
	}
	c.Tmb = tmb
	return
}

// Thumbprint generates Coze key thumbprint (`tmb`). For ECDSA, `tmb` is the
// digest of canon [alg, x, y] and for EdDSA `tmb` is the digest of canon [alg,
// x]
func Thumbprint(c *CozeKey) (tmb B64, err error) {
	b, err := Marshal(c)
	if err != nil {
		return nil, err
	}

	if c.Alg.SigAlg().Genus() == ce.Ecdsa {
		tmb, err = CH(b, &CozeKeyCanon{}, c.Alg.Hash())
	} else if c.Alg.SigAlg().Genus() == ce.Eddsa {
		tmb, err = CH(b, &CozeKeyCanon{}, c.Alg.Hash())
	} else {
		return nil, errors.New("coze: unknown coze key alg " + c.Alg.String() + " for thumbprint generation.")
	}

	return tmb, err
}

// Sign uses a private Coze key to sign a digest.
func (c *CozeKey) Sign(digest B64) (sig B64, err error) {
	if len(c.D) == 0 {
		return nil, errors.New("coze: `d` is not set.  Signing requires private key. ")
	}

	if c.Alg.SigAlg() == ce.Ed25519 {
		// TODO https://github.com/golang/go/issues/31804#issuecomment-1103824216
	}

	ck, err := c.ToCryptoKey()
	if err != nil {
		return nil, err
	}
	sig, err = ck.Sign(digest)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

// SignRaw uses a private Coze key to sign a pre-hashed, raw message.
func (c *CozeKey) SignRaw(msg []byte) (sig B64, err error) {
	if len(c.D) == 0 {
		return nil, errors.New("coze: `d` is not set.  Signing requires private key. ")
	}

	if c.Alg.SigAlg() == ce.Ed25519 {
		// TODO this is the wrong Ed signing function as Ed sings hashed messages and in Ed, the
		// pre-hash and the post-hash methods are different and produce different
		// results. See https://github.com/golang/go/issues/31804#issuecomment-1103824216
		return ed25519.Sign(ed25519.PrivateKey(c.D), msg), nil
	}

	ck, err := c.ToCryptoKey()
	if err != nil {
		return nil, err
	}
	sig, err = ck.SignRaw(msg)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

// SignHead signs head. Canon is optional.
//
// TODO:  Parse out MinCy and verify sanity.
func (c *CozeKey) SignHead(head interface{}, canon interface{}) (sig B64, err error) {
	b, err := Canon(head, canon)
	if err != nil {
		return nil, err
	}
	return c.SignRaw(b)
}

// SignCy signs a given Cy.Head and populates `sig`.
func (c *CozeKey) SignCy(cy *Cy, canon interface{}) (err error) {
	sig, err := c.SignHead(cy.Head, canon)
	if err != nil {
		return err
	}
	cy.Sig = sig
	return
}

// SignCyM is a convenience function for signing a Cy and returning the
// marshaled bytes of the Cy.
func SignCyM(cyer Cyer, key *CozeKey) ([]byte, error) {
	cy, err := cyer.Cy(nil)
	if err != nil {
		return nil, err
	}

	err = key.SignCy(cy, nil)
	if err != nil {
		return nil, err
	}

	return cy.MarshalJSON()
}

// Verify uses a public Coze key to verify a raw message.
func (c *CozeKey) VerifyRaw(msg []byte, sig []byte) (valid bool, err error) {
	if len(sig) == 0 {
		return false, errors.New("coze: sig is empty")
	}

	if c.Alg.SigAlg() == ce.Ed25519 {
		return ed25519.Verify(ed25519.PublicKey(c.X), msg, sig), nil
	}

	ck, err := c.ToCryptoKey()
	if err != nil {
		return false, err
	}
	return ck.Verify(msg, sig)
}

// VerifyDigest uses a public coze key to verify a digest.
//
// TODO Go's ed25519 package does not currently support verifying with a digest.
// https://pkg.go.dev/crypto/ed25519#Verify
func (c *CozeKey) VerifyDigest(digest []byte, sig []byte) (valid bool, err error) {
	if len(sig) == 0 {
		return false, errors.New("coze: sig is empty")
	}

	ck, err := c.ToCryptoKey()
	if err != nil {
		return false, err
	}

	return ck.VerifyDigest(digest, sig)
}

// Valid validates a private Coze Key and returns a bool.
//
// Valid works by
//  1. Ensuring required fields are present.
//  2. Signing a message and verifying a valid signature.
//
// Valid always returns false on public keys.  Use function "Verify" for public
// keys with signed message and "Correct" for public keys without signed
// messages.
func (c *CozeKey) Valid() (valid bool) {
	if c.D == nil || len(c.D) == 0 {
		return false
	}

	msg := []byte("Testing")
	sig, err := c.SignRaw(msg)

	if err != nil {
		return false
	}
	valid, err = c.VerifyRaw(msg, sig)
	if err != nil {
		return false
	}

	return valid
}

// Correct checks for the correct construction of a Coze key.  Correct may
// return "true" on cryptographically invalid public keys.  Use function
// "Verify" for public keys with signed message.  Correct is useful for public
// keys without signed messages.
//
// Correct:
//
// 1. Ensures required headers exist.
// 2. Checks the length of x.
// 3. Recalculates `tmb` and if incorrect throws an error.
// 4. If containing d, generates and verifies a signature, thus
//    verifying the key, by calling Valid()
func Correct(ck CozeKey) (bool, error) {
	var xLen = enum.SEAlg(ck.Alg).XSize()
	if xLen == 0 {
		return false, errors.New("coze.Correct: unknown alg")
	}

	if len(ck.X) != xLen {
		return false, fmt.Errorf("coze.Correct: x is incorrect length: %d for alg %s", len(ck.X), ck.Alg)
	}

	// Compare existing tmb
	oldTmb := ck.Tmb
	ck.Thumbprint()
	if bytes.Equal(oldTmb, ck.Tmb) {
		return false, fmt.Errorf("coze.Correct: given tmb is incorrect. Current: %X, Correct: %X", oldTmb, ck.Tmb)
	}

	// No keys from the future allowed.
	if ck.Iat > time.Now().Unix() {
		return false, errors.New("coze.Correct: keys cannot have iat greater than present time")
	}

	if ck.IsPrivate() {
		return ck.Valid(), nil
	}

	return true, nil
}

// IsPrivate reports if a Coze key should be considered private.  Any key with
// any value non-zero for `d` is considered private.
func (ck *CozeKey) IsPrivate() bool {
	if len(ck.D) > 0 {
		return true
	}
	return false
}

// ToCryptoKey takes a Coze Key object and returns a crypto key object.
func (cozekey *CozeKey) ToCryptoKey() (ck *ce.CryptoKey, err error) {
	// fmt.Printf("\n Ck Private: %+v \n", cozekey)
	if cozekey == nil {
		return nil, errors.New("coze: nil Coze Key")
	}
	if len(cozekey.X) == 0 {
		return nil, errors.New("coze: invalid CozeKey")
	}

	// TODO support Ed25519
	switch cozekey.Alg.SigAlg().Genus() {
	default:
		return nil, errors.New("unsupported alg")
	case ce.Ecdsa:
		ck, err = ecdsaCozeKeyToCryptoKey(cozekey)
		return
	}
}

// ecdsaCozeKeyToCryptoKey take a Coze Key (public or private) and returns a
// CryptoKey pair.
func ecdsaCozeKeyToCryptoKey(ck *CozeKey) (key *ce.CryptoKey, err error) {
	if ck.Alg.SigAlg().Genus() != ce.Ecdsa {
		return nil, errors.New("coze: unsupported alg for ecdsaCozeKeyToCryptoKey.")
	}

	key = new(ce.CryptoKey)
	key.Private = new(crypto.PrivateKey)
	key.Public = new(crypto.PublicKey)

	key.Alg = ck.Alg
	curve := ck.Alg.Curve().EllipticCurve()

	half := len(ck.X) / 2
	x := new(big.Int).SetBytes(ck.X[:half])
	y := new(big.Int).SetBytes(ck.X[half:])

	ec := ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}

	pub := crypto.PublicKey(ec) // set ecdsa.PublicKey to crypto.PublicKey
	key.Public = &pub

	if len(ck.D) == 0 {
		return key, err
	}

	d := new(big.Int).SetBytes(ck.D)
	var private crypto.PrivateKey
	private = ecdsa.PrivateKey{
		PublicKey: ec,
		D:         d,
	}
	key.Private = &private

	return key, err
}
