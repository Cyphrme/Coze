# Go Development

As always with Go, use `go test` and verify any changes. 


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


## Visualizer

```
gource -c 4 -s 1 -a 1
```

To record:


```
simplescreenrecorder
```