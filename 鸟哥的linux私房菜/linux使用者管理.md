sudo 只需要自己的密码，需要在/etc/sudoers中
（通过支持wheels群组执行所有命令）

crontab 周期性工作

& 在背景执行
jobs 查看所有进行中的工作
fg 取回到前景
bg
kill 默认值为-15，中断。-9强制关闭

以上的都是bash的背景，远程连接关闭时依旧会中断
应该使用 at 或 nohup

/proc 进程信息