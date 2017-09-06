package functionaltests

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/go-resty/resty"
	"github.com/nawa/back-friend/model"
	assert "github.com/stretchr/testify/require"
)

//Configure me through env variable
//const baseURL = "http://localhost:8080"

const baseURL = "http://backfriend-rest:8080"

func TestPing(t *testing.T) {
	resp, err := resty.R().
		Get(baseURL + "/ping")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "pong", string(resp.Body()))
}

func TestMainUseCase(t *testing.T) {
	resp, err := resetDb()
	assertResponse(t, resp, err, 200)

	resp, err = fund("P1", "300")
	s := resp.String()
	fmt.Println(string(s))
	assertResponse(t, resp, err, 200)

	resp, err = fund("P2", "300")
	assertResponse(t, resp, err, 200)

	resp, err = fund("P3", "300")
	assertResponse(t, resp, err, 200)

	resp, err = fund("P4", "500")
	assertResponse(t, resp, err, 200)

	resp, err = fund("P5", "1000")
	assertResponse(t, resp, err, 200)

	resp, err = announceTournament("1", "1000")
	assertResponse(t, resp, err, 200)

	resp, err = joinTournament("1", "P5")
	assertResponse(t, resp, err, 200)

	resp, err = joinTournament("1", "P1", "P2", "P3", "P4")
	assertResponse(t, resp, err, 200)

	result := `{"tournamentId": "1", "winners": [{"playerId": "P1", "prize": 2000}]}`
	resp, err = resultTournament(result)
	assertResponse(t, resp, err, 200)

	resp, player, err := getBalance("P1")
	assertResponse(t, resp, err, 200)
	assert.Equal(t, 550, player.Balance)

	resp, player, err = getBalance("P2")
	assertResponse(t, resp, err, 200)
	assert.Equal(t, 550, player.Balance)

	resp, player, err = getBalance("P3")
	assertResponse(t, resp, err, 200)
	assert.Equal(t, 550, player.Balance)

	resp, player, err = getBalance("P4")
	assertResponse(t, resp, err, 200)
	assert.Equal(t, 750, player.Balance)

	resp, player, err = getBalance("P5")
	assertResponse(t, resp, err, 200)
	assert.Equal(t, 0, player.Balance)
}

func assertResponse(t *testing.T, resp *resty.Response, err error, expectedStatus int) {
	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode())
}

func resetDb() (*resty.Response, error) {
	return resty.R().
		Get(baseURL + "/reset")
}

func getBalance(playerID string) (resp *resty.Response, player *model.Player, err error) {
	player = &model.Player{}
	resp, err = resty.R().
		SetQueryParam("playerId", playerID).
		SetResult(player).
		Get(baseURL + "/balance")
	return
}

func fund(playerID, points string) (*resty.Response, error) {
	return resty.R().
		SetQueryParam("playerId", playerID).
		SetQueryParam("points", points).
		Get(baseURL + "/fund")
}

func announceTournament(tournamentID, deposit string) (*resty.Response, error) {
	return resty.R().
		SetQueryParam("tournamentId", tournamentID).
		SetQueryParam("deposit", deposit).
		Get(baseURL + "/announceTournament")
}

func joinTournament(tournamentID, playerID string, backersID ...string) (*resty.Response, error) {
	return resty.R().
		SetQueryParam("tournamentId", tournamentID).
		SetQueryParam("playerId", playerID).
		SetMultiValueQueryParams(url.Values{
			"backerId": backersID,
		}).
		Get(baseURL + "/joinTournament")
}

func resultTournament(json string) (*resty.Response, error) {
	return resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte(json)).
		Post(baseURL + "/resultTournament")
}
