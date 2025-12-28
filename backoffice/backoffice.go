package backoffice

import "context"

type Client interface {
	ListTransactions(ctx context.Context, req ListTransactionsRequest) ([]Transaction, error)
	ListDeposits(ctx context.Context, req ListDepositsRequest) ([]Deposit, error)
	ListWithdrawals(ctx context.Context, req ListWithdrawalsRequest) ([]Withdrawal, error)

	ListRegisteredPlayers(ctx context.Context, req ListRegisteredPlayersRequest) ([]RegisteredPlayer, error)
	ListPlayers(ctx context.Context, req ListPlayersRequest) ([]ListPlayersPlayer, error)
	AddPaymentToPlayer(ctx context.Context, req AddPaymentToPlayerRequest) error

	ListSportBets(ctx context.Context, req ListSportBetsRequest) ([]SportBet, error)
}
