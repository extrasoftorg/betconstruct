package backoffice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const rgsBaseURL = "https://rgs-webadminapi.betconstruct.com/api"

type rgsDate struct {
	time.Time
}

func (d rgsDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format("02-01-06"))
}

type GetReportByPartnerRequest struct {
	FromDate   time.Time
	ToDate     time.Time
	CurrencyId string
}

type rgsReportByPartnerRequest struct {
	FromDate    rgsDate `json:"FromDate"`
	ToDate      rgsDate `json:"ToDate"`
	PlayerState int     `json:"PlayerState"`
	IsBonusBet  *bool   `json:"IsBonusBet"`
	CurrencyId  string  `json:"CurrencyId"`
	OrderKey    string  `json:"OrderKey"`
	OrderDir    string  `json:"OrderDir"`
	Offset      int     `json:"Offset"`
	Take        int     `json:"Take"`
}

type rgsReportByPartnerResponse struct {
	Result struct {
		TotalProfitByReportCurrency float64 `json:"TotalProfitByReportCurrency"`
	} `json:"Result"`
	HasError         bool    `json:"HasError"`
	ErrorDescription *string `json:"ErrorDescription"`
}

func (c *client) GetCasinoReportByPartner(ctx context.Context, req GetReportByPartnerRequest) (float64, error) {
	body, err := json.Marshal(rgsReportByPartnerRequest{
		FromDate:    rgsDate{req.FromDate},
		ToDate:      rgsDate{req.ToDate},
		PlayerState: 2,
		IsBonusBet:  nil,
		CurrencyId:  req.CurrencyId,
		OrderKey:    "BetAmountByReportCurrency",
		OrderDir:    "OrderByDescending",
		Offset:      0,
		Take:        20,
	})
	if err != nil {
		return 0, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/Reporting/getReportByPartner", rgsBaseURL), bytes.NewReader(body))
	if err != nil {
		return 0, err
	}

	httpReq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	var authToken string
	if c.pool != nil {
		at := c.pool.GetAuthToken(ctx)
		if at != nil {
			authToken = at.String()
		} else {
			return 0, ErrUnauthorized
		}
	} else {
		authToken = c.authToken
	}
	httpReq.Header.Set("Authentication", authToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data rgsReportByPartnerResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	if data.HasError && data.ErrorDescription != nil {
		return 0, errors.New(*data.ErrorDescription)
	}

	return data.Result.TotalProfitByReportCurrency, nil
}
