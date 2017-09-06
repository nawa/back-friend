package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func GetQuery(c *gin.Context, paramName string, required bool) (value string, ok bool) {
	if value, ok = c.GetQuery(paramName); !ok {
		if required {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				ErrorResponse{fmt.Sprintf("Missing '%v' query param", paramName)},
			)
		}
	}
	return
}

func GetQueryArray(c *gin.Context, paramName string, required bool) (value []string, ok bool) {
	if value, ok = c.GetQueryArray(paramName); !ok {
		if required {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				ErrorResponse{fmt.Sprintf("Missing '%v' query param", paramName)},
			)
		}
	}
	return
}

func GetQueryInt(c *gin.Context, paramName string, required bool) (value int, ok bool) {
	var sValue string
	if sValue, ok = GetQuery(c, paramName, required); !ok {
		return
	}
	value, err := strconv.Atoi(sValue)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			ErrorResponse{fmt.Sprintf("Incorrect '%v' query param", paramName)},
		)
		return 0, false
	}
	return
}
