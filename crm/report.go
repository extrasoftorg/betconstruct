package crm

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
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

type Report struct {
	Id          int32     `json:"ReportId"`
	Name        string    `json:"Name"`
	CreatorName string    `json:"CreatorName"`
	State       int       `json:"State"`
	ReportType  int       `json:"ReportType"`
	CreatedAt   time.Time `json:"CreatedDate"`
	HasResults  bool      `json:"HasResults"`
	IsExported  bool      `json:"IsExported"`
}

type listReportsRequest struct {
	Filters     []map[string]any `json:"Filters"`
	Pageing     map[string]any   `json:"Pageing"`
	Sorting     map[string]any   `json:"Sorting"`
	SortingThen map[string]any   `json:"SortingThen"`
}

type listReportsResponse struct {
	Data []Report `json:"Data"`
}

func (c *client) ListReports(ctx context.Context, pageSize, pageNumber int) ([]Report, error) {
	body, err := json.Marshal(listReportsRequest{
		Filters: []map[string]any{
			{"Name": "State", "Comparision": 0, "Values": []string{}},
		},
		Pageing:     map[string]any{"PageSize": pageSize, "PageNumber": pageNumber},
		Sorting:     map[string]any{"Name": "ArchivedDate", "Direction": "asc"},
		SortingThen: map[string]any{"Name": "CreatedDate", "Direction": "desc"},
	})
	if err != nil {
		return nil, err
	}

	results, err := makeRequest[listReportsResponse](ctx, http.MethodPost, "/Report/List", bytes.NewReader(body), c, nil)
	if err != nil {
		return nil, err
	}
	return results.Data, nil
}
