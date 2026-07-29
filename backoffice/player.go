package backoffice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ListRegisteredPlayersRequestDate struct {
	time.Time
}

func (d ListRegisteredPlayersRequestDate) MarshalJSON() ([]byte, error) {
	layout := "02-01-06"
	return json.Marshal(time.Time(d.Time).Format(layout))
}

type ListRegisteredPlayersRequest struct {
	Date time.Time `json:"DateLocal"`
}

type listRegisteredPlayersResponse []RegisteredPlayer

func (c *client) ListRegisteredPlayers(ctx context.Context, req ListRegisteredPlayersRequest) ([]RegisteredPlayer, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	players, err := makeRequest[listRegisteredPlayersResponse](
		ctx,
		http.MethodPost,
		"/Client/GetClientRegistrationStatisticsDetails",
		bytes.NewReader(body),
		c,
	)
	if err != nil {
		return nil, err
	}
	return *players, nil
}

type ListPlayersRequest struct {
	FromRegistrationDate time.Time
	ToRegistrationDate   time.Time
	MaxRows              int
	Username             string
	Phone                string
}

func (c *client) ListPlayers(ctx context.Context, req ListPlayersRequest) ([]*ListPlayersPlayer, error) {
	type payload struct {
		FromRegistrationDate *string `json:"MinCreatedLocal"`
		ToRegistrationDate   *string `json:"MaxCreatedLocal"`
		MaxRows              int     `json:"MaxRows"`
		Username             string  `json:"Login"`
		Phone                string  `json:"Phone"`
	}
	p := payload{
		MaxRows:  req.MaxRows,
		Username: req.Username,
		Phone:    req.Phone,
	}
	if !req.FromRegistrationDate.IsZero() {
		date := req.FromRegistrationDate.Format("02-01-06 - 15:04:05")
		p.FromRegistrationDate = &date
	}
	if !req.ToRegistrationDate.IsZero() {
		date := req.ToRegistrationDate.Format("02-01-06 - 15:04:05")
		p.ToRegistrationDate = &date
	}

	body, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	type player struct {
		ID             PlayerID `json:"Id"`
		CreatedAt      string   `json:"CreatedLocalDate"`
		Username       string   `json:"Login"`
		FirstName      string   `json:"FirstName"`
		MiddleName     string   `json:"MiddleName"`
		LastName       string   `json:"LastName"`
		Balance        float64  `json:"Balance"`
		PlayerCategory int      `json:"SportsbookProfileId"`
		FirstDepositAt *string  `json:"FirstDepositDateLocal"`
	}
	type response struct {
		Players []*player `json:"Objects"`
	}
	players, err := makeRequest[response](
		ctx,
		http.MethodPost,
		"/Client/GetClients",
		bytes.NewReader(body),
		c,
	)
	if err != nil {
		return nil, err
	}

	res := make([]*ListPlayersPlayer, len(players.Players))
	for i, p := range players.Players {
		var pc *PlayerCategory
		if p.PlayerCategory != 0 {
			cat := PlayerCategory(p.PlayerCategory)
			pc = &cat
		}

		createdAt, err := time.Parse("2006-01-02T15:04:05.999999", p.CreatedAt)
		if err != nil {
			return nil, err
		}

		var firstDepositAt *time.Time
		if p.FirstDepositAt != nil {
			date, err := time.Parse("2006-01-02T15:04:05.999999", *p.FirstDepositAt)
			if err != nil {
				return nil, err
			}
			firstDepositAt = &date
		}

		res[i] = &ListPlayersPlayer{
			ID:             p.ID,
			CreatedAt:      createdAt,
			Username:       p.Username,
			FirstName:      p.FirstName,
			MiddleName:     p.MiddleName,
			LastName:       p.LastName,
			Balance:        p.Balance,
			PlayerCategory: pc,
			FirstDepositAt: firstDepositAt,
		}
	}

	return res, nil
}

