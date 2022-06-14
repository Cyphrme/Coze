package coze

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/cyphrme/coze/enum"
)

// CozeKeyArrayCanon is the canonical form of a Coze Key in array form.
var CozeKeyArrayCanon = []string{"alg", "x"}

// CozeKeyCanon is the canonical form of a Coze Key in struct form.
type CozeKeyCanon struct {
	Alg string `json:"alg"`
	X   B64    `json:"x"`
}

// CozeKey is a Coze key. See `README.md` for details.
//
// Fields must be in order for correct JSON marshaling.
//
// Required Fields
//	- `alg` - Specific key algorithm. E.g. "ES256" or "Ed25519".
//
// Recommended and Optional Fields:
//	- `kid` - Human readable label and must not be used programmatically. E.g. "My Coze key".
//	- `iat` - Unix time of when the key was created. E.g. 1626069600.
//	- `tmb` - Key's thumbprint. E.g. "0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD".
//	- `typ` - The key's type and may be used by applications to identify the key.  "coze/key".
//
// Key Fields
//	- `d` - Private component.
//	- `x` - Public component.
//
// Revoked.  Key revocation should be done through a Coze message using the
// `rvk`` field in `pay` (See the Coze README).  The Coze key field `rvk` is
// useful for storing a key's revocation state.
//  - `rvk` - Unix time of key revocation. See docs on `rvk`. E.g. 1626069601.
//
type CozeKey struct {
	Alg enum.SEAlg `json:"alg"`
	D   B64        `json:"d,omitempty"`
	Iat int64      `json:"iat,omitempty"`
	Kid string     `json:"kid,omitempty"`
	Rvk int64      `json:"rvk,omitempty"`
	Tmb B64        `json:"tmb,omitempty"`
	Typ string     `json:"typ,omitempty"`
	X   B64        `json:"x,omitempty"`
}

// String returns the stringified Coze key.
func (c *CozeKey) String() string {
	b, err := Marshal(c)
	if err != nil {
		return ""
	}
	return string(b)
}

