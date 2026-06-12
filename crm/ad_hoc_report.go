package crm

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type AdHocReportColumnRef struct {
	Id int32 `json:"Id"`
}

type AdHocReportTableRef struct {
	Id int32 `json:"Id"`
}

type AdHocReportFilterValueData struct {
	Value              any     `json:"Value"` // string (e.g. "2026-06-08T00:00:00.000Z") or bool
	HasTime            bool    `json:"HasTime"`
	TimeZone           *string `json:"TimeZone"`
	ValueType          int     `json:"ValueType"`
	CascadeFiltersList any     `json:"CascadeFiltersList"` // null in practice
}

type AdHocReportFilterValueArgument struct {
	Value              *string `json:"Value"`
	HasTime            bool    `json:"HasTime"`
	TimeZone           *string `json:"TimeZone"`
	ValueType          int     `json:"ValueType"`
	CascadeFiltersList any     `json:"CascadeFiltersList"`
}

type AdHocReportFilterValue struct {
	Data      AdHocReportFilterValueData     `json:"Data"`
	Operation int                            `json:"Operation"`
	Argument  AdHocReportFilterValueArgument `json:"Argument"`
}

type AdHocReportFilter struct {
	Column      AdHocReportColumnRef   `json:"Column"`
	Comparision int                    `json:"Comparision"`
	Value       AdHocReportFilterValue `json:"Value"`
}

type AdHocReportQueryConfiguration struct {
	Columns         []AdHocReportColumnRef `json:"Columns"`
	Filters         [][]AdHocReportFilter  `json:"Filters"`
	Table           AdHocReportTableRef    `json:"Table"`
	UnifiedCurrency *string                `json:"UnifiedCurrency"`
}

type CreateAdHocReportRequest struct {
	Name               string                        `json:"Name"`
	Description        string                        `json:"Description"`
	QueryConfiguration AdHocReportQueryConfiguration `json:"QueryConfiguration"`
	DefinitionType     int                           `json:"DefinitionType"`
}

type AdHocReport struct {
	Id   int32  `json:"AdHocReportId"`
	Name string `json:"Name"`
}

// NewAdHocReportFilter builds one filter (one bound). For a date range, the
// caller emits two: Comparision 0 for the lower bound, Comparision 7 for the upper.
func NewAdHocReportFilter(columnID int32, comparision int, value any, hasTime bool, valueType int) AdHocReportFilter {
	return AdHocReportFilter{
		Column:      AdHocReportColumnRef{Id: columnID},
		Comparision: comparision,
		Value: AdHocReportFilterValue{
			Data: AdHocReportFilterValueData{
				Value:              value,
				HasTime:            hasTime,
				TimeZone:           nil,
				ValueType:          valueType,
				CascadeFiltersList: nil,
			},
			Operation: 0,
			Argument: AdHocReportFilterValueArgument{
				Value:              nil,
				HasTime:            false,
				TimeZone:           nil,
				ValueType:          0,
				CascadeFiltersList: nil,
			},
		},
	}
}

func (c *client) CreateAdHocReport(ctx context.Context, req CreateAdHocReportRequest) (*AdHocReport, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return makeRequest[AdHocReport](ctx, http.MethodPost, "/AdHocReport/Create", bytes.NewReader(body), c, nil)
}
