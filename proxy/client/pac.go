package main

import (
	"net/http"
	"log"
	"fmt"
	"net"
	"github.com/getlantern/pac"
	"sync/atomic"
)

const (
	pacTemplate = `function FindProxyForURL(url, host) {
			if (isPlainHostName(host) // including localhost
			|| shExpMatch(host, "*.local")) {
				return "DIRECT";
			}
			// only checks plain IP addresses to avoid leaking domain name
			if (/^[0-9.]+$/.test(host)) {
				if (isInNet(host, "10.0.0.0", "255.0.0.0") ||
				isInNet(host, "172.16.0.0",  "255.240.0.0") ||
				isInNet(host, "192.168.0.0",  "255.255.0.0") ||
				isInNet(host, "127.0.0.0", "255.255.255.0")) {
					return "DIRECT";
				}
			}
			return "PROXY %s; DIRECT";
		}`
)

const (
	localHttpPACServerAddr = "127.0.0.1:38250"			//get pacfile server
)

var (
	isPACOn = int32(0)
	pacFile []byte
)

func pacHandler(proxyURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(genPAC(proxyURL)); err != nil {
			log.Printf("Error writing response: %s", err)
		}
	}
}

func genPAC(proxyURL string) []byte {
	if pacFile == nil {
		pacFile = []byte(fmt.Sprintf(pacTemplate, proxyURL))
	}
	return pacFile
}

func setPAC(){
	openLocalhttpPACServer()
	initPACSetting()
}

func openLocalhttpPACServer(){
	localHttpListener, err := net.Listen("tcp", localHttpPACServerAddr)
	if err != nil {
		log.Fatalf("localHttpGetPACServer listen filed: %s", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/my.pac", pacHandler("127.0.0.1:8080"))

	localHttpServer := &http.Server{
		Handler: mux,
	}

	go func(){
		err := localHttpServer.Serve(localHttpListener)
		if err != nil {
			log.Fatalf("FATAL: localHttpGetPACServer stopped")
		}
	}()
}

func initPACSetting() error {
	err := pac.EnsureHelperToolPresent("pac-cmd", "PAC_SETUP", "")
	if err != nil {
		return fmt.Errorf("Unable to set up pac setting tool: %s", err)
	}
	enablePAC(localHttpPACServerAddr)

	return nil
}

func enablePAC(pacURL string) {
	log.Printf("Serving PAC file at %v", pacURL)
	err := pac.On("http://"+pacURL+"/my.pac")
	if err != nil {
		log.Printf("Unable to set system proxy: %s", err)
	}
	atomic.StoreInt32(&isPACOn, 1)
}

func disablePAC(pacURL string) {
	if atomic.CompareAndSwapInt32(&isPACOn, 1, 0) {
		err := pac.Off(pacURL)
		if err != nil {
			log.Printf("Unable to unset system proxy: %s", err)
		}
		log.Printf("Unset system proxy")
	}
}