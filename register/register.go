package register

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"log"
	"time"
)

type (
	Register struct {
		key           string
		client3       *clientv3.Client
		serverAddress string
		stop          chan bool
		interval      time.Duration
		leaseTime     int64
	}
)

func NewRegister(
	key string,
	client3 *clientv3.Client,
	serverAddress string,
) *Register {
	return &Register{
		key:           key,
		serverAddress: serverAddress,
		client3:       client3,
		interval:      10 * time.Second,
		leaseTime:     13,
		stop:          make(chan bool, 1),
	}
}

func (r *Register) Reg() {
	k := r.makeKey()
	go func() {
		t := time.NewTicker(r.interval)
		for {
			lgs, err := r.client3.Grant(context.TODO(), r.leaseTime)
			if nil != err {
				panic(err)
			}
			//get key
			if _, err := r.client3.Get(context.TODO(), k); nil != err {
				if err == rpctypes.ErrKeyNotFound {
					if _, err := r.client3.Put(context.TODO(),
						k, r.serverAddress,
						clientv3.WithLease(lgs.ID));
						nil != err {
						panic(err)
					}
				} else {
					panic(err)
				}
			} else {
				if _, err := r.client3.Put(context.TODO(),
					k,
					r.serverAddress,
					clientv3.WithLease(lgs.ID)); nil != err {
					panic(err)
				}
			}
			select {
			case ttl := <-t.C:
				log.Println(ttl)
			case <-r.stop:
				return
			}
		}
	}()
}

func (r *Register) makeKey() string {

	return fmt.Sprintf("%s_%s", r.key, r.serverAddress)
}
func (r *Register) GetServerAddress() string {
	return r.serverAddress
}
func (r *Register) UnReg() {
	r.stop <- true
	k := r.makeKey()
	r.stop = make(chan bool, 1) // 为了不让在多纯种下同现deadlock 情况
	if _, err := r.client3.Delete(context.TODO(), k); nil != err {
		panic(err)
	} else {
		log.Printf("%s UnReg Sucess", k)
	}
	return

}
