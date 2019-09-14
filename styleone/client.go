package main

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"grpclb/rpcfile"
	"log"
	"math/rand"
	"sync"
	"time"
)

type (
	RoundRobinConf struct {
		Key       string
		Endpoints []string //etcd 集群地址 ｛"10.189.231,45:2379","10.189.231,46:2379","10.189.231,47:2379","..."｝
	}
	CreateRoundIndexFunc func() (interface{}, error)// 用来创建grpc.ClientConn
	RoundRobin struct {  //LB 负载均衡策略
		index   int  // 用某一个 er ｛"10.189.231,45:2379","10.189.231,46:2379","10.189.231,47:2379"｝
		lock    sync.Mutex
		targets []interface{}  // ｛grpc.DialContext1(1),grpc.DialContext(2),grpc.DialContext(3)｝
	}
)

//  n => 有多少个rpc服务
//createRoundIndexFunc  用来返回grpcClient
func NewRoundRobin(n int, createRoundIndexFunc CreateRoundIndexFunc) (*RoundRobin, error) {
	//多个rpc服务grpc.Server
	targets := make([]interface{}, n)
	var err error
	for i := 0; i < n; i++ {
		targets[i], err = createRoundIndexFunc()
		if err != nil {
			return nil, err
		}
	}
	//1
	//获取随机数，不加随机种子，每次遍历获取都是重复的一些随机数据
	//
	//rand.Seed(time.Now().UnixNano())
	//1
	//设置随机数种子，加上这行代码，可以保证每次随机都是随机的
	//使用给定的seed将默认资源初始化到一个确定的状态；如未调用Seed，默认资源的行为就好像调用了Seed(1)
	rand.Seed(time.Now().UnixNano()) // 随机因子
	return &RoundRobin{
		targets: targets,
		index:   rand.Intn(n),
	}, nil
}

func (r *RoundRobin) Next() interface{} {
	r.lock.Lock()
	defer r.lock.Unlock()
	// index = 0 +1 = 2
	//2 % 3 = 1
	// 4 %3 = 1
	//0 %3 = 0
	r.index = (r.index + 1) % len(r.targets)
	return r.targets[r.index]
}

const (
	key string = "vector_rpc_server"
)

func main() {

	// 上下文

	conf := &RoundRobinConf{
		Key:       key,
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
		roundRobin, err := NewRoundRobin(l, func() (interface{}, error) {
			l--
			return serverMap[l], nil
		})
		client := demo.NewDemoServiceClient(roundRobin.Next().(*grpc.ClientConn))
		res, err := client.DemoHandler(
			context.TODO(),
			&demo.DemoRequest{Name: "vector"})
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("SUCCESS ：", res.Name)
		select {
		case <-t.C:

		}
	}

}
