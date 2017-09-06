package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nawa/back-friend/rest/utils"
)

func PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func ResetHandler(c *gin.Context) {
	err := getStorage(c).Reset()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		return
	}
}
