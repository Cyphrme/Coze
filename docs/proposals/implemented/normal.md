# Support Normal

Instead of forcing all cozies to be canonicalized with a unicode sorted canon,
permit cozies to have any arbitrary order and expand normalization.

This also fixes the "alg first" problem.  

Normalizations may be denoted by `typ` by applications.

See also: 
http://xml.coverpages.org/SchemaCentricCanonicalization-20020213.html

If JSON schema is supported in the future, we would add an option for order.  
https://github.com/iakovmarkov/json-schema-normalizer
(Also fun: https://rjsf-team.github.io/react-jsonschema-form/)