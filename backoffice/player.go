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
		makeRequestOptions{
			httpClient: c.httpClient,
			authToken:  c.authToken,
		},
	)
	if err != nil {
		return nil, err
	}
	return *players, nil
}

type ListPlayersRequestDate struct {
	time.Time
}

func (d ListPlayersRequestDate) MarshalJSON() ([]byte, error) {
	layout := "02-01-06 - 15:04:05"
	return json.Marshal(time.Time(d.Time).Format(layout))
}

type ListPlayersRequest struct {
	FromDate *ListPlayersRequestDate `json:"MinCreatedLocal"`
	ToDate   *ListPlayersRequestDate `json:"MaxCreatedLocal"`
	MaxRows  int                     `json:"MaxRows"`
	Username string                  `json:"Login"`
}

type listPlayersResponse struct {
	Players []ListPlayersPlayer `json:"Objects"`
}

func (c *client) ListPlayers(ctx context.Context, req ListPlayersRequest) ([]ListPlayersPlayer, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	players, err := makeRequest[listPlayersResponse](
		ctx,
		http.MethodPost,
		"/Client/GetClients",
		bytes.NewReader(body),
		makeRequestOptions{
			httpClient: c.httpClient,
			authToken:  c.authToken,
		},
	)
	if err != nil {
		return nil, err
	}
	return players.Players, nil
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
		makeRequestOptions{
			httpClient: c.httpClient,
			authToken:  c.authToken,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
