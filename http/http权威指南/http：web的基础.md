nc 比 telnet 功能更多，且支持 udp

url不支持二进制，需要转义

Content-length：实体的长度

- 一直发送小分组的话会影响性能，因此 Nagle 算法会等待数据到来或者ack到来
- tcp 延迟确认在没有数据要发送的情况下会等待40ms再发送ack
- 两者共同作用导致延迟


- 并行连接
- 持久连接：首部 Connection：keep-alive
- 管道化连接：在持久连接的基础上
    - 请求和响应的顺序必须一致
