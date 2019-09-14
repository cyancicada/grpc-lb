package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"time"
)

func main() {


}

func Demo() {
	endpoints := []string{
		"127.0.0.1:2379",
	}
	//v3
	client3, err := clientv3.New(
		clientv3.Config{
			Endpoints: endpoints,
		})
	if nil != err {
		panic(err)
	}
	//con success
	//ps, err := client3.Put(context.TODO(), "mine", "vector teacher")
	//if nil != err {
	//	panic(err)
	//}
	//fmt.Println(ps)

	//ps, err := client3.Get(context.TODO(), "foo",clientv3.WithPrefix())
	//if nil != err {
	//	panic(err)
	//}
	//for _,kv := range ps.Kvs {
	//	log.Println(string(kv.Key),string(kv.Value))
	//}

	//wc :=  client3.Watch(context.TODO(),"mine")
	//for w := range wc {
	//	for _,we := range w.Events {
	//		fmt.Println(we.Type,string(we.Kv.Key),string(we.Kv.Value))
	//	}
	//}

	//合同 => lease id
	lgs, err := client3.Grant(context.TODO(), 60)
	if nil != err {
		panic(err)
	}
	if _, err := client3.Put(context.TODO(), "mine_lease", "my_room", clientv3.WithLease(lgs.ID)); nil != err {
		panic(err)
	}
	fmt.Println("ok")
	i := 0
	for {
		ps, err := client3.Get(context.TODO(), "mine_lease")
		if nil != err {
			panic(err)
		}
		for _, kv := range ps.Kvs {
			log.Println(string(kv.Key), string(kv.Value))
		}
		i += 5
		fmt.Printf("%d s", i)
		time.Sleep(5 * time.Second)
	}
}
