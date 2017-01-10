介绍:

    natproxy用于2台服务器之间做端口映射代理，基于网关组件 https://github.com/funny/fastway 开发

使用场景:

    应用的开发版本部署在公司内部网络，公司外的客户想要连接到内部网络使用应用。一般可以通过连接公司vpn、公司路由做端口映射或者部署到外网实现这需求，
    但这样往往需要开很多vpn账号或者映射多个端口不方便管理，部署外网又不能跟内部版本同步。你可以选择natproxy来轻松实现

部署实例:

    假设有内网服务器server_in(192.168.1.110)和外网服务器:server_out(8.8.8.8) 2个服务器，server_in有个web服务开在80端口，需要让外部网络访问内网的web80端口服务器。

    配置如下：
    {
        "gate_way_addr" : "8.8.8.8",     //fastway地址
        "gate_way_client_port" : 20000,  //fastway客户端端口
        "gate_way_server_port" : 20001,  //fastway服务端端口
        "gate_way_auth_key" : "keyoffastway", //fastway秘钥
        "proxys" : [
            {
                "port" : 30000,
                "proxy" : "192.168.1.110:80"
            }
        ]
    }

    1.在server_out服务器开启fastway服务
        开启详情请查看：https://github.com/funny/fastway
    2.在server_in服务器开启natproxy服务端服务
        $ ./natproxy --Model=server
    3.在server_out服务器开启natproxy客户端服务
        $ ./natproxy --Model=client
    4.启动后通过http://8.8.8.8:80访问到内网http://192.168.1.110:80