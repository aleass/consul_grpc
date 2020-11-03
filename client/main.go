package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":6868", grpc.WithInsecure())
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	orderServiceClient := NewIp2AdderServiceClient(conn)
	//接收的参数  []string类型
	IpList := IpInfo{Ip: []string{"1.0.0.0","1.0.3.254"}}
	res, err := orderServiceClient.GetAdderToIp(context.Background(),&IpList)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.GetAdder())
}