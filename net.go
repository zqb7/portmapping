package portmapping

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"
)

type NetConn struct {
	*Item
	ClientNumber int64
	listener     net.Listener
	stopChan     chan struct{}
}

func (nc *NetConn) Stop() {
	if nc.Status == false {
		return
	}
	nc.Status = false
	nc.listener.Close()
	nc.stopChan <- struct{}{}
}

func NewNetConn(item *Item) *NetConn {
	net := &NetConn{Item: item, stopChan: make(chan struct{}, 1)}
	return net
}

func ListenNet(nc *NetConn) {
	for {
		listener, err := net.Listen(nc.Network, fmt.Sprintf("0.0.0.0:%d", nc.Port))
		if err != nil {
			log.Printf("listen port:%d err:%s\n", nc.Port, err)
			return
		}
		log.Printf("listen port:%d target_host:%s target_port:%d\n", nc.Port, nc.TargetHost, nc.TargetPort)
		nc.listener = listener
		nc.Status = true
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go dial(nc, conn, 20*time.Second)
		}
	}
}

func dial(nc *NetConn, from net.Conn, retryWait time.Duration) {
	retryChan := make(chan struct{}, 1)
	var needWait bool
	go func() { retryChan <- struct{}{} }()
	for {
		select {
		case <-nc.stopChan:
			return
		case <-retryChan:
			if needWait && retryWait > 0 {
				select {
				case <-time.NewTimer(retryWait).C:
				case <-nc.stopChan:
					return
				}
			} else {
				needWait = true
			}
			conn, err := net.Dial(nc.Network, fmt.Sprintf("%s:%d", nc.TargetHost, nc.TargetPort))
			if err != nil {
				retryChan <- struct{}{}
				continue
			}
			atomic.AddInt64(&nc.ClientNumber, 1)
			go func() { io.Copy(conn, from) }()
			go func() { io.Copy(from, conn); atomic.AddInt64(&nc.ClientNumber, -1) }()
		}
	}
}
