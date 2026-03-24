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

type adHocFilterValueData struct {
	HasTime   bool   `json:"HasTime"`
	Value     string `json:"Value"`
	ValueType int    `json:"ValueType"`
}

type adHocFilterValueArgument struct {
	Value     *string `json:"Value"`
	ValueType int     `json:"ValueType"`
}

type adHocFilterValueBound struct {
	Operation int                      `json:"Operation"`
	Data      adHocFilterValueData     `json:"Data"`
	Argument  adHocFilterValueArgument `json:"Argument"`
}

type AdHocReportFilterValue struct {
	From adHocFilterValueBound `json:"From"`
	To   adHocFilterValueBound `json:"To"`
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

func NewAdHocFilterValueBound(hasTime bool, value string, valueType int) adHocFilterValueBound {
	return adHocFilterValueBound{
		Operation: 0,
		Data: adHocFilterValueData{
			HasTime:   hasTime,
			Value:     value,
			ValueType: valueType,
		},
		Argument: adHocFilterValueArgument{
			Value:     nil,
			ValueType: 0,
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
