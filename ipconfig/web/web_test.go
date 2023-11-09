package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"im/ipconfig/web/router"
	"net/http"
	"testing"
	"time"
)

func TestGinWeb(t *testing.T) {
	router := gin.Default()
	router.Use(customMiddleware())
	router.GET("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"msg": "pong",
		})
	})

	router.Run(":8081")
}

func customMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) { //真正的中间件类型
		t := time.Now()
		c.Set("msg", "This is a test of middleware")
		//它执行调用处理程序内链中的待处理处理程序
		//让原本执行的逻辑继续执行
		c.Next()

		end := time.Since(t)
		fmt.Printf("耗时：%v\n", end.Seconds())
		status := c.Writer.Status()
		fmt.Println("状态监控:", status)
	}
}

func TestRunIpConfigWeb(t *testing.T) {
	router.Init()
}
