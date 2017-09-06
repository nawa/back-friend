package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nawa/back-friend/storage"
)

type Service struct {
	httpPort  int
	ginEngine *gin.Engine
	storage   storage.Storage
}

func NewService(port int, storage storage.Storage) *Service {
	router := gin.New()

	router.Use(gin.Recovery())
	service := &Service{
		httpPort:  port,
		ginEngine: router,
		storage:   storage,
	}
	router.Use(service.injectionsMiddleware)
	service.addRoutes()
	return service
}

func (s *Service) Start() error {
	return s.ginEngine.Run(fmt.Sprintf(":%v", s.httpPort))
}

func (s *Service) addRoutes() {
	router := s.ginEngine
	router.GET("/ping", PingHandler)
	router.GET("/take", TakeHandler)                             //GET /take?playerId=P1&points=300
	router.GET("/fund", FundHandler)                             //GET /fund?playerId=P2&points=300
	router.GET("/balance", PlayerBalanceHandler)                 //GET /balance?playerId=P1
	router.GET("/announceTournament", AnnounceTournamentHandler) //GET /announceTournament?tournamentId=1&deposit=1000
	router.GET("/joinTournament", JoinTournamentHandler)         //GET /joinTournament?tournamentId = 1&playerId = P1&backerId = P2&backerId = P3
	router.POST("/resultTournament", ResultTournamentHandler)    //POST /resultTournament with a POST
	router.GET("/reset", ResetHandler)
}

func (s *Service) injectionsMiddleware(c *gin.Context) {
	c.Set("storage", s.storage)
}

func getStorage(c *gin.Context) storage.Storage {
	return c.MustGet("storage").(storage.Storage)
}
