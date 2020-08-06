package models

// TeamStanding holds standing stats about a current Team
type TeamStanding struct {
	Team             Team
	OverallRecord    Record
	ConferenceRecord Record
	DivisionRecord   *Record
}

// Record holds Wins and Losses
type Record struct {
	Wins   int
	Losses int
}

// GamesPlayed calculates how many games this team has played
func (r Record) GamesPlayed() int {
	return r.Wins + r.Losses
}

// WinPercentage calculates the win percentage of the team
func (r Record) WinPercentage() float64 {
	gamesPlayed := r.GamesPlayed()
	if gamesPlayed == 0 {
		return 0.0
	}
	return float64(r.Wins) / float64(gamesPlayed)
}
