package ping

import (
	"launchpad/model/request"
	"launchpad/model/response"

	"github.com/gin-gonic/gin"
)

func GetPing(c *gin.Context, req *request.PingReq) (any, error) {
	data := response.PingResp{}
	data.Ping = "PONG"
	return data, nil
}
