package enum

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

// CryptoKey is a generalization of a singing or encryption cryptographic key:
// public, private, or a key pair.
type CryptoKey struct {
	Alg     SEAlg
	Public  crypto.PublicKey
	Private crypto.PrivateKey
}

// NewCryptoKey generates a new CryptoKey.
func NewCryptoKey(alg SEAlg) (ck *CryptoKey, err error) {

	var cryptoKey = new(CryptoKey)
	cryptoKey.Alg = alg

	switch SigAlg(alg) {
	case Ed25519, Ed25519ph:
		// The Go Private key is the seed || public key
		cryptoKey.Public, cryptoKey.Private, err = ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}

		return cryptoKey, nil

	case ES224, ES256, ES384, ES512:
		keyPair, err := ecdsa.GenerateKey(alg.Curve().EllipticCurve(), rand.Reader)
		if err != nil {
			return nil, err
		}

		cryptoKey.Public = keyPair.PublicKey
		cryptoKey.Private = keyPair
		return cryptoKey, nil
	default:
		return nil, errors.New("coze.enum.NewCryptoKey: Unknown Alg")
	}
}

// Sign signs a precalculated digest.  On error, returns zero bytes. Digest's
// length must match c.Alg.Hash().Size().
func (c CryptoKey) Sign(digest []byte) (sig []byte, err error) {
	if len(digest) != c.Alg.Hash().Size() {
		return nil, errors.New(fmt.Sprintf("coze.enum: digest length does not match alg.hash.size. Len: %d, Alg: %s.", len(digest), c.Alg.String()))
	}

	switch c.Alg.SigAlg() {
	default:
		return nil, errors.New("coze.enum.SignDigest: Unknown Alg")
	case ES224, ES256, ES384, ES512:

		v, ok := c.Private.(*ecdsa.PrivateKey)
		if !ok {
			return nil, errors.New("Not a valid ECDSA private key.")
		}

		// Note: ECDSA Sig is always R || S of a fixed size with left padding.  For
		// example, ES256 should always have a 64 byte signature.
		r, s, err := ecdsa.Sign(rand.Reader, v, digest)
		if err != nil {
			return nil, err
		}

		return PadCon(r, s, c.Alg.SigAlg().SigSize()), nil

	case Ed25519, Ed25519ph:
		v, ok := c.Private.(ed25519.PrivateKey)
		if !ok {
			return nil, errors.New("Not a valid EdDSA private key")
		}

		return ed25519.Sign(v, digest), nil
	}
}

// Verify verifies that a signature is valid with a given public CryptoKey
// and digest. `digest` should be the digest of the original msg to verify.
func (c CryptoKey) Verify(digest, sig []byte) (valid bool) {
	if len(sig) == 0 || len(digest) == 0 {
		return false
	}

	switch c.Alg.SigAlg() {
	default:
		return false
	case ES224, ES256, ES384, ES512:
		var size = c.Alg.SigAlg().SigSize() / 2
		r := big.NewInt(0).SetBytes(sig[:size])
		s := big.NewInt(0).SetBytes(sig[size:])

		v, ok := c.Public.(ecdsa.PublicKey)
		if !ok {
			return false
		}

		return ecdsa.Verify(&v, digest, r, s)
	case Ed25519, Ed25519ph:
		v, ok := c.Public.(ed25519.PublicKey)
		if !ok {
			return false
		}

		return ed25519.Verify(v, digest, sig)
	}
}

// SignMsg signs a pre-hash msg.  On error, returns zero bytes.
func (c CryptoKey) SignMsg(msg []byte) (sig []byte, err error) {
	return c.Sign(Hash(c.Alg.Hash(), msg))
}

// Verify verifies that a signature with a given public CryptoKey and
// signed message.
func (c CryptoKey) VerifyMsg(msg, sig []byte) (valid bool) {
	return c.Verify(Hash(c.Alg.Hash(), msg), sig)
}
