package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"grpclb/rpcfile"
	"grpclb/rpcserver"
	"grpclb/rpcserverimpl"
)

var (
	port = flag.Int("p", 50001, "server port")
)

const (
	key string = "vector_rpc_server"
)

func main() {
	flag.Parse()
	conf := &rpcserver.RpcServiceConf{
		Key:           key,
		ServerAddress: fmt.Sprintf("127.0.0.1:%d", *port),
		Endpoints:     []string{"127.0.0.1:2379"},
	}
	if server, err := rpcserver.NewRpcServer(conf, func(server *grpc.Server) {
		demo.RegisterDemoServiceServer(server,
			&rpcserverimpl.DemoServiceServer{ServerAddress: conf.ServerAddress})
	}); nil == err {
		if err := server.Run(); nil != err {
			panic(err)
		}
	}
}
