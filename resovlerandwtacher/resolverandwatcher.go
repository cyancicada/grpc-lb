package resovlerandwtacher

import (
	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/naming"
)

type (
	ResolverConf struct {
		Key       string
		Endpoints []string
	}
	Resolver struct {
		config  *ResolverConf
		client3 *clientv3.Client
	}
	Watcher struct {
		resolver *Resolver
		init     bool
	}
)

func NewResolver(conf *ResolverConf) (*Resolver, error) {
	cli, err := clientv3.New(
		clientv3.Config{
			Endpoints: conf.Endpoints,
		})
	if nil != err {
		return nil, err
	}
	return &Resolver{client3: cli, config: conf}, nil
}
func (r *Resolver) Resolve(target string) (naming.Watcher, error) {

	return NewWatcher(r), nil
}

func NewWatcher(resolver *Resolver) *Watcher {
	return &Watcher{resolver: resolver}
}

func (w *Watcher) Next() ([]*naming.Update, error) {

	// next()
	// next()
	// 没有初始化情况
	if !w.init {
		w.init = true
		lgs, err := w.resolver.client3.Get(context.TODO(), w.resolver.config.Key, clientv3.WithPrefix())
		if nil == err && lgs != nil && lgs.Kvs != nil {
			address := make([]*naming.Update, 0)
			for _, kv := range lgs.Kvs {
				if kv.Value != nil {
					address = append(address,
						&naming.Update{Op: naming.Add, Addr: string(kv.Value)})
				}
			}
			return address, nil
		}
	}
	wc := w.resolver.client3.Watch(context.TODO(), w.resolver.config.Key, clientv3.WithPrefix())
	{
		for w := range wc {
			for _, we := range w.Events {
				switch we.Type {
				case mvccpb.PUT:
					return []*naming.Update{{Op: naming.Add, Addr: string(we.Kv.Value)}}, nil
				case mvccpb.DELETE:
					return []*naming.Update{{Op: naming.Delete, Addr: string(we.Kv.Value)}}, nil
				}
			}
		}
	}
	return nil, nil
}

func (w *Watcher) Close() {

	return
}
