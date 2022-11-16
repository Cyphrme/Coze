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
	alg:x:
	alg:d::
	alg:izd:::


## Always Useful checksummed forms 
If always processed when possible.  

### Useful
These are the only forms that `chk` is useful.  All other forms are redundant.  

```
alg:izd~chk:::
alg:d~chk::
alg:x~chk:
alg:tmb~chk
```

### Redundant (don't do this)
Any combination of two or more Coze fields [`izd`,`d`,`x`,`tmb`], including the following:

- `alg:x:tmb~chk`
    - tmb already serves as the checksum of x.  
- `alg:d~chk:x:`
   - x is derived from d and so x may serve as a checksum.    
- `alg:d~chk::tmb`
    - tmb is derived from x, which x is derived from d and so tmb may serve as a checksum for d.
- `alg:d~chk::tmb~chk`
- `alg:d~chk:x~chk:`
- `alg:d~chk:x~chk:tmb~chk` 
- `alg:izd~chk:d~chk:x~chk:tmb~chk` 

The only reason to do these forms is if your system does not have the ability to
do signatures, but does have the ability to check digests.  This would be a
mostly unreasonable circumstance.  

