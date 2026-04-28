package backoffice

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// ClientSearchRecord is one element of Data.Objects from POST /Client/GetClients (full row shape).
type ClientSearchRecord struct {
	ID                                  PlayerID        `json:"Id"`
	CurrencyID                          string          `json:"CurrencyId"`
	Currencies                          json.RawMessage `json:"Currencies"`
	FirstName                           string          `json:"FirstName"`
	LastName                            string          `json:"LastName"`
	MiddleName                          string          `json:"MiddleName"`
	Login                               string          `json:"Login"`
	RegionID                            int             `json:"RegionId"`
	Gender                              int             `json:"Gender"`
	PersonalID                          *string         `json:"PersonalId"`
	Address                             string          `json:"Address"`
	Email                               string          `json:"Email"`
	Language                            string          `json:"Language"`
	Phone                               string          `json:"Phone"`
	MobilePhone                         string          `json:"MobilePhone"`
	BirthDate                           string          `json:"BirthDate"`
	TimeZone                            *string         `json:"TimeZone"`
	NickName                            *string         `json:"NickName"`
	DocNumber                           string          `json:"DocNumber"`
	IBAN                                *string         `json:"IBAN"`
	PromoCode                           *string         `json:"PromoCode"`
	ProfileID                           *float64        `json:"ProfileId"`
	MaximalDailyBet                     *float64        `json:"MaximalDailyBet"`
	MaximalSingleBet                    *float64        `json:"MaximalSingleBet"`
	CasinoMaximalDailyBet               *float64        `json:"CasinoMaximalDailyBet"`
	CasinoMaximalSingleBet              *float64        `json:"CasinoMaximalSingleBet"`
	PreMatchSelectionLimit              *float64        `json:"PreMatchSelectionLimit"`
	LiveSelectionLimit                  *float64        `json:"LiveSelectionLimit"`
	Excluded                            *float64        `json:"Excluded"`
	ExcludedLocalDate                   *string         `json:"ExcludedLocalDate"`
	IsSubscribedToNewsletter            bool            `json:"IsSubscribedToNewsletter"`
	IsVerified                          bool            `json:"IsVerified"`
	PartnerName                         string          `json:"PartnerName"`
	PartnerID                           int             `json:"PartnerId"`
	LastLoginIP                         string          `json:"LastLoginIp"`
	RegistrationIP                      *string         `json:"RegistrationIp"`
	YesterdayBalance                    *float64        `json:"YesterdayBalance"`
	CreditLimit                         float64         `json:"CreditLimit"`
	IsUsingCredit                       bool            `json:"IsUsingCredit"`
	LastLoginTime                       string          `json:"LastLoginTime"`
	LastLoginLocalDate                  string          `json:"LastLoginLocalDate"`
	Balance                             float64         `json:"Balance"`
	IsLocked                            bool            `json:"IsLocked"`
	IsCasinoBlocked                     *bool           `json:"IsCasinoBlocked"`
	IsSportBlocked                      *bool           `json:"IsSportBlocked"`
	IsRMTBlocked                        *bool           `json:"IsRMTBlocked"`
	Password                            *string         `json:"Password"`
	PasswordChangeDate                  *string         `json:"PasswordChangeDate"`
	PasswordChangeDateLocal             *string         `json:"PasswordChangeDateLocal"`
	SportsbookProfileID                 int             `json:"SportsbookProfileId"`
	CasinoProfileID                     *float64        `json:"CasinoProfileId"`
	GlobalLiveDelay                     *float64        `json:"GlobalLiveDelay"`
	Created                             string          `json:"Created"`
	CreatedLocalDate                    string          `json:"CreatedLocalDate"`
	RFID                                *string         `json:"RFId"`
	ResetExpireDate                     *string         `json:"ResetExpireDate"`
	ResetExpireDateLocal                *string         `json:"ResetExpireDateLocal"`
	DocIssuedBy                         *string         `json:"DocIssuedBy"`
	LoyaltyLevelID                      *float64        `json:"LoyaltyLevelId"`
	IsUsingLoyaltyProgram               bool            `json:"IsUsingLoyaltyProgram"`
	LoyaltyPoint                        float64         `json:"LoyaltyPoint"`
	AffilateID                          *float64        `json:"AffilateId"`
	BTag                                *float64        `json:"BTag"`
	TermsAndConditionsVersion           string          `json:"TermsAndConditionsVersion"`
	TCVersionAcceptanceDate             string          `json:"TCVersionAcceptanceDate"`
	TCVersionAcceptanceLocalDate        string          `json:"TCVersionAcceptanceLocalDate"`
	ExcludedLast                        *string         `json:"ExcludedLast"`
	ExcludedLastLocal                   *string         `json:"ExcludedLastLocal"`
	UnplayedBalance                     float64         `json:"UnplayedBalance"`
	IsTest                              bool            `json:"IsTest"`
	ExternalID                          *string         `json:"ExternalId"`
	AuthomaticWithdrawalAmount          *float64        `json:"AuthomaticWithdrawalAmount"`
	AuthomaticWithdrawalMinLeftAmount   *float64        `json:"AuthomaticWithdrawalMinLeftAmount"`
	IsAutomaticWithdrawalEnabled        *bool           `json:"IsAutomaticWithdrawalEnabled"`
	SwiftCode                           *string         `json:"SwiftCode"`
	Title                               *string         `json:"Title"`
	BirthCity                           *string         `json:"BirthCity"`
	BirthDepartment                     *string         `json:"BirthDepartment"`
	BirthRegionID                       *float64        `json:"BirthRegionId"`
	ZipCode                             *string         `json:"ZipCode"`
	BirthRegionCode2                    *string         `json:"BirthRegionCode2"`
	ActivationCode                      *string         `json:"ActivationCode"`
	ActivationCodeExpireDate            *string         `json:"ActivationCodeExpireDate"`
	ActivationCodeExpireDateLocal       *string         `json:"ActivationCodeExpireDateLocal"`
	LastSportBetTime                    string          `json:"LastSportBetTime"`
	LastSportBetTimeLocal               string          `json:"LastSportBetTimeLocal"`
	VerificationDate                    *string         `json:"VerificationDate"`
	VerificationDateLocal               *string         `json:"VerificationDateLocal"`
	LastCasinoBetTime                   string          `json:"LastCasinoBetTime"`
	LastCasinoBetTimeLocal              string          `json:"LastCasinoBetTimeLocal"`
	FirstDepositTime                    string          `json:"FirstDepositTime"`
	FirstDepositDateLocal               string          `json:"FirstDepositDateLocal"`
	LastDepositDateLocal                string          `json:"LastDepositDateLocal"`
	LastDepositTime                     string          `json:"LastDepositTime"`
	PasswordChangedLastLocal            *string         `json:"PasswordChangedLastLocal"`
	PasswordChangedLast                 *string         `json:"PasswordChangedLast"`
	ActivationState                     *float64        `json:"ActivationState"`
	ExcludeTypeID                       *float64        `json:"ExcludeTypeId"`
	DocIssueDate                        *string         `json:"DocIssueDate"`
	DocIssueCode                        *string         `json:"DocIssueCode"`
	Province                            *string         `json:"Province"`
	IsResident                          bool            `json:"IsResident"`
	RegistrationSource                  int             `json:"RegistrationSource"`
	IncomeSource                        json.RawMessage `json:"IncomeSource"`
	AccountHolder                       *string         `json:"AccountHolder"`
	CashDeskID                          *float64        `json:"CashDeskId"`
	ClientCashDeskName                  *string         `json:"ClientCashDeskName"`
	IsSubscribeToEmail                  bool            `json:"IsSubscribeToEmail"`
	IsSubscribeToSMS                    bool            `json:"IsSubscribeToSMS"`
	IsSubscribeToBonus                  bool            `json:"IsSubscribeToBonus"`
	IsSubscribeToInternalMessage        bool            `json:"IsSubscribeToInternalMessage"`
	IsSubscribeToPushNotification       bool            `json:"IsSubscribeToPushNotification"`
	IsSubscribeToDirectMail             bool            `json:"IsSubscribeToDirectMail"`
	IsSubscribeToPhoneCall              bool            `json:"IsSubscribeToPhoneCall"`
	IsSubscribeToCasinoNewsletter       bool            `json:"IsSubscribeToCasinoNewsletter"`
	IsSubscribeToCasinoEmail            bool            `json:"IsSubscribeToCasinoEmail"`
	IsSubscribeToCasinoSMS              bool            `json:"IsSubscribeToCasinoSMS"`
	IsSubscribeToCasinoBonus            bool            `json:"IsSubscribeToCasinoBonus"`
	IsSubscribeToCasinoInternalMessage  bool            `json:"IsSubscribeToCasinoInternalMessage"`
	IsSubscribeToCasinoPushNotification bool            `json:"IsSubscribeToCasinoPushNotification"`
	IsSubscribeToCasinoPhoneCall        bool            `json:"IsSubscribeToCasinoPhoneCall"`
	NotificationOptions                 int             `json:"NotificationOptions"`
	IsLoggedIn                          bool            `json:"IsLoggedIn"`
	City                                string          `json:"City"`
	CountryName                         string          `json:"CountryName"`
	ClientVerificationDate              *string         `json:"ClientVerificationDate"`
	BankName                            *string         `json:"BankName"`
	Status                              int             `json:"Status"`
	IsNoBonus                           bool            `json:"IsNoBonus"`
	IsTwoFactorAuthenticationEnabled    bool            `json:"IsTwoFactorAuthenticationEnabled"`
	IsQRCodeUsed                        bool            `json:"IsQRCodeUsed"`
	PartnerClientCategoryID             int             `json:"PartnerClientCategoryId"`
	WrongLoginBlockLocalTime            *string         `json:"WrongLoginBlockLocalTime"`
	WrongLoginAttempts                  int             `json:"WrongLoginAttempts"`
	LastWrongLoginTimeLocalDate         string          `json:"LastWrongLoginTimeLocalDate"`
	PepStatusID                         *float64        `json:"PepStatusId"`
	DocRegionID                         *float64        `json:"DocRegionId"`
	DocRegionName                       *string         `json:"DocRegionName"`
	DocType                             *string         `json:"DocType"`
	DocExpirationDate                   *string         `json:"DocExpirationDate"`
	AMLRisk                             *string         `json:"AMLRisk"`
	ExclusionReason                     *string         `json:"ExclusionReason"`
	Citizenship                         *string         `json:"Citizenship"`
	IsPhoneVerified                     bool            `json:"IsPhoneVerified"`
	IsMobilePhoneVerified               bool            `json:"IsMobilePhoneVerified"`
	IsEkengVerified                     bool            `json:"IsEkengVerified"`
	IsEmailVerified                     bool            `json:"IsEmailVerified"`
	OwnerID                             *float64        `json:"OwnerId"`
	ChildID                             *float64        `json:"ChildId"`
	BirthName                           *string         `json:"BirthName"`
	StatusActiveDate                    *string         `json:"StatusActiveDate"`
	StatusActiveDateLocalTime           *string         `json:"StatusActiveDateLocalTime"`
	PartnerFlag                         *string         `json:"PartnerFlag"`
	AdditionalAddress                   *string         `json:"AdditionalAddress"`
	PepStatuses                         json.RawMessage `json:"PepStatuses"`
}

