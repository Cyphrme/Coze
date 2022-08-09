# Serialized form

Colon delimit values. 

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


With just seed:

alg:seed:::

Full
alg:seed:d:x:tmb


# Seed

Although we have not implemented seeds, this documents is our thoughts
on how we would.  Ed25519's naming differences in implementations, (see
https://github.com/Cyphrme/ed25519_applet#naming-differences-in-implementations)
highlights the possible future need.  


# (Bad) Option 1: New Key field
All logic now has to worry about two fields, `d` and `seed`

# Option 2: Deliminator and Concatenate

```
"d":"seed:private d"
```


With just seed:
```
"d":"seed:"
```

And with just private d:

```
"d":"private d"
```


# Checksums and Reveal.

// TODO think about this more.  
If tmb is first class, there is no checksum for tmb itself.  Suggested form for checksums:

	alg:tmb~checksum

Alternatively:

	alg:::tmb:checksum




## See also:
https://github.com/multiformats/multiformats
https://multiformats.io/multihash/