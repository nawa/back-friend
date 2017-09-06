package model

type Tournament struct {
	ID      string   `json:"tournamentId" binding:"required"`
	Deposit *int     `json:"deposit"`
	Winners []Winner `json:"winners" binding:"required"`
}

type Winner struct {
	ID    string `json:"playerId" binding:"required"`
	Prize int    `json:"prize" binding:"required"`
}
