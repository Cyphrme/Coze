package coze

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cyphrme/coze/enum"
)

// Cy is for signed Coze objects (cozies).  See the Coze docs (README.md) for
// more on the construction of `coze`.
//
// (Struct fields must appear in sorted order for JSON marshaling.)
//
// Fields:
//
// Can: "Canon" Pay's fields in order of appearance.
//
// Cad: "Canonical Digest" Pay's compactified digest.
//
// Cyd: "Cy digest" with canon ["cad","sig"].
//
// Pay: Payload.
//
// Key: Key used to sign the message. Must be pointer, otherwise json.Marshal
// (and by extension coze.Marshal) will not marshal on zero type. See
// https://github.com/golang/go/issues/11939.
//
// Sig: signature over pay.
//
// Parsed: The parsed standard Coze pay fields ["alg","iat","tmb","typ"].
// Populated by Coze functions like "SetMeta", and is ignored by
// marshaling/unmarshaling.
type Cy struct {
	Can []string        `json:"can,omitempty"`
	Cad B64             `json:"cad,omitempty"`
	Cyd B64             `json:"cyd,omitempty"`
	Pay json.RawMessage `json:"pay,omitempty"`
	Key *CozeKey        `json:"key,omitempty"` // Must be pointer for Marshal.
	Sig B64             `json:"sig,omitempty"`

	Parsed Pay `json:"-"`
}

// Coze is a JSON encapsulator for coze.Cy and is used to wrap struct Cy into a
// JSON `coze`.
type Coze struct {
	Cy Cy `json:"coze"`
}

// CydCanon is the canon for a `cyd`.
var CydCanon = []string{"cad", "sig"}

// SetMeta canonicalizes Pay, according to Can if set, and recalculates meta for a given `coze`. Cy.Pay
// and Cy.Sig must be set.  SetMeta does no cryptographic verification.
//
// SeMeta recalculates [can, cad, cyd]
//
// This function is somewhat inefficient and can be written better.
func (cy *Cy) SetMeta() (err error) {
	if cy.Pay == nil {
		return errors.New("coze: cy.head is nil")
	}
	if cy.Sig == nil {
		return errors.New("coze: sig is nil")
	}

	// Set Parsed from Pay.
	err = json.Unmarshal(cy.Pay, &cy.Parsed)
	if err != nil {
		return err
	}

	// TODO fix this/doc it?
	if cy.Can == nil {
		// Existing fields are implicit canon.
		c, err := Canon(cy.Pay)
		if err != nil {
			return err
		}
		cy.Can = c
	} else {
		// Canonicalize `pay` and marshal.
		b, err := Canonical(cy.Pay, cy.Can)
		if err != nil {
			return err
		}
		cy.Pay = b
	}

	// Generate `cad``
	cy.Cad, err = CanonHash(cy.Pay, cy.Can, cy.Parsed.Alg.Hash())
	if err != nil {
		return err
	}

	// Calculate `cyd`
	cy.Cyd = GenCyd(cy.Parsed.Alg.Hash(), cy.Cad, cy.Sig)

	return nil
}

// GenCyd generates and returns `cyd`.
func GenCyd(hash enum.HashAlg, cad B64, sig B64) (cyd B64) {
	var cadSig = []byte(fmt.Sprintf(`{"cad":"%s","sig":"%s"}`, cad, sig))
	return enum.Hash(hash, cadSig)
}

// Verify cryptographically verifies `coze` with given `sig`.  Canon is optional.
func (cy *Cy) Verify(ck *CozeKey, canon interface{}) (bool, error) {
	if cy.Sig == nil {
		return false, errors.New("coze: sig is nil")
	}

	h := new(Pay)
	err := json.Unmarshal(cy.Pay, h)
	if err != nil {
		return false, err
	}
	if !bytes.Equal(h.Tmb, ck.Tmb) {
		return false, errors.New("coze: key tmb and cy tmb do not match")
	}

	b, err := Canonical(cy.Pay, canon)
	if err != nil {
		return false, err
	}
	return ck.VerifyMsg(b, cy.Sig)
}
