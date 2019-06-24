## api-gateway
- 认证
- 请求校验
- 路由请求
- 节流/限流
- 协议转换

#### pre-handling middleware
- cors
- validate
- auth
- user agent filter
- request id, tracing

#### post-handling middleware
- log
- stats

### 一些特性
- endpoint 级的配置支持覆盖 service 级的配置
- stg 自动更新


### 集成
- grpc
- http
- message queue
- stream

### 代码
迷之分布式限流器：quotas / count of instances
需要假设流量分布比较均匀

## TCP 长连接
- 实时发送位置
- 推送订单
- 推送消息

### code
- 为每一个 tcp 连接生成 uuid，并且将 instance 地址与 uuid 注册到 redis 上
- 业务上保证每建立一个新的 tcp 连接，首先必须有登录流程。登录服务从 redis 中获取 uuid 对应的 instance 地址，向该 instance 发送请求，关联 uuid 与用户 id（更新 redis）
- 客户端从 tcp 连接发送消息会推送到队列中
- 消费者消费队列更新状态
- 主动推送的消息从 redis 中根据用户 id 获得 instance 地址，确保发送到 tcp 连接所在的 instance 上
