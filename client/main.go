package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc/proto"
	"log"
	"time"
)

func getConn() *grpc.ClientConn {
	conn, err := grpc.Dial("localhost:8090", grpc.WithInsecure())
	if err != nil {
		time.Sleep(time.Second * 5)
		conn = getConn()
		log.Println("client:", err.Error())
	}
	return conn
}

func run(conn *grpc.ClientConn) {
	orderServiceClient := pb.NewIp2AdderServiceClient(conn)
	IpList := &pb.IpInfo{Ip: []string{"127.0.0.1"}}
	res, err := orderServiceClient.GetAdderToIp(context.Background(), IpList)
	if err != nil {
		log.Println("client:", err)
	}
	fmt.Println(res.GetAdder())
}

func main() {
	conn := getConn()
	for {
		go run(conn)
		time.Sleep(time.Millisecond * 500)
	}
	conn.Close()
}
