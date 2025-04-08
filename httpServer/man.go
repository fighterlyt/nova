package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	port             = `0`
	response         = ``
	responseFilePath = ``
)

func flagDeclare() {
	flag.StringVar(&port, "p", port, "http server port")
	flag.StringVar(&response, "r", response, "http server response")
	flag.StringVar(&responseFilePath, `f`, responseFilePath, "http server response file path")
}
func main() {
	flagDeclare()

	flag.Parse()

	var ()
	if responseFilePath != `` {
		data, err := os.ReadFile(responseFilePath)

		if err != nil {
			panic(`读取应答文件`)
		}

		response = string(data)
	}

	r := gin.Default()

	r.POST(`/`, func(context *gin.Context) {
		context.String(http.StatusOK, response)
	})

	ln, _ := net.Listen("tcp", `:`+port)

	_, port, _ := net.SplitHostPort(ln.Addr().String())
	fmt.Println("Listening on port", port)

	http.Serve(ln, r)
}
