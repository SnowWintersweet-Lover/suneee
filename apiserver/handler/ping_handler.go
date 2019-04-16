package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary ping/pong
// @Description ping/pong
// @Tags ping/pong
// @Accept  json
// @Produce  json
// @Success 200 {string} json "pong"
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
