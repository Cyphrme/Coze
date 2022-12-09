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

// KeyCanonSlice is the canonical form of a Coze key in slice form.
var KeyCanon = []string{"alg", "x"}

// Key is a Coze key. See `README.md` for details on Coze key. Fields must be in
// order for correct JSON marshaling.
//
// Standard Coze key Fields
//   - `alg` - Specific key algorithm. E.g. "ES256" or "Ed25519".
//   - `d`   - Private component. E.g. "bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVA".
//   - `iat` - Unix time of when the key was created. E.g. 1626069600.
//   - `kid` - Human readable, non-programmatic label. E.g. "My Coze key".
//   - `rvk` - Unix time of key revocation. See docs on `rvk`. E.g. 1626069601.
//   - `tmb` - Key thumbprint. E.g. "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk".
//   - `typ` - Application label for key. E.g. "coze/key".
//   - `x`   - Public component. E.g. "2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g".
type Key struct {
	Alg SEAlg  `json:"alg,omitempty"`
	D   B64    `json:"d,omitempty"`
	Iat int64  `json:"iat,omitempty"`
	Kid string `json:"kid,omitempty"`
	Rvk int64  `json:"rvk,omitempty"`
	Tmb B64    `json:"tmb,omitempty"`
	Typ string `json:"typ,omitempty"`
	X   B64    `json:"x,omitempty"`
}

// String implements Stringer. Returns empty on error.
func (c Key) String() string {
	b, err := Marshal(c)
	if err != nil {
		return ""
	}
	return string(b)
}

