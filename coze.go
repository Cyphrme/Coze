// See `coze.md` for details.
package coze

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"golang.org/x/crypto/sha3"
)

// Pay contains the standard Coze pay fields as well as a custom struct given by
// third party applications.  This allows embedding third party structs into Pay
// for creating custom cozies (see example ExampleCozeKey_SignPay).
//
// Note: The custom MarshalJSON() renders the JSON tags on [Alg, Iat, Tmb, Typ]
// ineffective, however they are present for documentation.  Field Struct will
// be marshaled on not empty, tag `json:"-"` is meaningless.
type Pay struct {
	Alg SEAlg  `json:"alg,omitempty"` // e.g. "ES256"
	Iat int64  `json:"iat,omitempty"` // e.g. 1623132000
	Tmb B64    `json:"tmb,omitempty"` // e.g. "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk"
	Typ string `json:"typ,omitempty"` // e.g. "cyphr.me/msg/create"

	Struct interface{} `json:"-"` // Custom arbitrary struct given by application. Custom marshaler promotes given structure and tag `json:"-"` is meaningless.
}

// Pay.Coze() returns a new Coze with only Pay populated.
func (p *Pay) Coze() (coze *Coze, err error) {
	coze = new(Coze)
	coze.Pay, err = Marshal(p)
	return coze, err
}

// String implements fmt.Stringer.
func (p Pay) String() string {
	b, err := p.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}

// Marshal promotes the embedded field "Struct" to top level JSON. Solution
// taken from Jonathan Hall
// https://jhall.io/posts/go-json-tricks-embedded-marshaler
func (p *Pay) MarshalJSON() ([]byte, error) {
	type pay2 Pay // Break infinite Marshal loop

	pay, err := Marshal((*pay2)(p))
	if err != nil {
		return nil, err
	}
	if p.Struct == nil {
		return pay, nil
	}

	s, err := json.Marshal(p.Struct)
	if err != nil {
		return nil, err
	}
	s[0] = ','
	return append(pay[:len(pay)-1], s...), nil
}

// Coze is for signed Coze objects (cozies).  See the Coze docs (README.md) for
// more on the construction of `coze`.
//
// Struct fields appear in sorted order for JSON marshaling.
//
// Fields:
//
// Can: "Canon" Pay's fields in order of appearance.
//
// Cad: "Canonical Digest" Pay's compactified form digest.
//
// Czd: "Coze digest" with canon ["cad","sig"].
//
// Pay: Payload.
//
// Key: Key used to sign the message. Must be pointer, otherwise json.Marshal
// (and by extension coze.Marshal) will not marshal on zero type. See
// https://github.com/golang/go/issues/11939.
//
// Sig: Signature over pay.
//
// Parsed: The parsed standard Coze pay fields ["alg","iat","tmb","typ"] from
// Pay.  Parsed is populated by Meta() and is JSON ignored. The source of truth
// is Pay, not Parsed; do not use Parsed until after calling Meta() which
// populates Parsed from Pay.
type Coze struct {
	Can []string        `json:"can,omitempty"`
	Cad B64             `json:"cad,omitempty"`
	Czd B64             `json:"czd,omitempty"`
	Pay json.RawMessage `json:"pay,omitempty"`
	Key *CozeKey        `json:"key,omitempty"`
	Sig B64             `json:"sig,omitempty"`

	Parsed *Pay `json:"-"`
}

// String implements fmt.Stringer.  Without this method `pay` prints as bytes.
// Errors are returned as a string.
func (cz Coze) String() string {
	b, err := Marshal(cz)
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}

// CzdCanon is the canon for a `czd`.
var CzdCanon = []string{"cad", "sig"}

// Meta recalculates meta, [can, cad, czd], for a given `coze`. Coze.Pay,
// Coze.Pay.Alg, and Coze.Sig must be set. Meta does no cryptographic
// verification.
func (cz *Coze) Meta() (err error) {
	if cz.Pay == nil || cz.Sig == nil {
		return errors.New("Meta: pay and/or sig is nil")
	}

	return cz.MetaWithAlg(0)
}

// MetaWithAlg is for contextual cozies that may be lacking `alg` in `pay`, but
// `alg` in otherwise known.  Errors if pay.alg is set and doesn't match
// parameter alg.  See notes on Meta()
func (cz *Coze) MetaWithAlg(alg SEAlg) (err error) {
	// Set Parsed from Pay.
	err = json.Unmarshal(cz.Pay, &cz.Parsed)
	if err != nil {
		return err
	}

	if alg == 0 {
		alg = cz.Parsed.Alg
	}
	if alg != cz.Parsed.Alg {
		return fmt.Errorf("MetaWithAlg: input alg %s doesn't match pay.alg %s", alg, cz.Parsed.Alg)
	}

	// Use Canonical to Compactify pay
	canonical, err := Canonical(cz.Pay, nil)
	if err != nil {
		return err
	}
	c, err := GetCanon(canonical)
	if err != nil {
		return err
	}
	cz.Can = c
	cz.Cad = Hash(cz.Parsed.Alg.Hash(), canonical)
	cz.Czd = GenCzd(cz.Parsed.Alg.Hash(), cz.Cad, cz.Sig)
	return nil
}

// GenCzd generates and returns `czd`.
func GenCzd(hash HashAlg, cad B64, sig B64) (czd B64) {
	var cadSig = []byte(fmt.Sprintf(`{"cad":"%s","sig":"%s"}`, cad, sig))
	return Hash(hash, cadSig)
}

// CozeMarshaler is a UTF-8 marshaler for Go structs. Go's `json.Marshal`
// removes the valid characters "&". "<", ">". See note on Marshal.
type CozeMarshaler interface {
	CozeMarshal() ([]byte, error)
}

// Marshal is a UTF-8 friendly marshaler. Go's json.Marshal is not UTF-8
// friendly because it replaces the valid JSON and valid UTF-8 characters "&".
// "<", ">" with the "slash u" unicode escaped forms (e.g. \u0026). It
// preemptively escapes for HTML friendliness. Where JSON may include these
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
// level. Spaces instead of tabs because some applications use 8 spaces per
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

// Hash hashes msg and returns the digest or a set size. Returns nil on error or invalid HashAlg.
//
// This function was written because it doesn't exist in the standard library.
// If in the future there is a standard lib function, use that and deprecate
// this.
//
// SHAKE128 returns 32 bytes. SHAKE256 returns 64 bytes.
func Hash(alg HashAlg, msg []byte) (digest B64) {
	if alg == SHAKE128 {
		h := make([]byte, 32)
		sha3.ShakeSum128(h, msg)
		return h
	}

	if alg == SHAKE256 {
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

	return hash.Sum(nil)
}

// PadInts creates a big-endian byte slice with given size that is the left
// padded concatenation of two input integers. Parameter `size` must be even.
// From Go's packages, X, Y, R, and S are type big.Int of varying size. Before
// encoding to fixed sized string, left padding of bytes is needed.
//
// NOTE: EdDSA is little-endian while ECDSA is big-endian. EdDSA should not be
// used with this function.
//
// For ECDSA, Coze's `x` and `sig` is left padded concatenation of X || Y and R
// || S respectively.
//
// Note: ES512's signature size is 132 bytes (and not 128, 131, or 130.25),
// because R and S are each respectively rounded up and padded to 528 and for a
// total signature size of 1056 bits.
// See https://datatracker.ietf.org/doc/html/rfc4754#section-7
func PadInts(r, s *big.Int, size int) (out B64) {
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
