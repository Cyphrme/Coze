package coze

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"
)

// CozeKeyArrayCanon is the canonical form of a Coze Key in array form.
var CozeKeyArrayCanon = []string{"alg", "x"}

// CozeKeyCanon is the canonical form of a Coze Key in struct form.
type CozeKeyCanon struct {
	Alg string `json:"alg"`
	X   B64    `json:"x"`
}

// CozeKey is a Coze key. See `README.md` for details. Fields must be in order
// for correct JSON marshaling.
//
// Standard Coze Key Fields
// - `alg` - Specific key algorithm. E.g. "ES256" or "Ed25519".
// - `d`   - Private component.
// - `iat` - Unix time of when the key was created. E.g. 1626069600.
// - `kid` - Human readable, non-programmatic label. E.g. "My Coze key".
// - `rvk` - Unix time of key revocation. See docs on `rvk`. E.g. 1626069601.
// - `tmb` - Key thumbprint. E.g. "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk".
// - `typ` - Application label for key.  "coze/key".
// - `x`   - Public component.
type CozeKey struct {
	Alg SEAlg  `json:"alg,omitempty"`
	D   B64    `json:"d,omitempty"`
	Iat int64  `json:"iat,omitempty"`
	Kid string `json:"kid,omitempty"`
	Rvk int64  `json:"rvk,omitempty"`
	Tmb B64    `json:"tmb,omitempty"`
	Typ string `json:"typ,omitempty"`
	X   B64    `json:"x,omitempty"`
}

// String returns the stringified Coze key.
func (c *CozeKey) String() string {
	b, err := Marshal(c)
	if err != nil {
		return ""
	}
	return string(b)
}

