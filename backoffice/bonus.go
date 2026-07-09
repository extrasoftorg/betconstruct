package backoffice

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type AddBonusToPlayerRequest struct {
	Amount   float64  `json:"Amount"`
	PlayerID PlayerID `json:"ClientId"`
	BonusID  int64    `json:"PartnerBonusId"`
	Type     int      `json:"Type"`
}

func (c *client) AddBonusToPlayer(ctx context.Context, req AddBonusToPlayerRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = makeRequest[any](
		ctx,
		http.MethodPost,
		"/Client/AddClientToBonus",
		bytes.NewReader(body),
		c,
	)
	if err != nil {
		return err
	}
	return nil
}

type BonusStatus string

const (
	BonusStatusCancelled BonusStatus = "cancelled"
	BonusStatusPending   BonusStatus = "pending"
	BonusStatusExpired   BonusStatus = "expired"
)

type PlayerBonus struct {
	ID        int64
	Amount    float64
	Name      string
	CreatedAt time.Time
	Status    BonusStatus
	BonusID   int64
}

func (b *PlayerBonus) UnmarshalJSON(data []byte) error {
	type wire struct {
		ID         int64   `json:"Id"`
		Amount     float64 `json:"Amount"`
		Name       string  `json:"Name"`
		CreatedAt  string  `json:"CreatedLocal"`
		Count      int     `json:"Count"`
		ResultType int     `json:"ResultType"`
		BonusID    int64   `json:"PartnerBonusId"`
	}
	var w wire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}

	b.ID = w.ID
	b.Name = w.Name
	b.Amount = w.Amount

	createdAt, err := time.Parse("2006-01-02T15:04:05.999999999", w.CreatedAt)
	if err != nil {
		return err
	}
	b.CreatedAt = createdAt

	if b.Amount <= 0 {
		b.Amount = float64(w.Count)
	}

	switch w.ResultType {
	// none
	case 0:
		b.Status = BonusStatusPending
	// cancelled
	case 3:
		b.Status = BonusStatusCancelled
	// expired
	case 4:
		b.Status = BonusStatusExpired
	}

	b.BonusID = w.BonusID

	return nil
}

func (c *client) ListPlayerBonuses(ctx context.Context, playerID PlayerID) ([]PlayerBonus, error) {
	type payload struct {
		PlayerID PlayerID `json:"ClientId"`
	}
	p := payload{
		PlayerID: playerID,
	}

	body, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	bonuses, err := makeRequest[[]PlayerBonus](
		ctx,
		http.MethodPost,
		"/Client/GetClientBonuses",
		bytes.NewReader(body),
		c,
	)
	if err != nil {
		return nil, err
	}
	return *bonuses, nil
}
