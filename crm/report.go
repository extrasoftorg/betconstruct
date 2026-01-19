package crm

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type executeReportRequest struct {
	Type     int   `json:"Type"`
	ReportID int32 `json:"CustomReportId"`
}

func (c *client) ExecuteReport(ctx context.Context, reportID int32) error {
	body, err := json.Marshal(executeReportRequest{
		ReportID: reportID,
	})
	if err != nil {
		return err
	}
	_, err = makeRequest[any](ctx, http.MethodPost, "/Report/Execute", bytes.NewReader(body), c, nil)
	return err
}
