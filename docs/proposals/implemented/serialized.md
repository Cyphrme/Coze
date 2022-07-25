# Serialized form

Coze key serialized:

	alg:d:x:tmb

Can be shorted to public:

	alg:x:tmb

Which can be shortened to just thumbprint:

	alg:tmb

With x and no tmb

	alg:x:

With just d:

	alg:d::


# Checksums and Reveal.

// TODO think about this more.  
If tmb is first class, there is no checksum for tmb itself.  Suggested form for checksums:

	alg:tmb~checksum

Alternatively:

	alg:::tmb:checksum




## See also:
https://github.com/multiformats/multiformats
https://multiformats.io/multihash/