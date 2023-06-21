package coze

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"golang.org/x/crypto/sha3"
)

// Coze is for signed Coze objects (cozies).  See the Coze docs (README.md) for
// more on the construction of `coze`.
//
//	Pay: The raw Payload.
//	Key: Key used to sign the message. Must be pointer, otherwise
//	  json.Marshal will not marshal on zero type.
//	  See: https://github.com/golang/go/issues/11939.
//	Can: "Canon" Pay's fields in order of appearance.
//	Cad: "Canonical Digest" Pay's compactified form digest.
//	Sig: Signature over `cad`.
//	Czd: "Coze digest" with canon ["cad","sig"].
//	Parsed: The standard Coze pay fields ["alg","iat","tmb","typ"] parsed
//	  from `Pay`.  `Parsed` is populated by Meta() and is JSON ignored.
type Coze struct {
	Pay json.RawMessage `json:"pay,omitempty"`
	Key *Key            `json:"key,omitempty"`
	Can []string        `json:"can,omitempty"`
	Cad B64             `json:"cad,omitempty"`
	Sig B64             `json:"sig,omitempty"`
	Czd B64             `json:"czd,omitempty"`

	Parsed *Pay `json:"-"`
}

// String implements fmt.Stringer.  Without this method `pay` prints as bytes.
// On error, returns the error as a string.
func (cz Coze) String() string {
	b, err := Marshal(cz)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

// Meta calculates [can, cad, czd] and sets Coze.Parsed
// ["alg","iat","tmb","typ"] from Pay. Coze.Pay, Coze.Pay.Alg, and Coze.Sig must
// be set.  Meta resets Parsed ("alg","iat","tmb","typ") to zero before
// populating Parsed from Pay. If needing to use for contextual cozies, use
// "MetaWithAlg".
//
// Meta does no cryptographic verification.
func (cz *Coze) Meta() (err error) {
	if cz.Pay == nil || cz.Sig == nil {
		return errors.New("Meta: pay and/or sig is nil")
	}
	// Reset coze.parsed to zero.
	cz.Parsed = new(Pay)
	return cz.MetaWithAlg("")
}

// MetaWithAlg is for contextual cozies that may be lacking `alg` in `pay`, but
// `alg` in otherwise known.  MetaWithAlg recalculates [can, cad, czd] and sets
// Coze.Parsed ("alg","iat","tmb","typ") from Pay.  Will not calculated `czd`
// if Coze.Sig is empty.
//
// Errors on
// 1. Invalid JSON.
// 2. No alg is given. (both coze.pay.alg and alg are empty).
// 3. Pay.Alg doesn't match parameter alg if both are set.
//
// MetaWithAlg does no cryptographic verification.
func (cz *Coze) MetaWithAlg(alg SEAlg) (err error) {
	// Set Parsed from Pay.
	err = json.Unmarshal(cz.Pay, &cz.Parsed)
	if err != nil {
		return err
	}
	if alg == "" {
		alg = cz.Parsed.Alg
	}
	if cz.Parsed.Alg != "" && alg != cz.Parsed.Alg {
		return fmt.Errorf("MetaWithAlg: input alg %q and pay.alg %q are unequal", alg, cz.Parsed.Alg)
	}
	b, err := compact(cz.Pay)
	if err != nil {
		return err
	}
	cz.Cad, err = Hash(alg.Hash(), b)
	if err != nil {
		return err
	}
	cz.Can, err = Canon(b)
	if err != nil {
		return err
	}
	cz.Czd = []byte{} // Zero `czd` in case Coze.Sig is not set.
	if len(cz.Sig) != 0 {
		cz.Czd, err = GenCzd(alg.Hash(), cz.Cad, cz.Sig)
		return err
	}
	return nil
}

// UnmarshalJSON unmarshals checks for duplicates and unmarshals `coze`.
// See notes on Pay.UnmarshalJSON.
func (cz *Coze) UnmarshalJSON(b []byte) error {
	err := checkDuplicate(json.NewDecoder(bytes.NewReader(b)))
	if err != nil {
		return err
	}

	type coze2 Coze // Break infinite unmarshal loop
	cz2 := new(coze2)
	cz2.Parsed = cz.Parsed
	cz2.Key = cz.Key
	err = json.Unmarshal(b, cz2)
	if err != nil {
		return err
	}
	*cz = *(*Coze)(cz2)
	return nil
}

// CzdCanon is the canon for a `czd`.
var CzdCanon = []string{"cad", "sig"}

const max_safe_integer = 9007199254740991

// GenCzd generates and returns `czd`.
func GenCzd(hash HshAlg, cad B64, sig B64) (czd B64, err error) {
	return Hash(hash, []byte(fmt.Sprintf(`{"cad":%q,"sig":%q}`, cad, sig)))
}

// Pay contains the standard Coze pay fields as well as custom struct given by
// third party applications.  This allows embedding third party structs into Pay
// for creating custom cozies (see example ExampleKey_SignPay).
//
// The JSON tags on [Alg, Iat, Tmb, Typ, Rvk, Struct] are ineffective due to the
// custom MarshalJSON(), however they are present for documentation.
//
// `Struct` will be marshaled when not empty. The custom marshaler promotes
// fields inside `Struct` to be top level fields inside of `pay`. The tag
// `json:"-"` is ignored by the custom marshaler, and  is set to "-" so that the
// default marshaler does not include it.
//
// iat and rvk are type int64 and not uint64 to follow the advised type for
// third party time fields.
type Pay struct {
	Alg SEAlg  `json:"alg,omitempty"` // e.g. "ES256"
	Iat int64  `json:"iat,omitempty"` // e.g. 1623132000
	Tmb B64    `json:"tmb,omitempty"` // e.g. "cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk"
	Typ string `json:"typ,omitempty"` // e.g. "cyphr.me/msg/create"

	// Rvk is only for revoke messages.
	Rvk int64 `json:"rvk,omitempty"` // e.g. 1623132000

	// Custom arbitrary struct given by application.
	Struct any `json:"-"`
}

// Pay.Coze() returns a new Coze with only Pay populated.
func (p *Pay) Coze() (coze *Coze, err error) {
	coze = new(Coze)
	coze.Pay, err = Marshal(p)
	return coze, err
}

// String implements fmt.Stringer.
// On error, returns the error as a string.
func (p Pay) String() string {
	b, err := p.MarshalJSON()
	if err != nil {
		return err.Error()
	}
	return string(b)
}

// MarshalJSON promotes the embedded field "Struct" to top level JSON.
// Solution from Jonathan Hall:
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
	// Concatenate the two:
	s[0] = ','
	return append(pay[:len(pay)-1], s...), nil
}

