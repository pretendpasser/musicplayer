package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	apihttp "player/api/http"
)

type config struct {
	HttpAddr string
}

var gConfig config

func GetEnvDefault(key, defalut string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defalut
	}
	return val
}

func initConfigFromENV() {
	gConfig = config{
		HttpAddr: GetEnvDefault("HTTP_LISTEN_ADDR", ":9582"),
	}
	fmt.Println("config defined in env", gConfig)
}

func init() {
	initConfigFromENV()
}

func main() {
	httpSrv := http.Server{Addr: gConfig.HttpAddr, Handler: apihttp.NewHTTPHandler()}
	errc := make(chan error)
	go func() {
		fmt.Println("HTTP server Start", gConfig.HttpAddr)
		errc <- httpSrv.ListenAndServe()
	}()

	err := <-errc
	log.Println(err)
}
