# 第1章
插件式存储引擎，存储引擎是基于表的

* InnoDB
    - 面向 OLTP 应用
    - 支持裸设备
* MyISAM
    - 不支持事务
    - 缓冲池只缓冲索引文件，不缓存数据库文件
* NDB
    - 全部内存
* Memory
    - 内存
    - 使用哈希索引
* Archive
    - 只支持 insert 和 select

数据库与传统文件的区别在于支持事务
