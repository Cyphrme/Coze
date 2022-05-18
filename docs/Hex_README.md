# Hex

# Why support Hex?
Hex is easy to debug and easy to implement. 

Excluding binaries from Coze makes base64 encoding less useful since any
efficiency gains are more relevant for long binary values than short digests.  

## Isn't Hex "inefficient"?
Binary is "inefficient" when encoded into ASCII/UTF-8.  Coze heavily uses
digests which are short and easily stored as bytes.


## Why Hex? 
Hex is widely supported and easy to implement. Coze requires that byte
information is always represented as upper case (majuscule), left padded Hex
string.  

Pros of Hex:
- Easy to implement across various systems.
- Easily identifiable as Hex. 
- Hex is a decent phonetic alphabet.  Phonetic concerns were a driving reasons
  for Bitcoin's creation of Base58.
- No special characters.  (Like dash "-" or equal "=")
- Hex is more ergonomic, like JSON itself.


## Why Hex as a string (in quotes)?
While JSON syntax permits small Javascript numbers, octal and hexadecimal are
[explicitly disallowed](https://www.json.org/json-en.html), as well as any other
form of large number/binary encoding.  

>... the octal and hexadecimal formats are not used.

A separate, future JSON5 release of Coze will use the quoteless `0X` notation
for Hex values, but JSON Coze will continue to use quoted Hex.   


## Why "Majuscule" Hex? 
Written "Hex" with an upper case "H" and spoken "upper case Hex" or "majuscule
Hex".  Hex is the name of the encoding, an alphabet paired with an encoding
method.

Pros:
 - Upper case characters are more compatible with various systems.
 - Upper case characters appear first in the ASCII table, and thus many other
   standards.
 - From a base conversion perspective, Hex is correct: 
   - We are fans of [arbitrary base conversion](https://convert.zamicol.com).  
   - By crafting a single master truncatable alphabet, a single alphabet can be
     specified for many applications.  For example, upper case Hex is a 16
     character truncation of Cyphr.me's "Base 64" or RFC 1924's base85.
   - Lower case hex is not a truncation of any widely used industry alphabet. 
   - Upper case Hex is a truncation of many industry alphabets. 
   - Both hex alphabets are extensions of binary, quaternary (base4), octal
     (base8), and base10 (decimal).  Dozenal is typically written in majuscule. 
   - Lower case hex not an extension of any industry alphabet.  
- Historians believe English majuscule letters were invented first.  Other
  languages, such as ancient and modern Greek, use majuscule letters as digits.  
- Numbers are "upper case".  (Oldstyle numbers are not used.) 
 
Cons:
- Some Unix people have an affinity for lower case. 


## Left Padded Hex?
It's just normal Hex. Don't elide beginning 0's. For decimal preceding 0's are
typically elided. Hex does not omit padding.  For example, in the case of
SHA-256 the digest is always 32 bytes and 64 characters in Hex.