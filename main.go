package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	fastway "github.com/funny/fastway/go"
	"github.com/funny/slab"
)

var (
	gModel   *string = flag.String("Model", "client", "client port")
	gConfig  *Config
	gClients []*Client
	gServers []*Server
)

const (
	modelClient string = "client"
	modelServer string = "server"
)

func init() {
	flag.Parse()
}

func main() {
	b, e := ioutil.ReadFile("conf.json")
	if e != nil {
		panic(e.Error())
	}

	gConfig = new(Config)

	if e = json.Unmarshal(b, gConfig); e != nil {
		panic(e.Error())
	}

	if *gModel == "client" {
		startClient()
	} else {
		startServer()
	}

	sigTERM := make(chan os.Signal, 1)

	signal.Notify(sigTERM, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

L:
	for {
		select {
		case <-sigTERM:
			if *gModel == "client" {
				stopClient()
			} else {
				stopServer()
			}
			break L
		}
	}
}

func startClient() {
	addr := fmt.Sprintf("%s:%d", gConfig.GateWayAddr, gConfig.GateWayClientPort)

	endPoint, e := fastway.DialClient("tcp", addr, fastway.EndPointCfg{
		MemPool:      slab.NewSyncPool(64, 4*1024, 64),
		MaxPacket:    128 * 1024,
		PingInterval: 60 * time.Second,
		PingTimeout:  5 * time.Second,
		SendChanSize: 512,
		RecvChanSize: 512,
		MsgFormat:    NewProtocolFormat(),
		TimeoutCallback: func() bool {
			log.Println("fastway disconnected")
			return false
		},
	})

	if e != nil {
		panic(e.Error())
	}

	for serverID, proxy := range gConfig.Proxys {
		l, e := net.Listen("tcp", fmt.Sprintf(":%d", proxy.Port))
		if e != nil {
			panic(e.Error())
		}
		c := NewClient(l, endPoint, uint32(serverID+1))
		gClients = append(gClients, c)
		log.Println(fmt.Sprintf("start proxy port %d to server %d", proxy.Port, serverID+1))
		go c.Serv()
	}
}

func stopClient() {
	for _, c := range gClients {
		c.Stop()
	}
}

func startServer() {
	for serverID, proxy := range gConfig.Proxys {
		endPoint, e := fastway.DialServer("tcp", fmt.Sprintf("%s:%d", gConfig.GateWayAddr, gConfig.GateWayServerPort), fastway.EndPointCfg{
			ServerID:     uint32(serverID + 1),
			MemPool:      slab.NewSyncPool(64, 4*1024, 64),
			MaxPacket:    128 * 1024,
			PingInterval: 60 * time.Second,
			PingTimeout:  5 * time.Second,
			AuthKey:      gConfig.GateWayAuthKey,
			SendChanSize: 100000,
			RecvChanSize: 100000,
			MsgFormat:    NewProtocolFormat(),
			TimeoutCallback: func() bool {
				log.Println("fastway disconnected")
				return false
			},
		})

		if e != nil {
			panic(e.Error())
		}

		s := NewServer(proxy.Proxy, endPoint)
		gServers = append(gServers, s)
		log.Println(fmt.Sprintf("start proxy server %d (%s)", serverID+1, proxy.Proxy))
		go s.Serv()
	}
}

func stopServer() {
	for _, s := range gServers {
		s.Stop()
	}
}
