介绍:

    natproxy用于2台服务器之间做端口映射代理，基于网关组件 https://github.com/funny/fastway 开发

使用场景:

    应用的开发版本部署在公司内部网络，公司外的客户想要连接到内部网络使用应用。一般可以通过连接公司vpn、公司路由做端口映射或者部署到外网实现这需求，
    但这样往往需要开很多vpn账号或者映射多个端口不方便管理，部署外网又不能跟内部版本同步。你可以选择natproxy来轻松实现

部署实例:

    假设有内网服务器server_in(192.168.1.110)和外网服务器:server_out(8.8.8.8) 2个服务器，server_in有个web服务开在80端口，需要让外部网络访问内网的web80端口服务器。

    配置如下：
    {
        "ServerAddr" : "8.8.8.8", //server addr
        "Port" : 20000,                 //port of gateway
        "GateModel" : "gc",             //value 'gs' or 'gc', 'gs':gateway with server, 'gc':gateway with client
        "AuthKey" : "keyofnatproxy",    //key
        "ReusePort" : false,            //Enable/Disable the reuseport feature.
        "MaxPacket" : 524288,           //Limit max packet size.
        "MemPoolType" : "atom",         //Type of memory pool ('sync', 'atom' or 'chan').
        "MemPoolFactor" : 2,            //Growth in chunk size in memory pool.
        "MemPoolMinChunk" : 64,         //Smallest chunk size in memory pool.
        "MemPoolMaxChunk" : 65536,      //Largest chunk size in memory pool.
        "MemPoolPageSize" : 1048576,    //Size of each slab in memory pool.
        "ClientMaxConn" : 40960,           //Limit max virtual connections for each client.
        "ClientBufferSize" : 2048,      //Setting bufio.Reader's buffer size.
        "ClientSendChanSize" : 1024,    //Tunning client session's async behavior.
        "Proxys" : [                    //proxy list
            {
                "ClientPort" : 80,           //client port
                "ServerProxy" : "192.168.1.110:80" //proxy target
            }
        ]
    }

    1.在server_out服务器开启natproxy客户端服务
        $ ./natproxy -m=client

    2.在server_in服务器开启natproxy服务端服务
        $ ./natproxy -m=server
    
    3.启动后通过http://8.8.8.8:80访问到内网http://192.168.1.110:80