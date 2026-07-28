package backoffice

import (
	"time"
)

type PlayerID int64

func (p PlayerID) Int64() int64 {
	return int64(p)
}

type PlayerCategory int

func (p PlayerCategory) Int() int {
	return int(p)
}

type PlayerCategoryEnum string

const (
	PlayerCategoryTestUser PlayerCategoryEnum = "testuser"
)

var playerCategoryToEnum = map[PlayerCategory]PlayerCategoryEnum{
	10: PlayerCategoryTestUser,
}

func (p PlayerCategory) Enum() PlayerCategoryEnum {
	return playerCategoryToEnum[p]
}

func (p PlayerCategory) Is(enum PlayerCategoryEnum) bool {
	return p.Enum() == enum
}

type ListPlayersPlayer struct {
	ID             PlayerID
	CreatedAt      time.Time
	Username       string
	FirstName      string
	MiddleName     string
	LastName       string
	Balance        float64
	PlayerCategory *PlayerCategory
	FirstDepositAt *time.Time
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
