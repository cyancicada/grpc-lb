package main

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"google.golang.org/grpc"
	"grpclb/rpcfile"
	"log"
	"time"
)

type (
	RoundRobinConf struct {
		Key       string
		LbKey     string
		Endpoints []string //etcd 集群地址 ｛"10.189.231,45:2379","10.189.231,46:2379","10.189.231,47:2379","..."｝
	}
)

const (
	key   string = "vector_rpc_server"
	lbKey string = "rpc_server_lb_key"
)

func main() {
	conf := &RoundRobinConf{
		Key:       key,
		LbKey:     lbKey,
		Endpoints: []string{"127.0.0.1:2379"},
	}
	client3, err := clientv3.New(clientv3.Config{
		Endpoints: conf.Endpoints,
	})
	if nil != err {
		panic(err)
	}
	t := time.NewTicker(3 * time.Second)
	for {

		allServerList, err := client3.Get(context.TODO(), conf.Key, clientv3.WithPrefix())
		if nil != err {
			panic(err)
		}
		lgs, err := client3.Get(context.TODO(), conf.LbKey)
		if nil != err {
			if err == rpctypes.ErrKeyNotFound {
				if _, err := client3.Put(context.TODO(), conf.LbKey, "1"); nil != err {
					panic(err)
				}
			}
			panic(err)
		}
		serverMap := make(map[int]*grpc.ClientConn)
		for i, kv := range allServerList.Kvs {
			ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)
			conn, err := grpc.DialContext(ctx, string(kv.Value), grpc.WithInsecure())
			if nil != err {
				continue
			}
			serverMap[i] = conn
		}
		l := len(serverMap)
		if l > 0 {
			for _, kv := range lgs.Kvs {
				versionId := int(kv.Version)
				index := versionId % l
				client := demo.NewDemoServiceClient(serverMap[index])
				res, err := client.DemoHandler(
					context.TODO(),
					&demo.DemoRequest{Name: "vector"})
				if err != nil {
					log.Println(err)
					continue
				}
				log.Println("SUCCESS ：", res.Name)
			}
			if _, err := client3.Put(context.TODO(), conf.LbKey, "1"); nil != err {
				panic(err)
			}
		} else {
			log.Println("There no rpc server ,serverMap is ", l)
		}
		select {
		case <-t.C:

		}
	}

}
