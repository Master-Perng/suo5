# Suo5

`suo5` 是一个全新的 HTTP 代理隧道，基于 `HTTP/1.1` 的 `Chunked-Encoding`
构建。相比 [Neo-reGeorg](https://github.com/L-codes/Neo-reGeorg) 等传统隧道工具, `suo5` 的性能可以达到其数十倍。 其主要特性如下：

- 一条连接实现数据的双向发送和接收，性能堪比 TCP 直连
- 同时支持全双工与半双工模式，并可自动选择最佳的模式
- 支持在 Nginx 反向代理的场景中使用
- 自有数据序列化协议，数据经过加密传输
- 服务端基于 `Servlet` 原生实现，JDK6~JDK19 全版本兼容
- 同时提供提供命令行和图形化界面，方便不同用户使用

> 免责声明：此工具仅限于安全研究，用户承担因使用此工具而导致的所有法律和相关责任！作者不承担任何法律责任！

具体原理参考 [博客]()

## 安装运行

前往 Releases 下载编译好的二进制，其中带 `gui` 的版本是界面版，不带 `gui` 的为命令行版。

使用前需上传 `proxy.jsp` 到目标环境中并确保可以执行。

### 界面版

> 界面版基于 [wails](https://github.com/wailsapp/wails) 实现，依赖 Webview2 框架。Windows 11 和 MacOS 已自带该组件，其他系统会弹框请允许下载安装，否则无法使用。

![gui.png](tests/img/gui.png)

### 命令行

```text
NAME:
   suo5 - A super http proxy tunnel

USAGE:
   suo5 command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --target value, -t value  set the memshell url, ex: http://localhost:8080/tomcat_debug_war_exploded/
   --listen value, -l value  set the socks server port (default: "127.0.0.1:1111")
   --no-auth                 disable socks5 authentication (default: true)
   --auth value              socks5 creds, username:password, leave empty to auto generate
   --mode value              connection mode, choices are auto, full, half (default: "auto")
   --ua value                the user-agent used to send request (default: "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.1.2.3")
   --timeout value           http request timeout in seconds (default: 10)
   --buf-size value          set the request max body size (default: 327680)
   --debug, -d               debug the traffic, print more details (default: false)
   --help, -h                show help (default: false)
   --version, -v             print the version (default: false)
```

命令行版本与界面版配置完全一致，可以对照界面版功能来使用，最简单的只需指定连接目标

```bash
$ ./suo5 -t https://example.com/proxy.jsp
```

自定义 socks5 监听在 `0.0.0.0:7788`，并自定义认证信息为 `test:test123`

```bash
$ ./suo5 -t https://example.com/proxy.jsp -l 0.0.0.0:7788 --auth test:test123
```
### 特别提醒
`User-Agent` (`ua`) 的配置本地端与服务端是绑定的，如果修改了其中一个，另一个也必须对应修改才能连接上。

## 使用体验
![experience](./tests/img/experience.gif)

除此之外，我做了一点模拟环境的性能测试，感兴趣的可以看下 [性能测试](./tests)

## 常见问题

1. 什么是全双工和半双工?
 
    **全双工** 仅需发送一个 HTTP 请求即可构建出一个 HTTP 隧道, 实现 TCP 的双向通信。可以理解成这个请求既是一个上传请求又是一个下载请求，只要连接不断开
    ，就会一直下载，一直上传, 便可以借此做双向通信。

    **半双工** 在部分场景下不支持 `全双工` 模式（比如有反代），可以退而求其次做半双工，即发送一个请求构建一个下行的隧道，同时

2. suo5 和 Neo-reGeorg 怎么选？
    
    如果目标是 Java 的站点，可以使用 `suo5` 来构建 http 隧道，大多数情况下 `suo5` 都要比 `neo` 更稳定速度更快。但 `neo` 提供了非常多种类的服务端支持，兼容性很好，而且也支持一些 `suo5` 当前还在开发的功能，比如负载均衡等，也支持更灵活的定制化。
 
## 接下来

- [ ] 支持配置上游 socks 代理
- [ ] 支持负载均衡的场景
- [ ] 支持 .Net 的类型