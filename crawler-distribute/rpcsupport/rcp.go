package rpcsupport

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

//启动rpc服务器,提供要注册的对象service
//则该service的方法均可用rpc进行调用
func ServeRpc(host string, service interface{}) error {
	rpc.Register(service)
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}
	log.Printf("Listening on %s", host)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}

		//开启goroutine来执行,否则在执行过程中,卡住无法再接受新的连接请求
		go jsonrpc.ServeConn(conn)
	}
	return nil
}

//供客户端连接rpc服务器
func NewClient(host string) (*rpc.Client, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	return jsonrpc.NewClient(conn), nil
}
