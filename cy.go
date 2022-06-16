package coze

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Cy is for signed Coze objects (cozies).  See the Coze docs (README.md) for
// more on the construction of `coze`.
//
// Struct fields appear in sorted order for JSON marshaling.
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
// Populated by Coze functions like "Meta", and is ignored by
// marshaling/unmarshaling.
type Cy struct {
	Can []string        `json:"can,omitempty"`
	Cad B64             `json:"cad,omitempty"`
	Cyd B64             `json:"cyd,omitempty"`
	Pay json.RawMessage `json:"pay,omitempty"`
	Key *CozeKey        `json:"key,omitempty"`
	Sig B64             `json:"sig,omitempty"`

	Parsed Pay `json:"-"`
}

// String implements fmt.Stringer because otherwise `pay` prints as bytes.
func (cy *Cy) String() string {

	b, err := MarshalPretty(cy)
	if err != nil {
		fmt.Println(err)
	}
	return string(b) + "\n"
}

// Coze is a JSON encapsulator for coze.Cy and is used to wrap struct Cy into a
// JSON `coze`.
type Coze struct {
	Cy Cy `json:"coze"`
}

// CydCanon is the canon for a `cyd`.
var CydCanon = []string{"cad", "sig"}

// Meta recalculates meta, [can, cad, cyd], for a given `coze`. Cy.Pay,
// Cy.Pay.Alg, and Cy.Sig must be set.  Meta does no cryptographic
// verification.
func (cy *Cy) Meta() (err error) {
	if cy.Pay == nil {
		return errors.New("coze: cy.Pay is nil")
	}
	if cy.Sig == nil {
		return errors.New("coze: sig is nil")
	}

	// Set Parsed from Pay.
	err = json.Unmarshal(cy.Pay, &cy.Parsed)
	if err != nil {
		return err
	}

	c, err := Canon(cy.Pay)
	if err != nil {
		return err
	}
	cy.Can = c

	canonical, err := Canonical(cy.Pay, nil)
	if err != nil {
		return err
	}
	cy.Cad = Hash(cy.Parsed.Alg.Hash(), canonical)
	cy.Cyd = GenCyd(cy.Parsed.Alg.Hash(), cy.Cad, cy.Sig)

	return nil
}

// GenCyd generates and returns `cyd`.
func GenCyd(hash HashAlg, cad B64, sig B64) (cyd B64) {
	var cadSig = []byte(fmt.Sprintf(`{"cad":"%s","sig":"%s"}`, cad, sig))
	return Hash(hash, cadSig)
}
