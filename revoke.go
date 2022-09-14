// See README.md https://github.com/Cyphrme/Coze#readme.
package coze

import (
	"encoding/json"
	"errors"
	"time"
)

// Revoke is a self revoke coze which contains the components necessary for self
// revoking a Coze key.  See the documentation section "Revoke".
type Revoke struct {
	Rvk int64  `json:"rvk"`           // Timestamp when key revoke occurred.
	Msg string `json:"msg,omitempty"` // Optional message describing why the key was revoked.
	Pay
}

// revoke is used for marshaling and unmarshalling Revoke.  Since Pay is
// embedded in Revoke, calling json.Marshal/json.Unmarshal will call
// Revoke.Pay.Marshal/Revoke.Pay.UnmarshalJSON, which does not include Rvk or
// Msg. To fix this, for marshaling: Revoke's unique fields, `rvk` and `msg`,
// are marshaled then Rvk.Pay is marshaled, and the two are concatenated. For
// unmarshalling, Revoke's unique fields, `rvk` and `msg`, are unmarshalled,
// Rvk.Pay is unmarshalled, and then revoke is set with the unmarshalled values.
//
// For comparison, see the notes on `func (p *Pay) MarshalJSON()`.
type revoke struct {
	Rvk int64  `json:"rvk"`
	Msg string `json:"msg,omitempty"`
}

// String implements Stringer. Returns empty on error.
func (r Revoke) String() string {
	b, err := Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}

// MarshalJSON marshals Revoke.
func (r *Revoke) MarshalJSON() ([]byte, error) {
	r2 := new(revoke)
	r2.Rvk = r.Rvk
	r2.Msg = r.Msg

	revoke, err := Marshal(r2)
	if err != nil {
		return nil, err
	}

	p, err := json.Marshal(r.Pay)
	if err != nil {
		return nil, err
	}
	// Concatenate Revoke and Pay:
	p[0] = ','
	return append(revoke[:len(revoke)-1], p...), nil
}

// UnmarshalJSON unmarshals a Revoke. See notes on revoke.
func (r *Revoke) UnmarshalJSON(b []byte) error {
	r2 := new(revoke)

	err := json.Unmarshal(b, r2)
	if err != nil {
		return err
	}
	r.Rvk = r2.Rvk
	r.Msg = r2.Msg
	return json.Unmarshal(b, &r.Pay)
}

// Revoke will return a signed revoke coze for the given key as well as setting
// `rvk` on the Coze key itself.
//
// The revoke coze may simply be discarded if not needed, e.g. `_, err :=
// cozeKey.Revoke(nil)`
func (c *Key) Revoke(msg string) (coze *Coze, err error) {
	correct, err := c.Correct()
	if !correct || err != nil {
		return nil, errors.New("revoke: Coze key is not correct")
	}

	now := time.Now().Unix()
	r := new(Revoke)
	r.Alg = c.Alg
	r.Iat = now
	r.Msg = msg
	r.Rvk = now
	r.Tmb = c.Tmb
	r.Typ = "cyphr.me/key/revoke"

	coze = new(Coze)
	coze.Pay, err = r.MarshalJSON()
	if err != nil {
		return
	}

	err = c.SignCoze(coze)
	if err != nil {
		return nil, err
	}
	c.Rvk = now // Sets `Key.Rvk` to the same value as the self-revoke coze.
	return coze, nil
}

// IsRevoked returns true if the given Key is marked as revoked.
func (c Key) IsRevoked() bool {
	return c.Rvk > 0
}

// IsRevoked will return whether is the given `pay` or Coze key is
// revoked. Returns true on error.
func IsRevoked(topLevel json.RawMessage) bool {
	r := &struct {
		Rvk int64 `json:"rvk"`
	}{}

	err := json.Unmarshal(topLevel, r)
	if err != nil {
		return true
	}
	return r.Rvk > 0
}
