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
