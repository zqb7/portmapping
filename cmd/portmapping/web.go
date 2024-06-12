package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/zqb7/portmapping"
)

var updateConfigChan = make(chan struct{}, 16)

func handleIndex(w http.ResponseWriter, req *http.Request) {
	tpl, _ := template.ParseFiles("index.tpl")
	tpl.Execute(w, NETS)
}

func handleAction(w http.ResponseWriter, req *http.Request) {
	index := req.URL.Query().Get("index")
	event := req.URL.Query().Get("event")
	var netConn *portmapping.NetConn
	var netConnIndex int = -1
	for i, nc := range NETS {
		if fmt.Sprintf("%d", i) == index {
			netConn = nc
			netConnIndex = i
			break
		}
	}
	if netConn == nil {
		http.Error(w, "invalid index", http.StatusBadRequest)
		return
	}
	switch event {
	case "start":
		go portmapping.ListenNet(netConn)
		updateConfigChan <- struct{}{}
	case "stop":
		netConn.Stop()
		updateConfigChan <- struct{}{}
	case "desc":
		v := req.URL.Query().Get("v")
		netConn.Desc = v
		updateConfigChan <- struct{}{}
	case "del":
		if netConnIndex > -1 {
			NETS = append(NETS[:netConnIndex], NETS[netConnIndex+1:]...)
			updateConfigChan <- struct{}{}
		}
	default:
		http.Error(w, "unkonwn event", http.StatusBadRequest)
	}
}

func updateConfig() {
	for {
		select {
		case <-updateConfigChan:
			portmapping.SaveCSV("mapping.csv", NETS)
		}
	}
}

func Web(port int) {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/action", handleAction)
	go updateConfig()
	log.Fatalln(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil))
}
