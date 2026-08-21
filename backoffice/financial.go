package backoffice

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type ListTransactionsRequestDate struct {
	time.Time
}

func (d ListTransactionsRequestDate) MarshalJSON() ([]byte, error) {
	layout := "02-01-06 - 15:04:05"
	return json.Marshal(time.Time(d.Time).Format(layout))
}

type ListTransactionsRequest struct {
	FromDate *ListTransactionsRequestDate `json:"FromCreatedDateLocal"`
	ToDate   *ListTransactionsRequestDate `json:"ToCreatedDateLocal"`
	MaxRows  int                          `json:"MaxRows"`
}

type listTransactionsResponse struct {
	Transactions []Transaction `json:"Objects"`
}

func (c *client) ListTransactions(ctx context.Context, req ListTransactionsRequest) ([]Transaction, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	transactions, err := makeRequest[listTransactionsResponse](
		ctx,
		http.MethodPost,
		"/Financial/GetDocumentsWithPaging",
		body,
		c,
	)
	if err != nil {
		return nil, err
	}
	return transactions.Transactions, nil
}
