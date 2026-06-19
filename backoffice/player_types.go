package backoffice

import (
	"encoding/json"
)

type PlayerID int64

func (p PlayerID) Int64() int64 {
	return int64(p)
}

type PlayerCategory string

func (p PlayerCategory) String() string {
	return string(p)
}

func (p *PlayerCategory) UnmarshalJSON(b []byte) error {
	var id int
	if err := json.Unmarshal(b, &id); err != nil {
		return err
	}
	switch id {
	case 10:
		*p = PlayerCategoryTestUser
	default:
		*p = PlayerCategoryUnknown
	}
	return nil
}

const (
	PlayerCategoryTestUser PlayerCategory = "testuser"
	PlayerCategoryUnknown  PlayerCategory = "unknown"
)

type ListPlayersPlayer struct {
	ID             PlayerID        `json:"Id"`
	CreatedAt      DateTime        `json:"CreatedLocalDate"`
	Username       string          `json:"Login"`
	FirstName      string          `json:"FirstName"`
	MiddleName     string          `json:"MiddleName"`
	LastName       string          `json:"LastName"`
	Balance        float64         `json:"Balance"`
	PlayerCategory *PlayerCategory `json:"SportsbookProfileId"`
}

type Player struct {
	ID         PlayerID `json:"Id"`
	FirstName  string   `json:"FirstName"`
	MiddleName string   `json:"MiddleName"`
	LastName   string   `json:"LastName"`
	Email      string   `json:"Email"`
	Phone      string   `json:"Phone"`
	Username   string   `json:"Login"`
}
