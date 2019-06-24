package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const indexHTML = `
<!doctype html>
<html><head></head><body><h1>this is index page</h1></body></html>
`

var client = getCallClient()

type meta struct {
	name    string
	timeout int
}

type endpoint func(ctx context.Context, req interface{}) (response interface{}, err error)

type middleware func(meta, endpoint) endpoint

func main() {
	http.HandleFunc("/index", makeHandler(meta{name: "index", timeout: 2}, indexEndpoint))
	http.ListenAndServe(":8888", nil)
}

func makeHandler(m meta, e endpoint) func(w http.ResponseWriter, req *http.Request) {
	middlewares := []middleware{cacheMiddleware, durationMiddleware, timeoutMiddleware}
	handler := e
	for _, v := range middlewares {
		handler = v(m, handler)
	}

	return func(w http.ResponseWriter, req *http.Request) {
		url := req.URL.Path
		ctx := context.Background()
		resp, err := handler(ctx, url)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		fmt.Fprintf(w, resp.(string))

	}
}

func getCallClient() func() {
	count := int32(0)
	return func() {
		ok := atomic.CompareAndSwapInt32(&count, 0, 1)
		if ok {
			return
		}
		time.Sleep(5 * time.Second)
	}
}

func indexEndpoint(ctx context.Context, req interface{}) (interface{}, error) {
	return index(ctx, req.(string))
}

func index(ctx context.Context, req string) (string, error) {
	done := make(chan struct{})
	go func() {
		client()
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return "", errors.New("get from cache")
	}

	return indexHTML, nil
}

func durationMiddleware(m meta, e endpoint) endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		now := time.Now()
		resp, err := e(ctx, req)
		duration := time.Now().Sub(now)
		go func() {
			fmt.Printf("[duration middleware] endpoint=%s, duration=%v\n", m.name, duration)
		}()
		return resp, err
	}
}

func timeoutMiddleware(m meta, e endpoint) endpoint {
	timeout := time.Duration(m.timeout) * time.Second
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_, ok := ctx.Deadline()
		if ok {
			return e(ctx, req)
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		fmt.Printf("[timeout middleware] endpoint=%s, timeout=%d\n", m.name, timeout)
		return e(ctx, req)
	}
}

func cacheMiddleware(m meta, e endpoint) endpoint {
	cacheMap := struct {
		store map[string]interface{}
		sync.RWMutex
	}{
		store: map[string]interface{}{},
	}

	return func(ctx context.Context, req interface{}) (interface{}, error) {
		resp, err := e(ctx, req)
		if err != nil {
			if err.Error() == "get from cache" {
				cacheMap.RLock()
				defer cacheMap.RUnlock()
				fmt.Printf("[cache middleware] endpoint=%s, key=%v\n", m.name, req.(string))
				return cacheMap.store[req.(string)], nil
			}

			return nil, err
		}

		cacheMap.Lock()
		defer cacheMap.Unlock()
		cacheMap.store[req.(string)] = resp

		return resp, err
	}

}
