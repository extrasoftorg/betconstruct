package backoffice

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
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

var sportBetStatusIDMap = map[SportBetStatus]int{
	SportBetStatusPending: 1,
}

func (s SportBetStatus) MarshalJSON() ([]byte, error) {
	id, ok := sportBetStatusIDMap[s]
	if !ok {
		return nil, fmt.Errorf("unknown sport bet status: %s", s)
	}
	return json.Marshal(id)
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

type NumericBool bool

func (n NumericBool) Bool() bool {
	return bool(n)
}

func (n NumericBool) MarshalJSON() ([]byte, error) {
	var v int
	switch n {
	case true:
		v = 1
	case false:
		v = 0
	}
	return json.Marshal(v)
}

func (n *NumericBool) UnmarshalJSON(b []byte) error {
	var v int
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*n = NumericBool(v != 0)
	return nil
}

type StringFloat64 float64

func (s StringFloat64) Float64() float64 {
	return float64(s)
}

func (s StringFloat64) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%.2f", s))
}

var ErrInvalidStringFloat64 = errors.New("invalid string float64")

func (s *StringFloat64) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	if f, ok := v.(float64); ok {
		*s = StringFloat64(f)
		return nil
	}

	if str, ok := v.(string); ok {
		str = strings.TrimSpace(str)

		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		*s = StringFloat64(f)
		return nil
	}

	return ErrInvalidStringFloat64
}

type PaymentMethodCurrencyProcessType string

func (p PaymentMethodCurrencyProcessType) String() string {
	return string(p)
}

const (
	PaymentMethodCurrencyProcessTypeText  PaymentMethodCurrencyProcessType = "text"
	PaymentMethodCurrencyProcessTypeHours PaymentMethodCurrencyProcessType = "hours"
)

type PaymentMethodConfigCurrency struct {
	Currency    string                           `json:"currency"`
	Fee         float64                          `json:"fee"`
	Min         StringFloat64                    `json:"min"`
	Max         StringFloat64                    `json:"max"`
	ProcessTime any                              `json:"process_time"`
	ProcessType PaymentMethodCurrencyProcessType `json:"process_type"`
	Blocked     bool                             `json:"blocked"`
	Enabled     bool                             `json:"enabled"`
	HasError    bool                             `json:"error"`
}

type PaymentMethodConfigFieldType string

func (t PaymentMethodConfigFieldType) String() string {
	return string(t)
}

const (
	PaymentMethodConfigFieldTypeSelect PaymentMethodConfigFieldType = "select"
)

type PaymentMethodConfigFieldOption struct {
	Value string `json:"value"`
	Label string `json:"text"`
}

type PaymentMethodConfigFieldOptions []PaymentMethodConfigFieldOption

var ErrInvalidPaymentMethodConfigFieldOption = errors.New("invalid payment method config field options")

func (f *PaymentMethodConfigFieldOptions) UnmarshalJSON(b []byte) error {
	opts := make([]PaymentMethodConfigFieldOption, 0)
	if err := json.Unmarshal(b, &opts); err != nil {
		var v any
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}

		if m, ok := v.(map[string]any); ok {
			for _, value := range m {
				if m2, ok := value.(map[string]any); ok {
					value, ok := m2["value"].(string)
					if !ok {
						return ErrInvalidPaymentMethodConfigFieldOption
					}
					label, ok := m2["text"].(string)
					if !ok {
						return ErrInvalidPaymentMethodConfigFieldOption
					}
					opts = append(opts, PaymentMethodConfigFieldOption{
						Value: value,
						Label: label,
					})
				}
			}
			return nil
		}

		return ErrInvalidPaymentMethodConfigFieldOption
	}
	*f = opts

	return nil
}

type PaymentMethodConfigField struct {
	Name     string                           `json:"name"`
	Label    string                           `json:"label"`
	Type     PaymentMethodConfigFieldType     `json:"type"`
	Options  *PaymentMethodConfigFieldOptions `json:"options"`
	Required bool                             `json:"required"`
}

type PaymentMethodConfig struct {
	Currencies []PaymentMethodConfigCurrency `json:"currencies"`
	Fields     []PaymentMethodConfigField    `json:"fields"`
}

type PaymentMethodDepositField struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Value string `json:"value"`
}

type PaymentMethodWithdrawField struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Value string `json:"value"`
}

type PaymentMethodGroupIDs []int8

func (g PaymentMethodGroupIDs) MarshalJSON() ([]byte, error) {
	return json.Marshal([]int8(g))
}

var ErrInvalidPaymentMethodGroupIDs = errors.New("invalid payment method group ids")

func (g *PaymentMethodGroupIDs) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	if slice, ok := v.([]any); ok {
		for _, v := range slice {
			if id, ok := v.(int8); ok {
				*g = append(*g, id)
			}
		}
		return nil
	}

	if _, ok := v.(string); ok {
		*g = make(PaymentMethodGroupIDs, 0)
		return nil
	}

	return ErrInvalidPaymentMethodGroupIDs
}

type PartnerID int32

type PaymentMethod struct {
	ID                     int16                        `json:"system_id"`
	SiteID                 int32                        `json:"site_system_id"`
	CanDeposit             bool                         `json:"can_deposit"`
	CanWithdraw            bool                         `json:"can_withdraw"`
	Name                   string                       `json:"system_name"`
	PartnerID              PartnerID                    `json:"partner_id"`
	IsActive               NumericBool                  `json:"is_active"`
	DepositConfig          PaymentMethodConfig          `json:"deposit"`
	WithdrawConfig         PaymentMethodConfig          `json:"withdraw"`
	BackendID              int16                        `json:"backend_id"`
	AllowCountries         []string                     `json:"allow_countries"`
	DepositFields          []PaymentMethodDepositField  `json:"deposit_config_fields"`
	WithdrawFields         []PaymentMethodWithdrawField `json:"withdraw_config_fields"`
	PartnerName            string                       `json:"site_name"`
	DepositInfoTextKey     string                       `json:"deposit_info_text_key"`
	WithdrawInfoTextKey    string                       `json:"withdraw_info_text_key"`
	IsManual               NumericBool                  `json:"is_manual"`
	Order                  int16                        `json:"order"`
	RestrictedCountries    []string                     `json:"restrict_countries"`
	DepositPrefill         bool                         `json:"deposit_prefilled_amount"`
	WithdrawPrefill        bool                         `json:"withdraw_prefilled_amount"`
	GroupIDs               PaymentMethodGroupIDs        `json:"payment_group"`
	StayInSameTabOnDeposit string                       `json:"stay_in_same_tab_on_deposit"`
	CryptoPayment          string                       `json:"crypto_payment"`
	HideInFooter           string                       `json:"hide_payment_in_footer"`
	OnlyInfoTextOnDeposit  string                       `json:"only_info_text_on_deposit"`
	OnlyInfoTextOnWithdraw string                       `json:"only_info_text_on_withdraw"`
	DisableConfirmPopup    string                       `json:"disable_confirm_pop_up"`
}

var ErrPaymentMethodNotFound = errors.New("payment method not found")
