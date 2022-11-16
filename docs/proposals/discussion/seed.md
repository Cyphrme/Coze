# Seed (`izd`) proposal
`d` must be derived from `izd`.  Already, `x` must be derived from `d`.  
Seed may be in Coze standard, not Coze core.  (Thinking about this).

Seed's three letter Coze reserved name ideas: 

"InitialiZation seeD":   `izd`
"Initialization SeeD":   `isd`
"seed":                  `sed`


We used to have `sed` in Coze, it was `czd`.  It stood for "SignaturE Digest".






# Ed25519
See also, Ed25519's naming differences in
implementations, (see
https://github.com/Cyphrme/ed25519_applet#naming-differences-in-implementations)

In Coze, "x" is the RFC's seed.  Secrete scalar s is recomputed.    


# Coze's `x` is derived from `d` assumption
Under some signing schemes, d's may not be knowably related to x. If only d is
known, and x is not related, d would directly need a checksum.  

Coze is going to assume x is derivable from d.  If we ever have to break this
assumption, we would have to figure out something.  Maybe a different Coze key
field name for such values.    

Some signing schemes may not have the relationship where the public component
can be derived from the private component.  Because of this, tmb cannot serve as
the checksum for d.  Coze is designed to avoid this situation by requiring that
x is always derivable from d, and so `x` should always be able to serve as
checksum for `izd` and `x`.   For example, this is avoided in Ed25519 by using
the RFC's "seed" instead of secrete scalar s (sss) for the value of `d`.  This
is common practice by other implementations as well. .  sss may be stored in
Coze keys separately, but perhaps field sanitization would be useful to build
into Coze (right now we think it is not).  