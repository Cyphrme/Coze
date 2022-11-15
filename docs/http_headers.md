# Coze and HTTP Cookies and HTTP Headers

When using Coze with HTTP cookies, Coze messages should be minified.  For
example, we've encountered no issues using the first example as a cookie:

```
token={"pay":{"msg": "Coze Rocks","alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w"}; Path=/;  Secure; Max-Age=999999999; SameSite=None
```

See https://www.rfc-editor.org/rfc/rfc7230#section-3.2.6

## Default Assumption: Assume unquoted JSON characters ```{}[]:,".+-``` are HTTP header compatible and that HTTP's `quoted-string` is sufficient for payload strings.

Well structured and minified JSON appears to have no compatibility issues, as
JSON strings that may contain arbitrary data are already in HTTP's
`quoted_string` mode. Minification removes unquoted new lines that might be used
at HTTP header delimiters.  

JSON escaping of `\` and `"` in strings is already compatible with HTTP's
escaping.  For example, the characters `"` `\` are already escaped when used in
JSON strings.  


# Testing
```sh
curl --insecure --cookie 'test1={"test1":"v1"}; test2={"test2":"v2"}'  https://localhost:8081/
```

Which results in a HTTP header like the following:

```
GET / HTTP/2.0
Host: localhost:8081
Accept: */*
Cookie: test1={"test1":"v1"}; test2={"test2":"v2"}
User-Agent: curl/7.0.0
```

Which is perfectly acceptable.  


# "If the above was wrong" solutions:
## Assuming some JSON characters `whitespace DQUOTE` and `;` may not be sufficiently escaped.

Since the cookie payload is likely to be in control of the web application, the
application should elect to not construct Coze messages that:

1. Use the characters ";". 
2. End strings in whitespace, as that would result in whitespace double quote.  

Arbitrary JSON may be out-of-the-box incompatible with some clients.  JSON
encoders should already escape `"` and `\`, but they do not escape `;`.  Since
the only valid place for `;` to appear in JSON is strings, simply do not use the
value or URL encode arbitrary strings in payloads.  

TODO created HTTP header escaper for this circumstance.  Also consider double
quote to single quote (with then single quote escaping) escaper.  

## Assuming some JSON characters ```{}[]:,".+-``` (and `\` used as an escape) are HTTP header incompatible.

### Use quoted-string on the JSON payload:
Simply encapsulate the JSON in HTTP's `quoted-string` format.  This maintains
human readability with minimal overhead.  

```HTTP
token="{\"pay\":{\"msg\": \"Coze Rocks\",\"alg\":\"ES256\",\"iat\":1623132000,\"tmb\":\"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk\",\"typ\":\"cyphr.me/msg\"},\"sig\":\"Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w\"}"; Path=/;  Secure; Max-Age=999999999; SameSite=None
```

The original string is **228** characters.  With quote escaping, it's **254**,
about 10% overhead.  Base64 URL truncated (b64ut) encoding would be **304** characters,
about 25% overhead.  

### Other options.  

If the advice of using minified and well structured JSON is insufficient, we
suggest one of the following strategies:  
1. Forget using HTTP headers and use HTTP POST. POST does not have URL character
   safety concern for payloads.  
	- Many JSON APIs already use post.  
	- There are also the JSON MIME type for content: `application/json`
2. URL Base64 encode the any Coze messages.  
	- This has the disadvantage of requiring special encoding, ballooning payload
  sizes, and increasing complexity, but it has the advantage guaranteeing
  safety.  Since Coze messages are already small, this isn't too much of a
  penalty.  
3. URL Encode your payloads.  
	- This is not only ugly, but may cause further compatibility issues with other
	HTTP systems. 

# Design philosophy (Especially in case of insufficiency of the preceding)

0. Every effort should be made to keep JSON unchanged.
	- Any sort of escaping/sanitization should be minimally invasive.  
1. JSON is webby.  
	- Any efforts to marry JSON and HTTP should strive to preserve JSON and be
	minimally invasive to JSON and HTTP.  
2. JSON is decently well suited for HTTP.  
	- Well structured JSON is already suitable for most HTTP headers.  
3. Many applications ignore the strict RFC character requirements for HTTP
   headers.  
	- For example, Chrome produces a header that ignores some HTTP header advice:  
		`Sec-Ch-Ua: "Chromium";v="92", " Not A;Brand";v="99", "Google Chrome";v="92"` 
4. Many web applications already use JSON in HTTP headers.
	- There are many in-the-wild examples of web applications storing JSON in headers such as cookies.  


If incompatibility is discovered, we should also push for clients to support JSON.  



# Further Reading

HTTP spec: https://www.rfc-editor.org/rfc/rfc7230#section-3.2.6 
Cookie spec: https://www.rfc-editor.org/rfc/rfc6265#section-4.1.1

Go's discussion:
Go's cookie sanitization code: https://cs.opensource.google/go/go/+/refs/tags/go1.19.3:src/net/http/cookie.go;drc=2041bde2b619c8e2cecaa72d986fc1f0d054c615;l=399
[net/http: support cookie values with comma](https://github.com/golang/go/issues/7243)

Cyphr.me also has custom logic in `cookies.go`.  


# Attempting to get answers about HTTP cookie acceptable characters
Zami tried to email the [11 year old proposed cookie RFC
6265](https://datatracker.ietf.org/doc/rfc6265) and it looks like they're gone.  

Stuck as a proposal forever? There are so many loose ends!

The original cookie RFC just says quoted strings as cookie values are okay,
which leaves a lot to be desired, but a reasonable assumption would be: "we
delegate the semantics to the HTTP standard".

So this is how I would "fix" the cookie RFC 6265. The cookie RFC should
explicitly defer to the HTTP RFC semantics regarding escapes as defined in RFC
7230 3.2.6.  https://www.rfc-editor.org/rfc/rfc7230#section-3.2.6

This means these characters should be included for `quoted-string`:

HTAB, SP
!#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~
%x80-FF

The characters " and / may be included if escaped.  

Done.  