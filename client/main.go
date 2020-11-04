package main

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	router := gin.Default()
	router.POST("/get_adder", UserLogin)
	_ = router.Run(":8080")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func UserLogin(c *gin.Context){
	Req := struct{
		Ips []string `json:"ips"`
	}{}
	if err := c.ShouldBindJSON(&Req); err != nil {
		log.Println(err)
	}
	res,err := GetAdder(Req.Ips)
	if err != nil {
		log.Println(err)
	}
	Res(c, res)

}
func Res(c *gin.Context,data interface{}){
	resp := &struct{
		Ret int
		Data interface{}
	}{
		200,
		data,
	}
	c.JSON(200, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
}
func GetAdder(ips []string)([]string, error) {
	conn, err := grpc.Dial("in96.cn:6868", grpc.WithInsecure())
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	orderServiceClient := NewIp2AdderServiceClient(conn)
	//接收的参数  []string类型
	IpList := IpInfo{Ip: ips}
	res, err := orderServiceClient.GetAdderToIp(context.Background(),&IpList)
	if err != nil {
		return []string{},err
	}
	return res.Adder,nil
}