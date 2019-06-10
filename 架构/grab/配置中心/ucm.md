UCM

同时支持多环境
模板：结构化，支持配置继承

模板：go yaml模式
值：toml
输出：json

模板支持自定义函数


通过 S3 提高读性能
- commit 到 git，git 改变生效之后，触发 lambda 函数更新 S3
- git 每个服务一个仓库
- 修改之后未 approve 的文件也已经提交到仓库中，只是作为一个 swap 文件，所以不能同时对一个文件修改

部署配置+应用配置

- 持久性
- 版本控制
- 权限
- 审计日志
- 多环境（参数化）
- 读吞吐量高，写吞吐量低

etcd/kafka 的作用？
- 高吞吐量的写接口
    - 顺序一致性
    - 持久性
    - 可以回滚
- 为什么要做写接口
    - 接入了 ab test 平台之后，很多写请求同时发生，而 git 不能并发处理

etcd 实现冲突检测和顺序一致性

将更改推送的队列中（kafka，并且保证和队列中的更改不会有冲突），通过 service 来做分区

两阶段提交和 fencing

只保证最终一致性（不能保证实时）

write sdk
- 初始化时，根据需要访问的 service 获取对应负责处理的 UCM 机器的 ip（通过 etcd 实现）
- sdk 发送请求时会带上版本号
- 冲突检测，版本号只能往上


代码：
- UCM sdk
    - 轮询获取最新配置
    - 直接调用 daemon
- UCM daemon
    - 支持 unix sock
- UCM service
    - kafka 和 etcd
    - 内存 git
    - 每次保存生成 swap 文件，触发 lambda 函数，更新 UCM 实例中的内存 git 仓库（如果两个请求同时更新一个文件，后一个在执行 git 操作时会失败）
    - approve 操作是一个 git 的 move 操作
- UCM write
    - fence 用于将各个 service 分配给不同的 UCM 实例
    - 流
        - 消费者：读取更改缓存到本地，再 flush 到 git
        - 生产者：
            - 读取 etcd 中的最新版本号进行冲突检测
            - 根据 service name 保存到单独的 partition，用来保证顺序一致性
    - write sdk
        - 由 sdk 轮询获取 service 对应的 UMC 实例地址