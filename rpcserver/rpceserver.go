package rpcserver

import (
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"grpclb/register"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type (
	RpcServiceFunc func(server *grpc.Server)
	RpcService struct {
		register       *register.Register
		rpcServiceFunc RpcServiceFunc
	}
	RpcServiceConf struct {
		Key           string
		ServerAddress string
		Endpoints     []string
	}
)

func NewRpcServer(conf *RpcServiceConf, rpcServiceFunc RpcServiceFunc) (*RpcService, error) {

	client3, err := clientv3.New(
		clientv3.Config{
			Endpoints: conf.Endpoints,
		})
	if nil != err {
		return nil, err
	}
	return &RpcService{
		register:       register.NewRegister(conf.Key, client3, conf.ServerAddress),
		rpcServiceFunc: rpcServiceFunc,
	}, nil
}

func (s *RpcService) Run(serverOptions ...grpc.ServerOption) error {
	listen, err := net.Listen("tcp", s.register.GetServerAddress())
	if nil != err {
		return err
	}
	log.Printf("Rpc server listen at %s", s.register.GetServerAddress())
	s.register.Reg()
	server := grpc.NewServer(serverOptions...)
	s.rpcServiceFunc(server)
	s.deadNotify()
	if err := server.Serve(listen); nil != err {
		return err
	}
	return nil

}

func (s *RpcService) deadNotify() error {
	ch := make(chan os.Signal, 1) //
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		log.Printf("signal.Notify %v", <-ch)
		s.register.UnReg()
		os.Exit(1) //
	}()
	return nil
}
