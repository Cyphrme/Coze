# Self-describing encoding
In the future, it would be nice to support other encoding and/or support
self-describing.

One of the less than universal design choices is that Coze opinionatedly uses
JSON, a text encoding, for cryptographic functions. Using text encoding may or
may not be ergonomic for other systems. Coze does this carefully, and we think
efficiently, but it's still a concern.  

Self-described encoding allows binary value to be used in future designs to
lessen the coupling with the text encoding design choice. 

For other future designs to be supported, Coze needs escape.  Since b64ut uses a
small subset of characters, escaping now doesn't seem to be a big issue.  Any
non-base64 character may be used to self-describe encoding as a value. 

To stay JSON compliant, quotes are needed around values.  Extending JSON and
using self-describing syntax would allow the dropping of the quotes.

Cyphr.me has an internal system of self-describing that might be used as well.
See the README in Cyphr.me's basecnv package section "Self-describing".



## Base and Encoding Notation
The following characters are good delimiter candidates:

```
~!.:=
```

The second character is the last character in the alphabet, thus denoting
both base and alphabet.  

For example, to represent Cyphr.me Base 64 (Not b64ut):

```
~_
```

Thus "Zami's Majuscule Key" (if using pure Base 64 and not b64ut) would be
represented by:


```
~_cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk
```


This scheme will already be compatible with Coze since b64ut is not compatible
with the character `~`.  


Further notation would be needed to denote bucket conversion, padding, and other
options.  (The work of generalizing bucket conversion is not yet completed.)

## Special
Perhaps there's a need for a "special" delimiter that denotes "legacy" bases.  
For b64ut

```
!u
```

Where "!" denotes a legacy/standard and the second character is the character
for the [multibase](https://github.com/multiformats/multibase#multibase-table)
encoding. 


Thus "Zami's Majuscule Key"  would be represented by:


```
!ucLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk
```

