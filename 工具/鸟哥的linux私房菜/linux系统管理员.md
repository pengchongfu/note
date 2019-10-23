服务开机启动：
- init
- systemd（通过systemctl来管理）

- rsyslog服务提供日志功能
- logratate 定时删除日志
（systemd 自带日志记录功能，默认只保留在内存中，通过journalctl来管理）


- bootloader可以转交控制权给其他的bootloader
- vmlinuz文件其实是内核文件
- initramfs是虚拟文件系统（包含了所需文件系统驱动）
- grub2：第一阶段，第二阶段