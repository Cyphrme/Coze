# Support Order

Instead of forcing all cozies to be canonicalized, permit cozies to not be
canonicalized and in any arbitrary order.  

This also fixes the "alg first" problem.  


# Classes
canon      (can)
require    (req)
order      (ord)
need       (ned)
needOrder  (nod)


A `canon` is an order and a need with no extra fields permitted. Additionally,
canon must always be Unicode sorted.  

A `require` is an order and a need with no extra fields permitted. 

An `order` specifies an order, but not a need.  Additional fields are
permitted after the order fields

A `need` specifies fields that are required to be present, but not an order.
Additional fields are permitted.  

A `needOrder` specifies an order with required fields.  Additional fields are
permitted after the needOrder fields.  


# `typ`
Field `typ` may denote a canon, require, order, need, or needOrder


# Needed Adjustment 
- `hd` needs to replace `cad`.  
- How does sig know what it's signing?  (`pay` is signed as is no matter what.
  If you need extra fields, put it outside of `pay`)