// NewKey generates a new Coze key.
func NewKey(alg SEAlg) (c *CozeKey, err error) {
	c = new(CozeKey)
	c.Alg = alg

	switch c.Alg.SigAlg() {
	default:
		return nil, errors.New("coze.NewKey: unsupported alg: " + alg.String())
	case ES224, ES256, ES384, ES512:
		eck, err := ecdsa.GenerateKey(c.Alg.Curve().EllipticCurve(), rand.Reader)
		if err != nil {
			return nil, err
		}

		d := make([]byte, alg.DSize())
		c.D = eck.D.FillBytes(d) // Left pads bytes
		c.X = PadCon(eck.X, eck.Y, alg.XSize())
	case Ed25519:
		pub, pri, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		// ed25519.GenerateKey returns "private key" that is the seed || publicKey.
		// Remove public key for 32 byte "seed", which is used as the private key.
		c.D = []byte(pri[:32])
		c.X = B64(pub)
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
	c.Tmb, err = Thumbprint(c)
	return
}

// Thumbprint generates `tmb` which is the digest of canon [alg, x]
func Thumbprint(c *CozeKey) (tmb B64, err error) {
	b, err := Marshal(c)
	if err != nil {
		return nil, err
	}

	return CanonHash(b, &CozeKeyCanon{}, c.Alg.Hash())
}

// Sign uses a private Coze key to sign a digest.
func (c *CozeKey) Sign(digest B64) (sig B64, err error) {
	if len(c.D) == 0 {
		return nil, errors.New("Sign: private key `d` is not set.")
	}

	switch c.Alg.SigAlg().Genus() {
	default:
		return nil, errors.New("Sign: unsupported alg: " + c.Alg.String())
	case Ecdsa:
		prk := ecdsa.PrivateKey{
			PublicKey: *cozeKeyToPubEcdsa(c),
			D:         new(big.Int).SetBytes(c.D),
		}
		// Note: ECDSA Sig is always R || S of a fixed size with left padding.  For
		// example, ES256 should always have a 64 byte signature.
		r, s, err := ecdsa.Sign(rand.Reader, &prk, digest)
		if err != nil {
			return nil, err
		}
		return PadCon(r, s, c.Alg.SigAlg().SigSize()), nil
	case Eddsa:
		pk := ed25519.NewKeyFromSeed(c.D)
		// Alternatively, concat d with x
		// b := make([]coze.B64, 64)
		// d := append(b, c.D, c.X)
		return ed25519.Sign(pk, digest), nil
	}
}

// cozeKeyToPubEcdsa converts a public Coze Key to ecdsa.PublicKey.
func cozeKeyToPubEcdsa(c *CozeKey) (key *ecdsa.PublicKey) {
	half := c.Alg.XSize() / 2
	x := new(big.Int).SetBytes(c.X[:half])
	y := new(big.Int).SetBytes(c.X[half:])

	a := ecdsa.PublicKey{
		Curve: c.Alg.Curve().EllipticCurve(),
		X:     x,
		Y:     y,
	}
	return &a
}

// SignCy signs Cy.Pay and populates `sig` and verifies the cy and key fields
// `alg` and `tmb` match.  If trying to sign without `alg` and/or `tmb`, use
// Verify instead.  Canon is optional.
func (c *CozeKey) SignCy(cy *Cy, canon any) (err error) {
	// Get Coze standard fields
	h := new(Pay)
	err = json.Unmarshal(cy.Pay, h)
	if err != nil {
		return err
	}
	if c.Alg != h.Alg {
		return errors.New(fmt.Sprintf("SignCy: key alg \"%s\" and cy alg \"%s\" do not match", c.Alg, h.Alg))
	}
	if !bytes.Equal(c.Tmb, h.Tmb) {
		return errors.New(fmt.Sprintf("SignCy: key tmb \"%s\" and cy tmb  \"%s\" do not match", c.Tmb, h.Tmb))
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
	return nil
}

// Verify uses a public Coze key to verify a digest.
func (c *CozeKey) Verify(digest, sig B64) (valid bool) {
	switch c.Alg.SigAlg() {
	default:
		return false
	case ES224, ES256, ES384, ES512:
		var size = c.Alg.SigAlg().SigSize() / 2
		r := big.NewInt(0).SetBytes(sig[:size])
		s := big.NewInt(0).SetBytes(sig[size:])
		return ecdsa.Verify(cozeKeyToPubEcdsa(c), digest, r, s)
	case Ed25519, Ed25519ph:
		return ed25519.Verify(ed25519.PublicKey(c.X), digest, sig)
	}
}

// VerifyCy cryptographically verifies `pay` with given `sig` and verifies the
// `pay` and `key` fields `alg` and `tmb` match.  If trying to verify without
// `alg` and/or `tmb`, use Verify instead.
func (c *CozeKey) VerifyCy(cy *Cy) (bool, error) {
	if cy.Sig == nil {
		return false, errors.New("coze: sig is nil")
	}

	h := new(Pay)
	err := json.Unmarshal(cy.Pay, h)
	if err != nil {
		return false, err
	}
	if c.Alg != h.Alg {
		return false, errors.New(fmt.Sprintf("VerifyCy: key alg \"%s\" and cy alg \"%s\" do not match", c.Alg, h.Alg))
	}

	if !bytes.Equal(c.Tmb, h.Tmb) {
		return false, errors.New(fmt.Sprintf("VerifyCy: key tmb \"%s\" and cy tmb  \"%s\" do not match", c.Tmb, h.Tmb))
	}

	b, err := Canonical(cy.Pay, nil)
	if err != nil {
		return false, err
	}

	return c.Verify(Hash(c.Alg.Hash(), b), cy.Sig), nil
}

// Valid cryptographically validates a private Coze Key by signing a message and
// verifying the resulting signature.
//
// Valid always returns false on public keys.  Use function "Verify" for public
// keys with signed message and "Correct" for public keys without signed
// messages.
func (c *CozeKey) Valid() (valid bool) {
	if len(c.D) == 0 {
		return false
	}

	digest := Hash(c.Alg.Hash(), []byte("7AtyaCHO2BAG06z0W1tOQlZFWbhxGgqej4k9-HWP3DE-zshRbrE-69DIfgY704_FDYez7h_rEI1WQVKhv5Hd5Q"))
	sig, err := c.Sign(digest)
	if err != nil {
		return false
	}

	return c.Verify(digest, sig)
}

// Correct checks for the correct construction of a Coze key.  Key must have alg and tmb xCorrect may
// return "true" on cryptographically invalid public keys.  Use function
// "Verify" for public keys with signed message.  Correct is useful for public
// keys without signed messages, and thumb only keys.
//
// Correct:
//
// 1. Ensures required fields exist.
// 2. Checks the length of x.
// 3. Recalculates `tmb` and if incorrect throws an error.
// 4. If containing d, generates and verifies a signature, thus
//    verifying the key, by calling Valid()
func Correct(c CozeKey) (bool, error) {
	if len(c.X) == 0 {
		// Thumb only key
		if len(c.Tmb) != SEAlg(c.Alg).Hash().Size() {
			return false, errors.New("Correct: incorrect tmb size")
		}
	} else {
		// x is given
		if len(c.X) != SEAlg(c.Alg).XSize() {
			return false, errors.New("coze.Correct: incorrect x size")
		}

		// If tmb is set, recompute and compare.
		if len(c.Tmb) != 0 {
			oldTmb := c.Tmb
			c.Thumbprint()
			if bytes.Equal(oldTmb, c.Tmb) {
				return false, fmt.Errorf("coze.Correct: incorrect given tmb. Current: %X, Calculated: %X", oldTmb, c.Tmb)
			}
		}
	}

	if len(c.D) > 0 { // If private
		return c.Valid(), nil
	}

	return true, nil
}
