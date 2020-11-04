package main

import (
	"bufio"
	"bytes"
	//"consul_grpc/service/common"
	"context"
	"encoding/binary"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"os"
	"strconv"
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
//func init1() {
//	NewIpArea := make([]IpArea,0)
//	common.MainDbEngine.Sql("SELECT `ip_start`, `ip_end`, `area` FROM `swhc`.`ip_area` limit 20").Find(&NewIpArea)
//	Ips = make(IpData,0)
//	for _,v := range NewIpArea {
//		ir := IpRange{
//			Begin: v.IpStart,
//			End: v.IpEnd,
//			Data:  []byte(v.Area),
//		}
//		Ips = append(Ips,ir)
//	}
//	l := len(Ips)
//	log.Println(l)
//	if l == 0 {
//		panic("mysql error")
//	}
//
//}
func init() {
	file, err := os.Open("./ip.txt")
	Ips = make(IpData,0)
	if err != nil {
		fmt.Println("文件打开失败 = ", err)
	}
	//及时关闭 file 句柄，否则会有内存泄漏
	defer file.Close()
	//创建一个 *Reader ， 是带缓冲的
	reader := bufio.NewReader(file)
	for {
		str, err := reader.ReadString('\n') //读到一个换行就结束
		if err == io.EOF {                  //io.EOF 表示文件的末尾
			break
		}
		s := strings.Fields(str)
		b,_ := strconv.Atoi(s[0])
		e,_ := strconv.Atoi(s[1])
		ir := IpRange{
			Begin: uint32(b),
			End: uint32(e),
			Data:  []byte(s[2]),
		}
		Ips = append(Ips,ir)
	}
	l := len(Ips)
	log.Println(l)
	if l == 0 {
		panic("mysql error")
	}
}



func ip2Long(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}


type GetAdderToIp struct{}

func(g GetAdderToIp) GetAdderToIp(c context.Context,req *IpInfo) (*AdderInfo, error){
	AdderData := &AdderInfo{}
	IpList := req.GetIp()
	if len(IpList) == 0  {
		return AdderData, nil
	}

	str,err := GetClientIP(c)
	if err != nil {
		log.Println("err:",err)
	}
	log.Println("request client:",str)
	log.Println("request data:",IpList)

	for _,v := range IpList{
		r := ip2Long(v)
		lIps := len(Ips)
		left, right := 1, lIps
		if r < Ips[0].Begin {
			AdderData.Adder = append(AdderData.Adder,"IP地址有误")
			continue
		}
		if r > Ips[lIps-1].End{
			AdderData.Adder = append(AdderData.Adder,"IP地址有误")
			continue
		}
		if r >= 0 && r <= Ips[0].End {
			AdderData.Adder = append(AdderData.Adder,"IANA保留地址")
			continue
		}
		_is := false
		for t:=0;t<lIps;t++ {
			//二分查找?
			mid := (left + right)/2
			if r < Ips[mid].Begin {
				right = mid
				continue
			}
			if r > Ips[mid].End {
				left = mid
				continue
			}
			if r >= Ips[mid].Begin && r <= Ips[mid].End {
				AdderData.Adder = append(AdderData.Adder,string(Ips[mid].Data))
				_is = true
				break
			}
			mid ++
		}
		if !_is {
			AdderData.Adder = append(AdderData.Adder,"unknow")
		}
	}

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