# Signature Malleability
Coze requires non-malleable signature schemes because `sig` and `czd` are used
as identifiers.  Signatures must not be mutatable by third parties. Prohibiting
signature malleability makes `czd` useful in preventing replay attacks and helps
prevent applications from making other bad assumptions.

Without consideration for malleability, elliptic curve signatures scheme may be
mutated by third parties.  The no-malleability constraint is already adopted by
some existing standards and it is easily applied to the remaining standards.
Non-malleability is expected to apply to future standards as it is now
considered best practice.  

[Modern Ed25519](https://www.rfc-editor.org/rfc/rfc8032#section-8.4) already
makes a malleability prohibition. However, be aware that some older and RFC
non-compliant libraries do not implement this prohibition.  Ed libraries should
be tested for low-S when implementing Coze to make sure they are RFC compliant.  

For ECDSA, the "low-S" rule must be implemented over most existing libraries.
Bitcoin and Ethereum have both implemented the "low-S" rule and Paul's
noble/curves library supports "low-S".  
- https://github.com/bitcoin/bips/blob/master/bip-0146.mediawiki#low_s
- https://eips.ethereum.org/EIPS/eip-2

# Why Malleability is Dangerous: Replay Attack Scenario
A user signs a message, "sign me into example.com".  Unaware of the replay
potential, the user shares the signed message publicly.  A third see the message
 and mutates the signature to another form and sends the message to example.com.
Example.com uses the `czd` to prevent re-login attempts, and with a new
signature the value of `czd` is changed so it accepts the second message as a
new sign in request.  The third party is now also logged into example.com using
a valid message.  

Coze prevents this scenario by requiring signatures to be non-malleable which
makes `czd` useful in identify previously processed messages.  Third parties
cannot mutate an existing signature to another valid form.  


### Future considerations
If for some reason a future algorithm cannot make a no malleability guarantee,
the suggestion is to leave `sig` and `czd` empty and populate a new fields
in `coze` specially designated for malleable signatures.  However, this is
expected to be unlikely, and we'd most likely advocate for non-adoption of such
standard.  

### Go Code to generate Malleable Signatures
```golang
// Example_GenHighSCoze generates high s.  Must comment out S canonicalization
// in verify and sign for this to work.
func Example_GenHighSCoze() {
	goEcdsa := KeyToPubEcdsa(&GoldenKey)

	for i := 0; i < 10; i++ {
		cz := new(Coze)
		err := json.Unmarshal([]byte(GoldenCoze), cz)
		if err != nil {
			panic(err)
		}

		err = GoldenKey.SignCoze(cz)
		if err != nil {
			panic(err)
		}

		size := GoldenKey.Alg.SigAlg().SigSize() / 2
		s := big.NewInt(0).SetBytes(cz.Sig[size:])

		ls, _ := IsLowS(goEcdsa, s)
		if !ls {
			fmt.Printf("High-S coze: %s\n", cz)
		}
		fmt.Printf("Low-S coze: %s\n", cz)
	}
	// Output:
}
```


# Other Links
 - [rfc 6979 "Deterministic Usage of the Digital Signature Algorithm (DSA) and
   Elliptic Curve Digital Signature Algorithm
   (ECDSA)"](https://www.rfc-editor.org/rfc/rfc6979)
 - [Signature Malleability on Wikipedia](https://en.wikipedia.org/wiki/Malleability_(cryptography))
- [Non-modern Ed malleability demonstration](https://slowli.github.io/ed25519-quirks/malleability)




