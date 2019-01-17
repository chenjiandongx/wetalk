# 📠 wetalk
> wetalk is a chatroom application for coders based on websockets

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

### 🔰 安装

```shell
$ go get github.com/chenjiandongx/wetalk
```

### 📝 使用

$ wetalk
```shell
wetalk is a chatroom application for coders

Example:
  start server: wetalk server -p 8086
  start client：wetalk client ws://127.0.0.1:8086 -u somebody

Usage:
  wetalk [command]

Available Commands:
  client      start websockets client
  help        Help about any command
  server      start websockets server

Flags:
  -h, --help      help for wetalk
      --version   version for wetalk

Use "wetalk [command] --help" for more information about a command.
```

$ wetalk server -h
```shell
start websockets server

Usage:
  wetalk server  [flags]

Flags:
  -h, --help       help for server
  -p, --port int   server port (default 8087
```

$ wetalk server -h
```
start websockets client

Usage:
  wetalk client <addr> [flags]

Flags:
  -h, --help          help for client
  -u, --name string   chatroom's nickname
```

### 📺 示例
![example](https://user-images.githubusercontent.com/19553554/51330669-e7627100-1ab2-11e9-9586-5fb383b6817d.gif)

### 📃 License
MIT [©chenjiandongx](http://github.com/chenjiandongx)
