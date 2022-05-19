package coze

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	ce "github.com/cyphrme/coze/enum"
)

// Cy is for signed Coze objects.
//
// See the Coze docs (README.md) for more on the construction of `cy`.
//
// This library expects ordered Cy struct fields so sorting after marshaling
// isn't needed.
//
// Fields:
// Cad: head's canon digest.
// Can: Explicit canon of head.  In not present, `can` is assumed implicitly from `head`'s present fields.
// cyd: Cy digest.  Thumbprint of `cy` with canon ["cad","sig"].
// Head: Head serialized for input or after canonicalization/normalization.
// Key: Key used to sign the message. Must be pointer, otherwise json.Marshal will not marshal on zero type. See https://github.com/golang/go/issues/11939.
// Sig: signature over head.
// Sigs: Slice of signatures. (currently not implemented)
//
// SCH: The "Standard Coze Head" fields ["alg","iat","tmb","typ"] in a Coze object. Populated by Coze functions like "SetMeta", and is ignored by marshaling/unmarshaling.
type Cy struct {
	Cad  B64             `json:"cad,omitempty"`
	Can  []string        `json:"can,omitempty"`
	Cyd  B64             `json:"cyd,omitempty"`
	Head json.RawMessage `json:"head"`
	Key  *CozeKey        `json:"key,omitempty"` // Must be pointer for json.Marshal. https://github.com/golang/go/issues/11939
	Sig  B64             `json:"sig,omitempty"`
	Sigs json.RawMessage `json:"sigs,omitempty"`

	Sch Head `json:"-"`
}

// CydCanon is the canon for a `cyd`.
var CydCanon = []string{"head", "sig"}

// CyEn a JSON encapsulator for coze.Cy. It's useful for wrapping struct Cy into
// a JSON `cy`.
type CyEn struct {
	Cy Cy `json:"cy"`
}

// MarshalJSON is a custom marshaler that ensures keys are UTF-8 sorted.
func (cy *Cy) MarshalJSON() (b []byte, err error) {
	// Cast to a proxy type so that default "UnmarshalJSON" will be called.
	type cy2 Cy
	b, err = Marshal((*cy2)(cy))
	if err != nil {
		return nil, err
	}

	b, err = Canonical(b, nil)
	if err != nil {
		return nil, err
	}

	return b, err
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

	// Set Sch from Head ["alg", "iat", etc...]
	err = json.Unmarshal(cy.Head, &cy.Sch)
	if err != nil {
		return err
	}

	if cy.Can != nil {
		sort.Strings(cy.Can) // sorts in place
	} else {
		// Existing fields are implicit canon.
		c, err := CanonB(cy.Head)
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
	cy.Cad, err = CH(cy.Head, cy.Can, cy.Sch.Alg.Hash())
	if err != nil {
		return err
	}

	// Calculate `cyd`
	cy.Cyd = GenCyd(cy.Sch.Alg.Hash(), cy.Cad, cy.Sig)

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

	b, err := Canonical(cy.Head, canon)
	if err != nil {
		return false, err
	}
	return ck.VerifyRaw(b, cy.Sig)
}

// Cyer allows Head types to be converted into a Coze.Cy.  Libraries
// implementing Coze should not define their own "cy" types, but their own
// "head" types.  Those head types may then be converted into Coze.Cy types.
type Cyer interface {
	Cy(sig B64) (*Cy, error)
}

// Method Cy implements the Cyer interface.
func (cy *Cy) Cy(sig B64) (*Cy, error) {
	cy.Sig = sig
	return cy, nil
}

// Verify is a convenience function that is is equivalent to calling
// `Cyer.Cy().Verify(ck, canon)`. See docs on cy.Verify.
func Verify(cy Cyer, ck *CozeKey, sig B64, canon interface{}) (bool, error) {
	c, err := cy.Cy(sig)
	if err != nil {
		return false, err
	}
	return c.Verify(ck, canon)
}
