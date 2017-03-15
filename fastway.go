package main

import (
	"log"
	"net"
	"time"

	"fmt"

	fastway "github.com/funny/fastway/go"
	"github.com/funny/reuseport"
	"github.com/funny/slab"
	snet "github.com/funny/snet/go"
)

func startFastway() (*fastway.Gateway, net.Addr) {
	var pool slab.Pool
	switch gConf.MemPoolType {
	case "sync":
		pool = slab.NewSyncPool(gConf.MemPoolMinChunk, gConf.MemPoolMaxChunk, gConf.MemPoolFactor)
	case "atom":
		pool = slab.NewAtomPool(gConf.MemPoolMinChunk, gConf.MemPoolMaxChunk, gConf.MemPoolFactor, gConf.MemPoolPageSize)
	case "chan":
		pool = slab.NewChanPool(gConf.MemPoolMinChunk, gConf.MemPoolMaxChunk, gConf.MemPoolFactor, gConf.MemPoolPageSize)
	default:
		println(`unsupported memory pool type, must be "sync", "atom" or "chan"`)
	}

	gw := fastway.NewGateway(pool, gConf.MaxPacket)

	var serverAddr, clientAddr string
	if gConf.GateModel == GateModelsClient {
		serverAddr = fmt.Sprintf(":%d", gConf.Port)
		clientAddr = fmt.Sprintf("%s:0", gConf.Addr)
	} else {
		clientAddr = fmt.Sprintf(":%d", gConf.Port)
		serverAddr = fmt.Sprintf("%s:0", gConf.Addr)
	}

	clientL := listen("client", clientAddr, gConf.ReusePort,
		false,
		false,
		0,
		0,
		0,
	)

	go gw.ServeClients(
		clientL,
		fastway.GatewayCfg{
			MaxConn:      gConf.ClientMaxConn,
			BufferSize:   gConf.ClientBufferSize,
			SendChanSize: gConf.ClientSendChanSize,
			IdleTimeout:  29 * time.Minute,
		},
	)

	serverL := listen("server", serverAddr, gConf.ReusePort,
		false,
		false,
		0,
		0,
		0,
	)

	go gw.ServeServers(
		serverL,
		fastway.GatewayCfg{
			AuthKey:      gConf.AuthKey,
			BufferSize:   512 * 1024,
			SendChanSize: 512,
			IdleTimeout:  29 * time.Minute,
		},
	)

	var l net.Listener

	if gConf.GateModel == GateModelsClient {
		l = clientL
	} else {
		l = serverL
	}

	return gw, l.Addr()
}

func listen(who, addr string, reuse, snetEnable, snetEncrypt bool, snetBuffer int, snetInitTimeout, snetWaitTimeout time.Duration) net.Listener {
	var lsn net.Listener
	var err error

	if reuse {
		lsn, err = reuseport.NewReusablePortListener("tcp", addr)
	} else {
		lsn, err = net.Listen("tcp", addr)
	}

	if err != nil {
		log.Fatalf("setup %s listener at %s failed - %s", who, addr, err)
	}

	if snetEnable {
		lsn, _ = snet.Listen(snet.Config{
			EnableCrypt:        snetEncrypt,
			RewriterBufferSize: snetBuffer,
			HandshakeTimeout:   snetInitTimeout,
			ReconnWaitTimeout:  snetWaitTimeout,
		}, func() (net.Listener, error) {
			return lsn, nil
		})
	}

	return lsn
}
