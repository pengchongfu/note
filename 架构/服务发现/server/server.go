package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go.etcd.io/etcd/clientv3"
)

const (
	serviceName = "/sd/service/"
	host        = "localhost:2379"
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

	// grant
	lease, err := cli.Grant(context.Background(), 5)
	if err != nil {
		log.Fatalf("failed to get lease, err=%s", err.Error())
	}
	leaseID := lease.ID

	// put
	rand.Seed(time.Now().Unix())
	id := strconv.Itoa(rand.Intn(10000))
	host := serviceName + id
	_, err = cli.Put(context.Background(), host, id, clientv3.WithLease(leaseID))
	if err != nil {
		log.Fatalf("failed to put, err=%s", err.Error())
	}

	// keep-alive lease
	keepAlive, err := cli.KeepAlive(context.Background(), leaseID)
	if err != nil {
		log.Fatalf("failed to start keep alive, err=%s", err.Error())
	}
	go func() {
		for {
			select {
			case _, ok := <-keepAlive:
				if !ok {
					log.Fatalf("error")
				}
			}
		}
	}()

	// exit
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// revoke
	_, err = cli.Revoke(context.Background(), leaseID)
	if err != nil {
		log.Fatalf("failed to revoke lease, err=%s", err.Error())
	}
}
