package backoffice

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type ListDepositsInput struct {
	FromDate time.Time
	ToDate   time.Time
	Limit    int
	Offset   int
}

type listDepositsPayload struct {
	FromDate  *string `json:"FromCreatedDateLocal"`
	ToDate    *string `json:"ToCreatedDateLocal"`
	MaxRows   int     `json:"MaxRows"`
	SkeepRows int     `json:"SkeepRows"`
}

func (in ListDepositsInput) wire(loc *time.Location) listDepositsPayload {
	p := listDepositsPayload{
		MaxRows:   20,
		SkeepRows: 0,
	}
	if !in.FromDate.IsZero() {
		in.FromDate = in.FromDate.In(loc)
		fromDate := in.FromDate.Format("02-01-06 - 15:04:05")
		p.FromDate = &fromDate
	}
	if !in.ToDate.IsZero() {
		in.ToDate = in.ToDate.In(loc)
		toDate := in.ToDate.Format("02-01-06 - 15:04:05")
		p.ToDate = &toDate
	}

	if in.Limit > 0 {
		p.MaxRows = in.Limit
	}
	if in.Offset > 0 {
		p.SkeepRows = in.Offset
	}

	return p
}

type ListDepositsOutput struct {
	Deposits []Deposit
	Count    int
}

type Deposit struct {
	ID              int64
	Amount          float64
	PlayerID        PlayerID
	CreatedAt       time.Time
	PaymentMethod   string
	Currency        string
	PartnerID       int64
	PaymentMethodID int32
}

func (c *client) ListDeposits(ctx context.Context, in ListDepositsInput) (*ListDepositsOutput, error) {
	body, err := json.Marshal(in.wire(c.timeLocation))
	if err != nil {
		return nil, err
	}

	type responseDeposit struct {
		ID              int64   `json:"Id"`
		Amount          float64 `json:"Amount"`
		PlayerID        int64   `json:"ClientId"`
		CreatedAt       string  `json:"CreatedLocal"`
		PaymentMethod   string  `json:"PaymentSystemName"`
		Curreny         string  `json:"CurrencyId"`
		PartnerID       int64   `json:"PartnerId"`
		PaymentMethodID int32   `json:"PaymentSystemId"`
	}
	type response struct {
		Documents struct {
			Deposits []responseDeposit `json:"Objects"`
			Count    int               `json:"Count"`
		} `json:"Documents"`
	}
	resp, err := makeRequest[response](
		ctx,
		http.MethodPost,
		"/Financial/GetDepositsWithdrawalsWithPaging",
		body,
		c,
	)
	if err != nil {
		return nil, err
	}

	deposits := make([]Deposit, len(resp.Documents.Deposits))
	for i, d := range resp.Documents.Deposits {
		createdAt, err := time.ParseInLocation("2006-01-02T15:04:05.999", d.CreatedAt, c.timeLocation)
		if err != nil {
			return nil, err
		}

		deposit := Deposit{
			ID:              d.ID,
			Amount:          d.Amount,
			PlayerID:        PlayerID(d.PlayerID),
			PaymentMethod:   d.PaymentMethod,
			CreatedAt:       createdAt.UTC(),
			Currency:        d.Curreny,
			PartnerID:       d.PartnerID,
			PaymentMethodID: d.PaymentMethodID,
		}
		deposits[i] = deposit
	}

	return &ListDepositsOutput{
		Deposits: deposits,
		Count:    resp.Documents.Count,
	}, nil
}
