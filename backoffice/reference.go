package backoffice

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type ListPaymentMethodsStatus string

func (s ListPaymentMethodsStatus) MarshalJSON() ([]byte, error) {
	var v int
	switch s {
	case ListPaymentMethodsStatusActive:
		v = 1
	case ListPaymentMethodsStatusInactive:
		v = 0
	}
	return json.Marshal(v)
}

const (
	ListPaymentMethodsStatusActive   ListPaymentMethodsStatus = "active"
	ListPaymentMethodsStatusInactive ListPaymentMethodsStatus = "inactive"
)

type paymentAPIRequestBase struct {
	RequestType int8 `json:"PaymentRequestType"`
}

type listPaymentMethodsRequest struct {
	paymentAPIRequestBase
	Status *ListPaymentMethodsStatus `json:"status,omitempty"`
}

type ListPaymentMethodsRequest struct {
	Status ListPaymentMethodsStatus
}

func (c *client) ListPaymentMethods(ctx context.Context, req ListPaymentMethodsRequest) ([]*PaymentMethod, error) {
	papiReq := listPaymentMethodsRequest{
		paymentAPIRequestBase: paymentAPIRequestBase{
			RequestType: 1,
		},
	}
	if req.Status != "" {
		papiReq.Status = &req.Status
	}

	body, err := json.Marshal(papiReq)
	if err != nil {
		return nil, err
	}

	methods, err := makeRequest[[]*PaymentMethod](ctx, http.MethodPost, "/Reference/PaymentAPI", bytes.NewReader(body), makeRequestOptions{
		httpClient: c.httpClient,
		authToken:  c.authToken,
	})
	if err != nil {
		return nil, err
	}
	return *methods, nil
}

type findPaymentMethodByNameRequest struct {
	paymentAPIRequestBase
	Name string `json:"system_name"`
}

func (c *client) FindPaymentMethodByName(ctx context.Context, name string) (*PaymentMethod, error) {
	papiReq := findPaymentMethodByNameRequest{
		paymentAPIRequestBase: paymentAPIRequestBase{
			RequestType: 1,
		},
		Name: name,
	}

	body, err := json.Marshal(papiReq)
	if err != nil {
		return nil, err
	}

	methods, err := makeRequest[[]PaymentMethod](ctx, http.MethodPost, "/Reference/PaymentAPI", bytes.NewReader(body), makeRequestOptions{
		httpClient: c.httpClient,
		authToken:  c.authToken,
	})
	if err != nil {
		return nil, err
	}

	if len(*methods) == 0 {
		return nil, ErrPaymentMethodNotFound
	}

	method := (*methods)[0]

	return &method, nil
}

type updatePaymentMethodRequest struct {
	paymentAPIRequestBase
	PaymentMethod
}

func (c *client) UpdatePaymentMethod(ctx context.Context, method PaymentMethod) error {
	papiReq := updatePaymentMethodRequest{
		paymentAPIRequestBase: paymentAPIRequestBase{
			RequestType: 6,
		},
		PaymentMethod: method,
	}
	body, err := json.Marshal(papiReq)
	if err != nil {
		return err
	}

	_, err = makeRequest[any](ctx, http.MethodPost, "/Reference/PaymentAPI", bytes.NewReader(body), makeRequestOptions{
		httpClient: c.httpClient,
		authToken:  c.authToken,
	})
	return err
}

type listPartnerDomainsResponse []PartnerDomain

func (c *client) ListPartnerDomains(ctx context.Context, partnerID PartnerID) ([]PartnerDomain, error) {
	domains, err := makeRequest[listPartnerDomainsResponse](ctx, http.MethodGet, "/Reference/GetPartnerDomains", nil, makeRequestOptions{
		httpClient: c.httpClient,
		authToken:  c.authToken,
	})
	if err != nil {
		return nil, err
	}
	return *domains, nil
}

type setActiveDomainRequest struct {
	ID int32 `json:"Id"`
}

func (c *client) SetActiveDomain(ctx context.Context, domainID int32) error {
	req := setActiveDomainRequest{
		ID: domainID,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = makeRequest[any](ctx, http.MethodPost, "/Reference/SetActiveDomain", bytes.NewReader(body), makeRequestOptions{
		httpClient: c.httpClient,
		authToken:  c.authToken,
	})
	return err
}
