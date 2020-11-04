package main

import (
	"consul_grpc/service/common"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	servName   = "ip2adder_service"
	port       = 6868
	healthPort = 9999
	healthType = "http"
	ip = "127.0.0.1"
)

func main() {
	//监听本地端口
	listen, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	s := getGrpcServer()
	go func() {
		err = s.Serve(listen)
		if err != nil {
			panic(err)
		}
	}()

	//注册consul
	consul, err := RegisterConsul()
	if err != nil {
		panic(err)
	}
	//	监听关闭服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	err = DeRegisterConsul(consul)
	if err != nil {
		panic(err)
	}

	fmt.Println("服务关闭")

}
// 注册服务到consul
func RegisterConsul() (*common.Consul, error) {

	//localip 应该是内网ip，不然容易导致健康监测出错
	consul, err := common.NewConsulService(servName, servName, port, healthPort, healthType,ip)

	if err != nil {
		return nil, err
	}
	_, err = consul.Register()

	if err != nil {
		return nil, err
	}

	return consul, nil
}


// 将服务进行注销
func DeRegisterConsul(c *common.Consul) error {
	return c.Deregister()
}