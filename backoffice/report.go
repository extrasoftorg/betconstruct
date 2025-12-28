package backoffice

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type ListSportBetsRequestDate struct {
	time.Time
}

func (d ListSportBetsRequestDate) MarshalJSON() ([]byte, error) {
	layout := "02-01-06 - 15:04:05"
	return json.Marshal(time.Time(d.Time).Format(layout))
}

type ListSportBetsRequest struct {
	FromDate *ListSportBetsRequestDate `json:"StartDateLocal"`
	ToDate   *ListSportBetsRequestDate `json:"EndDateLocal"`
	PlayerID *PlayerID                 `json:"ClientId"`
	Status   *SportBetStatus           `json:"State"`
}

type listSportBetsResponse struct {
	Data struct {
		Bets []SportBet `json:"Objects"`
	} `json:"BetData"`
}

func (c *client) ListSportBets(ctx context.Context, req ListSportBetsRequest) ([]SportBet, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	bets, err := makeRequest[listSportBetsResponse](
		ctx,
		http.MethodPost,
		"/Report/GetBetHistory",
		bytes.NewReader(body),
		makeRequestOptions{
			httpClient: c.httpClient,
			authToken:  c.authToken,
		},
	)
	if err != nil {
		return nil, err
	}
	return bets.Data.Bets, nil
}
