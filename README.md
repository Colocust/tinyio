# tinyio

`tinyio`是一个简洁的网络IO事件处理器。它不同于标准的Go net包，而是直接使用epoll系统调用，将每个连接当作事件进行处理。

## 特征

- 单线程事件循环，没有线程并发的安全问题，适用于一些内存操作的服务。
- 简单的API
- 内存使用率低

## 入门

### 安装
```sh
go get -u github.com/Colocust/tinyio
```

### 用法
启动`tinyio`十分简单，只需要将你绑定的地址以及具体逻辑的实现传递给app包中的Boot函数就好了。
以下是一个简单的示例：
```go
package main

import (
	"github.com/Colocust/tinyio/app"
)

func main() {
	app.Boot("127.0.0.1:8877", func(in []byte) (out []byte) {
		out = in
		return
	})
}

```

## 下个版本规划

- 支持多线程处理read以及write逻辑
- 支持自定义事件


