package crm

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type listReportResultsRequest struct {
	ReportID    string         `json:"ReportId"`
	Type        string         `json:"Type"`
	SearchModel map[string]any `json:"Searchmodel"`
}

type listReportResultsResponse struct {
	Data []ReportResult `json:"Data"`
}

func (c *client) ListReportResults(ctx context.Context, reportID int32) ([]ReportResult, error) {
	reportIDStr := strconv.Itoa(int(reportID))
	body, err := json.Marshal(listReportResultsRequest{
		ReportID: reportIDStr,
		Type:     "0",
		SearchModel: map[string]any{
			"Filters": []map[string]any{
				{
					"Comparision": 2,
					"Name":        "Name",
					"Values":      []string{""},
				},
			},
			"Pageing": map[string]any{
				"PageSize":   20,
				"PageNumber": 1,
			},
			"Sorting": map[string]any{
				"Name":      "ArchivedDate",
				"Direction": "asc",
			},
			"SortingThen": map[string]any{
				"Name":      "CreatedDate",
				"Direction": "desc",
			},
		},
	})
	if err != nil {
		return nil, err
	}

	results, err := makeRequest[listReportResultsResponse](ctx, http.MethodPost, "/ReportResult/List", bytes.NewReader(body), c, nil)
	if err != nil {
		return nil, err
	}
	return results.Data, nil
}
