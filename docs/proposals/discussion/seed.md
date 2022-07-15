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