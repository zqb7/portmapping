package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zqb7/portmapping"
)

func main() {
	items, err := portmapping.LoadCSV("mapping.csv")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, item := range items {
		go portmapping.ListenNet(ctx, item)
	}
	c1 := make(chan os.Signal, 1)
	signal.Notify(c1, syscall.SIGHUP, syscall.SIGTERM)
	select {
	case <-c1:
	}
}
