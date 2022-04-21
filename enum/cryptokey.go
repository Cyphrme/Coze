package enum

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"math/big"
)

// CryptoKey is a generalization of a singing or encryption cryptographic key:
// public, private, or a key pair.
type CryptoKey struct {
	Alg     SEAlg
	Public  *crypto.PublicKey
	Private *crypto.PrivateKey
}

// NewCryptoKey generates a new cryptokey.
func NewCryptoKey(alg SEAlg) (ck *CryptoKey, err error) {
	var cryptoKey CryptoKey
	cryptoKey.Alg = alg

	if SigAlg(alg) == Ed25519 {
		var pub, pri []byte
		pub, pri, err = ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		var public crypto.PublicKey = pub
		var private crypto.PrivateKey = pri

		cryptoKey.Public = &public
		cryptoKey.Private = &private
		return &cryptoKey, nil
	}

	var keyPair *ecdsa.PrivateKey
	switch SigAlg(alg) {
	case ES224:
		keyPair, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case ES256:
		keyPair, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case ES384:
		keyPair, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case ES512: // ES512 and **NOT "ES521"** The curve != the alg.
		keyPair, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		return nil, errors.New("cozeenum.NewCryptoKey: Unknown Alg")
	}

	if err != nil {
		return nil, err
	}

	var public crypto.PublicKey = keyPair.PublicKey
	var private crypto.PrivateKey = keyPair

	cryptoKey.Public = &public
	cryptoKey.Private = &private
	return &cryptoKey, nil
}

// Sign signs a msg.  `msg` should be pre-hash.  On error, returns zero bytes.
func (c CryptoKey) SignRaw(msg []byte) (sig []byte, err error) {
	return c.Sign(Hash(c.Alg.Hash(), msg))
}

// SignDigest signs a precalculated digest.  On error, returns zero bytes.
// Digest's size must match c.Alg.Hash().Size().
func (c CryptoKey) Sign(digest []byte) (sig []byte, err error) {
	switch c.Alg.SigAlg() {
	default:
		return nil, errors.New("cozeenum.SignDigest: Unknown Alg")
	case ES224, ES256, ES384, ES512:
		if len(digest) != c.Alg.Hash().Size() {
			return nil, errors.New("cryptokey: digest size does not match cryptokey.alg.hash size")
		}

		priv := *c.Private
		v, ok := priv.(ecdsa.PrivateKey)
		if !ok { // check to see if inner type is a pointer.
			var vv *ecdsa.PrivateKey
			vv, ok = priv.(*ecdsa.PrivateKey)
			v = *vv
		}
		if !ok {
			return nil, errors.New("Not a valid ecdsa private key")
		}

		// Note: ECDSA Sig is always R || S of a fixed size with left padding.  For
		// example, ES256 should always have a 64 byte signature.
		r, s, err := ecdsa.Sign(rand.Reader, &v, digest)
		if err != nil {
			return nil, err
		}

		sig = PadSig(r, s, c.Alg.SigAlg().SigSize())
	}

	return sig, nil
}

// Verify verifies that a signature with a given public cryptokey and
// signed message.
func (c CryptoKey) Verify(msg, sig []byte) (valid bool, err error) {
	digest := Hash(c.Alg.Hash(), msg)
	return c.VerifyDigest(digest, sig)
}

// VerifyDigest verifies that a signature is valid with a given public cryptokey
// and digest. `digest` should be the digest of the original msg to verify.
func (c CryptoKey) VerifyDigest(digest, sig []byte) (valid bool, err error) {
	var size = c.Alg.SigAlg().SigSize() / 2

	r := big.NewInt(0).SetBytes(sig[:size])
	s := big.NewInt(0).SetBytes(sig[size:])

	pub := *c.Public
	v, ok := pub.(ecdsa.PublicKey)
	if !ok {
		return false, errors.New("cryptokey: public key is invalid")
	}

	if !ecdsa.Verify(&v, digest, r, s) {
		return false, nil
	}
	return true, nil
}
