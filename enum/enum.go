package enum

import (
	"bytes"
	"encoding/json"
	"math/big"

	"golang.org/x/crypto/sha3"
)

// Hash hashes msg and returns the digest or a set size. Returns nil on error.
//
// This function was written because it doesn't exist in the standard lib.  If
// in the future there is a standard lib function, use that and deprecate this.
//
// Shake128 returns 32 bytes. Shake256 returns 64 bytes.
func Hash(c HashAlg, msg []byte) (digest []byte) {
	// If HashAlg is zero type, return nil.
	if c == 0 {
		return nil
	}

	if c == Shake128 {
		h := make([]byte, 32)
		sha3.ShakeSum128(h, msg)
		return h
	}

	if c == Shake256 {
		h := make([]byte, 64)
		sha3.ShakeSum256(h, msg)
		return h
	}

	hash := c.goHash()
	if hash == nil {
		return nil
	}
	_, err := hash.Write(msg)
	if err != nil {
		return nil
	}

	digest = hash.Sum(nil)
	return
}

// PadSig creates a signature that is the padded concatendation of R and S.
// Padding is needed because R and S are integers that can vary in size.
//
// For ECDSA sig is always R || S of a fixed size with left
// padding.  For example, ES256 should always have a 64 byte signature. [0,0, 1
// .... || 0,0,1 ...]
//
// Note: ES512's signature size is 132 bytes (and not 128, 131, or 130.25),
// because R and S are each respectively rounded up and padded to 528 and for a
// total signature size of 1056 bits.
// See https://datatracker.ietf.org/doc/html/rfc4754#section-7
func PadSig(r, s *big.Int, size int) (sig []byte) {
	sig = make([]byte, size)
	half := size / 2

	// R
	rb := r.Bytes() // get r bytes
	n := copy(sig[half-len(rb):], rb)
	if n != len(rb) {
		return nil
	}

	// S
	sb := s.Bytes() // get s bytes
	n = copy(sig[half+(half-(len(sb))):], sb)
	if n != len(sb) {
		return nil
	}
	return sig
}

// Marshal is a UTF-8 friendly marshaler.  Go's json.Marshal is not UTF-8
// friendly because it replaces the valid UTF-8 and JSON characters "&". "<",
// ">" with the "slash u" unicode escaped forms (e.g. \u0026).  It preemptively
// escapes for HTML friendliness.  Where text may include any of these
// characters, json.Marshal should not be used. Playground of Go breaking a
// book title: https://play.golang.org/p/o2hiX0c62oN.  Taken from package `coze`.
func Marshal(i interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(i)
	if err != nil {
		return nil, err
	}
	return bytes.TrimRight(buffer.Bytes(), "\n"), nil
}

// MarshalPretty is the pretty version of Marshal. It uses 4 spaces for each
// level.  Spaces instead of tabs because many applications still use 8 spaces
// per tab, which is excessive. Taken from package `coze`.
func MarshalPretty(i interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "    ")
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(i)
	if err != nil {
		return nil, err
	}
	return bytes.TrimRight(buffer.Bytes(), "\n"), nil
}
