// Normal - A normal is an arrays of fields specifying the normalization of
// payload.  Normals are implemented in Go as []string.  There are six types of
// normals, and a nil normal is valid.
//
//  canon       (can)
//  only        (ony)
//  need        (ned)
//  option      (opt)
//  extra       (ext)
//  (nil)
//
// An nil normal matches all payloads.
//
// `canon` requires specified fields in the given order.
//
// `only` specifies fields that are required to be present, does not specify any
// order.
//
// `option` specifies permissible fields in a given order. (no extras)
//
// `need` specifies fields that are required to be present, but does not specify
// any order and allows extra fields (after need?). 
//
// `extra` specifies optional fields in a given order and allows extra fields
// are permitted after the extra fields.  A nil extra allows extra fields.
//
// (nil) 000
// (no Canon/Normal) only extras


// A an empty canon of [] matches only an a empty (i.e. `{}`) payload.  
//
// # Using with Option
// When a need is used with an option all fields are unordered.
//
// When an order is used with an option, all fields are ordered.
//
// ## Normal, Require, and Option
//
// `canon`, `only`, `need`, and `order` are valid `require` in that they specify
// required fields.  An option is distinct in that option specifies optional
// fields and precludes other optional fields.
//

(required:x, order:x, alphabetical:x, [extra:0, ])

```Go
// Future addition, JSON Schema
type Schema struct{
	form json.RawMessage //JSON schema (i.e.{type: "object",properties: {})
	// Additonal flags
	Order bool
	UTFSorted bool
}
```

Note that schemas have no order (since they are JSON objects), but jsonfrom does
have an order.  
Order need to be added over JSON. 



See also: 
http://xml.coverpages.org/SchemaCentricCanonicalization-20020213.html

If JSON schema is supported in the future, we would add an option for order.  
https://github.com/iakovmarkov/json-schema-normalizer
(Also fun: https://rjsf-team.github.io/react-jsonschema-form/)