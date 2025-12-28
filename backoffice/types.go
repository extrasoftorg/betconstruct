package backoffice

import (
	"encoding/json"
	"time"
)

type PlayerID int64

func (p PlayerID) Int64() int64 {
	return int64(p)
}

type DateTime time.Time

func (t *DateTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}
	*t = DateTime(parsed)
	return nil
}

type TransactionType string

func (t TransactionType) String() string {
	return string(t)
}

func (t *TransactionType) UnmarshalJSON(b []byte) error {
	var id int
	if err := json.Unmarshal(b, &id); err != nil {
		return err
	}
	switch id {
	case 15:
		*t = TransactionTypeWinning
	case 10:
		*t = TransactionTypeBet
	case 301:
		*t = TransactionTypeCorrectionUp
	case 302:
		*t = TransactionTypeCorrectionDown
	case 3:
		*t = TransactionTypeDeposit
	default:
		*t = TransactionTypeUnknown
	}
	return nil
}

const (
	TransactionTypeWinning        TransactionType = "winning"
	TransactionTypeBet            TransactionType = "bet"
	TransactionTypeCorrectionUp   TransactionType = "correctionUp"
	TransactionTypeCorrectionDown TransactionType = "correctionDown"
	TransactionTypeDeposit        TransactionType = "deposit"
	TransactionTypeUnknown        TransactionType = "unknown"
)

type Transaction struct {
	ID        int64           `json:"Id"`
	Amount    float64         `json:"Amount"`
	PlayerID  PlayerID        `json:"ClientId"`
	Type      TransactionType `json:"TypeId"`
	Note      *string         `json:"Note"`
	CreatedAt DateTime        `json:"CreatedLocal"`
}

type Deposit struct {
	ID            int64    `json:"Id"`
	Amount        float64  `json:"Amount"`
	PlayerID      PlayerID `json:"ClientId"`
	CreatedAt     DateTime `json:"CreatedLocal"`
	PaymentMethod string   `json:"PaymentSystemName"`
}

type WithdrawalStatus string

func (w WithdrawalStatus) String() string {
	return string(w)
}

func (w *WithdrawalStatus) UnmarshalJSON(b []byte) error {
	var id int
	if err := json.Unmarshal(b, &id); err != nil {
		return err
	}
	switch id {
	case 0:
		*w = WithdrawalStatusPending
	case 3:
		*w = WithdrawalStatusPaid
	case -2:
		*w = WithdrawalStatusRejected
	case -1:
		*w = WithdrawalStatusCancelled
	default:
		*w = WithdrawalStatusUnknown
	}
	return nil
}

const (
	WithdrawalStatusPending   WithdrawalStatus = "pending"
	WithdrawalStatusPaid      WithdrawalStatus = "paid"
	WithdrawalStatusRejected  WithdrawalStatus = "rejected"
	WithdrawalStatusCancelled WithdrawalStatus = "cancelled"
	WithdrawalStatusUnknown   WithdrawalStatus = "unknown"
)

type Withdrawal struct {
	ID            int64            `json:"Id"`
	Amount        float64          `json:"Amount"`
	PlayerID      PlayerID         `json:"ClientId"`
	RequestedAt   DateTime         `json:"RequestTimeLocal"`
	PaymentMethod string           `json:"PaymentSystemName"`
	AllowedAt     *DateTime        `json:"AllowTimeLocal"`
	Info          string           `json:"Info"`
	Status        WithdrawalStatus `json:"State"`
}

type RegisteredPlayer struct {
	ID        PlayerID `json:"ClientId"`
	CreatedAt DateTime `json:"CreatedLocal"`
	Username  string   `json:"Login"`
	FullName  string   `json:"Name"`
}

type ListPlayersPlayer struct {
	ID         PlayerID `json:"Id"`
	CreatedAt  DateTime `json:"CreatedLocalDate"`
	Username   string   `json:"Login"`
	FirstName  string   `json:"FirstName"`
	MiddleName string   `json:"MiddleName"`
	LastName   string   `json:"LastName"`
	Balance    float64  `json:"Balance"`
}

type SportBetStatus string

func (s SportBetStatus) String() string {
	return string(s)
}

func (s *SportBetStatus) UnmarshalJSON(b []byte) error {
	var id int
	if err := json.Unmarshal(b, &id); err != nil {
		return err
	}
	switch id {
	case 1:
		*s = SportBetStatusPending
	default:
		*s = SportBetStatusUnknown
	}
	return nil
}

const (
	SportBetStatusPending SportBetStatus = "pending"
	SportBetStatusUnknown SportBetStatus = "unknown"
)

type SportBet struct {
	Status SportBetStatus `json:"State"`
}
