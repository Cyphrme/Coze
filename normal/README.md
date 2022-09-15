# Normal

We consider Normal to be Coze Standard and not Coze core.


Coze Core
Coze Standard
Coze Experimental


# Example
Use with Coze

Example code (`ExampleIsNormalNeedOption`) can be found in `normal_test.go`.

The following example demonstrates use with Coze Standard for a user account
form.

The form requires the standard Coze fields ("alg", "iat", "tmb", "typ") along
with the field 'id', which is the user's account id for the application.

The Coze Normal for the form would be:
``` Go
Need{"alg", "iat", "tmb", "typ", "id"}
```
Now let's say we allow other optional fields in the form, such as: "display_name",
"first_name", "last_name", "email", "address_1", "address_2", "phone_1", "phone_2",
"city", "state", "zip", and "country".

As a Coze Normal, the options would be the following:

``` Go
Option{"display_name", "first_name", "last_name", "email", "address_1", "address_2", "phone_1", "phone_2", "city", "state", "zip", "country"}
```

The full Coze Normal (normals chained) would become:


``` Go
[
Need{"alg", "iat", "tmb", "typ", "id"},
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

Normal does not check that a payload is cryptographically valid, but is
useful for JSON field validation.

On the other hand, a cryptographically verified payload may not guarantee that
the 'pay' has the required fields present for an application's specific endpoint.

Pairing verification, and normal allows for better and potentially more 
optimized validation of incoming payloads.