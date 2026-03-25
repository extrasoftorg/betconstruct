package backoffice

import "context"

type Client interface {
	ListTransactions(ctx context.Context, req ListTransactionsRequest) ([]Transaction, error)
	ListDeposits(ctx context.Context, req ListDepositsRequest) ([]Deposit, error)
	ListWithdrawals(ctx context.Context, req ListWithdrawalsRequest) ([]Withdrawal, error)

	ListRegisteredPlayers(ctx context.Context, req ListRegisteredPlayersRequest) ([]RegisteredPlayer, error)
	ListPlayers(ctx context.Context, req ListPlayersRequest) ([]ListPlayersPlayer, error)
	GetClientKPI(ctx context.Context, playerID PlayerID) (*PlayerKPI, error)
	SaveClientRestriction(ctx context.Context, req SaveClientRestrictionRequest) error
	AddPaymentToPlayer(ctx context.Context, req AddPaymentToPlayerRequest) error
	AddBonusToPlayer(ctx context.Context, req AddBonusToPlayerRequest) error
	ListPlayerTransactions(ctx context.Context, req ListPlayerTransactionsRequest) ([]Transaction, error)
	ListPlayerCasinoGames(ctx context.Context, req ListPlayerCasinoGamesRequest) ([]PlayerCasinoGame, error)

	GetSportKindReport(ctx context.Context, req GetSportKindReportRequest) ([]SportKindReport, error)
	ListSportBets(ctx context.Context, req ListSportBetsRequest) ([]SportBet, error)
	GetBetHistory(ctx context.Context, req ListSportBetsRequest) (*GetBetHistoryResult, error)
	GetCasinoReportByPartner(ctx context.Context, req GetReportByPartnerRequest) (float64, error)

	ListPaymentMethods(ctx context.Context, req ListPaymentMethodsRequest) ([]*PaymentMethod, error)
	FindPaymentMethodByName(ctx context.Context, name string) (*PaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, method PaymentMethod) error
	ListPartnerDomains(ctx context.Context, partnerID PartnerID) ([]PartnerDomain, error)
	SetActiveDomain(ctx context.Context, domainID int32) error
}