// UnmarshalJSON unmarshals both Pay and if given custom Pay.Struct.
//
// UnmarshalJSON handles deduplicate and throws an error on duplicate.
// See the Go issue: https://github.com/golang/go/issues/48298
func (p *Pay) UnmarshalJSON(b []byte) error {
	err := checkDuplicate(json.NewDecoder(bytes.NewReader(b)))
	if err != nil {
		return err
	}

	type pay2 Pay // Break infinite unmarshal loop
	p2 := new(pay2)
	err = json.Unmarshal(b, p2)
	if err != nil {
		return err
	}
	if p.Struct != nil { // Inner custom application struct.
		str := p.Struct
		err = json.Unmarshal(b, str)
		if err != nil {
			return err
		}
		p2.Struct = str
	}

	if p2.Iat > max_safe_integer || p2.Rvk > max_safe_integer || p2.Iat < 0 || p2.Rvk < 0 {
		return fmt.Errorf("Pay.UnmarshalJSON: values for iat and rvk must be between 0 and 2^53 - 1")
	}

	*p = *(*Pay)(p2)
	return nil
}

// PadInts creates a big-endian byte slice with given size that is the left
// padded concatenation of two input integers. Parameter `size` must be even.
// From Go's packages, X, Y, R, and S are type big.Int of varying size. Before
// encoding to fixed sized string, left padding of bytes is needed.
//
// Algorithm notes: EdDSA is little-endian while ECDSA is big-endian. EdDSA
// should not be used with this function.
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

// Marshaler is a UTF-8 marshaler for Go structs. Go's `json.Marshal`
// removes the valid characters "&". "<", ">". See note on Marshal.
type Marshaler interface {
	CozeMarshal() ([]byte, error)
}

