package router

import (
	"github.com/gin-gonic/gin"
	"im/ipconfig/web/controllers"
)

func WebRouterInit() (routers *gin.Engine) {
	/*TODO:Example Initialize all web routes*/
	routers = gin.Default()

	/*TODO:ip router group*/

	ip := routers.Group("ip")
	{
		ip.GET("list", controllers.IpsGetHandler)
	}
	return
}

func Init() {
	routers := WebRouterInit()

	routers.Run(":8084")
}
