## GDown

a graceful shutdown tool.

## Install

```shell
go get -u github.com/bostin/gdown
```

## Usage

```go
// declare shutdown
ctx, cancel := context.WithCancel(context.Background())
shutdown := gdown.NewGraceful(ctx, cancel)

// register callbacks
shutdown.Register(gdown.PriorityLevel10, func() {
	// called on shutdown
})

// listening
shutdown.Listen()
```
