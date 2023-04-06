# Install

```
go get github.com/cyphrme/coze@master
```

# Go Development

As always with Go, use `go test` and verify any changes. 

```
go test && go test -bench=. && (cd normal && go test)
```


## Go mod

See 

[Requiring module code in a local directory](https://go.dev/doc/modules/managing-dependencies#local_directory)
[Coding against an unpublished module](https://go.dev/doc/modules/release-workflow#unpublished)

Go 1.18 adds [workspace mode](https://go.dev/blog/get-familiar-with-workspaces)
to Go, which lets you work on multiple modules simultaneously. See [Tutorial:
Getting started with multi-module
workspaces](https://go.dev/doc/tutorial/workspaces) which details "you can tell
the Go command that you’re writing code in multiple modules at the same time and
easily build and run code in those modules".


We use the module during local development (That should be "no duh", but Go mod
has a gotcha.)

Add the following line to `go.mod` in your other projects for local changes to
apply while doing local development.
```go.mod
replace github.com/cyphrme/coze => ../coze
```

Alternatively, use Go workspaces.  


Also do a 

```
go get 
```

For development on untagged commit or a particular branch:

```sh
go get github.com/cyphrme/coze@master
# Or
go get github.com/cyphrme/coze@base64
```

## gofumpt and go-critic for linting

```sh
gofumpt -l -w .
 ```

This project uses https://github.com/go-critic/go-critic and its companion cli
tool: https://github.com/golangci/golangci-lint


It can be installed with Go:

```sh
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2
cd $COZE
golangci-lint run
```

## Go Doc

```
godoc -http=:6060
```

http://localhost:6060/pkg/github.com/cyphrme/coze/


## Visualizer

```
gource -c 4 -s 1 -a 1
```

To record:


```
simplescreenrecorder
```


Wanted to say thank you again, ran `golangci-lint run` again today.

# Screenshots, Gifs

Zami uses `peek` for gif generation. 