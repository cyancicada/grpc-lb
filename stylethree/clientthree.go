package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpclb/resovlerandwtacher"
	"grpclb/rpcfile"
	"strings"
	"time"
)

const (
	key string = "vector_rpc_server"
)

func main() {
	endpoints := []string{"127.0.0.1:2379"}
	conf := &resovlerandwtacher.ResolverConf{Key: key, Endpoints: endpoints}
	r, err := resovlerandwtacher.NewResolver(conf)
	if nil != err {
		panic(err)
	}
	b := grpc.RoundRobin(r)

	conn, err := grpc.DialContext(
		context.TODO(),
		strings.Join(endpoints, ","),
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
		fmt.Println("I has get result ", res.Name," time is [",tc,"]")
	}
}
