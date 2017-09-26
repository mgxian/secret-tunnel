## go语言实现的基于TCP的加密隧道


secret-tunnel是一个透明的TCP加密隧道，包含客户端与服务端。
secret-tunnel会把发送到secret-tunnel客户端的数据加密之后发送给secret-tunnel服务端，
secret-tunnel服务端接受到数据之后解密，然后发送给后端真实的用户服务。

数据如下图方式流动：
```
------------     加密    ------------------------    解密    ------------
|           |---------->|                       |---------> |           |
|    A      |           |    secret-tunnel      |           |     B     |
|           |<----------|                       |<----------|           |
------------     解密    ------------------------    加密    ------------
```

这样AB之间的通信就是加密的，不需要再做加密处理。

## 安装
```sh
go get github.com/maogx8/secret-tunnel
```

## 帮助
```
secret-tunnel --help
```

## 启动客户端
```sh
secret-tunnel -p 1027 -k maogx8 -s 192.168.12.211:1026
```

## 启动服务端
```sh
secret-tunnel -server -p 1026 -k maogx8 -r 192.168.12.211:1025
```