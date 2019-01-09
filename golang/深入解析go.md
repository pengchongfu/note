new 返回一个指向已清零内存的指针，make 返回一个复杂的结构

布尔类型是false，整型是0，字符串是""，而指针，函数，interface，slice，channel和map的零值都是nil

return 不是原子指令：先给返回值赋值，调用 defer 指令，最后返回到调用函数

返回闭包时其实返回了一个结构体

调度
- 结构体G
- 结构体M（关联到操作系统线程）
- 结构体P（处理器）

垃圾回收：
- 标记清除
- stop the world

channel
- 缓冲区紧跟着 channel 结构体分配
- recvq 和 sendq 两个链表

interface
- 实际上是一个两个成员的结构体：指针指向数据，类型