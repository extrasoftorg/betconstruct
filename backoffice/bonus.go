package backoffice

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
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
