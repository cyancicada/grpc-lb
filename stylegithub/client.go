package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	"grpclb/rpcfile"
	"time"

	"google.golang.org/grpc"
)

func main() {
	endpoints := []string{"127.0.0.1:2379"}
	cli, err := clientv3.New(
		clientv3.Config{
			Endpoints: endpoints,
		})
	if nil != err {
		panic(err)
	}

	r := &etcdnaming.GRPCResolver{Client: cli}
	b := grpc.RoundRobin(r)
	conn, err := grpc.DialContext(
		context.TODO(),
		"my-service",
		grpc.WithInsecure(),
		grpc.WithBalancer(b),
		grpc.WithBlock(),
	)

	t := time.NewTicker(2 * time.Second)
	for tc := range t.C {
		cli := demo.NewDemoServiceClient(conn)
		res, err := cli.DemoHandler(context.TODO(), &demo.DemoRequest{Name: "vector"})
		if nil != err {
			panic(err)
		}
		fmt.Println("I has get result ", res.Name, " time is [", tc, "]")
	}
}
