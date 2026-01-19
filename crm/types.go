package crm

import (
	"time"
)

type ReportResult struct {
	ID          int32     `json:"AdHocReportResultId"`
	CreatorName string    `json:"CreatorName"`
	CreatedAt   time.Time `json:"CreatedDate"`
}
