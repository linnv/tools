package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linnv/logx"
	"github.com/pkg/errors"
)

const (
	uriPrefixRegister = "/_set_"
	uriPrefixMock     = "/_mock_"
)

var uriMap map[string][]byte

//Register implements set mock uri and it's content to memory
func Register(c *gin.Context) {
	r := c.Request
	setUri := strings.TrimPrefix(r.RequestURI, uriPrefixRegister)
	bs, err := ioutil.ReadAll(r.Body)
	logx.Debugfln(" %s ->  [%s]\n", setUri, bs)
	if err != nil {
		err = errors.Wrapf(err, " body:%v", bs)
		logx.Errorf(" getReq err %v bs [%s]\n", err, bs)
		c.AbortWithStatusJSON(401, err)
		return
	}
	uriMap[setUri] = bs

	c.Writer.WriteHeader(200)
	c.Writer.WriteString(fmt.Sprintf(" %s ->  [%s]\n", setUri, bs))
	return
}

//Mock implements return the content of mock uri registered
func Mock(c *gin.Context) {
	r := c.Request
	mockUri := strings.TrimPrefix(r.RequestURI, uriPrefixMock)
	if getUriBs, ok := uriMap[mockUri]; ok {
		c.Writer.WriteHeader(200)
		c.Writer.Write(getUriBs)
		return
	}
	c.AbortWithStatusJSON(401, "not foud uri "+mockUri)
	return
}

func main() {
	pp := flag.String("p", "9081", "port listen")
	flag.Parse()
	if !flag.Parsed() {
		os.Stderr.Write([]byte("ERROR: logging before flag.Parse"))
		return
	}

	port := ":" + *pp
	routers := gin.Default()
	routers.POST(uriPrefixRegister+"/:any", Register)
	routers.GET(uriPrefixMock+"/:any", Mock)
	routers.POST(uriPrefixMock+"/:any", Mock)
	server := http.Server{
		Addr:    port,
		Handler: routers,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
			} else {
				panic(err)
			}
		}
	}()

	uriMap = make(map[string][]byte, 1)
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	logx.Debugfln("http server listen on port %s", port)
	logx.Debugfln("use c-c to exit ")
	<-sigChan
	server.Shutdown(nil)
}
