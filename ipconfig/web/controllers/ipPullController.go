package controllers

import (
	"github.com/gin-gonic/gin"
	"im/common/response"
	"im/ipconfig"
	"im/ipconfig/pkg/domain"
	"net"
	"net/http"
	"strconv"
)

func IpsGetHandler(c *gin.Context) {
	iplist := ipconfig.Ms.GetIps()

	var data []*domain.EndPoint

	for _, ip := range iplist {
		tmp := &domain.EndPoint{}
		host, port, err := net.SplitHostPort(ip)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable,
				response.Response().Err().End())
			return
		}
		tmp.Ip = host
		tmpPort, err := strconv.ParseInt(port, 10, 32)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable,
				response.Response().Err().End())
			return
		}
		tmp.Port = uint16(tmpPort)
		data = append(data, tmp)
	}

	c.JSON(http.StatusOK,
		response.Response().Ok().Put("data", data).End())
}
