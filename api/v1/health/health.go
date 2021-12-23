package health

import (
	"github.com/gin-gonic/gin"
	msghealth "github.com/romberli/go-template/pkg/message/health"
	"github.com/romberli/go-template/pkg/resp"
)

const (
	pongString = "{\"ping\": \"pong\"}"
)

// @Tags health
// @Summary ping
// @Produce  application/json
// @Success 200 {string} string "{"code": 200, "data": {"ping": "pong"}}"
// @Router /api/v1/health/ping [get]
func Ping(c *gin.Context) {
	resp.ResponseOK(c, pongString, msghealth.InfoHealthPing)
}
