package main

import (
	"context"
	"google.golang.org/grpc"
	pb "grpc/proto"
	"net"
)

type AdderInfo struct {
}

func (s AdderInfo) GetAdderToIp(c context.Context, src *pb.IpInfo) (*pb.AdderInfo, error) {
	AdderData := &pb.AdderInfo{
		Adder: src.Ip,
	}
	return AdderData, nil
}

func main() {
	server := grpc.NewServer()
	pb.RegisterIp2AdderServiceServer(server, &AdderInfo{})
	lis, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err.Error())
	}
	server.Serve(lis)
}
