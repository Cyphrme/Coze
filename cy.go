package coze

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	ce "github.com/cyphrme/coze/enum"
)

// Cy is for signed Coze objects.  See the Coze docs (README.md) for more on the
// construction of `cy`.  Fields must appear in correct order for JSON
// marshaling.
//
// Fields: Cad: head's canon digest.
//
// Can: Explicit canon of head.  In not present, `can` is assumed implicitly
// from `head`'s present fields.
//
// Cyd: Cy digest.  Thumbprint of `cy` with canon ["cad","sig"].
//
// Head: Head serialized for input or after canonicalization/normalization.
//
// Key: Key used to sign the message. Must be pointer, otherwise json.Marshal
// (and by extension coze.Marshal) will not marshal on zero type. See
// https://github.com/golang/go/issues/11939.
//
// Sig: signature over head.
//
// Sigs: Slice of signatures. (currently not implemented)
//
// Parsed: The parsed standard Coze head fields ["alg","iat","tmb","typ"].
// Populated by Coze functions like "SetMeta", and is ignored by
// marshaling/unmarshaling.
type Cy struct {
	Cad  B64             `json:"cad,omitempty"`
	Can  []string        `json:"can,omitempty"`
	Cyd  B64             `json:"cyd,omitempty"`
	Head json.RawMessage `json:"head"`
	Key  *CozeKey        `json:"key,omitempty"` // Must be pointer for Marshal.
	Sig  B64             `json:"sig,omitempty"`
	Sigs json.RawMessage `json:"sigs,omitempty"`

	Parsed Head `json:"-"`
}

// CydCanon is the canon for a `cyd`.
var CydCanon = []string{"head", "sig"}

// CyEn a JSON encapsulator for coze.Cy. It's useful for wrapping struct Cy into
// a JSON `cy`.
type CyEn struct {
	Cy Cy `json:"cy"`
}

// SetMeta canonicalizes Head and recalculates Meta for a given `cy`.
// Cy.Head and Cy.Sig must be set.
//
// SetMeta does not verify cy.
//
// SeMeta recalculates [can, cad, cyd]
//
// This function is somewhat inefficient and can be written better.
func (cy *Cy) SetMeta() (err error) {
	if cy.Head == nil {
		return errors.New("coze: cy.head is nil")
	}
	if cy.Sig == nil {
		return errors.New("coze: sig is nil")
	}

	// Set Parsed from Head.
	err = json.Unmarshal(cy.Head, &cy.Parsed)
	if err != nil {
		return err
	}

	if cy.Can != nil {
		sort.Strings(cy.Can) // sorts in place
	} else {
		// Existing fields are implicit canon.
		c, err := Canon(cy.Head)
		if err != nil {
			return err
		}
		cy.Can = c
	}

	// Canonicalize `head` and marshal.
	b, err := Canonical(cy.Head, cy.Can)
	if err != nil {
		return err
	}
	cy.Head = b

	// Generate `cad``
	cy.Cad, err = CanonHash(cy.Head, cy.Can, cy.Parsed.Alg.Hash())
	if err != nil {
		return err
	}

	// Calculate `cyd`
	cy.Cyd = GenCyd(cy.Parsed.Alg.Hash(), cy.Cad, cy.Sig)

	return nil
}

// GenCyd generates and returns `cyd`.
func GenCyd(hash ce.HashAlg, cad B64, sig B64) (cyd B64) {
	var cadSig = []byte(fmt.Sprintf(`{"cad":"%s","sig":"%s"}`, cad, sig))
	return ce.Hash(hash, cadSig)
}

// Verify cryptographically verifies `cy` with given `sig`.  Canon is optional.
func (cy *Cy) Verify(ck *CozeKey, canon interface{}) (bool, error) {
	if cy.Sig == nil {
		return false, errors.New("coze: sig is nil")
	}

	h := new(Head)
	err := json.Unmarshal(cy.Head, h)
	if err != nil {
		return false, err
	}
	if !bytes.Equal(h.Tmb, ck.Tmb) {
		return false, errors.New("coze: key tmb and cy tmb do not match")
	}

	// Canonical removes spaces or canonicalizes
	b, err := Canonical(cy.Head, canon)
	if err != nil {
		return false, err
	}
	return ck.VerifyRaw(b, cy.Sig)
}
