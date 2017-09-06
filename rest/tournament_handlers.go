package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nawa/back-friend/model"
	"github.com/nawa/back-friend/rest/utils"
	"github.com/nawa/back-friend/storage"
)

func AnnounceTournamentHandler(c *gin.Context) {
	tournamentID, ok := utils.GetQuery(c, "tournamentId", true)
	if !ok {
		return
	}

	deposit, ok := utils.GetQueryInt(c, "deposit", true)
	if !ok {
		return
	}
	s := getStorage(c)
	err := s.TournamentStorage().Create(tournamentID, deposit)
	if err != nil {
		if illegalArgumentError, ok := err.(*storage.IllegalArgumentError); ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{Error: illegalArgumentError.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		}
	}
}

func JoinTournamentHandler(c *gin.Context) {
	tournamentID, ok := utils.GetQuery(c, "tournamentId", true)
	if !ok {
		return
	}

	playerID, ok := utils.GetQuery(c, "playerId", true)
	if !ok {
		return
	}

	backersID, _ := utils.GetQueryArray(c, "backerId", false)

	err := getStorage(c).TournamentStorage().Join(tournamentID, playerID, backersID...)
	if err != nil {
		if illegalArgumentError, ok := err.(*storage.IllegalArgumentError); ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{Error: illegalArgumentError.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		}
	}
}

func ResultTournamentHandler(c *gin.Context) {
	tournament := model.Tournament{}
	err := c.ShouldBindWith(&tournament, binding.JSON)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			utils.ErrorResponse{Error: "Incorrect json body: " + err.Error()},
		)
	}

	err = getStorage(c).TournamentStorage().Result(tournament)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
	}

}
