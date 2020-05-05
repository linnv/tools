package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/linnv/logx"
)

func main() {
	dir := flag.String("d", "./", "dir to serve")
	pp := flag.String("p", "9081", "port listen")
	flag.Parse()
	if !flag.Parsed() {
		os.Stderr.Write([]byte("ERROR: logging before flag.Parse"))
		return
	}

	port := ":" + *pp
	routers := gin.Default()
	// curl xxx:9081/
	routers.Use(static.Serve("/", static.LocalFile(*dir, true)))

	// curl xxx:9081/a // e.g. /$prefix   , static.Serve("/$prefix", static.LocalFile($localDirAbs, true))
	// routers.Use(static.Serve("/a", static.LocalFile(*dir, true)))

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

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	logx.Debugfln("SimpleHttp serving dir : %s", *dir)
	logx.Debugfln("http server listen on port %s", port)
	logx.Debugfln("use c-c to exit ")
	<-sigChan
	server.Shutdown(nil)
}
