package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/zqb7/portmapping"
)

var NETS []*portmapping.NetConn

func main() {
	Web()
	items, err := portmapping.LoadCSV("mapping.csv")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, item := range items {
		net := portmapping.NewNetConn(item)
		NETS = append(NETS, net)
		if item.Status == false {
			continue
		}
		go portmapping.ListenNet(net)

	}
	go func() {
		if err := http.ListenAndServe("127.0.0.1:8081", nil); err != nil {
		}
	}()
	c1 := make(chan os.Signal, 1)
	signal.Notify(c1, syscall.SIGHUP, syscall.SIGTERM)
	select {
	case <-c1:
	}
}