// NewKey generates a new Coze key.
func NewKey(alg SEAlg) (c *Key, err error) {
	c = new(Key)
	c.Alg = alg

	switch c.Alg.SigAlg() {
	default:
		return nil, fmt.Errorf("NewKey: unsupported alg: %s", alg)
	case ES224, ES256, ES384, ES512:
		eck, err := ecdsa.GenerateKey(c.Alg.Curve().EllipticCurve(), rand.Reader)
		if err != nil {
			return nil, err
		}

		d := make([]byte, alg.DSize())
		c.D = eck.D.FillBytes(d) // Left pads bytes
		c.X = PadInts(eck.X, eck.Y, alg.XSize())
	case Ed25519, Ed25519ph:
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
	return c, c.Thumbprint()
}

// Thumbprint generates and sets the Coze key thumbprint (`tmb`) from `x` and
// `alg`.
func (c *Key) Thumbprint() (err error) {
	c.Tmb, err = Thumbprint(c)
	return err
}

// Thumbprint generates `tmb` which is the digest of canon [alg, x].
func Thumbprint(c *Key) (tmb B64, err error) {
	b, err := Marshal(c)
	if err != nil {
		return nil, err
	}
	return CanonicalHash(b, KeyCanon, c.Alg.Hash())
}

// Sign uses a private Coze key to sign a digest.
//
// Sign() and Verify() do not check if the Coze is correct, such as checking
// pay.alg and pay.tmb matches with Key.  Use SignPay, SignCoze, SignPayJSON,
// and/or VerifyCoze if needing Coze validation.
func (c *Key) Sign(digest B64) (sig B64, err error) {
	if len(c.D) != c.Alg.DSize() {
		return nil, fmt.Errorf("Sign: Invalid `d` length %d", len(c.D))
	}

	switch c.Alg.SigAlg().Genus() {
	default:
		return nil, fmt.Errorf("Sign: unsupported alg: %s", c.Alg)
	case ECDSA:
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
	case EdDSA:
		pk := ed25519.NewKeyFromSeed(c.D)
		// Alternatively, concat d with x
		// b := make([]coze.B64, 64)
		// d := append(b, c.D, c.X)
		return ed25519.Sign(pk, digest), nil
	}
}

// SignPay signs coze.Pay and returns a new Coze with coze.Sig populated. If set
// SignPay checks that `pay.alg` and `key.alg` match and that `pay.tmb` is
// correct according to `key`.
//
// SignPay works with contextual cozies that lack pay.alg and/or pay.tmb and
// uses key as a source of truth.
func (c *Key) SignPay(p *Pay) (coze *Coze, err error) {
	if p.Alg != "" && c.Alg != p.Alg {
		return nil, fmt.Errorf("SignPay: key alg \"%s\" and coze alg \"%s\" do not match", c.Alg, p.Alg)
	}
	if len(p.Tmb) != 0 && !bytes.Equal(c.Tmb, p.Tmb) {
		return nil, fmt.Errorf("SignPay: key tmb \"%s\" and coze tmb  \"%s\" do not match", c.Tmb, p.Tmb)
	}

	b, err := Marshal(p)
	if err != nil {
		return nil, err
	}

	coze, err = p.Coze()
	if err != nil {
		return nil, err
	}

	d, err := Hash(c.Alg.Hash(), b)
	if err != nil {
		return nil, err
	}

	coze.Sig, err = c.Sign(d)
	return coze, err
}

// SignPayJSON signs a json `coze.pay`.  See documentation on SignPay.
func (c *Key) SignPayJSON(pay json.RawMessage) (coze *Coze, err error) {
	p := new(Pay)
	err = json.Unmarshal(pay, p)
	if err != nil {
		return nil, err
	}

	if p.Alg != "" && c.Alg != p.Alg {
		return nil, fmt.Errorf("SignPay: key alg \"%s\" and coze alg \"%s\" do not match", c.Alg, p.Alg)
	}
	if len(p.Tmb) != 0 && !bytes.Equal(c.Tmb, p.Tmb) {
		return nil, fmt.Errorf("SignPay: key tmb \"%s\" and coze tmb  \"%s\" do not match", c.Tmb, p.Tmb)
	}

	b, err := compact(pay)
	if err != nil {
		return nil, err
	}

	d, err := Hash(c.Alg.Hash(), b)
	if err != nil {
		return nil, err
	}

	coze = new(Coze)
	coze.Pay = b
	coze.Sig, err = c.Sign(d)
	return coze, err
}

// SignCoze signs `coze.pay` and sets `coze.sig`.  See documentation on SignPay.
func (c *Key) SignCoze(cz *Coze) (err error) {
	coze, err := c.SignPayJSON(cz.Pay)
	if err != nil {
		return err
	}
	cz.Sig = coze.Sig
	return nil
}

// Verify uses a Coze key to verify a digest.
//
// Sign() and Verify() do not check if the Coze is correct, such as checking
// pay.alg and pay.tmb matches with Key.  Use SignPay, SignCoze, SignPayJSON,
// and/or VerifyCoze if needing Coze validation.
func (c *Key) Verify(digest, sig B64) (valid bool) {
	if len(c.X) != c.Alg.XSize() {
		return false
	}
	switch c.Alg.SigAlg() {
	default:
		return false
	case ES224, ES256, ES384, ES512:
		size := c.Alg.SigAlg().SigSize() / 2
		r := big.NewInt(0).SetBytes(sig[:size])
		s := big.NewInt(0).SetBytes(sig[size:])
		return ecdsa.Verify(KeyToPubEcdsa(c), digest, r, s)
	case Ed25519, Ed25519ph:
		return ed25519.Verify(ed25519.PublicKey(c.X), digest, sig)
	}
}

// VerifyCoze cryptographically verifies `pay` with given `sig`.  If set
// VerifyCoze checks that `pay.alg` and `key.alg` match and that `pay.tmb` is
// correct according to `key`. Always returns false on error.
//
// VerifyCoze works with contextual cozies that lack pay.alg and/or
// pay.tmb and uses key as a source of truth.
func (c *Key) VerifyCoze(cz *Coze) (bool, error) {
	p := new(Pay)
	err := json.Unmarshal(cz.Pay, p)
	if err != nil {
		return false, err
	}
	if p.Alg != "" && c.Alg != p.Alg {
		return false, fmt.Errorf("VerifyCoze: key alg \"%s\" and coze alg \"%s\" do not match", c.Alg, p.Alg)
	}
	if len(p.Tmb) != 0 && !bytes.Equal(c.Tmb, p.Tmb) {
		return false, fmt.Errorf("VerifyCoze: key tmb \"%s\" and coze tmb  \"%s\" do not match", c.Tmb, p.Tmb)
	}

	b, err := compact(cz.Pay)
	if err != nil {
		return false, err
	}

	d, err := Hash(c.Alg.Hash(), b)
	if err != nil {
		return false, err
	}

	return c.Verify(d, cz.Sig), nil
}

// Valid cryptographically validates a private Coze Key by signing a message and
// verifying the resulting signature with the given "x".
//
// Valid always returns false on public keys.  Use function "Verify" for public
// keys with signed message.  See also function Correct.
func (c *Key) Valid() (valid bool) {
	d, err := Hash(c.Alg.Hash(), []byte("7AtyaCHO2BAG06z0W1tOQlZFWbhxGgqej4k9-HWP3DE-zshRbrE-69DIfgY704_FDYez7h_rEI1WQVKhv5Hd5Q"))
	if err != nil {
		return false
	}
	sig, err := c.Sign(d)
	if err != nil {
		return false
	}
	return c.Verify(d, sig)
}

// Correct checks for the correct construction of a Coze key, but may return
// true on cryptographically invalid public keys.  Key must have `alg` and at
// least one of `tmb`, `x`, and `d`. Using input information, if possible to
// definitively know the given key is incorrect, Correct returns false, but if
// plausibly correct, Correct returns true. Correct answers the question: "Is
// the given Coze key reasonable using the information provided?". Correct is
// useful for sanity checking public keys without signed messages, sanity
// checking `tmb` only keys, and validating private keys. Use function "Verify"
// instead for verifying public keys when a signed message is available. Correct
// is considered an advanced function. Please understand it thoroughly before
// use.
//
// Correct:
//
//  1. Checks the length of `x` and/or `tmb` against `alg`.
//  2. If `x` and `tmb` are present, verifies correct `tmb`.
//  3. If `d` is present, verifies correct `tmb` and `x` if present, and
//     verifies the key by verifying a generated signature.
func (c *Key) Correct() (bool, error) {
	if c.Alg == "" {
		return false, errors.New("Correct: Alg must be set")
	}

	if len(c.Tmb) == 0 && len(c.X) == 0 && len(c.D) == 0 {
		return false, errors.New("Correct: At least one of [x, tmb, d] must be set")
	}

	// tmb only key
	if len(c.X) == 0 && len(c.D) == 0 {
		if len(c.Tmb) != c.Alg.Hash().Size() {
			return false, fmt.Errorf("Correct: incorrect tmb size: %d", len(c.Tmb))
		}
		return true, nil
	}

	// d is not set
	if len(c.D) == 0 {
		if len(c.X) != 0 && len(c.X) != c.Alg.XSize() {
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

	// If d and (x and/or tmb) is given, recompute from d and compare.
	x := c.recalcX()
	if len(c.X) != 0 && !bytes.Equal(c.X, x) {
		return false, fmt.Errorf("Correct: incorrect X. Current: %s, Calculated: %s", c.X, x)
	}
	ck := Key{Alg: c.Alg, X: x}
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

// recalcX recalculates 'x' from 'd' and returns 'x'. 'x' will not be set on the
// key from here. Algorithms are constant-time.
// https://cs.opensource.google/go/go/+/refs/tags/go1.18.3:src/crypto/elliptic/elliptic.go;l=455;drc=7f9494c277a471f6f47f4af3036285c0b1419816
func (c *Key) recalcX() B64 {
	switch c.Alg.SigAlg() {
	default:
		return nil
	case ES224, ES256, ES384, ES512:
		pukx, puky := c.Alg.Curve().EllipticCurve().ScalarBaseMult(c.D)
		return PadInts(pukx, puky, c.Alg.XSize())
	case Ed25519, Ed25519ph:
		return []byte(ed25519.NewKeyFromSeed(c.D)[32:])
	}
}

// KeyToPubEcdsa converts a Coze Key to ecdsa.PublicKey.
func KeyToPubEcdsa(c *Key) (key *ecdsa.PublicKey) {
	size := c.Alg.XSize() / 2
	return &ecdsa.PublicKey{
		Curve: c.Alg.Curve().EllipticCurve(),
		X:     new(big.Int).SetBytes(c.X[:size]),
		Y:     new(big.Int).SetBytes(c.X[size:]),
	}
}
