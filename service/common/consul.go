package common

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"log"
	"net/http"
	"net/url"
	"strings"
)
type Consul struct {
	registration *consulapi.AgentServiceRegistration
	client       *consulapi.Client
}
func (c *Consul) MakeClient() (*consulapi.Client, error) {

	uriStr := "127.0.0.1:8500"
	token := ""

	config := consulapi.DefaultConfig()
	config.Address = uriStr

	if len(token) > 0 {
		config.Token = token
	} else {
		config.Token = "defaultToken"
	}

	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
// @param svrName string: 要注册当服务名
// @param useType string: 对应consul中的tag, 可用于过滤
// @param svrPort int: 服务对应的端口号
// @param healthPort int: http检测时需要对应端口号，tcp检测默认当前端口
// @param healthType string: http或tcp,
// @param localIp string: 当前节点的内网IP，即其他服务能访问到的IP

func NewConsulService(svrName string, useType string, svrPort int, healthPort int, healthType string,ip string) (*Consul, error) {
	var err error
	// 注册配置信息
	reg := &consulapi.AgentServiceRegistration{
		ID:      strings.ToLower(fmt.Sprintf("%s_%s", svrName, ip)), // 生成一个唯一当服务ID
		Name:    strings.ToLower(fmt.Sprintf("%s", svrName)),        // 注册服务名
		Tags:    []string{strings.ToLower(useType)},                 // 标签
		Port:    svrPort,                                            // 端口号
		Address: ip,                                                 // 所在节点ip地址
	}

	// 健康检测配置信息
	reg.Check = &consulapi.AgentServiceCheck{
		TCP:                            fmt.Sprintf("%s:%d", ip, svrPort),
		Timeout:                        "1s",
		Interval:                      "15s",
		DeregisterCriticalServiceAfter: "30s", // 30秒服务不可达时，注销服务
		Status:                         consulapi.HealthPassing,                        // 服务启动时，默认正常
	}

	if healthType == "http" {
		// http检测默认/health路径
		reg.Check.HTTP = fmt.Sprintf("http://%s:%d%s", reg.Address, healthPort, "/health")
		reg.Check.TCP = ""
	}

	if len(reg.Check.HTTP) > 0 {
		err = RunHealthCheck(reg.Check.HTTP)
	}
	//准备agent
	temConsul := new(Consul)
	temConsul.registration = reg

	//准备client
	temConsul.client, err = temConsul.MakeClient()

	return temConsul, err

}

// 服务注册
func (c *Consul) Register() (string, error) {

	var err error

	err = c.client.Agent().ServiceRegister(c.registration)

	return c.registration.ID, err
}

// 服务注销
func (c *Consul) Deregister() error {

	var err error

	err = c.client.Agent().ServiceDeregister(c.registration.ID)

	return err
}

//	TODO 获取到服务的健康状态
func RunHealthCheck(addr string) error {
	// 实现一个接口类

	uri, err := url.Parse(addr)
	if err != nil {
		return err
	}

	http.HandleFunc(uri.Path, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
		log.Println("porttttt")
	})
	go func() {
		err := http.ListenAndServe(uri.Host, nil)
	log.Println("porttttt")
		if err != nil {
			fmt.Println(err)
		}
	}()
	return nil
}

// TODO 获取服务列表
func GetAvailableService() {

	uriStr := "127.0.0.1:8500"
	token := ""

	config := consulapi.DefaultConfig()
	config.Address = uriStr

	if len(token) > 0 {
		config.Token = token
	} else {
		config.Token = "defaultToken"
	}

	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Println(err)
	}

	if err != nil {
		log.Println(err)
	}
	fmt.Println(client)
}