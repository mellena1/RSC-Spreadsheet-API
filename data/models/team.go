package models

import "fmt"

// Team an RSC Team
// swagger:response models_Team
type Team struct {
	TeamID     string  `json:"id"`
	Name       string  `json:"name"`
	Franchise  string  `json:"franchise"`
	Tier       string  `json:"tier"`
	Conference string  `json:"conference"`
	Division   *string `json:"division,omitempty"`
}

func (t Team) String() string {
	return fmt.Sprintf("[%s] %s (%s)", t.Franchise, t.Name, t.Tier)
}

// Player holds data about a player
type Player struct {
	RSCID  string
	Name   string
	TeamID string
	Stats  Stats
}

// Stats holds stats about an entity
type Stats struct {
	// Record  Record
	Goals   int
	Assists int
	Saves   int
	Shots   int
}
