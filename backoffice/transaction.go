package backoffice

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type ListPlayerTransactionsRequestDate struct {
	time.Time
}

func (d ListPlayerTransactionsRequestDate) MarshalJSON() ([]byte, error) {
	layout := "02-01-06"
	return json.Marshal(time.Time(d.Time).Format(layout))
}

type ListPlayerTransactionsRequest struct {
	PlayerID        PlayerID                          `json:"ClientId"`
	FromDate        ListPlayerTransactionsRequestDate `json:"StartTimeLocal"`
	ToDate          ListPlayerTransactionsRequestDate `json:"EndTimeLocal"`
	Currency        string                            `json:"CurrencyId"`
	DocumentTypeIDs []int                             `json:"DocumentTypeIds"`
}

func (r *ListPlayerTransactionsRequest) MarshalJSON() ([]byte, error) {
	type wire struct {
		FromDate        *ListPlayerTransactionsRequestDate `json:"StartTimeLocal"`
		ToDate          *ListPlayerTransactionsRequestDate `json:"EndTimeLocal"`
		Currency        *string                            `json:"CurrencyId"`
		PlayerID        *PlayerID                          `json:"ClientId"`
		DocumentTypeIDs []int                              `json:"DocumentTypeIds"`
	}
	w := wire{}
	if !r.FromDate.IsZero() {
		w.FromDate = &r.FromDate
	}
	if !r.ToDate.IsZero() {
		w.ToDate = &r.ToDate
	}
	if r.Currency != "" {
		w.Currency = &r.Currency
	}
	if r.PlayerID != 0 {
		w.PlayerID = &r.PlayerID
	}
	if len(r.DocumentTypeIDs) > 0 {
		w.DocumentTypeIDs = r.DocumentTypeIDs
	}
	return json.Marshal(w)
}

type listPlayerTransactionsResponse struct {
	Transactions []Transaction `json:"Objects"`
}

func (c *client) ListPlayerTransactions(ctx context.Context, req ListPlayerTransactionsRequest) ([]Transaction, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	transactions, err := makeRequest[listPlayerTransactionsResponse](
		ctx,
		http.MethodPost,
		"/Client/GetClientTransactionsV1",
		bytes.NewReader(body),
		c,
	)
	if err != nil {
		return nil, err
	}
	return transactions.Transactions, nil
}

type ListPlayerCasinoGamesRequest struct {
	PlayerID PlayerID                          `json:"ClientId"`
	FromDate ListPlayerTransactionsRequestDate `json:"FromDateLocal"`
	ToDate   ListPlayerTransactionsRequestDate `json:"ToDateLocal"`
	Currency string                            `json:"Currency"`
}

func (c *client) ListPlayerCasinoGames(ctx context.Context, req ListPlayerCasinoGamesRequest) ([]PlayerCasinoGame, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := makeRequest[[]PlayerCasinoGame](
		ctx,
		http.MethodPost,
		"/Client/GetClientCasinoGames",
		bytes.NewReader(body),
		c,
	)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}
