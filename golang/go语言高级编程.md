- main.main 函数之前的所有代码都在一个 goroutine 中，如果在 init 函数中调用 go 关键字，与 main.main 并发执行
- 因为在一个线程，满足顺序一致性


## 分布式 id 生成器
- twitter 的 snowflake 算法
- 64 位，1位未使用，41位毫秒时间戳，5位数据中心，5位实例id，12位循环自增key

## 分布式锁
- redis setnx
- zookeeper 会阻塞，通过 watch 实现
- etcd：往某个路径下写入 key，也是分布式阻塞锁

etcd 和 zookeeper 没有伸缩性，无法水平扩展。
- 只能增加集群
- proxy 根据 id 分片；或者做 client 分片

## 延时任务系统
### 定时器
1. 用最小堆实现
2. 时间轮
### 任务分发


## 分布式搜索引擎
### elasticsearch
1. 基于时间戳增量同步
2. binlog 同步到 kafka

## 负载平衡
```
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	m1 := map[int]int{}
	m2 := map[int]int{}

	for i := 0; i < 100000; i++ {
		nodes := []int{1, 2, 3, 4, 5, 6, 7, 8}
		suffle1(nodes)
		m1[nodes[0]]++
	}
	for i := 0; i < 100000; i++ {
		nodes := []int{1, 2, 3, 4, 5, 6, 7, 8}
		suffle2(nodes)
		m2[nodes[0]]++
	}

	fmt.Println(m1)
	fmt.Println(m2)
}

func suffle1(nums []int) {
	for i := 0; i < len(nums); i++ {
		a := rand.Intn(len(nums))
		b := rand.Intn(len(nums))
		nums[a], nums[b] = nums[b], nums[a]
	}
}

func suffle2(nums []int) {
	for i := len(nums) - 1; i > 0; i-- {
		a := rand.Intn(i + 1)
		nums[i], nums[a] = nums[a], nums[i]
	}
}
```
`map[1:21196 5:11205 6:11298 3:11123 8:11290 2:11186 7:11256 4:11446]
map[7:12473 5:12599 1:12593 2:12387 3:12329 6:12540 8:12672 4:12407]`

## 分布式配置管理
### etcd