// NewKey generates a new Coze Key.
func NewKey(alg enum.SEAlg) (c *CozeKey, err error) {
	c = new(CozeKey)
	c.Alg = alg

	if c.Alg.SigAlg().Genus() == enum.Ecdsa {
		eck := new(ecdsa.PrivateKey)
		switch enum.SigAlg(alg) {
		case enum.ES224:
			eck, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
		case enum.ES256:
			eck, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		case enum.ES384:
			eck, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		case enum.ES512:
			eck, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		}

		if err != nil {
			return nil, err
		}

		d := make([]byte, alg.DSize())
		c.D = eck.D.FillBytes(d) // Left pads bytes
		c.X = enum.PadCon(eck.X, eck.Y, alg.XSize())

	} else if c.Alg == enum.SEAlg(enum.Ed25519) {
		var pub, pri []byte
		pub, pri, err = ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		c.D = pri
		c.X = pub
	} else {
		return nil, errors.New("coze.NewKey:unsupported alg")
	}

	c.Iat = time.Now().Unix()
	err = c.Thumbprint()
	if err != nil {
		return nil, err
	}

	return c, nil
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

// Thumbprint generates Coze key thumbprint `tmb` which is the digest of canon
// [alg, x]
func Thumbprint(c *CozeKey) (tmb B64, err error) {
	b, err := Marshal(c)
	if err != nil {
		return nil, err
	}

	tmb, err = CanonHash(b, &CozeKeyCanon{}, c.Alg.Hash())

	return tmb, err
}

// Sign uses a private Coze key to sign a digest.
func (c *CozeKey) Sign(digest B64) (sig B64, err error) {
	if len(c.D) == 0 {
		return nil, errors.New("coze Sign: private key `d` is not set.")
	}

	if c.Alg.SigAlg() == enum.Ed25519 {
		// TODO Coze signs hashed messages and "pure" Ed signs messages.  Ed's
		// pre-hash and the post-hash methods are different and produce different
		// TODO https://github.com/golang/go/issues/31804#issuecomment-1103824216

		// TODO Go's ed25519 package currently does not currently support verifying with a digest.
		// https://pkg.go.dev/crypto/ed25519#Verify
		return nil, errors.New("Ed25519 is currently unsupported")
	}

	ck, err := c.ToCryptoKey()
	if err != nil {
		return nil, err
	}
	// Cryptokey handles all error checking.
	sig, err = ck.Sign(digest)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

// SignCy signs Cy.Pay and populates `sig`.  Canon is optional.
func (c *CozeKey) SignCy(cy *Cy, canon any) (err error) {
	// Get Coze standard fields
	h := new(Pay)
	err = json.Unmarshal(cy.Pay, h)
	if err != nil {
		return err
	}
	if !bytes.Equal(h.Tmb, c.Tmb) {
		return errors.New("coze: key tmb and cy tmb do not match")
	}

	b, err := Canonical(cy.Pay, canon) // compactify
	if err != nil {
		return err
	}

	sig, err := c.Sign(Hash(c.Alg.Hash(), b))
	if err != nil {
		return err
	}
	cy.Sig = sig
	return
}

// Verify uses a public coze key to verify a digest.
//
func (c *CozeKey) Verify(digest, sig B64) (valid bool) {
	ck, err := c.ToCryptoKey()
	if err != nil {
		return false
	}

	return ck.Verify(digest, sig)
}

// Verify cryptographically verifies `coze` with given `sig`.  Canon is optional.
func (ck *CozeKey) VerifyCy(cy *Cy) (bool, error) {
	if cy.Sig == nil {
		return false, errors.New("coze: sig is nil")
	}

	h := new(Pay)
	err := json.Unmarshal(cy.Pay, h)
	if err != nil {
		return false, err
	}
	if !bytes.Equal(h.Tmb, ck.Tmb) {
		return false, errors.New("coze: key tmb and cy tmb do not match")
	}

	b, err := Canonical(cy.Pay, nil)
	if err != nil {
		return false, err
	}

	return ck.Verify(Hash(ck.Alg.Hash(), b), cy.Sig), nil
}

// Valid validates a private Coze Key by signing a message and verifying a valid
// signature.
//
// Valid always returns false on public keys.  Use function "Verify" for public
// keys with signed message and "Correct" for public keys without signed
// messages.
func (c *CozeKey) Valid() (valid bool) {
	if c.D == nil || len(c.D) == 0 {
		return false
	}
	// Random message
	msg := []byte("a0HwToVezVCBrucf3RiBW4xDWSnap1GZTvNfkI7q77k")
	digest := Hash(c.Alg.Hash(), msg)

	sig, err := c.Sign(digest)
	if err != nil {
		return false
	}
	valid = c.Verify(digest, sig)

	return valid
}

// Correct checks for the correct construction of a Coze key.  Correct may
// return "true" on cryptographically invalid public keys.  Use function
// "Verify" for public keys with signed message.  Correct is useful for public
// keys without signed messages.
//
// Correct:
//
// 1. Ensures required fields exist.
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

// ToCryptoKey takes a Coze Key and returns a crypto key.  Organizationally,
// this function would be better in the enum package but that would result in an
// import cycle.
func (cozekey *CozeKey) ToCryptoKey() (ck *enum.CryptoKey, err error) {
	//fmt.Printf("\n Ck Private: %+v \n", cozekey)
	if cozekey == nil {
		return nil, errors.New("coze: nil Coze Key")
	}
	if len(cozekey.X) == 0 {
		return nil, errors.New("coze: invalid CozeKey")
	}

	switch cozekey.Alg.SigAlg().Genus() {
	default:
		return nil, errors.New("unsupported alg")
	case enum.Ecdsa:
		ck, err = ecDSACozeKeyToCryptoKey(cozekey)
		return
	case enum.Eddsa:
		ck, err = edDSACozeKeyToCryptoKey(cozekey)
		return
	}
}

func edDSACozeKeyToCryptoKey(ck *CozeKey) (key *enum.CryptoKey, err error) {
	// TODO support Ed25519
	return nil, nil
}

// ecdsaCozeKeyToCryptoKey take a Coze Key (public or private) and returns a
// CryptoKey pair.  Organizationally, this function would be better in the enum
// package but that would result in an import cycle.
func ecDSACozeKeyToCryptoKey(ck *CozeKey) (key *enum.CryptoKey, err error) {
	if ck.Alg.SigAlg().Genus() != enum.Ecdsa {
		return nil, errors.New("coze: unsupported alg for ecdsaCozeKeyToCryptoKey.")
	}

	key = new(enum.CryptoKey)
	key.Private = new(crypto.PrivateKey)
	key.Public = new(crypto.PublicKey)

	key.Alg = ck.Alg
	curve := ck.Alg.Curve().EllipticCurve()

	half := ck.Alg.XSize() / 2
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
