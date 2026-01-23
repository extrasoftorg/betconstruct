package backoffice

import (
	"bytes"
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
		bytes.NewReader(body),
		c,
	)
	if err != nil {
		return nil, err
	}
	return transactions.Transactions, nil
}

type ListDepositsRequestDate struct {
	time.Time
}

func (d ListDepositsRequestDate) MarshalJSON() ([]byte, error) {
	layout := "02-01-06 - 15:04:05"
	return json.Marshal(time.Time(d.Time).Format(layout))
}

type ListDepositsRequest struct {
	FromDate *ListDepositsRequestDate `json:"FromCreatedDateLocal"`
	ToDate   *ListDepositsRequestDate `json:"ToCreatedDateLocal"`
}

type listDepositsResponse struct {
	Documents struct {
		Deposits []Deposit `json:"Objects"`
	} `json:"Documents"`
}

func (c *client) ListDeposits(ctx context.Context, req ListDepositsRequest) ([]Deposit, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	deposits, err := makeRequest[listDepositsResponse](
		ctx,
		http.MethodPost,
		"/Financial/GetDepositsWithdrawalsWithPaging",
		bytes.NewReader(body),
		c,
	)
	if err != nil {
		return nil, err
	}
	return deposits.Documents.Deposits, nil
}
