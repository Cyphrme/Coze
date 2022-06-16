// See `coze.md` for details.
package coze

import (
	"bytes"
	"encoding/json"
	"math/big"

	"golang.org/x/crypto/sha3"
)

// Pay contains the standard fields in a signed Coze object.
type Pay struct {
	Alg SEAlg  `json:"alg,omitempty"` // e.g. "ES256"
	Iat int64  `json:"iat,omitempty"` // e.g. 1623132000
	Tmb B64    `json:"tmb,omitempty"` // e.g. "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk"
	Typ string `json:"typ,omitempty"` // e.g. "cyphr.me/msg/create"
}

// CozeMarshaler is a UTF-8 marshaler for Go structs.  Go's `json.Marshal`
// removes the valid characters "&". "<", ">".  See note on Marshal.
type CozeMarshaler interface {
	CozeMarshal() ([]byte, error)
}

// Marshal is a UTF-8 friendly marshaler.  Go's json.Marshal is not UTF-8
// friendly because it replaces the valid JSON and valid UTF-8 characters "&".
// "<", ">" with the "slash u" unicode escaped forms (e.g. \u0026).  It
// preemptively escapes for HTML friendliness.  Where JSON may include these
// characters, json.Marshal should not be used. Playground of Go breaking a book
// title: https://play.golang.org/p/o2hiX0c62oN
func Marshal(i any) ([]byte, error) {
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
// level.  Spaces instead of tabs because some applications use 8 spaces per
// tab, which is excessive.
func MarshalPretty(i any) ([]byte, error) {
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

// Hash hashes msg and returns the digest or a set size. Returns nil on error.
//
// This function was written because it doesn't exist in the standard library.
// If in the future there is a standard lib function, use that and deprecate
// this.
//
// Shake128 returns 32 bytes. Shake256 returns 64 bytes.
func Hash(alg HashAlg, msg []byte) (digest B64) {
	if alg == Shake128 {
		h := make([]byte, 32)
		sha3.ShakeSum128(h, msg)
		return h
	}

	if alg == Shake256 {
		h := make([]byte, 64)
		sha3.ShakeSum256(h, msg)
		return h
	}

	hash := alg.goHash()
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

// PadCon creates a big-endian byte slice with given size that is the left
// padded concatenation of two input integers.  Parameter `size` must be even.
// From Go's packages, X, Y, R, and S are type big.Int of varying size. Before
// encoding to fixed sized string, left padding of bytes is needed.
//
// NOTE: EdDSA is little-endian while ECDSA is big-endian.  EdDSA should not be
// used with this function.
//
// For ECDSA, Coze's `x` and `sig` is left padded concatenation of X || Y and R
// || S respectively.
//
// Note: ES512's signature size is 132 bytes (and not 128, 131, or 130.25),
// because R and S are each respectively rounded up and padded to 528 and for a
// total signature size of 1056 bits.
// See https://datatracker.ietf.org/doc/html/rfc4754#section-7
func PadCon(r, s *big.Int, size int) (out B64) {
	if !(size%2 == 0) {
		panic("size must be even.")
	}

	out = make([]byte, size)
	half := size / 2

	rb := r.Bytes()
	copy(out[half-len(rb):], rb)

	sb := s.Bytes()
	copy(out[half+(half-(len(sb))):], sb)

	return out
}
