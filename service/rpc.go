package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strings"
)

type IpArea struct {
	IpStart uint32 `xorm:"index(ip) BIGINT(32)"`
	IpEnd   uint32 `xorm:"index(ip) BIGINT(32)"`
	Area    string `xorm:"VARCHAR(255)"`
}
type IpRange struct {
	Begin uint32
	End   uint32
	Data  []byte
	Index   int
}
var  Ips  IpData
type IpData []IpRange

func ip2Long(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}

type GetAdderToIp struct{}

func(g GetAdderToIp) GetAdderToIp(c context.Context,req *IpInfo) (*AdderInfo, error){
	str,err := GetClientIP(c)
	if err != nil {
		log.Println("err:",err)
	}
	log.Println("request client:",str)
	AdderData := &AdderInfo{Adder: []string{"1.1.1.1"}}
	//IpList := req.GetIp()
	//if len(IpList) == 0  {
	//	return AdderData, nil
	//}
	//for _,v := range IpList{
	//	r := ip2Long(v)
	//	lIps := len(Ips)
	//	left, right := 1, lIps
	//	if r < Ips[0].Begin {
	//		AdderData.Adder = append(AdderData.Adder,"IP地址有误")
	//		continue
	//	}
	//	if r > Ips[lIps-1].End{
	//		AdderData.Adder = append(AdderData.Adder,"IP地址有误")
	//		continue
	//	}
	//	if r >= 0 && r <= Ips[0].End {
	//		AdderData.Adder = append(AdderData.Adder,"IANA保留地址")
	//		continue
	//	}
	//	_is := false
	//	for t:=0;t<lIps;t++ {
	//		//二分查找?
	//		mid := (left + right)/2
	//		if r < Ips[mid].Begin {
	//			right = mid
	//			continue
	//		}
	//		if r > Ips[mid].End {
	//			left = mid
	//			continue
	//		}
	//		if r >= Ips[mid].Begin && r <= Ips[mid].End {
	//			AdderData.Adder = append(AdderData.Adder,string(Ips[mid].Data))
	//			_is = true
	//			break
	//		}
	//		mid ++
	//	}
	//	if !_is {
	//		AdderData.Adder = append(AdderData.Adder,"unknow")
	//	}
	//}

	return AdderData, nil
}


func getGrpcServer() *grpc.Server {
	s := grpc.NewServer()
	// 在gRPC服务器注册服务
	RegisterIp2AdderServiceServer(s, new(GetAdderToIp))
	reflection.Register(s)
	return s
}
func GetClientIP(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("[getClinetIP] invoke FromContext() failed")
	}
	if pr.Addr == net.Addr(nil) {
		return "", fmt.Errorf("[getClientIP] peer.Addr is nil")
	}

	addSlice := strings.Split(pr.Addr.String(), ":")
	if addSlice[0] == "[" {
		//本机地址
		return "localhost", nil
	}
	return addSlice[0], nil
}