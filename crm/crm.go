package crm

import "context"

type Client interface {
	Login(ctx context.Context) error

	CreateAdHocReport(ctx context.Context, req CreateAdHocReportRequest) (*AdHocReport, error)

	ListReports(ctx context.Context) ([]Report, error)

	ExecuteReport(ctx context.Context, reportID int32) error

	DownloadReportAsExcel(ctx context.Context, reportResultID int32) ([]byte, error)

	ListReportResults(ctx context.Context, reportID int32) ([]ReportResult, error)
}
