package crm

import (
	"fmt"
	"time"
)

type ReportColumn string

func (c ReportColumn) ID() int {
	return reportColumnIDMap[c]
}

const (
	ReportColumnPlayerID           ReportColumn = "PlayerId"
	ReportColumnRegistrationDate   ReportColumn = "RegistrationDate"
	ReportColumnLastDepositDate    ReportColumn = "LastDepositDate"
	ReportColumnTotalDepositAmount ReportColumn = "TotalDepositAmount"
	ReportColumnDepositActivity    ReportColumn = "DepositActivity"
)

var reportColumnIDMap = map[ReportColumn]int{
	ReportColumnPlayerID:           38,
	ReportColumnRegistrationDate:   61,
	ReportColumnLastDepositDate:    712,
	ReportColumnTotalDepositAmount: 223,
	ReportColumnDepositActivity:    634,
}

type ReportFilterOp string

const (
	ReportFilterOpEq  ReportFilterOp = "eq"
	ReportFilterOpLt  ReportFilterOp = "lt"
	ReportFilterOpLte ReportFilterOp = "lte"
)

var reportFilterOpMap = map[ReportFilterOp]int{
	ReportFilterOpEq:  0,
	ReportFilterOpLt:  6,
	ReportFilterOpLte: 7,
}

type ReportFilterValue interface {
	isReportFilterValue()
}

type ReportFilterValueBool struct {
	Value bool
}

func (ReportFilterValueBool) isReportFilterValue() {}

type ReportFilterValueAmount struct {
	Amount   float64
	Currency string
}

func (ReportFilterValueAmount) isReportFilterValue() {}

type ReportFilterValueDate struct {
	Time time.Time
}

func (ReportFilterValueDate) isReportFilterValue() {}

type ReportFilter struct {
	Column int
	Op     ReportFilterOp
	Value  ReportFilterValue
}

type CreateReportInput struct {
	Name    string
	Columns []int
	Filters []ReportFilter
}

type column struct {
	ID int `json:"Id"`
}

type queryConfiguration struct {
	Columns []column `json:"Columns"`
	Table   struct {
		ID int `json:"Id"`
	} `json:"Table"`
	Filters         [][]filter `json:"Filters"`
	UnifiedCurrency any        `json:"UnifiedCurrency"`
}

type filter struct {
	Column     column `json:"Column"`
	Comparison int    `json:"Comparision"`
	Value      any    `json:"Value"`
}

type filterValueArgument struct {
	Value any `json:"Value"`
	Type  int `json:"ValueType"`
}

type filterValueDateData struct {
	HasTime   bool      `json:"HasTime"`
	Value     time.Time `json:"Value"`
	ValueType int       `json:"ValueType"`
}
type filterValueDate struct {
	Operation int                 `json:"Operation"`
	Argument  filterValueArgument `json:"Argument"`
	Data      filterValueDateData `json:"Data"`
}

type filterValueBoolData struct {
	HasTime   bool `json:"HasTime"`
	Value     bool `json:"Value"`
	ValueType int  `json:"ValueType"`
}
type filterValueBool struct {
	Operation int                 `json:"Operation"`
	Argument  filterValueArgument `json:"Argument"`
	Data      filterValueBoolData `json:"Data"`
}

type filterValueAmountData struct {
	CurrencyCode  string  `json:"CurrencyCode"`
	OriginalValue float64 `json:"OriginalValue"`
	IsEquivalent  bool    `json:"IsEquivalent"`
}
type filterValueAmount struct {
	Data filterValueAmountData `json:"Data"`
}

type createReportPayload struct {
	Name               string             `json:"Name"`
	QueryConfiguration queryConfiguration `json:"QueryConfiguration"`
	DefinitionType     int                `json:"DefinitionType"`
	Description        string             `json:"Description"`
}

func (in CreateReportInput) toPayload() (*createReportPayload, error) {
	columns := make([]column, len(in.Columns))
	for i, col := range in.Columns {
		columns[i] = column{
			ID: col,
		}
	}

	filters := make([]filter, len(in.Filters))
	for i, filt := range in.Filters {
		var filterVal any
		switch val := filt.Value.(type) {
		case ReportFilterValueAmount:
			filterVal = filterValueAmount{
				Data: filterValueAmountData{
					CurrencyCode:  val.Currency,
					OriginalValue: val.Amount,
					IsEquivalent:  false,
				},
			}

		case ReportFilterValueDate:
			t := val.Time
			hasTime := false
			if t.Hour() != 0 || t.Minute() != 0 || t.Second() != 0 {
				hasTime = true
			}

			filterVal = filterValueDate{
				Argument: filterValueArgument{
					Value: nil,
					Type:  0,
				},
				Data: filterValueDateData{
					HasTime:   hasTime,
					Value:     t,
					ValueType: 0,
				},
				Operation: 0,
			}

		case ReportFilterValueBool:
			filterVal = filterValueBool{
				Operation: 0,
				Argument: filterValueArgument{
					Value: nil,
					Type:  0,
				},
				Data: filterValueBoolData{
					HasTime:   false,
					Value:     val.Value,
					ValueType: 0,
				},
			}

		default:
			return nil, fmt.Errorf("invalid filter value type: %T", val)
		}

		filters[i] = filter{
			Column:     column{ID: filt.Column},
			Comparison: reportFilterOpMap[filt.Op],
			Value:      filterVal,
		}
	}

	return &createReportPayload{
		Name: in.Name,
		QueryConfiguration: queryConfiguration{
			Columns: columns,
			Table: struct {
				ID int `json:"Id"`
			}{
				ID: 4,
			},
			Filters: [][]filter{filters},
		},
	}, nil
}

type CreateReportResponse struct {
	ReportID int32
}
