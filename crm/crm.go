package crm

import "context"

type Client interface {
	Login(ctx context.Context) error
	AuthToken() string

	CreateReport(ctx context.Context, in CreateReportInput) (*CreateReportResponse, error)
	ExecuteReport(ctx context.Context, reportID int32) error

	ListReports(ctx context.Context, pageSize, pageNumber int) ([]Report, error)

	DownloadReportAsExcel(ctx context.Context, reportResultID int32) ([]byte, error)

	ListReportResults(ctx context.Context, reportID int32) ([]ReportResult, error)
}
