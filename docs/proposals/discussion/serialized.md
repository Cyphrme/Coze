# Serialized Form for Coze Keys
Colon `:` is used as the serialized component delimiter.  Although "url unsafe",
as it is widely used in URLs without issue.

# Serialized Forms 

`chk` is checksum.  

	alg:izd:d:x:tmb
	alg:d:x:tmb
	alg:d:x:
	alg:x:tmb
	alg:tmb
	alg:tmb~chk
	alg:d::tmb
	alg:izd:::


# Why Specify Alg?  (Weakest link problem and "bad" crypto agility.  )



An alternative is to:
 - Assume that all digests of a similar size have a similar strength. All
   digests in a given class are considered secure.  If a class gets a "weak"
   member, applications can require explicit alg denotation (This is
   hard/tricky/messy, so this is not something we'd advocate for. )