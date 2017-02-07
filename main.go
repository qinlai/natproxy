package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var gModel = flag.String("m", "model:server/client", "client")

const (
	ModelsClient string = "client"
	ModelsServer string = "server"
)

func init() {
	flag.Parse()
}

func main() {
	pid := syscall.Getpid()
	ioutil.WriteFile("natproxy.pid", []byte(strconv.Itoa(pid)), 0644)
	defer os.Remove("natproxy.pid")

	addr := fmt.Sprintf("%s:%d", gConf.Addr, gConf.Port)
	if *gModel == ModelsClient {
		if gConf.GateModel == GateModelsClient {
			gw, serverAddr := startFastway()
			defer gw.Stop()
			addr = serverAddr.String()
		}
		startClient(addr)
	} else {
		if gConf.GateModel == GateModelsServer {
			gw, serverAddr := startFastway()
			defer gw.Stop()
			addr = serverAddr.String()
		}
		startServer(addr)
	}

	sigTERM := make(chan os.Signal, 1)

	signal.Notify(sigTERM, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

L:
	for {
		select {
		case <-sigTERM:
			if *gModel == ModelsClient {
				stopClient()
			} else {
				stopServer()
			}
			break L
		}
	}
}
