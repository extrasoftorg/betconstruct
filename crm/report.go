package crm

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type ExecuteReportRequest struct {
	Type     int   `json:"Type"`
	ReportID int32 `json:"CustomReportId"`
}

func (c *client) ExecuteReport(ctx context.Context, req ExecuteReportRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = makeRequest[any](ctx, http.MethodPost, "/Report/Execute", bytes.NewReader(body), c, nil)
	return err
}