// SearchClientsRequest is the full filter body for POST /Client/GetClients (advanced client search).
// Field names and spelling match the backoffice API (e.g. AffilateId, SkeepRows).
type SearchClientsRequest struct {
	ID                         string                  `json:"Id,omitempty"`
	FirstName                  string                  `json:"FirstName,omitempty"`
	MiddleName                 string                  `json:"MiddleName,omitempty"`
	LastName                   string                  `json:"LastName,omitempty"`
	PersonalID                 string                  `json:"PersonalId,omitempty"`
	Email                      string                  `json:"Email,omitempty"`
	Phone                      string                  `json:"Phone,omitempty"`
	MobilePhone                string                  `json:"MobilePhone,omitempty"`
	ZipCode                    *string                 `json:"ZipCode,omitempty"`
	AMLRisk                    string                  `json:"AMLRisk,omitempty"`
	AffilateID                 *float64                `json:"AffilateId,omitempty"`
	AffiliatePlayerType        *float64                `json:"AffiliatePlayerType,omitempty"`
	BTag                       *float64                `json:"BTag,omitempty"`
	BetShopGroupID             string                  `json:"BetShopGroupId,omitempty"`
	BirthDate                  *string                 `json:"BirthDate,omitempty"`
	CashDeskID                 *float64                `json:"CashDeskId,omitempty"`
	CasinoProfileID            *float64                `json:"CasinoProfileId,omitempty"`
	CasinoProfitnessFrom       *float64                `json:"CasinoProfitnessFrom,omitempty"`
	CasinoProfitnessTo         *float64                `json:"CasinoProfitnessTo,omitempty"`
	City                       string                  `json:"City,omitempty"`
	ClientCategory             *float64                `json:"ClientCategory,omitempty"`
	CurrencyID                 *float64                `json:"CurrencyId,omitempty"`
	DocumentNumber             string                  `json:"DocumentNumber,omitempty"`
	ExternalID                 string                  `json:"ExternalId,omitempty"`
	Gender                     *float64                `json:"Gender,omitempty"`
	IBAN                       *string                 `json:"IBAN,omitempty"`
	IsEmailSubscribed          *bool                   `json:"IsEmailSubscribed,omitempty"`
	IsLocked                   *bool                   `json:"IsLocked,omitempty"`
	IsOrderedDesc              bool                    `json:"IsOrderedDesc"`
	IsSMSSubscribed            *bool                   `json:"IsSMSSubscribed,omitempty"`
	IsSelfExcluded             *bool                   `json:"IsSelfExcluded,omitempty"`
	IsStartWithSearch          bool                    `json:"IsStartWithSearch"`
	IsTest                     *bool                   `json:"IsTest,omitempty"`
	IsVerified                 *bool                   `json:"IsVerified,omitempty"`
	Login                      string                  `json:"Login,omitempty"`
	MaxBalance                 *float64                `json:"MaxBalance,omitempty"`
	MaxCreatedLocal            *ListPlayersRequestDate `json:"MaxCreatedLocal,omitempty"`
	MaxCreatedLocalDisable     bool                    `json:"MaxCreatedLocalDisable"`
	MaxFirstDepositDateLocal   *ListPlayersRequestDate `json:"MaxFirstDepositDateLocal,omitempty"`
	MaxLastTimeLoginDateLocal  *ListPlayersRequestDate `json:"MaxLastTimeLoginDateLocal,omitempty"`
	MaxLastWrongLoginDateLocal *ListPlayersRequestDate `json:"MaxLastWrongLoginDateLocal,omitempty"`
	MaxLoyaltyPointBalance     *float64                `json:"MaxLoyaltyPointBalance,omitempty"`
	MaxRows                    int                     `json:"MaxRows"`
	MaxVerificationDateLocal   *ListPlayersRequestDate `json:"MaxVerificationDateLocal,omitempty"`
	MaxWrongLoginAttempts      *float64                `json:"MaxWrongLoginAttempts,omitempty"`
	MinBalance                 *float64                `json:"MinBalance,omitempty"`
	MinCreatedLocal            *ListPlayersRequestDate `json:"MinCreatedLocal,omitempty"`
	MinCreatedLocalDisable     bool                    `json:"MinCreatedLocalDisable"`
	MinFirstDepositDateLocal   *ListPlayersRequestDate `json:"MinFirstDepositDateLocal,omitempty"`
	MinLastTimeLoginDateLocal  *ListPlayersRequestDate `json:"MinLastTimeLoginDateLocal,omitempty"`
	MinLastWrongLoginDateLocal *ListPlayersRequestDate `json:"MinLastWrongLoginDateLocal,omitempty"`
	MinLoyaltyPointBalance     *float64                `json:"MinLoyaltyPointBalance,omitempty"`
	MinVerificationDateLocal   *ListPlayersRequestDate `json:"MinVerificationDateLocal,omitempty"`
	MinWrongLoginAttempts      *float64                `json:"MinWrongLoginAttempts,omitempty"`
	NickName                   string                  `json:"NickName,omitempty"`
	OrderedItem                int                     `json:"OrderedItem"`
	OwnerID                    *float64                `json:"OwnerId,omitempty"`
	PartnerClientCategoryID    *float64                `json:"PartnerClientCategoryId,omitempty"`
	RegionID                   *float64                `json:"RegionId,omitempty"`
	RegistrationSource         *float64                `json:"RegistrationSource,omitempty"`
	SelectedPepStatuses        string                  `json:"SelectedPepStatuses,omitempty"`
	SkeepRows                  int                     `json:"SkeepRows"`
	SportProfitnessFrom        *float64                `json:"SportProfitnessFrom,omitempty"`
	SportProfitnessTo          *float64                `json:"SportProfitnessTo,omitempty"`
	Status                     *float64                `json:"Status,omitempty"`
	Time                       string                  `json:"Time,omitempty"`
	TimeZone                   string                  `json:"TimeZone,omitempty"`
}

// SearchClientsResult is Data from the GetClients response envelope (Count + Objects).
type SearchClientsResult struct {
	Count   int
	Clients []ClientSearchRecord
}

type searchClientsData struct {
	Count   int                  `json:"Count"`
	Objects []ClientSearchRecord `json:"Objects"`
}

func (c *client) SearchClients(ctx context.Context, req SearchClientsRequest) (*SearchClientsResult, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	data, err := makeRequest[searchClientsData](
		ctx,
		http.MethodPost,
		"/Client/GetClients",
		bytes.NewReader(body),
		c,
	)
	if err != nil {
		return nil, err
	}
	return &SearchClientsResult{Count: data.Count, Clients: data.Objects}, nil
}
