package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nawa/back-friend/rest/utils"
	"github.com/nawa/back-friend/storage"
)

func TakeHandler(c *gin.Context) {
	playerID, ok := utils.GetQuery(c, "playerId", true)
	if !ok {
		return
	}
	points, ok := utils.GetQueryInt(c, "points", true)
	if !ok {
		return
	}

	s := getStorage(c)
	err := s.PlayerStorage().Take(playerID, points)
	if err != nil {
		if err == storage.ErrNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, utils.ErrorResponse{Error: "player not found"})
		} else if illegalArgumentError, ok := err.(*storage.IllegalArgumentError); ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{Error: illegalArgumentError.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		}
		return
	}
}

func FundHandler(c *gin.Context) {
	playerID, ok := utils.GetQuery(c, "playerId", true)
	if !ok {
		return
	}
	points, ok := utils.GetQueryInt(c, "points", true)
	if !ok {
		return
	}
	s := getStorage(c)
	err := s.PlayerStorage().FundOrCreate(playerID, points)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		return
	}
}

func PlayerBalanceHandler(c *gin.Context) {
	playerID, ok := utils.GetQuery(c, "playerId", true)
	if !ok {
		return
	}

	s := getStorage(c)
	player, err := s.PlayerStorage().Get(playerID)
	if err != nil {
		if err == storage.ErrNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, utils.ErrorResponse{Error: "player not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, player)
}
