package goredis_server

import (
	. "../goredis"
	"./storage"
	//"./uuid"
	"fmt"
	"strings"
)

var (
	WrongKindReply = ErrorReply("Wrong kind opration")
)

// GoRedisServer
type GoRedisServer struct {
	CommandHandler
	RedisServer
	// 数据源
	datasource storage.DataSource
	// 从库
	slaveMgr *SlaveServerManager
	// 当前实例名字
	uid string
	// 从库状态
	ReplicationInfo ReplicationInfo
}

func NewGoRedisServer() (server *GoRedisServer) {
	server = &GoRedisServer{}
	// set as itself
	server.SetHandler(server)
	// default datasource
	server.datasource = storage.NewMemoryDataSource()
	// slave
	server.slaveMgr = NewSlaveServerManager(server)
	server.ReplicationInfo = ReplicationInfo{}
	return
}

func (server *GoRedisServer) Listen(host string) {
	port := strings.Split(host, ":")[1]
	var e1 error
	server.datasource, e1 = storage.NewLevelDBDataSource("/tmp/goredis_" + port + ".ldb")
	if e1 != nil {
		panic(e1)
	}
	server.initUID()
	server.RedisServer.Listen(host)
}

func (server *GoRedisServer) initUID() {
	// uuidKey := "__goredis_uuid__"
	// data, e1 := server.Storages.StringStorage.Get(uuidKey)
	// if e1 != nil {
	// 	panic(e1)
	// }
	// if data != nil {
	// 	switch data.(type) {
	// 	case string:
	// 		server.uid = data.(string)
	// 	case []byte:
	// 		server.uid = string(data.([]byte))
	// 	default:
	// 		panic("Bad UUID")
	// 	}
	// } else {
	// 	server.uid = uuid.NewV4().String()
	// 	server.Storages.StringStorage.Set(uuidKey, server.uid)
	// }
	// fmt.Println("GoRedis UUID:", server.UID())
}

func (server *GoRedisServer) UID() string {
	return server.uid
}

// for CommandHandler
func (server *GoRedisServer) On(name string, cmd *Command) (reply *Reply) {
	go func() {
		fmt.Println("Slave Send:", cmd.String())
		server.slaveMgr.PublishCommand(cmd)
	}()
	return ErrorReply("Not Supported: " + cmd.String())
}
