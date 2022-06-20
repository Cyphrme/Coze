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

// NewKey generates a new Coze key.
func NewKey(alg SEAlg) (c *CozeKey, err error) {
	c = new(CozeKey)
	c.Alg = alg

	switch c.Alg.SigAlg() {
	default:
		return nil, errors.New("NewKey: unsupported alg: " + alg.String())
	case ES224, ES256, ES384, ES512:
		eck, err := ecdsa.GenerateKey(c.Alg.Curve().EllipticCurve(), rand.Reader)
		if err != nil {
			return nil, err
		}

		d := make([]byte, alg.DSize())
		c.D = eck.D.FillBytes(d) // Left pads bytes
		c.X = PadInts(eck.X, eck.Y, alg.XSize())
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

// Thumbprint generates and sets the Coze key thumbprint (`tmb`) from `x` and `alg``.
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
			// ecdsa.Sign only needs PublicKey.Curve, not it's value.
			PublicKey: ecdsa.PublicKey{Curve: c.Alg.Curve().EllipticCurve()},
			D:         new(big.Int).SetBytes(c.D),
		}
		r, s, err := ecdsa.Sign(rand.Reader, &prk, digest)
		if err != nil {
			return nil, err
		}
		// ECDSA Sig is R || S rounded up to byte left padded.
		return PadInts(r, s, c.Alg.SigAlg().SigSize()), nil
	case Eddsa:
		pk := ed25519.NewKeyFromSeed(c.D)
		// Alternatively, concat d with x
		// b := make([]coze.B64, 64)
		// d := append(b, c.D, c.X)
		return ed25519.Sign(pk, digest), nil
	}
}

// cozeKeyToPubEcdsa converts a Coze Key to ecdsa.PublicKey.
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
	if len(c.X) == 0 {
		return false
	}

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

// Correct checks for the correct construction of a Coze key.  Correct is useful
// for sanity checking public keys without signed messages, sanity checking tmb
// only keys, and validating private keys.
//
//
// Key must have `alg` and any, but at
// least one, of `tmb`, `x`, and `d`. Correct may return "true" on
// cryptographically invalid public keys.  See also function "Verify" for
// verifying public keys when a signed message is available.
//
// Correct:
//
// 1. Checks the length of x against alg.
// 2. If x and tmb are present, verifies correct `tmb`.
// 3. If d is present, verifies correct tmb and x if present, and generates and
// verifies a signature, thus verifying the key
func (c *CozeKey) Correct() (bool, error) {
	if c.Alg == 0 {
		return false, errors.New("Correct: Alg must be set")
	}

	if len(c.Tmb) == 0 && len(c.X) == 0 && len(c.D) == 0 {
		return false, errors.New("Correct: At least one of [x, tmb, d] must be set")
	}

	// tmb only key
	if len(c.X) == 0 && len(c.D) == 0 {
		if len(c.Tmb) != SEAlg(c.Alg).Hash().Size() {
			return false, fmt.Errorf("Correct: incorrect tmb size: %d", len(c.Tmb))
		}
		return true, nil
	}

	//  d is not set
	if len(c.D) == 0 {
		if len(c.X) != 0 && len(c.X) != SEAlg(c.Alg).XSize() {
			return false, fmt.Errorf("Correct: incorrect x size: %d", len(c.X))
		}
		// If tmb is set, recompute and compare.
		if len(c.Tmb) != 0 {
			tmb, err := Thumbprint(c)
			if err != nil {
				return false, err
			}
			if !bytes.Equal(c.Tmb, tmb) {
				return false, fmt.Errorf("Correct: incorrect given tmb. Current: %s, Calculated: %s", c.Tmb, tmb)
			}
		}
		return true, nil
	}

	// If d and (x and/or tmb) is given, recompute from d and compare:
	x := c.recalcX()
	if len(c.X) != 0 && !bytes.Equal(c.X, x) {
		return false, fmt.Errorf("Correct: incorrect X. Current: %s, Calculated: %s", c.X, x)
	}
	var ck = CozeKey{Alg: c.Alg, X: x}
	// If tmb is set, recompute and compare with existing.
	if len(c.Tmb) != 0 {
		tmb, err := Thumbprint(&ck)
		if err != nil {
			return false, err
		}
		if !bytes.Equal(c.Tmb, tmb) {
			return false, fmt.Errorf("Correct: incorrect given tmb. Current: %s, Calculated: %s", c.Tmb, tmb)
		}
	}
	ck.D = c.D
	return ck.Valid(), nil
}

// recalcX recalculates x from d.  Algorithms are constant-time.
// https://cs.opensource.google/go/go/+/refs/tags/go1.18.3:src/crypto/elliptic/elliptic.go;l=455;drc=7f9494c277a471f6f47f4af3036285c0b1419816
func (c *CozeKey) recalcX() (x B64) {
	switch c.Alg.SigAlg() {
	default:
		return nil
	case ES224, ES256, ES384, ES512:
		pukx, puky := c.Alg.Curve().EllipticCurve().ScalarBaseMult(c.D)
		x = PadInts(pukx, puky, c.Alg.XSize())
	case Ed25519, Ed25519ph:
		prk := ed25519.NewKeyFromSeed(c.D)
		x = []byte(prk[:32])
	}
	return x
}

// String returns the stringified Coze key.
func (c CozeKey) String() string {
	b, err := Marshal(c)
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}
