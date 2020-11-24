package state

import (
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/gin/response"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	response.Response(c, "pong")
}
