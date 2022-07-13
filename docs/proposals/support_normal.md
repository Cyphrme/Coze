# Support Normal

Instead of forcing all cozies to be canonicalized with a unicode sorted canon,
permit cozies to posses canon with any arbitrary order.  

This also fixes the "alg first" problem.  


# Normalizations denoted by `typ`.
API's using Coze may denote 

canon      (can)
only       (ony)
need       (ned)
order      (ord)


A `canon` requires specified fields in the given order and no extra fields
permitted. 

A `only` specifies fields that are required to be present, does not specify
any order, and no extra fields permitted. 

A `need` specifies fields that are required to be present, but does not specify
any order. Additional fields are permitted.  

An `order` requires specified fields in the given order and additional fields
are permitted after the order fields.




# `typ`
Field `typ` may denote a canon, order, or need.

# Implementation
```go
type Normal []string

type Canon Normal //(Rename Canon() to GetCanon())
type Only Normal
type Need Normal
type Order Normal


IsNormal(coze Coze, norm any)bool{
	var ms = MapSlice{}
	err = json.Unmarshal(raw, &ms)
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