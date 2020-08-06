package models

// Team describes a team
type Team struct {
	Name       string
	Franchise  string
	Tier       string
	Conference string
	Division   string
	Players    []Player
	Stats      Stats
}

// Player holds data about a player
type Player struct {
	RSCID string
	Name  string
	Team  Team
	Stats Stats
}

// Stats holds stats about an entity
type Stats struct {
	Record  Record
	Goals   int
	Assists int
	Saves   int
	Shots   int
}
