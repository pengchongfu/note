# 常用例子
---

## 作用域，编译时属性
```
package main

import "fmt"

var s = "undefined"

func main() {
	m := map[string]string{
		"": "defined",
	}
	s, ok := m[""]
	if ok {
		fmt.Println(s)
	}
	fmt.Println(getGlobalS())
}

func getGlobalS() string {
	return s
}
```
`defined
undefined`

## unicode 字符串
```
package main

import "fmt"

func main() {
	world := "ah,世界はやばい"
	var c, rC int

	for i := 0; i < len(world); i++ {
		c++
	}
	for _ = range world {
		rC++
	}

	fmt.Printf("c=%d, rC=%d\n", c, rC)
}
```
`c=21, rC=9`

## 数组，编译时确定长度
```
package main

import "fmt"

func main() {
	a := [...]int{10: 1}
	changeA(a)
	fmt.Printf("%v\n", a)
}

func changeA(a [11]int) {
	a[0] = 100
}
```
`[0 0 0 0 0 0 0 0 0 0 1]`

## slice, 两倍扩容
```
package main

import "fmt"

func main() {
	a := []int{1, 2}
	b := append(a, 3)

	c := append(b, 4)
	d := append(b, 5)

	fmt.Printf("c[3]=%v, d[3]=%v\n", c[3], d[3])
}
```
`c[3]=5, d[3]=5`
```
package main

import "fmt"

func main() {
	a := []int{1}
	b := append(a, 2)

	c := append(b, 3)
	d := append(b, 4)

	fmt.Printf("c[2]=%v, d[2]=%v\n", c[2], d[2])
}
```
`c[2]=3, d[2]=4`

## map panic
```
package main

import "fmt"

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	var m map[int]bool
	m[0] = true
}
```
`assignment to entry in nil map`

## struct, 匿名成员
```
package main

import "fmt"

type a struct {
	aa string
	b  b
}

type b struct {
	bb string
}

func main() {
	var i = a{}
	fmt.Println(i.bb)
}
```
`panic`

## defer stack, panic
```
package main

import "fmt"

func main() {
	fmt.Println(d())
}

func d() (res int) {
	p := &res

	defer func() {
		*p = 3
	}()
	defer func() {
		res = 2
	}()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover 1")
		}
	}()
	defer func() {
		panic("panic 1")
	}()

	return 1
}
```
`
recover 1
3
`

## 匿名成员可以直接调用方法
```
package main

import "fmt"

type A interface {
	name() string
}

func getName(a A) string {
	return a.name()
}

type M struct {
	N
}

type N struct{}

func (n N) name() string {
	return "n"
}

func main() {
	m := M{}
	fmt.Println(m.name())
	fmt.Println(getName(m))
}
```

## 类型断言：具体类型和接口
```
package main

import (
	"fmt"
	"strconv"
)

type A struct {
	a int
}

func (a A) String() string {
	return strconv.Itoa(a.a)
}

func main() {
	a := &A{a: 10}
	var b interface{}
	b = a
	if v, ok := b.(fmt.Stringer); ok {
		fmt.Println(v.String())
	}
}
```
`10`

## 并发控制
```
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	tokens := make(chan struct{}, 3)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			tokens <- struct{}{}
			defer func() {
				<-tokens
			}()
			time.Sleep(time.Duration(id)*time.Millisecond + time.Second)
			fmt.Printf("task: %d\n", id)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
```
```
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	tasks := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			tasks <- i
		}
		close(tasks)
	}()

	wg := sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			for id := range tasks {
				time.Sleep(time.Duration(id)*time.Millisecond + time.Second)
				fmt.Printf("task: %d\n", id)
			}
			wg.Done()
			return
		}()
	}
	wg.Wait()
}
```

## 类型分支
```
package main

import "fmt"

func main() {
	b := true
	var i interface{}
	i = b
	switch t := i.(type) {
	case int:
		fmt.Println("int")
	case bool:
		if t {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}
	}
}
```
`true`

## select
- case 里只能是 chan 相关操作
- 判断完case都会阻塞之后，才执行default
- 没有 default 的话阻塞
```
package main

import "fmt"

func main() {
	c := make(chan int, 1)
	for i := 0; i < 10; i++ {
		select {
		case x := <-c:
			fmt.Println(x)
		case c <- i:
		default:
		}
	}
}
```
`0
2
4
6
8`

