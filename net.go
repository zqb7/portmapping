package portmapping

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

func ListenNet(ctx context.Context, item *Item) {
	for {
		listener, err := net.Listen(item.Network, fmt.Sprintf("0.0.0.0:%d", item.Port))
		if err != nil {
			return
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go dial(ctx, conn, item.Network, item.TargetIP, item.TargetPort, 20*time.Second)
		}
	}
}

func dial(ctx context.Context, from net.Conn, network, targetIP string, targetPort int, retryWait time.Duration) {
	retryChan := make(chan struct{}, 1)
	var needWait bool
	go func() { retryChan <- struct{}{} }()
	for {
		select {
		case <-ctx.Done():
			return
		case <-retryChan:
			if needWait && retryWait > 0 {
				select {
				case <-time.NewTimer(retryWait).C:
				case <-ctx.Done():
					return
				}
			} else {
				needWait = true
			}
			conn, err := net.Dial(network, fmt.Sprintf("%s:%d", targetIP, targetPort))
			if err != nil {
				retryChan <- struct{}{}
				continue
			}
			go func() { io.Copy(conn, from) }()
			go func() { io.Copy(from, conn) }()
		}
	}
}
