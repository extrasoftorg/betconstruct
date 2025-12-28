package backoffice

import "context"

type Client interface {
	ListTransactions(ctx context.Context, req ListTransactionsRequest) ([]Transaction, error)
	ListDeposits(ctx context.Context, req ListDepositsRequest) ([]Deposit, error)
	ListWithdrawals(ctx context.Context, req ListWithdrawalsRequest) ([]Withdrawal, error)

	ListRegisteredPlayers(ctx context.Context, req ListRegisteredPlayersRequest) ([]RegisteredPlayer, error)
}