## 轮流打印 a b c
```
package main

import (
	"fmt"
	"time"
)

func main() {
	turns := 3
	outputs := []string{"a", "b", "c"}
	chans := make([]chan struct{}, len(outputs))
	for i := 0; i < len(chans); i++ {
		chans[i] = make(chan struct{})
	}

	done := make(chan struct{}, 1)
	for i := 0; i < len(outputs); i++ {
		go func(id int) {
			for i := 0; i < turns; i++ {
				<-chans[id]
				time.Sleep(time.Second)
				fmt.Printf("routine: %d, turn: %d, %s\n", id, i, outputs[id])

				if i == turns-1 && id == len(outputs)-1 {
					done <- struct{}{}
					break
				}
				next := (id + 1) % len(outputs)
				chans[next] <- struct{}{}
			}
		}(i)
	}
	chans[0] <- struct{}{}
	<-done
	fmt.Println("done")
}
```
```
routine: 0, turn: 0, a
routine: 1, turn: 0, b
routine: 2, turn: 0, c
routine: 0, turn: 1, a
routine: 1, turn: 1, b
routine: 2, turn: 1, c
routine: 0, turn: 2, a
routine: 1, turn: 2, b
routine: 2, turn: 2, c
done
```

## 实现 sync.Once
```
package main

import (
	"fmt"
	"sync"
)

type once struct {
	mutex sync.RWMutex
	done  bool
}

func (o *once) Do(f func()) {
	o.mutex.RLock()
	if o.done {
		o.mutex.RUnlock()
		return
	}
	o.mutex.RUnlock()

	o.mutex.Lock()
	if !o.done {
		f()
		o.done = true
	}
	o.mutex.Unlock()
}

func main() {
	o := once{}

	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			o.Do(func() {
				fmt.Println("print once")
			})
			wg.Done()
		}()
	}
	wg.Wait()
}
```

优化版，不使用读锁
```
type once struct {
	sync.Mutex
	done uint32
}

func (o *once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}

	o.Lock()
	defer o.Unlock()
	if o.done == 0 {
		f()
		atomic.StoreUint32(&o.done, 1)
	}
}
```

## golang 将返回值放在内存栈中，实现多值返回
```
package main

import "fmt"

func main() {
	stack()
}

func stack() (a, b int) {
	fmt.Printf("%p\n", &a)
	fmt.Printf("%p\n", &b)
	var c int
	fmt.Printf("%p\n", &c)
	return
}
```
`0xc0000cc008
0xc0000cc010
0xc0000cc018`

## 方法调用只是函数调用的语法糖
- 方法编译时会决定绑定的实例
- 有点类似闭包的结构体实现
```
package main

import "fmt"

type T struct {
	name string
}

func (t T) getName() string {
	return t.name
}

func main() {
	f1 := T.getName
	t := T{}
	t.name = "t1"
	f2 := t.getName
	fmt.Println(f1(t))
	t.name = "t2"
	fmt.Println(f2())
	fmt.Println(t.getName())
}
```
`t1
t1
t2`

## panic 之后会改变执行流程，执行 defer；向调用方冒泡，直到 recover
```
package main

import "fmt"

func main() {
	fmt.Println(d1())
}

func d1() int {
	defer fmt.Println("d1 done")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("d1 panic")
		}
	}()
	fmt.Print(d2())
	defer fmt.Println("d1 return")
	return 100
}

func d2() int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("d2 panic")
			panic(r)
		}
	}()
	defer fmt.Println("d2 return")

	panic("error")
}
```
`d2 return
d2 panic
d1 panic
d1 done
0`

## slice 一个 array 指向同一内存地址
```
package main

import "fmt"

func main() {
	a := [...]int{3: 10, 0: 9}

	b := a[1:]
	b[0] = 33

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(&b[0] == &a[1])

	fmt.Printf("%p\n", b)
	fmt.Printf("%p\n", &a)
}
```
`[9 33 0 10]
[33 0 10]
true
0xc0000fe008
0xc0000fe000`

## in-place filter
```
package main

import "fmt"

func main() {
	nums := []int{1, 3, -3, 8, 9, 0, 1}
	fmt.Printf("%p\n", nums)
	nums = filter(nums)
	fmt.Println(nums)
	fmt.Printf("%p\n", nums)
}

func filter(nums []int) []int {
	res := nums[:0]
	for _, v := range nums {
		if v > 0 {
			res = append(res, v)
		}
	}
	return res
}
```
`0xc0000f8000
[1 3 8 9 1]
0xc0000f8000`


## recover 和有异常的栈帧必须只隔一个栈帧（刚好跨越 defer 函数，祖父级）
```
package main

import "fmt"

func main() {
	defer func() {
		func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
	}()

	panic("error")
}
```
```panic: error

goroutine 1 [running]:
main.main()
        /Users/chongfu.peng/Desktop/note/golang/main.go:14 +0x63
exit status 2
```

## select 随机数发生器
```
package main

import "fmt"

func main() {
	c := make(chan int, 1)

	for i := 0; i < 10; i++ {
		select {
		case c <- 0:
		case c <- 1:
		}
		v := <-c
		fmt.Println(v)
	}
}
```