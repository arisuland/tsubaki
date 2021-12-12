package main

import (
	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"context"
	"os"
)

var (
	version    string
	commitHash string
	log        slog.Logger
)

func init() {
	log = slog.Make(sloghuman.Sink(os.Stdout))
}

func main() {
	log.Info(context.Background(), "hello world :3")
}
