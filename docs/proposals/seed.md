# Seed

Although we have no plan in implementing seeds, this documents is our thoughts on how we
would.  

Ed25519's naming differences in implementations highlights the possible future need.  





# (Bad) Option 1: New CozeKey field
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
"d":"seed:"
```