package models

import "fmt"

// Team describes a team
type Team struct {
	Name       string
	Franchise  string
	Tier       string
	Conference string
	Division   *string
	PlayerIDs  []string
	Stats      Stats
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
	Record  Record
	Goals   int
	Assists int
	Saves   int
	Shots   int
}
