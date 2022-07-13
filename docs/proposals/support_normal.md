# Support Normal

Instead of forcing all cozies to be canonicalized with a unicode sorted canon,
permit cozies to posses canon with any arbitrary order.  

This also fixes the "alg first" problem.  

Normalizations may be denoted by `typ` by applications.

# Normal
All normals are arrays of fields.  In Go this is implemented at a slice
[]string.

	canon      (can)
	only       (ony)
	need       (ned)
	order      (ord)
	option     (opt)


`canon` requires specified fields in the given order and no extra fields
permitted. 

`only` specifies fields that are required to be present, does not specify
any order, and no extra fields permitted.

`need` specifies fields that are required to be present, but does not specify
any order. Additional fields are permitted.

`order` requires specified fields in the given order and additional fields
are permitted after the order fields.

`option` specifies permissable optional fields and is used with a `need` or an
`order`. All fields not in `option` and the respective `need` or `order` are
invalid. If option is nil, all extra fields are valid.  

## Normal, Require, and Option

`canon`, `only`, `need`, and `order` are a `require` in that they specify
required fields.  An option is distinct in that option specifies optional fields
and precludes other optional fields.  

             ┌────────────────┐
             │     Normal     │
             └───────┬────────┘
             ┌───────┴────────┐
       ┌─────┴────┐     ┌─────┴──────┐
       │ Require  │     │   Option   │
       └──────────┘     └────────────┘

Normal
	Requires:
		canon
		only
		need
		order
	Option:
		option

# `typ`
Field `typ` may denote a canon, order, or need.

# Implementation
```go
type Normal []string

type Canon Normal 
type Only Normal
type Need Normal
type Order Normal
type Option Normal

func IsNormal(coze Coze, norm any, opt Option) bool {
	var ms = MapSlice{}
	err := json.Unmarshal(coze.Pay, &ms)
	if err != nil {
		return nil, err
	}

	switch v := norm.(type) {
	case Canon:

	case Only:

	case Need:

	case Order:

	default:
		return false
	}
}
```