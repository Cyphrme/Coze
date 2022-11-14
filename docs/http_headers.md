# Coze and HTTP Cookies and HTTP Headers

When using Coze with HTTP cookies, Coze messages should be minified.  For
example, we've encountered no issues using the first example as a cookie:

```
token={"pay":{"msg": "Coze Rocks","alg":"ES256","iat":1623132000,"tmb":"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk","typ":"cyphr.me/msg"},"sig":"Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w"}; Path=/;  Secure; Max-Age=999999999; SameSite=None
```

See https://www.rfc-editor.org/rfc/rfc7230#section-3.2.6 (further reading at
https://www.rfc-editor.org/rfc/rfc8187)

## Default Assumption: Assume unquoted JSON characters ```{}[]:,".+-``` are HTTP header compatible and that HTTP's `quoted-string` is sufficient for payload strings.

Well structured and minified JSON appears to have no compatibility issues, as
JSON strings that may contain arbitrary data are already in HTTP's
`quoted_string` mode. Minification removes unquoted new lines that might be used
at HTTP header delimiters.  

JSON escaping of `\` and `"` in strings is already compatible with HTTP's
escaping.  For example, the characters `"` `\` are already escaped when used in
JSON strings.  


## Assuming some JSON characters `whitespace DQUOTE` and `;` may not be sufficiently escaped.

Since the cookie payload is likely to be in control of the web application, the
application should elect to not construct Coze messages that:

1. Use the characters ";". 
2. End strings in whitespace.

Arbitrary JSON may be out-of-the-box incompatible with some clients.  JSON encoders should already escape
`"` and `\`, but they do not escape `;`.  Since the only valid place for `;` to
appear in JSON is strings, simply do not use the value or URL encode arbitrary
strings in payloads.  

TODO created HTTP header escaper for this circumstance.

## Assuming some JSON characters ```{}[]:,".+-``` (and `\` used as an escape) are HTTP header incompatible.

### Use quoted-string on the JSON payload:
Simply encapsulate the JSON in HTTP's `quoted-string` format.  This maintains
human readability with minimal overhead.  

```HTTP
token="{\"pay\":{\"msg\": \"Coze Rocks\",\"alg\":\"ES256\",\"iat\":1623132000,\"tmb\":\"cLj8vsYtMBwYkzoFVZHBZo6SNL8wSdCIjCKAwXNuhOk\",\"typ\":\"cyphr.me/msg\"},\"sig\":\"Jl8Kt4nznAf0LGgO5yn_9HkGdY3ulvjg-NyRGzlmJzhncbTkFFn9jrwIwGoRAQYhjc88wmwFNH5u_rO56USo_w\"}"; Path=/;  Secure; Max-Age=999999999; SameSite=None
```

The original string is **228** characters.  With `"` escaping, it's **254**,
about 10% overhead.  Base64 URL truncated encoding would be **304** characters,
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

# Design philosophy in case of insufficiency

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

