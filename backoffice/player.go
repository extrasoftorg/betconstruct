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
}

func (r ListPlayersRequest) MarshalJSON() ([]byte, error) {
	type wire struct {
		FromRegistrationDate *string `json:"MinCreatedLocal"`
		ToRegistrationDate   *string `json:"MaxCreatedLocal"`
		MaxRows              int     `json:"MaxRows"`
		Username             string  `json:"Login"`
	}
	w := wire{
		MaxRows:  r.MaxRows,
		Username: r.Username,
	}
	if !r.FromRegistrationDate.IsZero() {
		date := r.FromRegistrationDate.Format("02-01-06 - 15:04:05")
		w.FromRegistrationDate = &date
	}
	if !r.ToRegistrationDate.IsZero() {
		date := r.ToRegistrationDate.Format("02-01-06 - 15:04:05")
		w.ToRegistrationDate = &date
	}
	return json.Marshal(w)
}

func (c *client) ListPlayers(ctx context.Context, req ListPlayersRequest) ([]*ListPlayersPlayer, error) {
	body, err := req.MarshalJSON()
	if err != nil {
		return nil, err
	}
	players, err := makeRequest[struct {
		Players []*ListPlayersPlayer `json:"Objects"`
	}](
		ctx,
		http.MethodPost,
		"/Client/GetClients",
		bytes.NewReader(body),
		c,
	)
	if err != nil {
		return nil, err
	}
	return players.Players, nil
}

func (c *client) GetClientKPI(ctx context.Context, playerID PlayerID) (*PlayerKPI, error) {
	kpi, err := makeRequest[PlayerKPI](
		ctx,
		http.MethodGet,
		fmt.Sprintf("/Client/GetClientKpi?id=%d", playerID),
		nil,
		c,
	)
	if err != nil {
		return nil, err
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

type AddPaymentToPlayerRequestAmount struct {
	Value float64
}

func (a AddPaymentToPlayerRequestAmount) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("%.2f", a.Value)
	return json.Marshal(str)
}

type AddPaymentToPlayerRequestType string

func (t AddPaymentToPlayerRequestType) MarshalJSON() ([]byte, error) {
	var v int
	switch t {
	case AddPaymentToPlayerRequestTypeCorrectionUp:
		v = 3
	default:
		return nil, fmt.Errorf("unknown type: %s", t)
	}
	return json.Marshal(v)
}

const (
	AddPaymentToPlayerRequestTypeCorrectionUp AddPaymentToPlayerRequestType = "correctionUp"
)

type AddPaymentToPlayerRequest struct {
	PlayerID PlayerID                        `json:"ClientId"`
	Amount   AddPaymentToPlayerRequestAmount `json:"Amount"`
	Note     string                          `json:"Info"`
	Type     AddPaymentToPlayerRequestType   `json:"DocTypeInt"`
	Currency string                          `json:"CurrencyId"`
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
