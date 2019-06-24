package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

const (
	serviceName = "/sd/service/"
	host        = "localhost:2379"
)

var (
	m     sync.RWMutex
	hosts = []string{}
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{host},
		DialTimeout: time.Second,
	})
	if err != nil {
		log.Fatalf("failed to create client, err=%s", err.Error())
	}
	defer cli.Close()

	// init
	resp, err := cli.Get(context.Background(), serviceName, clientv3.WithPrefix())
	if err != nil {
		log.Fatalf("failed to init, err=%s", err.Error())
	}
	watchRevision := resp.Header.Revision
	for i := 0; i < len(resp.Kvs); i++ {
		hosts = append(hosts, string(resp.Kvs[i].Key))
	}

	go func() {
		tick := time.Tick(time.Second)

		for {
			select {
			case <-tick:
				func() {
					m.RLock()
					defer m.RUnlock()
					fmt.Printf("%v\n", hosts)
				}()
			}
		}
	}()

	// watch
	c := cli.Watch(context.Background(), serviceName, clientv3.WithPrefix(), clientv3.WithRev(watchRevision))
	for {
		select {
		case watchResp := <-c:
			updateHandler(watchResp.Events)
		}
	}
}

func updateHandler(events []*clientv3.Event) {
	for _, event := range events {
		switch event.Type {
		case mvccpb.PUT:
			putHandler(event.Kv)
		case mvccpb.DELETE:
			deleteHandler(event.Kv)
		}
	}
}

func putHandler(kv *mvccpb.KeyValue) {
	fmt.Printf("trigger put, %v\n", kv)
	m.Lock()
	defer m.Unlock()
	for i := 0; i < len(hosts); i++ {
		if hosts[i] == string(kv.Key) {
			return
		}
	}
	hosts = append(hosts, string(kv.Key))
}

func deleteHandler(kv *mvccpb.KeyValue) {
	fmt.Printf("trigger delete, %v\n", kv)
	m.Lock()
	defer m.Unlock()
	hs := hosts[:0]
	for i := 0; i < len(hosts); i++ {
		if hosts[i] != string(kv.Key) {
			hs = append(hs, hosts[i])
		}
	}
	hosts = hs
}