func (c *client) GetPlayerKPI(ctx context.Context, playerID PlayerID) (*PlayerKPI, error) {
	type response struct {
		PlayerID              PlayerID `json:"ClientId"`
		TotalDepositAmount    float64  `json:"DepositAmount"`
		DepositCount          int32    `json:"DepositCount"`
		TotalWithdrawalAmount float64  `json:"WithdrawalAmount"`
		WithdrawalCount       int32    `json:"WithdrawalCount"`
		LastDepositAmount     float64  `json:"LastDepositAmount"`
		LastWithdrawalAmount  float64  `json:"LastWithdrawalAmount"`
		FirstDepositAt        *string  `json:"FirstDepositTimeLocal"`
		LastDepositAt         *string  `json:"LastDepositTimeLocal"`
		LastWithdrawalAt      *string  `json:"LastWithdrawalTimeLocal"`
		LastSportBetAt        *string  `json:"LastSportBetTimeLocal"`
		LastCasinoBetAt       *string  `json:"LastCasinoBetTimeLocal"`
	}

	resp, err := makeRequest[response](
		ctx,
		http.MethodGet,
		fmt.Sprintf("/Client/GetClientKpi?id=%d", playerID),
		nil,
		c,
	)
	if err != nil {
		return nil, err
	}

	kpi := &PlayerKPI{
		PlayerID:              resp.PlayerID,
		TotalDepositAmount:    resp.TotalDepositAmount,
		DepositCount:          resp.DepositCount,
		TotalWithdrawalAmount: resp.TotalWithdrawalAmount,
		WithdrawalCount:       resp.WithdrawalCount,
		LastDepositAmount:     resp.LastDepositAmount,
		LastWithdrawalAmount:  resp.LastWithdrawalAmount,
	}

	if resp.FirstDepositAt != nil {
		date, err := time.Parse("2006-01-02T15:04:05.999999", *resp.FirstDepositAt)
		if err != nil {
			return nil, err
		}
		kpi.FirstDepositAt = &date
	}
	if resp.LastDepositAt != nil {
		date, err := time.Parse("2006-01-02T15:04:05.999999", *resp.LastDepositAt)
		if err != nil {
			return nil, err
		}
		kpi.LastDepositAt = &date
	}
	if resp.LastWithdrawalAt != nil {
		date, err := time.Parse("2006-01-02T15:04:05.999999", *resp.LastWithdrawalAt)
		if err != nil {
			return nil, err
		}
		kpi.LastWithdrawalAt = &date
	}
	if resp.LastSportBetAt != nil {
		date, err := time.Parse("2006-01-02T15:04:05.999999", *resp.LastSportBetAt)
		if err != nil {
			return nil, err
		}
		kpi.LastSportBetAt = &date
	}
	if resp.LastCasinoBetAt != nil {
		date, err := time.Parse("2006-01-02T15:04:05.999999", *resp.LastCasinoBetAt)
		if err != nil {
			return nil, err
		}
		kpi.LastCasinoBetAt = &date
	}

	return kpi, nil
}

type GetClientRestrictionResult struct {
	PlayerID                PlayerID  `json:"ClientId"`
	CanLogin                bool      `json:"CanLogin"`
	CanBet                  bool      `json:"CanBet"`
	CanDeposit              bool      `json:"CanDeposit"`
	CanWithdraw             bool      `json:"CanWithdraw"`
	CanIncreaseLimit        bool      `json:"CanIncreaseLimit"`
	CanClaimBonus           bool      `json:"CanClaimBonus"`
	CanCasinoLogin          bool      `json:"CanCasinoLogin"`
	CanUploadDocument       bool      `json:"CanUploadDocument"`
	IsWithdrawalAutoConfirm bool      `json:"IsWithdrawalAutoConfirm"`
	UpdatedAt               *DateTime `json:"ModifedLocal"`
}

func (c *client) GetClientRestriction(ctx context.Context, playerID PlayerID) (*GetClientRestrictionResult, error) {
	restriction, err := makeRequest[GetClientRestrictionResult](
		ctx,
		http.MethodGet,
		fmt.Sprintf("/Client/GetClientRestriction?clientId=%d", playerID),
		nil,
		c,
	)
	if err != nil {
		return nil, err
	}
	return restriction, nil
}

type SaveClientRestrictionRequest struct {
	PlayerID          PlayerID `json:"ClientId"`
	CanLogin          bool     `json:"CanLogin"`
	CanBet            bool     `json:"CanBet"`
	CanDeposit        bool     `json:"CanDeposit"`
	CanWithdraw       bool     `json:"CanWithdraw"`
	CanIncreaseLimit  bool     `json:"CanIncreaseLimit"`
	CanClaimBonus     bool     `json:"CanClaimBonus"`
	CanCasinoLogin    bool     `json:"CanCasinoLogin"`
	CanUploadDocument bool     `json:"CanUploadDocument"`
	UserName          string   `json:"UserName,omitempty"`
}

func (c *client) SaveClientRestriction(ctx context.Context, req SaveClientRestrictionRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = makeRequest[any](
		ctx,
		http.MethodPost,
		"/Client/SaveClientRestriction",
		bytes.NewReader(body),
		c,
	)
	return err
}

type AddPaymentToPlayerRequestType string

const (
	AddPaymentToPlayerRequestTypeCorrectionUp AddPaymentToPlayerRequestType = "correctionUp"
)

type AddPaymentToPlayerRequest struct {
	PlayerID PlayerID
	Amount   float64
	Note     string
	Type     AddPaymentToPlayerRequestType
	Currency string
}

func (a AddPaymentToPlayerRequest) MarshalJSON() ([]byte, error) {
	w := struct {
		PlayerID PlayerID `json:"ClientId"`
		Amount   string   `json:"Amount"`
		Note     string   `json:"Info"`
		Type     int      `json:"DocTypeInt"`
		Currency string   `json:"CurrencyId"`
	}{
		PlayerID: a.PlayerID,
		Note:     a.Note,
		Currency: a.Currency,
		Amount:   fmt.Sprintf("%.2f", a.Amount),
	}
	switch a.Type {
	case AddPaymentToPlayerRequestTypeCorrectionUp:
		w.Type = 3
	default:
		return nil, fmt.Errorf("unknown type: %s", a.Type)
	}
	return json.Marshal(w)
}

func (c *client) AddPaymentToPlayer(ctx context.Context, req AddPaymentToPlayerRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = makeRequest[any](
		ctx,
		http.MethodPost,
		"/Client/CreateClientPaymentDocument",
		bytes.NewReader(body),
		c,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) GetPlayer(ctx context.Context, playerID PlayerID) (*Player, error) {
	return makeRequest[Player](
		ctx,
		http.MethodGet,
		fmt.Sprintf("/Client/GetClientById?id=%d", playerID),
		nil,
		c,
	)
}
