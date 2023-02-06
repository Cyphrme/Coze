# Signature Malleability
Coze requires that signature schemes be non-malleable, meaning that it must not
be possible for signature mutation by third parties after signing.  Prohibiting
signature malleability makes `czd` useful in preventing replay attacks and helps
prevent applications from making bad assumptions when using `sig` as an
identifier.

Without consideration for malleability, elliptic curve signatures scheme may be
mutated by third parties.  The no-malleability constraint is already adopted by
some existing standards, it is easily to apply to the remaining standards, and
it is expected to apply to future standards as well as it is now considered to
be best practice.  

[Modern Ed25519](https://www.rfc-editor.org/rfc/rfc8032#section-8.4) already
makes a malleability prohibition. However, be aware that there are older
libraries and RFC non-compliant that do not implement this prohibition.  Ed
libraries should be tested for low s when implementing Coze to make sure they
are RFC compliant.  

For ECDSA, the "low s" rule must be implemented over most existing libraries.
For more detail, see
- https://github.com/bitcoin/bips/blob/master/bip-0146.mediawiki#low_s
- https://eips.ethereum.org/EIPS/eip-2

# Replay Attack Scenario
A user signs a message, "sign me into example.com".  Unaware of the replay
potential, the user shares the signed message publicly.  A third see the message
 and mutates the signature to another form and sends the message to example.com.
Example.com uses the `czd` to prevent re-login attempts, and with a new
signature the value of czd is changed so it accepts the second message as a new
sign in request.  The third party is now also logged into example.com using a
valid message.  

Coze prevents this scenario by requiring signatures to be non-malleable.  This
allows `czd` identify previously processed messages.  The third party then
cannot mutate an existing signature to any other valid form.  


### Future considerations
If for some reason a future algorithm cannot make no malleability guarantee,
then the suggestion is to leave `czd` empty and populate a new field in `coze`.
However, this is expected to be unlikely, and we'd most likely advocate for
non-adoption of such standard.  




# Other Links
Non-modern Ed malleability demonstration:
https://slowli.github.io/ed25519-quirks/malleability




