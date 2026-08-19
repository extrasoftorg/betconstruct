package backoffice

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type PromoCodeType string

const (
	PromoCodeTypeCash    PromoCodeType = "cash"
	PromoCodeTypeFreeBet PromoCodeType = "free_bet"
)

var promoCodeTypes = map[PromoCodeType]int{
	PromoCodeTypeCash:    2,
	PromoCodeTypeFreeBet: 3,
}

func (t PromoCodeType) ID() int {
	return promoCodeTypes[t]
}

type PromoCodeItem interface {
	isPromoCodeItem()
}

type PromoCodeItemCash struct {
	Amount   float64
	Currency string
}

func (i PromoCodeItemCash) isPromoCodeItem() {}

type PromoCodeItemBonus struct {
	BonusID int
	Amount  float64
}

func (i PromoCodeItemBonus) isPromoCodeItem() {}

type CreatePromoCodeInput struct {
	Code      string
	MaxUses   int
	Type      int
	StartDate time.Time
	EndDate   time.Time
	Items     []PromoCodeItem
}

type createPromoCodePayloadItem struct {
	Amount   string `json:"Amount"`
	Currency string `json:"CurrencyId"`
}

type createPromoCodePayload struct {
	Code      string                       `json:"Code"`
	StartDate string                       `json:"StartDateLocal"`
	EndDate   string                       `json:"EndDateLocal"`
	MaxUses   string                       `json:"MaxCount"`
	IsActive  bool                         `json:"IsActive"`
	Type      int                          `json:"TypeId"`
	Items     []createPromoCodePayloadItem `json:"PromoItems"`
	BonusID   *int                         `json:"BonusId"`
}

func (in CreatePromoCodeInput) toPayload() (*createPromoCodePayload, error) {
	timeLayout := "02-01-06 - 15:04:05"
	payload := &createPromoCodePayload{
		Code:      in.Code,
		StartDate: in.StartDate.Format(timeLayout),
		EndDate:   in.EndDate.Format(timeLayout),
		MaxUses:   strconv.Itoa(in.MaxUses),
		IsActive:  true,
		Type:      in.Type,
		Items:     make([]createPromoCodePayloadItem, len(in.Items)),
	}

	for i, item := range in.Items {
		switch it := item.(type) {
		case PromoCodeItemCash:
			payload.Items[i] = createPromoCodePayloadItem{
				Amount:   strconv.FormatFloat(it.Amount, 'f', -1, 64),
				Currency: it.Currency,
			}
		case PromoCodeItemBonus:
			payload.BonusID = &it.BonusID
			payload.Items[i] = createPromoCodePayloadItem{
				Amount: strconv.FormatFloat(it.Amount, 'f', -1, 64),
			}
		}
	}

	return payload, nil
}

