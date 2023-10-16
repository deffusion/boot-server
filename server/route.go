package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Report struct {
	Node string
	Info string
}

func (r Report) String() string {
	return fmt.Sprintf("%s received %s", r.Node, r.Info)
}

func Route(rg *gin.RouterGroup) {
	rg.POST("/report", func(ctx *gin.Context) {
		var report Report
		ctx.BindJSON(&report)
		log.Println(report)
		ctx.String(http.StatusOK, "ok")
	})
	rg.POST("/stopContainers", func(ctx *gin.Context) {
		server := ctx.Value("server").(*Server)
		server.StopContainers()
		ctx.String(http.StatusOK, "ok")
	})
}
