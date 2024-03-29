package coze

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"
)

// KeyCanon is the canonical form of a Coze key.
var KeyCanon = []string{"alg", "x"}

// Key is a Coze key. See `README.md` for details on Coze key. Fields `alg` and
// `tmb` must be in correct relative order for thumbprint canon because JSON
// marshal uses struct order.
//
// Standard Coze key Fields
//
//	`alg` - Specific key algorithm. E.g. "ES256" or "Ed25519".
//	`d`   - Private component. E.g. "bNstg4_H3m3SlROufwRSEgibLrBuRq9114OvdapcpVA".
//	`iat` - Unix time of when the key was created. E.g. 1626069600.
//	`kid` - Human readable, non-programmatic label. E.g. "My Coze key".
//	`rvk` - Unix time of key revocation. See docs on `rvk`. E.g. 1626069601.
//	`tmb` - Key thumbprint. E.g. "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk".
//	`typ` - Application label for key. E.g. "coze/key".
//	`x`   - Public component. E.g. "2nTOaFVm2QLxmUO_SjgyscVHBtvHEfo2rq65MvgNRjORojq39Haq9rXNxvXxwba_Xj0F5vZibJR3isBdOWbo5g".
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
		return nil, fmt.Errorf("NewKey: unsupported alg %q", alg)
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

// Thumbprint generates `tmb` which is the digest of canon [alg, x].  X must be
// set and be a valid length.  On error, tmb is set to nil.
func (c *Key) Thumbprint() (err error) {
	c.Tmb, err = Thumbprint(c)
	return err
}

// Thumbprint generates `tmb` which is the digest of canon [alg, x].  X must be
// set and be a valid length.  On error, tmb is set to nil.
func Thumbprint(c *Key) (tmb B64, err error) {
	if len(c.X) != c.Alg.XSize() {
		return nil, fmt.Errorf("Thumbprint: incorrect x length for alg %q; expected %q; given %q", c.Alg, c.Alg.Hash().Size(), len(tmb))
	}
	b, err := Marshal(c)
	if err != nil {
		return nil, err
	}
	return CanonicalHash(b, KeyCanon, c.Alg.Hash())
}

// UnmarshalJSON always populates `tmb` even if it isn't given.
func (c *Key) UnmarshalJSON(b []byte) error {
	err := checkDuplicate(json.NewDecoder(bytes.NewReader(b)))
	if err != nil {
		return err
	}

	type key2 Key // Break infinite unmarshal loop
	czk2 := new(key2)
	err = json.Unmarshal(b, czk2)
	if err != nil {
		return err
	}

	*c = *(*Key)(czk2)
	err = c.Correct() // Correct sets tmb.
	if err != nil {
		return err
	}
	return nil
}

