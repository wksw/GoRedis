GoRedis
=======

### RedisServer Implemented by Go
#### 说明
	1、围绕Redis协议衍生出的数据处理框架

#### 已实现
	1、核心RedisServer协议层(src/goredis，5个文件)，提供高性能简易API（src/main/simple_server.go，GET/SET 10w+/s）
	2、面向应用的GoRedisServer（src/goredis_server），提供内存String/List/Hash操作
	3、GoRedisServer作为原生Redis的从库

#### 开发中
	1、SlaveOf和Sync指令，和原生Redis之间的主从同步
	2、双主模式GoRedisServer
	3、MultiSlaveOf，一个GoRedisServer作为n个原生Redis的从库，汇总备份数据到第三方存储（MongoDB/MySQL/HBase）

#### vi ~/.profile 

	export GOPATH=/User/lptmoon/Downloads/go/gopath/

#### Install:

	go get github.com/latermoon/GoRedis/src/goredis

#### Update:

	go get -u github.com/latermoon/GoRedis/src/goredis

#### RedisServer Demo:

	package main

	import (
		"fmt"
		. "github.com/latermoon/GoRedis/src/goredis"
		"runtime"
	)

	// ==============================
	// 简单的Redis服务器处理类
	// ==============================
	type SimpleServerHandler struct {
		CommandHandler
		kvCache map[string]interface{} // KeyValue
		kvLock  chan int               // Set操作的写锁
	}

	func NewSimpleServerHandler() (handler *SimpleServerHandler) {
		handler = &SimpleServerHandler{}
		handler.kvCache = make(map[string]interface{})
		handler.kvLock = make(chan int, 1)
		return
	}

	func (s *SimpleServerHandler) On(name string, cmd *Command) (reply *Reply) {
		reply = ErrorReply("Not Supported: " + cmd.String())
		return
	}

	func (s *SimpleServerHandler) OnGET(cmd *Command) (reply *Reply) {
		key := cmd.StringAtIndex(1)
		value := s.kvCache[key]
		reply = BulkReply(value)
		return
	}

	func (s *SimpleServerHandler) OnSET(cmd *Command) (reply *Reply) {
		key := cmd.StringAtIndex(1)
		value := cmd.StringAtIndex(2)
		s.kvLock <- 0
		s.kvCache[key] = value
		<-s.kvLock
		reply = StatusReply("OK")
		return
	}

	func (s *SimpleServerHandler) OnINFO(cmd *Command) (reply *Reply) {
		lines := "Powerby GoRedis" + "\n"
		lines += "SimpleRedisServer" + "\n"
		lines += "Support GET/SET/INFO" + "\n"
		reply = BulkReply(lines)
		return
	}

	func main() {
		runtime.GOMAXPROCS(2)
		fmt.Println("SimpleServer start, listen 1603 ...")
		server := NewRedisServer(NewSimpleServerHandler())
		server.Listen(":1603")
	}

