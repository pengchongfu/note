## grab-kit 支持的功能

---
- 服务端客户端
- dto dao
- 性能监控、日志、tracing
- 可用性：断路器、重试机制、超时控制、限流、缓存、服务降级
- 数据库、缓存、队列

## 中间件
- meta 元数据中间件
- chaos 混沌实验
- quota 流量限额
- throttling 根据 cpu 负载进行节流

### 重试机制
为了提高 resilience
1. 非幂等接口
    - 只有部分错误重试（资源紧张）
2. 幂等接口
    - 部分明显重试无效的错误不重试（不存在，无权限）

### 缓存机制
客户端缓存请求结果，并且设置过期时间
- 保存在内存中，gob格式

## code
endpoint 函数声明
```
type Endpoint func(ctx context.Context, request interface{}) (metadata metadata.Metadata, response interface{}, err error)
```

middleware 函数声明
```
type Middleware func(Meta, Endpoint) Endpoint
```