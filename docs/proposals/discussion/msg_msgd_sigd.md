Adding three more fields to Coze:


`msg`: utf-8 encoded text.  
`msgd`: base64 encoded digest of a message.  (May be used alone and not with `msgd`)
`sigd`: digest of a signature.  Allows future crypto to keep small sizes when transmitting proof of signature. 

`msg`, `msgd`, and `sigd` can be implemented in Coze Standard.  While msg and
msgd are conveniences, `sigd` is relevant for post-quantum and will eventually
belong in Coze Core if not implemented immediately in Coze Core.  

For now, `czd` can be to be used in place of `sigd`.




-------------------------------

As a fun thought, to sign tweets, the last tweet can be the following, where the
previous tweet's thread is the content of `msg`, and `msgd` is the digest of the
omitted `msg`.

```json
cyphr.me:
{"pay":{"msgd":"2Nw_gaosyBHwWvSmIOyKzd3UOLdC-Koog8BAakYv3tI","alg":"ES256","iat":1678306319,"tmb":"9PcBWntvjAktwfiPp8WxgOyQOwc1h6Lo1UnB_gkWXKk","typ":"cyphr.me/tweet/sign"},"sig":"Ep6jR8bWPUWD7RfYC-XkG1nnaP2VgBMr1NCP3d9D8uxvmTvdqLNrNkCQ-zfMr9pN6gdtUYmfL5uaY4jglxt4gw"}
```

So then on Twitter you can sign any tweet/series of tweets you make. That
payload is 278 characters, so it will fit in a single tweet.  


Also, the tweet can be referred to by `czd`.

```
Tweet signed by https://cyphr.me/e/bHP8F_J8OYn6I3l80C13RkHLjU9O84ZFKAvYLkr04bA
```

To avoid the weakest link problem, the alg directly can be specified.  (TODO needs to be supported on Cyphr.me)
```
cyphr.me/ES256/bHP8F_J8OYn6I3l80C13RkHLjU9O84ZFKAvYLkr04bA
```
