package crm

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type downloadReportAsExcelRequest struct {
	ReportResultID int32  `json:"ReportResultId"`
	Currency       string `json:"CurrencyCode"`
	DocumentType   string `json:"DocumentType"`
	ReportType     string `json:"ReportType"`
	TZ             int    `json:"UserTimeZone"`
	FileName       string `json:"fileName"`
}

func (c *client) DownloadReportAsExcel(ctx context.Context, reportResultID int32) ([]byte, error) {
	req := downloadReportAsExcelRequest{
		ReportResultID: reportResultID,
		Currency:       "TRY",
		DocumentType:   "xlsx",
		ReportType:     "0",
		TZ:             3,
		FileName:       "report",
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var bts []byte
	if _, err := makeRequest[[]byte](ctx, http.MethodPost, "/AdHocReportResult/GetExcel", bytes.NewReader(body), c, func(r io.Reader) error {
		bts, err = io.ReadAll(r)
		return err
	}); err != nil {
		return nil, err
	}

	return bts, nil
}