// Marshal is a JSON friendly marshaler.  json.Marshal preemptively replaces the
// valid JSON and UTF-8 characters "&". "<", ">" with the "slash u" unicode
// escapes (e.g. \u0026) in the name of HTML friendliness, for example:
// https://play.golang.org/p/o2hiX0c62oN. As JSON is not HTML, preemptive HTML
// escaping is incorrect.  (The JSON spec calls for no such measures and other
// industry encoders do no such preemptive escaping.)  Where JSON may include
// these legitimate characters, like sanitized arbitrary user data, json.Marshal
// should not be used. The Go team is aware that the existing implementation has
// this and other concerns but has not yet release a new version.  Joe Tsai is
// working on fixes in a yet-to-be-released Tailscale's "JSONv2" package, which
// we hope to use upon release:
// https://pkg.go.dev/github.com/go-json-experiment/json  The package also
// enumerates other JSON best practices that may be of concern for anyone
// concerned with the aforementioned issue.
// https://github.com/go-json-experiment/json#behavior-changes
//
// See https://github.com/Cyphrme/Coze/issues/15 for other JSON encoding concerns.
//
// Go structs already require unique fields, so unlike coze.UnmarshalJSON or
// pay.UnmarshalJSON, marshaling will not sanitize for duplicates.
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

// MarshalPretty uses 4 spaces for each level. Spaces are used instead of tabs
// because some applications display tabs as 8 spaces, which is excessive.
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

// Hash hashes msg and returns the digest. Returns nil on error. Errors on
// invalid HshAlg or if digest is empty.
//
// For algorithms that support arbitrary sized digests, Hash only returns a
// static size.  SHAKE128 returns 32 bytes. SHAKE256 returns 64 bytes.
func Hash(h HshAlg, msg []byte) (digest B64, err error) {
	switch h {
	case SHAKE128:
		digest = make([]byte, 32)
		sha3.ShakeSum128(digest, msg)
	case SHAKE256:
		digest = make([]byte, 64)
		sha3.ShakeSum256(digest, msg)
	default:
		hash := h.goHash()
		if hash == nil {
			return nil, fmt.Errorf("Hash: invalid HashAlg %q", h)
		}
		_, err = hash.Write(msg)
		if err != nil {
			return nil, err
		}
		digest = hash.Sum(nil)
	}

	if len(digest) == 0 { // sanity check
		return nil, fmt.Errorf("Hash: digest is empty; given HashAlg %q", h)
	}
	return digest, nil
}

// IsRevoke returns true if the given Key is marked as revoked.
func (p *Pay) IsRevoke() bool {
	return isRevoke(p.Rvk)
}

func isRevoke(rvk int64) bool {
	// rvk is not allowed to be larger than 2^53 -1.  This library assumes that
	// Unmarshal will error on rvk's that do not meet the specification
	// requirements, so no error is needed here.
	if rvk > max_safe_integer {
		return false
	}
	return rvk > 0
}

// checkDuplicate checks if the JSON string has a duplicate. Go has an issue
// regarding duplicates: https://github.com/golang/go/issues/48298. Another
// solution is being created by Joe Tsai (see notes on Marshal).  When he's
// done, we'll take this out and use v2.
//
// Inspire by Cerise Limón.
// https://stackoverflow.com/questions/50107569/detect-duplicate-in-json-string-golang
func checkDuplicate(d *json.Decoder) error {
	t, err := d.Token()
	if err != nil {
		return err
	}

	// Is it a delimiter?
	delim, ok := t.(json.Delim)
	if !ok {
		return nil // scaler type, nothing to do
	}

	switch delim {
	case '{':
		keys := make(map[string]bool)
		for d.More() {
			// Get field key.
			t, err := d.Token()
			if err != nil {
				return err
			}

			key := t.(string)
			if keys[key] { // Check for duplicates.
				return ErrJSONDuplicate
			}
			keys[key] = true

			// Recursive, Check value in case value is object.
			err = checkDuplicate(d)
			if err != nil {
				return err
			}
		}
		// consume trailing }
		if _, err := d.Token(); err != nil {
			return err
		}

	case '[':
		for d.More() {
			if err := checkDuplicate(d); err != nil {
				return err
			}
		}
		// consume trailing ]
		if _, err := d.Token(); err != nil {
			return err
		}
	}
	return nil
}

// ErrJSONDuplicate is for applications that need to check for the JSON
// duplicate error.  Alternatively, applications need to check for the error
// string, which may change.
var ErrJSONDuplicate = errors.New("Coze: JSON duplicate field name")
