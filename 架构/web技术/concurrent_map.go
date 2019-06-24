package main

import (
	"fmt"
	"hash/fnv"
	"math"
	"strconv"
	"sync"
)

const count uint32 = 8

func main() {
	m := newConcurrentMap()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Set(strconv.Itoa(i), i)
		}(i)
	}
	wg.Wait()

	for i := 0; i < 100; i++ {
		if !isPrime(i) {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				m.Del(strconv.Itoa(i))
			}(i)
		}
	}
	wg.Wait()

	fmt.Println(m.Count())
	fmt.Println(m)
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n == 2 {
		return true
	}

	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func (cm concurrentMap) String() string {
	s := "\n"
	for _, shard := range cm {
		func() {
			shard.RLock()
			defer shard.RUnlock()
			s += fmt.Sprintln(shard.m)
		}()
	}
	return s
}

type concurrentMap []*concurrentMapShared

type concurrentMapShared struct {
	sync.RWMutex
	m map[string]interface{}
}

func newConcurrentMap() concurrentMap {
	m := make(concurrentMap, count)
	for i := 0; i < len(m); i++ {
		m[i] = &concurrentMapShared{
			m: map[string]interface{}{},
		}
	}
	return m
}

func (cm concurrentMap) getShard(key string) *concurrentMapShared {
	hasher := fnv.New32()
	_, err := hasher.Write([]byte(key))
	if err != nil {
		return nil
	}
	return cm[hasher.Sum32()%count]
}

func (cm concurrentMap) Set(key string, val interface{}) {
	shard := cm.getShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.m[key] = val
}

func (cm concurrentMap) Get(key string) (interface{}, bool) {
	shard := cm.getShard(key)
	shard.RLock()
	defer shard.RUnlock()
	val, ok := shard.m[key]
	return val, ok
}

func (cm concurrentMap) Del(key string) {
	shard := cm.getShard(key)
	shard.Lock()
	defer shard.Unlock()
	delete(shard.m, key)
}

func (cm concurrentMap) Count() int {
	count := 0
	for _, shard := range cm {
		func() {
			shard.RLock()
			defer shard.RUnlock()
			count += len(shard.m)
		}()
	}
	return count
}
