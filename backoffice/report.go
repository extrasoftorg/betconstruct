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
	FromDate   *ListSportBetsRequestDate `json:"StartDateLocal"`
	ToDate     *ListSportBetsRequestDate `json:"EndDateLocal"`
	PlayerID   *PlayerID                 `json:"ClientId"`
	Status     *SportBetStatus           `json:"State"`
	ToCurrency string                    `json:"ToCurrencyId"`
}

type listSportBetsResponse struct {
	Data struct {
		Bets []SportBet `json:"Objects"`
	} `json:"BetData"`
}

type GetBetHistoryResult struct {
	Bets   []SportBet
	Totals SportBetTotals
}

type getBetHistoryResponse struct {
	BetData struct {
		Bets []SportBet `json:"Objects"`
	} `json:"BetData"`
	BetTotals SportBetTotals `json:"BetTotals"`
}

func (c *client) GetBetHistory(ctx context.Context, req ListSportBetsRequest) (*GetBetHistoryResult, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := makeRequest[getBetHistoryResponse](
		ctx,
		http.MethodPost,
		"/Report/GetBetHistory",
		bytes.NewReader(body),
		c,
	)
	if err != nil {
		return nil, err
	}
	return &GetBetHistoryResult{
		Bets:   resp.BetData.Bets,
		Totals: resp.BetTotals,
	}, nil
}

type SportKindReport struct {
	Id         int      `json:"Id"`
	Name       string   `json:"Name"`
	BetCount   *int     `json:"BetCount"`
	Stakes     float64  `json:"Stakes"`
	Winnings   float64  `json:"Winnings"`
	Profit     float64  `json:"Profit"`
	Profitness float64  `json:"Profitness"`
}

type GetSportKindReportRequest struct {
	StartTime time.Time
	EndTime   time.Time
	Currency  string
}

type getSportKindReportRequest struct {
	IsLive            string                    `json:"IsLive"`
	Sources           []any                     `json:"Sources"`
	IsTest            string                    `json:"IsTest"`
	IsBonus           *bool                     `json:"IsBonus"`
	IsCalculated      bool                      `json:"IsCalculated"`
	FilterByCurrency  *string                   `json:"FilterByCurrency"`
	StartTimeLocal    ListSportBetsRequestDate  `json:"StartTimeLocal"`
	EndTimeLocal      ListSportBetsRequestDate  `json:"EndTimeLocal"`
	Currency          string                    `json:"Currency"`
}

func (c *client) GetSportKindReport(ctx context.Context, req GetSportKindReportRequest) ([]SportKindReport, error) {
	body, err := json.Marshal(getSportKindReportRequest{
		IsLive:           "",
		Sources:          []any{},
		IsTest:           "false",
		IsBonus:          nil,
		IsCalculated:     false,
		FilterByCurrency: nil,
		StartTimeLocal:   ListSportBetsRequestDate{req.StartTime},
		EndTimeLocal:     ListSportBetsRequestDate{req.EndTime},
		Currency:         req.Currency,
	})
	if err != nil {
		return nil, err
	}
	results, err := makeRequest[[]SportKindReport](ctx, http.MethodPost, "/Report/GetSportKindReport", bytes.NewReader(body), c)
	if err != nil {
		return nil, err
	}
	return *results, nil
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
		c,
	)
	if err != nil {
		return nil, err
	}
	return bets.Data.Bets, nil
}
