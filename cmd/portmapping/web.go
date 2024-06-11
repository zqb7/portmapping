package main

import (
	"fmt"
	"html/template"
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
	for i, nc := range NETS {
		if fmt.Sprintf("%d", i) == index {
			netConn = nc
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

func Web() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/action", handleAction)
	go updateConfig()
}
