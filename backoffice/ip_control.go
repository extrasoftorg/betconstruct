package backoffice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type IPConflictClient struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type IPConflictDetail struct {
	IP           string             `json:"ip"`
	SameLastName bool               `json:"same_last_name"`
	Clients      []IPConflictClient `json:"clients"`
}

type IPConflictResult struct {
	HasConflict bool               `json:"has_conflict"`
	Conflicts   []IPConflictDetail `json:"conflicts"`
}

type ipControlGetClientsRequest struct {
	MaxRows       int    `json:"MaxRows"`
	Login         string `json:"Login,omitempty"`
	LoginIP       string `json:"LoginIP,omitempty"`
	IsOrderedDesc bool   `json:"IsOrderedDesc"`
	OrderedItem   int    `json:"OrderedItem"`
}

type ipControlClient struct {
	ID         int64  `json:"Id"`
	Login      string `json:"Login"`
	FirstName  string `json:"FirstName"`
	MiddleName string `json:"MiddleName"`
	LastName   string `json:"LastName"`
}

type ipControlGetClientsResponse struct {
	Objects []ipControlClient `json:"Objects"`
}

type ipControlGetLoginsRequest struct {
	ClientId      int64  `json:"ClientId"`
	FromDateLocal string `json:"FromDateLocal"`
	ToDateLocal   string `json:"ToDateLocal"`
}

type ipControlLoginRecord struct {
	LoginIP string `json:"LoginIP"`
}

type ipControlGetLoginsResponse struct {
	Objects []ipControlLoginRecord `json:"Objects"`
}

func (c *client) CheckIPConflict(ctx context.Context, username string, days int) (*IPConflictResult, error) {
	if days <= 0 {
		days = 30
	}

	clientsReq := ipControlGetClientsRequest{
		MaxRows:       20,
		Login:         username,
		IsOrderedDesc: true,
		OrderedItem:   1,
	}

	body, err := json.Marshal(clientsReq)
	if err != nil {
		return nil, err
	}

	clientsResp, err := makeRequest[ipControlGetClientsResponse](
		ctx,
		http.MethodPost,
		"/Client/GetClients",
		bytes.NewReader(body),
		c,
	)
	if err != nil {
		return nil, err
	}

	if clientsResp == nil || len(clientsResp.Objects) == 0 {
		return nil, errors.New("client not found")
	}

	targetClient := clientsResp.Objects[0]

	loc := time.FixedZone("UTC+3", 3*60*60)
	now := time.Now().In(loc)
	fromDate := now.AddDate(0, 0, -days).Format("02-01-06") + " - 00:00:00"
	toDate := now.AddDate(0, 0, 1).Format("02-01-06") + " - 00:00:00"

	loginsReq := ipControlGetLoginsRequest{
		ClientId:      targetClient.ID,
		FromDateLocal: fromDate,
		ToDateLocal:   toDate,
	}

	bodyLogins, err := json.Marshal(loginsReq)
	if err != nil {
		return nil, err
	}

	loginsResp, err := makeRequest[ipControlGetLoginsResponse](
		ctx,
		http.MethodPost,
		"/Client/GetLogins",
		bytes.NewReader(bodyLogins),
		c,
	)
	if err != nil {
		return nil, err
	}

	if loginsResp == nil || len(loginsResp.Objects) == 0 {
		return &IPConflictResult{HasConflict: false, Conflicts: []IPConflictDetail{}}, nil
	}

	ipMap := make(map[string]bool)
	var ipList []string
	for _, login := range loginsResp.Objects {
		if !ipMap[login.LoginIP] && login.LoginIP != "" {
			ipMap[login.LoginIP] = true
			ipList = append(ipList, login.LoginIP)
		}
	}

	result := &IPConflictResult{
		HasConflict: false,
		Conflicts:   make([]IPConflictDetail, 0),
	}

	for _, ip := range ipList {
		ipClientsReq := ipControlGetClientsRequest{
			MaxRows:       20,
			LoginIP:       ip,
			IsOrderedDesc: true,
			OrderedItem:   1,
		}

		bodyIP, err := json.Marshal(ipClientsReq)
		if err != nil {
			continue
		}

		ipClientsResp, err := makeRequest[ipControlGetClientsResponse](
			ctx,
			http.MethodPost,
			"/Client/GetClients",
			bytes.NewReader(bodyIP),
			c,
		)
		if err != nil || ipClientsResp == nil {
			continue
		}

		var conflictingClients []IPConflictClient
		sameLastName := false

		for _, otherClient := range ipClientsResp.Objects {
			if otherClient.ID != targetClient.ID {
				conflictingClients = append(conflictingClients, IPConflictClient{
					ID:        otherClient.ID,
					Username:  otherClient.Login,
					FirstName: otherClient.FirstName,
					LastName:  otherClient.LastName,
				})
				if strings.EqualFold(otherClient.LastName, targetClient.LastName) {
					sameLastName = true
				}
			}
		}

		if len(conflictingClients) > 0 {
			result.HasConflict = true
			result.Conflicts = append(result.Conflicts, IPConflictDetail{
				IP:           ip,
				SameLastName: sameLastName,
				Clients:      conflictingClients,
			})
		}
	}

	return result, nil
}
