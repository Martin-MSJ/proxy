package main

import (
	"fmt"
	"net/http"

	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().DoFunc(func(r *http.Request,ctx *goproxy.ProxyCtx) (*http.Request, *http.Response){
		fmt.Println("haha");
		return r, nil
	})
	proxy.Verbose = true
	http.ListenAndServe("127.0.0.1:8080", proxy)
}

