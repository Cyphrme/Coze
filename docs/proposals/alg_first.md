# This is deprecated by changes made to Canon.  Canon no longer has to be in Unicode order.  
The "tilda encapsulated payload" is still relevant for other situations. 

# "Alg First" 
## Problem Concerning Long Messages

*Coze wasn't designed for long messages*, but this proposal considers solutions
accommodating efficient processing of long messages.

In order to hash `pay`, `alg` must be parsed before hashing.  As soon as `alg`
is known hashing `pay` can begin.  It's beneficial to have `alg`, and sometimes
the standard fields ["iat", "tmb", "typ"] as soon as possible for long messages.
For short messages this isn't really relevant.  


## Background
Coze is ASCII and Unicode sorted.  The following is valid JSON demonstration of
correct sorting.

```json
{
	"": "blank",
	" ": "space",
	" alg": "spaceAlg",
	"!": "bang",
	"alg": "ES256", // normal "alg"
	"~": "tilda"
}
 ```

## Potential solutions: 

0. Do nothing. (Leave it alone.)
1. Sort "alg" first all the time.
2. A flag for "alg" sorted first.
3. "!" or "space" characters prepended to Coze standard fields. e.g. " alg".  No
   repeat.
4. "" an alias for `alg`. 
5. Repeat sorted first `alg` in `pay`, (e.g. " alg" and "alg") 
6. Repeat `alg` in `coze`. (`alg` in `coze` alone isn't considered a viable option
   since it wouldn't be signed)
7. Use "pay" for payload. 
	a. pay in `coze` 
	b. pay in `pay`
8. Use `"~"` for payload.  (Last ASCII character)


## Pros/cons
 - (1, 2) break sorting.
 - (3, 4, 5, 6, 7.a) is compatible with existing logic, only additional is
   needed. 
 - (5, 6) repeating results in two sources of truth and is not only bad practice
   but also increases the sizes of the payloads. 
 - (7.b) needs no addition, but `typ` and `tmb` will appear after pay.  
 - (8) Needs no additions, no required special treatment, and fields are sorted
   correctly.  Cons: a strange name/json key.  It is a symbol with special
   meaning which may need additional logic for implementing applications.


# Conclusion:
 Option 8 appears best.  

 No change to Coze.  


## Example Solution (tilda encapsulated payload):

```json
 "pay": {
  "alg": "ES256",
  "iat": 1623132000,
  "tmb": "0148F4CD9093C9CBE3E8BF78D3E6C9B824F11DD2F29E2B1A630DD1CE1E176CDD",
  "typ": "cyphr.me/msg/create",
  "~": {
   "msg": "tilde encapsulated payload"
  }
 }
 ```








