package crm

import "context"

type Client interface {
	Login(ctx context.Context) error

	ExecuteReport(ctx context.Context, reportID int32) error

	DownloadReportAsExcel(ctx context.Context, reportResultID int32) ([]byte, error)

	ListReportResults(ctx context.Context, reportID int32) ([]ReportResult, error)
}
