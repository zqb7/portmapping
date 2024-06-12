package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zqb7/portmapping"
)

var NETS []*portmapping.NetConn

var flagPort = flag.Int("port", 0, "")

func main() {
	flag.Parse()
	if *flagPort > 0 {
		go Web(*flagPort)
	}

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
	c1 := make(chan os.Signal, 1)
	signal.Notify(c1, syscall.SIGHUP, syscall.SIGTERM)
	select {
	case <-c1:
	}
}
