package backoffice

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type ListWithdrawalsRequest struct {
	FromDate time.Time `json:"FromDateLocal"`
	ToDate   time.Time `json:"ToDateLocal"`
	ID       int64     `json:"Id"`
}

func (r *ListWithdrawalsRequest) MarshalJSON() ([]byte, error) {
	type wire struct {
		FromDate *time.Time `json:"FromDateLocal"`
		ToDate   *time.Time `json:"ToDateLocal"`
		ID       *int64     `json:"Id"`
	}
	w := wire{}
	if !r.FromDate.IsZero() {
		w.FromDate = &r.FromDate
	}
	if !r.ToDate.IsZero() {
		w.ToDate = &r.ToDate
	}
	if r.ID != 0 {
		w.ID = &r.ID
	}
	return json.Marshal(w)
}

type listWithdrawalsResponse struct {
	Withdrawals []Withdrawal `json:"ClientRequests"`
}

func (c *client) ListWithdrawals(ctx context.Context, req ListWithdrawalsRequest) ([]Withdrawal, error) {
	body, err := req.MarshalJSON()
	if err != nil {
		return nil, err
	}
	withdrawals, err := makeRequest[listWithdrawalsResponse](
		ctx,
		http.MethodPost,
		"/Client/GetClientWithdrawalRequestsWithTotals",
		bytes.NewReader(body),
		c,
	)
	if err != nil {
		return nil, err
	}
	return withdrawals.Withdrawals, nil
}

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
