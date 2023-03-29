# Normal

Normal is design to make normalizing fields in a Coze payload easy.

A normal is an arrays of fields specifying the normalization of a
payload. Normals may be chained to represent various combinations of
normalization.  There are five
types of normals plus a nil normal.

See the documentation in [normal.go](normal.go)

## Example

Example code (`ExampleIsNormalNeedOption`) can be found in `normal_test.go`.

### Form Example
Given an arbitrary form with the following requirements.  

The application requires "id", which is the user's account id for the
application.

The Normal is:
``` Go
Need{"id"}
```

Also, the form has optional fields in the form, such as:
"display_name", "first_name", "last_name", "email", "address_1", "address_2",
"phone_1", "phone_2", "city", "state", "zip", and "country".

The optional fields as a Normal would is the following:

``` Go
Option{"display_name", "first_name", "last_name", "email", "address_1", "address_2", "phone_1", "phone_2", "city", "state", "zip", "country"}
```

Since the application uses Coze to sign messages, the application decides to require the standard Coze fields ("alg", "iat", "tmb", "typ")

``` Go
Need{"alg", "iat", "tmb", "typ"}
```


The two needs can be summed, and the resulting need is concatenated to the
option.  The full Normal chain becomes:


``` Go
[
Need{"id", "alg", "iat", "tmb", "typ"},
Option{"display_name", "first_name", "last_name", "email", "address_1", "address_2", "phone_1", "phone_2", "city", "state", "zip", "country"}
]
```

The following payload would return 'true', indicating the payload is considered
to be a normal payload:

``` JSON
{
	"alg": "ES256",
	"iat": 1647357960,
	"tmb": "L0SS81e5QKSUSu-17LTQsvwKpUhBxe6ZZIEnSRV73o8",
	"typ": "cyphr.me/user/profile/update",
	"id": "L0SS81e5QKSUSu-17LTQsvwKpUhBxe6ZZIEnSRV73o8",
	"city": "Pueblo",
	"country": "ISO 3166-2:US",
	"display_name": "Mr. Dev",
	"first_name": "Dev Test",
	"last_name": "1"
 }

```

While the following payload would return 'false', since it is missing the
required 'id' field.

``` JSON
{
	"alg": "ES256",
	"iat": 1647357960,
	"tmb": "L0SS81e5QKSUSu-17LTQsvwKpUhBxe6ZZIEnSRV73o8",
	"typ": "cyphr.me/user/profile/update",
	"city": "Pueblo",
	"country": "ISO 3166-2:US",
	"display_name": "Mr. Dev",
	"first_name": "Dev Test",
	"last_name": "1"
 }

```

Normal does not check that a payload is cryptographically valid, but is useful
for JSON field validation. On the other hand, a cryptographically verified
payload may not guarantee that the 'pay' has the required fields present for an
application's specific endpoint.

Pairing verification and normal allows for better and potentially more 
optimized validation of incoming payloads.


# How does Normal relate to Coze?

We consider Normal to be apart of Coze Standard and not Coze core.


1. Coze Core
2. Coze Standard
3. Coze Experimental
