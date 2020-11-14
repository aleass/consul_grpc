package route

import (
	"client/service"
	"github.com/gin-gonic/gin"
)

func RouteInit() {
	router := gin.Default()
	router.POST("/get_adder",service.UserLogin)
	router.POST("/get_way", service.GetWay)
	_ = router.Run(":8080")
}