// Sign uses a private Coze key to sign a digest.
//
// Sign() and Verify() do not check if the Coze is correct, such as checking
// pay.alg and pay.tmb matches with Key.  Use SignPay, SignCoze, SignPayJSON,
// and/or VerifyCoze if needing Coze validation.
func (c *Key) Sign(digest B64) (sig B64, err error) {
	if len(c.D) != c.Alg.DSize() {
		return nil, fmt.Errorf("Sign: incorrect d length for alg %q; expected %q, given %q", c.Alg, c.Alg.DSize(), len(c.D))
	}

	switch c.Alg.SigAlg().Genus() {
	default:
		return nil, fmt.Errorf("Sign: unsupported alg %q", c.Alg)
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

		// S canonicalization generates signature with low-S.
		err = ToLowS(c, s)
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

// SignPay signs coze.Pay and returns a new Coze with coze.Sig populated. If set,
// SignPay checks that `pay.alg` and `key.alg` match and that `pay.tmb` is
// correct according to `key`.
//
// SignPay works with contextual cozies that lack pay.alg and/or pay.tmb and
// uses key as a source of truth.
func (c *Key) SignPay(p *Pay) (coze *Coze, err error) {
	return c.signPayJSON(p, nil)
}

// SignPayJSON signs a json `coze.pay`.  See documentation on SignPay.
func (c *Key) SignPayJSON(pay json.RawMessage) (coze *Coze, err error) {
	p := new(Pay)
	err = json.Unmarshal(pay, p)
	if err != nil {
		return nil, err
	}
	return c.signPayJSON(p, pay)
}

// signPayJSON efficiently consolidates common code between SignPay and
// SignPayJSON. Parameter p must be given and b is optional.  If b is nil, b is
// generated from p. If b is not nil b is compacted.
func (c *Key) signPayJSON(p *Pay, b json.RawMessage) (coze *Coze, err error) {
	if p.Alg != "" && c.Alg != p.Alg {
		return nil, fmt.Errorf("SignPay: key alg %q and coze alg %q do not match", c.Alg, p.Alg)
	}
	if len(p.Tmb) != 0 && !bytes.Equal(c.Tmb, p.Tmb) {
		return nil, fmt.Errorf("SignPay: key tmb %q and coze tmb %q do not match", c.Tmb, p.Tmb)
	}

	if b == nil {
		b, err = Marshal(p)
		if err != nil {
			return nil, err
		}
	} else {
		b, err = compact(b)
		if err != nil {
			return nil, err
		}
	}

	d, err := Hash(c.Alg.Hash(), b)
	if err != nil {
		return nil, err
	}
	sig, err := c.Sign(d)
	if err != nil {
		return nil, err
	}

	coze = new(Coze)
	coze.Pay = b
	coze.Sig = sig
	return coze, nil
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

// Verify uses a Coze key to verify a digest.  Typically digest is `cad`.
//
// Sign() and Verify() do not check if the coze is correct, such as checking
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

		// S canonicalization. Only accept low-S.
		lowS, err := IsLowS(c, s)
		if !lowS || err != nil {
			return false
		}

		return ecdsa.Verify(c.ToPubEcdsa(), digest, r, s)
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
		return false, fmt.Errorf("VerifyCoze: key.alg %q and coze.alg %q do not match", c.Alg, p.Alg)
	}
	if len(p.Tmb) != 0 && !bytes.Equal(c.Tmb, p.Tmb) {
		return false, fmt.Errorf("VerifyCoze: key tmb %q and coze tmb %q do not match", c.Tmb, p.Tmb)
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
	// fmt.Printf("Valid key: %v\n", c)
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

// Correct is an advanced function for checking for the correct construction of
// a Coze key if it can be known from the given inputs. Key must have at least
// one of [`tmb`, `x`,`d`] and `alg` set.  Correct may return no error on
// cryptographically invalid public keys.  Using input information, if possible
// to definitively know the given key is incorrect, Correct returns an error,
// but if plausibly correct, Correct returns no error. Correct answers the
// question: "Is the given Coze key reasonable using the information provided?".
// Correct is useful for sanity checking public keys without signed messages,
// sanity checking `tmb` only keys, and validating private keys. Use function
// "Verify" instead for verifying public keys when a signed message is
// available. Correct is considered an advanced function. Please understand it
// thoroughly before use.
//
// Correct:
//
//  1. Checks the length of `x` and/or `tmb` against `alg`.
//  2. If `x` and `tmb` are present, verifies correct `tmb`.
//  3. If `d` is present, verifies correct `tmb` and `x` if present, and
//     verifies the key by verifying a generated signature.
//  4. If possible, sets tmb and/or x.
//
// Functions that call correct can check for correctness by `if key.Correct() != nil`
func (c *Key) Correct() (err error) {
	if c.Alg == "" {
		return errors.New("Correct: alg must be set")
	}
	if len(c.Tmb) == 0 && len(c.X) == 0 && len(c.D) == 0 {
		return errors.New("Correct: at least one of [x, tmb, d] must be set")
	}

	// d is set.
	// Calculate x from d and compare with given value.
	if len(c.D) != 0 {
		givenX := c.X
		c.X = c.calcX()
		if len(givenX) != 0 && !bytes.Equal(c.X, givenX) {
			return fmt.Errorf("Correct: incorrect X; expected %q, given %q, ", c.X, givenX)
		}
		if !c.Valid() {
			return fmt.Errorf("Correct: key is invalid")
		}
	}

	// x is set.
	// Calculate tmb from x and compare with given value.
	if len(c.X) != 0 {
		if len(c.X) != c.Alg.XSize() {
			return fmt.Errorf("Correct: incorrect x length for alg %q; expected %q, given %q", c.Alg, c.Alg.XSize(), len(c.X))
		}
		givenTmb := c.Tmb
		err := c.Thumbprint()
		if err != nil {
			return err
		}
		if len(givenTmb) != 0 && !bytes.Equal(c.Tmb, givenTmb) {
			return fmt.Errorf("Correct: incorrect tmb; expected %q, given %q", c.Tmb, givenTmb)
		}
	}

	// tmb only key.  (Coze assumes `x` is calculable from `d`, so at this point
	// `tmb` should always be set. See `checksum_and_seed.md` for exposition.
	if len(c.Tmb) != c.Alg.Hash().Size() {
		return fmt.Errorf("Correct: incorrect tmb length for alg %q; expected %q, given %q", c.Alg, c.Alg.Hash().Size(), len(c.Tmb))
	}
	return nil
}

// Revoke returns a signed revoke coze and sets `rvk` on the key itself.
func (c *Key) Revoke() (coze *Coze, err error) {
	err = c.Correct()
	if err != nil {
		return nil, fmt.Errorf("Revoke: Coze key is not correct; %s", err)
	}

	r := new(Pay)
	r.Alg = c.Alg
	r.Iat = time.Now().Unix()
	r.Rvk = r.Iat
	r.Tmb = c.Tmb
	// If needing "typ" populated, use Sign.

	coze = new(Coze)
	coze.Pay, err = r.MarshalJSON()
	if err != nil {
		return
	}

	err = c.SignCoze(coze)
	if err != nil {
		return nil, err
	}
	c.Rvk = r.Iat // Sets `Key.Rvk` to the same value as the self-revoke coze.
	return coze, nil
}

// IsRevoked returns true if the given Key is marked as revoked.
func (c Key) IsRevoked() bool {
	return isRevoke(c.Rvk)
}

// recalcX recalculates 'x' from 'd' and returns 'x'. 'x' will not be set on the
// key from here. Algorithms are constant-time.
// https://cs.opensource.google/go/go/+/refs/tags/go1.18.3:src/crypto/elliptic/elliptic.go;l=455;drc=7f9494c277a471f6f47f4af3036285c0b1419816
func (c *Key) calcX() B64 {
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

// ToPubEcdsa converts a Coze Key to ecdsa.PublicKey.
func (c *Key) ToPubEcdsa() (key *ecdsa.PublicKey) {
	size := c.Alg.XSize() / 2
	return &ecdsa.PublicKey{
		Curve: c.Alg.Curve().EllipticCurve(),
		X:     new(big.Int).SetBytes(c.X[:size]),
		Y:     new(big.Int).SetBytes(c.X[size:]),
	}
}

// curveOrders contains curve group orders.
var curveOrders = map[SigAlg]*big.Int{
	ES224: elliptic.P224().Params().N,
	ES256: elliptic.P256().Params().N,
	ES384: elliptic.P384().Params().N,
	ES512: elliptic.P521().Params().N,
}

// curveHalfOrders contains curve group orders halved for ToLowS.  From
// https://github.com/golang/go/issues/54549
var curveHalfOrders = map[SigAlg]*big.Int{
	// Logical right shift divides a number by 2 discreetly.
	ES224: new(big.Int).Rsh(elliptic.P224().Params().N, 1),
	ES256: new(big.Int).Rsh(elliptic.P256().Params().N, 1),
	ES384: new(big.Int).Rsh(elliptic.P384().Params().N, 1),
	ES512: new(big.Int).Rsh(elliptic.P521().Params().N, 1),
}

// IsLowS checks if S is a low-S for ECDSA.  See Coze docs on low-S.
func IsLowS(c *Key, s *big.Int) (bool, error) {
	if c.Alg.Genus() != ECDSA {
		return false, fmt.Errorf("IsLowS: alg %q is not ECDSA", c.Alg)
	}
	return s.Cmp(curveHalfOrders[c.Alg.SigAlg()]) != 1, nil
}

// ToLowS converts high-S to low-S or if already low-S returns itself.
// It does this by (N - S) where N is the order.  See Coze docs on low-S.
func ToLowS(c *Key, s *big.Int) error {
	lowS, err := IsLowS(c, s)
	if err != nil {
		return err
	}
	if !lowS {
		s.Sub(c.Alg.Curve().EllipticCurve().Params().N, s)
		return nil
	}
	return nil
}

// ECDSAToLowSSig generates low-S signature from existing ecdsa signatures (high
// or low-S).  This is useful for migrating signatures from non-Coze systems
// that may have high S signatures. See Coze docs on low-S.
func ECDSAToLowSSig(c *Key, coze *Coze) (err error) {
	if c.Alg.Genus() != ECDSA {
		return nil
	}
	size := c.Alg.SigAlg().SigSize() / 2
	r := big.NewInt(0).SetBytes(coze.Sig[:size])
	s := big.NewInt(0).SetBytes(coze.Sig[size:])

	// low-S
	err = ToLowS(c, s)
	if err != nil {
		return err
	}
	coze.Sig = PadInts(r, s, c.Alg.SigSize())

	// Make sure the possible mutation of the signature is valid.
	valid, err := c.VerifyCoze(coze)
	if !valid {
		return err
	}

	return nil
}
