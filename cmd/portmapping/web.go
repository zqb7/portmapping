package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zqb7/portmapping"
)

var (
	updateConfigChan = make(chan struct{}, 16)
	server           = &Server{}
)

type Server struct {
}

func (s *Server) toClient(nets []*portmapping.NetConn) (d [][]interface{}) {
	for _, v := range nets {
		var v2 []interface{}
		v2 = append(v2, v.Port, v.Network, v.Status, v.TargetHost, v.TargetPort, v.ClientNumber)
		d = append(d, v2)
	}
	return
}

func (s *Server) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tpl", NETS)
}

func (s *Server) Add(c *gin.Context) {
	item := portmapping.Item{Network: "tcp"}
	NETS = append(NETS, portmapping.NewNetConn(&item))
	updateConfigChan <- struct{}{}
	c.JSON(http.StatusOK, s.toClient(NETS))
}

func (s *Server) Start(c *gin.Context) {
	for i, nc := range NETS {
		if fmt.Sprintf("%d", i) == c.Param("index") {
			go portmapping.ListenNet(nc)
			updateConfigChan <- struct{}{}
			c.JSON(http.StatusOK, gin.H{})
			return
		}
	}
}

func (s *Server) Stop(c *gin.Context) {
	for i, nc := range NETS {
		if fmt.Sprintf("%d", i) == c.Param("index") {
			nc.Stop()
			updateConfigChan <- struct{}{}
			c.JSON(http.StatusOK, gin.H{})
			return
		}
	}
}

func (s *Server) Delete(c *gin.Context) {
	for i, nc := range NETS {
		if fmt.Sprintf("%d", i) == c.Param("index") {
			nc.Stop()
			NETS = append(NETS[:i], NETS[i+1:]...)
			updateConfigChan <- struct{}{}
			c.JSON(http.StatusOK, gin.H{})
			return
		}
	}
}

func (s *Server) Update(c *gin.Context) {
	var reqBody = make(map[string]string, 0)
	err := c.ShouldBindBodyWithJSON(&reqBody)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var netConn *portmapping.NetConn
	for i, nc := range NETS {
		if fmt.Sprintf("%d", i) == c.Param("index") {
			netConn = nc
			break
		}
	}
	if netConn == nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var isUpdate bool
	if desc, ok := reqBody["desc"]; ok {
		netConn.Desc = desc
		isUpdate = true
	}
	if isUpdate {
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

func Web(port int) {
	e := gin.Default()
	e.LoadHTMLFiles("index.tpl")
	e.GET("", server.Index)
	e.POST("/start/:index", server.Start)
	e.POST("/stop/:index", server.Stop)
	e.POST("/del/:index", server.Delete)
	e.POST("/update/:index", server.Update)
	e.POST("/add", server.Add)
	go updateConfig()
	log.Fatalln(e.Run(fmt.Sprintf("127.0.0.1:%d", port)))
}
