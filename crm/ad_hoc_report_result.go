package crm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var ErrReportResultNotReady = errors.New("report result not ready")

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
		if err != nil {
			return err
		}
		var data response[any]
		if err := json.Unmarshal(bts, &data); err == nil {
			return ErrReportResultNotReady
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return bts, nil
}
