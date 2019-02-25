# ssh 免密登录

## 背景

> 使用 ssh-keygen 和 ssh_copy_id 两个命令来达成免密登录远程机器

## 详细设计

**本机操作**

```shell
cd ~
ssh-keygen -t rsa
# 此处可以更改 pub 文件名，我改为 id_rsa_login，命令要改为 ssh-copy-id -i ~/id_rsa_login.pub ....
ssh-copy-id -i ~/.ssh/id_rsa.pub -p 57522 root@10.22.51.94
#然后输入一次密码，之后就不用输入密码了，直接输入下面命令就可以登录 10.22.51.94
ssh -p 57522 root@10.22.51.94
```