func (c *client) CreatePromoCode(ctx context.Context, in CreatePromoCodeInput) error {
	payload, err := in.toPayload()
	if err != nil {
		return err
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = makeRequest[any](
		ctx,
		http.MethodPost,
		"/Reference/SavePromoCodeWithItemsAsync",
		body,
		c,
	)
	return err
}

type PromoCode struct {
	ID        int64
	Code      string
	CreatedAt time.Time
	EndDate   time.Time
	StartDate time.Time
	MaxUses   int
	UsedCount int
	Type      int
}

type ListPromoCodesInput struct {
	Code string
}

type listPromoCodesPayload struct {
	Code string `json:"Code"`
}

func (in ListPromoCodesInput) toPayload() (*listPromoCodesPayload, error) {
	return &listPromoCodesPayload{
		Code: in.Code,
	}, nil
}

func (c *client) ListPromoCodes(ctx context.Context, in ListPromoCodesInput) ([]PromoCode, error) {
	payload, err := in.toPayload()
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	type responsePromoCode struct {
		ID        int64  `json:"Id"`
		Code      string `json:"Code"`
		CreatedAt string `json:"Created"`
		EndDate   string `json:"EndDateLocal"`
		StartDate string `json:"StartDateLocal"`
		MaxCount  int    `json:"MaxCount"`
		UsedCount int    `json:"UsedCount"`
		Type      int    `json:"TypeId"`
	}

	type response struct {
		Count      int                 `json:"Count"`
		PromoCodes []responsePromoCode `json:"Objects"`
	}

	resp, err := makeRequest[response](
		ctx,
		http.MethodPost,
		"/Reference/GetPromoCodesPagingAsync",
		body,
		c,
	)
	if err != nil {
		return nil, err
	}

	promoCodes := make([]PromoCode, len(resp.PromoCodes))
	for i, promoCode := range resp.PromoCodes {
		pc := PromoCode{
			ID:        promoCode.ID,
			Code:      promoCode.Code,
			MaxUses:   promoCode.MaxCount,
			UsedCount: promoCode.UsedCount,
			Type:      promoCode.Type,
		}

		pc.CreatedAt, err = time.Parse("2006-01-02T15:04:05.99999Z07:00", promoCode.CreatedAt)
		if err != nil {
			return nil, err
		}
		pc.EndDate, err = time.Parse("02-01-06 - 15:04:05", promoCode.EndDate)
		if err != nil {
			return nil, err
		}
		pc.StartDate, err = time.Parse("02-01-06 - 15:04:05", promoCode.StartDate)
		if err != nil {
			return nil, err
		}

		promoCodes[i] = pc
	}

	return promoCodes, nil
}

type ListPromoCodeUsagesInput struct {
	PromoCodeID int64
	PromoCode   string
	StartDate   time.Time
	EndDate     time.Time
}

type PromoCodeUsage struct {
	Amount           float64
	PlayerID         PlayerID
	PlayerUsername   string
	PromoCode        string
	ID               int64
	PromoCodeID      int64
	PromoCodeType    int
	PromoCodeBonusID *int64
	CreatedAt        time.Time
}

type listPromoCodeUsagesPayload struct {
	PromoCodeID string `json:"PromoCodeId"`
	PromoCode   string `json:"Code"`
	StartDate   string `json:"StartTimeLocal"`
	EndDate     string `json:"EndTimeLocal"`
}

func (in ListPromoCodeUsagesInput) toPayload() (*listPromoCodeUsagesPayload, error) {
	var promoCodeID string
	if in.PromoCodeID != 0 {
		promoCodeID = strconv.FormatInt(in.PromoCodeID, 10)
	}
	return &listPromoCodeUsagesPayload{
		PromoCodeID: promoCodeID,
		PromoCode:   in.PromoCode,
		StartDate:   in.StartDate.Format("02-01-06 - 15:04:05"),
		EndDate:     in.EndDate.Format("02-01-06 - 15:04:05"),
	}, nil
}

func (c *client) ListPromoCodeUsages(ctx context.Context, in ListPromoCodeUsagesInput) ([]*PromoCodeUsage, error) {
	payload, err := in.toPayload()

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	type usage struct {
		Amount           float64  `json:"Amount"`
		PlayerID         PlayerID `json:"PlayerId"`
		PlayerUsername   string   `json:"Login"`
		PromoCode        string   `json:"Code"`
		ID               int64    `json:"Id"`
		PromoCodeID      int64    `json:"PromoCodeId"`
		PromoCodeType    int      `json:"PromoCodeType"`
		PromoCodeBonusID *int64   `json:"PromoCodeBonusId"`
		CreatedAt        string   `json:"CreatedDateLocal"`
	}

	resp, err := makeRequest[[]*usage](
		ctx,
		http.MethodPost,
		"/Report/GetClientPromoCodes",
		body,
		c,
	)
	if err != nil {
		return nil, err
	}

	usages := make([]*PromoCodeUsage, len(*resp))
	for i, usage := range *resp {
		usg := &PromoCodeUsage{
			Amount:           usage.Amount,
			PlayerID:         usage.PlayerID,
			PlayerUsername:   usage.PlayerUsername,
			PromoCode:        usage.PromoCode,
			ID:               usage.ID,
			PromoCodeID:      usage.PromoCodeID,
			PromoCodeType:    usage.PromoCodeType,
			PromoCodeBonusID: usage.PromoCodeBonusID,
		}
		usg.CreatedAt, err = time.Parse("2006-01-02T15:04:05.99999", usage.CreatedAt)
		if err != nil {
			return nil, err
		}
		usages[i] = usg
	}

	return usages, nil
}
