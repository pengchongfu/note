- 账号信息：/etc/passwd
- 密码信息：/etc/shadow
- 组信息：/etc/group

修改文件权限：chgrp chown

chmod a+x filename
(u g o a) (+ -) (r w x)

umask 新建文件或目录时的初始权限
- 文件：默认无执行权限 666
- 目录：默认可以执行 777
- 扣除 umask 的值

- SUID:4 执行的时候具有所有者权限
- SGID:2 群组
- SBIT:1
（例子： chmod 4775 filename）


FAT：没有inode

inode
- 文件
    - 大小固定
    - 一个文件占用一个 inode
    - 使用 block 做间接指针
- 目录
    - 一个inode和至少一个block
    - block记录文件名和对应inode地址

block bitmap
- 标记 block 是否可用，删除文件时需要更新

硬链接
- 不能跨文件系统
- 不能链接目录


- 打包 tar
- 压缩 gzip

