package crm

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

func (c *client) Login(ctx context.Context) error {
	body, err := json.Marshal(c.betconstructToken)
	if err != nil {
		return err
	}
	authToken, err := makeRequest[string](ctx, http.MethodPost, "/User/LoginWithPlatform", bytes.NewReader(body), c, nil)
	if err != nil {
		return err
	}

	token := strings.TrimPrefix(*authToken, "Bearer ")
	c.authToken = token

	return nil
}
